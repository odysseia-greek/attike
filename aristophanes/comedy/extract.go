package comedy

import (
	"context"
	"fmt"
	"strings"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc/metadata"
)

func ExtractRequestIds(ctx context.Context) (traceID string, spanID string, traceCall bool) {
	// 1) Prefer the current tracing value in ctx (set by interceptor)
	if v := ctx.Value(config.DefaultTracingName); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return TraceFromString(s)
		}
	}

	// 2) Fallback: incoming metadata header
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", false
	}
	vals := md.Get(service.HeaderKey)
	if len(vals) == 0 {
		return "", "", false
	}
	return TraceFromString(vals[0])
}

func TraceFromCtx(ctx context.Context) *v1.TraceBare {
	v, _ := ctx.Value(config.DefaultTracingName).(string)
	if v == "" {
		return &v1.TraceBare{}
	}
	return TraceBareFromString(v)
}

func TraceBareFromString(requestId string) *v1.TraceBare {
	splitID := strings.Split(requestId, "+")

	trace := &v1.TraceBare{}

	if len(splitID) >= 3 {
		trace.Save = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		trace.TraceId = splitID[0]
	}
	if len(splitID) >= 2 {
		trace.SpanId = splitID[1]
	}

	return trace
}

func TraceFromString(requestId string) (traceID string, spanID string, traceCall bool) {
	splitID := strings.Split(requestId, "+")
	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		spanID = splitID[1]
	}
	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}
	return
}

func CreateRequestId(traceID, spanID string, traceCall bool) string {
	flag := "0"
	if traceCall {
		flag = "1"
	}
	return traceID + "+" + spanID + "+" + flag
}

func CreateCombinedId(trace *v1.TraceBare) string {
	saveTrace := 1
	if !trace.Save {
		saveTrace = 0
	}
	return fmt.Sprintf("%s+%s+%d", trace.TraceId, trace.SpanId, saveTrace)
}
