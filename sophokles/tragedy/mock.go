package tragedy

import (
	"context"
	pb "github.com/odysseia-greek/attike/sophokles/proto"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockMetricsService is a mock implementation of the TraceService interface
type MockMetricsService struct {
	mock.Mock
}

func (m *MockMetricsService) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.HealthCheckMetrics(context.Background(), &pb.EmptyMetrics{})
		if err == nil && response.Status {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MockMetricsService) HealthCheckMetrics(ctx context.Context, request *pb.EmptyMetrics) (*pb.HealthCheckResponseMetrics, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.HealthCheckResponseMetrics), args.Error(1)
}

func (m *MockMetricsService) FetchMetrics(ctx context.Context, request *pb.EmptyMetrics) (*pb.MetricsResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.MetricsResponse), args.Error(1)
}
