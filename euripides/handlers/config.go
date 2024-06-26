package handlers

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"os"
	"strings"
	"time"
)

// CreateNewConfig creates a new configuration based on the provided environment.
//
// The function performs the following steps:
//  1. Determines whether health check should be enabled based on the environment.
//  2. Retrieves configuration values from the Vault if health check is enabled.
//  3. Initializes the Elasticsearch service based on the TLS setting.
//  4. Creates the configuration object with the Elasticsearch service, username, password, and certificate.
//     - If health check is disabled, the configuration is created using the environment, testOverWrite, and TLS settings.
//     - If health check is enabled, the configuration is created using the retrieved Vault configuration values.
//  5. Creates a new Elasticsearch client using the configuration.
//  6. Performs a health check on the Elasticsearch client if health check is enabled.
//  7. Retrieves the index name from the environment variables or uses the default index name.
//  8. Returns the created configuration with the Elasticsearch client and index name.
//
// Parameters:
//   - env: The environment name (e.g., "LOCAL", "TEST").
//
// Returns:
//   - *Config: The created configuration containing the Elasticsearch client and index name.
//   - error: An error if any occurred during the configuration creation process.
func CreateNewConfig(env string) (*EuripidesHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador := diplomat.NewClientAmbassador()

	healthy := ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	vaultConfig, err := ambassador.GetSecret(ctx, &pb.VaultRequest{})
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

	index := config.StringFromEnv(config.EnvIndex, "tracing;metrics")
	splitIndex := strings.Split(index, ";")
	if len(splitIndex) == 1 {
		return nil, fmt.Errorf("cannot find second index")
	}

	return &EuripidesHandler{
		TracingIndex: splitIndex[0],
		MetricsIndex: splitIndex[1],
		Elastic:      elastic,
	}, nil
}
