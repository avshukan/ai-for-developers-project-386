<script setup lang="ts">
import Button from 'primevue/button';
import { useRouter } from 'vue-router';

const router = useRouter();
</script>

<template>
  <section class="stack" style="align-items: flex-start">
    <h1 class="page__title">Page not found</h1>
    <p class="page__subtitle">The page you're looking for doesn't exist.</p>
    <Button label="Go to booking" icon="pi pi-home" @click="router.push('/')" />
  </section>
</template>
