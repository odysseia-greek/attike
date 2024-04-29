package tragedy

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/sophokles/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (m *MetricServiceImpl) HealthCheck(ctx context.Context, health *pb.Empty) (*pb.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pod, err := m.Kube.CoreV1().Pods(m.Namespace).Get(ctx, m.PodName, v1.GetOptions{})

	response := &pb.HealthCheckResponse{
		Healthy: true,
		Status:  string(pod.Status.Phase),
	}

	if err != nil {
		logging.Error(err.Error())
		response.Healthy = false
		response.Status = err.Error()
	}

	return response, nil
}

func (m *MetricServiceImpl) FetchMetrics(ctx context.Context, request *pb.Empty) (*pb.MetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	response := &pb.MetricsResponse{
		MemoryUnits: "mb",
		CpuUnits:    "millicores",
		Pod:         &pb.PodMetrics{},
	}

	podMetrics, err := m.Kube.MetricsClient().MetricsV1beta1().PodMetricses(m.Namespace).Get(ctx, m.PodName, v1.GetOptions{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	var totalMemory int64
	var totalCpu int64

	for _, container := range podMetrics.Containers {
		totalMemory = totalMemory + container.Usage.Memory().Value()/1024/1024
		totalCpu = totalCpu + container.Usage.Cpu().MilliValue()
		containerMetrics := &pb.ContainerMetrics{
			ContainerName:                container.Name,
			ContainerCpuRaw:              container.Usage.Cpu().MilliValue(),
			ContainerMemoryRaw:           container.Usage.Memory().Value() / 1024 / 1024,
			ContainerCpuHumanReadable:    fmt.Sprintf("%vm", container.Usage.Cpu().MilliValue()),
			ContainerMemoryHumanReadable: fmt.Sprintf("%vMi", container.Usage.Memory().Value()/1024/1024),
		}

		response.Pod.Containers = append(response.Pod.Containers, containerMetrics)
	}

	response.Pod.CpuRaw = totalCpu
	response.Pod.MemoryRaw = totalMemory
	response.Pod.CpuHumanReadable = fmt.Sprintf("%vm", totalCpu)
	response.Pod.MemoryHumanReadable = fmt.Sprintf("%vMi", totalMemory)

	metricLog := fmt.Sprintf("name: %s | cpu: %s | memory: %s |", response.Pod.Name, response.Pod.CpuHumanReadable, response.Pod.MemoryHumanReadable)
	logging.Info(metricLog)
	return response, nil
}
