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

        // Create a mapping of spanID to podName for all spans
        let spanIdToPodName = {};
        traceDataParsed.items.filter(item => item.itemType === 'span').forEach(span => {
            spanIdToPodName[span.spanID] = span.podName.split('-')[0];
        });

        let firstServicePodName = "";
        // Initialize the diagram with the first service
        if (traceDataParsed.items.length > 0 && traceDataParsed.items[0].itemType === 'trace') {
            firstServicePodName = traceDataParsed.items[0].podName.split('-')[0];
            if (firstServicePodName) {
                diagram += `    [*] --> ${firstServicePodName}\n`;
            }
        }

        let firstTransitionAdded = false;

        // Iterate through traces to create transitions based on parentSpanID -> spanID mapping
        traceDataParsed.items.filter(item => item.itemType === 'trace').forEach(trace => {
            const parentService = spanIdToPodName[trace.parentSpanID];
            const childService = trace.podName.split('-')[0]; // Directly use trace podName if no spanID -> podName mapping

            // Add transition if parentService is identified
            if (parentService && childService) {
                const transition = `    ${parentService} --> ${childService}`;
                diagram += `${transition}\n`;

                // Add the missing link between the firstServicePodName and the first parentService
                if (!firstTransitionAdded && firstServicePodName && firstServicePodName !== childService) {
                    diagram += `    ${firstServicePodName} --> ${parentService}\n`;
                    firstTransitionAdded = true; // Ensure this transition is only added once
                }
            }

            if (!firstTransitionAdded && parentService === undefined  && firstServicePodName !== childService) {
                diagram += `    ${firstServicePodName} --> ${childService}\n`;
                firstTransitionAdded = true; // Ensure this transition is only added once
            }
        });

        console.log(diagram)
        return diagram;
    }

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

        let spanItems = traceDataParsed.items.filter(item => item.itemType === 'span');

        const groupedSpans = spanItems.reduce((acc, item) => {
            if (!acc[item.spanID]) {
                acc[item.spanID] = { start: null, end: null, podName: item.podName, parentSpanID: item.parentSpanID };
            }

            // Assign start time from the timestamp
            if (item.timestamp) {
                acc[item.spanID].start = item.timestamp;

                // Check if the 'took' field is present, then calculate and set the end time
                if (item.took) {
                    const correctedTimestamp = item.timestamp.replace(/'/g, '');
                    const durationMs = parseFloat(item.took.replace('ms', ''));
                    // Parse the timestamp to a Date object and add the duration
                    const startTime = new Date(correctedTimestamp);
                    const endTime = new Date(startTime.getTime() + durationMs);

                    // Convert the end time back to ISO 8601 format string
                    let localEndTime = `${endTime.getFullYear()}-${('0' + (endTime.getMonth()+1)).slice(-2)}-${('0' + endTime.getDate()).slice(-2)}'T'${('0' + endTime.getHours()).slice(-2)}:${('0' + endTime.getMinutes()).slice(-2)}:${('0' + endTime.getSeconds()).slice(-2)}.${('00' + endTime.getMilliseconds()).slice(-3)}`;
                    acc[item.spanID].end = localEndTime;
                }
            }

            return acc;
        }, {});

        const mergedSpans = Object.entries(groupedSpans).map(([spanID, { start, end, podName, parentSpanID }]) => ({
            spanID,
            podName,
            parentSpanID,
            timeStarted: start,
            timeFinished: end
        }));

        diagram += `section Span\n`;
        mergedSpans.forEach((item, index) => {
            const taskLabel = `${item.spanID} - ${item.podName}`;
            // Now including parentID and the podName of the starter in the output
            diagram += `    ${taskLabel} (${item.parentSpanID}):${index + 1}, ${item.timeStarted }, ${item.timeFinished}\n`;
        });

        diagram += `section Database\n`;
        let dataBaseItem = traceDataParsed.items.filter(item => item.itemType === 'database_span')
        dataBaseItem.forEach((item, index) => {
            const taskStart = item.timestamp;
            let taskEnd = null
            const taskLabel = item.podName;

            // Check if the 'took' field is present, then calculate and set the end time
            if (item.took) {
                const correctedTimestamp = item.timestamp.replace(/'/g, '');
                const durationMs = parseFloat(item.took.replace('ms', ''));
                // Parse the timestamp to a Date object and add the duration
                const startTime = new Date(correctedTimestamp);
                const endTime = new Date(startTime.getTime() + durationMs);

                // Convert the end time back to ISO 8601 format string
                taskEnd = `${endTime.getFullYear()}-${('0' + (endTime.getMonth()+1)).slice(-2)}-${('0' + endTime.getDate()).slice(-2)}'T'${('0' + endTime.getHours()).slice(-2)}:${('0' + endTime.getMinutes()).slice(-2)}:${('0' + endTime.getSeconds()).slice(-2)}.${('00' + endTime.getMilliseconds()).slice(-3)}`;
                diagram += `    ${taskLabel} :${index + 1}, ${taskStart}, ${taskEnd}\n`;
        }
        });

        return diagram;
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
