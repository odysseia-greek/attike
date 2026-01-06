package comedy

import (
	"context"
	"fmt"
	"os"

	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
)

const (
	DefaultAddress string = "localhost:50052"
)

func NewTraceServiceImpl() (*TraceServiceImpl, error) {
	podName := os.Getenv("POD_NAME")

	if podName == "" {
		podName = "aristophanes-0"
	}
	// Get the Namespace from the environment
	namespace := config.StringFromEnv(config.EnvNamespace, "attike")

	commands := make(chan MapCommand, 100)

	// Queue
	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)
	logging.Debug(fmt.Sprintf("creating new eupalinos client: %s", eupalinosAddress))
	queue, err := stomion.NewEupalinosClient(eupalinosAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	logging.Debug("waiting for queue to be ready")
	queueHealthy := queue.WaitForHealthyState()
	if !queueHealthy {
		logging.Debug("no queue that is healthy")
	}

	channel := config.StringFromEnv(config.EnvChannel, "aristophanes")
	ctx, cancel := context.WithCancel(context.Background())

	return &TraceServiceImpl{
		PodName:   podName,
		Namespace: namespace,
		commands:  commands,
		Eupalinos: queue,
		Channel:   channel,
		baseCtx:   ctx,
		cancel:    cancel,
	}, nil
}
