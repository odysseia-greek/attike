package comedy

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

// ManageStartTimeMap runs a single-goroutine registry that tracks trace start times.
// It serializes all access to an internal map via commands, allowing concurrent
// goroutines to safely set, retrieve, and delete start times without locks.
// This is used by the tracing pipeline to correlate trace lifecycle events
// (e.g. computing durations when spans or traces complete).
func (t *TraceServiceImpl) ManageStartTimeMap() {
	startTimeMap := make(map[string]time.Time)
	for cmd := range t.commands {
		switch cmd.Action {
		case "set":
			startTimeMap[cmd.TraceID] = cmd.Time
		case "get":
			startTime, found := startTimeMap[cmd.TraceID]
			cmd.Response <- MapResponse{Time: startTime, Found: found}
		case "delete":
			delete(startTimeMap, cmd.TraceID)
		}
	}
}

func (t *TraceServiceImpl) HealthCheck(ctx context.Context, start *v1.Empty) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: true}, nil
}

func (t *TraceServiceImpl) Chorus(stream v1.TraceService_ChorusServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.ObserveResponse{Ack: "Received"})
		}
		if err != nil {
			return err
		}

		switch req := in.Kind.(type) {

		case *v1.ObserveRequest_TraceStart:
			go t.safeExecute(func() {
				t.HandleTraceStart(req.TraceStart, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_TraceHop:
			go t.safeExecute(func() {
				t.HandleTraceHop(req.TraceHop, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_Graphql:
			go t.safeExecute(func() {
				t.HandleGraphQL(req.Graphql, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_Action:
			go t.safeExecute(func() {
				t.HandleAction(req.Action, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_DbSpan:
			go t.safeExecute(func() {
				t.HandleDbSpan(req.DbSpan, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_TraceStop:
			go t.safeExecute(func() {
				t.HandleTraceStop(req.TraceStop, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		case *v1.ObserveRequest_TraceHopStop:
			go t.safeExecute(func() {
				t.HandleTraceStopHop(req.TraceHopStop, in.TraceId, in.SpanId, in.ParentSpanId)
			})

		default:
			logging.Debug("Unhandled observe request kind")
		}
	}
}

func (t *TraceServiceImpl) HandleTraceStart(start *v1.ObserveTraceStart, traceID, spanID, parentSpanID string) {
	traceTime := time.Now().UTC()

	t.commands <- MapCommand{
		Action:  "set",
		TraceID: traceID,
		Time:    traceTime,
	}

	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: parentSpanID, // or spanID if you want root parent == root
			Timestamp:    traceTime.Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_TRACE_START,
		},
		Payload: &v1.AttikeEvent_TraceStart{
			TraceStart: &v1.TraceStartEvent{
				Method:        start.Method,
				Url:           start.Url,
				Host:          start.Host,
				RemoteAddress: start.RemoteAddress,
				RootQuery:     start.RootQuery,
				Operation:     start.Operation,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
	}
}

func (t *TraceServiceImpl) HandleTraceStop(close *v1.ObserveTraceStop, traceID, spanID, parentSpanID string) {
	traceTime := time.Now().UTC()

	// get start time from your registry
	responseChan := make(chan MapResponse, 1)
	t.commands <- MapCommand{Action: "get", TraceID: traceID, Response: responseChan}
	resp := <-responseChan

	timeStarted := "1970-01-01T00:00:00.000"
	totalTime := int64(0)
	if resp.Found {
		timeStarted = resp.Time.UTC().Format("2006-01-02T15:04:05.000")
		totalTime = traceTime.Sub(resp.Time).Milliseconds()
	}

	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: parentSpanID, // or spanID if you want root parent == root
			Timestamp:    traceTime.Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_TRACE_STOP,
		},
		Payload: &v1.AttikeEvent_TraceStop{
			TraceStop: &v1.TraceStopEvent{
				ResponseBody: close.ResponseBody,
				ResponseCode: close.ResponseCode,
				TimeStarted:  timeStarted,
				TimeEnded:    traceTime.Format("2006-01-02T15:04:05.000"),
				TotalTimeMs:  totalTime,
				IsClosed:     true,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
	}

	t.commands <- MapCommand{Action: "delete", TraceID: traceID}
}

func (t *TraceServiceImpl) HandleTraceHop(in *v1.ObserveTraceHop, traceID, spanID, parentSpanID string) {
	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: parentSpanID, // or spanID if you want root parent == root
			Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_TRACE_HOP,
		},
		Payload: &v1.AttikeEvent_TraceHop{
			TraceHop: &v1.TraceHopEvent{
				Method: in.Method,
				Url:    in.Url,
				Host:   in.Host,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

func (t *TraceServiceImpl) HandleTraceStopHop(close *v1.ObserveTraceHopStop, traceID, spanID, parentSpanID string) {
	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: parentSpanID, // or spanID if you want root parent == root
			Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_TRACE_HOP_STOP,
		},
		Payload: &v1.AttikeEvent_TraceHopStop{
			TraceHopStop: &v1.TraceHopStopEvent{
				ResponseCode: close.ResponseCode,
				TookMs:       close.TookMs,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

func (t *TraceServiceImpl) HandleAction(in *v1.ObserveAction, traceID, spanID, ParentSpanID string) {
	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: ParentSpanID, // or spanID if you want root parent == root
			Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_ACTION,
		},
		Payload: &v1.AttikeEvent_Action{
			Action: &v1.ActionEvent{
				Action: in.Action,
				Status: in.Status,
				TookMs: in.TookMs,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

func (t *TraceServiceImpl) HandleDbSpan(in *v1.ObserveDbSpan, traceID, spanID, ParentSpanID string) {
	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: ParentSpanID, // or spanID if you want root parent == root
			Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_DB_SPAN,
		},
		Payload: &v1.AttikeEvent_DbSpan{
			DbSpan: &v1.DatabaseSpanEvent{
				Action: in.Action,
				Query:  in.Query,
				Hits:   in.Hits,
				TookMs: in.TookMs,
				Target: in.Target,
				Index:  in.Index,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

func (t *TraceServiceImpl) HandleGraphQL(in *v1.ObserveGraphQL, traceID, spanID, ParentSpanID string) {
	ev := &v1.AttikeEvent{
		Common: &v1.TraceCommon{
			TraceId:      traceID,
			SpanId:       spanID,
			ParentSpanId: ParentSpanID, // or spanID if you want root parent == root
			Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05.000"),
			PodName:      t.PodName,
			Namespace:    t.Namespace,
			ItemType:     v1.ItemType_GRAPHQL,
		},
		Payload: &v1.AttikeEvent_Graphql{
			Graphql: &v1.GraphQLEvent{
				Operation: in.Operation,
				RootQuery: in.RootQuery,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.baseCtx, 3*time.Second)
	defer cancel()

	err := t.enqueueTask(ctx, ev)
	if err != nil {
		logging.Error(err.Error())
		return
	}
}

func (t *TraceServiceImpl) safeExecute(action func()) {
	defer func() {
		if r := recover(); r != nil {
			logging.Warn(fmt.Sprintf("Recovered from panic: %v\n%s", r, debug.Stack()))
		}
	}()
	action()
}
