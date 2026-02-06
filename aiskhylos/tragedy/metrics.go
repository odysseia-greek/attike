package tragedy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
)

type MetricBulkConfig struct {
	Index       string        // e.g. "metric" (rollover alias) or concrete index
	FlushEvery  time.Duration // e.g. 1 * time.Minute
	MaxDocs     int           // e.g. 2000
	MaxBytes    int           // e.g. 5_000_000
	IdleBackoff time.Duration // e.g. 100 * time.Millisecond (when queue empty/errors)
}

func (g *GathererImpl) StartMetrics(ctx context.Context) {
	// Defaults
	if g.MetricCfg.FlushEvery == 0 {
		g.MetricCfg.FlushEvery = 1 * time.Minute
	}
	if g.MetricCfg.MaxDocs == 0 {
		g.MetricCfg.MaxDocs = 2000
	}
	if g.MetricCfg.MaxBytes == 0 {
		g.MetricCfg.MaxBytes = 5_000_000
	}
	if g.MetricCfg.IdleBackoff == 0 {
		g.MetricCfg.IdleBackoff = 100 * time.Millisecond
	}
	if g.MetricCfg.Index == "" {
		g.MetricCfg.Index = g.MetricsIndex
	}

	// Channels
	if g.metricCh == nil {
		g.metricCh = make(chan *pb.Epistello, 10_000)
	}
	if g.metricSamples == nil {
		g.metricSamples = make(chan MetricSampleDoc, 10_000)
	}
	if g.rollupSamples == nil {
		g.rollupSamples = make(chan MetricSampleDoc, 10_000)
	}

	// Start rollups (consumes parsed docs)
	g.StartRollups(ctx, RollupConfig{
		Index:      g.MetricsRollupIndex,
		Bucket:     10 * time.Minute,
		FlushEvery: 1 * time.Minute,
		Aggregator: "aiskhylos",
		WindowKind: "fixed",
	}, g.rollupSamples)

	// Stages
	go g.consumeMetricQueue(ctx)
	go g.decodeAndFanoutMetrics(ctx)
	go g.bulkIndexMetrics(ctx)

	<-ctx.Done()
}

func (g *GathererImpl) consumeMetricQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := g.Eupalinos.DequeueMessage(ctx, &pb.ChannelInfo{Name: g.MetricsChannel})
		if err != nil {
			time.Sleep(g.MetricCfg.IdleBackoff)
			continue
		}

		select {
		case g.metricCh <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func (g *GathererImpl) decodeAndFanoutMetrics(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-g.metricCh:
			parsed, err := ParseMetricSample(msg)
			if err != nil {
				logging.Warn("failed to parse metric sample: " + err.Error())
				continue
			}

			// Raw metrics: do not drop; apply backpressure if needed
			select {
			case g.metricSamples <- parsed:
			case <-ctx.Done():
				return
			}

			// Rollups: OK to drop if rollup loop is slow (tune later)
			select {
			case g.rollupSamples <- parsed:
			default:
			}
		}
	}
}

func (g *GathererImpl) bulkIndexMetrics(ctx context.Context) {
	ticker := time.NewTicker(g.MetricCfg.FlushEvery)
	defer ticker.Stop()

	var (
		buf      bytes.Buffer
		docCount int
		wg       sync.WaitGroup
	)

	reset := func() {
		buf.Reset()
		docCount = 0
	}

	flush := func() {
		if docCount == 0 {
			return
		}

		wg.Wait()

		_, err := g.Elastic.Document().Bulk(buf, g.MetricCfg.Index)
		if err != nil {
			logging.Error("metrics bulk failed: " + err.Error())
		}

		reset()
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return

		case <-ticker.C:
			flush()

		case parsed := <-g.metricSamples:
			jsonified, err := json.Marshal(parsed)
			if err != nil {
				logging.Warn("failed to marshal metric sample: " + err.Error())
				continue
			}

			meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))

			buf.Grow(len(meta) + len(jsonified) + 1)
			buf.Write(meta)
			buf.Write(jsonified)
			buf.WriteByte('\n')
			docCount++

			if docCount >= g.MetricCfg.MaxDocs || buf.Len() >= g.MetricCfg.MaxBytes {
				logging.Debug(fmt.Sprintf("flushing metrics bulk (%d docs, %d bytes):", docCount, buf.Len()))
				flush()
			}
		}
	}
}
