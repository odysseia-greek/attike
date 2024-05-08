<script>
    // Import necessary dependencies for GraphQL
    import { gql } from '@apollo/client';
    import metricsQuery from '../graphql/queries/getMetrics.graphql';
    import metricsQueryRaw from '../graphql/queries/getMetricsRaw.graphql';
    import JSONViewWrapper from './JSONViewWrapper.svelte'; // Import the JSONViewWrapper component
    import client from "../graphql/client.js";
    import {onMount} from "svelte";
    import Diagram from "./Diagram.svelte";
    // Define your GraphQL query
    const GET_METRICS = gql`${metricsQuery}`;
    const GET_METRICS_RAW =gql`${metricsQueryRaw}`;

    // Declare a variable to store the query result
    let queryResult = null; // Declare queryResult here
    let rawResult = null; // set the raw result
    let pieChartMemory = null;
    let pieChartCpu = null;
    let barChartMemory = null;
    let barChartCpu = null;

    let order = "desc"; // Default order
    let timeSpan = "hour"; // Default time span
    let selectedTimeStamp = '';
    let timeStamps = [];

    // Reactive statement to update data based on user input
    $: {
        fetchData(order, timeSpan);
    }

    $: if (selectedTimeStamp) {
        updatePieCharts(selectedTimeStamp);
    }

    async function fetchData(order, timeSpan) {
        console.log(order, timeSpan)
        const variables = {
            order: order,
            timeSpan: timeSpan,
        };

        const result = await client.query({
            query: GET_METRICS,
            variables,
        });
        queryResult = result.data;

        const rawResultQuery = await client.query({
            query: GET_METRICS_RAW,
            variables,
        });
        rawResult = rawResultQuery.data;
        timeStamps = rawResult.metrics.timeStamps;

        convertToPieChart(rawResult);
        convertToMemoryChart(rawResult);
        convertToCpuChart(rawResult);
    }

    function updatePieCharts(timeStamp) {
        const timestampIndex = timeStamps.findIndex(ts => ts === timeStamp);
        convertToPieChart(rawResult, timestampIndex);
    }

    function convertToPieChart(rawResult, timestampIndex = 0) {
        let diagramMemory = 'pie title Pie chart for Memory\n';
        let diagramCpu = 'pie title Pie chart for CPU\n';

        if (rawResult && rawResult.metrics && rawResult.metrics.grouped) {
            rawResult.metrics.grouped.forEach(item => {
                const memoryValue = item.memoryRaw && item.memoryRaw[timestampIndex] ? item.memoryRaw[timestampIndex] : 0;
                const cpuValue = item.cpuRaw && item.cpuRaw[timestampIndex] ? item.cpuRaw[timestampIndex] : 0;

                if (cpuValue !== 0) {
                    diagramCpu += `    "${item.name || 'Unknown'} - ${cpuValue}Mi" : ${cpuValue}\n`;
                }
                if (memoryValue !== 0) {
                    diagramMemory += `    "${item.name || 'Unknown'} - ${memoryValue}Mbi" : ${memoryValue}\n`;
                }
            });
        }

        pieChartMemory = diagramMemory;
        pieChartCpu = diagramCpu;
    }


    function convertToCpuChart(rawResult) {
        let timestamps = formatTimestamps(rawResult.metrics.timeStamps);
        let diagram = `xychart-beta\n`;
        diagram += '    title "CPU usage per Node (in CPU units)"\n';
        diagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
        diagram += '    y-axis "CPU Units"\n';

        rawResult.metrics.nodes.forEach(node => {
            // Reverse the cpuRaw array
            let reversedCpuData = node.cpuRaw.slice().reverse();
            diagram += `    line [${reversedCpuData.join(', ')}]\n`;
        });

        barChartCpu = diagram;
    }

    function convertToMemoryChart(rawResult) {
        let timestamps = formatTimestamps(rawResult.metrics.timeStamps);
        let diagram = `xychart-beta\n`;
        diagram += '    title "Memory usage per Node (in MiB)"\n';
        diagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
        diagram += '    y-axis "Memory in MiB"\n';

        rawResult.metrics.nodes.forEach(node => {
            // Reverse the memoryRaw array
            let reversedMemoryData = node.memoryRaw.slice().reverse();
            diagram += `    line [${reversedMemoryData.join(', ')}]\n`;
        });

        barChartMemory = diagram;
    }

    function formatTimestamps(timestamps) {
        let reversedTimestamps = timestamps.slice().reverse();
        return reversedTimestamps.map(ts => {
            let date = new Date(ts);
            // Use padStart to ensure minutes are always two digits
            let hours = date.getHours();
            let minutes = date.getMinutes().toString().padStart(2, '0');
            return `${hours}:${minutes}`;
        });
    }

    onMount(async () => {
        await fetchData(order, timeSpan)
    });


</script>

<style>
    .metrics-container,
    .visual-container {
        padding: 1em; /* Add padding in em units for spacing */
    }

    .metrics-container {
        margin-left: 10em;
        width: 60%;
    }


</style>

<div class="container">
    <div class="metrics-container">
        <h2 id="metrics">Metrics</h2>
        <h4>Order:</h4>
        <label>
            <input type="radio" bind:group={order} value="desc"> Descending
        </label>
        <label>
            <input type="radio" bind:group={order} value="asc"> Ascending
        </label>

        <!-- Inserted Radio Buttons for 'timeSpan' -->
        <h4>Time Span:</h4>
        <label>
            <input type="radio" bind:group={timeSpan} value="hour"> Hour
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="3hour"> 3 Hours
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="6hour"> 6 Hours
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="12hour"> 12 Hours
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="day"> 24 Hours
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="week"> Week
        </label>
        <label>
            <input type="radio" bind:group={timeSpan} value="month"> Month
        </label>
        {#if queryResult}
            <h3>Results:</h3>
            <JSONViewWrapper jsonString={JSON.stringify(queryResult.metrics, null, 2)} />
        {:else}
            <p>Waiting for metrics...</p>
        {/if}
    </div>
    <div class="visual-container">
        <h2 id="metrics-visual">Diagrams</h2>
        {#if pieChartCpu}
            <div class="options-dropdown">
                <select bind:value={selectedTimeStamp}>
                    <option value="">Select a timeStamp</option>
                    {#each timeStamps as timeStamp }
                        <option value={timeStamp}>{timeStamp}</option>
                    {/each}
                </select>
            </div>
            <Diagram mermaidDiagram={pieChartCpu} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
        {#if pieChartMemory}
            <Diagram mermaidDiagram={pieChartMemory} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
        {#if barChartCpu}
            <Diagram mermaidDiagram={barChartCpu} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
        {#if barChartMemory}
            <Diagram mermaidDiagram={barChartMemory} />
        {:else}
            <p>No data available for visualization.</p>
        {/if}
    </div>
</div>
