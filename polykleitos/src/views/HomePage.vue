<template>
  <div>
    <v-app id="homepage" class="paper-card-light">
      <v-container fluid>
      <!-- Hero Section with Parallax -->
      <v-parallax height="100em" :src="attike">
        <div class="d-flex flex-column justify-center align-center text-white" style="height: 100%;">
          <h1 class="text-h2 mb-4 text-center"
              style="text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);">
            Polykleitos
          </h1>

          <h2 class="text-h4 mb-4 text-center"
              style="text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);">
            Distributed Tracing & Observability
          </h2>

          <p class="text-h6 mb-8 text-center px-4"
             style="max-width: 820px; text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);">
            Named after <strong>Polykleitos</strong>,
            the classical sculptor who sought harmony through proportion,
            this interface presents the hidden structure of your system.
            <br><br>
            Polykleitos is a single-page observability view for exploring traces,
            spans, service hops, and performance metrics across the entire
            Odysseia-Greek platform.
          </p>
          <v-icon class="scroll-icon" size="48" @click="scrollToSearch">
            mdi-chevron-down
          </v-icon>
          <p class="mt-4 text-center">Scroll down to search traces</p>
        </div>
      </v-parallax>

      <!-- Main Content -->
      <v-container fluid class="paper-card-wrapper" id="search-section">
        <v-card color="triadic" class="mb-6" elevation="10" rounded="lg">
          <v-card-text>
            <div class="text-h4 mb-2">Search Traces</div>
            <div class="text-body-1">
              Paste a trace ID to fetch the full journey through your distributed system.
            </div>
          </v-card-text>
        </v-card>

        <!-- Search (no prefetch; only query on submit) -->
        <v-card class="paper-card-light" rounded="lg">
          <v-card-text>
            <v-row align="center">
              <v-col cols="12" md="9">
                <v-autocomplete
                    v-model="search"
                    :items="searchHistory"
                    :loading="loading"
                    label="Trace ID"
                    placeholder="424ea68f-28a6-412a-bd47-2d05bf6fdf8d"
                    prepend-inner-icon="mdi-magnify"
                    clearable
                    hide-no-data
                    @keyup.enter="runSearch(search)"
                    @update:modelValue="onSelectFromHistory"
                />
              </v-col>

              <v-col cols="12" md="3">
                <v-btn color="primary" block :loading="loading" @click="runSearch(search)">
                  Search
                </v-btn>
              </v-col>
            </v-row>

            <v-alert v-if="errorText" type="warning" class="mt-3" rounded="lg">
              {{ errorText }}
            </v-alert>
          </v-card-text>
        </v-card>

        <!-- ES searched Traces -->
        <v-card class="paper-card-light mt-6" rounded="lg" elevation="6">
          <v-card-title class="d-flex align-center">
            Trace Search (Elasticsearch)
            <v-spacer />
            <v-btn size="small" variant="text" icon="mdi-refresh" @click="runTraceSearch" />
          </v-card-title>

          <v-card-subtitle class="d-flex align-center">
    <span class="text-caption">
      Total matched: <strong>{{ searchTotal }}</strong>
    </span>
          </v-card-subtitle>

          <v-card-text>
            <v-card class="mb-4" variant="tonal" rounded="lg">
              <v-card-text>
                <v-row dense>
                  <v-col cols="12" md="3">
                    <v-select
                        v-model="searchWindow"
                        :items="windowOptions"
                        item-title="title"
                        item-value="value"
                        label="Time window"
                        density="compact"
                    />
                  </v-col>

                  <v-col cols="12" md="3">
                    <v-text-field
                        v-model="searchOperation"
                        label="Operation (optional)"
                        placeholder="mediaAnswer"
                        density="compact"
                        clearable
                    />
                  </v-col>

                  <v-col cols="12" md="2">
                    <v-text-field
                        v-model.number="searchResponseCode"
                        label="Response code"
                        placeholder="200"
                        density="compact"
                        clearable
                        type="number"
                    />
                  </v-col>

                  <v-col cols="12" md="2">
                    <v-text-field
                        v-model.number="searchTimeTookGreaterThan"
                        label="Min total time (ms)"
                        placeholder="25"
                        density="compact"
                        clearable
                        type="number"
                    />
                  </v-col>

                  <v-col cols="12" md="1">
                    <v-text-field
                        v-model.number="searchLimit"
                        label="Limit"
                        density="compact"
                        type="number"
                    />
                  </v-col>

                  <v-col cols="12" md="1" class="d-flex align-center">
                    <v-btn color="primary" block :loading="searchLoading" @click="runTraceSearch">
                      Run
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <v-alert v-if="searchError" type="error" rounded="lg" class="mb-3">
              {{ searchError.message }}
            </v-alert>

            <v-data-table
                :headers="traceSearchHeaders"
                :items="searchRows"
                :loading="searchLoading"
                item-key="id"
                density="compact"
                class="elevation-1"
                :items-per-page="10"
                hover
                @click:row="onTraceRowClick"
            >
              <template #item.id="{ item }">
                <span class="mono">{{ item.id }}</span>
              </template>

              <template #item.isActive="{ item }">
                <v-chip size="small" :color="item.isActive ? 'orange' : 'green'">
                  {{ item.isActive ? "active" : "closed" }}
                </v-chip>
              </template>

              <template #item.hasDbSpan="{ item }">
                <v-chip size="small" :color="item.hasDbSpan ? 'green' : 'grey'">
                  DB
                </v-chip>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>

        <!-- Live Traces -->
        <v-card class="paper-card-light mb-6" rounded="lg" elevation="6">
          <v-card-title class="d-flex align-center">
            Live Traces
            <v-spacer />
            <v-btn size="small" variant="text" icon="mdi-broom" @click="clearLive" />
          </v-card-title>

          <v-card-subtitle class="d-flex align-center">
            <span class="text-caption">updated: {{ liveUpdatedAt || "—" }}</span>
            <v-spacer />
            <span class="text-caption">polls every 10s (max {{ LIVE_POLL_LIMIT }}/poll)</span>
          </v-card-subtitle>

          <v-card-text>
            <v-data-table
                density="compact"
                :headers="liveHeaders"
                :items="liveTraces"
                :items-per-page="10"
                item-key="id"
                class="elevation-1"
                hover
                @click:row="onLiveRowClick"
            >
              <template #item.id="{ item }">
                <span class="mono">{{ item.id }}</span>
              </template>

              <template #item.isActive="{ item }">
                <v-chip size="small" :color="item.isActive ? 'red' : 'green'">
                  {{ item.isActive ? "active" : "done" }}
                </v-chip>
              </template>

              <template #item.hasDbSpan="{ item }">
                <v-chip size="small" :color="item.hasDbSpan ? 'orange' : 'grey'">
                  {{ item.hasDbSpan ? "db" : "-" }}
                </v-chip>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>


    <!-- Errors from backend -->
    <v-alert v-if="queryError" type="error" class="mb-6" rounded="lg">
      {{ queryError.message }}
    </v-alert>

    <!-- No result -->
    <v-alert v-else-if="searched && !loading && !trace" type="info" class="mb-6" rounded="lg">
      No trace found for ID: <strong>{{ lastSearchedId }}</strong>
    </v-alert>

    <!-- Trace -->
    <v-card class="paper-card-report" v-if="trace" rounded="lg" elevation="6"  id="summary-section">
      <v-card-title>Trace Summary</v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="12" md="4"><strong>Operation:</strong> {{ trace.operation }}</v-col>
          <v-col cols="12" md="4"><strong>Namespace:</strong> {{ trace.namespace }}</v-col>
          <v-col cols="12" md="4"><strong>Pod:</strong> {{ trace.podName }}</v-col>
        </v-row>

        <v-row class="mt-2">
          <v-col cols="12" md="4"><strong>Total time:</strong> {{ trace.totalTimeMs }} ms</v-col>
          <v-col cols="12" md="8">
            <v-chip :color="trace.hasDbSpan ? 'green' : 'grey'" class="mr-2">DB</v-chip>
            <v-chip :color="trace.hasAction ? 'blue' : 'grey'">Action</v-chip>
            <span class="ml-4 text-caption"><strong>Trace:</strong> {{ trace.id }}</span>
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <v-btn-toggle v-model="viewMode" mandatory class="mb-4">
          <v-btn value="raw">Raw JSON</v-btn>
          <v-btn value="state">Mermaid State</v-btn>
          <v-btn value="gantt">Mermaid Gantt</v-btn>
        </v-btn-toggle>

        <div v-if="viewMode === 'raw'">
          <JsonViewer :data="trace" title="trace" :default-collapsed="false" :deep="30" />
        </div>

        <div v-if="viewMode === 'state'">
          <TraceMermaid :trace="trace" />
        </div>

        <div v-if="viewMode === 'gantt'">
          <TraceGantt v-if="trace" :trace="trace" />
        </div>


        </v-card-text>
        </v-card>

        <MetricsSummaryPanel />
      </v-container>
    </v-container>
    </v-app>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useLazyQuery, useQuery } from "@vue/apollo-composable";
