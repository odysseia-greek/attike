import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client/core';
import { provideApolloClient, DefaultApolloClient } from '@vue/apollo-composable';

// Define the URL for the GraphQL server
let url = document.location.origin + '/euripides/graphql';
if (process.env.NODE_ENV === 'development') {
    url = 'http://localhost:8080/euripides/graphql';
}

if (process.env.NODE_ENV === 'local') {
    url = 'http://attike.byzantium.odysseia-greek:8080/euripides/graphql';
}

const inMemoryCache = new InMemoryCache({});

// Create the HTTP link for the Apollo client
const httpLink = new HttpLink({
    uri: url,
    method: 'POST',
    mode: 'cors',
    headers: {
        'Content-Type': 'application/json',
    },
    credentials: 'same-origin',
});

// Create the Apollo client
export const apolloClient = new ApolloClient({
    link: httpLink,
    cache: inMemoryCache,
    connectToDevTools: true,
});

// Provide the Apollo client for use with the Composition API
provideApolloClient(apolloClient);

export default {
    install: (app) => {
        app.provide(DefaultApolloClient, apolloClient);
    }
};
