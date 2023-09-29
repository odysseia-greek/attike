package app

import (
	"context"
	"github.com/odysseia-greek/aristoteles"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type TraceService interface {
	HealthCheck(ctx context.Context, start *pb.Empty) (*pb.HealthCheckResponse, error)
	StartTrace(ctx context.Context, request *pb.StartTraceRequest) (*pb.TraceResponse, error)
	Trace(ctx context.Context, request *pb.TraceRequest) (*pb.TraceResponse, error)
	StartNewSpan(ctx context.Context, request *pb.StartSpanRequest) (*pb.TraceResponse, error)
	Span(ctx context.Context, request *pb.SpanRequest) (*pb.TraceResponse, error)
	DatabaseSpan(ctx context.Context, request *pb.DatabaseSpanRequest) (*pb.TraceResponse, error)
	CloseTrace(ctx context.Context, request *pb.CloseTraceRequest) (*pb.TraceResponse, error)
	WaitForHealthyState() bool
}

type TraceServiceImpl struct {
	PodName      string
	Namespace    string
	Index        string
	Elastic      aristoteles.Client
	StartTimeMap map[string]time.Time
	pb.UnimplementedTraceServiceServer
	mu sync.Mutex // Mutex to protect the task queue
}

type TraceServiceClient struct {
	Impl TraceService
}

type ClientTracer struct {
	tracer pb.TraceServiceClient
}

func NewClientTracer() *ClientTracer {
	// Initialize the gRPC client for the tracing service
	conn, _ := grpc.Dial(DefaultAddress, grpc.WithInsecure())
	client := pb.NewTraceServiceClient(conn)
	return &ClientTracer{tracer: client}
}

func (c *ClientTracer) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := c.HealthCheck(context.Background(), &pb.Empty{})
		if err == nil && response.Status {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (c *ClientTracer) HealthCheck(ctx context.Context, request *pb.Empty) (*pb.HealthCheckResponse, error) {
	return c.tracer.HealthCheck(ctx, request)
}

func (c *ClientTracer) StartTrace(ctx context.Context, request *pb.StartTraceRequest) (*pb.TraceResponse, error) {
	return c.tracer.StartTrace(ctx, request)
}

func (c *ClientTracer) Trace(ctx context.Context, request *pb.TraceRequest) (*pb.TraceResponse, error) {
	return c.tracer.Trace(ctx, request)
}

func (c *ClientTracer) StartNewSpan(ctx context.Context, request *pb.StartSpanRequest) (*pb.TraceResponse, error) {
	return c.tracer.StartNewSpan(ctx, request)
}

func (c *ClientTracer) Span(ctx context.Context, request *pb.SpanRequest) (*pb.TraceResponse, error) {
	return c.tracer.Span(ctx, request)
}

func (c *ClientTracer) DatabaseSpan(ctx context.Context, request *pb.DatabaseSpanRequest) (*pb.TraceResponse, error) {
	return c.tracer.DatabaseSpan(ctx, request)
}

func (c *ClientTracer) CloseTrace(ctx context.Context, request *pb.CloseTraceRequest) (*pb.TraceResponse, error) {
	return c.tracer.CloseTrace(ctx, request)
}
