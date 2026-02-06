<template>
  <v-card rounded="lg" elevation="6">
    <v-card-title class="d-flex align-center">
      Hop Timeline (gantt)
      <v-spacer />
      <v-btn
          size="small"
          variant="text"
          icon="mdi-refresh"
          @click="render()"
          :disabled="!blocks.length"
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
      <v-alert v-if="!blocks.length" type="info" rounded="lg">
        No spans found to render.
      </v-alert>

      <div class="mermaid-shell">
        <div ref="host" class="mermaid"></div>
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
  showTotal: { type: Boolean, default: true },
});

const host = ref(null);

const rawItems = computed(() => props.items ?? props.trace?.items ?? []);
const traceStartIso = computed(() => props.trace?.timeStarted ?? null);
const traceTotalMs = computed(() => props.trace?.totalTimeMs ?? null);

const blocks = computed(() => buildSpanBlocks(rawItems.value));

const mermaidCode = computed(() =>
    buildGanttTimeline({
      blocks: blocks.value,
      traceStartIso: traceStartIso.value,
      traceTotalMs: traceTotalMs.value,
      rootPodName: props.trace?.podName ?? null,
      operation: props.trace?.operation ?? "",
      showZeroVisual: true,
    })
);

let mermaidInitialized = false;

onMounted(() => {
  if (!mermaidInitialized) {
    mermaid.initialize({
      startOnLoad: false,
      theme: "neutral",
      securityLevel: "loose",
      gantt: {
        titleTopMargin: 50,
        barHeight: 30,
        barGap: 20,
        topPadding: 75,
        rightPadding: 75,
        leftPadding: 75,
        fontSize: 15,
        sectionFontSize: 25,
        gridLineStartPadding: 20,
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
    host.value.innerHTML = `<pre style="color:#ffb4b4; white-space:pre-wrap;">${String(
        e
    )}</pre>`;
  }
}

async function copyMermaid() {
  try {
    await navigator.clipboard.writeText(mermaidCode.value || "");
  } catch {
    // ignore
  }
}

// ---------- core logic (UPDATED) ----------

function buildSpanBlocks(items) {
  if (!items || items.length === 0) return [];

  const sorted = [...items].sort((a, b) =>
      String(a.timestamp).localeCompare(String(b.timestamp))
  );

  const isStartType = (t) => {
    const x = String(t || "").toUpperCase();
    return x === "TRACE_START" || x === "GRAPHQL" || x === "TRACE_HOP";
  };


  const spans = new Map(); // spanId -> { start, stop }
  const hopStops = []; // all hop stops for fallback matching

  for (const it of sorted) {
    const sid = it.spanId;
    if (!sid) continue;

    const t = String(it.itemType || "").toUpperCase();
    if (t === "TRACE_HOP_STOP") {
      hopStops.push(it);
    }

    if (!spans.has(sid)) spans.set(sid, { start: null, stop: null });

    if (isStartType(t)) {
      // keep earliest start per span
      if (!spans.get(sid).start) spans.get(sid).start = it;
    }

    if (t === "TRACE_HOP_STOP") {
      // keep latest stop per span (usually correct end)
      spans.get(sid).stop = it;
    }
  }

  const out = [];

  for (const [sid, v] of spans.entries()) {
    if (!v.start) continue;

    const startMs = parseIsoMs(v.start.timestamp);
    let stop = v.stop;

    // Fallback: if no hop_stop with same spanId exists,
    // match by same ns+pod and pick the LATEST hop_stop after start.
    if (!stop) {
      stop =
          hopStops
              .filter((hs) => {
                if (!hs) return false;
                if (hs.namespace !== v.start.namespace) return false;
                if (hs.podName !== v.start.podName) return false;

                const hsMs = parseIsoMs(hs.timestamp);
                if (startMs == null || hsMs == null) return false;
                return hsMs >= startMs;
              })
              .sort((a, b) => parseIsoMs(a.timestamp) - parseIsoMs(b.timestamp))
              .pop() || null;
    }

    const stopMs = stop ? parseIsoMs(stop.timestamp) : null;

    // prefer tookMs if present; otherwise wall-clock diff
    const tookMs = stop?.payload?.tookMs;
    let durMs = isFiniteNumber(tookMs) ? Number(tookMs) : null;

    if (durMs == null && startMs != null && stopMs != null && stopMs >= startMs) {
      durMs = stopMs - startMs;
    }

    // Label method: prefer hop method; else GraphQL op; else url
    const method =
        v.start.payload?.method ||
        v.start.payload?.operation ||
        v.start.payload?.url ||
        "span";

    out.push({
      spanId: sid,
      namespace: v.start.namespace || "unknown",
      podName: v.start.podName || "unknown",
      method,
      startMs,
      durMs,
      startType: String(v.start.itemType || "").toUpperCase(),
    });
  }

  // Sort by start time
  out.sort((a, b) => (a.startMs ?? 1e18) - (b.startMs ?? 1e18));

  return out;
}

function buildGanttTimeline({
                              blocks,
                              traceStartIso,
                              traceTotalMs,
                              rootPodName,
                              operation,
                              showZeroVisual = true, // set false if you want literal 0ms bars
                            }) {
  if (!blocks || blocks.length === 0) return "";

  const baseStart =
      traceStartIso || firstNonNullIso(blocks) || new Date().toISOString();

  const op = sanitizeText(operation || "trace");
  const rootSvc = serviceFromPod(rootPodName);

  // Sort by start time (nulls last)
  const sorted = [...blocks].sort((a, b) => (a.startMs ?? 1e18) - (b.startMs ?? 1e18));

  const lines = [];
  lines.push("gantt");
  lines.push(`  title Trace hop timings ${op}`);
  lines.push("  dateFormat  YYYY-MM-DDTHH:mm:ss.SSSZ");
  lines.push("  axisFormat  %H:%M:%S.%L");
  lines.push("");

  lines.push("  section timeline");

  // Root bar first (homeros total)
  if (isFiniteNumber(traceTotalMs) && traceTotalMs > 0 && rootSvc) {
    const total = Math.round(traceTotalMs);
    const rootLabel = sanitizeText(`pod ${rootSvc} — ${total}ms (root)`);
    lines.push(`  ${rootLabel} :active, root, ${baseStart}, ${total}ms`);
  }

  for (const b of sorted) {
    const service = serviceFromPod(b.podName);

    // skip root service blocks so we don't duplicate the root bar
    if (rootSvc && service === rootSvc) continue;

    const startIso = b.startMs != null ? new Date(b.startMs).toISOString() : baseStart;

    const realMs = isFiniteNumber(b.durMs) ? Math.max(0, Math.round(b.durMs)) : 0;

    // Optional: make 0ms visible but keep label truthful
    let visualMs = realMs;
    if (showZeroVisual && realMs === 0) {
      visualMs = isFiniteNumber(traceTotalMs)
          ? Math.max(1, Math.floor(traceTotalMs / 10))
          : 1;
    }

    const label = sanitizeText(`pod ${service} — ${realMs}ms`);
    const id = `h_${b.spanId}`;

    lines.push(`  ${label} :active, ${id}, ${startIso}, ${visualMs}ms`);
  }

  lines.push("");
  return lines.join("\n");
}

function buildGantt({
                      blocks,
                      traceStartIso,
                      traceTotalMs,
                      rootPodName,
                      showTotal,
                      operation,
                    }) {
  if (!blocks || blocks.length === 0) return "";

  const baseStart =
      traceStartIso || firstNonNullIso(blocks) || new Date().toISOString();

  const op = sanitizeText(operation || "trace");
  const rootSvc = serviceFromPod(rootPodName);

  const lines = [];
  lines.push("gantt");
  lines.push(`  title Trace hop timings ${op}`);
  lines.push("  dateFormat  YYYY-MM-DDTHH:mm:ss.SSSZ");
  lines.push("  axisFormat  %H:%M:%S.%L");
  lines.push("");

  // Group blocks by namespace
  const byNs = new Map();
  for (const b of blocks) {
    const ns = b.namespace || "unknown";
    if (!byNs.has(ns)) byNs.set(ns, []);
    byNs.get(ns).push(b);
  }

  for (const [ns, arr] of byNs.entries()) {
    lines.push(`  section ns: ${sanitizeText(ns)}`);

    // 👉 Root service (homeros) goes here
    if (
        showTotal &&
        rootSvc &&
        ns === (blocks.find(b => serviceFromPod(b.podName) === rootSvc)?.namespace)
    ) {
      const totalMs = isFiniteNumber(traceTotalMs)
          ? Math.round(traceTotalMs)
          : 0;

      const rootLabel = sanitizeText(`pod ${rootSvc} — ${totalMs}ms`);
      lines.push(
          `  ${rootLabel} :active, root, ${baseStart}, ${totalMs}ms`
      );
    }

    for (const b of arr) {
      const service = serviceFromPod(b.podName);

      // Skip root hop (we already rendered it as TOTAL)
      if (service === rootSvc) continue;

      const startIso =
          b.startMs != null ? new Date(b.startMs).toISOString() : baseStart;

      const realMs = isFiniteNumber(b.durMs)
          ? Math.max(0, Math.round(b.durMs))
          : 0;

// Visual duration: make 0ms hops visible but tiny
      const visualMs =
          realMs > 0
              ? realMs
              : isFiniteNumber(traceTotalMs)
                  ? Math.max(1, Math.floor(traceTotalMs / 10))
                  : 1;

      const label = sanitizeText(`pod ${service} — ${realMs}ms`);
      const dur = `${visualMs}ms`;
      const id = `h_${b.spanId}`;

      lines.push(
          `  ${label} :active, ${id}, ${startIso}, ${dur}`
      );
    }

    lines.push("");
  }

  return lines.join("\n");
}


function parseIsoMs(s) {
  const t = Date.parse(String(s || ""));
  return Number.isFinite(t) ? t : null;
}

function isFiniteNumber(x) {
  return typeof x === "number" && Number.isFinite(x);
}

function serviceFromPod(podName) {
  const p = String(podName || "").trim();
  if (!p) return "unknown";
  return p.split("-")[0] || p;
}

function shortMethod(m) {
  const s = String(m || "");
  const parts = s.split("/");
  return parts.length ? parts[parts.length - 1] || s : s;
}

function sanitizeText(s) {
  return String(s).replace(/[:;]/g, " ").replace(/\s+/g, " ").trim();
}

function firstNonNullIso(blocks) {
  for (const b of blocks) {
    if (b.startMs != null) return new Date(b.startMs).toISOString();
  }
  return null;
}

</script>

<style scoped>
.mermaid-shell {
  overflow: auto;
}

.mermaid {
  margin: 0 auto;
  overflow: auto;
}

/* Most important for gantt: do NOT force svg to fit */
.mermaid :deep(svg) {
  max-width: none !important;
  height: auto;
  width: max-content !important; /* gantt should be wide and scroll */
}
</style>