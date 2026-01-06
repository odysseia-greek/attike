package tragedy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
)

type Collector struct {
	NodeName          string
	ScrapeInterval    time.Duration
	ExcludeNamespaces map[string]struct{}

	// Kubelet stats via apiserver proxy
	HTTPClient   *http.Client
	BearerToken  string
	APIServerURL string

	// Queue
	Eupalinos *stomion.QueueClient
	Channel   string
}

func NewCollector() (*Collector, error) {
	nodeName := config.StringFromEnv("NODE_NAME", "local-node")
	intervalStr := config.StringFromEnv("SCRAPE_INTERVAL", "10s")

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		interval = 10 * time.Second
	}

	exclude := map[string]struct{}{
		"kube-system":     {},
		"kube-public":     {},
		"kube-node-lease": {},
	}

	// Allow override: EXCLUDE_NAMESPACES="kube-system,monitoring"
	if v := os.Getenv("EXCLUDE_NAMESPACES"); v != "" {
		exclude = map[string]struct{}{}
		for _, ns := range splitCSV(v) {
			if ns != "" {
				exclude[ns] = struct{}{}
			}
		}
	}

	apiServer := config.StringFromEnv(config.KubernetesApiEndPointEnv, config.KubernetesApiEndpoint)

	tokenBytes, err := os.ReadFile(saTokenPath)
	if err != nil {
		return nil, fmt.Errorf("read serviceaccount token: %w", err)
	}

	caBytes, err := os.ReadFile(saCAPath)
	if err != nil {
		return nil, fmt.Errorf("read serviceaccount ca: %w", err)
	}
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caBytes); !ok {
		return nil, fmt.Errorf("failed to parse serviceaccount ca cert")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caPool},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   25 * time.Second,
	}

	eupalinosAddress := config.StringFromEnv("EUPALINOS_TRACING_SERVICE", config.DefaultEupalinosService)
	logging.Debug(fmt.Sprintf("creating new eupalinos client: %s", eupalinosAddress))
	queue, err := stomion.NewEupalinosClient(eupalinosAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	logging.Debug("waiting for queue to be ready")
	queueHealthy := queue.WaitForHealthyState()
	if !queueHealthy {
		logging.Debug("no queue that is healthy")
	}

	logging.Debug("queue healthy starting up")

	channel := config.StringFromEnv(config.EnvChannel, "sophokles")

	return &Collector{
		NodeName:          nodeName,
		ScrapeInterval:    interval,
		ExcludeNamespaces: exclude,

		HTTPClient:   client,
		BearerToken:  strings.TrimSpace(string(tokenBytes)),
		APIServerURL: apiServer,

		Eupalinos: queue,
		Channel:   channel,
	}, nil
}

// tiny helper to avoid pulling in extra deps
func splitCSV(s string) []string {
	out := make([]string, 0, 8)
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, trimSpaces(cur))
			cur = ""
			continue
		}
		cur += string(s[i])
	}
	out = append(out, trimSpaces(cur))
	return out
}

func trimSpaces(s string) string {
	// cheap trim for this use case
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\n' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\n' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