import { TraceByIdQuery } from "../constants/traceById.js";
import { TracePollQuery } from "../constants/tracePoll.js";
import { TraceSearchQuery } from "../constants/traceSearch.js";
import JsonViewer from "@/components/JSONViewWrapper.vue";
import TraceMermaid from "@/components/TraceMermaidState.vue";
import TraceGantt from "@/components/TraceMermaidGantt.vue";
import MetricsSummaryPanel from "@/components/MetricsSummaryPanel.vue";

const search = ref("");               // current input / selection
const searchHistory = ref([]);        // string[] of trace IDs
const searched = ref(false);
const lastSearchedId = ref("");
const viewMode = ref("raw");

const attike = ref('');

const LIVE_KEEP = 100;
const LIVE_POLL_LIMIT = 5;

// live state
const liveTraces = ref([]);
const liveUpdatedAt = ref("");

const searchRows = ref([]);         // TraceSummary rows
const searchTotal = ref(0);
const searchLoading = ref(false);
const searchError = ref(null);

// UI controls (defaults you can tweak)
const searchWindow = ref("M30");          // GraphQL enum string
const searchOperation = ref("");
const searchResponseCode = ref(null);
const searchTimeTookGreaterThan = ref(null);
const searchLimit = ref(20);

// headers for v-data-table
const liveHeaders = [
  { title: "Trace ID", key: "id" },
  { title: "RootQuery", key: "rootQuery" },
  { title: "Total (ms)", key: "totalTimeMs" },
  { title: "Code", key: "responseCode" },
  { title: "DB", key: "hasDbSpan" },
  { title: "Status", key: "isActive" },
  { title: "Started", key: "timeStarted" },
  { title: "Ended", key: "timeEnded", sortable: true },
];

