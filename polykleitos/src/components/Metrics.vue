<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="6">
        <h2>Metrics</h2>

        <v-card class="mb-4">
          <v-card-text>
            <h4>Order:</h4>
            <v-radio-group v-model="order" inline>
              <v-radio label="Descending" value="desc"></v-radio>
              <v-radio label="Ascending" value="asc"></v-radio>
            </v-radio-group>

            <h4>Time Span:</h4>
            <v-radio-group v-model="timeSpan" inline>
              <v-radio label="Hour" value="hour"></v-radio>
              <v-radio label="3 Hours" value="3hour"></v-radio>
              <v-radio label="6 Hours" value="6hour"></v-radio>
              <v-radio label="12 Hours" value="12hour"></v-radio>
              <v-radio label="24 Hours" value="day"></v-radio>
              <v-radio label="Week" value="week"></v-radio>
              <v-radio label="Month" value="month"></v-radio>
            </v-radio-group>
          </v-card-text>
        </v-card>

        <div v-if="queryResult">
          <h3>Results:</h3>
          <JSONViewWrapper :jsonString="JSON.stringify(queryResult.metrics, null, 2)" />
        </div>
        <div v-else>
          <p>Waiting for metrics...</p>
        </div>
      </v-col>

      <v-col cols="12" md="6">
        <h2>Diagrams</h2>

        <v-card v-if="pieChartCpu" class="mb-4">
          <v-card-title>CPU Pie Chart</v-card-title>
          <v-card-text>
            <v-select
              v-model="selectedTimeStamp"
              :items="timeStamps"
              label="Select a timestamp"
              clearable
            ></v-select>
            <Diagram :mermaidDiagram="pieChartCpu" />
          </v-card-text>
        </v-card>

        <v-card v-if="pieChartMemory" class="mb-4">
          <v-card-title>Memory Pie Chart</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="pieChartMemory" />
          </v-card-text>
        </v-card>

        <v-card v-if="barChartCpu" class="mb-4">
          <v-card-title>CPU Usage per Node</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="barChartCpu" />
          </v-card-text>
        </v-card>

        <v-card v-if="barChartMemory" class="mb-4">
          <v-card-title>Memory Usage per Node</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="barChartMemory" />
          </v-card-text>
        </v-card>

        <div v-if="!pieChartCpu">
          <p>No data available for visualization.</p>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { gql } from '@apollo/client/core';
import client from '../graphql/client.js';
import JSONViewWrapper from './JSONViewWrapper.vue';
import Diagram from './Diagram.vue';

const queryResult = ref(null);
const rawResult = ref(null);
const pieChartMemory = ref(null);
const pieChartCpu = ref(null);
const barChartMemory = ref(null);
const barChartCpu = ref(null);

const order = ref('desc');
const timeSpan = ref('hour');
const selectedTimeStamp = ref('');
const timeStamps = ref([]);

const GET_METRICS = gql`
  query GetMetrics($order: String!, $timeSpan: String!) {
    metrics(order: $order, timeSpan: $timeSpan) {
      nodes {
        nodeName
        cpuPercentageHumanReadable
        memoryPercentageHumanReadable
      }
      grouped {
        name
        cpuHumanReadable
        memoryHumanReadable
      }
      pods {
        name
        cpuHumanReadable
        memoryHumanReadable
      }
    }
  }
`;

const GET_METRICS_RAW = gql`
  query GetMetrics($order: String!, $timeSpan: String!) {
    metrics(order: $order, timeSpan: $timeSpan) {
      timeStarted
      timeEnded
      timeStamps
      nodes {
        nodeName
        cpuRaw
        memoryRaw
        cpuPercentage
        memoryPercentage
      }
      grouped {
        name
        cpuRaw
        memoryRaw
      }
      pods {
        name
        cpuRaw
        memoryRaw
      }
    }
  }
`;

watch([order, timeSpan], () => {
  fetchData();
});

watch(selectedTimeStamp, (newVal) => {
  if (newVal) {
    updatePieCharts(newVal);
  }
});

async function fetchData() {
  try {
    const variables = {
      order: order.value,
      timeSpan: timeSpan.value,
    };

    const result = await client.query({
      query: GET_METRICS,
      variables,
    });
    queryResult.value = result.data;

    const rawResultQuery = await client.query({
      query: GET_METRICS_RAW,
      variables,
    });
    rawResult.value = rawResultQuery.data;
    timeStamps.value = rawResult.value.metrics.timeStamps;

    convertToPieChart(rawResult.value);
    convertToMemoryChart(rawResult.value);
    convertToCpuChart(rawResult.value);
  } catch (error) {
    console.error('Error fetching metrics:', error);
  }
}

function updatePieCharts(timeStamp) {
  const timestampIndex = timeStamps.value.findIndex(ts => ts === timeStamp);
  convertToPieChart(rawResult.value, timestampIndex);
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

  pieChartMemory.value = diagramMemory;
  pieChartCpu.value = diagramCpu;
}

function convertToCpuChart(rawResult) {
  const timestamps = formatTimestamps(rawResult.metrics.timeStamps);
  let diagram = `xychart-beta\n`;
  diagram += '    title "CPU usage per Node (in CPU units)"\n';
  diagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
  diagram += '    y-axis "CPU Units"\n';

  rawResult.metrics.nodes.forEach(node => {
    const reversedCpuData = node.cpuRaw.slice().reverse();
    diagram += `    line [${reversedCpuData.join(', ')}]\n`;
  });

  barChartCpu.value = diagram;
}

function convertToMemoryChart(rawResult) {
  const timestamps = formatTimestamps(rawResult.metrics.timeStamps);
  let diagram = `xychart-beta\n`;
  diagram += '    title "Memory usage per Node (in MiB)"\n';
  diagram += `    x-axis ${JSON.stringify(timestamps)}\n`;
  diagram += '    y-axis "Memory in MiB"\n';

  rawResult.metrics.nodes.forEach(node => {
    const reversedMemoryData = node.memoryRaw.slice().reverse();
    diagram += `    line [${reversedMemoryData.join(', ')}]\n`;
  });

  barChartMemory.value = diagram;
}

function formatTimestamps(timestamps) {
  const reversedTimestamps = timestamps.slice().reverse();
  return reversedTimestamps.map(ts => {
    const date = new Date(ts);
    const hours = date.getHours();
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
  });
}

onMounted(() => {
  fetchData();
});
</script>
