<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="6">
        <h2>Traces</h2>
        <SearchFilter @search="handleSearch" />

        <div v-if="queryResult">
          <h3>Results:</h3>
          <div v-for="(trace, index) in queryResult.traces" :key="index">
            <JSONViewWrapper :jsonString="JSON.stringify(trace, null, 2)" />
          </div>
        </div>
        <div v-else>
          <p>Enter query...</p>
        </div>
      </v-col>

      <v-col cols="12" md="6">
        <h2>Diagrams</h2>
        <v-card v-if="ganttContent" class="mb-4">
          <v-card-title>Gantt Chart</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="ganttContent" />
          </v-card-text>
        </v-card>

        <v-card v-if="stateChartContent" class="mb-4">
          <v-card-title>State Diagram</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="stateChartContent" />
          </v-card-text>
        </v-card>

        <v-card v-if="podMetricsCpu" class="mb-4">
          <v-card-title>CPU Metrics</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="podMetricsCpu" />
          </v-card-text>
        </v-card>

        <v-card v-if="podMetricsMemory" class="mb-4">
          <v-card-title>Memory Metrics</v-card-title>
          <v-card-text>
            <Diagram :mermaidDiagram="podMetricsMemory" />
          </v-card-text>
        </v-card>

        <div v-if="!ganttContent && !stateChartContent">
          <p>No data available for visualization.</p>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';
import { gql } from '@apollo/client/core';
import client from '../graphql/client.js';
import SearchFilter from './SearchFilter.vue';
import JSONViewWrapper from './JSONViewWrapper.vue';
import Diagram from './Diagram.vue';

const queryResult = ref(null);
const ganttContent = ref(null);
const stateChartContent = ref(null);
const podMetricsMemory = ref(null);
const podMetricsCpu = ref(null);

const GET_TRACES = gql`
  query GetTraces($input: TraceQueryInput!) {
    traces(input: $input) {
      traceID
      isActive
      timeEnded
      timeStarted
      totalTime
      responseCode
      items {
        ... on StartTrace {
          parentSpanID
          spanID
          method
          url
          host
          remoteAddress
          timestamp
          podName
          namespace
          itemType
          operation
          rootQuery
          metrics {
            cpuHumanReadable
            memoryHumanReadable
            cpuRaw
            memoryRaw
          }
        }
        ... on CloseTrace {
          parentSpanID
          timestamp
          podName
          namespace
          itemType
          responseBody
          metrics {
            cpuHumanReadable
            memoryHumanReadable
            cpuRaw
            memoryRaw
          }
        }
        ... on Trace {
          parentSpanID
          spanID
          method
          url
          host
          timestamp
          podName
          namespace
          itemType
          metrics {
            cpuHumanReadable
            memoryHumanReadable
            cpuRaw
            memoryRaw
          }
        }
        ... on Span {
          parentSpanID
          spanID
          namespace
          timestamp
          podName
          itemType
          action
          status
          took
        }
        ... on DatabaseSpan {
          parentSpanID
          spanID
          itemType
          query
          namespace
          timestamp
          podName
          took
          hits
        }
      }
    }
  }
`;

const handleSearch = async (filters) => {
  try {
    const result = await client.query({
      query: GET_TRACES,
      variables: { input: filters }
    });

    queryResult.value = result.data;

    if (result.data.traces && result.data.traces.length > 0) {
      const jsonResult = JSON.stringify(result.data.traces[0], null, 2);
      ganttContent.value = convertTraceToGantt(jsonResult);
      stateChartContent.value = convertTraceToStateDiagram(jsonResult);
    }

    if (filters.podName) {
      const allItems = result.data.traces.flatMap(trace => trace.items);
      const filteredItems = allItems.filter(item =>
        item.itemType === 'trace' &&
        item.podName === filters.podName &&
        'metrics' in item &&
        item.metrics != null
      );

      const sortedItems = filteredItems.sort((a, b) => {
        const dateA = new Date(a.timestamp.replace('T', ' '));
        const dateB = new Date(b.timestamp.replace('T', ' '));
        return dateA - dateB;
      });

      if (sortedItems.length > 0) {
        createLineDiagram(sortedItems);
      }
    }
  } catch (error) {
    console.error('Error fetching traces:', error);
  }
};

