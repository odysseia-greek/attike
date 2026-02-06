# Attike

Attike is an observability pipeline for traces and metrics with a GraphQL read
gateway and a dedicated frontend.

## Why Attike

The goal is to build a first-principles tracing and metrics flow. You can use
Tempo, OpenTelemetry, or Jaeger, but building a pipeline from scratch changes
how you understand event modeling, sampling, storage, and query design.

## System Overview

### Components

| Component | Runs as | Responsibility |
| --- | --- | --- |
| aristophanes | Sidecar | Receives calls from the main container and enqueues trace and metric events in eupalinos. |
| sophokles | DaemonSet | Periodic host-level collector that pushes metrics into eupalinos. |
| aiskhylos | Deployment | Consumes eupalinos streams and indexes trace and metric documents in Elastic. |
| euripides | Deployment | GraphQL gateway for querying indexed data from Elastic. |
| polykleitos | Deployment | Frontend UI for exploring traces and metrics. |

### Flow

```
[Main container in pod] -> [aristophanes sidecar] -> [eupalinos queue] <- [sophokles daemonset]
                                                           |
                                                           v
                                                     [aiskhylos] -> [Elastic]
                                                           |
                                                           v
                                                     [euripides] -> [polykleitos]
```

### Data Path Detail (Text)

- Main container emits traces to aristophanes.
- Sophokles pushes metrics to eupalinos.
- Aiskhylos reads from eupalinos and indexes into Elastic.
- Euripides queries Elastic and serves GraphQL.
- Polykleitos renders data from euripides.

## Deployment

Attike is deployed via Kustomize. Base manifests live in:
`github.com/odysseia-greek/mykenai/thrasyboulos/hydor/base/attike`

The development overlay is:
`github.com/odysseia-greek/mykenai/thrasyboulos/hydor/overlays/romaioi/attike`

Note: `romaioi` is the development environment.

### Access

- Euripides GraphQL: `https://attike.byzantium.odysseia-greek/euripides/graphql`
- Euripides REST: `https://attike.byzantium.odysseia-greek/euripides/v1`
- Polykleitos UI: `https://attike.byzantium.odysseia-greek`

## Local Development (mirrord in GoLand)

This setup is intentionally specific to the GoLand plugin workflow.

- Install and enable the mirrord plugin in GoLand.
- Open Run / Debug Configurations and add a `go build` configuration.
- Set `package` to the service you want to run locally.
- Add an env var for the target workload, for example:

```bash
TARGET=deploy/euripides/container/euripides
```

When you Run or Debug, mirrord attaches to the target workload and routes
traffic so your local process behaves as if it is inside the cluster.

## GraphQL Examples

### Metrics

```graphql
query {
  metricsSummary(
    input: {
      window: M10
    }
  ) {
    window
    start
    end
    nodes {
      total
      items {
        sampleCount
        docCount
        node
        cpu {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
          totalMax
          totalMaxHuman
        }
        mem {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
          totalMax
          totalMaxHuman
        }
        sortKey
      }
    }
    namespaces {
      total
      items {
        sampleCount
        docCount
        namespace
        cpu {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
          totalMax
          totalMaxHuman
        }
        mem {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
          totalMax
          totalMaxHuman
        }
        sortKey
      }
    }
    pods {
      total
      items {
        sampleCount
        docCount
        podName
        cpu {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
        }
        mem {
          avg
          max
          p95
          avgHuman
          maxHuman
          p95Human
        }
        sortKey
      }
    }
  }
}
```

### traceById

```graphql
query {
  trace(id: "6e571918-d5dd-4824-8a58-83772842277d") {
    id
    operation
    podName
    namespace
    timeStarted
    totalTimeMs
    hasDbSpan
    hasAction
    items {
      timestamp
      itemType
      podName
      namespace
      spanId
      parentSpanId
      payload {
        ... on DatabaseSpanEvent { action hits query tookMs }
        ... on ActionEvent { action status tookMs }
        ... on TraceHopEvent { method host url }
        ... on TraceStartEvent { operation url host remoteAddress rootQuery }
        ... on TraceStopEvent { responseBody }
        ... on GraphQLEvent { operation rootQuery }
        ... on TraceHopStopEvent { responseCode tookMs }
      }
    }
  }
}
```

### traceSearch

```graphql
query {
  traceSearch(
    input: {
      limit: 100
      window: M30
      responseCode: 200
      operation: "mediaAnswer"
      timeTookGreaterThan: 25
    }
  ) {
    total
    items {
      id
      numberOfItems
      isActive
      hasDbSpan
      responseCode
      totalTimeMs
      rootQuery
      timeStarted
      timeEnded
    }
  }
}
```

### tracePoll

```graphql
query {
  tracePoll(limit: 1) {
    updatedAt
    traces {
      id
      numberOfItems
      isActive
      hasDbSpan
      responseCode
      totalTimeMs
      rootQuery
      timeStarted
      timeEnded
    }
  }
}
```
