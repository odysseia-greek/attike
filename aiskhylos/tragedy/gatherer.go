package tragedy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/thales"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
	"sort"
	"strings"
	"time"
)

type MetricsGathererImpl struct {
	Index     string
	Namespace string
	Interval  time.Duration
	Kube      *thales.KubeClient
	Elastic   aristoteles.Client
}

func (m *MetricsGathererImpl) GatherMetricsOnTimerFull() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	metrics := &ClusterMetrics{
		MemoryUnits: "Mi",
		CpuUnits:    "millicores",
	}

	podMetrics, err := m.Kube.MetricsClient().MetricsV1beta1().PodMetricses(m.Namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	// Initialize a map to hold aggregated metrics for pod groups.
	podGroupMetrics := make(map[string]*PodCluster)

	for _, pod := range podMetrics.Items {
		var totalMemory int64
		var totalCpu int64

		// Aggregate container metrics for each pod.
		for _, container := range pod.Containers {
			totalMemory += container.Usage.Memory().Value() / (1024 * 1024) // Convert to MiB
			totalCpu += container.Usage.Cpu().MilliValue()                  // MilliCPUs
		}

		// Create a PodCluster for the individual pod.
		individualPod := &PodCluster{
			Name:                pod.Name,
			CpuRaw:              totalCpu,
			MemoryRaw:           totalMemory,
			MemoryHumanReadable: fmt.Sprintf("%vMi", totalMemory),
			CpuHumanReadable:    fmt.Sprintf("%vm", totalCpu),
		}
		metrics.Pods = append(metrics.Pods, individualPod)

		// Determine the group name (prefix) for the pod.
		groupName := strings.Split(pod.Name, "-")[0]

		if _, exists := podGroupMetrics[groupName]; !exists {
			podGroupMetrics[groupName] = &PodCluster{
				Name:      groupName,
				CpuRaw:    0,
				MemoryRaw: 0,
			}
		}

		// Aggregate metrics for the pod group.
		podGroup := podGroupMetrics[groupName]
		podGroup.CpuRaw += totalCpu
		podGroup.MemoryRaw += totalMemory
	}

	// Convert the aggregated group metrics into a slice of PodClusters.
	for _, group := range podGroupMetrics {
		group.CpuHumanReadable = fmt.Sprintf("%vm", group.CpuRaw)
		group.MemoryHumanReadable = fmt.Sprintf("%vMi", group.MemoryRaw)
		metrics.Grouped = append(metrics.Grouped, group)
	}

	nodeMetrics, err := m.Kube.MetricsClient().MetricsV1beta1().NodeMetricses().List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, nodeMetric := range nodeMetrics.Items {
		nodePlaceholder := &NodeCluster{
			NodeName: nodeMetric.Name,
		}

		nodeCPUUsage := nodeMetric.Usage.Cpu().MilliValue()
		nodeMemUsage := nodeMetric.Usage.Memory().Value() / 1024 / 1024

		nodePlaceholder.CpuRaw = nodeCPUUsage
		nodePlaceholder.MemoryRaw = nodeMemUsage
		nodePlaceholder.CpuHumanReadable = fmt.Sprintf("%vm", nodeCPUUsage)
		nodePlaceholder.MemoryHumanReadable = fmt.Sprintf("%vMi", nodeMemUsage)

		node, err := m.Kube.CoreV1().Nodes().Get(ctx, nodeMetric.Name, v1.GetOptions{})
		if err != nil {
			return err
		}

		resourceCPU := node.Status.Capacity[corev1.ResourceCPU]
		resourceMemory := node.Status.Capacity[corev1.ResourceMemory]
		cpuCapacityRaw := resourceCPU.MilliValue()
		memCapacityMiB := resourceMemory.Value() / (1024 * 1024)

		cpuUsagePercentage := float64(nodeCPUUsage) / float64(cpuCapacityRaw) * 100
		memoryUsagePercentage := float64(nodeMemUsage) / float64(memCapacityMiB) * 100

		roundedCpuUsagePercentage := math.Round(cpuUsagePercentage)
		roundedMemoryUsagePercentage := math.Round(memoryUsagePercentage)

		nodePlaceholder.CpuPercentage = cpuUsagePercentage
		nodePlaceholder.MemoryPercentage = memoryUsagePercentage
		nodePlaceholder.CpuPercentageHumanReadable = fmt.Sprintf("%d%%", int(roundedCpuUsagePercentage))
		nodePlaceholder.MemoryPercentageHumanReadable = fmt.Sprintf("%d%%", int(roundedMemoryUsagePercentage))

		metrics.Nodes = append(metrics.Nodes, nodePlaceholder)

	}
	jsonData := map[string]interface{}{
		"metrics":   metrics,
		"timeStamp": time.Now().UTC().Format("2006-01-02T15:04:05.000"),
	}

	data, err := json.Marshal(&jsonData)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	doc, err := m.Elastic.Document().CreateWithId(m.Index, id, data)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created new document: %s", doc.ID))
	createSummary(metrics)

	return nil
}

