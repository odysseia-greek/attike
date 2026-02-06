# euripides

GraphQL gateway for querying trace and metric data stored in Elastic.

## How It Works

- Resolves GraphQL queries into Elastic searches and aggregations.
- Exposes a REST interface for supporting endpoints.
- Can publish report data to the `euripides` channel on eupalinos at
  `eupalinos.attike.svc.cluster.local:50060`.

## Deployment

Base configuration:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/base/attike/euripides`

Development overlay (Romaioi) deployment and ingress:
- `github.com/odysseia-greek/mykenai/thrasyboulos/hydor/overlays/romaioi/attike/euripides`

## Access

- GraphQL: `https://attike.byzantium.odysseia-greek/euripides/graphql`
- REST: `https://attike.byzantium.odysseia-greek/euripides/v1`
