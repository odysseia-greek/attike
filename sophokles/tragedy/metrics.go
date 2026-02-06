package tragedy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
)

const (
	saTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	saCAPath    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

func (c *Collector) Run(ctx context.Context) error {
	if err := c.scrapeOnce(ctx); err != nil {
		logging.Error(fmt.Sprintf("initial scrape failed: %v", err))
	}

	ticker := time.NewTicker(c.ScrapeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logging.System("Sophokles shutting down.")
			return nil
		case <-ticker.C:
			if err := c.scrapeOnce(ctx); err != nil {
				logging.Error(fmt.Sprintf("scrape failed: %v", err))
			}
		}
	}
}

func (c *Collector) scrapeOnce(parent context.Context) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()

	summary, err := c.fetchStatsSummary(ctx)
	if err != nil {
		return err
	}
	nodeCPU := int64(0)
	if summary.Node.CPU.UsageNanoCores != nil {
		nodeCPU = int64(*summary.Node.CPU.UsageNanoCores / 1_000_000)
	}
	nodeMem := int64(0)
	if summary.Node.Memory.WorkingSetBytes != nil {
		nodeMem = int64(*summary.Node.Memory.WorkingSetBytes)
	}
	nodeTotals := &ResourceTotals{CPUMillicores: nodeCPU, MemoryBytes: nodeMem}

	// Emit one MetricSample per pod
	now := time.Now().UTC()
	for _, p := range summary.Pods {
		ns := p.PodRef.Namespace
		if _, skip := c.ExcludeNamespaces[ns]; skip {
			continue
		}

		var totalcpuMcores int64
		var totalmemBytes int64
		containers := make([]ContainerSample, 0, len(p.Containers))

		for _, ctr := range p.Containers {
			// Prefer instantaneous usageNanoCores (if present).
			cpuMcores := int64(0)
			if ctr.CPU.UsageNanoCores != nil {
				// 1 millicore = 1e6 nanocores
				cpuMcores = int64(*ctr.CPU.UsageNanoCores / 1_000_000)
			}

			memBytes := int64(0)
			if ctr.Memory.WorkingSetBytes != nil {
				memBytes = int64(*ctr.Memory.WorkingSetBytes)
			}

			totalcpuMcores += cpuMcores
			totalmemBytes += memBytes

			containers = append(containers, ContainerSample{
				Name: ctr.Name,
				Totals: ResourceTotals{
					CPUMillicores: cpuMcores,
					MemoryBytes:   memBytes,
				},
			})
		}

		sample := MetricSample{
			SchemaVersion: 1,
			Timestamp:     now,
			Node:          NodeRef{Name: c.NodeName},
			NodeTotals:    nodeTotals,
			Workload: WorkloadRef{
				Namespace: ns,
				PodName:   p.PodRef.Name,
				PodUID:    p.PodRef.UID,
			},
			PodTotals: ResourceTotals{
				CPUMillicores: totalcpuMcores,
				MemoryBytes:   totalmemBytes,
			},
			Containers: containers,
			Source: SourceMeta{
				Collector:        "sophokles",
				Method:           "kubelet_stats_summary_via_apiserver_proxy",
				ScrapeDurationMS: time.Since(start).Milliseconds(),
			},
		}

		b, err := json.Marshal(sample)
		if err != nil {
			logging.Error(fmt.Sprintf("marshal sample failed for %s/%s: %v", ns, p.PodRef.Name, err))
			continue
		}
		err = c.enqueueTask(ctx, string(b))
		if err != nil {
			logging.Error(fmt.Sprintf("queue failed %s/%s: %v", ns, p.PodRef.Name, err))
			continue
		}
	}

	return nil
}

func (c *Collector) fetchStatsSummary(ctx context.Context) (*StatsSummary, error) {
	url := fmt.Sprintf("%s/api/v1/nodes/%s/proxy/stats/summary", c.APIServerURL, c.NodeName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+string(c.BearerToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET stats summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("stats summary non-2xx: %s body=%s", resp.Status, string(body))
	}

	var summary StatsSummary
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&summary); err != nil {
		return nil, fmt.Errorf("decode stats summary: %w", err)
	}

	return &summary, nil
}
