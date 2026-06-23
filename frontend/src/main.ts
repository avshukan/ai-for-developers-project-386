import { createApp } from 'vue';
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import 'primeicons/primeicons.css';

import App from './App.vue';
import { router } from './router';
import './styles.css';

const app = createApp(App);

app.use(router);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    // Keep a stable light theme for the MVP (selector is never applied).
    options: { darkModeSelector: '.app-dark' },
  },
});

app.mount('#app');
