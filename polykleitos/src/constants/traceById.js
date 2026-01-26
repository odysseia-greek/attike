import gql from "graphql-tag";

export const TraceByIdQuery = gql`
  query traceById($id: ID!) {
    trace(id: $id) {
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
`;
