package comedy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DatabaseSpan(query map[string]interface{}, hits, timeTook int64, ctx context.Context, streamer arv1.TraceService_ChorusClient) {
	traceID, spanID, traceCall := ExtractRequestIds(ctx)

	if !traceCall {
		return
	}

	parsedQuery, _ := json.Marshal(query)

	dataBaseSpan := &arv1.ObserveRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       GenerateSpanID(),
		Kind: &arv1.ObserveRequest_DbSpan{DbSpan: &arv1.ObserveDbSpan{
			Action: "search",
			Query:  string(parsedQuery),
			Hits:   hits,
			TookMs: timeTook,
		}},
	}

	err := streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func CacheSpan(response string, sessionId string, ctx context.Context, streamer arv1.TraceService_ChorusClient) {
	traceID, spanID, traceCall := ExtractRequestIds(ctx)

	if !traceCall {
		return
	}

	span := &arv1.ObserveRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       GenerateSpanID(),
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: fmt.Sprintf("taken from cache with key: %s", sessionId),
			Status: response,
		}},
	}

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func ServiceToServiceSpan(span *arv1.ObserveRequest, ctx context.Context, streamer arv1.TraceService_ChorusClient) {
	traceID, spanID, traceCall := ExtractRequestIds(ctx)

	if !traceCall {
		return
	}

	span.TraceId = traceID
	span.SpanId = GenerateSpanID()
	span.ParentSpanId = spanID

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}

}

func ServiceToServiceSpanWithCtx(ctx context.Context, span *arv1.ObserveRequest, streamer arv1.TraceService_ChorusClient) (context.Context, func(err error, tookMs int64)) {
	traceID, parentSpanID, traceCall := ExtractRequestIds(ctx)
	if !traceCall {
		return ctx, func(error, int64) {}
	}

	childSpanID := GenerateSpanID()

	// emit "start" hop/action using childSpanID
	span.TraceId = traceID
	span.ParentSpanId = parentSpanID
	span.SpanId = childSpanID
	_ = streamer.Send(span)

	combined := fmt.Sprintf("%s+%s+%d", traceID, childSpanID, 1)

	// propagate to downstream AND keep parent context
	outCtx := context.WithValue(ctx, config.DefaultTracingName, combined)
	outCtx = metadata.AppendToOutgoingContext(outCtx, service.HeaderKey, combined)

	// return a finisher to emit tookMs + response_code
	finish := func(err error, tookMs int64) {
		code := int32(0)
		if err != nil {
			code = int32(status.Code(err))
		}

		done := &arv1.ObserveRequest{
			TraceId:      traceID,
			ParentSpanId: parentSpanID,
			SpanId:       childSpanID,
			Kind: &arv1.ObserveRequest_TraceHopStop{
				TraceHopStop: &arv1.ObserveTraceHopStop{
					ResponseCode: code,
					TookMs:       tookMs,
				},
			},
		}
		_ = streamer.Send(done)
	}

	return outCtx, finish
}
