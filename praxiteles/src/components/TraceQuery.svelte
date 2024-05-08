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
    let ganntContent = null;
    let stateChartContent = null;
    let podMetricsMemory = null;
    let podMetricsCpu = null;

    // Function to handle the query result
    function handleQueryResult(result, filters) {
        queryResult = result.data; // Assign the result to queryResult
        if (result.data.traces && result.data.traces.length > 0) {
            // Handle other diagram generation that does not depend on podName
            let jsonResult = JSON.stringify(result.data.traces[0], null, 2);
            ganntContent = convertTraceToGannt(jsonResult);
            stateChartContent = convertTraceToStateDiagram(jsonResult);
        }

        if (filters.podName) {
            let allItems = result.data.traces.flatMap(trace => trace.items);

            let filteredItems = allItems.filter(item =>
                item.itemType === 'trace' &&
                item.podName === filters.podName &&
                'metrics' in item &&
                item.metrics != null
            );

            let sortedItems = filteredItems.sort((a, b) => {
                const dateA = new Date(a.timestamp.replace('T', ' '));
                const dateB = new Date(b.timestamp.replace('T', ' '));
                return dateA - dateB;
            });

            // Generate line chart if there are traces for the specified podName
            if (sortedItems.length > 0) {
                createLineDiagram(sortedItems); // Generate line chart
            }
        }
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
        handleQueryResult(result, filters);
    }

    // Listen for the custom event and handle it
    document.addEventListener('searchfiltersubmit', (event) => {
        // Access the filter criteria from the event's detail property
        const filters = event.detail;

        // Call the Promise-based function to handle the GraphQL query
        handleQueryWithFilters(filters);
    });
    function convertTraceToStateDiagram(traceData) {
        let traceDataParsed = JSON.parse(traceData);

        let diagram = `---\nState Diagram for ${traceDataParsed.traceID}\n---\n`;
        diagram += 'stateDiagram-v2\n';

        // Create a mapping of spanID to podName for all traces
        let traceIdToPodName = {};
        traceDataParsed.items.filter(item => item.itemType === 'trace').forEach(trace => {
            traceIdToPodName[trace.spanID] = trace.podName.split('-')[0];
        });

        let firstServicePodName = "";
        // Initialize the diagram with the first service
        if (traceDataParsed.items.length > 0 && traceDataParsed.items[0].itemType === 'trace_start') {
            firstServicePodName = traceDataParsed.items[0].podName.split('-')[0];
            traceIdToPodName[traceDataParsed.items[0].spanID] = firstServicePodName;
            if (firstServicePodName) {
                diagram += `    [*] --> ${firstServicePodName}\n`;
            }
        }

        // Iterate through traces to create transitions based on parentSpanID -> spanID mapping
        traceDataParsed.items.filter(item => item.itemType === 'trace').forEach(trace => {
            const parentService = traceIdToPodName[trace.parentSpanID];
            const childService = trace.podName.split('-')[0]; // Directly use trace podName if no spanID -> podName mapping

            // Add transition if parentService is identified
            if (parentService && childService && parentService !== childService) {
                const transition = `    ${parentService} --> ${childService}`;
                diagram += `${transition}\n`;

            }
        });

        return diagram;
    }
    function convertTraceToGannt(traceData) {
        let traceDataParsed = JSON.parse(traceData);

        let traceStartItem = traceDataParsed.items.find(item => item.itemType === 'trace_start');
        if (!traceStartItem) {
            return; // Exit if no trace_start item is found
        }

        let diagram = 'gantt\n';
        diagram += `dateFormat YYYY-MM-DD'T'HH:mm:ss.SSS\n`; // Set the date format
        diagram += `axisFormat %Y-%m-%d | %H:%M:%S.%L\n`
        diagram += `title Trace Gantt Chart for ${traceDataParsed.traceID}\n`;
        diagram += `section Parent\n`;
        diagram += `Overall Timeline :milestone, ${traceDataParsed.timeStarted}, ${traceDataParsed.timeEnded}\n`;

        let parsedEndTime = traceDataParsed.timeEnded.replace('T', "'T'");
        diagram += `section Trace\n`;
        diagram += `    ${traceStartItem.podName} :${1}, ${traceStartItem.timestamp}, ${parsedEndTime}\n`;

        let traceItems = traceDataParsed.items.filter(item => item.itemType === 'trace')
            .sort((a, b) => {
                const dateA = new Date(a.timestamp.replace("'T'", " "));
                const dateB = new Date(b.timestamp.replace("'T'", " "));
                return dateA - dateB;
            });

        traceItems.forEach((item, index) => {
            const taskEnd = index < traceItems.length - 1 ? traceItems[index + 1].timestamp : traceDataParsed.timeEnded;
            diagram += `    ${item.podName} :${index + 2}, ${item.timestamp}, ${taskEnd}\n`;
        });

        let spanItems = traceDataParsed.items.filter(item => item.itemType === 'span')
            .map(item => ({
                ...item,
                startTime: calculateStartTime(item.timestamp, item.took)
            }))
            .sort((a, b) => {
                const dateA = new Date(a.startTime.replace("'T'", " "));
                const dateB = new Date(b.startTime.replace("'T'", " "));
                return dateA - dateB;
            });

        diagram += `section Span\n`;
        spanItems.forEach((item, index) => {
            diagram += `    ${item.spanID} - ${item.podName} (${item.parentSpanID}):${index + 1}, ${item.startTime}, ${item.timestamp}\n`;
        });

        let databaseItems = traceDataParsed.items.filter(item => item.itemType === 'database_span')
            .map(item => ({
                ...item,
                startTime: calculateStartTime(item.timestamp, item.took)
            }))
            .sort((a, b) => {
                const dateA = new Date(a.startTime.replace("'T'", " "));
                const dateB = new Date(b.startTime.replace("'T'", " "));
                return dateA - dateB;
            });

        diagram += `section Database\n`;
        databaseItems.forEach((item, index) => {
            diagram += `    ${item.podName} :${index + 1}, ${item.startTime}, ${item.timestamp}\n`;
        });

        return diagram;
    }

    // Helper function to calculate start time
    function calculateStartTime(startTime, duration) {
        if (!duration) {
            return startTime;
        }
        const correctedTimestamp = startTime.replace(/'/g, '');
        const durationMs = parseFloat(duration.replace('ms', ''));
        const endDateTime = new Date(correctedTimestamp);
        const startDateTime = new Date(endDateTime.getTime() - durationMs);

        return `${startDateTime.getFullYear()}-${('0' + (startDateTime.getMonth() + 1)).slice(-2)}-${('0' + startDateTime.getDate()).slice(-2)}'T'${('0' + startDateTime.getHours()).slice(-2)}:${('0' + startDateTime.getMinutes()).slice(-2)}:${('0' + startDateTime.getSeconds()).slice(-2)}.${('00' + startDateTime.getMilliseconds()).slice(-3)}`;
    }

    function createLineDiagram(filteredTraces) {
        let timestamps = filteredTraces.map((_, index) => index + 1);
        let cpuValues = filteredTraces.map(trace => trace.metrics.cpuRaw);
        let memoryValues = filteredTraces.map(trace => trace.metrics.memoryRaw);

        let startTime = new Date(filteredTraces[0].timestamp.replace("'T'", " "));
        let endTime = new Date(filteredTraces[filteredTraces.length - 1].timestamp.replace("'T'", " "));

        let cpuDiagram = `xychart-beta\n`;
        cpuDiagram += `    title "CPU usage for pod: ${filteredTraces[0].podName} from ${startTime} to ${endTime}"\n`;
        cpuDiagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
        cpuDiagram += '    y-axis "CPU in Units"\n';
        cpuDiagram += `    bar [${cpuValues.join(', ')}]\n`;

        let memoryDiagram = `xychart-beta\n`;
        memoryDiagram += `    title "Memory usage for pod: ${filteredTraces[0].podName} from ${startTime.toLocaleString()} to ${endTime.toLocaleString()}"\n`;
        memoryDiagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
        memoryDiagram += '    y-axis "Memory in Mi"\n';
        memoryDiagram += `    bar [${memoryValues.join(', ')}]\n`;

        podMetricsMemory = memoryDiagram;
        podMetricsCpu = cpuDiagram;

    }


</script>

<style>
    .traces-container,
    .visual-container {
        padding: 1em; /* Add padding in em units for spacing */
    }

    .traces-container {
        margin-left: 10em;
        width: 60%;
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
        <h2 id="traces-visual">Diagrams</h2>
        {#if ganntContent}
            <Diagram mermaidDiagram={ganntContent} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
        {#if stateChartContent}
            <Diagram mermaidDiagram={stateChartContent} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
        {#if podMetricsCpu}
            <Diagram mermaidDiagram={podMetricsCpu} />
        {/if}
        {#if podMetricsMemory}
            <Diagram mermaidDiagram={podMetricsMemory} />
        {/if}
    </div>
</div>
