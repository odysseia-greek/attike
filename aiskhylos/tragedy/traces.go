package tragedy

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

type TraceConfig struct {
	Index       string
	GracePeriod time.Duration
	IdleBackoff time.Duration
	MaxOpenAge  time.Duration
}

type traceState struct {
	doc       TraceDoc
	lastSeen  time.Time
	stoppedAt *time.Time
	timer     *time.Timer
}

const attikeTSLayout = "2006-01-02T15:04:05.000"

func parseAttikeTS(s string) (time.Time, error) {
	t, err := time.Parse(attikeTSLayout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func (g *GathererImpl) StartTraces(ctx context.Context, cfg TraceConfig) {
	if cfg.Index == "" {
		cfg.Index = g.TraceIndex
	}
	if cfg.GracePeriod == 0 {
		cfg.GracePeriod = 5 * time.Second
	}
	if cfg.IdleBackoff == 0 {
		cfg.IdleBackoff = 100 * time.Millisecond
	}
	if cfg.MaxOpenAge == 0 {
		cfg.MaxOpenAge = 30 * time.Minute
	}

	go g.traceLoop(ctx, cfg)
}

func (g *GathererImpl) traceLoop(ctx context.Context, cfg TraceConfig) {
	var (
		mu     sync.Mutex
		traces = map[string]*traceState{}
	)

	finalize := func(traceID string) {
		mu.Lock()
		st := traces[traceID]
		if st == nil {
			mu.Unlock()
			return
		}
		delete(traces, traceID)
		// stop timer if any
		if st.timer != nil {
			st.timer.Stop()
		}
		doc := st.doc
		mu.Unlock()

		sortTraceItems(doc.Items)
		body, err := json.Marshal(doc)
		if err != nil {
			logging.Error("marshal trace doc failed: " + err.Error())
			return
		}

		_, err = g.Elastic.Document().CreateWithId(cfg.Index, traceID, body)
		if err != nil {
			logging.Error("index trace doc failed: " + err.Error())
			return
		}

		logging.Trace("trace doc created: " + traceID)
	}

	janitor := time.NewTicker(1 * time.Minute)
	defer janitor.Stop()

	for {
		select {
		case <-ctx.Done():
			// flush everything best-effort
			mu.Lock()
			ids := make([]string, 0, len(traces))
			for id := range traces {
				ids = append(ids, id)
			}
			mu.Unlock()

			for _, id := range ids {
				finalize(id)
			}
			return

		case <-janitor.C:
			now := time.Now().UTC()
			var evict []string

			mu.Lock()
			for id, st := range traces {
				if now.Sub(st.lastSeen) > cfg.MaxOpenAge {
					evict = append(evict, id)
				}
			}
			mu.Unlock()

			for _, id := range evict {
				finalize(id)
			}

		default:
			// Block on dequeue, but still respond to ctx cancellation:
			msg, err := g.Eupalinos.DequeueMessageBytes(ctx, &pb.ChannelInfo{Name: g.TraceChannel})
			if err != nil {
				// If ctx is cancelled, we’ll hit ctx.Done on next loop.
				time.Sleep(cfg.IdleBackoff)
				continue
			}

			var ev v1.AttikeEvent
			if err := proto.Unmarshal(msg.Data, &ev); err != nil {
				logging.Warn("failed to unmarshal AttikeEvent: " + err.Error())
				continue
			}

			common := ev.GetCommon()
			if common == nil {
				continue
			}
			traceID := common.GetTraceId()
			if traceID == "" {
				continue
			}

			// Parse common timestamp for internal ordering/timers.
			evTime, err := parseAttikeTS(common.GetTimestamp())
			if err != nil {
				evTime = time.Now().UTC()
			}

			// Store the full protobuf event JSON as payload (debug-friendly).
			evJSON, _ := protojson.MarshalOptions{EmitUnpopulated: false}.Marshal(&ev)

			item := TraceItem{
				Timestamp:    common.GetTimestamp(), // string in your schema
				ItemType:     common.GetItemType().String(),
				SpanID:       common.GetSpanId(),
				ParentSpanID: common.GetParentSpanId(),
				PodName:      common.GetPodName(),
				Namespace:    common.GetNamespace(),
				Payload:      json.RawMessage(evJSON),
			}

			mu.Lock()
			st := traces[traceID]
			if st == nil {
				st = &traceState{
					doc: TraceDoc{
						IsActive: true,
						Items:    make([]TraceItem, 0, 64),
						// TimeStarted is filled when we see TRACE_START or TRACE_STOP with time_started
					},
					lastSeen: evTime,
				}
				traces[traceID] = st
			}

			// Update state + append item
			st.lastSeen = evTime
			st.doc.Items = append(st.doc.Items, item)

			// Apply semantics based on payload kind
			switch ev.GetPayload().(type) {
			case *v1.AttikeEvent_TraceStart:
				st.doc.IsActive = true

				// Set start time if not already set (or keep earliest if you prefer)
				if st.doc.TimeStarted == "" {
					st.doc.TimeStarted = common.GetTimestamp()
				}

			case *v1.AttikeEvent_TraceStop:
				stop := ev.GetTraceStop()
				st.doc.IsActive = false

				// Prefer stop-provided time_started/time_ended if present.
				if ts := stop.GetTimeStarted(); ts != "" {
					st.doc.TimeStarted = ts
				} else if st.doc.TimeStarted == "" {
					// fallback: if never got start, use common timestamp
					st.doc.TimeStarted = common.GetTimestamp()
				}

				endedStr := stop.GetTimeEnded()
				if endedStr == "" {
					endedStr = common.GetTimestamp()
				}
				st.doc.TimeEnded = &endedStr

				// Response code
				if stop.GetResponseCode() != 0 {
					st.doc.ResponseCode = int16(stop.GetResponseCode())
				}

				// Total time
				if stop.GetTotalTimeMs() != 0 {
					st.doc.TotalTime = stop.GetTotalTimeMs()
				} else {
					// compute if we can parse both timestamps
					var (
						t0, t1 time.Time
						ok0    bool
						ok1    bool
					)

					if st.doc.TimeStarted != "" {
						if tt, err := parseAttikeTS(st.doc.TimeStarted); err == nil {
							t0, ok0 = tt, true
						}
					}
					if endedStr != "" {
						if tt, err := parseAttikeTS(endedStr); err == nil {
							t1, ok1 = tt, true
						}
					}
					if ok0 && ok1 {
						st.doc.TotalTime = t1.Sub(t0).Milliseconds()
					}
				}

				// internal stoppedAt used for grace logic
				endedTime := evTime
				if tt, err := parseAttikeTS(endedStr); err == nil {
					endedTime = tt
				}
				stoppedAt := endedTime
				st.stoppedAt = &stoppedAt

				// Reset grace timer: after GracePeriod with no later events, finalize.
				if st.timer != nil {
					st.timer.Stop()
				}

				traceIDCopy := traceID
				st.timer = time.AfterFunc(cfg.GracePeriod, func() {
					mu.Lock()
					cur := traces[traceIDCopy]
					if cur == nil {
						mu.Unlock()
						return
					}
					lastSeen := cur.lastSeen
					stopSeen := cur.stoppedAt
					mu.Unlock()

					// If we saw events after stop time, don't finalize yet
					if stopSeen != nil && lastSeen.After(*stopSeen) {
						return
					}
					finalize(traceIDCopy)
				})
			}

			mu.Unlock()
		}
	}
}

func sortTraceItems(items []TraceItem) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.Timestamp != b.Timestamp {
			return a.Timestamp < b.Timestamp
		}
		// tie-breakers (pick whatever you like)
		if a.SpanID != b.SpanID {
			return a.SpanID < b.SpanID
		}
		return a.ItemType < b.ItemType
	})
}
