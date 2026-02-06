<template>
  <div class="json-viewer">
    <div class="row">
      <div class="title" v-if="title">{{ title }}</div>
      <v-spacer />
      <v-btn
          v-if="copyable"
          size="x-small"
          variant="text"
          icon="mdi-content-copy"
          @click="copyJson"
      />
      <v-btn
          v-if="collapsible"
          size="x-small"
          variant="text"
          icon="mdi-unfold-more-horizontal"
          @click="toggleAll"
      />
    </div>

    <VueJsonPretty
        :data="data ?? {}"
        :deep="deep"
        :showLine="true"
        :showLength="true"
        :collapsed="collapsed"
    />
  </div>
</template>

<script setup>
import { ref } from "vue";
import VueJsonPretty from "vue-json-pretty";
import "vue-json-pretty/lib/styles.css";

const props = defineProps({
  data: { type: [Object, Array, String, Number, Boolean, null], default: null },
  title: { type: String, default: "" },
  deep: { type: Number, default: 6 },
  copyable: { type: Boolean, default: true },
  collapsible: { type: Boolean, default: true },
  defaultCollapsed: { type: Boolean, default: true },
});

const collapsed = ref(props.defaultCollapsed);

const toggleAll = () => {
  collapsed.value = !collapsed.value;
};

const copyJson = async () => {
  try {
    const text = JSON.stringify(props.data ?? {}, null, 2);
    await navigator.clipboard.writeText(text);
  } catch (e) {
    // ignore; clipboard may be blocked in some contexts
  }
};
</script>

<style scoped>
.json-viewer {
  border-radius: 8px;
  padding: 8px 10px;
  background: rgba(0,0,0,0.04);
}
.row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}
.title {
  font-weight: 600;
  font-size: 0.85rem;
}
</style>