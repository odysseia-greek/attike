package comedy

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"strings"
)

func TraceFromCtx(ctx context.Context) *pb.TraceBare {
	requestId := ctx.Value(config.HeaderKey).(string)
	return traceFromString(requestId)
}

func traceFromString(requestId string) *pb.TraceBare {
	splitID := strings.Split(requestId, "+")

	trace := &pb.TraceBare{}

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

func CreateCombinedId(trace *pb.TraceBare) string {
	saveTrace := 1
	if !trace.Save {
		saveTrace = 0
	}
	return fmt.Sprintf("%s+%s+%d", trace.TraceId, trace.SpanId, saveTrace)
}
