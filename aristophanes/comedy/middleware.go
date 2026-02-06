package comedy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func UnaryServerInterceptor(streamer Streamer, opts ...Option) grpc.UnaryServerInterceptor {
	c := &cfg{
		HeaderKey:      "x-request-id", // override with WithHeaderKey(config.HeaderKey)
		ContextKeyName: "tracing",      // override with WithContextKeyName(config.DefaultTracingName)
		EmitCloseHop:   true,
	}
	for _, o := range opts {
		o(c)
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if streamer == nil {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var requestID string
		if v := md.Get(c.HeaderKey); len(v) > 0 {
			requestID = v[0]
		}
		if requestID == "" {
			return handler(ctx, req)
		}

		meta := TraceBareFromString(requestID)
		if !meta.Save || meta.TraceId == "" || meta.SpanId == "" {
			return handler(ctx, req)
		}

		// who called us
		host := ""
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			host = p.Addr.String()
		}

		// create a new span for THIS hop
		hopSpanID := GenerateSpanID()

		// set new request id in ctx: traceID + newSpanID + 1
		nextReqID := CreateRequestId(meta.TraceId, hopSpanID, true)
		ctx = context.WithValue(ctx, c.ContextKeyName, nextReqID)

		// emit hop start
		_ = streamer.Send(&arv1.ObserveRequest{
			TraceId:      meta.TraceId,
			ParentSpanId: meta.SpanId, // the caller span becomes our parent
			SpanId:       hopSpanID,   // our hop span
			Kind: &arv1.ObserveRequest_TraceHop{
				TraceHop: &arv1.ObserveTraceHop{
					Method: info.FullMethod,
					Url:    info.FullMethod,
					Host:   host,
				},
			},
		})

		// if you need to send a header back, do it here (optional)
		// NOTE: this should probably be the *nextReqID* not just traceID
		// but keep your current behavior if other code expects traceID only.
		_ = grpc.SendHeader(ctx, metadata.New(map[string]string{
			c.HeaderKey: nextReqID,
		}))

		start := time.Now()
		resp, err := handler(ctx, req)
		dur := time.Since(start)

		if c.EmitCloseHop {
			st := status.Convert(err)
			code := int32(st.Code())

			_ = streamer.Send(&arv1.ObserveRequest{
				TraceId:      meta.TraceId,
				ParentSpanId: meta.SpanId,
				SpanId:       hopSpanID,
				Kind: &arv1.ObserveRequest_TraceHopStop{
					TraceHopStop: &arv1.ObserveTraceHopStop{
						ResponseCode: code,
						TookMs:       dur.Milliseconds(),
					},
				},
			})
		}

		return resp, err
	}
}

func TraceWithHopStop(tracer arv1.TraceService_ChorusClient) Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			requestId := r.Header.Get(config.HeaderKey)
			w.Header().Set(config.HeaderKey, requestId) // keep existing behavior
			trace := TraceBareFromString(requestId)

			// Fall back to normal behavior when not tracing
			if tracer == nil || !trace.Save || trace.TraceId == "" || trace.SpanId == "" {
				f(w, r)
				// your api logging stays as-is
				duration := time.Since(startTime)
				statusCode := http.StatusOK
				if responseWriter, ok := w.(*middleware.StatusRecorder); ok && responseWriter.Status != 0 {
					statusCode = responseWriter.Status
				}
				logging.Api(statusCode, r.Method, requestId, r.RemoteAddr, r.URL.Path, duration)
				return
			}

			// IMPORTANT: preserve incoming span as parent BEFORE generating a new hop span.
			parentSpanId := trace.SpanId

			// New span for *this* hop (HTTP handler)
			hopSpanId := GenerateSpanID()

			// Update trace for propagation downstream from this service
			trace.SpanId = hopSpanId
			combinedId := CreateCombinedId(trace)

			// Put the new request id into context (this is what downstream uses)
			ctx := context.WithValue(r.Context(), config.DefaultTracingName, combinedId)
			// if other code reads from config.HeaderKey in ctx, keep this too:
			ctx = context.WithValue(ctx, config.HeaderKey, combinedId)

			// Optionally return the updated id to the caller (debugging / chaining)
			// If you *don’t* want to expose it, keep the old requestId header instead.
			w.Header().Set(config.HeaderKey, combinedId)

			// Emit hop start (sync is usually better for ordering)
			if err := tracer.Send(&arv1.ObserveRequest{
				TraceId:      trace.TraceId,
				ParentSpanId: parentSpanId, // caller span
				SpanId:       hopSpanId,    // this hop span
				Kind: &arv1.ObserveRequest_TraceHop{
					TraceHop: &arv1.ObserveTraceHop{
						Method: r.Method,
						Url:    r.URL.RequestURI(),
						Host:   r.Host,
					},
				},
			}); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace hop: %v", err))
			}

			// ---- CALL THE HANDLER (preserved) ----
			f(w, r.WithContext(ctx))
			// -------------------------------------

			// Now we can assert everything after it returns
			duration := time.Since(startTime)

			statusCode := http.StatusOK
			if responseWriter, ok := w.(*middleware.StatusRecorder); ok && responseWriter.Status != 0 {
				statusCode = responseWriter.Status
			}

			// Emit hop stop (this replaces CloseSpan Action)
			if err := tracer.Send(&arv1.ObserveRequest{
				TraceId:      trace.TraceId,
				ParentSpanId: parentSpanId, // stop is child of the hop span
				SpanId:       hopSpanId,    // or reuse hopSpanId if you prefer "close same span"
				Kind: &arv1.ObserveRequest_TraceHopStop{
					TraceHopStop: &arv1.ObserveTraceHopStop{
						ResponseCode: int32(statusCode),
						TookMs:       duration.Milliseconds(),
					},
				},
			}); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace hop stop: %v", err))
			}

			logging.Api(statusCode, r.Method, requestId, r.RemoteAddr, r.URL.Path, duration)
		}
	}
}