const traceSearchHeaders = [
  { title: "Trace ID", key: "id" },
  { title: "RootQuery", key: "rootQuery" },
  { title: "Total (ms)", key: "totalTimeMs" },
  { title: "Code", key: "responseCode" },
  { title: "DB", key: "hasDbSpan" },
  { title: "Status", key: "isActive" },
  { title: "Started", key: "timeStarted" },
];

const windowOptions = [
  { title: "Last 5 minutes", value: "M5" },
  { title: "Last 10 minutes", value: "M10" },
  { title: "Last 30 minutes", value: "M30" },
  { title: "Last 1 hour", value: "H1" },
  { title: "Last 2 hours", value: "H2" },
  { title: "Last 12 hours", value: "H12" },
  { title: "Last 24 hours", value: "H24" },
];

const onTraceRowClick = async (_event, row) => {
  const id = row?.item?.id;
  if (!id) return;

  // add it to the autocomplete history + prefill
  search.value = id;
  addToHistory(id);

  // optionally auto-load the full trace immediately:
  await runSearch(id);
};

const {
  result: pollResult,
  loading: pollLoading,
  error: pollError,
} = useQuery(
    TracePollQuery,
    () => ({ limit: LIVE_POLL_LIMIT }),
    {
      fetchPolicy: "network-only",
      pollInterval: 10_000,
      notifyOnNetworkStatusChange: true,
    }
);

const loadImage = () => {
  import('@/assets/akropolis_klenze.webp').then((module) => {
    attike.value = module.default;
  });
};

function mergeLive(incoming) {
  if (!incoming || incoming.length === 0) return;

  for (const t of incoming) {
    addToHistory(t.id);
  }

  const map = new Map();
  // put incoming first so they win
  for (const t of incoming) map.set(t.id, t);
  for (const t of liveTraces.value) {
    if (!map.has(t.id)) map.set(t.id, t);
  }

  // Keep order: incoming first, then existing
  const merged = Array.from(map.values());

  // Optional: sort active first, then by totalTimeMs desc
  merged.sort((a, b) => {
    if (a.isActive !== b.isActive) return a.isActive ? -1 : 1;
    return (b.totalTimeMs ?? 0) - (a.totalTimeMs ?? 0);
  });

  liveTraces.value = merged.slice(0, LIVE_KEEP);
}

watch(
    () => pollResult.value,
    (val) => {
      const payload = val?.tracePoll;
      if (!payload) return;
      liveUpdatedAt.value = payload.updatedAt || "";
      mergeLive(payload.traces || []);
    }
);

