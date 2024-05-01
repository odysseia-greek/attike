package comedy

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	sophokles "github.com/odysseia-greek/attike/sophokles/tragedy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type TraceService interface {
	HealthCheck(ctx context.Context, start *pb.Empty) (*pb.HealthCheckResponse, error)
	Chorus(ctx context.Context) (pb.TraceService_ChorusClient, error)
	WaitForHealthyState() bool
}

type MapCommand struct {
	Action   string
	TraceID  string
	Time     time.Time
	Response chan<- MapResponse
}

type MapResponse struct {
	Time  time.Time
	Found bool
}

type TraceServiceImpl struct {
	PodName       string
	Namespace     string
	Index         string
	Elastic       aristoteles.Client
	Metrics       *sophokles.ClientMetrics
	GatherMetrics bool
	commands      chan MapCommand
	pb.UnimplementedTraceServiceServer
}

type TraceServiceClient struct {
	Impl TraceService
}

type ClientTracer struct {
	tracer pb.TraceServiceClient
}

func NewClientTracer() (*ClientTracer, error) {
	// Initialize the gRPC client for the tracing service
	conn, err := grpc.Dial(DefaultAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewTraceServiceClient(conn)
	return &ClientTracer{tracer: client}, nil
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

func (c *ClientTracer) Chorus(ctx context.Context) (pb.TraceService_ChorusClient, error) {
	return c.tracer.Chorus(ctx)
}
