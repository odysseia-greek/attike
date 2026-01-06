package comedy

import (
	"context"
	"fmt"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"

	"net/http"
	"time"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func Trace(tracer v1.TraceService_ChorusClient) Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(config.HeaderKey)
			w.Header().Set(config.HeaderKey, requestId)
			trace := traceFromString(requestId)

			if trace.Save {
				newSpan := GenerateSpanID()
				trace.SpanId = newSpan
				combinedId := CreateCombinedId(trace)
				ctx := context.WithValue(r.Context(), config.DefaultTracingName, combinedId)

				go func() {
					parabasis := &v1.ObserveRequest{
						TraceId:      trace.TraceId,
						SpanId:       newSpan,
						ParentSpanId: trace.SpanId,
						Kind: &v1.ObserveRequest_TraceHop{
							TraceHop: &v1.ObserveTraceHop{
								Method: r.Method,
								Url:    r.URL.RequestURI(),
								Host:   r.Host,
							},
						},
					}
					if err := tracer.Send(parabasis); err != nil {
						logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
					}

					logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", requestId, newSpan))
				}()

				f(w, r.WithContext(ctx))
				return
			}

			f(w, r)
		}
	}
}

func TraceWithLogAndSpan(tracer v1.TraceService_ChorusClient) Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			requestId := r.Header.Get(config.HeaderKey)
			w.Header().Set(config.HeaderKey, requestId)
			trace := traceFromString(requestId)

			if trace.Save {
				newSpan := GenerateSpanID()
				originalSpanId := trace.SpanId
				trace.SpanId = newSpan
				combinedId := CreateCombinedId(trace)
				ctx := context.WithValue(r.Context(), config.DefaultTracingName, combinedId)

				go func() {
					parabasis := &v1.ObserveRequest{
						TraceId:      trace.TraceId,
						SpanId:       newSpan,
						ParentSpanId: originalSpanId,
						Kind: &v1.ObserveRequest_TraceHop{
							TraceHop: &v1.ObserveTraceHop{
								Method: r.Method,
								Url:    r.URL.RequestURI(),
								Host:   r.Host,
							},
						},
					}
					if err := tracer.Send(parabasis); err != nil {
						logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
					}

					traceTime := time.Since(startTime)
					logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s and took: %v", requestId, newSpan, traceTime))
				}()

				f(w, r.WithContext(ctx))

				duration := time.Since(startTime)
				clientIp := r.RemoteAddr
				method := r.Method
				path := r.URL.Path
				var statusCode int
				responseWriter, ok := w.(*middleware.StatusRecorder)
				if ok {
					statusCode = responseWriter.Status
				} else {
					// if w is not our wrapped response writer, we cannot get the status
					// so, let's set the status to StatusOK for this case
					statusCode = http.StatusOK
				}
				go func() {
					parabasis := &v1.ObserveRequest{
						TraceId:      trace.TraceId,
						SpanId:       GenerateSpanID(),
						ParentSpanId: newSpan,
						Kind: &v1.ObserveRequest_Action{
							Action: &v1.ObserveAction{
								Action: "CloseSpan",
								Status: fmt.Sprintf("status code: %d", statusCode),
								TookMs: duration.Milliseconds(),
							},
						},
					}
					if err := tracer.Send(parabasis); err != nil {
						logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
					}
				}()

				logging.Api(statusCode, method, requestId, clientIp, path, duration)
				return
			}

			f(w, r)

			duration := time.Since(startTime)
			statusCode := http.StatusOK
			if responseWriter, ok := w.(*middleware.StatusRecorder); ok {
				statusCode = responseWriter.Status
			}
			logging.Api(statusCode, r.Method, requestId, r.RemoteAddr, r.URL.Path, duration)
		}
	}
}
