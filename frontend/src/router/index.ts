import { createRouter, createWebHistory } from 'vue-router';

/**
 * Routes map to the user scenarios in docs/product.md:
 *  - guest booking flow (default route);
 *  - host event types / availability / bookings pages.
 * Pages are lazy-loaded to keep the initial bundle small.
 */
export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'book',
      component: () => import('../pages/GuestBookingPage.vue'),
      meta: { title: 'Book a meeting' },
    },
    {
      path: '/host/event-types',
      name: 'host-event-types',
      component: () => import('../pages/HostEventTypesPage.vue'),
      meta: { title: 'Event types' },
    },
    {
      path: '/host/availability',
      name: 'host-availability',
      component: () => import('../pages/HostAvailabilityPage.vue'),
      meta: { title: 'Availability' },
    },
    {
      path: '/host/bookings',
      name: 'host-bookings',
      component: () => import('../pages/HostBookingsPage.vue'),
      meta: { title: 'Bookings' },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../pages/NotFoundPage.vue'),
      meta: { title: 'Not found' },
    },
  ],
});
