package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/euripides/graph/model"
)

func (e *EuripidesHandler) TraceById(ctx context.Context, id string) (*model.Trace, error) {
	res, err := e.Elastic.Query().GetById(ctx, e.TraceIndex, id)
	if err != nil {
		return nil, err
	}

	logging.Debug(fmt.Sprintf("trace with id: %s found %v", res.Id, res.Found))

	var trace model.Trace
	source, _ := json.Marshal(res.Source)
	if err := json.Unmarshal(source, &trace); err != nil {
		return nil, err
	}

	return &trace, nil
}
