import { createApp } from 'vue';
import App from './App.vue';
import vuetify from './plugins/vuetify';
import apolloProvider from './apollo'; // Adjust the path as needed
import 'vuetify/styles';
import '@mdi/font/css/materialdesignicons.css';
import './assets/styles/global.css'

createApp(App)
    .use(vuetify)
    .use(apolloProvider)
    .mount('#app');
