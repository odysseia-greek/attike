package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/attike/aristophanes/config"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"log"
	"sync"
	"time"
)

const (
	TRACE        string = "trace"
	SPAN         string = "span"
	DATABASESPAN string = "database_span"
	ITEMS        string = "items"
)

// AristophanesHandler is the gRPC server handling task queue operations
type AristophanesHandler struct {
	pb.UnimplementedTraceServiceServer
	mu     sync.Mutex // Mutex to protect the task queue
	Config *config.Config
}

func (a *AristophanesHandler) StartTrace(ctx context.Context, start *pb.StartTraceRequest) (*pb.TraceResponse, error) {
	traceId := uuid.New().String()
	spanId := a.generateSpanID()

	if start.SaveTrace {
		trace := pb.Trace{
			SpanId:       spanId,
			ParentSpanId: spanId,
			Method:       start.Method,
			Url:          start.Url,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
			PodName:      a.Config.PodName,
			Namespace:    a.Config.Namespace,
			ItemType:     TRACE,
		}

		jsonData := map[string][]pb.Trace{
			"items": {trace},
		}

		data, err := json.Marshal(&jsonData)
		if err != nil {
			return nil, err
		}

		doc, err := a.Config.Elastic.Document().CreateWithId(a.Config.Index, traceId, data)
		if err != nil {
			return nil, err
		}

		log.Printf("created trace with id: %s", doc.ID)
	}

	combinedId := fmt.Sprintf("%s+%s", traceId, spanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (a *AristophanesHandler) StartNewSpan(ctx context.Context, start *pb.StartSpanRequest) (*pb.TraceResponse, error) {
	spanId := a.generateSpanID()
	combinedId := fmt.Sprintf("%s+%s", start.TraceId, spanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (a *AristophanesHandler) Trace(ctx context.Context, traceRequest *pb.TraceRequest) (*pb.TraceResponse, error) {
	spanId := a.generateSpanID()

	if traceRequest.SaveTrace {
		trace := pb.Trace{
			SpanId:       spanId,
			ParentSpanId: traceRequest.ParentSpanId,
			Method:       traceRequest.Method,
			Url:          traceRequest.Url,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
			PodName:      a.Config.PodName,
			Namespace:    a.Config.Namespace,
			ItemType:     TRACE,
		}
		data, err := json.Marshal(&trace)
		if err != nil {
			return nil, err
		}

		initialItemArray := fmt.Sprintf(`[%s]`, data)

		doc, err := a.Config.Elastic.Document().AddItemToDocument(a.Config.Index, traceRequest.TraceId, initialItemArray, ITEMS)
		if err != nil {
			return nil, err
		}

		log.Printf("created trace with id: %s", doc.ID)
	}

	combinedId := fmt.Sprintf("%s+%s", traceRequest.TraceId, traceRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (a *AristophanesHandler) Span(ctx context.Context, spanRequest *pb.SpanRequest) (*pb.TraceResponse, error) {
	spanId := a.generateSpanID()

	if spanRequest.SaveTrace {
		span := pb.Span{
			SpanId:       spanId,
			ParentSpanId: spanRequest.ParentSpanId,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
			PodName:      a.Config.PodName,
			Namespace:    a.Config.Namespace,
			Action:       spanRequest.Action,
			RequestBody:  spanRequest.RequestBody,
			ResponseBody: spanRequest.ResponseBody,
			ItemType:     SPAN,
		}

		data, err := json.Marshal(&span)
		if err != nil {
			return nil, err
		}

		initialItemArray := fmt.Sprintf(`[%s]`, data)

		doc, err := a.Config.Elastic.Document().AddItemToDocument(a.Config.Index, spanRequest.TraceId, initialItemArray, ITEMS)
		if err != nil {
			return nil, err
		}

		log.Printf("created trace with id: %s", doc.ID)
	}

	combinedId := fmt.Sprintf("%s+%s", spanRequest.TraceId, spanRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (a *AristophanesHandler) DatabaseSpan(ctx context.Context, spanRequest *pb.DatabaseSpanRequest) (*pb.TraceResponse, error) {
	spanId := a.generateSpanID()

	if spanRequest.SaveTrace {
		span := pb.DatabaseSpan{
			SpanId:       spanId,
			ParentSpanId: spanRequest.ParentSpanId,
			Timestamp:    time.Now().UTC().Format("2006-01-02'T'15:04:05"),
			PodName:      a.Config.PodName,
			Namespace:    a.Config.Namespace,
			Query:        spanRequest.Query,
			ResultJson:   spanRequest.ResultJson,
			ItemType:     DATABASESPAN,
		}

		data, err := json.Marshal(&span)
		if err != nil {
			return nil, err
		}

		initialItemArray := fmt.Sprintf(`[%s]`, data)

		doc, err := a.Config.Elastic.Document().AddItemToDocument(a.Config.Index, spanRequest.TraceId, initialItemArray, ITEMS)
		if err != nil {
			return nil, err
		}

		log.Printf("created trace with id: %s", doc.ID)
	}

	combinedId := fmt.Sprintf("%s+%s", spanRequest.TraceId, spanRequest.ParentSpanId)
	return &pb.TraceResponse{
		CombinedId: combinedId,
	}, nil
}

func (a *AristophanesHandler) generateSpanID() string {
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
