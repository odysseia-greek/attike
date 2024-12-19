package comedy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	pbm "github.com/odysseia-greek/attike/sophokles/proto"
	"io"
	"time"
)

const (
	TRACE        string = "trace"
	SPAN         string = "span"
	TRACESTART   string = "trace_start"
	TRACECLOSE   string = "trace_close"
	DATABASESPAN string = "database_span"
	ITEMS        string = "items"
)

func (t *TraceServiceImpl) ManageStartTimeMap() {
	startTimeMap := make(map[string]time.Time)
	for cmd := range t.commands {
		switch cmd.Action {
		case "set":
			startTimeMap[cmd.TraceID] = cmd.Time
		case "get":
			startTime, found := startTimeMap[cmd.TraceID]
			cmd.Response <- MapResponse{Time: startTime, Found: found}
		case "delete":
			delete(startTimeMap, cmd.TraceID)
		}
	}
}

func (t *TraceServiceImpl) HealthCheck(ctx context.Context, start *pb.Empty) (*pb.HealthCheckResponse, error) {
	elasticHealth := t.Elastic.Health().Info()
	return &pb.HealthCheckResponse{Status: elasticHealth.Healthy}, nil
}

func (t *TraceServiceImpl) Chorus(stream pb.TraceService_ChorusServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.TraceResponse{Ack: "Received"})
		}
		if err != nil {
			return err
		}

		switch req := in.RequestType.(type) {
		case *pb.ParabasisRequest_StartTrace:
			go t.StartTrace(req, in.TraceId, in.ParentSpanId)
		case *pb.ParabasisRequest_Trace:
			go t.Trace(req, in.TraceId, in.ParentSpanId, in.SpanId)
		case *pb.ParabasisRequest_CloseTrace:
			go t.CloseTrace(req, in.TraceId, in.ParentSpanId)
		case *pb.ParabasisRequest_Span:
			go t.Span(req, in.TraceId, in.ParentSpanId)
		case *pb.ParabasisRequest_DatabaseSpan:
			go t.DatabaseSpan(req, in.TraceId, in.ParentSpanId)
		default:
			logging.Debug(fmt.Sprintf("Unhandled trace request type: %s", req))
		}
	}
}

func (t *TraceServiceImpl) StartTrace(start *pb.ParabasisRequest_StartTrace, traceID, ParentSpanID string) {
	traceTime := time.Now().UTC()
	t.commands <- MapCommand{
		Action:  "set",
		TraceID: traceID,
		Time:    traceTime,
	}

	traceData := &pb.TraceStart{
		Method:        start.StartTrace.Method,
		Url:           start.StartTrace.Url,
		Host:          start.StartTrace.Host,
		RemoteAddress: start.StartTrace.RemoteAddress,
		Operation:     start.StartTrace.Operation,
		RootQuery:     start.StartTrace.RootQuery,
		Common: &pb.TraceCommon{
			SpanId:       ParentSpanID,
			ParentSpanId: ParentSpanID,
			Timestamp:    traceTime.Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     TRACESTART,
		},
	}

	if t.GatherMetrics && t.Metrics != nil {
		metrics, err := t.gatherMetrics()
		if err != nil {
			logging.Error(fmt.Sprintf("cannot set metrics because an error was returned: %s", err.Error()))
		} else {
			traceData.Metrics = metrics
		}
	}

	jsonData := map[string]interface{}{
		"items":       []pb.TraceStart{*traceData},
		"isActive":    true,
		"timeStarted": traceTime.Format("2006-01-02T15:04:05.000"), // Include milliseconds
		"timeEnded":   "1970-01-01T00:00:00.000",                   // Default null-like value
		"totalTime":   0,
	}

	data, err := json.Marshal(&jsonData)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))
		return
	}

	doc, err := t.Elastic.Document().CreateWithId(t.Index, traceID, data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error creating document for trace ID %s: %s", traceID, err))
		return
	}

	logging.Debug(fmt.Sprintf("created trace with id: %s", doc.ID))

}

func (t *TraceServiceImpl) CloseTrace(close *pb.ParabasisRequest_CloseTrace, traceID, ParentSpanID string) {
	traceTime := time.Now().UTC()

	traceData := pb.TraceStop{
		ResponseBody: close.CloseTrace.ResponseBody,

		Common: &pb.TraceCommon{
			ParentSpanId: ParentSpanID,
			Timestamp:    traceTime.Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     TRACECLOSE,
		},
	}

	if t.GatherMetrics && t.Metrics != nil {
		metrics, err := t.gatherMetrics()
		if err != nil {
			logging.Error(fmt.Sprintf("cannot set metrics because an error was returned: %s", err.Error()))
		} else {
			traceData.Metrics = metrics
		}
	}

	data, err := json.Marshal(&traceData)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))

	}
	docID, err := t.UpdateDocumentWithRetry(traceID, ITEMS, data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error updating document for trace ID %s: %s", traceID, err))
	}

	var totalTime int64
	responseChan := make(chan MapResponse)
	t.commands <- MapCommand{
		Action:   "get",
		TraceID:  traceID,
		Response: responseChan,
	}
	response := <-responseChan

	if !response.Found {
		totalTime = 0
	} else {
		totalTime = traceTime.Sub(response.Time).Milliseconds()
	}

	closingTrace := map[string]interface{}{
		"isActive":     false,
		"timeEnded":    traceTime.Format("2006-01-02T15:04:05.000"), // Include milliseconds
		"totalTime":    totalTime,
		"responseCode": close.CloseTrace.ResponseCode,
	}

	update, err := json.Marshal(closingTrace)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))
	}

	doc, err := t.Elastic.Document().Update(t.Index, traceID, update)
	if err != nil {
		logging.Error(fmt.Sprintf("Error updating document for trace ID %s: %s", traceID, err))
		return
	}

	t.commands <- MapCommand{
		Action:  "delete",
		TraceID: doc.ID,
	}

	logging.Debug(fmt.Sprintf("closed trace with id: %s", docID))
}

