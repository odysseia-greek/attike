<template>
  <div ref="mermaidContainer" class="mermaid-diagram"></div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue';
import mermaid from 'mermaid';

const props = defineProps({
  mermaidDiagram: {
    type: String,
    required: true
  }
});

const mermaidContainer = ref(null);

onMounted(() => {
  mermaid.initialize({ startOnLoad: false });
  renderDiagram();
});

watch(() => props.mermaidDiagram, () => {
  renderDiagram();
});

async function renderDiagram() {
  if (mermaidContainer.value && props.mermaidDiagram) {
    mermaidContainer.value.innerHTML = '';
    await nextTick();

    try {
      const { svg } = await mermaid.render('mermaid-' + Date.now(), props.mermaidDiagram);
      mermaidContainer.value.innerHTML = svg;
    } catch (error) {
      console.error('Mermaid rendering error:', error);
      mermaidContainer.value.innerHTML = '<p>Error rendering diagram</p>';
    }
  }
}
</script>

<style scoped>
.mermaid-diagram {
  width: 100%;
  overflow-x: auto;
}
</style>
