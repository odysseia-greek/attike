package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/odysseia-greek/attike/euripides/elastic"
	"github.com/odysseia-greek/attike/euripides/graph/model"
	"github.com/odysseia-greek/attike/euripides/models"
)

func (e *EuripidesHandler) MetricsSummary(ctx context.Context, input *model.MetricsSummaryInput) (*model.MetricsSummary, error) {
	summary := &model.MetricsSummary{}

	query := elastic.BuildMetricsSummary(input)

	res, err := e.Elastic.Query().Match(e.MetricsRollupIndex, query)
	if err != nil {
		return nil, err
	}
	if len(res.Hits.Hits) == 0 {
		return summary, nil
	}

	docs := make([]models.RollupDoc, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		b, _ := json.Marshal(hit.Source)

		var d models.RollupDoc
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}

	// Aggregate
	nodeAgg := map[string]*metricAgg{}
	nsAgg := map[string]*metricAgg{}
	podAgg := map[string]*metricAgg{}

	// optional: keep labels for pods (node+ns+pod)
	type podLabel struct{ node, ns, pod string }
	podLabels := map[string]podLabel{}

	for _, d := range docs {
		node := d.Node.Name
		ns := d.Workload.Namespace
		pod := d.Workload.PodName

		// nodes
		if nodeAgg[node] == nil {
			nodeAgg[node] = &metricAgg{}
		}
		nodeAgg[node].add(d)

		// namespaces
		if nsAgg[ns] == nil {
			nsAgg[ns] = &metricAgg{}
		}
		nsAgg[ns].add(d)

		pk := ns + "/" + pod
		if podAgg[pk] == nil {
			podAgg[pk] = &metricAgg{}
			podLabels[pk] = podLabel{node: node, ns: ns, pod: pod}
		}
		podAgg[pk].add(d)
	}

	// Convert to GraphQL models
	summary.Nodes = &model.MetricsGroupNode{}
	for node, a := range nodeAgg {
		sampleCount := int32(a.sampleCount)
		item := &model.NodeMetricsAgg{
			Node:        node,
			DocCount:    int32(a.docCount),
			SampleCount: &sampleCount,
			CPU: &model.MetricStats{
				Avg:      a.cpuAvg.valueAvg(),
				Max:      a.cpuMax.valueMax(),
				P95:      a.cpuP95.valueAvg(),
				TotalMax: a.totalMaxCPUVal(),
			},
			Mem: &model.MetricStats{
				Avg:      a.memAvg.valueAvg(),
				Max:      a.memMax.valueMax(),
				P95:      a.memP95.valueAvg(),
				TotalMax: a.totalMaxMemVal(),
			},
		}
		withHumanStats(item.Mem, bytesHuman)
		withHumanStats(item.CPU, mcoresHuman)
		withHumanTotalMax(item.Mem, bytesHuman)
		withHumanTotalMax(item.CPU, mcoresHuman)

		summary.Nodes.Items = append(summary.Nodes.Items, item)
	}

	sortByMemMaxDesc(summary.Nodes.Items,
		func(n *model.NodeMetricsAgg) *float64 {
			if n == nil || n.Mem == nil {
				return nil
			}
			return n.Mem.TotalMax
		},
		func(n *model.NodeMetricsAgg) string {
			if n == nil {
				return ""
			}
			return n.Node
		},
	)

	summary.Nodes.Total = int32(len(summary.Nodes.Items))

	summary.Namespaces = &model.MetricsGroupNamespace{}
	for ns, a := range nsAgg {
		sampleCount := int32(a.sampleCount)
		item := &model.NamespaceMetricsAgg{
			Namespace:   ns,
			DocCount:    int32(a.docCount),
			SampleCount: &sampleCount,
			CPU: &model.MetricStats{
				Avg:      a.cpuAvg.valueAvg(),
				Max:      a.cpuMax.valueMax(),
				P95:      a.cpuP95.valueAvg(),
				TotalMax: a.totalMaxCPUVal(),
			},
			Mem: &model.MetricStats{
				Avg:      a.memAvg.valueAvg(),
				Max:      a.memMax.valueMax(),
				P95:      a.memP95.valueAvg(),
				TotalMax: a.totalMaxMemVal(),
			},
		}
		withHumanStats(item.Mem, bytesHuman)
		withHumanStats(item.CPU, mcoresHuman)
		withHumanTotalMax(item.Mem, bytesHuman)
		withHumanTotalMax(item.CPU, mcoresHuman)

		summary.Namespaces.Items = append(summary.Namespaces.Items, item)
	}

	sortByMemMaxDesc(summary.Namespaces.Items,
		func(n *model.NamespaceMetricsAgg) *float64 {
			if n == nil || n.Mem == nil {
				return nil
			}
			return n.Mem.TotalMax
		},
		func(n *model.NamespaceMetricsAgg) string {
			if n == nil {
				return ""
			}
			return n.Namespace
		},
	)

	summary.Namespaces.Total = int32(len(summary.Namespaces.Items))

	summary.Pods = &model.MetricsGroupPod{}
	for pk, a := range podAgg {
		lbl := podLabels[pk]
		sampleCount := int32(a.sampleCount)
		item := &model.PodMetricsAgg{
			PodName:     lbl.pod,
			Namespace:   &lbl.ns,
			Node:        &lbl.node,
			DocCount:    int32(a.docCount),
			SampleCount: &sampleCount,
			CPU: &model.MetricStats{
				Avg: a.cpuAvg.valueAvg(),
				Max: a.cpuMax.valueMax(),
				P95: a.cpuP95.valueAvg(),
			},
			Mem: &model.MetricStats{
				Avg: a.memAvg.valueAvg(),
				Max: a.memMax.valueMax(),
				P95: a.memP95.valueAvg(),
			},
		}
		withHumanStats(item.Mem, bytesHuman)
		withHumanStats(item.CPU, mcoresHuman)

		summary.Pods.Items = append(summary.Pods.Items, item)
	}
	sortByMemMaxDesc(summary.Pods.Items,
		func(p *model.PodMetricsAgg) *float64 {
			if p == nil || p.Mem == nil {
				return nil
			}
			return p.Mem.Max
		},
		func(p *model.PodMetricsAgg) string {
			if p == nil {
				return ""
			}
			ns := ""
			if p.Namespace != nil {
				ns = *p.Namespace
			}
			return ns + "/" + p.PodName
		},
	)

	summary.Pods.Total = int32(len(summary.Pods.Items))

	return summary, nil
}

