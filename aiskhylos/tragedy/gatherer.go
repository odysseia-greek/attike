package tragedy

import (
	"context"
	"time"

	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
	pb "github.com/odysseia-greek/agora/eupalinos/v1"
)

type GathererImpl struct {
	MetricsIndex       string
	MetricsRollupIndex string
	TraceIndex         string
	Elastic            aristoteles.Client

	Eupalinos      *stomion.QueueClient
	TraceChannel   string
	MetricsChannel string
	ReportChannel  string

	MetricCfg     MetricBulkConfig
	metricCh      chan *pb.Epistello   // raw queue messages
	metricSamples chan MetricSampleDoc // parsed docs for raw indexing
	rollupSamples chan MetricSampleDoc // parsed docs for rollups
}

func (g *GathererImpl) Start(ctx context.Context) {
	go g.StartMetrics(ctx)
	go g.StartTraces(ctx, TraceConfig{
		Index:       g.TraceIndex,
		GracePeriod: 5 * time.Second,
		MaxOpenAge:  15 * time.Minute,
	})
	<-ctx.Done()
}
