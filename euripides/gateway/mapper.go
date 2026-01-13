package gateway

import (
	"encoding/json"
	"time"

	"github.com/odysseia-greek/attike/euripides/graph/model"
)

type esTraceDoc struct {
	IsActive     bool          `json:"isActive"`
	TimeStarted  string        `json:"timeStarted"`
	TimeEnded    *string       `json:"timeEnded,omitempty"`
	TotalTime    int64         `json:"totalTime"`
	ResponseCode int           `json:"responseCode"`
	Items        []esTraceItem `json:"items"`
}

type esTraceItem struct {
	Timestamp    string          `json:"timestamp"`
	ItemType     string          `json:"item_type"`
	SpanID       *string         `json:"span_id,omitempty"`
	ParentSpanID *string         `json:"parent_span_id,omitempty"`
	PodName      *string         `json:"pod_name,omitempty"`
	Namespace    *string         `json:"namespace,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"` // object in your ES source
}

type esPayload struct {
	Common json.RawMessage `json:"common,omitempty"`

	TraceStart json.RawMessage `json:"traceStart,omitempty"`
	TraceHop   json.RawMessage `json:"traceHop,omitempty"`
	Graphql    json.RawMessage `json:"graphql,omitempty"`
	Action     json.RawMessage `json:"action,omitempty"`
	DbSpan     json.RawMessage `json:"dbSpan,omitempty"`
	TraceStop  json.RawMessage `json:"traceStop,omitempty"`
}

type esDbSpan struct {
	Action string          `json:"action"`
	Query  string          `json:"query"`
	Hits   json.RawMessage `json:"hits"`   // "284" or 284
	TookMs json.RawMessage `json:"tookMs"` // "4" or 4
	Target string          `json:"target,omitempty"`
	Index  string          `json:"index,omitempty"`
}

type esAction struct {
	Action string          `json:"action"`
	Status string          `json:"status"`
	TookMs json.RawMessage `json:"tookMs"`
}

func parseAttikeTime(s string) (time.Time, error) {
	// common cases you already emit
	layouts := []string{
		time.RFC3339Nano,          // 2026-01-12T19:54:17.687459278Z
		time.RFC3339,              // 2026-01-12T19:54:17Z
		"2006-01-02T15:04:05.000", // 2026-01-12T19:51:32.052
		"2006-01-02T15:04:05",     // if you ever send seconds only
	}

	var lastErr error
	for _, l := range layouts {
		t, err := time.Parse(l, s)
		if err == nil {
			return t.UTC(), nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func toModelTrace(id string, d esTraceDoc) (*model.Trace, error) {
	start, err := parseAttikeTime(d.TimeStarted)
	if err != nil {
		return nil, err
	}

	var end *time.Time
	if d.TimeEnded != nil && *d.TimeEnded != "" {
		te, err := parseAttikeTime(*d.TimeEnded)
		if err != nil {
			return nil, err
		}
		end = &te
	}

	items := make([]*model.TraceItem, 0, len(d.Items))
	hasDb := false
	hasAction := false
	var topNS, topPod, op *string

	for _, it := range d.Items {
		ts, err := parseAttikeTime(it.Timestamp)
		if err != nil {
			return nil, err
		}

		mi := &model.TraceItem{
			Timestamp:    ts.Format(time.RFC3339Nano),
			ItemType:     model.TraceItemType(it.ItemType),
			SpanID:       it.SpanID,
			ParentSpanID: it.ParentSpanID,
			PodName:      it.PodName,
			Namespace:    it.Namespace,
		}

		if topNS == nil && it.Namespace != nil {
			topNS = it.Namespace
		}
		if topPod == nil && it.PodName != nil {
			topPod = it.PodName
		}

		if len(it.Payload) > 0 {
			var p esPayload
			if err := json.Unmarshal(it.Payload, &p); err == nil {
				switch it.ItemType {
				case "TRACE_START":
					var v model.TraceStartEvent
					_ = json.Unmarshal(p.TraceStart, &v)
					mi.Payload = &v
					// if v.Operation is *string:
					if op == nil && v.Operation != nil {
						op = v.Operation
					}
				case "TRACE_STOP":
					var v model.TraceStopEvent
					_ = json.Unmarshal(p.TraceStop, &v)
					mi.Payload = &v
				case "TRACE_HOP":
					var v model.TraceHopEvent
					_ = json.Unmarshal(p.TraceHop, &v)
					mi.Payload = &v
				case "GRAPHQL":
					var v model.GraphQLEvent
					_ = json.Unmarshal(p.Graphql, &v)
					mi.Payload = &v
				case "DB_SPAN":
					hasDb = true

					var esv esDbSpan
					if err := json.Unmarshal(p.DbSpan, &esv); err == nil {
						hits32, _ := parseOptInt32(esv.Hits)
						took32, _ := parseOptInt32(esv.TookMs)

						v := model.DatabaseSpanEvent{
							Action: &esv.Action,
							Query:  &esv.Query,
							Hits:   hits32,
							TookMs: took32,
						}

						mi.Payload = &v
					}

				case "ACTION":
					hasAction = true

					var esv esAction
					if err := json.Unmarshal(p.Action, &esv); err == nil {
						took32, _ := parseOptInt32(esv.TookMs)

						v := model.ActionEvent{
							Action: &esv.Action,
							Status: &esv.Status,
							TookMs: took32,
						}
						mi.Payload = &v
					}
				}
			}
		}

		items = append(items, mi)
	}

	var endStr *string
	if end != nil {
		s := end.Format(time.RFC3339Nano)
		endStr = &s
	}

	startStr := start.Format(time.RFC3339Nano)

	return &model.Trace{
		ID:           id,
		IsActive:     d.IsActive,
		TimeStarted:  startStr,
		TimeEnded:    endStr,
		TotalTimeMs:  int32(d.TotalTime),
		ResponseCode: int32(d.ResponseCode),

		Namespace: topNS,
		PodName:   topPod,
		Operation: op,
		HasDbSpan: hasDb,
		HasAction: hasAction,
		Items:     items,
	}, nil
}
