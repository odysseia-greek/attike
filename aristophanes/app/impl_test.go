package app

import (
	"context"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartTraceMocked(t *testing.T) {
	mockService := new(MockTraceService)

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
}
