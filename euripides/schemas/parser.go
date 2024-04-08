package schemas

import "github.com/odysseia-greek/agora/aristoteles/models"

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

		sourceItems := make([]map[string]interface{}, 0)
		sources, ok := hit.Source["items"].([]interface{})
		if !ok {
			continue
		} else {
			for _, item := range sources {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				} else {
					sourceItems = append(sourceItems, itemMap)
				}
			}
		}

		var graphqlObjects []interface{}

		for _, item := range sourceItems {
			if common, exists := item["common"].(map[string]interface{}); exists {
				if itemType, exists := common["item_type"].(string); exists {
					switch itemType {
					case "trace":
						var trace Trace
						trace.ParentSpanID = getStringFromMap(common, "parent_span_id")
						trace.Method = getStringFromMap(item, "method")
						trace.URL = getStringFromMap(item, "url")
						trace.Host = getStringFromMap(item, "host")
						trace.RemoteAddress = getStringFromMap(item, "remote_address")
						trace.Operation = getStringFromMap(item, "operation")
						trace.RootQuery = getStringFromMap(item, "root_query")
						trace.Timestamp = getStringFromMap(common, "timestamp")
						trace.PodName = getStringFromMap(common, "pod_name")
						trace.Namespace = getStringFromMap(common, "namespace")
						trace.ItemType = itemType
						graphqlObjects = append(graphqlObjects, &trace)

					case "span":
						var span Span
						span.ParentSpanID = getStringFromMap(common, "parent_span_id")
						span.Namespace = getStringFromMap(common, "namespace")
						span.Timestamp = getStringFromMap(common, "timestamp")
						span.PodName = getStringFromMap(common, "pod_name")
						span.TimeStarted = getStringFromMap(item, "time_started")
						span.TimeFinished = getStringFromMap(item, "time_finished")
						span.RequestBody = getStringFromMap(item, "request_body")
						span.ResponseCode = getIntFromMap(item, "response_code")
						span.ItemType = itemType
						span.Action = getStringFromMap(item, "action")
						graphqlObjects = append(graphqlObjects, &span)

					case "database_span":
						var dbSpan DatabaseSpan
						dbSpan.ParentSpanID = getStringFromMap(common, "parent_span_id")
						dbSpan.SpanID = getStringFromMap(common, "span_id")
						dbSpan.ItemType = itemType
						dbSpan.Query = getStringFromMap(item, "query")
						dbSpan.Took = getStringFromMap(item, "took")
						dbSpan.Hits = int64(getIntFromMap(item, "hits"))
						dbSpan.Namespace = getStringFromMap(common, "namespace")
						dbSpan.Timestamp = getStringFromMap(common, "timestamp")
						dbSpan.PodName = getStringFromMap(common, "pod_name")
						span.TimeStarted = getStringFromMap(item, "time_started")
						span.TimeFinished = getStringFromMap(item, "time_finished")
						graphqlObjects = append(graphqlObjects, &dbSpan)

					default:
						// Handle unknown item types or ignore them.
					}
				} else {
					// "item_type" is missing, treat as a "closer" or parse as a Span.
					var span Span
					span.ParentSpanID = getStringFromMap(common, "parent_span_id")
					span.Namespace = getStringFromMap(common, "namespace")
					span.Timestamp = getStringFromMap(common, "timestamp")
					span.PodName = getStringFromMap(common, "pod_name")
					span.ItemType = "trace"
					span.Action = getStringFromMap(item, "action")
					graphqlObjects = append(graphqlObjects, &span)
				}
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
	RequestBody  string `json:"request_body"`
	TimeStarted  string `json:"time_started"`
	TimeFinished string `json:"time_finished"`
	ResponseCode int    `json:"response_code"`
}

type DatabaseSpan struct {
	ParentSpanID string `json:"parent_span_id"`
	Hits         int64  `json:"hits"`
	Took         string `json:"took"`
	TimeFinished string `json:"time_finished"`
	TimeStarted  string `json:"time_started"`
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
