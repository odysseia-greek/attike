package tragedy

import (
	"encoding/json"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
)

// TraceDoc is the top-level document for your "trace" index.
type TraceDoc struct {
	Items          []TraceItem `json:"items,omitempty"` // nested in ES
	IsActive       bool        `json:"isActive"`
	TimeStarted    string      `json:"timeStarted"`
	TimeEnded      *string     `json:"timeEnded,omitempty"`
	TotalTime      int64       `json:"totalTime"`
	ResponseCode   int16       `json:"responseCode"`
	Operation      string      `json:"operation,omitempty"`
	ContainsDBSpan bool        `json:"containsDBSpan"`
	NumberOfItems  int32       `json:"numberOfItems"`
}

// TraceItem is stored inside TraceDoc.items (nested)
type TraceItem struct {
	Timestamp    string          `json:"timestamp"`
	ItemType     string          `json:"item_type"`
	SpanID       string          `json:"span_id,omitempty"`
	ParentSpanID string          `json:"parent_span_id,omitempty"`
	PodName      string          `json:"pod_name,omitempty"`
	Namespace    string          `json:"namespace,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"` // JSON form of the protobuf event
}

type EnqueueTraceItem struct {
	TraceID        string  `json:"traceID"`
	IsActive       bool    `json:"isActive"`
	TimeStarted    string  `json:"timeStarted"`
	Operation      string  `json:"operation"`
	TimeEnded      *string `json:"timeEnded,omitempty"`
	TotalTime      int64   `json:"totalTime"`
	ResponseCode   int16   `json:"responseCode"`
	ContainsDBSpan bool    `json:"containsDBSpan"`
	NumberOfItems  int32   `json:"numberOfItems"`
}

// ParseMetricSample parses the JSON payload inside Epistello.Data into a typed doc.
func ParseMetricSample(msg *pb.Epistello) (MetricSampleDoc, error) {
	var doc MetricSampleDoc
	if err := json.Unmarshal([]byte(msg.GetData()), &doc); err != nil {
		return MetricSampleDoc{}, err
	}
	return doc, nil
}

// MetricsSnapshot captures point-in-time resource usage context.
type MetricsSnapshot struct {
	Timestamp time.Time          `json:"@timestamp"`
	Pod       PodUsage           `json:"pod"`
	Node      *NodeUsage         `json:"node,omitempty"`
	Source    MetricsSourceTrace `json:"source"`
}

type PodUsage struct {
	CPUMCores int64 `json:"cpu_mcores"`
	MemBytes  int64 `json:"mem_bytes"`
}

type NodeUsage struct {
	Name      string `json:"name"`
	CPUMCores int64  `json:"cpu_mcores"`
	MemBytes  int64  `json:"mem_bytes"`
}

type MetricsSourceTrace struct {
	Collector string `json:"collector"` // e.g. sophokles
	Method    string `json:"method"`    // e.g. kubelet_stats_summary
}

//
// METRIC SAMPLES DOCUMENT
//

// MetricSampleDoc is the top-level document for your "metric" index.
type MetricSampleDoc struct {
	SchemaVersion int                `json:"schema_version"`
	Timestamp     time.Time          `json:"@timestamp"`
	Node          MetricNodeRef      `json:"node"`
	NodeTotals    ResourceTotals     `json:"node_totals"`
	Workload      WorkloadRef        `json:"workload"`
	PodTotals     ResourceTotals     `json:"pod_totals"`
	Containers    []ContainerSample  `json:"containers,omitempty"` // nested in ES
	Source        MetricSampleSource `json:"source"`
}

type MetricNodeRef struct {
	Name string `json:"name"`
}

type WorkloadRef struct {
	Namespace string `json:"namespace"`
	PodName   string `json:"pod_name"`
	PodUID    string `json:"pod_uid"`
}

type ResourceTotals struct {
	CPUMCores int64 `json:"cpu_mcores"`
	MemBytes  int64 `json:"mem_bytes"`
}

type ContainerSample struct {
	Name   string         `json:"name"`
	Totals ResourceTotals `json:"totals"`
}

type MetricSampleSource struct {
	Collector        string `json:"collector"`
	Method           string `json:"method"`
	ScrapeDurationMS int64  `json:"scrape_duration_ms"`
}

// MetricsRollupDoc is the top-level document for your "metrics-rollup" index.
type MetricsRollupDoc struct {
	SchemaVersion int                `json:"schema_version"`
	Timestamp     time.Time          `json:"@timestamp"`
	BucketMS      int64              `json:"bucket_ms"`
	Count         int64              `json:"count"`
	Node          MetricNodeRef      `json:"node"`
	Workload      WorkloadRef        `json:"workload"`
	CPUMCores     RollupStatsNumeric `json:"cpu_mcores"`
	MemBytes      RollupStatsNumeric `json:"mem_bytes"`
	Source        RollupSource       `json:"source"`
}

type RollupStatsNumeric struct {
	Avg float64 `json:"avg"`
	Max int64   `json:"max"`
	P95 float64 `json:"p95"`
}

type RollupSource struct {
	Aggregator    string `json:"aggregator"`
	Collector     string `json:"collector"`
	Method        string `json:"method"`
	WindowKind    string `json:"window_kind"`
	WindowSamples int64  `json:"window_samples"`
}
