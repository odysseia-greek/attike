package comedy

import (
	"context"
	"fmt"
	"strings"

	"github.com/odysseia-greek/agora/plato/config"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

func TraceFromCtx(ctx context.Context) *v1.TraceBare {
	requestId := ctx.Value(config.HeaderKey).(string)
	return traceFromString(requestId)
}

func traceFromString(requestId string) *v1.TraceBare {
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
