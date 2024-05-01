package comedy

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
)

const (
	SPANCTX string = "SPANCTX"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func Trace(tracer pb.TraceService_ChorusClient) Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(config.HeaderKey)
			w.Header().Set(config.HeaderKey, requestId)
			trace := traceFromString(requestId)

			if !trace.Save {
				f.ServeHTTP(w, r.WithContext(r.Context()))
			}

			if trace.Save {
				go func() {
					newSpan := GenerateSpanID()
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

					trace.SpanId = newSpan

					err := tracer.Send(parabasis)
					if err != nil {
						logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
					}

					combinedId := CreateCombinedId(trace)

					ctx := context.WithValue(r.Context(), SPANCTX, combinedId)
					f.ServeHTTP(w, r.WithContext(ctx))
				}()
			}
		}
	}
}
