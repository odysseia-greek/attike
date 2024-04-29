// src/graphql/client.js
import { ApolloClient, InMemoryCache } from '@apollo/client/core';
import { HttpLink } from '@apollo/client/link/http';

let graphqlUrl = document.location.origin + '/graphql'
if (process.env.ENV === 'development') {
    graphqlUrl = "http://localhost:8080/graphql"
}

if (process.env.ENV === 'k3d') {
    graphqlUrl = "http://k3d-odysseia.api.greek:8080/graphql"
}

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