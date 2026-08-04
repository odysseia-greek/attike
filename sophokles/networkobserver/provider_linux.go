//go:build linux

package networkobserver

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type procfsProvider struct {
	nodeName    string
	interval    time.Duration
	remotePorts map[uint16]struct{}
}

type procfsConnection struct {
	Key           string
	Family        string
	Protocol      string
	State         string
	LocalIP       string
	LocalPort     uint16
	RemoteIP      string
	RemotePort    uint16
	SocketInode   string
	PID           int
	Comm          string
	RemoteMatched bool
}

type socketOwner struct {
	PID  int
	Comm string
}

func newProvider(providerName string, interval time.Duration, nodeName string, remotePorts map[uint16]struct{}) (provider, error) {
	switch providerName {
	case "", "procfs":
		return &procfsProvider{
			nodeName:    nodeName,
			interval:    interval,
			remotePorts: remotePorts,
		}, nil
	case "ebpf":
		return nil, fmt.Errorf("NETWORK_OBSERVER_PROVIDER=ebpf is not implemented yet; real eBPF will need a kernel program and a Go loader")
	default:
		return nil, fmt.Errorf("unknown NETWORK_OBSERVER_PROVIDER %q", providerName)
	}
}

func (p *procfsProvider) Run(ctx context.Context, emit func(Event) error) error {
	current, err := p.snapshot()
	if err != nil {
		return err
	}

	// Emit the initial snapshot as well as newly observed connections. This is
	// intentionally useful for the test image: operators get immediate proof
	// that collection works without waiting for a fresh TCP connection.
	for _, conn := range current {
		if err := emit(p.toEvent(conn)); err != nil {
			return err
		}
	}

	seen := snapshotKeys(current)
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			current, err = p.snapshot()
			if err != nil {
				return err
			}

			nextSeen := make(map[string]struct{}, len(current))
			for _, conn := range current {
				nextSeen[conn.Key] = struct{}{}
				if _, ok := seen[conn.Key]; ok {
					continue
				}

				if err := emit(p.toEvent(conn)); err != nil {
					return err
				}
			}
			seen = nextSeen
		}
	}
}

func (p *procfsProvider) snapshot() ([]procfsConnection, error) {
	connections, err := readProcNetTCP("/proc/net/tcp", "ipv4", p.remotePorts)
	if err != nil {
		return nil, err
	}

	v6, err := readProcNetTCP("/proc/net/tcp6", "ipv6", p.remotePorts)
	if err != nil {
		return nil, err
	}

	connections = append(connections, v6...)
	if len(connections) == 0 {
		return connections, nil
	}

	owners, err := lookupSocketOwners(connections)
	if err != nil {
		return nil, err
	}

	for i := range connections {
		owner, ok := owners[connections[i].SocketInode]
		if !ok {
			continue
		}
		connections[i].PID = owner.PID
		connections[i].Comm = owner.Comm
	}

	return connections, nil
}

func (p *procfsProvider) toEvent(conn procfsConnection) Event {
	return Event{
		SchemaVersion: 1,
		Timestamp:     time.Now().UTC(),
		Node:          NodeRef{Name: p.nodeName},
		Process: ProcessRef{
			PID:  conn.PID,
			Comm: conn.Comm,
		},
		Connection: Connection{
			Family:      conn.Family,
			Protocol:    conn.Protocol,
			State:       conn.State,
			LocalIP:     conn.LocalIP,
			LocalPort:   conn.LocalPort,
			RemoteIP:    conn.RemoteIP,
			RemotePort:  conn.RemotePort,
			SocketInode: conn.SocketInode,
		},
		Source: SourceMeta{
			Collector: "sophokles",
			Method:    "procfs_tcp_snapshot",
			Provider:  "procfs",
		},
	}
}

func readProcNetTCP(path string, family string, remotePorts map[uint16]struct{}) ([]procfsConnection, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	connections := make([]procfsConnection, 0, len(lines)-1)
	for _, line := range lines[1:] {
		conn, ok, err := parseProcfsTCPLine(line, family, remotePorts)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		if !ok {
			continue
		}
		connections = append(connections, conn)
	}

	return connections, nil
}

