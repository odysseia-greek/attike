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
    let pieChartMemory = null;
    let pieChartCpu = null;
    let barChartMemory = null;
    let barChartCpu = null;


    function convertToPieChart(rawResult) {
        let diagramMemory = 'pie title Pie chart for Memory\n';
        let diagramCpu = 'pie title Pie chart for CPU\n';

        if (rawResult && rawResult.metrics && rawResult.metrics.grouped) {
            rawResult.metrics.grouped.forEach((item, index) => {
                const memoryValue = (item.memoryRaw && item.memoryRaw.length > 0) ? item.memoryRaw[0] : 0;
                const cpuValue = (item.cpuRaw && item.cpuRaw.length > 0) ? item.cpuRaw[0] : 0;

                // Add to pie chart data only if there's a non-zero value to display
                if (cpuValue !== 0) {
                    diagramCpu += `    "${item.name || 'Unknown'} - ${cpuValue}Mi" : ${cpuValue}\n`;
                }
                if (memoryValue !== 0) {
                    diagramMemory += `    "${item.name || 'Unknown'} - ${memoryValue}Mbi" : ${memoryValue}\n`;
                }
            });
        }

        pieChartMemory = diagramMemory
        pieChartCpu = diagramCpu
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
        // Trigger the GraphQL query with the updated variables
        const variables = {
            input: {},
        };

        // Trigger the GraphQL query with the updated variables
        const result = await client.query({
            query: GET_METRICS,
            variables: variables,
        });

        queryResult = result.data;

        const rawResultQuery = await client.query({
            query: GET_METRICS_RAW,
            variables: variables,
        })

        const rawResult = rawResultQuery.data;
        convertToPieChart(rawResult)
        convertToMemoryChart(rawResult)
        convertToCpuChart(rawResult)
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
