<template>
  <v-card class="paper-card-light mt-6" rounded="lg" elevation="6" id="metrics-section">
    <v-card-title class="d-flex align-center">
      Metrics Summary
      <v-spacer />
      <v-btn size="small" variant="text" icon="mdi-refresh" @click="runMetrics" />
    </v-card-title>

    <v-card-subtitle class="d-flex align-center">
      <span class="text-caption">
        Window: <strong>{{ metricsWindow }}</strong>
      </span>
      <v-spacer />
      <span class="text-caption" v-if="summary">
        {{ summary.start }} → {{ summary.end }}
      </span>
    </v-card-subtitle>

    <v-card-text>
      <v-card class="mb-4" variant="tonal" rounded="lg">
        <v-card-text>
          <v-row dense>
            <v-col cols="12" md="3">
              <v-select
                  v-model="metricsWindow"
                  :items="windowOptions"
                  item-title="title"
                  item-value="value"
                  label="Time window"
                  density="compact"
              />
            </v-col>

            <v-col cols="12" md="2" class="d-flex align-center">
              <v-btn color="primary" block :loading="loading" @click="runMetrics">
                Run
              </v-btn>
            </v-col>

            <v-col cols="12" md="7" class="d-flex align-center">
              <span class="text-caption" v-if="summary">
                Nodes: <strong>{{ summary.nodes?.total ?? 0 }}</strong> •
                Namespaces: <strong>{{ summary.namespaces?.total ?? 0 }}</strong> •
                Pods: <strong>{{ summary.pods?.total ?? 0 }}</strong>
              </span>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <v-alert v-if="error" type="error" rounded="lg" class="mb-3">
        {{ error.message }}
      </v-alert>

      <v-row dense>
        <v-col cols="12" md="4">
          <div class="text-subtitle-1 mb-2">Nodes</div>
          <v-data-table
              :headers="nodeHeaders"
              :items="summary?.nodes?.items ?? []"
              :loading="loading"
              density="compact"
              class="elevation-1"
              :items-per-page="10"
              hover
          >
            <template #item.mem="{ item }">
              <span class="mono">{{ item.mem?.totalMaxHuman ?? item.mem?.totalMax ?? "—" }}</span>
            </template>
            <template #item.memMax="{ item }">
              <span class="mono">{{ item.mem?.max  }}</span>
            </template>
            <template #item.cpu="{ item }">
              <span class="mono">{{ item.cpu?.totalMaxHuman ?? item.cpu?.totalMax ?? "—" }}</span>
            </template>
          </v-data-table>
        </v-col>

        <v-col cols="12" md="4">
          <div class="text-subtitle-1 mb-2">Namespaces</div>
          <v-data-table
              :headers="nsHeaders"
              :items="summary?.namespaces?.items ?? []"
              :loading="loading"
              density="compact"
              class="elevation-1"
              :items-per-page="10"
              hover
          >
            <template #item.mem="{ item }">
              <span class="mono">{{ item.mem?.totalMaxHuman ?? item.mem?.totalMax ?? "—" }}</span>
            </template>
            <template #item.cpu="{ item }">
              <span class="mono">{{ item.cpu?.totalMaxHuman ?? item.cpu?.totalMax ?? "—" }}</span>
            </template>
          </v-data-table>
        </v-col>

        <v-col cols="12" md="4">
          <div class="text-subtitle-1 mb-2">Pods</div>
          <v-data-table
              :headers="podHeaders"
              :items="summary?.pods?.items ?? []"
              :loading="loading"
              density="compact"
              class="elevation-1"
              :items-per-page="10"
              hover
          >
            <template #item.podName="{ item }">
              <span class="mono">{{ item.podName }}</span>
            </template>
            <template #item.mem="{ item }">
              <span class="mono">{{ item.mem?.maxHuman ?? item.mem?.max ?? "—" }}</span>
            </template>
            <template #item.cpu="{ item }">
              <span class="mono">{{ item.cpu?.maxHuman ?? item.cpu?.max ?? "—" }}</span>
            </template>
          </v-data-table>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { useLazyQuery } from "@vue/apollo-composable";
import { MetricsSummaryQuery } from "@/constants/metricsSummary.js";

const metricsWindow = ref("M30");

const windowOptions = [
  { title: "Last 10 minutes", value: "M10" },
  { title: "Last 30 minutes", value: "M30" },
  { title: "Last 1 hour", value: "H1" },
  { title: "Last 2 hours", value: "H2" },
  { title: "Last 12 hours", value: "H12" },
  { title: "Last 24 hours", value: "H24" },
];

const nodeHeaders = [
  { title: "Node", key: "node" },
  { title: "Mem max", key: "mem" },
  { title: "CPU max", key: "cpu" },
];

const nsHeaders = [
  { title: "Namespace", key: "namespace" },
  { title: "Mem max", key: "mem" },
  { title: "CPU max", key: "cpu" },
];

const podHeaders = [
  { title: "Pod", key: "podName" },
  { title: "Mem max", key: "mem" },
  { title: "CPU max", key: "cpu" },
];

const buildInput = () => ({
  window: metricsWindow.value,
});

const { load, result, loading, error, refetch } = useLazyQuery(
    MetricsSummaryQuery,
    () => ({ input: buildInput() }),
    { fetchPolicy: "network-only" }
);

const summary = computed(() => result.value?.metricsSummary ?? null);

const runMetrics = async () => {
  if (!result.value) {
    await load();
  } else {
    await refetch();
  }
};

// Auto refresh when window changes (optional)
watch(metricsWindow, async () => {
  await runMetrics();
});

// expose for parent if you want
defineExpose({ runMetrics });
</script>

<style scoped>
.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}
</style>