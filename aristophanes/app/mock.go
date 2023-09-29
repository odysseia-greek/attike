package app

import (
	"context"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockTraceService is a mock implementation of the TraceService interface
type MockTraceService struct {
	mock.Mock
}

func (m *MockTraceService) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.HealthCheck(context.Background(), &pb.Empty{})
		if err == nil && response.Status {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MockTraceService) HealthCheck(ctx context.Context, request *pb.Empty) (*pb.HealthCheckResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.HealthCheckResponse), args.Error(1)
}

func (m *MockTraceService) StartTrace(ctx context.Context, request *pb.StartTraceRequest) (*pb.TraceResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}

func (m *MockTraceService) Trace(ctx context.Context, request *pb.TraceRequest) (*pb.TraceResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}

func (m *MockTraceService) StartNewSpan(ctx context.Context, request *pb.StartSpanRequest) (*pb.TraceResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}

func (m *MockTraceService) Span(ctx context.Context, request *pb.SpanRequest) (*pb.TraceResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}

func (m *MockTraceService) DatabaseSpan(ctx context.Context, request *pb.DatabaseSpanRequest) (*pb.TraceResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}

func (m *MockTraceService) CloseTrace(ctx context.Context, request *pb.CloseTraceRequest) (*pb.TraceResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*pb.TraceResponse), args.Error(1)
}
