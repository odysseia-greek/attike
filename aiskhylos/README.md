# aiskhylos

Ingest worker that gathers events from eupalinos and indexes them in Elastic.

## How It Works

- Consumes trace events from the `aristophanes` channel.
- Consumes metric events from the `sophokles` channel.
- Writes normalized documents into Elastic for fast query access.
- Connects to eupalinos at `eupalinos.attike.svc.cluster.local:50060`.

## Deployment

Base configuration:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/base/attike/aiskhylos`

Development overlay (Romaioi) deployment and jobs:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/overlays/romaioi/attike/aiskhylos`
