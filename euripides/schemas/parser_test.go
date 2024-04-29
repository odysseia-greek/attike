package schemas

import (
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHitsToGraphql(t *testing.T) {
	// Create a sample input for testing
	hits := &models.Hits{
		Hits: []models.Hit{
			{
				ID: "1",
				Source: map[string]interface{}{
					"timeEnded":    "2023-09-04T09:40:15.773",
					"timeStarted":  "2023-09-04T09:40:15.633",
					"totalTime":    140.0, // Should be converted to int
					"responseCode": 200.0, // Should be converted to int
					"isActive":     false,
					"items": []map[string]interface{}{
						{
							"parent_span_id": "af4090be495a68ba",
							"method":         "POST",
							"url":            "/graphql",
							"host":           "odysseia-greek.com",
							"remote_address": "10.42.1.57:41792",
							"timestamp":      "2023-09-04'T'09:40:15.633",
							"pod_name":       "homeros-564fcc58d7",
							"namespace":      "odysseia",
							"item_type":      "trace",
							"operation":      "methods",
							"root_query":     "query methods {...}",
						},
					},
				},
			},
		},
	}

	// Call the function to be tested
	traces := parseTracesToGraphql(hits)

	// Use testify/assert for assertions
	assert := assert.New(t)

	assert.Equal(1, len(traces), "Expected 1 trace")

	trace := traces[0]

	assert.Equal("1", trace.TraceID, "Expected TraceID to be '1'")
	assert.Equal("2023-09-04T09:40:15.773", trace.TimeEnded, "Expected TimeEnded to be '2023-09-04T09:40:15.773'")
	assert.Equal("2023-09-04T09:40:15.633", trace.TimeStarted, "Expected TimeStarted to be '2023-09-04T09:40:15.633'")
	assert.Equal(140, trace.TotalTime, "Expected TotalTime to be 140")
	assert.Equal(200, trace.ResponseCode, "Expected ResponseCode to be 200")
	assert.Equal(false, trace.IsActive, "Expected IsActive to be false")
	assert.Equal(1, len(trace.Items), "Expected 1 item in the trace")

	// Perform similar assertions for the items within the trace as needed

}
