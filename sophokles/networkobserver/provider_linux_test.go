//go:build linux

package networkobserver

import "testing"

func TestParseProcfsTCPLineIPv4(t *testing.T) {
	line := "   0: 0100007F:CC8A 0102000A:23F0 01 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 20 4 30 10 -1"

	conn, ok, err := parseProcfsTCPLine(line, "ipv4", map[uint16]struct{}{9200: struct{}{}})
	if err != nil {
		t.Fatalf("parseProcfsTCPLine returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected connection to match filter")
	}
	if conn.LocalIP != "127.0.0.1" {
		t.Fatalf("unexpected local ip: %s", conn.LocalIP)
	}
	if conn.RemoteIP != "10.0.2.1" {
		t.Fatalf("unexpected remote ip: %s", conn.RemoteIP)
	}
	if conn.RemotePort != 9200 {
		t.Fatalf("unexpected remote port: %d", conn.RemotePort)
	}
}

func TestParseProcfsTCPLineSkipsNonMatchingPort(t *testing.T) {
	line := "   0: 0100007F:CC8A 0102000A:01BB 01 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 20 4 30 10 -1"

	_, ok, err := parseProcfsTCPLine(line, "ipv4", map[uint16]struct{}{9200: struct{}{}})
	if err != nil {
		t.Fatalf("parseProcfsTCPLine returned error: %v", err)
	}
	if ok {
		t.Fatal("expected connection to be filtered out")
	}
}
