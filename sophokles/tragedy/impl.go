package tragedy

import (
	"context"
	"github.com/odysseia-greek/agora/thales"
	pb "github.com/odysseia-greek/attike/sophokles/proto"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type MetricService interface {
	HealthCheck(ctx context.Context, start *pb.Empty) (*pb.HealthCheckResponse, error)
	WaitForHealthyState() bool
	FetchMetrics(ctx context.Context, request *pb.Empty) (*pb.MetricsResponse, error)
}

type MetricServiceImpl struct {
	PodName   string
	Namespace string
	Kube      *thales.KubeClient
	pb.UnimplementedMetricsServiceServer
	mu sync.Mutex // Mutex to protect the task queue
}

type MetricServiceClient struct {
	Impl MetricService
}

type ClientMetrics struct {
	metrics pb.MetricsServiceClient
}

func NewMetricsClient() *ClientMetrics {
	conn, _ := grpc.Dial(DefaultAddress, grpc.WithInsecure())
	client := pb.NewMetricsServiceClient(conn)
	return &ClientMetrics{metrics: client}
}

func (c *ClientMetrics) WaitForHealthyState() bool {
	timeout := 10 * time.Second
	checkInterval := 500 * time.Millisecond
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := c.HealthCheck(context.Background(), &pb.Empty{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (c *ClientMetrics) HealthCheck(ctx context.Context, request *pb.Empty) (*pb.HealthCheckResponse, error) {
	return c.metrics.HealthCheck(ctx, request)
}

func (c *ClientMetrics) FetchMetrics(ctx context.Context, request *pb.Empty) (*pb.MetricsResponse, error) {
	return c.metrics.FetchMetrics(ctx, request)
}
