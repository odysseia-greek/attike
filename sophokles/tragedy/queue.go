package tragedy

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/service"
)

// EnqueueTask sends a task to the Eupalinos queue
func (c *Collector) enqueueTask(ctx context.Context, data string) error {
	message := &pb.Epistello{
		Data:    data,
		Channel: c.Channel,
	}

	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())
	if err != nil {
		return err
	}

	_, err = c.Eupalinos.EnqueueMessage(ctx, message)
	if err != nil {
		return err
	}
	return err
}
