<template>
  <v-card rounded="lg" elevation="6">
    <v-card-title class="d-flex align-center">
      Trace Flow (stateDiagram)
      <v-spacer />
      <v-btn
          size="small"
          variant="text"
          icon="mdi-refresh"
          @click="render()"
          :disabled="!filteredItems.length"
      />
      <v-btn
          size="small"
          variant="text"
          icon="mdi-content-copy"
          @click="copyMermaid"
          :disabled="!mermaidCode"
      />
    </v-card-title>

    <v-card-text>
      <v-alert v-if="!filteredItems.length" type="info" rounded="lg">
        No relevant trace items to render.
      </v-alert>

      <div v-else class="mermaid-shell">
        <div ref="host" class="mermaid-host"></div>
      </div>

      <v-expansion-panels class="mt-4" variant="accordion">
        <v-expansion-panel>
          <v-expansion-panel-title color="primary">Mermaid source</v-expansion-panel-title>
          <v-expansion-panel-text>
            <pre class="code">{{ mermaidCode }}</pre>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import mermaid from "mermaid";

const props = defineProps({
  trace: { type: Object, default: null },
  items: { type: Array, default: null },
});

const host = ref(null);

const rawItems = computed(() => props.items ?? props.trace?.items ?? []);

const filteredItems = computed(() => {
  const allow = new Set(["TRACE_START", "TRACE_HOP", "GRAPHQL", "TRACE_STOP"]);
  return (rawItems.value || []).filter((it) =>
      allow.has(String(it.itemType || "").toUpperCase())
  );
});

const mermaidCode = computed(() => buildStateDiagram(filteredItems.value));

let mermaidInitialized = false;

onMounted(() => {
    if (!mermaidInitialized) {
      mermaid.initialize({
        startOnLoad: false,
        look: "handDrawn",
        theme: "loose",
        securityLevel: "loose",
        themeVariables: {
          fontFamily: "Roboto, Arial, sans-serif",
          fontSize: "18px",
        },
      });
      mermaidInitialized = true;
    }
  render();
});

watch(mermaidCode, () => render());

async function render() {
  if (!host.value) return;

  if (!mermaidCode.value || mermaidCode.value.trim().length === 0) {
    host.value.innerHTML = "";
    return;
  }

  try {
    const id = "mmd-" + Math.random().toString(36).slice(2);
    const { svg } = await mermaid.render(id, mermaidCode.value);
    host.value.innerHTML = svg;
  } catch (e) {
    host.value.innerHTML = `<pre style="color:#ffb4b4; white-space:pre-wrap;">${String(e)}</pre>`;
  }
}

async function copyMermaid() {
  try {
    await navigator.clipboard.writeText(mermaidCode.value || "");
  } catch {
    // ignore
  }
}

// ---- core logic ----

function buildStateDiagram(items) {
  if (!items || items.length === 0) return "";

  // Sort by time so "root" detection is stable
  const sorted = [...items].sort((a, b) =>
      String(a.timestamp).localeCompare(String(b.timestamp))
  );

  // We build spanId -> representative serviceName (from podName)
  const spanService = new Map(); // spanId -> "homeros"
  const spanParent = new Map();  // spanId -> parentSpanId
  const spanFirstType = new Map(); // spanId -> first itemType seen (for root hint)

  for (const it of sorted) {
    const sid = it.spanId || "";
    if (!sid) continue;

    if (!spanService.has(sid)) {
      spanService.set(sid, serviceFromPod(it.podName));
      spanParent.set(sid, it.parentSpanId || "");
      spanFirstType.set(sid, String(it.itemType || "").toUpperCase());
    }
  }

  // Root: prefer the span whose first event is TRACE_START
  let rootSpan = "";
  for (const [sid, t] of spanFirstType.entries()) {
    if (t === "TRACE_START") {
      rootSpan = sid;
      break;
    }
  }
  // fallback: span with no parent
  if (!rootSpan) {
    for (const [sid, p] of spanParent.entries()) {
      if (!p) {
        rootSpan = sid;
        break;
      }
    }
  }

  const rootService = rootSpan ? spanService.get(rootSpan) : "root";

  // Build edges parentService -> childService based on span relationships
  const edges = new Set();
  for (const [sid, parentSid] of spanParent.entries()) {
    const childSvc = spanService.get(sid);
    if (!childSvc) continue;

    // If no parent, attach to rootService (as you described)
    if (!parentSid || !spanService.has(parentSid)) {
      if (childSvc && rootService && childSvc !== rootService) {
        edges.add(`${rootService} --> ${childSvc}`);
      }
      continue;
    }

    const parentSvc = spanService.get(parentSid);
    if (!parentSvc) continue;

    // Avoid self loops (often happens when multiple spans map to same service)
    if (parentSvc === childSvc) continue;

    edges.add(`${parentSvc} --> ${childSvc}`);
  }

  // Ensure root is present even if there are no edges yet
  const lines = [];
  lines.push("stateDiagram-v2");
  lines.push(`  %% Root service inferred from TRACE_START`);
  lines.push(`  state "${rootService}" as ${safeId(rootService)}`);

  // Also declare states so Mermaid labels are nice
  const allServices = new Set([...spanService.values()]);
  for (const svc of allServices) {
    if (svc === rootService) continue;
    lines.push(`  state "${svc}" as ${safeId(svc)}`);
  }

  // If we have no edges, at least show root
  if (edges.size === 0) return lines.join("\n");

  // Render edges (use aliases)
  for (const e of [...edges]) {
    const [a, b] = e.split(" --> ").map((s) => s.trim());
    lines.push(`  ${safeId(a)} --> ${safeId(b)}`);
  }

  return lines.join("\n");
}

function serviceFromPod(podName) {
  const p = String(podName || "").trim();
  if (!p) return "unknown";

  // Common k8s patterns:
  // homeros-56d8fc58d9-lv97g -> homeros
  // aspasia-6c8c7f6df-schtx -> aspasia
  let base = p.split("-")[0];
  return base || p;
}

function safeId(name) {
  // Mermaid state IDs must be simple tokens
  return String(name)
      .toLowerCase()
      .replace(/[^a-z0-9_]/g, "_");
}
</script>

<style scoped>
.mermaid-shell {
  overflow: auto;
}

.mermaid-host :deep(svg) {
  width: 25% !important;
  height: auto !important;
  max-width: none !important;
  display: block;
}

/* Most important for gantt: do NOT force svg to fit */
.mermaid :deep(svg) {
  max-width: none !important;
  height: auto;
}
</style>