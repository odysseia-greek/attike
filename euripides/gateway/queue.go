package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/euripides/models"
)

func (e *EuripidesHandler) StartTraceReportReader(ctx context.Context) {
	e.mu.Lock()
	if e.latest == nil {
		e.latest = make(map[string]models.TraceRootSource, 1024)
	}
	if e.pendingIn == nil {
		e.pendingIn = make(map[string]struct{}, 1024)
	}
	e.mu.Unlock()

	go e.traceReportLoop(ctx)
}

func (e *EuripidesHandler) traceReportLoop(ctx context.Context) {
	cfg := e.traceCfg
	if cfg.PollEvery <= 0 {
		cfg.PollEvery = 10 * time.Second
	}
	if cfg.DequeueWait <= 0 {
		cfg.DequeueWait = 300 * time.Millisecond
	}
	if cfg.MaxDrainPerPoll <= 0 {
		cfg.MaxDrainPerPoll = 5000
	}

	poll := time.NewTicker(cfg.PollEvery)
	defer poll.Stop()

	e.drainTraceReportsOnce(ctx, cfg)

	for {
		select {
		case <-ctx.Done():
			logging.Info("trace report reader stopped: " + ctx.Err().Error())
			return
		case <-poll.C:
			e.drainTraceReportsOnce(ctx, cfg)
		}
	}
}

func (e *EuripidesHandler) drainTraceReportsOnce(ctx context.Context, cfg TraceReportReaderConfig) {
	drained := 0

	for drained < cfg.MaxDrainPerPoll {
		dctx, cancel := context.WithTimeout(ctx, cfg.DequeueWait)
		msg, err := e.Eupalinos.DequeueMessageBytes(dctx, &pb.ChannelInfo{Name: e.ReportChannel})
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				break
			}
			// queue empty
			break
		}

		if msg == nil || len(msg.Data) == 0 {
			break
		}

		logging.Trace("trace report received: " + string(msg.Data))
		
		var item models.TraceRootSource
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			logging.Warn("failed to unmarshal TraceRootSource: " + err.Error())
			continue
		}
		if item.TraceID == "" {
			continue
		}

		e.upsertTraceReport(item)
		drained++
	}

	if drained > 0 {
		logging.Trace("trace report reader drained: " + strconv.Itoa(drained))
	}
}

func (e *EuripidesHandler) upsertTraceReport(item models.TraceRootSource) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.latest == nil {
		e.latest = make(map[string]models.TraceRootSource, 1024)
	}
	if e.pendingIn == nil {
		e.pendingIn = make(map[string]struct{}, 1024)
	}

	e.latest[item.TraceID] = item

	if _, ok := e.pendingIn[item.TraceID]; !ok {
		e.pendingIn[item.TraceID] = struct{}{}
		e.pendingQ = append(e.pendingQ, item.TraceID)
	}
}
