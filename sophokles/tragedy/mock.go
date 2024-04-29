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
		response, err := m.HealthCheck(context.Background(), &pb.Empty{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MockMetricsService) HealthCheck(ctx context.Context, request *pb.Empty) (*pb.HealthCheckResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.HealthCheckResponse), args.Error(1)
}

func (m *MockMetricsService) FetchMetrics(ctx context.Context, request *pb.Empty) (*pb.MetricsResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.MetricsResponse), args.Error(1)
}
