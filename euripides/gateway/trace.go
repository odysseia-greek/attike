package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/euripides/elastic"
	"github.com/odysseia-greek/attike/euripides/graph/model"
	"github.com/odysseia-greek/attike/euripides/models"
)

func (e *EuripidesHandler) TraceById(ctx context.Context, id string) (*model.Trace, error) {
	res, err := e.Elastic.Query().GetById(ctx, e.TraceIndex, id)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	logging.Debug(fmt.Sprintf("trace with id: %s found %v", res.Id, res.Found))

	var esModel models.EsTraceDoc
	source, _ := json.Marshal(res.Source)
	if err := json.Unmarshal(source, &esModel); err != nil {
		return nil, err
	}

	traceModel, err := toModelTrace(id, esModel)
	jsonModel, _ := json.Marshal(traceModel)
	logging.Debug(fmt.Sprintf("trace model: %s", jsonModel))

	if err != nil {
		return nil, err
	}

	return traceModel, nil
}

func (e *EuripidesHandler) TraceSearch(ctx context.Context, input *model.TraceSearchInput) (*model.TracePage, error) {
	page := &model.TracePage{
		Items: []*model.TraceSummary{},
		Total: 0,
	}

	limit := int32(20)
	if input != nil && input.Limit != nil && *input.Limit > 0 {
		limit = *input.Limit
	}
	if limit > 500 {
		limit = 500
	}

	query := elastic.BuildTraceSearchQuery(input, limit)

	// Prefer a raw-search method if you have it.
	res, err := e.Elastic.Query().Match(e.TraceIndex, query)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	if len(res.Hits.Hits) == 0 {
		return page, nil
	}

	for _, hit := range res.Hits.Hits {
		var esSum models.TraceRootSource

		source, _ := json.Marshal(hit.Source)
		if err := json.Unmarshal(source, &esSum); err != nil {
			return nil, err
		}

		page.Items = append(page.Items, toModelTraceSummary(hit.ID, esSum))
	}

	page.Total = int32(len(page.Items))

	return page, nil
}

func (e *EuripidesHandler) TracePoll(ctx context.Context, limit *int) (*model.EnqueueItems, error) {
	n := 20
	if limit != nil && *limit > 0 {
		n = *limit
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.pendingQ) == 0 {
		return &model.EnqueueItems{
			Traces:    []*model.TraceSummary{},
			UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano), // optional; can be omitted from schema too
		}, nil
	}

	if n > len(e.pendingQ) {
		n = len(e.pendingQ)
	}

	ids := e.pendingQ[:n]
	e.pendingQ = e.pendingQ[n:]

	out := make([]*model.TraceSummary, 0, len(ids))
	for _, id := range ids {
		delete(e.pendingIn, id)

		src, ok := e.latest[id]
		if !ok {
			continue
		}

		// Option 1: also remove from latest so "live feed" doesn't retain memory forever:
		delete(e.latest, id)

		// map -> gqlgen model
		dst := &model.TraceSummary{
			ID:            src.TraceID,
			IsActive:      src.IsActive,
			HasDbSpan:     src.ContainsDBSpan,
			TotalTimeMs:   safeInt32FromInt64(src.TotalTime),
			ResponseCode:  int32(src.ResponseCode),
			NumberOfItems: src.NumberOfItems,
			RootQuery:     src.Operation,
		}
		if src.TimeStarted != "" {
			ts := src.TimeStarted
			dst.TimeStarted = &ts
		}
		dst.TimeEnded = src.TimeEnded

		out = append(out, dst)
	}

	return &model.EnqueueItems{
		Traces:    out,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}, nil
}

func toModelTraceSummary(id string, s models.TraceRootSource) *model.TraceSummary {
	var ts *string
	if s.TimeStarted != "" {
		ts = &s.TimeStarted
	}

	return &model.TraceSummary{
		ID:            id,
		IsActive:      s.IsActive,
		HasDbSpan:     s.ContainsDBSpan,
		TimeStarted:   ts,
		TimeEnded:     s.TimeEnded,
		RootQuery:     firstNonEmpty(s.Operation, "unknown"),
		TotalTimeMs:   int32(s.TotalTime),
		ResponseCode:  int32(s.ResponseCode),
		NumberOfItems: s.NumberOfItems,
	}
}

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func safeInt32FromInt64(v int64) int32 {
	if v > int64(^uint32(0)>>1) { // MaxInt32
		return int32(^uint32(0) >> 1)
	}
	if v < -int64(^uint32(0)>>1)-1 { // MinInt32
		return -int32(^uint32(0)>>1) - 1
	}
	return int32(v)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
