# sophokles

Daemonset collector that runs on a loop and emits host-level metrics.

## How It Works

- Scrapes node and cluster context on an interval.
- Excludes system namespaces by default.
- Publishes metrics into the `sophokles` channel on eupalinos at
  `eupalinos.attike.svc.cluster.local:50060`.

## Deployment

Base configuration and RBAC:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/base/attike/sophokles`

Development overlay (Romaioi) daemonset:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/overlays/romaioi/attike/sophokles`
