package models

import "time"

type TraceRootSource struct {
	TraceID        string  `json:"traceID"`
	IsActive       bool    `json:"isActive"`
	TimeStarted    string  `json:"timeStarted"`
	TimeEnded      *string `json:"timeEnded,omitempty"`
	Operation      string  `json:"operation"`
	TotalTime      int64   `json:"totalTime"`
	ResponseCode   int16   `json:"responseCode"`
	ContainsDBSpan bool    `json:"containsDBSpan"`
	NumberOfItems  int32   `json:"numberOfItems"`
}

type Trace struct {
	ID           string
	IsActive     bool
	TimeStarted  time.Time
	TimeEnded    *time.Time
	TotalTimeMs  int32
	ResponseCode int32
	Items        []TraceItem
}

type TraceItem struct {
	Timestamp    time.Time
	ItemType     string
	SpanID       *string
	ParentSpanID *string
	PodName      *string
	Namespace    *string
	Payload      any // one of the typed payload structs below (or nil)
}

type TraceHopStop struct {
	ResponseCode *int32
	TookMs       *int32
}
