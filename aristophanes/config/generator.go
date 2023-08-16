package config

import (
	"github.com/odysseia-greek/aristoteles"
	"github.com/odysseia-greek/aristoteles/models"
	"github.com/odysseia-greek/plato/config"
	"os"
)

func CreateNewConfig(env string) (*Config, error) {
	healthCheck := true
	if env == "LOCAL" || env == "TEST" {
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

	return &Config{
		PodName:   podName,
		Namespace: namespace,
		Elastic:   elastic,
		Index:     index,
	}, nil
}
