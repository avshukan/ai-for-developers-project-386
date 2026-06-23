<script setup lang="ts">
import { onMounted, ref } from 'vue';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import Message from 'primevue/message';
import Tag from 'primevue/tag';
import AsyncSection from '../components/AsyncSection.vue';
import {
  ApiError,
  createEventType,
  listEventTypes,
  type EventType,
} from '../api';

/** MVP event types are fixed to 30 minutes (see docs/domain.md). */
const MVP_DURATION = 30;

const eventTypes = ref<EventType[] | null>(null);
const loading = ref(true);
const loadError = ref<string | null>(null);

const title = ref('');
const description = ref('');
const submitting = ref(false);
const formError = ref<string | null>(null);
const formDetails = ref<string[]>([]);
const created = ref<EventType | null>(null);

async function load() {
  loading.value = true;
  loadError.value = null;
  try {
    eventTypes.value = await listEventTypes();
  } catch (error) {
    loadError.value =
      error instanceof ApiError ? error.message : 'Failed to load event types.';
  } finally {
    loading.value = false;
  }
}

async function submit() {
  formError.value = null;
  formDetails.value = [];
  created.value = null;

  if (!title.value.trim() || !description.value.trim()) {
    formError.value = 'Title and description are required.';
    return;
  }

  submitting.value = true;
  try {
    const eventType = await createEventType({
      title: title.value.trim(),
      description: description.value.trim(),
      durationMinutes: MVP_DURATION,
    });
    created.value = eventType;
    title.value = '';
    description.value = '';
    eventTypes.value = eventTypes.value
      ? [eventType, ...eventTypes.value]
      : [eventType];
  } catch (error) {
    if (error instanceof ApiError) {
      formError.value = error.message;
      formDetails.value = error.details ?? [];
    } else {
      formError.value = 'Failed to create event type.';
    }
  } finally {
    submitting.value = false;
  }
}

onMounted(load);
</script>

<template>
  <section>
    <h1 class="page__title">Event types</h1>
    <p class="page__subtitle">
      Create the meeting types guests can book. Host pages are public in the MVP.
    </p>

    <div class="section">
      <h2 class="section__title">Create an event type</h2>
      <form class="form" @submit.prevent="submit">
        <div class="field">
          <label for="et-title">Title</label>
          <InputText
            id="et-title"
            v-model="title"
            placeholder="Intro call"
            autocomplete="off"
          />
        </div>

        <div class="field">
          <label for="et-desc">Description</label>
          <Textarea
            id="et-desc"
            v-model="description"
            rows="3"
            auto-resize
            placeholder="A short introduction call."
          />
        </div>

        <div class="field">
          <label for="et-duration">Duration</label>
          <InputText id="et-duration" model-value="30 minutes" disabled />
          <span class="field__hint">Fixed at 30 minutes in the MVP.</span>
        </div>

        <Message v-if="formError" severity="error" :closable="false">
          {{ formError }}
          <ul v-if="formDetails.length" style="margin: 0.5rem 0 0; padding-left: 1.1rem">
            <li v-for="detail in formDetails" :key="detail">{{ detail }}</li>
          </ul>
        </Message>

        <Message v-if="created" severity="success" :closable="false">
          Created "{{ created.title }}".
        </Message>

        <div class="form__actions">
          <Button
            type="submit"
            label="Create event type"
            icon="pi pi-plus"
            :loading="submitting"
          />
        </div>
      </form>
    </div>

    <div class="section">
      <h2 class="section__title">Existing event types</h2>
      <AsyncSection
        :loading="loading"
        :error="loadError"
        :empty="eventTypes !== null && eventTypes.length === 0"
        empty-message="No event types yet. Create the first one above."
        @retry="load"
      >
        <div class="grid">
          <article
            v-for="eventType in eventTypes ?? []"
            :key="eventType.id"
            class="selectable-card"
            style="cursor: default"
          >
            <div class="selectable-card__title">{{ eventType.title }}</div>
            <div class="selectable-card__desc">{{ eventType.description }}</div>
            <Tag :value="`${eventType.durationMinutes} min`" severity="secondary" />
          </article>
        </div>
      </AsyncSection>
    </div>
  </section>
</template>
