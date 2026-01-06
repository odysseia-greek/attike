package comedy

import (
	"context"
	"fmt"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/protobuf/proto"
)

// EnqueueTask sends a task to the Eupalinos queue
func (t *TraceServiceImpl) enqueueTask(ctx context.Context, ev *v1.AttikeEvent) error {
	data, err := proto.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal trace %s trace_id=%s: %v", ev.Common.ItemType.String(), ev.Common.TraceId, err)
	}

	ctx = context.WithValue(ctx, service.HeaderKey, ev.Common.TraceId)

	message := &pb.EpistelloBytes{
		Data:    data,
		Channel: t.Channel,
	}

	queue, err := t.Eupalinos.EnqueueMessageBytes(ctx, message)
	if err != nil {
		return fmt.Errorf("Error creating queueing for trace ID %s: %s", ev.Common.TraceId, err)
	}

	logging.Debug(fmt.Sprintf("queued %s with id: %s", ev.Common.ItemType.String(), queue.Id))

	return nil
}
