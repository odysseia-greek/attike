package comedy

import (
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"os"
	"time"
)

const (
	DefaultAddress string = "localhost:50052"
)

func NewTraceServiceImpl(env string) (*TraceServiceImpl, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}

	podName := os.Getenv("POD_NAME")

	if podName == "" {
		podName = "aristophanes-0"
	}
	// Get the Namespace from the environment
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "odysseia"
	}

	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	cfg = aristoteles.ElasticConfig(env, testOverWrite, tls)

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if healthCheck {
		err := aristoteles.HealthCheck(elastic)
		if err != nil {
			return nil, err
		}
	}

	index := config.StringFromEnv(config.EnvIndex, config.TracingElasticIndex)
	startTimeMap := make(map[string]time.Time)

	return &TraceServiceImpl{
		StartTimeMap: startTimeMap,
		PodName:      podName,
		Namespace:    namespace,
		Elastic:      elastic,
		Index:        index,
	}, nil
}
