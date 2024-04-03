<!-- TraceQuery.svelte -->
<script>
    // Import necessary dependencies for GraphQL
    import { gql } from '@apollo/client';
    import traceQuery from '../graphql/queries/getTraces.graphql';
    import JSONViewWrapper from './JSONViewWrapper.svelte'; // Import the JSONViewWrapper component
    import SearchFilter from "./SearchFilter.svelte";
    import client from "../graphql/client.js";
    import Diagram from "./Diagram.svelte"
    // Define your GraphQL query
    const GET_TRACES = gql`${traceQuery}`;

    // Declare a variable to store the query result
    let queryResult = null; // Declare queryResult here
    let ganntContent = ``

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

    function convertTraceToGannt(traceData) {
        let traceDataParsed = JSON.parse(traceData);

        let diagram = 'gantt\n';
        diagram += `dateFormat YYYY-MM-DD'T'HH:mm:ss.SSS\n`; // Set the date format to match your timestamp format
        diagram += `axisFormat %Y-%m-%d | %H:%M:%S.%L\n`
        diagram += `title Trace Gantt Chart for ${traceDataParsed.traceID}\n`;

        // Adding a section for the overall ParentTrace timeline
        diagram += `section Parent\n`;
        diagram += `Overall Timeline :milestone, ${traceDataParsed.timeStarted}, ${traceDataParsed.timeEnded}\n`;

        // Filter out 'span' items and sort by 'timestamp'
        let traceItems = traceDataParsed.items.filter(item => item.itemType === 'trace')
            .sort((a, b) => {
                // Replace 'T' with a space for correct date parsing
                const dateA = new Date(a.timestamp.replace('T', ' '));
                const dateB = new Date(b.timestamp.replace('T', ' '));
                return dateA - dateB;
            });

        // Determine the overall start and end times for the Trace section
        let traceEnd = traceItems[traceItems.length - 1].timestamp;
        if (new Date(traceDataParsed.timeEnded) > new Date(traceEnd)) {
            traceEnd = traceDataParsed.timeEnded; // Use parent's end time if it's later
        }

        // Adding the main section for Trace items
        diagram += `section Trace\n`;

        traceItems.forEach((item, index) => {
            const taskStart = item.timestamp;
            const taskEnd = index < traceItems.length - 1 ? traceItems[index + 1].timestamp : traceEnd; // Use next item's timestamp or the overall trace end time
            const taskLabel = item.podName;

            diagram += `    ${taskLabel} :${index + 1}, ${taskStart}, ${taskEnd}\n`;
        });

        return diagram;
    }

    $: if (queryResult && queryResult.traces && queryResult.traces.length > 0) {
        ganntContent = convertTraceToGannt(JSON.stringify(queryResult.traces[0], null, 2));
    } else {
        ganntContent = '';
    }

</script>

<style>
    .container {
    }

    .traces-container,
    .visual-container {
        padding: 1em; /* Add padding in em units for spacing */
    }

    .traces-container {
        margin-left: 10em;
        width: 60%;
    }

    .visual-container {
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
    <div class="visual-container">
        <h2 id="visual">FlowChart</h2>
        {#if ganntContent}
            <Diagram mermaidDiagram={ganntContent} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
    </div>
</div>
