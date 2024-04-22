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
	HealthCheckMetrics(ctx context.Context, start *pb.EmptyMetrics) (*pb.HealthCheckResponseMetrics, error)
	WaitForHealthyState() bool
	FetchMetrics(ctx context.Context, request *pb.EmptyMetrics) (*pb.MetricsResponse, error)
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
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := c.HealthCheckMetrics(context.Background(), &pb.EmptyMetrics{})
		if err == nil && response.Status {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (c *ClientMetrics) HealthCheckMetrics(ctx context.Context, request *pb.EmptyMetrics) (*pb.HealthCheckResponseMetrics, error) {
	return c.metrics.HealthCheckMetrics(ctx, request)
}

func (c *ClientMetrics) FetchMetrics(ctx context.Context, request *pb.EmptyMetrics) (*pb.MetricsResponse, error) {
	return c.metrics.FetchMetrics(ctx, request)
}
