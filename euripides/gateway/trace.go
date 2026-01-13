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
		logging.Error(err.Error())
		return nil, err
	}

	logging.Debug(fmt.Sprintf("trace with id: %s found %v", res.Id, res.Found))

	var esModel esTraceDoc
	source, _ := json.Marshal(res.Source)
	if err := json.Unmarshal(source, &esModel); err != nil {
		return nil, err
	}

	return toModelTrace(id, esModel)
}
