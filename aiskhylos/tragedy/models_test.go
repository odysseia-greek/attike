package tragedy

import (
	"testing"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
)

func TestParseMetricSampleFromEpistello(t *testing.T) {
	msg := &pb.Epistello{
		Id:      "74be63dd-258d-4dd2-bc44-a5039728fcaf",
		Data:    `{"schema_version":1,"@timestamp":"2026-01-08T15:56:01.529497354Z","node":{"name":"lima-k0s-byzantium"},"node_totals":{"cpu_mcores":526,"mem_bytes":8522702848},"workload":{"namespace":"delphi","pod_name":"solon-7f6d88949f-2djqc","pod_uid":"b858346d-08b3-425c-9450-64b3a6605b5e"},"pod_totals":{"cpu_mcores":0,"mem_bytes":21962752},"containers":[{"name":"solon","totals":{"cpu_mcores":0,"mem_bytes":16531456}},{"name":"aristophanes","totals":{"cpu_mcores":0,"mem_bytes":5431296}}],"source":{"collector":"sophokles","method":"kubelet_stats_summary_via_apiserver_proxy","scrape_duration_ms":56}}`,
		Channel: "sophokles",
	}

	doc, err := ParseMetricSample(msg)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if doc.SchemaVersion != 1 {
		t.Fatalf("expected schema_version=1, got %d", doc.SchemaVersion)
	}
	if doc.Node.Name != "lima-k0s-byzantium" {
		t.Fatalf("expected node.name lima-k0s-byzantium, got %q", doc.Node.Name)
	}
	if doc.Workload.Namespace != "delphi" {
		t.Fatalf("expected workload.namespace delphi, got %q", doc.Workload.Namespace)
	}
	if len(doc.Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(doc.Containers))
	}

	wantTS, err := time.Parse(time.RFC3339Nano, "2026-01-08T15:56:01.529497354Z")
	if err != nil {
		t.Fatalf("failed to parse want timestamp: %v", err)
	}
	if !doc.Timestamp.Equal(wantTS) {
		t.Fatalf("timestamp mismatch: got %s want %s",
			doc.Timestamp.Format(time.RFC3339Nano),
			wantTS.Format(time.RFC3339Nano),
		)
	}

	if doc.Source.ScrapeDurationMS != 56 {
		t.Fatalf("expected scrape_duration_ms=56, got %d", doc.Source.ScrapeDurationMS)
	}
}

func TestParseTraceSample(t *testing.T) {

}
