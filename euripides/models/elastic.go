package models

import "encoding/json"

type EsTraceDoc struct {
	IsActive     bool          `json:"isActive"`
	TimeStarted  string        `json:"timeStarted"`
	TimeEnded    *string       `json:"timeEnded,omitempty"`
	TotalTime    int64         `json:"totalTime"`
	ResponseCode int           `json:"responseCode"`
	Items        []EsTraceItem `json:"items"`
}

type EsTraceItem struct {
	Timestamp    string          `json:"timestamp"`
	ItemType     string          `json:"item_type"`
	SpanID       *string         `json:"span_id,omitempty"`
	ParentSpanID *string         `json:"parent_span_id,omitempty"`
	PodName      *string         `json:"pod_name,omitempty"`
	Namespace    *string         `json:"namespace,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"`
}

type EsPayload struct {
	Common       json.RawMessage `json:"common,omitempty"`
	TraceStart   json.RawMessage `json:"traceStart,omitempty"`
	TraceHop     json.RawMessage `json:"traceHop,omitempty"`
	Graphql      json.RawMessage `json:"graphql,omitempty"`
	Action       json.RawMessage `json:"action,omitempty"`
	DbSpan       json.RawMessage `json:"dbSpan,omitempty"`
	TraceStop    json.RawMessage `json:"traceStop,omitempty"`
	TraceHopStop json.RawMessage `json:"traceHopStop,omitempty"`
}

type EsDbSpan struct {
	Action string          `json:"action"`
	Query  string          `json:"query"`
	Hits   json.RawMessage `json:"hits"`
	TookMs json.RawMessage `json:"tookMs"`
	Target string          `json:"target,omitempty"`
	Index  string          `json:"index,omitempty"`
}

type EsAction struct {
	Action string          `json:"action"`
	Status string          `json:"status"`
	TookMs json.RawMessage `json:"tookMs"`
}

type EsHopStop struct {
	ResponseCode *int32          `json:"responseCode,omitempty"` // NOTE pointer
	TookMs       json.RawMessage `json:"tookMs"`
}
