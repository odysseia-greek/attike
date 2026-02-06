package gateway

import (
	"sync"
	"time"

	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/attike/euripides/models"
)

type TraceReportReaderConfig struct {
	PollEvery       time.Duration
	DequeueWait     time.Duration
	MaxDrainPerPoll int
}

type EuripidesHandler struct {
	MetricsIndex       string
	MetricsRollupIndex string
	TraceIndex         string
	Version            string
	Environment        string
	Elastic            aristoteles.Client

	Eupalinos     *stomion.QueueClient
	ReportChannel string

	mu        sync.Mutex
	traceCfg  TraceReportReaderConfig
	latest    map[string]models.TraceRootSource // last known per traceID
	pendingQ  []string                          // FIFO of traceIDs that changed
	pendingIn map[string]struct{}               // de-dupe so each ID only once in pendingQ
}
