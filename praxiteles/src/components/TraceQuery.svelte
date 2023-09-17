<!-- TraceQuery.svelte -->
<script>
    // Import necessary dependencies for GraphQL
    import { gql } from '@apollo/client';
    import traceQuery from '../graphql/queries/getTraces.graphql';
    import JSONViewWrapper from './JSONViewWrapper.svelte'; // Import the JSONViewWrapper component
    import SearchFilter from "./SearchFilter.svelte";
    import client from "../graphql/client.js"; // Import the client

    // Define your GraphQL query
    const GET_TRACES = gql`${traceQuery}`;

    // Declare a variable to store the query result
    let queryResult = null; // Declare queryResult here

    // Function to handle the query result
    function handleQueryResult(result) {
        // The query result is available in the result variable
        // You can access data using result.data
        queryResult = result.data; // Assign the result to queryResult
    }

    // Create a Promise-based function to handle the GraphQL query
    async function handleQueryWithFilters(filters) {
        // Use the filter criteria to update the GraphQL query variables
        const variables = {
            input: filters,
        };

        // Trigger the GraphQL query with the updated variables
        const result = await client.query({
            query: GET_TRACES,
            variables: variables,
        });

        // Call the function to handle the query result
        handleQueryResult(result);
    }

    // Listen for the custom event and handle it
    document.addEventListener('searchfiltersubmit', (event) => {
        // Access the filter criteria from the event's detail property
        const filters = event.detail;

        // Call the Promise-based function to handle the GraphQL query
        handleQueryWithFilters(filters);
    });

</script>

<style>
    .container {
        display: flex; /* Use flexbox layout */
        flex-direction: row; /* By default, display containers side by side */
    }

    .traces-container,
    .metrics-container {
        padding: 1em; /* Add padding in em units for spacing */
        border-right: 0.0625em solid #ccc; /* Add a border in em units to separate containers */
    }

    .traces-container {
        flex: 1; /* Take up 50% of the screen */
    }

    .metrics-container {
        flex: 1; /* Take up 50% of the screen */
    }

    /* Media query for screens larger than a tablet (min-width: 51.25em) */
    @media (min-width: 51.25em) {
        .container {
            flex-direction: row; /* Display containers side by side */
        }

        .traces-container,
        .metrics-container {
            flex: 1; /* Take up 50% of the screen */
            border-right: 0.0625em solid #ccc; /* Add a border in em units to separate containers */
        }
    }

    /* Media query for tablet screens (max-width: 51.25em) */
    @media (max-width: 51.25em) {
        .container {
            flex-direction: column; /* Stack containers vertically for smaller screens */
        }

        .traces-container,
        .metrics-container {
            flex: initial; /* Reset flex value to allow natural width */
            width: 100%; /* Expand to full width of the container */
            border-right: none; /* Remove the border for smaller screens */
        }
    }
</style>

<div class="container">
    <div class="traces-container">
        <h2 id="traces">Traces</h2>
        <SearchFilter />
        <!-- Display the query result data here -->
        {#if queryResult}
            <h3>Results:</h3>
            {#each queryResult.traces as trace, index}
                <JSONViewWrapper jsonString={JSON.stringify(trace, null, 2)} />
            {/each}
        {:else}
            <p>Enter query...</p>
        {/if}
    </div>
    <div class="metrics-container">
        <h2 id="metrics">Metrics</h2>
        <!-- You can add content to the metrics container here -->
    </div>
</div>
