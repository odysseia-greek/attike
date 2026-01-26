package gateway

import (
	"encoding/json"
	"time"

	"github.com/odysseia-greek/attike/euripides/graph/model"
	"github.com/odysseia-greek/attike/euripides/models"
)

func toModelTrace(id string, d models.EsTraceDoc) (*model.Trace, error) {
	start, err := models.ParseAttikeTime(d.TimeStarted)
	if err != nil {
		return nil, err
	}

	var end *time.Time
	if d.TimeEnded != nil && *d.TimeEnded != "" {
		te, err := models.ParseAttikeTime(*d.TimeEnded)
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
		ts, err := models.ParseAttikeTime(it.Timestamp)
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

		if !models.IsEmptyJSON(it.Payload) {
			var p models.EsPayload
			if err := json.Unmarshal(it.Payload, &p); err == nil {
				switch it.ItemType {
				case "TRACE_START":
					if !models.IsEmptyJSON(p.TraceStart) {
						var v model.TraceStartEvent
						_ = json.Unmarshal(p.TraceStart, &v)
						mi.Payload = &v
						if op == nil && v.Operation != nil {
							op = v.Operation
						}
					}

				case "TRACE_STOP":
					if !models.IsEmptyJSON(p.TraceStop) {
						var v model.TraceStopEvent
						_ = json.Unmarshal(p.TraceStop, &v)
						mi.Payload = &v
					}

				case "TRACE_HOP":
					if !models.IsEmptyJSON(p.TraceHop) {
						var v model.TraceHopEvent
						_ = json.Unmarshal(p.TraceHop, &v)
						mi.Payload = &v
					}

				case "GRAPHQL":
					if !models.IsEmptyJSON(p.Graphql) {
						var v model.GraphQLEvent
						_ = json.Unmarshal(p.Graphql, &v)
						mi.Payload = &v
					}

				case "DB_SPAN":
					hasDb = true
					if !models.IsEmptyJSON(p.DbSpan) {
						var esv models.EsDbSpan
						if err := json.Unmarshal(p.DbSpan, &esv); err == nil {
							hits32, _ := models.ParseOptInt32(esv.Hits)
							took32, _ := models.ParseOptInt32(esv.TookMs)

							v := model.DatabaseSpanEvent{
								Action: &esv.Action,
								Query:  &esv.Query,
								Hits:   hits32,
								TookMs: took32,
							}
							mi.Payload = &v
						}
					}

				case "ACTION":
					hasAction = true
					if !models.IsEmptyJSON(p.Action) {
						var esv models.EsAction
						if err := json.Unmarshal(p.Action, &esv); err == nil {
							took32, _ := models.ParseOptInt32(esv.TookMs)
							v := model.ActionEvent{
								Action: &esv.Action,
								Status: &esv.Status,
								TookMs: took32,
							}
							mi.Payload = &v
						}
					}

				case "TRACE_HOP_STOP":
					if !models.IsEmptyJSON(p.TraceHopStop) {
						var esv models.EsHopStop
						if err := json.Unmarshal(p.TraceHopStop, &esv); err == nil {
							took32, _ := models.ParseOptInt32(esv.TookMs)
							v := model.TraceHopStopEvent{
								ResponseCode: esv.ResponseCode, // already *int32
								TookMs:       took32,
							}
							mi.Payload = &v
						}
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
