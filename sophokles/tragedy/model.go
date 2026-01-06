package tragedy

import "time"

// MetricSample is the unit Sophokles emits (to logs for now; later to Eupalinos).
type MetricSample struct {
	SchemaVersion int       `json:"schema_version"`
	Timestamp     time.Time `json:"@timestamp"`

	Node       NodeRef         `json:"node"`
	NodeTotals *ResourceTotals `json:"node_totals,omitempty"`
	Workload   WorkloadRef     `json:"workload"`

	PodTotals  ResourceTotals    `json:"pod_totals"`
	Containers []ContainerSample `json:"containers,omitempty"`

	Source SourceMeta `json:"source,omitempty"`
}

type NodeRef struct {
	Name string `json:"name"`
}

type WorkloadRef struct {
	Namespace string `json:"namespace"`
	PodName   string `json:"pod_name"`
	PodUID    string `json:"pod_uid,omitempty"`
}

type ResourceTotals struct {
	CPUMillicores int64 `json:"cpu_mcores"`
	MemoryBytes   int64 `json:"mem_bytes"`
}

type ContainerSample struct {
	Name   string         `json:"name"`
	Totals ResourceTotals `json:"totals"`
}

type SourceMeta struct {
	Collector        string `json:"collector"`
	Method           string `json:"method"` // "metrics_api" (for now)
	ScrapeDurationMS int64  `json:"scrape_duration_ms"`
}

// ---- kubelet /stats/summary JSON structures ----

type StatsSummary struct {
	Node NodeStats  `json:"node"`
	Pods []PodStats `json:"pods"`
}

type NodeStats struct {
	NodeName string      `json:"nodeName"`
	CPU      CPUStats    `json:"cpu"`
	Memory   MemoryStats `json:"memory"`
}

type PodStats struct {
	PodRef     PodRef           `json:"podRef"`
	Containers []ContainerStats `json:"containers"`
}

type PodRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

type ContainerStats struct {
	Name   string      `json:"name"`
	CPU    CPUStats    `json:"cpu"`
	Memory MemoryStats `json:"memory"`
}

type CPUStats struct {
	// Instantaneous CPU usage in nanocores (not always present)
	UsageNanoCores *uint64 `json:"usageNanoCores,omitempty"`
	// Cumulative CPU usage in core-nanoseconds (present more often)
	UsageCoreNanoSeconds *uint64 `json:"usageCoreNanoSeconds,omitempty"`
}

type MemoryStats struct {
	// Working set bytes (good default for “memory used”)
	WorkingSetBytes *uint64 `json:"workingSetBytes,omitempty"`
}
