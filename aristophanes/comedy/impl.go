package comedy

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/eupalinos/stomion"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TraceService interface {
	HealthCheck(ctx context.Context, start *v1.Empty) (*v1.HealthCheckResponse, error)
	Chorus(ctx context.Context) (v1.TraceService_ChorusClient, error)
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
	PodName   string
	Namespace string
	commands  chan MapCommand
	Eupalinos *stomion.QueueClient
	Channel   string
	baseCtx   context.Context
	cancel    context.CancelFunc
	v1.UnimplementedTraceServiceServer
}

type Streamer interface {
	Send(*v1.ObserveRequest) error
}

type Option func(*cfg)

type cfg struct {
	HeaderKey      string
	ContextKeyName any
	EmitCloseHop   bool
}

func WithHeaderKey(k string) Option {
	return func(c *cfg) { c.HeaderKey = k }
}

func WithContextKeyName(k any) Option {
	return func(c *cfg) { c.ContextKeyName = k }
}

func WithCloseHop() Option {
	return func(c *cfg) { c.EmitCloseHop = true }
}

type TraceServiceClient struct {
	Impl TraceService
}

type ClientTracer struct {
	tracer v1.TraceServiceClient
}

func NewClientTracer(serviceAddress string) (*ClientTracer, error) {
	// Initialize the gRPC client for the tracing service
	if serviceAddress == "" {
		serviceAddress = DefaultAddress
	}
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewTraceServiceClient(conn)
	return &ClientTracer{tracer: client}, nil
}

func (c *ClientTracer) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := c.HealthCheck(context.Background(), &v1.Empty{})
		if err == nil && response.Status {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (c *ClientTracer) HealthCheck(ctx context.Context, request *v1.Empty) (*v1.HealthCheckResponse, error) {
	return c.tracer.HealthCheck(ctx, request)
}

func (c *ClientTracer) Chorus(ctx context.Context) (v1.TraceService_ChorusClient, error) {
	return c.tracer.Chorus(ctx)
}
