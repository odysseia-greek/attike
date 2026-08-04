# sophokles

Daemonset collector that runs on a loop and emits host-level metrics.

## How It Works

- Scrapes node and cluster context on an interval.
- Excludes system namespaces by default.
- Publishes metrics into the `sophokles` channel on eupalinos at
  `eupalinos.attike.svc.cluster.local:50060`.
- Can optionally start a separate log-only network observer process path from
  `main`, independent from the metrics collector.

## Network Observer

The network observer is a standalone package and currently uses a simple
Linux `procfs` snapshot provider so you can inspect outbound TCP traffic
without touching the normal collector flow.

Environment variables:

- `ENABLE_NETWORK_OBSERVER=true` (default for the experimental image; set to
  `false` to disable)
- `NETWORK_OBSERVER_PROVIDER=procfs`
- `NETWORK_OBSERVER_POLL_INTERVAL=2s`
- `NETWORK_OBSERVER_REMOTE_PORTS=9200`

Notes:

- If `NETWORK_OBSERVER_REMOTE_PORTS` is unset, it logs all outbound TCP
  connections it sees.
- It logs the initial connection snapshot and then each newly observed TCP
  connection as JSON prefixed with `network event:`.
- The observer runs in its own goroutine. Configuration, permission, and
  runtime errors are logged but never stop or delay the metrics collector.
- `NETWORK_OBSERVER_PROVIDER=ebpf` is not implemented yet. Real eBPF will
  still need a kernel program plus a Go loader.

## Deployment

Base configuration and RBAC:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/base/attike/sophokles`

Development overlay (Romaioi) daemonset:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/overlays/romaioi/attike/sophokles`
