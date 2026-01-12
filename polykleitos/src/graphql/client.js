// src/graphql/client.js
import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client/core';

let graphqlUrl = document.location.origin + '/graphql'
if (import.meta.env.MODE === 'development') {
    graphqlUrl = "http://localhost:8080/graphql"
}

if (import.meta.env.VITE_ENV === 'k3d') {
    graphqlUrl = "http://k3d-odysseia.api.greek:8080/graphql"
}

const httpLink = new HttpLink({
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
