import gql from "graphql-tag";

export const TracePollQuery = gql`
  query tracePoll($limit: Int!) {
    tracePoll(limit: $limit) {
      updatedAt
      traces {
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