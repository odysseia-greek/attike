package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/odysseia-greek/aristoteles"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestTraceServiceStartAndClose(t *testing.T) {
	startTimeMap := make(map[string]time.Time)

	t.Run("SuccessfulStart", func(t *testing.T) {
		fixtureFile := "info"
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		request := &pb.StartTraceRequest{
			Method:        "GET",
			Url:           "http://test.com",
			Host:          "http://iamamock.com",
			RemoteAddress: "http://remotecaller.com",
			RootQuery:     "/graphql",
			Operation:     "Something",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.StartTrace(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		split := ValidateInput(response.CombinedId)
		assert.Nil(t, split)
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "serviceDown"
		mockCode := 502
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		request := &pb.StartTraceRequest{
			Method:        "GET",
			Url:           "http://test.com",
			Host:          "http://iamamock.com",
			RemoteAddress: "http://remotecaller.com",
			RootQuery:     "/graphql",
			Operation:     "Something",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.StartTrace(context.Background(), request)

		// Check expectations and assertions
		assert.NotNil(t, err)
		assert.Nil(t, response)
	})

	t.Run("SuccessfulClose", func(t *testing.T) {
		fixtureFiles := []string{"updated", "updated"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		request := &pb.CloseTraceRequest{
			TraceId:      "841a4f73-ba5b-4c38-9237-e1ad91459028",
			ParentSpanId: "70b993de1e2f879d",
			ResponseBody: "{ \"sentenceId\": \"EmUfwX8BFAEUBcEdfSb7Ty\", \"answerSentence\": \"i have tried\",\"author\": \"plato\" }",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.CloseTrace(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		// Check if the first part is a valid UUID
		if _, err := uuid.Parse(response.CombinedId); err != nil {
			assert.Nil(t, err)
		}
	})
}

func TestCreateItemTypes(t *testing.T) {
	startTimeMap := make(map[string]time.Time)

	t.Run("StartNewSpan", func(t *testing.T) {
		fixtureFiles := []string{"info"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		request := &pb.StartSpanRequest{
			TraceId: "841a4f73-ba5b-4c38-9237-e1ad91459028",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.StartNewSpan(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		split := ValidateInput(response.CombinedId)
		assert.Nil(t, split)
	})

	t.Run("CreateDatabaseSpan", func(t *testing.T) {
		fixtureFiles := []string{"info"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		request := &pb.DatabaseSpanRequest{
			TraceId:      "841a4f73-ba5b-4c38-9237-e1ad91459028",
			ParentSpanId: "70b993de1e2f879d",
			Action:       "query",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.DatabaseSpan(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		split := ValidateInput(response.CombinedId)
		assert.Nil(t, split)
	})

	t.Run("CreateSpan", func(t *testing.T) {
		fixtureFiles := []string{"info"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		request := &pb.SpanRequest{
			TraceId:      "841a4f73-ba5b-4c38-9237-e1ad91459028",
			ParentSpanId: "70b993de1e2f879d",
			Action:       "query",
			RequestBody:  "{}",
			ResponseBody: "{ \"sentenceId\": \"EmUfwX8BFAEUBcEdfSb7Ty\", \"answerSentence\": \"i have tried\",\"author\": \"plato\" }",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.Span(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		split := ValidateInput(response.CombinedId)
		assert.Nil(t, split)
	})

	t.Run("CreateTrace", func(t *testing.T) {
		fixtureFiles := []string{"info"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		request := &pb.TraceRequest{
			TraceId:      "841a4f73-ba5b-4c38-9237-e1ad91459028",
			ParentSpanId: "70b993de1e2f879d",
			Method:       "GET",
			Url:          "http://test.com",
			Host:         "http://iamamock.com",
		}

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		response, err := handler.Trace(context.Background(), request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)

		split := ValidateInput(response.CombinedId)
		assert.Nil(t, split)
	})
}

func TestHealthChecks(t *testing.T) {
	startTimeMap := make(map[string]time.Time)

	t.Run("Healthy", func(t *testing.T) {
		fixtureFiles := []string{"info"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)

		// Create the handler using the mock
		handler := &TraceServiceImpl{
			PodName:      "testpod",
			Namespace:    "testnamespace",
			Index:        "test",
			Elastic:      mockElasticClient,
			StartTimeMap: startTimeMap,
		}

		// Call the method being tested
		request := pb.Empty{}
		response, err := handler.HealthCheck(context.Background(), &request)

		// Check expectations and assertions
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})
}

func ValidateInput(input string) error {
	parts := strings.Split(input, "+")

	if len(parts) != 2 && len(parts) != 3 {
		return errors.New("input must have either two or three parts")
	}

	// Check if the first part is a valid UUID
	if _, err := uuid.Parse(parts[0]); err != nil {
		return errors.New("invalid UUID format")
	}

	// Check if the second part is a valid hexadecimal string
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", parts[1])
	if !match {
		return errors.New("invalid hexadecimal string format")
	}

	// If there's a third part, check if it is either "0" or "1"
	if len(parts) == 3 && parts[2] != "0" && parts[2] != "1" {
		return errors.New("third part should be either 0 or 1")
	}

	return nil
}
