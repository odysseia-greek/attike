package comedy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"log"
	"time"
)

const (
	TRACE        string = "trace"
	SPAN         string = "span"
	DATABASESPAN string = "database_span"
	ITEMS        string = "items"
)

func (t *TraceServiceImpl) HealthCheck(ctx context.Context, start *pb.Empty) (*pb.HealthCheckResponse, error) {
	elasticHealth := t.Elastic.Health().Info()
	return &pb.HealthCheckResponse{Status: elasticHealth.Healthy}, nil
}

func (t *TraceServiceImpl) StartTrace(ctx context.Context, start *pb.StartTraceRequest) (*pb.TraceResponse, error) {
	traceId := uuid.New().String()
	spanId := t.generateSpanID()

	traceTime := time.Now().UTC()
	t.StartTimeMap[traceId] = traceTime

	trace := pb.TraceStart{
		Method:        start.Method,
		Url:           start.Url,
		Host:          start.Host,
		RemoteAddress: start.RemoteAddress,
		Operation:     start.Operation,
		RootQuery:     start.RootQuery,
		Common: &pb.TraceCommon{
			SpanId:       spanId,
			ParentSpanId: spanId,
			Timestamp:    traceTime.Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     TRACE,
		},
	}

	jsonData := map[string]interface{}{
		"items":       []pb.TraceStart{trace},
		"isActive":    true,
		"timeStarted": traceTime.Format("2006-01-02T15:04:05.000"), // Include milliseconds
		"timeEnded":   "1970-01-01T00:00:00.000",                   // Default null-like value
		"totalTime":   0,
	}

	data, err := json.Marshal(&jsonData)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().CreateWithId(t.Index, traceId, data)
	if err != nil {
		return nil, fmt.Errorf("an error was returned by elasticSearch: %v", err)
	}

	log.Printf("created trace with id: %s", doc.ID)

	combinedId := fmt.Sprintf("%s+%s+%d", traceId, spanId, 1)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) CloseTrace(ctx context.Context, stop *pb.CloseTraceRequest) (*pb.TraceResponse, error) {
	traceTime := time.Now().UTC()

	trace := pb.TraceStop{
		Common: &pb.TraceCommon{
			ParentSpanId: stop.ParentSpanId,
			Timestamp:    traceTime.Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
		},
	}

	if stop.ResponseBody != "" {
		trace.ResponseBody = stop.ResponseBody
	}

	data, err := json.Marshal(&trace)
	if err != nil {
		return nil, err
	}

	_, err = t.UpdateDocumentWithRetry(stop.TraceId, ITEMS, data)
	if err != nil {
		return nil, err
	}

	var totalTime int64
	originalTimeStart, found := t.StartTimeMap[stop.TraceId]
	if !found {
		totalTime = 0
	} else {
		totalTime = traceTime.Sub(originalTimeStart).Milliseconds()
	}

	closingTrace := map[string]interface{}{
		"isActive":     false,
		"timeEnded":    traceTime.Format("2006-01-02T15:04:05.000"), // Include milliseconds
		"totalTime":    totalTime,
		"responseCode": stop.ResponseCode,
	}

	update, err := json.Marshal(closingTrace)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().Update(t.Index, stop.TraceId, update)
	if err != nil {
		return nil, fmt.Errorf("an error was returned by elasticSearch: %v", err)
	}

	// Remove the entry from the map
	delete(t.StartTimeMap, doc.ID)

	log.Printf("closed trace with id: %s", doc.ID)

	return &pb.TraceResponse{
		CombinedId: stop.TraceId,
	}, nil
}

func (t *TraceServiceImpl) StartNewSpan(ctx context.Context, start *pb.StartSpanRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()
	combinedId := fmt.Sprintf("%s+%s", start.TraceId, spanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) Trace(ctx context.Context, traceRequest *pb.TraceRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()

	trace := pb.Trace{
		Method: traceRequest.Method,
		Url:    traceRequest.Url,
		Host:   traceRequest.Host,
		Common: &pb.TraceCommon{
			SpanId:       spanId,
			ParentSpanId: traceRequest.ParentSpanId,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     TRACE,
		},
	}
	data, err := json.Marshal(&trace)
	if err != nil {
		return nil, err
	}

	// if this update fails the id might not yet exists. Perhaps its best to verify and if it doesnt exist start a trace here
	// else it will error out leading to difficult debug situations
	docID, err := t.UpdateDocumentWithRetry(traceRequest.TraceId, ITEMS, data)
	if err != nil {
		return nil, fmt.Errorf("an error was returned by elasticSearch: %v", err)
	}

	log.Printf("added trace with id: %s", docID)

	combinedId := fmt.Sprintf("%s+%s", traceRequest.TraceId, traceRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) Span(ctx context.Context, spanRequest *pb.SpanRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()

	span := pb.Span{
		Action:       spanRequest.Action,
		RequestBody:  spanRequest.RequestBody,
		ResponseBody: spanRequest.ResponseBody,
		Common: &pb.TraceCommon{
			SpanId:       spanId,
			ParentSpanId: spanRequest.ParentSpanId,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     SPAN,
		},
	}

	data, err := json.Marshal(&span)
	if err != nil {
		return nil, err
	}

	docID, err := t.UpdateDocumentWithRetry(spanRequest.TraceId, ITEMS, data)
	if err != nil {
		return nil, fmt.Errorf("an error was returned by elasticSearch: %v", err)
	}

	log.Printf("added span with id: %s to trace: %s", spanId, docID)

	combinedId := fmt.Sprintf("%s+%s", spanRequest.TraceId, spanRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) DatabaseSpan(ctx context.Context, spanRequest *pb.DatabaseSpanRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()

	span := pb.DatabaseSpan{
		Query:      spanRequest.Query,
		ResultJson: spanRequest.ResultJson,
		Common: &pb.TraceCommon{
			SpanId:       spanId,
			ParentSpanId: spanRequest.ParentSpanId,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     DATABASESPAN,
		},
	}

	data, err := json.Marshal(&span)
	if err != nil {
		return nil, err
	}

	docID, err := t.UpdateDocumentWithRetry(spanRequest.TraceId, ITEMS, data)
	if err != nil {
		return nil, fmt.Errorf("an error was returned by elasticSearch: %v", err)
	}

	log.Printf("added database_span with id: %s to trace: %s", spanId, docID)

	combinedId := fmt.Sprintf("%s+%s", spanRequest.TraceId, spanRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) generateSpanID() string {
	// Create a byte slice to store the random data
	randomBytes := make([]byte, 8)

	// Read random data from the crypto/rand package
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return ""
	}

	// Convert the random data to a hexadecimal string
	spanID := hex.EncodeToString(randomBytes)

	return spanID
}

func (t *TraceServiceImpl) UpdateDocumentWithRetry(traceID, itemName string, data []byte) (string, error) {
	maxRetries := 10
	retryDelay := 50 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		doc, err := t.Elastic.Document().AddItemToDocument(t.Index, traceID, string(data), itemName)
		if err == nil {
			return doc.ID, nil
		}

		if attempt < maxRetries {
			// Sleep before the next retry
			time.Sleep(retryDelay)
		}
	}

	return "", fmt.Errorf("failed after %d attempts", maxRetries)
}