function clearLive() {
  liveTraces.value = [];
}

// Click row behavior: copy + optionally auto-search
async function onLiveRowClick(_event, row) {
  const id = row?.item?.id ?? row?.id; // depending on vuetify event shape
  if (!id) return;

  // Put in input
  search.value = id;
  addToHistory(id);

  // Copy to clipboard (best-effort)
  try {
    await navigator.clipboard.writeText(id);
  } catch (e) {
    // ignore (clipboard may require https)
  }

  // Optional: auto-run search (uncomment if desired)
  // await runSearch(id);
}

const errorText = computed(() => {
  const v = (search.value || "").trim();
  if (!v && searched.value) return "Please enter a trace ID.";
  if (v && v.length < 8) return "That trace ID looks too short.";
  return "";
});

// no prefetch: lazy query triggers only when you call load()
const { load, result, loading, error: queryError } = useLazyQuery(
    TraceByIdQuery,
    () => ({ id: lastSearchedId.value }),
    { fetchPolicy: "network-only" }
);

const trace = computed(() => result.value?.trace ?? null);

const sortedItems = computed(() => {
  const items = trace.value?.items ?? [];
  return [...items].sort((a, b) => String(a.timestamp).localeCompare(String(b.timestamp)));
});

const formatPayload = (payload) => (payload ? JSON.stringify(payload, null, 2) : "");

const addToHistory = (id) => {
  const v = (id || "").trim();
  if (!v) return;
  // de-dupe, most-recent-first
  searchHistory.value = [v, ...searchHistory.value.filter((x) => x !== v)].slice(0, 20);
};

const runSearch = async (id) => {
  const v = (id || "").trim();
  searched.value = true;

  if (!v || v.length < 8) return;

  lastSearchedId.value = v;
  addToHistory(v);

  scrollToSummary()

  // triggers the query (first time) and refetches (subsequent times)
  await load();
};

// If user picked from dropdown, you can auto-search immediately (optional)
const onSelectFromHistory = (val) => {
  // comment out if you *only* want to search on Enter/click
  // runSearch(val);
};


const { load: loadTraceSearch, result: traceSearchResult, loading: traceSearchLoading, error: traceSearchError, refetch: refetchTraceSearch } =
    useLazyQuery(
        TraceSearchQuery,
        () => ({ input: buildTraceSearchInput() }),
        { fetchPolicy: "network-only" }
    );

watch(
    () => traceSearchResult.value,
    (v) => {
      const page = v?.traceSearch;
      searchRows.value = page?.items ?? [];
      searchTotal.value = page?.total ?? 0;
    }
);

watch(traceSearchError, (e) => (searchError.value = e));
watch(traceSearchLoading, (v) => (searchLoading.value = v));

const runTraceSearch = async () => {
  // first call uses load(); subsequent calls can use refetch()
  if (!traceSearchResult.value) {
    await loadTraceSearch();
  } else {
    await refetchTraceSearch();
  }
};

const buildTraceSearchInput = () => {
  const input = {
    limit: searchLimit.value,
    window: searchWindow.value,
  };

  const op = (searchOperation.value || "").trim();
  if (op) input.operation = op;

  if (typeof searchResponseCode.value === "number" && !Number.isNaN(searchResponseCode.value)) {
    input.responseCode = searchResponseCode.value;
  }

  if (typeof searchTimeTookGreaterThan.value === "number" && !Number.isNaN(searchTimeTookGreaterThan.value)) {
    input.timeTookGreaterThan = searchTimeTookGreaterThan.value;
  }

  return input;
};

const scrollToSearch = () => {
  const searchSection = document.getElementById('search-section');
  if (searchSection) {
    searchSection.scrollIntoView({ behavior: 'smooth' });
  }
};

const scrollToSummary = () => {
  const summarySection = document.getElementById('summary-section');
  if (summarySection) {
    summarySection.scrollIntoView({ behavior: 'smooth' });
  }
};

onMounted(() => {
  loadImage();
  runTraceSearch();
});
</script>

<style scoped>
.paper-card-wrapper {
  max-width: 80%;
  background: #fdf6e3; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: 'Roboto', serif;
}

.paper-card-report {
  background: #fdf6e3; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: 'Roboto', serif;
}

.paper-card-light {
  background: #fefcf5; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: 'Roboto', serif;
}

.scroll-icon {
  cursor: pointer;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
  60% {
    transform: translateY(-5px);
  }
}

.scroll-icon:hover {
  color: #1c61d1;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

</style>