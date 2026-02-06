# aristophanes

Sidecar that receives trace and metric events from the main container in a pod
and forwards them to the eupalinos queue.

## How It Works

- Accepts events from the co-located application container.
- Pushes trace events into the `aristophanes` channel on eupalinos.
- Keeps the app path lightweight by offloading queue transport.

## Deployment

Aristophanes runs as a sidecar inside application pods. Attike's Kustomize
manifests do not deploy it directly; it is added to workloads that need tracing.
It publishes to the eupalinos service in the `attike` namespace.
