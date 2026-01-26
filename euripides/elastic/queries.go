package elastic

import (
	"strings"
	"time"

	"github.com/odysseia-greek/attike/euripides/graph/model"
)

func traceWindowToDuration(w model.TraceWindow) time.Duration {
	switch w {
	case model.TraceWindowM5:
		return 5 * time.Minute
	case model.TraceWindowM10:
		return 10 * time.Minute
	case model.TraceWindowM30:
		return 30 * time.Minute
	case model.TraceWindowH1:
		return 1 * time.Hour
	case model.TraceWindowH2:
		return 2 * time.Hour
	case model.TraceWindowH24:
		return 24 * time.Hour
	default:
		return 10 * time.Minute
	}
}

func metricsWindowToDuration(w model.MetricsWindow) time.Duration {
	switch w {
	case model.MetricsWindowM10:
		return 10 * time.Minute
	case model.MetricsWindowM30:
		return 30 * time.Minute
	case model.MetricsWindowH1:
		return 1 * time.Hour
	case model.MetricsWindowH2:
		return 2 * time.Hour
	case model.MetricsWindowH12:
		return 12 * time.Hour
	case model.MetricsWindowH24:
		return 24 * time.Hour
	default:
		return 10 * time.Minute
	}
}

func esMillis(t time.Time) string {
	// 2006-01-02T15:04:05.000Z07:00 gives milliseconds
	return t.UTC().Format("2006-01-02T15:04:05.000")
}

func BuildTraceSearchQuery(input *model.TraceSearchInput, limit int32) map[string]interface{} {
	if limit <= 0 {
		limit = 20
	}
	if limit > 500 {
		limit = 500
	}

	// window → duration
	win := 1 * time.Hour
	if input != nil {
		win = traceWindowToDuration(input.Window)
	}

	now := time.Now().UTC()
	from := now.Add(-win)

	filter := make([]interface{}, 0, 8)

	// time window
	filter = append(filter, map[string]interface{}{
		"range": map[string]interface{}{
			"timeStarted": map[string]interface{}{
				"gte": esMillis(from),
				"lte": esMillis(now),
			},
		},
	})

	// responseCode
	if input != nil && input.ResponseCode != nil && *input.ResponseCode != 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"responseCode": *input.ResponseCode,
			},
		})
	}

	// operation (exact match; if you want prefix/contains we can swap this to match/wildcard)
	if input != nil && input.Operation != nil {
		op := strings.TrimSpace(*input.Operation)
		if op != "" {
			// If "operation" is mapped as keyword: use operation.keyword
			// If it's plain keyword already, just "operation" is fine.
			filter = append(filter, map[string]interface{}{
				"match": map[string]interface{}{
					"operation": op,
				},
			})
		}
	}

	if input != nil && input.TimeTookGreaterThan != nil && *input.TimeTookGreaterThan > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"totalTime": map[string]interface{}{
					"gte": *input.TimeTookGreaterThan,
				},
			},
		})
	}

	return map[string]interface{}{
		"size":             int(limit),
		"track_total_hits": true,
		"sort": []interface{}{
			map[string]interface{}{
				"timeStarted": map[string]interface{}{"order": "desc"},
			},
		},
		"_source": []interface{}{
			"isActive",
			"timeStarted",
			"timeEnded",
			"totalTime",
			"responseCode",
			"operation",
			"containsDBSpan",
			"numberOfItems",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filter,
			},
		},
	}
}

func BuildMetricsSummary(input *model.MetricsSummaryInput) map[string]interface{} {
	// default window if nil or unspecified
	win := 1 * time.Hour
	if input != nil {
		win = metricsWindowToDuration(input.Window)
	}

	now := time.Now().UTC()
	from := now.Add(-win)

	filter := make([]interface{}, 0, 2)

	// Range on @timestamp (rollup docs use @timestamp)
	filter = append(filter, map[string]interface{}{
		"range": map[string]interface{}{
			"@timestamp": map[string]interface{}{
				"gte": esMillis(from),
				"lte": esMillis(now),
			},
		},
	})

	// Final ES query body
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filter,
			},
		},

		// Optional but typically useful for graphing / debugging:
		"sort": []interface{}{
			map[string]interface{}{"@timestamp": "asc"},
		},
		"size": 1000,
	}

	return body
}
