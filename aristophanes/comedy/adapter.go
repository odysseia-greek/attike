package comedy

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func Trace(tracer pb.TraceService_ChorusClient) Adapter {
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
					parabasis := &pb.ParabasisRequest{
						TraceId:      trace.TraceId,
						ParentSpanId: trace.SpanId,
						SpanId:       newSpan,
						RequestType: &pb.ParabasisRequest_Trace{
							Trace: &pb.TraceRequest{
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

				// Serve the request with the updated context
				f(w, r.WithContext(ctx))
				return
			}

			f(w, r)
		}
	}
}
