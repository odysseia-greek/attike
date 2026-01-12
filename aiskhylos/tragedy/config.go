package tragedy

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
)

func NewGatherer() (*GathererImpl, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)
	if err != nil {
		return nil, err
	}

	healthy := ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	ambassadorCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	vaultConfig, err := ambassador.GetSecret(ambassadorCtx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	elasticService := aristoteles.ElasticService(tls)

	cfg = models.Config{
		Service:     elasticService,
		Username:    vaultConfig.ElasticUsername,
		Password:    vaultConfig.ElasticPassword,
		ElasticCERT: vaultConfig.ElasticCERT,
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	var traceIndex string
	var metricsIndex string
	var metricsRollupIndex string

	rootIndex := config.StringFromEnv(config.EnvIndex, "")
	splitIndex := strings.Split(rootIndex, ";")
	for _, index := range splitIndex {
		switch index {
		case config.TracingElasticIndex:
			traceIndex = index
		case config.MetricsElasticIndex:
			metricsIndex = index
		case "metrics_rollup":
			metricsRollupIndex = index
		}
	}

	eupalinosAddress := config.StringFromEnv("EUPALINOS_TRACING_SERVICE", config.DefaultEupalinosService)
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

	logging.Debug("queue healthy starting up")

	metricsChannel := config.StringFromEnv("METRICS_CHANNEL", "sophokles")
	traceChannel := config.StringFromEnv("TRACE_CHANNEL", "aristophanes")

	return &GathererImpl{
		MetricsIndex:       metricsIndex,
		MetricsRollupIndex: metricsRollupIndex,
		TraceIndex:         traceIndex,
		Elastic:            elastic,
		Eupalinos:          queue,
		TraceChannel:       traceChannel,
		MetricsChannel:     metricsChannel,
	}, nil
}
