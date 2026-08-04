package networkobserver

import "time"

type Event struct {
	SchemaVersion int        `json:"schema_version"`
	Timestamp     time.Time  `json:"@timestamp"`
	Node          NodeRef    `json:"node"`
	Process       ProcessRef `json:"process"`
	Connection    Connection `json:"connection"`
	Source        SourceMeta `json:"source"`
}

type NodeRef struct {
	Name string `json:"name"`
}

type ProcessRef struct {
	PID  int    `json:"pid"`
	Comm string `json:"comm,omitempty"`
}

type Connection struct {
	Family      string `json:"family"`
	Protocol    string `json:"protocol"`
	State       string `json:"state"`
	LocalIP     string `json:"local_ip,omitempty"`
	LocalPort   uint16 `json:"local_port,omitempty"`
	RemoteIP    string `json:"remote_ip,omitempty"`
	RemotePort  uint16 `json:"remote_port,omitempty"`
	SocketInode string `json:"socket_inode,omitempty"`
}

type SourceMeta struct {
	Collector string `json:"collector"`
	Method    string `json:"method"`
	Provider  string `json:"provider"`
}
