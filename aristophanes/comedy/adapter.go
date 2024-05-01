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
			trace := traceFromString(requestId)

			if trace.Save {
				go func(traceCopy *pb.TraceBare) {
					newSpan := GenerateSpanID()
					parabasis := &pb.ParabasisRequest{
						TraceId:      traceCopy.TraceId,
						ParentSpanId: traceCopy.SpanId,
						SpanId:       newSpan,
						RequestType: &pb.ParabasisRequest_Trace{
							Trace: &pb.TraceRequest{
								Method: r.Method,
								Url:    r.URL.RequestURI(),
								Host:   r.Host,
							},
						},
					}

					traceCopy.SpanId = newSpan

					err := tracer.Send(parabasis)
					if err != nil {
						logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
					}

					combinedId := CreateCombinedId(traceCopy)
					ctx := context.WithValue(r.Context(), config.HeaderKey, combinedId)
					f(w, r.WithContext(ctx))
				}(trace)
				return
			}

			f(w, r)
		}
	}
}
