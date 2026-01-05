<template>
  <v-container fluid>
    <v-app-bar color="primary" prominent>
      <v-app-bar-title>Polykleitos</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon>mdi-magnify</v-icon>
      </v-btn>
    </v-app-bar>

    <v-container class="mt-8">
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card elevation="2">
            <v-card-title class="text-h4 text-center pa-6">
              Welcome to Polykleitos
            </v-card-title>
            <v-card-text class="text-center">
              <p class="text-h6 mb-4">A new UI built with Vue 3 + Vuetify 3</p>
              <v-divider class="my-4"></v-divider>
              <p class="text-body-1">This is your homepage scaffold. Start building your application here.</p>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row class="mt-6">
        <v-col cols="12" sm="6" md="4" v-for="n in 3" :key="n">
          <v-card hover>
            <v-card-title>Feature {{ n }}</v-card-title>
            <v-card-text>
              Add your content here for feature {{ n }}
            </v-card-text>
            <v-card-actions>
              <v-btn color="primary" variant="text">Learn More</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-container>
</template>

<script setup>
// Add your component logic here
</script>
