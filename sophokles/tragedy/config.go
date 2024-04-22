package tragedy

import (
	"github.com/odysseia-greek/agora/thales"
	"os"
)

const (
	DefaultAddress string = "localhost:50053"
)

func NewMetricServiceImpl() (*MetricServiceImpl, error) {
	podName := os.Getenv("POD_NAME")

	if podName == "" {
		podName = "sophokles-0"
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "odysseia"
	}

	kube, err := thales.CreateKubeClient(false)
	if err != nil {
		return nil, err
	}

	return &MetricServiceImpl{
		PodName:   podName,
		Namespace: namespace,
		Kube:      kube,
	}, nil
}
