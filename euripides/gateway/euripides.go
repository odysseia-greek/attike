package gateway

import (
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
)

type EuripidesHandler struct {
	MetricsIndex       string
	MetricsRollupIndex string
	TraceIndex         string
	Elastic            aristoteles.Client

	Eupalinos      *stomion.QueueClient
	TraceIdChannel string
}
