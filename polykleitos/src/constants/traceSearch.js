// src/constants/traceSearch.js
import gql from "graphql-tag";

export const TraceSearchQuery = gql`
  query traceSearch($input: TraceSearchInput!) {
    traceSearch(input: $input) {
      total
      items {
        id
        numberOfItems
        rootQuery
        isActive
        hasDbSpan
        responseCode
        totalTimeMs
        timeStarted
        timeEnded
      }
    }
  }
`;