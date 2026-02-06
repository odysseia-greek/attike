package models

import "time"

type RollupDoc struct {
	SchemaVersion int       `json:"schema_version"`
	Timestamp     time.Time `json:"@timestamp"`
	BucketMs      int       `json:"bucket_ms"`
	Count         int       `json:"count"`

	Node struct {
		Name string `json:"name"`
	} `json:"node"`

	Workload struct {
		Namespace string `json:"namespace"`
		PodName   string `json:"pod_name"`
		PodUid    string `json:"pod_uid"`
	} `json:"workload"`

	CpuMcores struct {
		Avg float64 `json:"avg"`
		Max float64 `json:"max"` // use float64 to be safe
		P95 float64 `json:"p95"`
	} `json:"cpu_mcores"`

	MemBytes struct {
		Avg float64 `json:"avg"`
		Max float64 `json:"max"`
		P95 float64 `json:"p95"`
	} `json:"mem_bytes"`

	Source struct {
		Aggregator    string `json:"aggregator"`
		Collector     string `json:"collector"`
		Method        string `json:"method"`
		WindowKind    string `json:"window_kind"`
		WindowSamples int    `json:"window_samples"`
	} `json:"source"`
}
