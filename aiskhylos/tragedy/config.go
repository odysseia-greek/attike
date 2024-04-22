package tragedy

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/thales"
	"os"
	"time"
)

const (
	DefaultWaitTime string = "180"
)

func NewMetricGathererImpl(env string) (*MetricsGathererImpl, error) {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "odysseia"
	}

	kube, err := thales.CreateKubeClient(false)
	if err != nil {
		return nil, err
	}

	var interval time.Duration
	wait := os.Getenv(config.EnvWaitTime)

	if wait == "" {
		interval, _ = time.ParseDuration(DefaultWaitTime + "s")
	} else {
		interval, _ = time.ParseDuration(wait + "s")
	}

	logging.Info(fmt.Sprintf("gathering set to an interval of: %s", interval.String()))

	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	cfg = aristoteles.ElasticConfig(env, false, tls)

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	index := config.StringFromEnv(config.EnvIndex, config.MetricsElasticIndex)

	return &MetricsGathererImpl{
		Namespace: namespace,
		Index:     index,
		Kube:      kube,
		Interval:  interval,
		Elastic:   elastic,
	}, nil
}
