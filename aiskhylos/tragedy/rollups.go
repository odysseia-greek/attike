package tragedy

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
)

type RollupConfig struct {
	Index      string        // "metrics-rollup" alias
	Bucket     time.Duration // 10 * time.Minute
	FlushEvery time.Duration // how often to check for closed windows (e.g. 1m)
	MaxBytes   int           // e.g. 5_000_000
	MaxDocs    int           // optional, e.g. 2000
	Schema     int           // schema_version
	Aggregator string        // e.g. "aiskhylos"
	WindowKind string        // e.g. "fixed"
}

type rollupKey struct {
	Node string
	NS   string
	UID  string
	Name string
}

type windowKey struct {
	BucketEnd time.Time
	Key       rollupKey
}

type windowAgg struct {
	cpu []int64
	mem []int64

	collector string
	method    string
}

func (g *GathererImpl) StartRollups(ctx context.Context, rollCfg RollupConfig, samples <-chan MetricSampleDoc) {
	if rollCfg.Bucket == 0 {
		rollCfg.Bucket = 10 * time.Minute
	}
	if rollCfg.FlushEvery == 0 {
		rollCfg.FlushEvery = 1 * time.Minute
	}
	if rollCfg.MaxBytes == 0 {
		rollCfg.MaxBytes = 5_000_000
	}
	if rollCfg.MaxDocs == 0 {
		rollCfg.MaxDocs = 2000
	}
	if rollCfg.Schema == 0 {
		rollCfg.Schema = 1
	}
	if rollCfg.Aggregator == "" {
		rollCfg.Aggregator = "aiskhylos"
	}
	if rollCfg.WindowKind == "" {
		rollCfg.WindowKind = "fixed"
	}
	if rollCfg.Index == "" {
		rollCfg.Index = g.MetricsRollupIndex
	}

	go g.rollupLoop(ctx, rollCfg, samples)
}

func bucketEnd(ts time.Time, bucket time.Duration) time.Time {
	start := ts.Truncate(bucket)
	return start.Add(bucket)
}

// closedWindowsUpTo returns the time threshold such that any window with BucketEnd <= threshold is closed.
// We intentionally do NOT flush the current in-progress bucket.
// Example: now=12:03, bucket=10m -> current bucket is 12:00-12:10, so flushUpTo=12:00
func closedWindowsUpTo(now time.Time, bucket time.Duration) time.Time {
	return now.Truncate(bucket)
}

func (g *GathererImpl) rollupLoop(ctx context.Context, cfg RollupConfig, samples <-chan MetricSampleDoc) {
	ticker := time.NewTicker(cfg.FlushEvery)
	defer ticker.Stop()

	aggs := make(map[windowKey]*windowAgg)

	flushClosed := func(now time.Time) {
		flushUpTo := closedWindowsUpTo(now, cfg.Bucket)

		toFlush := make([]windowKey, 0, 128)
		for wk := range aggs {
			if !wk.BucketEnd.After(flushUpTo) {
				toFlush = append(toFlush, wk)
			}
		}
		if len(toFlush) == 0 {
			return
		}

		for _, wk := range toFlush {
			agg := aggs[wk]
			delete(aggs, wk)

			if agg == nil || len(agg.cpu) == 0 {
				continue
			}

			doc := MetricsRollupDoc{
				SchemaVersion: cfg.Schema,
				Timestamp:     wk.BucketEnd,
				BucketMS:      int64(cfg.Bucket / time.Millisecond),
				Count:         int64(len(agg.cpu)),
				Node:          MetricNodeRef{Name: wk.Key.Node},
				Workload: WorkloadRef{
					Namespace: wk.Key.NS,
					PodName:   wk.Key.Name,
					PodUID:    wk.Key.UID,
				},
				CPUMCores: RollupStatsNumeric{
					Avg: avgI64(agg.cpu),
					Max: maxI64(agg.cpu),
					P95: pctl95I64(agg.cpu),
				},
				MemBytes: RollupStatsNumeric{
					Avg: avgI64(agg.mem),
					Max: maxI64(agg.mem),
					P95: pctl95I64(agg.mem),
				},
				Source: RollupSource{
					Aggregator:    cfg.Aggregator,
					Collector:     agg.collector,
					Method:        agg.method,
					WindowKind:    cfg.WindowKind,
					WindowSamples: int64(len(agg.cpu)),
				},
			}

			body, err := json.Marshal(doc)
			if err != nil {
				logging.Warn("marshal rollup failed: " + err.Error())
				continue
			}

			// deterministic + idempotent
			id := fmt.Sprintf("%d:%s:%s:%s",
				wk.BucketEnd.UnixMilli(),
				wk.Key.Node,
				wk.Key.NS,
				wk.Key.UID,
			)

			_, err = g.Elastic.Document().CreateWithId(cfg.Index, id, body)
			if err != nil {
				logging.Error("rollup index failed: " + err.Error())
			}

			logging.Info("rollup doc created: " + id)
		}
	}

	for {
		select {
		case <-ctx.Done():
			flushClosed(time.Now())
			return

		case <-ticker.C:
			flushClosed(time.Now())

		case s, ok := <-samples:
			if !ok {
				return
			}

			wk := windowKey{
				BucketEnd: bucketEnd(s.Timestamp, cfg.Bucket),
				Key: rollupKey{
					Node: s.Node.Name,
					NS:   s.Workload.Namespace,
					UID:  s.Workload.PodUID,
					Name: s.Workload.PodName,
				},
			}

			agg := aggs[wk]
			if agg == nil {
				agg = &windowAgg{
					cpu:       make([]int64, 0, 64),
					mem:       make([]int64, 0, 64),
					collector: s.Source.Collector,
					method:    s.Source.Method,
				}
				aggs[wk] = agg
			}

			agg.cpu = append(agg.cpu, s.PodTotals.CPUMCores)
			agg.mem = append(agg.mem, s.PodTotals.MemBytes)
		}
	}
}

func avgI64(v []int64) float64 {
	if len(v) == 0 {
		return 0
	}
	var sum int64
	for _, x := range v {
		sum += x
	}
	return float64(sum) / float64(len(v))
}

func maxI64(v []int64) int64 {
	if len(v) == 0 {
		return 0
	}
	m := v[0]
	for _, x := range v[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func pctl95I64(v []int64) float64 {
	n := len(v)
	if n == 0 {
		return 0
	}
	cp := make([]int64, n)
	copy(cp, v)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })

	// Nearest-rank percentile (ceil(0.95*n))
	rank := int((95*n + 99) / 100)
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return float64(cp[rank-1])
}
