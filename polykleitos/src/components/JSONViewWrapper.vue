<template>
  <v-card class="json-viewer">
    <v-card-text>
      <pre>{{ formattedJson }}</pre>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  jsonString: {
    type: String,
    required: true
  }
});

const formattedJson = computed(() => {
  try {
    const parsed = JSON.parse(props.jsonString);
    return JSON.stringify(parsed, null, 2);
  } catch (e) {
    return props.jsonString;
  }
});
</script>

<style scoped>
.json-viewer {
  background-color: #f5f5f5;
  margin: 1em 0;
}

pre {
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
