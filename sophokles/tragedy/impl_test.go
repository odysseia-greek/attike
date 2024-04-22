package tragedy

import (
	"context"
	pb "github.com/odysseia-greek/attike/sophokles/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStartMetricsMocked(t *testing.T) {
	mockService := new(MockMetricsService)

	t.Run("Metrics", func(t *testing.T) {
		request := &pb.Empty{}

		expectedResponse := &pb.MetricsResponse{
			// Set fields of the response here
		}

		client := &MetricServiceClient{
			Impl: mockService, // Set the mock service as the implementation
		}

		// Set expectations on the mock
		mockService.On("FetchMetrics", mock.Anything, request).Return(expectedResponse, nil)

		// Call the method being tested
		response, err := client.Impl.FetchMetrics(context.Background(), request)

		// Check expectations and assertions
		mockService.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		expectedResponse := &pb.HealthCheckResponseMetrics{
			Status: true,
		}

		request := pb.Empty{}

		client := &MetricServiceClient{
			Impl: mockService, // Set the mock service as the implementation
		}
		// Set expectations on the mock
		mockService.On("HealthCheck", mock.Anything, &request).Return(expectedResponse, nil)

		// Call the method being tested
		response := client.Impl.WaitForHealthyState()

		// Check expectations and assertions
		mockService.AssertExpectations(t)
		assert.True(t, response)
	})
}
