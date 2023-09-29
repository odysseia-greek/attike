package schemas

import "github.com/odysseia-greek/aristoteles/models"

func parseHitsToGraphql(hits *models.Hits) []TraceObject {
	var traces []TraceObject

	for _, hit := range hits.Hits {
		innerObject := TraceObject{
			TraceID:      hit.ID,
			TimeEnded:    getStringFromSource(hit.Source, "timeEnded"),
			TimeStarted:  getStringFromSource(hit.Source, "timeStarted"),
			TotalTime:    getIntFromSource(hit.Source, "totalTime"),
			ResponseCode: getIntFromSource(hit.Source, "responseCode"),
			IsActive:     getBoolFromSource(hit.Source, "isActive"),
		}

		sourceItems, ok := hit.Source["items"].([]map[string]interface{})
		if !ok {
			return traces
		}

		var graphqlObjects []interface{}

		for _, item := range sourceItems {
			if itemType, exists := item["item_type"].(string); exists {

				switch itemType {
				case "trace":
					var trace Trace
					trace.ParentSpanID = getStringFromMap(item, "parent_span_id")
					trace.Method = getStringFromMap(item, "method")
					trace.URL = getStringFromMap(item, "url")
					trace.Host = getStringFromMap(item, "host")
					trace.RemoteAddress = getStringFromMap(item, "remote_address")
					trace.Timestamp = getStringFromMap(item, "timestamp")
					trace.PodName = getStringFromMap(item, "pod_name")
					trace.Namespace = getStringFromMap(item, "namespace")
					trace.ItemType = itemType
					trace.Operation = getStringFromMap(item, "operation")
					trace.RootQuery = getStringFromMap(item, "root_query")
					graphqlObjects = append(graphqlObjects, &trace)

				case "span":
					var span Span
					span.ParentSpanID = getStringFromMap(item, "parent_span_id")
					span.Namespace = getStringFromMap(item, "namespace")
					span.Timestamp = getStringFromMap(item, "timestamp")
					span.PodName = getStringFromMap(item, "pod_name")
					span.ItemType = itemType
					span.Action = getStringFromMap(item, "action")
					graphqlObjects = append(graphqlObjects, &span)

				case "database_span":
					var dbSpan DatabaseSpan
					dbSpan.ParentSpanID = getStringFromMap(item, "parent_span_id")
					dbSpan.ResultJSON = getStringFromMap(item, "result_json")
					dbSpan.SpanID = getStringFromMap(item, "span_id")
					dbSpan.ItemType = itemType
					dbSpan.Query = getStringFromMap(item, "query")
					dbSpan.Namespace = getStringFromMap(item, "namespace")
					dbSpan.Timestamp = getStringFromMap(item, "timestamp")
					dbSpan.PodName = getStringFromMap(item, "pod_name")
					graphqlObjects = append(graphqlObjects, &dbSpan)

				default:
					// Handle unknown item types or ignore them.
				}
			} else {
				// "item_type" is missing, treat as a "closer" or parse as a Span.
				var span Span
				span.ParentSpanID = getStringFromMap(item, "parent_span_id")
				span.Namespace = getStringFromMap(item, "namespace")
				span.Timestamp = getStringFromMap(item, "timestamp")
				span.PodName = getStringFromMap(item, "pod_name")
				span.ItemType = "trace"
				span.ResponseBody = getStringFromMap(item, "response_body")
				graphqlObjects = append(graphqlObjects, &span)
			}
		}
		innerObject.Items = graphqlObjects
		traces = append(traces, innerObject)
	}

	return traces
}

type TraceObject struct {
	TraceID      string        `json:"traceID"`
	TimeEnded    string        `json:"timeEnded"`
	TimeStarted  string        `json:"timeStarted"`
	TotalTime    int           `json:"totalTime"`
	ResponseCode int           `json:"responseCode"`
	IsActive     bool          `json:"isActive"`
	Items        []interface{} `json:"items"`
}

type Trace struct {
	ParentSpanID  string `json:"parent_span_id"`
	Method        string `json:"method"`
	URL           string `json:"url"`
	Host          string `json:"host"`
	RemoteAddress string `json:"remote_address"`
	Timestamp     string `json:"timestamp"`
	PodName       string `json:"pod_name"`
	Namespace     string `json:"namespace"`
	ItemType      string `json:"item_type"`
	Operation     string `json:"operation"`
	RootQuery     string `json:"root_query"`
}

type Span struct {
	ParentSpanID string `json:"parent_span_id"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
	PodName      string `json:"pod_name"`
	ItemType     string `json:"item_type"`
	Action       string `json:"action"`
	ResponseBody string `json:"response_body"`
}

type DatabaseSpan struct {
	ParentSpanID string `json:"parent_span_id"`
	ResultJSON   string `json:"result_json"`
	SpanID       string `json:"span_id"`
	ItemType     string `json:"item_type"`
	Query        string `json:"query"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
	PodName      string `json:"pod_name"`
}

// Utility function to safely retrieve a string value from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, exists := m[key].(string); exists {
		return val
	}
	return ""
}

// Utility function to safely retrieve an integer value from a map
func getIntFromMap(m map[string]interface{}, key string) int {
	if val, exists := m[key].(float64); exists {
		return int(val)
	}
	return 0
}

// Utility function to safely retrieve a boolean value from a map
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, exists := m[key].(bool); exists {
		return val
	}
	return false
}

// Utility function to safely retrieve a string value from the source
func getStringFromSource(source map[string]interface{}, key string) string {
	if val, exists := source[key].(string); exists {
		return val
	}
	return ""
}

// Utility function to safely retrieve an integer value from the source
func getIntFromSource(source map[string]interface{}, key string) int {
	if val, exists := source[key].(float64); exists {
		return int(val)
	}
	return 0
}

// Utility function to safely retrieve a boolean value from the source
func getBoolFromSource(source map[string]interface{}, key string) bool {
	if val, exists := source[key].(bool); exists {
		return val
	}
	return false
}