type statsAgg struct {
	// weighted sums
	wSum float64
	w    float64

	// max tracking
	max float64
	has bool
}

func (a *statsAgg) addAvg(value float64, weight int) {
	if weight <= 0 {
		weight = 1
	}
	a.wSum += value * float64(weight)
	a.w += float64(weight)
	a.has = true
}

func (a *statsAgg) addMax(value float64) {
	if !a.has || value > a.max {
		a.max = value
	}
	a.has = true
}

func (a *statsAgg) valueAvg() *float64 {
	if a.w == 0 {
		return nil
	}
	v := a.wSum / a.w
	return &v
}

func (a *statsAgg) valueMax() *float64 {
	if !a.has {
		return nil
	}
	v := a.max
	return &v
}

type bucketTotals struct {
	cpuSum float64
	memSum float64
}

type metricAgg struct {
	docCount    int
	sampleCount int

	cpuAvg statsAgg
	cpuMax statsAgg
	cpuP95 statsAgg

	memAvg statsAgg
	memMax statsAgg
	memP95 statsAgg

	// NEW: totals across pods per time bucket
	byBucket    map[int64]*bucketTotals
	totalMaxCPU float64
	totalMaxMem float64
	hasTotalMax bool
}

func (m *metricAgg) add(d models.RollupDoc) {
	m.docCount++
	m.sampleCount += d.Count

	// per-doc aggregates (pod-level within the group)
	m.cpuAvg.addAvg(d.CpuMcores.Avg, d.Count)
	m.cpuP95.addAvg(d.CpuMcores.P95, d.Count)
	m.cpuMax.addMax(d.CpuMcores.Max)

	m.memAvg.addAvg(d.MemBytes.Avg, d.Count)
	m.memP95.addAvg(d.MemBytes.P95, d.Count)
	m.memMax.addMax(d.MemBytes.Max)

	// totals per bucket (group-level totals)
	if m.byBucket == nil {
		m.byBucket = make(map[int64]*bucketTotals)
	}

	// rollup docs are already bucketed; use timestamp as bucket key
	k := d.Timestamp.UnixMilli()

	bt := m.byBucket[k]
	if bt == nil {
		bt = &bucketTotals{}
		m.byBucket[k] = bt
	}

	// Choose what "usage in bucket" means.
	// For totals I'd use AVG (stable). If you want "peaky", sum Max instead.
	bt.cpuSum += d.CpuMcores.Avg
	bt.memSum += d.MemBytes.Avg

	if !m.hasTotalMax || bt.cpuSum > m.totalMaxCPU {
		m.totalMaxCPU = bt.cpuSum
	}
	if !m.hasTotalMax || bt.memSum > m.totalMaxMem {
		m.totalMaxMem = bt.memSum
	}
	m.hasTotalMax = true
}

