package app

import (
	"context"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/stretchr/testify/mock"
)

// MockTraceService is a mock implementation of the TraceService interface
type MockTraceService struct {
	mock.Mock
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