func (t *TraceServiceImpl) Trace(in *pb.ParabasisRequest_Trace, traceID, ParentSpanID, spanID string) {
	traceData := &pb.Trace{
		Method: in.Trace.Method,
		Url:    in.Trace.Url,
		Host:   in.Trace.Host,
		Common: &pb.TraceCommon{
			SpanId:       spanID,
			ParentSpanId: ParentSpanID,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     TRACE,
		},
	}

	if t.GatherMetrics && t.Metrics != nil {
		metrics, err := t.gatherMetrics()
		if err != nil {
			logging.Error(fmt.Sprintf("cannot set metrics because an error was returned: %s", err.Error()))
		} else {
			traceData.Metrics = metrics
		}
	}
	data, err := json.Marshal(&traceData)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))

	}
	docID, err := t.UpdateDocumentWithRetry(traceID, ITEMS, data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error updating document for trace ID %s: %s", traceID, err))
		return
	}

	logging.Debug(fmt.Sprintf("added trace with id: %s", docID))
}

func (t *TraceServiceImpl) Span(in *pb.ParabasisRequest_Span, traceID, ParentSpanID string) {
	spanID := GenerateSpanID()

	span := &pb.Span{
		Action: in.Span.Action,
		Status: in.Span.Status,
		Took:   in.Span.Took,
		Common: &pb.TraceCommon{
			ParentSpanId: ParentSpanID,
			SpanId:       spanID,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     SPAN,
		},
	}

	data, err := json.Marshal(&span)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))

	}
	docID, err := t.UpdateDocumentWithRetry(traceID, ITEMS, data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error updating document for trace ID %s: %s", traceID, err))
		return
	}

	logging.Debug(fmt.Sprintf("added span with id: %s to trace: %s", spanID, docID))
}

func (t *TraceServiceImpl) DatabaseSpan(in *pb.ParabasisRequest_DatabaseSpan, traceID, ParentSpanID string) {
	spanId := GenerateSpanID()

	timeEndedStr := time.Now().UTC().Format("2006-01-02'T'15:04:05.000")
	timeEnded, err := time.Parse("2006-01-02'T'15:04:05.000", timeEndedStr)
	if err != nil {
		logging.Error(fmt.Sprintf("Error parsing time for trace ID %s: %s", traceID, err))
	}

	// Convert TimeTook to a time.Duration based on MS
	duration := time.Duration(in.DatabaseSpan.TimeTook) * time.Millisecond
	timeStartedStr := timeEnded.Add(-duration).Format("2006-01-02'T'15:04:05.000")

	dbSpan := &pb.DatabaseSpan{
		Query:        in.DatabaseSpan.Query,
		TimeStarted:  timeStartedStr,
		TimeFinished: timeEndedStr,
		Hits:         in.DatabaseSpan.Hits,
		Took:         fmt.Sprintf("%vms", in.DatabaseSpan.TimeTook),
		Common: &pb.TraceCommon{
			SpanId:       spanId,
			ParentSpanId: ParentSpanID,
			Timestamp:    timeEndedStr,
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     DATABASESPAN,
		},
	}

	data, err := json.Marshal(dbSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling document for trace ID %s: %s", traceID, err))
	}
	docID, err := t.UpdateDocumentWithRetry(traceID, ITEMS, data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error updating document for trace ID %s: %s", traceID, err))
		return
	}

	logging.Debug(fmt.Sprintf("added database span to id: %s", docID))
}

func (t *TraceServiceImpl) gatherMetrics() (*pb.TracingMetrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metrics, err := t.Metrics.FetchMetrics(ctx, &pbm.Empty{})
	if err != nil {
		return nil, err
	}

	return &pb.TracingMetrics{
		CpuUnits:            metrics.CpuUnits,
		MemoryUnits:         metrics.MemoryUnits,
		Name:                metrics.Pod.Name,
		CpuRaw:              metrics.Pod.CpuRaw,
		MemoryRaw:           metrics.Pod.MemoryRaw,
		CpuHumanReadable:    metrics.Pod.CpuHumanReadable,
		MemoryHumanReadable: metrics.Pod.MemoryHumanReadable,
	}, nil
}

func (t *TraceServiceImpl) UpdateDocumentWithRetry(traceID, itemName string, data []byte) (string, error) {
	maxRetries := 10
	retryDelay := 100 * time.Millisecond
	var tenTriesError error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		doc, err := t.Elastic.Document().AddItemToDocument(t.Index, traceID, string(data), itemName)

		if err == nil {
			return doc.ID, nil
		}

		if attempt < maxRetries {
			tenTriesError = err
			// Sleep before the next retry
			time.Sleep(retryDelay)
		}
	}

	return "", fmt.Errorf("error updating document for trace ID %s: %s", traceID, tenTriesError.Error())
}
