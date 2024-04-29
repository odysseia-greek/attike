package tragedy

type ClusterMetrics struct {
	CpuUnits    string         `json:"cpuUnits"`
	MemoryUnits string         `json:"memoryUnits"`
	Pods        []*PodCluster  `json:"pods"`
	Nodes       []*NodeCluster `json:"nodes"`
	Grouped     []*PodCluster  `json:"grouped"`
}

type PodCluster struct {
	Name                string `json:"name"`
	CpuRaw              int64  `json:"cpuRaw"`
	MemoryRaw           int64  `json:"memoryRaw"`
	CpuHumanReadable    string `json:"cpuHumanReadable"`
	MemoryHumanReadable string `json:"memoryHumanReadable"`
}

type NodeCluster struct {
	NodeName                      string  `json:"nodeName"`
	CpuRaw                        int64   `json:"cpuRaw"`
	MemoryRaw                     int64   `json:"memoryRaw"`
	CpuPercentage                 float64 `json:"cpuPercentage"`
	MemoryPercentage              float64 `json:"memoryPercentage"`
	CpuHumanReadable              string  `json:"cpuHumanReadable"`
	MemoryHumanReadable           string  `json:"memoryHumanReadable"`
	CpuPercentageHumanReadable    string  `json:"cpuPercentageHumanReadable"`
	MemoryPercentageHumanReadable string  `json:"memoryPercentageHumanReadable"`
}
