<script setup lang="ts">
/**
 * Renders the standard async UI states so pages don't repeat them:
 *  - loading: a spinner;
 *  - error: an error message (with optional retry);
 *  - empty: a friendly empty state;
 *  - success: the default slot.
 * Precedence: loading > error > empty > success.
 */
import ProgressSpinner from 'primevue/progressspinner';
import Message from 'primevue/message';
import Button from 'primevue/button';

defineProps<{
  loading: boolean;
  error?: string | null;
  empty?: boolean;
  loadingMessage?: string;
  emptyMessage?: string;
}>();

const emit = defineEmits<{ (e: 'retry'): void }>();
</script>

<template>
  <div class="async">
    <div v-if="loading" class="async__state">
      <ProgressSpinner
        style="width: 40px; height: 40px"
        stroke-width="4"
        aria-label="Loading"
      />
      <span class="async__muted">{{ loadingMessage ?? 'Loading…' }}</span>
    </div>

    <div v-else-if="error" class="async__state">
      <Message severity="error" :closable="false" class="async__message">
        {{ error }}
      </Message>
      <Button
        label="Try again"
        icon="pi pi-refresh"
        severity="secondary"
        outlined
        size="small"
        @click="emit('retry')"
      />
    </div>

    <div v-else-if="empty" class="async__state">
      <i class="pi pi-inbox async__empty-icon" aria-hidden="true" />
      <span class="async__muted">{{ emptyMessage ?? 'Nothing here yet.' }}</span>
    </div>

    <slot v-else />
  </div>
</template>
