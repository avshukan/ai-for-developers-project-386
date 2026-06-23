<script setup lang="ts">
import { apiBaseUrl } from '../api';

const hostLinks = [
  { to: '/host/event-types', label: 'Event types' },
  { to: '/host/availability', label: 'Availability' },
  { to: '/host/bookings', label: 'Bookings' },
];
</script>

<template>
  <div class="layout">
    <header class="layout__header">
      <div class="layout__bar">
        <router-link to="/" class="layout__brand">
          <i class="pi pi-calendar" aria-hidden="true" />
          <span>Call&nbsp;Booking</span>
        </router-link>

        <nav class="layout__nav" aria-label="Main navigation">
          <router-link to="/" class="layout__link">Book a meeting</router-link>
          <span class="layout__group-label">Host</span>
          <router-link
            v-for="link in hostLinks"
            :key="link.to"
            :to="link.to"
            class="layout__link"
          >
            {{ link.label }}
          </router-link>
        </nav>
      </div>
    </header>

    <main class="layout__main">
      <slot />
    </main>

    <footer class="layout__footer">
      Call Booking MVP — guest &amp; host flows over the public API.
      <span class="layout__api">API: {{ apiBaseUrl }}</span>
    </footer>
  </div>
</template>
