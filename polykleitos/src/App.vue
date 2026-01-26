<script setup>
import HomePage from "@/views/HomePage.vue";
import { ref, onMounted } from 'vue';

const frontendVersion = ref(import.meta.env.VERSION || 'v0.0.1');
const gatewayVersion = ref('loading...');
const gatewayHealth = ref(false);
const lastCheckedTime = ref('');

const fetchGatewayHealth = async () => {
  let url = document.location.origin + '/euripides/v1/health';
  if (process.env.NODE_ENV === 'local') {
    url = 'http://attike.byzantium.odysseia-greek:8080/euripides/v1/health';
  }

  try {
    const response = await fetch(url);
    if (response.ok) {
      const data = await response.json();
      gatewayVersion.value = data.version || 'unknown';
      gatewayHealth.value = data.healthy;
      lastCheckedTime.value = new Date().toLocaleTimeString();
    } else {
      gatewayVersion.value = 'unknown';
      gatewayHealth.value = false;
      lastCheckedTime.value = new Date().toLocaleTimeString();
    }
  } catch (error) {
    console.error('Error fetching gateway health:', error);
    gatewayVersion.value = 'unknown';
    gatewayHealth.value = false;
    lastCheckedTime.value = new Date().toLocaleTimeString();
  }
};

onMounted(() => {
  fetchGatewayHealth();
  setInterval(() => {
    fetchGatewayHealth();
  }, 30000); // 30 seconds
});

const scrollToSearch = () => {
  const searchSection = document.getElementById('search-section');
  if (searchSection) {
    searchSection.scrollIntoView({ behavior: 'smooth' });
  }
};
</script>

<template>
  <v-app id="polykleitos">
    <v-app-bar floating color="footer">
      <v-app-bar-title>Polykleitos</v-app-bar-title>
      <v-btn variant="text" @click="scrollToSearch">
        <v-icon>mdi-graph</v-icon>
        Traces
      </v-btn>
      <v-btn variant="text" @click="scrollToSearch">
        <v-icon>mdi-chart-bar</v-icon>
        Metrics
      </v-btn>
    </v-app-bar>

    <v-main>
      <HomePage />
    </v-main>

    <v-footer color="surface-variant">
      <v-card flat width="100%" class="text-center" color="surface-variant">
        <v-card-text>
          <v-btn
              class="ma-2"
              href="https://github.com/odysseia-greek"
              icon="mdi-github"
              variant="text"
          />
        </v-card-text>

        <v-divider></v-divider>

        <v-card-text>
          {{ new Date().getFullYear() }} — <strong>Odysseia-greek</strong>
        </v-card-text>

        <v-card-text class="mt-4">
          <v-row justify="center">
            <v-col cols="12" md="6">
              <strong>Polykleitos Version: {{ frontendVersion }}</strong>
            </v-col>
            <v-col cols="12" md="6">
              <strong>Euripides Version: {{ gatewayVersion }}</strong>
              <v-icon
                  :color="gatewayHealth ? 'green' : 'red'"
                  class="ml-2"
                  size="16"
              >
                mdi-circle
              </v-icon>
              <small v-if="gatewayHealth"> (Checked: {{ lastCheckedTime }})</small>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-footer>
  </v-app>
</template>

<style>
#polykleitos .v-main {
  padding: 0 !important;
}

#polykleitos .v-main__wrap {
  max-width: 100% !important;
}
</style>
