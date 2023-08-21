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
)

func TestTraceService(t *testing.T) {
	t.Run("StartTrace", func(t *testing.T) {
		request := &pb.StartTraceRequest{
			Method:        "GET",
			Url:           "http://test.com",
			Host:          "http://iamamock.com",
			RemoteAddress: "http://remotecaller.com",
			RootQuery:     "/graphql",
			Operation:     "Something",
		}

		t.Run("SuccessfulStart", func(t *testing.T) {
			fixtureFile := "info"
			mockCode := 200
			mockElasticClient, err := aristoteles.NewMockClient(fixtureFile, mockCode)
			assert.Nil(t, err)

			// Create the handler using the mock
			handler := &TraceServiceImpl{
				PodName:   "testpod",
				Namespace: "testnamespace",
				Index:     "test",
				Elastic:   mockElasticClient,
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

			// Create the handler using the mock
			handler := &TraceServiceImpl{
				PodName:   "testpod",
				Namespace: "testnamespace",
				Index:     "test",
				Elastic:   mockElasticClient,
			}

			// Call the method being tested
			response, err := handler.StartTrace(context.Background(), request)

			// Check expectations and assertions
			assert.NotNil(t, err)
			assert.Nil(t, response)
		})
	})
}

func TestTraceClose(t *testing.T) {
	t.Run("CloseTrace", func(t *testing.T) {
		request := &pb.CloseTraceRequest{
			TraceId:      "841a4f73-ba5b-4c38-9237-e1ad91459028",
			ParentSpanId: "70b993de1e2f879d",
			ResponseBody: "{ \"sentenceId\": \"EmUfwX8BFAEUBcEdfSb7Ty\", \"answerSentence\": \"i have tried\",\"author\": \"plato\" }",
		}

		t.Run("SuccessfulClose", func(t *testing.T) {
			fixtureFiles := []string{"updated", "updated"}
			mockCode := 200
			mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
			assert.Nil(t, err)

			// Create the handler using the mock
			handler := &TraceServiceImpl{
				PodName:   "testpod",
				Namespace: "testnamespace",
				Index:     "test",
				Elastic:   mockElasticClient,
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
	})

}

func ValidateInput(input string) error {
	parts := strings.Split(input, "+")

	if len(parts) != 3 {
		return errors.New("input does not have three parts")
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

	// Check if the third part is either "0" or "1"
	if parts[2] != "0" && parts[2] != "1" {
		return errors.New("third part should be either 0 or 1")
	}

	return nil
}
