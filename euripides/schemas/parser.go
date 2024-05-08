package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
)

func parseMetricsToGraphql(hits *models.Hits) MetricsObject {
	startTime := hits.Hits[0].Source["timeStamp"].(string)
	endTime := hits.Hits[len(hits.Hits)-1].Source["timeStamp"].(string)
	metrics := MetricsObject{
		TimeEnded:   endTime,
		TimeStarted: startTime,
		TimeStamps:  make([]string, 0),
	}

	for _, hit := range hits.Hits {
		metricHits, err := parseMetrics(hit.Source)
		metrics.TimeStamps = append(metrics.TimeStamps, hit.Source["timeStamp"].(string))
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		metrics.CpuUnits = metricHits.CpuUnits
		metrics.MemoryUnits = metricHits.MemoryUnits
		metrics.Pods = aggregateMetrics(metrics.Pods, metricHits.Pods)
		metrics.Nodes = aggregateNodes(metrics.Nodes, metricHits.Nodes)
		metrics.Grouped = aggregateMetrics(metrics.Grouped, metricHits.Grouped)
	}

	return metrics
}

func parseMetrics(hitSource map[string]interface{}) (*MetricsElastic, error) {
	metricsData, exists := hitSource["metrics"]
	if !exists {
		return nil, fmt.Errorf("metrics key not found in source")
	}

	metricsMap, ok := metricsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metrics data is not a map")
	}

	jsonData, err := json.Marshal(metricsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics data: %v", err)
	}

	var metrics MetricsElastic
	err = json.Unmarshal(jsonData, &metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics data: %v", err)
	}

	return &metrics, nil
}

func aggregateMetrics(existingSlices []GenericMetricsSlice, newMetrics []GenericMetrics) []GenericMetricsSlice {
	for _, metric := range newMetrics {
		found := false
		for i, existing := range existingSlices {
			if existing.Name == metric.Name {
				existingSlices[i].CpuRaw = append(existing.CpuRaw, metric.CpuRaw)
				existingSlices[i].MemoryRaw = append(existing.MemoryRaw, metric.MemoryRaw)
				existingSlices[i].CpuHumanReadable = append(existing.CpuHumanReadable, metric.CpuHumanReadable)
				existingSlices[i].MemoryHumanReadable = append(existing.MemoryHumanReadable, metric.MemoryHumanReadable)
				found = true
				break
			}
		}
		if !found {
			newSlice := GenericMetricsSlice{
				Name:                metric.Name,
				CpuRaw:              []int{metric.CpuRaw},
				MemoryRaw:           []int{metric.MemoryRaw},
				CpuHumanReadable:    []string{metric.CpuHumanReadable},
				MemoryHumanReadable: []string{metric.MemoryHumanReadable},
			}
			existingSlices = append(existingSlices, newSlice)
		}
	}
	return existingSlices
}

func aggregateNodes(existingSlices []GenericNodeSlice, newMetrics []GenericNode) []GenericNodeSlice {
	for _, metric := range newMetrics {
		found := false
		for i, existing := range existingSlices {
			if existing.NodeName == metric.NodeName {
				existingSlices[i].CpuRaw = append(existing.CpuRaw, metric.CpuRaw)
				existingSlices[i].MemoryRaw = append(existing.MemoryRaw, metric.MemoryRaw)
				existingSlices[i].CpuHumanReadable = append(existing.CpuHumanReadable, metric.CpuHumanReadable)
				existingSlices[i].MemoryHumanReadable = append(existing.MemoryHumanReadable, metric.MemoryHumanReadable)
				existingSlices[i].CpuPercentage = append(existing.CpuPercentage, metric.CpuPercentage)
				existingSlices[i].MemoryPercentage = append(existing.MemoryPercentage, metric.MemoryPercentage)
				existingSlices[i].CpuPercentageHumanReadable = append(existing.CpuPercentageHumanReadable, metric.CpuPercentageHumanReadable)
				existingSlices[i].MemoryPercentageHumanReadable = append(existing.MemoryPercentageHumanReadable, metric.MemoryPercentageHumanReadable)
				found = true
				break
			}
		}
		if !found {
			newSlice := GenericNodeSlice{
				NodeName:                      metric.NodeName,
				CpuRaw:                        []int{metric.CpuRaw},
				MemoryRaw:                     []int{metric.MemoryRaw},
				CpuPercentage:                 []float64{metric.CpuPercentage},
				MemoryPercentage:              []float64{metric.MemoryPercentage},
				CpuHumanReadable:              []string{metric.CpuHumanReadable},
				MemoryHumanReadable:           []string{metric.MemoryHumanReadable},
				CpuPercentageHumanReadable:    []string{metric.CpuPercentageHumanReadable},
				MemoryPercentageHumanReadable: []string{metric.MemoryPercentageHumanReadable},
			}
			existingSlices = append(existingSlices, newSlice)
		}
	}
	return existingSlices
}

