<template>
  <v-card class="mb-4">
    <v-card-title>Search Traces</v-card-title>
    <v-card-text>
      <v-row>
        <v-col cols="12" md="6">
          <v-autocomplete
            v-model="filters.ids"
            :items="ids"
            label="Trace IDs"
            multiple
            chips
            clearable
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-autocomplete
            v-model="filters.podName"
            :items="podNames"
            label="Pod Name"
            clearable
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" md="6">
          <v-autocomplete
            v-model="filters.operation"
            :items="operations"
            label="Operation"
            clearable
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-autocomplete
            v-model="filters.statusCode"
            :items="responseCodes"
            label="Status Code"
            clearable
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="filters.totalTimeHigherThan"
            label="Total Time Higher Than"
            type="text"
            clearable
          />
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="filters.beginTime"
            label="Begin Time"
            type="datetime-local"
            clearable
          />
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="filters.endTime"
            label="End Time"
            type="datetime-local"
            clearable
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-btn color="primary" @click="search" class="mr-2">Search</v-btn>
          <v-btn color="secondary" @click="clear">Clear</v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <v-card v-if="showQuery">
    <v-card-title>Query Parameters</v-card-title>
    <v-card-text>
      <JSONViewWrapper :jsonString="JSON.stringify(exampleQuery, null, 2)" />
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { gql } from '@apollo/client/core';
import client from '../graphql/client.js';
import JSONViewWrapper from './JSONViewWrapper.vue';

const emit = defineEmits(['search']);

const filters = reactive({
  ids: [],
  podName: '',
  operation: '',
  statusCode: '',
  totalTimeHigherThan: '',
  beginTime: '',
  endTime: '',
});

const ids = ref([]);
const responseCodes = ref([]);
const operations = ref([]);
const podNames = ref([]);
const showQuery = ref(true);

const exampleQuery = computed(() => {
  const query = { ...filters };
  Object.keys(query).forEach(key => {
    if (!query[key] || (Array.isArray(query[key]) && query[key].length === 0)) {
      query[key] = '';
    }
  });
  return query;
});

const formatDateTime = (dateTime) => {
  if (dateTime) {
    const date = new Date(dateTime);
    return date.toISOString().slice(0, 23);
  }
  return undefined;
};

const search = () => {
  const searchFilters = {};

  if (filters.ids && filters.ids.length > 0) {
    searchFilters.ids = filters.ids;
  }
  if (filters.podName) {
    searchFilters.podName = filters.podName;
  }
  if (filters.operation) {
    searchFilters.operation = filters.operation;
  }
  if (filters.statusCode) {
    searchFilters.statusCode = filters.statusCode;
  }
  if (filters.totalTimeHigherThan) {
    searchFilters.totalTimeHigherThan = filters.totalTimeHigherThan;
  }
  if (filters.beginTime) {
    searchFilters.beginTime = formatDateTime(filters.beginTime);
  }
  if (filters.endTime) {
    searchFilters.endTime = formatDateTime(filters.endTime);
  }

  emit('search', searchFilters);
};

const clear = () => {
  filters.ids = [];
  filters.podName = '';
  filters.operation = '';
  filters.statusCode = '';
  filters.totalTimeHigherThan = '';
  filters.beginTime = '';
  filters.endTime = '';
  search();
};

onMounted(async () => {
  const BUILDER_TRACES = gql`
    query BuildTraces($input: TraceQueryInput!) {
      traces(input: $input) {
        traceID
        responseCode
        items {
          ... on StartTrace {
            podName
            operation
          }
          ... on Trace {
            podName
          }
          ... on CloseTrace {
            podName
          }
          ... on Span {
            podName
          }
          ... on DatabaseSpan {
            podName
          }
        }
      }
    }
  `;

  try {
    const result = await client.query({
      query: BUILDER_TRACES,
      variables: { input: {} }
    });

    if (result.data && result.data.traces) {
      const traces = result.data.traces;
      const idsSet = new Set();
      const responseCodesSet = new Set();
      const operationsSet = new Set();
      const podNamesSet = new Set();

      traces.forEach(trace => {
        if (trace.traceID) {
          idsSet.add(trace.traceID);
        }
        if (trace.responseCode) {
          responseCodesSet.add(trace.responseCode);
        }
        trace.items.forEach(item => {
          if (item.operation) {
            operationsSet.add(item.operation);
          }
          if (item.podName) {
            podNamesSet.add(item.podName);
          }
        });
      });

      ids.value = Array.from(idsSet);
      responseCodes.value = Array.from(responseCodesSet);
      operations.value = Array.from(operationsSet);
      podNames.value = Array.from(podNamesSet);
    }
  } catch (error) {
    console.error('Error fetching builder traces:', error);
  }

  search();
});
</script>