function convertTraceToStateDiagram(traceData) {
  const traceDataParsed = JSON.parse(traceData);
  let diagram = `---\nState Diagram for ${traceDataParsed.traceID}\n---\n`;
  diagram += 'stateDiagram-v2\n';

  const traceIdToPodName = {};
  traceDataParsed.items.filter(item => item.itemType === 'trace').forEach(trace => {
    traceIdToPodName[trace.spanID] = trace.podName.split('-')[0];
  });

  let firstServicePodName = "";
  if (traceDataParsed.items.length > 0 && traceDataParsed.items[0].itemType === 'trace_start') {
    firstServicePodName = traceDataParsed.items[0].podName.split('-')[0];
    traceIdToPodName[traceDataParsed.items[0].spanID] = firstServicePodName;
    if (firstServicePodName) {
      diagram += `    [*] --> ${firstServicePodName}\n`;
    }
  }

  traceDataParsed.items.filter(item => item.itemType === 'trace').forEach(trace => {
    const parentService = traceIdToPodName[trace.parentSpanID];
    const childService = trace.podName.split('-')[0];

    if (parentService && childService && parentService !== childService) {
      diagram += `    ${parentService} --> ${childService}\n`;
    }
  });

  return diagram;
}

function convertTraceToGantt(traceData) {
  const traceDataParsed = JSON.parse(traceData);
  const traceStartItem = traceDataParsed.items.find(item => item.itemType === 'trace_start');

  if (!traceStartItem) {
    return '';
  }

  let diagram = 'gantt\n';
  diagram += `dateFormat YYYY-MM-DD'T'HH:mm:ss.SSS\n`;
  diagram += `axisFormat %Y-%m-%d | %H:%M:%S.%L\n`;
  diagram += `title Trace Gantt Chart for ${traceDataParsed.traceID}\n`;
  diagram += `section Parent\n`;
  diagram += `Overall Timeline :milestone, ${traceDataParsed.timeStarted}, ${traceDataParsed.timeEnded}\n`;

  const parsedEndTime = traceDataParsed.timeEnded.replace('T', "'T'");
  diagram += `section Trace\n`;
  diagram += `    ${traceStartItem.podName} :${1}, ${traceStartItem.timestamp}, ${parsedEndTime}\n`;

  const traceItems = traceDataParsed.items.filter(item => item.itemType === 'trace')
    .sort((a, b) => {
      const dateA = new Date(a.timestamp.replace("'T'", " "));
      const dateB = new Date(b.timestamp.replace("'T'", " "));
      return dateA - dateB;
    });

  traceItems.forEach((item, index) => {
    const taskEnd = index < traceItems.length - 1 ? traceItems[index + 1].timestamp : traceDataParsed.timeEnded;
    diagram += `    ${item.podName} :${index + 2}, ${item.timestamp}, ${taskEnd}\n`;
  });

  const spanItems = traceDataParsed.items.filter(item => item.itemType === 'span')
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

  const databaseItems = traceDataParsed.items.filter(item => item.itemType === 'database_span')
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
  const timestamps = filteredTraces.map((_, index) => index + 1);
  const cpuValues = filteredTraces.map(trace => trace.metrics.cpuRaw);
  const memoryValues = filteredTraces.map(trace => trace.metrics.memoryRaw);

  const startTime = new Date(filteredTraces[0].timestamp.replace("'T'", " "));
  const endTime = new Date(filteredTraces[filteredTraces.length - 1].timestamp.replace("'T'", " "));

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

  podMetricsMemory.value = memoryDiagram;
  podMetricsCpu.value = cpuDiagram;
}
</script>
