<script setup lang="ts">
import { ref } from 'vue';
import DatePicker from 'primevue/datepicker';
import Button from 'primevue/button';
import Message from 'primevue/message';
import {
  ApiError,
  createAvailability,
  type Availability,
} from '../api';
import { formatDateTime } from '../lib/datetime';

const startAt = ref<Date | null>(null);
const endAt = ref<Date | null>(null);
const submitting = ref(false);
const formError = ref<string | null>(null);
const formDetails = ref<string[]>([]);
const created = ref<Availability | null>(null);

async function submit() {
  formError.value = null;
  formDetails.value = [];
  created.value = null;

  if (!startAt.value || !endAt.value) {
    formError.value = 'Both start and end time are required.';
    return;
  }
  if (endAt.value.getTime() <= startAt.value.getTime()) {
    formError.value = 'End time must be after start time.';
    return;
  }

  submitting.value = true;
  try {
    created.value = await createAvailability({
      startAt: startAt.value.toISOString(),
      endAt: endAt.value.toISOString(),
    });
    startAt.value = null;
    endAt.value = null;
  } catch (error) {
    if (error instanceof ApiError) {
      formError.value = error.message;
      formDetails.value = error.details ?? [];
    } else {
      formError.value = 'Failed to publish availability.';
    }
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <section>
    <h1 class="page__title">Availability</h1>
    <p class="page__subtitle">
      Publish a time range. The backend splits it into bookable 30-minute slots.
    </p>

    <form class="form" @submit.prevent="submit">
      <div class="field">
        <label for="start">Start time</label>
        <DatePicker
          input-id="start"
          v-model="startAt"
          show-time
          hour-format="24"
          show-icon
          fluid
          placeholder="Select start"
        />
      </div>

      <div class="field">
        <label for="end">End time</label>
        <DatePicker
          input-id="end"
          v-model="endAt"
          show-time
          hour-format="24"
          show-icon
          fluid
          placeholder="Select end"
        />
        <span class="field__hint">
          Times are entered in your local zone and sent to the API as UTC.
        </span>
      </div>

      <Message v-if="formError" severity="error" :closable="false">
        {{ formError }}
        <ul v-if="formDetails.length" style="margin: 0.5rem 0 0; padding-left: 1.1rem">
          <li v-for="detail in formDetails" :key="detail">{{ detail }}</li>
        </ul>
      </Message>

      <Message v-if="created" severity="success" :closable="false">
        Availability published: {{ formatDateTime(created.startAt) }} →
        {{ formatDateTime(created.endAt) }}.
      </Message>

      <div class="form__actions">
        <Button
          type="submit"
          label="Publish availability"
          icon="pi pi-calendar-plus"
          :loading="submitting"
        />
      </div>
    </form>
  </section>
</template>
