package comedy

import (
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	sophokles "github.com/odysseia-greek/attike/sophokles/tragedy"
	"os"
	"strconv"
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

	gatherMetricsString := os.Getenv("GATHER_METRICS")
	gatherMetrics, err := strconv.ParseBool(gatherMetricsString)
	if err != nil {
		logging.Error(err.Error())
	}

	var metrics *sophokles.ClientMetrics
	if gatherMetrics {
		metrics = sophokles.NewClientTracer()
		if healthCheck {
			healthy := metrics.WaitForHealthyState()
			if !healthy {
				logging.Info("metrics service not ready - restarting seems the only option")
				os.Exit(1)
			}
		}
	}

	return &TraceServiceImpl{
		StartTimeMap:  startTimeMap,
		PodName:       podName,
		Namespace:     namespace,
		Elastic:       elastic,
		Index:         index,
		Metrics:       metrics,
		GatherMetrics: gatherMetrics,
	}, nil
}
