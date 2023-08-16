package app

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

func (t *TraceServiceImpl) StartTrace(ctx context.Context, start *pb.StartTraceRequest) (*pb.TraceResponse, error) {
	traceId := uuid.New().String()
	spanId := t.generateSpanID()

	trace := pb.TraceStart{
		ParentSpanId:  spanId,
		Method:        start.Method,
		Url:           start.Url,
		Host:          start.Host,
		RemoteAddress: start.RemoteAddress,
		Timestamp:     time.Now().UTC().Format("2006-01-02'T'15:04:05"),
		PodName:       t.PodName,
		Namespace:     t.Namespace,
		ItemType:      TRACE,
		Operation:     start.Operation,
		RootQuery:     start.RootQuery,
	}

	jsonData := map[string]interface{}{
		"items":    []pb.TraceStart{trace},
		"isActive": true,
	}

	data, err := json.Marshal(&jsonData)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().CreateWithId(t.Index, traceId, data)
	if err != nil {
		return nil, err
	}

	log.Printf("created trace with id: %s", doc.ID)

	combinedId := fmt.Sprintf("%s+%s+%d", traceId, spanId, 1)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) CloseTrace(ctx context.Context, stop *pb.CloseTraceRequest) (*pb.TraceResponse, error) {
	trace := pb.TraceStop{
		ParentSpanId: stop.ParentSpanId,
		Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
		PodName:      t.PodName,
		Namespace:    t.Namespace,
		ResponseBody: stop.ResponseBody,
	}

	data, err := json.Marshal(&trace)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().AddItemToDocument(t.Index, stop.TraceId, string(data), ITEMS)
	if err != nil {
		return nil, err
	}

	closingTrace := map[string]interface{}{
		"isActive": false,
	}

	update, err := json.Marshal(closingTrace)
	if err != nil {
		return nil, err
	}

	updatedDoc, err := t.Elastic.Document().Update(t.Index, stop.TraceId, update)
	if err != nil {
		return nil, err
	}

	log.Printf("current version: %v", updatedDoc.Version)
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
		SpanId:       spanId,
		ParentSpanId: traceRequest.ParentSpanId,
		Method:       traceRequest.Method,
		Url:          traceRequest.Url,
		Host:         traceRequest.Host,
		Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
		PodName:      t.PodName,
		Namespace:    t.Namespace,
		ItemType:     TRACE,
	}
	data, err := json.Marshal(&trace)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().AddItemToDocument(t.Index, traceRequest.TraceId, string(data), ITEMS)
	if err != nil {
		return nil, err
	}

	log.Printf("added trace with id: %s", doc.ID)

	combinedId := fmt.Sprintf("%s+%s", traceRequest.TraceId, traceRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) Span(ctx context.Context, spanRequest *pb.SpanRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()

	span := pb.Span{
		SpanId:       spanId,
		ParentSpanId: spanRequest.ParentSpanId,
		Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
		PodName:      t.PodName,
		Namespace:    t.Namespace,
		Action:       spanRequest.Action,
		RequestBody:  spanRequest.RequestBody,
		ResponseBody: spanRequest.ResponseBody,
		ItemType:     SPAN,
	}

	data, err := json.Marshal(&span)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().AddItemToDocument(t.Index, spanRequest.TraceId, string(data), ITEMS)
	if err != nil {
		return nil, err
	}

	log.Printf("added span with id: %s to trace: %s", spanId, doc.ID)

	combinedId := fmt.Sprintf("%s+%s", spanRequest.TraceId, spanRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (t *TraceServiceImpl) DatabaseSpan(ctx context.Context, spanRequest *pb.DatabaseSpanRequest) (*pb.TraceResponse, error) {
	spanId := t.generateSpanID()

	span := pb.DatabaseSpan{
		SpanId:       spanId,
		ParentSpanId: spanRequest.ParentSpanId,
		Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
		PodName:      t.PodName,
		Namespace:    t.Namespace,
		Query:        spanRequest.Query,
		ResultJson:   spanRequest.ResultJson,
		ItemType:     DATABASESPAN,
	}

	data, err := json.Marshal(&span)
	if err != nil {
		return nil, err
	}

	doc, err := t.Elastic.Document().AddItemToDocument(t.Index, spanRequest.TraceId, string(data), ITEMS)
	if err != nil {
		return nil, err
	}

	log.Printf("added database_span with id: %s to trace: %s", spanId, doc.ID)

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
