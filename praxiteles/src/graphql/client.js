// src/graphql/client.js
import { ApolloClient, InMemoryCache } from '@apollo/client/core';
import { HttpLink } from '@apollo/client/link/http';

const graphqlUrl = process.env.GRAPHQL_URL ||  document.location.origin + '/graphql';

const httpLink = new HttpLink({
    // You should use an absolute URL here
    uri: graphqlUrl,
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    headers: {
        "Content-Type": "application/json",
    },
    redirect: "follow",
    referrerPolicy: "no-referrer",
})


const client = new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
    connectToDevTools: true,
})

export default client;