func parseProcfsTCPLine(line string, family string, remotePorts map[uint16]struct{}) (procfsConnection, bool, error) {
	fields := strings.Fields(line)
	if len(fields) < 10 {
		return procfsConnection{}, false, fmt.Errorf("expected at least 10 fields, got %d", len(fields))
	}

	localIP, localPort, err := parseProcfsAddress(fields[1], family)
	if err != nil {
		return procfsConnection{}, false, err
	}

	remoteIP, remotePort, err := parseProcfsAddress(fields[2], family)
	if err != nil {
		return procfsConnection{}, false, err
	}

	state := tcpState(fields[3])
	if remoteIP == "0.0.0.0" || remoteIP == "::" || remotePort == 0 {
		return procfsConnection{}, false, nil
	}
	if state == "LISTEN" {
		return procfsConnection{}, false, nil
	}

	if len(remotePorts) > 0 {
		if _, ok := remotePorts[remotePort]; !ok {
			return procfsConnection{}, false, nil
		}
	}

	inode := fields[9]
	key := fmt.Sprintf("%s:%s:%d:%s:%d:%s:%s", family, localIP, localPort, remoteIP, remotePort, state, inode)

	return procfsConnection{
		Key:         key,
		Family:      family,
		Protocol:    "tcp",
		State:       state,
		LocalIP:     localIP,
		LocalPort:   localPort,
		RemoteIP:    remoteIP,
		RemotePort:  remotePort,
		SocketInode: inode,
	}, true, nil
}

func parseProcfsAddress(value string, family string) (string, uint16, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid procfs address %q", value)
	}

	portValue, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return "", 0, fmt.Errorf("parse port %q: %w", parts[1], err)
	}

	ip, err := decodeProcfsIP(parts[0], family)
	if err != nil {
		return "", 0, err
	}

	return ip, uint16(portValue), nil
}

func decodeProcfsIP(raw string, family string) (string, error) {
	switch family {
	case "ipv4":
		return decodeProcfsIPv4(raw)
	case "ipv6":
		return decodeProcfsIPv6(raw)
	default:
		return "", fmt.Errorf("unknown family %q", family)
	}
}

func decodeProcfsIPv4(raw string) (string, error) {
	if len(raw) != 8 {
		return "", fmt.Errorf("invalid ipv4 hex length %d", len(raw))
	}

	bytes, err := hex.DecodeString(raw)
	if err != nil {
		return "", err
	}

	return net.IP([]byte{bytes[3], bytes[2], bytes[1], bytes[0]}).String(), nil
}

func decodeProcfsIPv6(raw string) (string, error) {
	if len(raw) != 32 {
		return "", fmt.Errorf("invalid ipv6 hex length %d", len(raw))
	}

	bytes, err := hex.DecodeString(raw)
	if err != nil {
		return "", err
	}

	reordered := make([]byte, 16)
	for i := 0; i < 16; i += 4 {
		reordered[i] = bytes[i+3]
		reordered[i+1] = bytes[i+2]
		reordered[i+2] = bytes[i+1]
		reordered[i+3] = bytes[i]
	}

	return net.IP(reordered).String(), nil
}

func tcpState(raw string) string {
	switch raw {
	case "01":
		return "ESTABLISHED"
	case "02":
		return "SYN_SENT"
	case "03":
		return "SYN_RECV"
	case "04":
		return "FIN_WAIT1"
	case "05":
		return "FIN_WAIT2"
	case "06":
		return "TIME_WAIT"
	case "07":
		return "CLOSE"
	case "08":
		return "CLOSE_WAIT"
	case "09":
		return "LAST_ACK"
	case "0A":
		return "LISTEN"
	case "0B":
		return "CLOSING"
	default:
		return raw
	}
}

func lookupSocketOwners(connections []procfsConnection) (map[string]socketOwner, error) {
	inodes := make(map[string]struct{}, len(connections))
	for _, conn := range connections {
		inodes[conn.SocketInode] = struct{}{}
	}

	owners := make(map[string]socketOwner, len(inodes))
	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}

	for _, entry := range procEntries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		commBytes, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "comm"))
		if err != nil {
			continue
		}
		comm := strings.TrimSpace(string(commBytes))

		fds, err := os.ReadDir(filepath.Join("/proc", entry.Name(), "fd"))
		if err != nil {
			continue
		}

		for _, fd := range fds {
			linkTarget, err := os.Readlink(filepath.Join("/proc", entry.Name(), "fd", fd.Name()))
			if err != nil {
				continue
			}
			if !strings.HasPrefix(linkTarget, "socket:[") || !strings.HasSuffix(linkTarget, "]") {
				continue
			}

			inode := strings.TrimSuffix(strings.TrimPrefix(linkTarget, "socket:["), "]")
			if _, ok := inodes[inode]; !ok {
				continue
			}

			if _, exists := owners[inode]; !exists {
				owners[inode] = socketOwner{PID: pid, Comm: comm}
			}
		}
	}

	return owners, nil
}

func snapshotKeys(connections []procfsConnection) map[string]struct{} {
	seen := make(map[string]struct{}, len(connections))
	for _, conn := range connections {
		seen[conn.Key] = struct{}{}
	}
	return seen
}
