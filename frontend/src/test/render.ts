import { render } from '@testing-library/vue';
import type { Component } from 'vue';
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import {
  createMemoryHistory,
  createRouter,
  type RouteRecordRaw,
} from 'vue-router';

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'book', component: { template: '<div />' } },
  { path: '/host/event-types', name: 'host-event-types', component: { template: '<div />' } },
  { path: '/host/availability', name: 'host-availability', component: { template: '<div />' } },
  { path: '/host/bookings', name: 'host-bookings', component: { template: '<div />' } },
];

/** Render a component with the PrimeVue plugin and a memory router installed. */
export function renderWithApp(component: Component) {
  const router = createRouter({ history: createMemoryHistory(), routes });
  return render(component, {
    global: {
      plugins: [router, [PrimeVue, { theme: { preset: Aura } }]],
    },
  });
}
