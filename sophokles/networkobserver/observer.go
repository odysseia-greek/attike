package networkobserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
)

type Observer struct {
	Enabled      bool
	NodeName     string
	ProviderName string
	PollInterval time.Duration
	RemotePorts  map[uint16]struct{}
	provider     provider
}

type provider interface {
	Run(context.Context, func(Event) error) error
}

func NewFromEnv() (*Observer, error) {
	// Enabled by default for the experimental test image. It remains strictly
	// best-effort and can be disabled without affecting the metrics collector.
	enabled := parseBoolEnv("ENABLE_NETWORK_OBSERVER", true)
	nodeName := getenvDefault("NODE_NAME", "local-node")
	providerName := getenvDefault("NETWORK_OBSERVER_PROVIDER", "procfs")
	pollIntervalValue := getenvDefault("NETWORK_OBSERVER_POLL_INTERVAL", "2s")

	pollInterval, err := time.ParseDuration(pollIntervalValue)
	if err != nil {
		return nil, fmt.Errorf("parse NETWORK_OBSERVER_POLL_INTERVAL: %w", err)
	}

	remotePorts, err := parseRemotePorts(os.Getenv("NETWORK_OBSERVER_REMOTE_PORTS"))
	if err != nil {
		return nil, err
	}

	p, err := newProvider(providerName, pollInterval, nodeName, remotePorts)
	if err != nil {
		return nil, err
	}

	return &Observer{
		Enabled:      enabled,
		NodeName:     nodeName,
		ProviderName: providerName,
		PollInterval: pollInterval,
		RemotePorts:  remotePorts,
		provider:     p,
	}, nil
}

func (o *Observer) Run(ctx context.Context) error {
	if !o.Enabled {
		return nil
	}

	logging.System(fmt.Sprintf(
		"starting network observer provider=%s interval=%s remote_ports=%s",
		o.ProviderName,
		o.PollInterval,
		remotePortsSummary(o.RemotePorts),
	))

	return o.provider.Run(ctx, func(event Event) error {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal network event: %w", err)
		}

		logging.System(fmt.Sprintf("network event: %s", string(payload)))
		return nil
	})
}

func parseBoolEnv(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	case "":
		return fallback
	default:
		return fallback
	}
}

func parseRemotePorts(value string) (map[uint16]struct{}, error) {
	ports := map[uint16]struct{}{}
	if value == "" {
		return ports, nil
	}

	for _, raw := range strings.Split(value, ",") {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}

		port, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("parse remote port %q: %w", part, err)
		}
		ports[uint16(port)] = struct{}{}
	}

	return ports, nil
}

func remotePortsSummary(ports map[uint16]struct{}) string {
	if len(ports) == 0 {
		return "all"
	}

	values := make([]string, 0, len(ports))
	for port := range ports {
		values = append(values, strconv.Itoa(int(port)))
	}

	return strings.Join(values, ",")
}

func getenvDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
