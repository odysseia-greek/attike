package app

import (
	"context"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStartTraceMocked(t *testing.T) {
	mockService := new(MockTraceService)

	t.Run("StartTrace", func(t *testing.T) {
		request := &pb.StartTraceRequest{
			Method:        "GET",
			Url:           "http://test.com",
			Host:          "http://iamamock.com",
			RemoteAddress: "http://remotecaller.com",
			RootQuery:     "/graphql",
			Operation:     "Something",
		}

		expectedResponse := &pb.TraceResponse{
			// Set fields of the response here
		}

		client := &TraceServiceClient{
			Impl: mockService, // Set the mock service as the implementation
		}

		// Set expectations on the mock
		mockService.On("StartTrace", mock.Anything, request).Return(expectedResponse, nil)

		// Call the method being tested
		response, err := client.Impl.StartTrace(context.Background(), request)

		// Check expectations and assertions
		mockService.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		expectedResponse := &pb.HealthCheckResponse{
			Status: true,
		}

		request := pb.Empty{}

		client := &TraceServiceClient{
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
