package pneuma

import (
	"fmt"
	"strings"

	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
)

type AnaximenesHandler struct {
	Indices                 []string
	MaxAge                  string
	PolicyNameTrace         string
	PolicyNameMetrics       string
	PolicyNameMetricsRollup string
	Elastic                 elastic.Client
	Ambassador              *diplomat.ClientAmbassador
}

func (a *AnaximenesHandler) CreateAttikeIndices() {
	for _, index := range a.Indices {
		deleted, err := a.Elastic.Index().Delete(index)

		logging.Info(fmt.Sprintf("delete response for %s: success=%v, err=%v",
			index, deleted, err))

		if err != nil {
			if strings.Contains(err.Error(), "index_not_found_exception") {
				logging.Info(fmt.Sprintf("index %s not found, creating it", index))
				err = a.createIndexAtStartup(index)
				if err != nil {
					logging.Error(fmt.Sprintf("failed to create index: %v", err))
				}
				continue
			}

			// Log any other error in detail
			logging.Debug(fmt.Sprintf("error deleting index %s: %v", index, err))
		}

		if deleted {
			logging.Info(fmt.Sprintf("recreating index %s after deletion", index))
			err = a.createIndexAtStartup(index)
			if err != nil {
				logging.Error(fmt.Sprintf("failed to recreate index: %v", err))
			}
		}
	}
}

func (a *AnaximenesHandler) createIndexAtStartup(index string) error {
	request := a.createMapping(index)
	created, err := a.Elastic.Index().CreateWithAlias(index, request)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", index, created.Acknowledged))

	return nil
}

func (a *AnaximenesHandler) createMapping(indexName string) map[string]interface{} {
	switch indexName {
	case config.TracingElasticIndex:
		return a.createTraceIndexMapping(a.PolicyNameTrace)
	case config.MetricsElasticIndex:
		return a.createMetricSamplesIndexMapping(a.PolicyNameMetrics)
	case config.MetricsRollupElasticIndex:
		return a.createMetricsRollupIndexMapping(a.PolicyNameMetricsRollup)
	}

	return nil
}

func (a *AnaximenesHandler) createTraceIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"items": map[string]interface{}{
					"type": "nested", // Use "nested" type for arrays of complex objects
				},
				"isActive": map[string]interface{}{
					"type": "boolean",
				},
				"timeStarted": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ss.SSS",
				},
				"timeEnded": map[string]interface{}{
					"type":       "date",
					"format":     "yyyy-MM-dd'T'HH:mm:ss.SSS",
					"null_value": "1970-01-01T00:00:00.000",
				},
				"totalTime": map[string]interface{}{
					"type": "long",
				},
				"responseCode": map[string]interface{}{
					"type": "short",
				},
				"metrics_snapshot": map[string]interface{}{
					"properties": map[string]interface{}{
						// when the snapshot was taken
						"@timestamp": map[string]interface{}{
							"type":   "date",
							"format": "strict_date_optional_time||epoch_millis",
						},

						// pod-level usage
						"pod": map[string]interface{}{
							"properties": map[string]interface{}{
								"cpu_mcores": map[string]interface{}{"type": "long"},
								"mem_bytes":  map[string]interface{}{"type": "long"},
							},
						},

						// optional node-level context
						"node": map[string]interface{}{
							"properties": map[string]interface{}{
								"name":       map[string]interface{}{"type": "keyword"},
								"cpu_mcores": map[string]interface{}{"type": "long"},
								"mem_bytes":  map[string]interface{}{"type": "long"},
							},
						},

						// how this snapshot was produced
						"source": map[string]interface{}{
							"properties": map[string]interface{}{
								"collector": map[string]interface{}{"type": "keyword"}, // sophokles
								"method":    map[string]interface{}{"type": "keyword"}, // kubelet_stats_summary
							},
						},
					},
				},
			},
		},
		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "trace",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}

func (a *AnaximenesHandler) createMetricSamplesIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"dynamic": "strict",
			"properties": map[string]interface{}{
				"schema_version": map[string]interface{}{"type": "integer"},
				"@timestamp": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"node": map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "keyword"},
					},
				},
				"node_totals": map[string]interface{}{
					"properties": map[string]interface{}{
						"cpu_mcores": map[string]interface{}{"type": "long"},
						"mem_bytes":  map[string]interface{}{"type": "long"},
					},
				},
				"workload": map[string]interface{}{
					"properties": map[string]interface{}{
						"namespace": map[string]interface{}{"type": "keyword"},
						"pod_name":  map[string]interface{}{"type": "keyword"},
						"pod_uid":   map[string]interface{}{"type": "keyword"},
					},
				},
				"pod_totals": map[string]interface{}{
					"properties": map[string]interface{}{
						"cpu_mcores": map[string]interface{}{"type": "long"},
						"mem_bytes":  map[string]interface{}{"type": "long"},
					},
				},
				"containers": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "keyword"},
						"totals": map[string]interface{}{
							"properties": map[string]interface{}{
								"cpu_mcores": map[string]interface{}{"type": "long"},
								"mem_bytes":  map[string]interface{}{"type": "long"},
							},
						},
					},
				},
				"source": map[string]interface{}{
					"properties": map[string]interface{}{
						"collector":          map[string]interface{}{"type": "keyword"},
						"method":             map[string]interface{}{"type": "keyword"},
						"scrape_duration_ms": map[string]interface{}{"type": "long"},
					},
				},
			},
		},
		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "metric",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}

func (a *AnaximenesHandler) createMetricsRollupIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"dynamic": "strict",
			"properties": map[string]interface{}{
				"schema_version": map[string]interface{}{"type": "integer"},

				"@timestamp": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},

				"bucket_ms": map[string]interface{}{"type": "long"},
				"count":     map[string]interface{}{"type": "long"},

				"node": map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "keyword"},
					},
				},

				"workload": map[string]interface{}{
					"properties": map[string]interface{}{
						"namespace": map[string]interface{}{"type": "keyword"},
						"pod_name":  map[string]interface{}{"type": "keyword"},
						"pod_uid":   map[string]interface{}{"type": "keyword"},
					},
				},

				"cpu_mcores": map[string]interface{}{
					"properties": map[string]interface{}{
						"avg": map[string]interface{}{"type": "double"},
						"max": map[string]interface{}{"type": "long"},
						"p95": map[string]interface{}{"type": "double"},
					},
				},

				"mem_bytes": map[string]interface{}{
					"properties": map[string]interface{}{
						"avg": map[string]interface{}{"type": "double"},
						"max": map[string]interface{}{"type": "long"},
						"p95": map[string]interface{}{"type": "double"},
					},
				},

				"source": map[string]interface{}{
					"properties": map[string]interface{}{
						"aggregator":     map[string]interface{}{"type": "keyword"},
						"collector":      map[string]interface{}{"type": "keyword"},
						"method":         map[string]interface{}{"type": "keyword"},
						"window_kind":    map[string]interface{}{"type": "keyword"},
						"window_samples": map[string]interface{}{"type": "long"},
					},
				},
			},
		},

		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "metrics-rollup",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}