func parseTracesToGraphql(hits *models.Hits) []TraceObject {
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
					case "trace_start":
						var trace TraceStart
						traceMetrics := getMetrics(item)
						trace.ParentSpanID = getStringFromMap(common, "parent_span_id")
						trace.SpanID = getStringFromMap(common, "span_id")
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
						if traceMetrics != nil {
							trace.Metrics = traceMetrics
						}
						graphqlObjects = append(graphqlObjects, &trace)

					case "trace_close":
						var trace TraceClose
						traceMetrics := getMetrics(item)
						trace.ParentSpanID = getStringFromMap(common, "parent_span_id")
						trace.ResponseBody = getStringFromMap(item, "response_body")
						trace.Timestamp = getStringFromMap(common, "timestamp")
						trace.PodName = getStringFromMap(common, "pod_name")
						trace.Namespace = getStringFromMap(common, "namespace")
						trace.ItemType = itemType
						if traceMetrics != nil {
							trace.Metrics = traceMetrics
						}
						graphqlObjects = append(graphqlObjects, &trace)

					case "trace":
						var trace Trace
						traceMetrics := getMetrics(item)
						trace.ParentSpanID = getStringFromMap(common, "parent_span_id")
						trace.SpanID = getStringFromMap(common, "span_id")
						trace.Method = getStringFromMap(item, "method")
						trace.URL = getStringFromMap(item, "url")
						trace.Host = getStringFromMap(item, "host")
						trace.Timestamp = getStringFromMap(common, "timestamp")
						trace.PodName = getStringFromMap(common, "pod_name")
						trace.Namespace = getStringFromMap(common, "namespace")
						trace.ItemType = itemType
						if traceMetrics != nil {
							trace.Metrics = traceMetrics
						}
						graphqlObjects = append(graphqlObjects, &trace)

					case "span":
						var span Span
						span.ParentSpanID = getStringFromMap(common, "parent_span_id")
						span.SpanID = getStringFromMap(common, "span_id")
						span.Namespace = getStringFromMap(common, "namespace")
						span.Timestamp = getStringFromMap(common, "timestamp")
						span.PodName = getStringFromMap(common, "pod_name")
						span.Took = getStringFromMap(item, "took")
						span.Status = getStringFromMap(item, "status")
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
						dbSpan.TimeStarted = getStringFromMap(item, "time_started")
						dbSpan.TimeFinished = getStringFromMap(item, "time_finished")
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

type MetricsElastic struct {
	CpuUnits    string           `json:"cpuUnits"`
	MemoryUnits string           `json:"memoryUnits"`
	Pods        []GenericMetrics `json:"pods"`
	Nodes       []GenericNode    `json:"nodes"`
	Grouped     []GenericMetrics `json:"grouped"`
}

type MetricsObject struct {
	TimeEnded   string                `json:"timeEnded"`
	TimeStarted string                `json:"timeStarted"`
	TimeStamps  []string              `json:"timeStamps"`
	CpuUnits    string                `json:"cpuUnits"`
	MemoryUnits string                `json:"memoryUnits"`
	Pods        []GenericMetricsSlice `json:"pods"`
	Nodes       []GenericNodeSlice    `json:"nodes"`
	Grouped     []GenericMetricsSlice `json:"grouped"`
}

type GenericMetrics struct {
	Name                string `json:"name"`
	CpuRaw              int    `json:"cpuRaw"`
	MemoryRaw           int    `json:"memoryRaw"`
	CpuHumanReadable    string `json:"cpuHumanReadable"`
	MemoryHumanReadable string `json:"memoryHumanReadable"`
}

type GenericNodeSlice struct {
	NodeName                      string    `json:"nodeName"`
	CpuRaw                        []int     `json:"cpuRaw"`
	MemoryRaw                     []int     `json:"memoryRaw"`
	CpuPercentage                 []float64 `json:"cpuPercentage"`
	MemoryPercentage              []float64 `json:"memoryPercentage"`
	CpuHumanReadable              []string  `json:"cpuHumanReadable"`
	MemoryHumanReadable           []string  `json:"memoryHumanReadable"`
	CpuPercentageHumanReadable    []string  `json:"cpuPercentageHumanReadable"`
	MemoryPercentageHumanReadable []string  `json:"memoryPercentageHumanReadable"`
}

type GenericNode struct {
	NodeName                      string  `json:"nodeName"`
	CpuRaw                        int     `json:"cpuRaw"`
	MemoryRaw                     int     `json:"memoryRaw"`
	CpuPercentage                 float64 `json:"cpuPercentage"`
	MemoryPercentage              float64 `json:"memoryPercentage"`
	CpuHumanReadable              string  `json:"cpuHumanReadable"`
	MemoryHumanReadable           string  `json:"memoryHumanReadable"`
	CpuPercentageHumanReadable    string  `json:"cpuPercentageHumanReadable"`
	MemoryPercentageHumanReadable string  `json:"memoryPercentageHumanReadable"`
}

type GenericMetricsSlice struct {
	Name                string   `json:"name"`
	CpuRaw              []int    `json:"cpuRaw"`
	MemoryRaw           []int    `json:"memoryRaw"`
	CpuHumanReadable    []string `json:"cpuHumanReadable"`
	MemoryHumanReadable []string `json:"memoryHumanReadable"`
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
	ParentSpanID string   `json:"parent_span_id"`
	SpanID       string   `json:"span_id"`
	Method       string   `json:"method"`
	URL          string   `json:"url"`
	Host         string   `json:"host"`
	Timestamp    string   `json:"timestamp"`
	PodName      string   `json:"pod_name"`
	Namespace    string   `json:"namespace"`
	ItemType     string   `json:"item_type"`
	Metrics      *Metrics `json:"metrics,omitempty"`
}

type TraceStart struct {
	ParentSpanID  string   `json:"parent_span_id"`
	SpanID        string   `json:"span_id"`
	Method        string   `json:"method"`
	URL           string   `json:"url"`
	Host          string   `json:"host"`
	RemoteAddress string   `json:"remote_address"`
	Timestamp     string   `json:"timestamp"`
	PodName       string   `json:"pod_name"`
	Namespace     string   `json:"namespace"`
	ItemType      string   `json:"item_type"`
	Operation     string   `json:"operation"`
	RootQuery     string   `json:"root_query"`
	Metrics       *Metrics `json:"metrics,omitempty"`
}

type TraceClose struct {
	ParentSpanID string   `json:"parent_span_id"`
	Timestamp    string   `json:"timestamp"`
	PodName      string   `json:"pod_name"`
	Namespace    string   `json:"namespace"`
	ItemType     string   `json:"item_type"`
	ResponseBody string   `json:"response_body"`
	Metrics      *Metrics `json:"metrics,omitempty"`
}

type Metrics struct {
	CpuRaw              int    `json:"cpu_raw"`
	CpuHumanReadable    string `json:"cpu_human_readable"`
	MemoryHumanReadable string `json:"memory_human_readable"`
	MemoryUnits         string `json:"memory_units"`
	CpuUnits            string `json:"cpu_units"`
	MemoryRaw           int    `json:"memory_raw"`
}

type Span struct {
	ParentSpanID string `json:"parent_span_id"`
	SpanID       string `json:"span_id"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
	PodName      string `json:"pod_name"`
	ItemType     string `json:"item_type"`
	Action       string `json:"action"`
	Status       string `json:"status"`
	Took         string `json:"took"`
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

func getMetrics(m map[string]interface{}) *Metrics {
	if val, exists := m["metrics"].(map[string]interface{}); exists {
		return &Metrics{
			CpuRaw:              getIntFromMap(val, "cpu_raw"),
			CpuHumanReadable:    getStringFromMap(val, "cpu_human_readable"),
			MemoryHumanReadable: getStringFromMap(val, "memory_human_readable"),
			MemoryUnits:         getStringFromMap(val, "memory_units"),
			CpuUnits:            getStringFromMap(val, "cpu_units"),
			MemoryRaw:           getIntFromMap(val, "memory_raw"),
		}
	}

	return nil
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