func createSummary(metrics *ClusterMetrics) {
	logging.Debug("metrics summary")
	for _, node := range metrics.Nodes {
		logging.Debug(fmt.Sprintf("node %s:\nCPU %s %s | MEMORY %s %s", node.NodeName, node.CpuHumanReadable, node.CpuPercentageHumanReadable, node.MemoryHumanReadable, node.MemoryPercentageHumanReadable))
	}

	sort.Slice(metrics.Grouped, func(i, j int) bool {
		return metrics.Grouped[i].CpuRaw > metrics.Grouped[j].CpuRaw
	})

	topCPUByGroup := metrics.Grouped
	if len(topCPUByGroup) > 3 {
		topCPUByGroup = topCPUByGroup[:3]
	}

	group := "top 3 CPU Users by Group:\n"
	for _, pod := range topCPUByGroup {
		group += fmt.Sprintf("%s: %d (%s)\n", pod.Name, pod.CpuRaw, pod.CpuHumanReadable)
	}

	logging.Debug(group)

	sort.Slice(metrics.Grouped, func(i, j int) bool {
		return metrics.Grouped[i].MemoryRaw > metrics.Grouped[j].MemoryRaw
	})

	topMemoryByGroup := metrics.Grouped
	if len(topMemoryByGroup) > 3 {
		topMemoryByGroup = topMemoryByGroup[:3]
	}

	memoryGroup := "top 3 Memory Users by Group:\n"
	for _, pod := range topMemoryByGroup {
		memoryGroup += fmt.Sprintf("%s: %d (%s)\n", pod.Name, pod.MemoryRaw, pod.MemoryHumanReadable)
	}

	logging.Debug(memoryGroup)

	sort.Slice(metrics.Pods, func(i, j int) bool {
		return metrics.Pods[i].CpuRaw > metrics.Pods[j].CpuRaw
	})

	topCPU := metrics.Pods
	if len(topCPU) > 5 {
		topCPU = topCPU[:5]
	}

	topByPod := "top 5 CPU Users by POD:\n"
	for _, pod := range topCPU {
		topByPod += fmt.Sprintf("%s: %d (%s)\n", pod.Name, pod.CpuRaw, pod.CpuHumanReadable)
	}

	logging.Debug(topByPod)

	sort.Slice(metrics.Pods, func(i, j int) bool {
		return metrics.Pods[i].MemoryRaw > metrics.Pods[j].MemoryRaw
	})

	topMemory := metrics.Pods
	if len(topMemory) > 5 {
		topMemory = topMemory[:5]
	}

	topByPodMemory := "top 5 Memory Users by POD:\n"
	for _, pod := range topMemory {
		topByPodMemory += fmt.Sprintf("%s: %d (%s)\n", pod.Name, pod.MemoryRaw, pod.MemoryHumanReadable)
	}
	logging.Debug(topByPodMemory)
}

func (t *TraceServiceImpl) UpdateDocumentWithRetry(traceID, itemName string, data []byte) (string, error) {
	maxRetries := 3
	retryDelay := 100 * time.Millisecond
	var tenTriesError error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		doc, err := t.Elastic.Document().AddItemToDocument(t.Index, traceID, string(data), itemName)

		if err == nil {
			return doc.ID, nil
		}

		if attempt < maxRetries {
			tenTriesError = err
			// Sleep before the next retry
			time.Sleep(retryDelay)
		}
	}

	return "", fmt.Errorf("error updating document for trace ID %s: %s", traceID, tenTriesError.Error())
}
