package comedy

import (
	"context"
	"time"

	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

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