func sortByMemMaxDesc[T any](items []*T, memMax func(*T) *float64, tie func(*T) string) {
	sort.Slice(items, func(i, j int) bool {
		mi := 0.0
		mj := 0.0

		if items[i] != nil {
			if v := memMax(items[i]); v != nil {
				mi = *v
			}
		}
		if items[j] != nil {
			if v := memMax(items[j]); v != nil {
				mj = *v
			}
		}

		if mi == mj {
			return tie(items[i]) < tie(items[j])
		}
		return mi > mj
	})
}

func bytesHuman(b float64) string {
	if b < 0 {
		return ""
	}
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case b >= TB:
		return fmt.Sprintf("%.2f TiB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2f GiB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MiB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KiB", b/KB)
	default:
		return fmt.Sprintf("%.0f B", b)
	}
}

func mcoresHuman(m float64) string {
	if m < 0 {
		return ""
	}
	// Nice display:
	// - show mCPU if < 1000
	// - show cores if >= 1000
	if m >= 1000 {
		return fmt.Sprintf("%.2f cores", m/1000.0)
	}
	// keep some precision for tiny numbers
	if m < 10 {
		return fmt.Sprintf("%.2f mCPU", m)
	}
	return fmt.Sprintf("%.0f mCPU", m)
}

func withHumanStats(s *model.MetricStats, format func(float64) string) {
	if s == nil || format == nil {
		return
	}
	if s.Avg != nil {
		s.AvgHuman = ptr(format(*s.Avg))
	}
	if s.Max != nil {
		s.MaxHuman = ptr(format(*s.Max))
	}
	if s.P95 != nil {
		s.P95Human = ptr(format(*s.P95))
	}
}

func ptr[T any](v T) *T { return &v }

func (m *metricAgg) totalMaxCPUVal() *float64 {
	if !m.hasTotalMax {
		return nil
	}
	v := m.totalMaxCPU
	return &v
}

func (m *metricAgg) totalMaxMemVal() *float64 {
	if !m.hasTotalMax {
		return nil
	}
	v := m.totalMaxMem
	return &v
}

func withHumanTotalMax(s *model.MetricStats, format func(float64) string) {
	if s == nil || format == nil {
		return
	}
	if s.TotalMax != nil {
		s.TotalMaxHuman = ptr(format(*s.TotalMax))
	}
}
