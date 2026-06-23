<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Message from 'primevue/message';
import Tag from 'primevue/tag';
import AsyncSection from '../components/AsyncSection.vue';
import {
  ApiError,
  createBooking,
  listEventTypes,
  listSlots,
  type Booking,
  type EventType,
  type Slot,
} from '../api';
import { dayKey, formatDateTime, formatDay, formatTime } from '../lib/datetime';

/** Minimal client-side email check; the backend is the source of truth. */
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Step 1 — event types
const eventTypes = ref<EventType[] | null>(null);
const eventTypesLoading = ref(false);
const eventTypesError = ref<string | null>(null);

// Step 2 — slots for the selected event type
const selectedEventType = ref<EventType | null>(null);
const slots = ref<Slot[] | null>(null);
const slotsLoading = ref(false);
const slotsError = ref<string | null>(null);

// Step 3 — booking
const selectedSlot = ref<Slot | null>(null);
const guestName = ref('');
const guestEmail = ref('');
const submitting = ref(false);
const bookingError = ref<string | null>(null);
const bookingDetails = ref<string[]>([]);
const confirmed = ref<Booking | null>(null);

const availableSlots = computed(() =>
  (slots.value ?? []).filter((slot) => slot.status === 'available'),
);

interface DayGroup {
  key: string;
  label: string;
  slots: Slot[];
}

const slotsByDay = computed<DayGroup[]>(() => {
  const groups = new Map<string, DayGroup>();
  const sorted = [...availableSlots.value].sort(
    (a, b) => new Date(a.startAt).getTime() - new Date(b.startAt).getTime(),
  );
  for (const slot of sorted) {
    const key = dayKey(slot.startAt);
    const group = groups.get(key);
    if (group) {
      group.slots.push(slot);
    } else {
      groups.set(key, { key, label: formatDay(slot.startAt), slots: [slot] });
    }
  }
  return [...groups.values()];
});

async function loadEventTypes() {
  eventTypesLoading.value = true;
  eventTypesError.value = null;
  try {
    eventTypes.value = await listEventTypes();
  } catch (error) {
    eventTypesError.value =
      error instanceof ApiError ? error.message : 'Failed to load event types.';
  } finally {
    eventTypesLoading.value = false;
  }
}

async function loadSlots(eventTypeId: string) {
  slotsLoading.value = true;
  slotsError.value = null;
  try {
    slots.value = await listSlots(eventTypeId);
  } catch (error) {
    slotsError.value =
      error instanceof ApiError ? error.message : 'Failed to load slots.';
  } finally {
    slotsLoading.value = false;
  }
}

function selectEventType(eventType: EventType) {
  selectedEventType.value = eventType;
  selectedSlot.value = null;
  confirmed.value = null;
  bookingError.value = null;
  slots.value = null;
  void loadSlots(eventType.id);
}

function selectSlot(slot: Slot) {
  selectedSlot.value = slot;
  bookingError.value = null;
}

function isSelectedSlot(slot: Slot): boolean {
  return selectedSlot.value?.startAt === slot.startAt;
}

async function submitBooking() {
  bookingError.value = null;
  bookingDetails.value = [];
  if (!selectedEventType.value || !selectedSlot.value) return;

  if (!guestName.value.trim()) {
    bookingError.value = 'Your name is required.';
    return;
  }
  if (!EMAIL_RE.test(guestEmail.value.trim())) {
    bookingError.value = 'Please enter a valid email address.';
    return;
  }

  submitting.value = true;
  try {
    confirmed.value = await createBooking({
      eventTypeId: selectedEventType.value.id,
      startAt: selectedSlot.value.startAt,
      guestName: guestName.value.trim(),
      guestEmail: guestEmail.value.trim(),
    });
    guestName.value = '';
    guestEmail.value = '';
  } catch (error) {
    if (error instanceof ApiError) {
      bookingError.value = error.message;
      bookingDetails.value = error.details ?? [];
      // The slot is gone (already booked / overlaps): refresh and reselect.
      if (error.status === 409) {
        selectedSlot.value = null;
        void loadSlots(selectedEventType.value.id);
      }
    } else {
      bookingError.value = 'Failed to create booking.';
    }
  } finally {
    submitting.value = false;
  }
}

function bookAnother() {
  confirmed.value = null;
  selectedSlot.value = null;
  bookingError.value = null;
  if (selectedEventType.value) void loadSlots(selectedEventType.value.id);
}

function changeEventType() {
  selectedEventType.value = null;
  selectedSlot.value = null;
  slots.value = null;
  confirmed.value = null;
  bookingError.value = null;
}

onMounted(loadEventTypes);
</script>

<template>
  <section>
    <h1 class="page__title">Book a meeting</h1>
    <p class="page__subtitle">
      Pick an event type, choose a free slot, and leave your details.
    </p>

    <!-- Confirmation -->
    <div v-if="confirmed" class="summary-card stack">
      <div style="display: flex; align-items: center; gap: 0.5rem">
        <i class="pi pi-check-circle" style="color: #10b981; font-size: 1.4rem" />
        <h2 class="section__title" style="margin: 0">Booking confirmed</h2>
      </div>
      <dl class="summary-list">
        <div>
          <dt>Event type</dt>
          <dd>{{ selectedEventType?.title ?? confirmed.eventTypeId }}</dd>
        </div>
        <div>
          <dt>When</dt>
          <dd>
            {{ formatDateTime(confirmed.startAt) }} – {{ formatTime(confirmed.endAt) }}
          </dd>
        </div>
        <div>
          <dt>Guest</dt>
          <dd>{{ confirmed.guestName }} ({{ confirmed.guestEmail }})</dd>
        </div>
        <div>
          <dt>Booking id</dt>
          <dd>{{ confirmed.id }}</dd>
        </div>
      </dl>
      <div class="form__actions">
        <Button label="Book another" icon="pi pi-plus" @click="bookAnother" />
        <Button
          label="Change event type"
          severity="secondary"
          outlined
          @click="changeEventType"
        />
      </div>
    </div>

    <template v-else>
      <!-- Step 1: event types -->
      <div class="section">
        <h2 class="section__title">1. Choose an event type</h2>
        <AsyncSection
          :loading="eventTypesLoading"
          :error="eventTypesError"
          :empty="(eventTypes?.length ?? 0) === 0"
          empty-message="No event types are available yet. Please check back later."
          @retry="loadEventTypes"
        >
          <div class="grid">
            <button
              v-for="eventType in eventTypes ?? []"
              :key="eventType.id"
              type="button"
              class="selectable-card"
              :class="{ 'is-selected': selectedEventType?.id === eventType.id }"
              @click="selectEventType(eventType)"
            >
              <div class="selectable-card__title">{{ eventType.title }}</div>
              <div class="selectable-card__desc">{{ eventType.description }}</div>
              <Tag :value="`${eventType.durationMinutes} min`" severity="secondary" />
            </button>
          </div>
        </AsyncSection>
      </div>

      <!-- Step 2: slots -->
      <div v-if="selectedEventType" class="section">
        <div style="display: flex; align-items: baseline; gap: 0.75rem; flex-wrap: wrap">
          <h2 class="section__title" style="margin: 0">2. Pick a slot</h2>
          <Button
            label="Change event type"
            icon="pi pi-arrow-left"
            text
            size="small"
            @click="changeEventType"
          />
        </div>
        <p class="page__subtitle" style="margin: 0.5rem 0 1rem">
          Available 30-minute slots for "{{ selectedEventType.title }}" over the next
          14 days.
        </p>
        <AsyncSection
          :loading="slotsLoading"
          :error="slotsError"
          :empty="availableSlots.length === 0"
          empty-message="No available slots in the next 14 days."
          @retry="loadSlots(selectedEventType.id)"
        >
          <div class="slots__day" v-for="day in slotsByDay" :key="day.key">
            <h3 class="slots__day-title">{{ day.label }}</h3>
            <div class="slots__list">
              <Button
                v-for="slot in day.slots"
                :key="slot.startAt"
                :label="formatTime(slot.startAt)"
                size="small"
                :outlined="!isSelectedSlot(slot)"
                @click="selectSlot(slot)"
              />
            </div>
          </div>
        </AsyncSection>
      </div>

      <!-- Step 3: guest details -->
      <div v-if="selectedSlot" class="section">
        <h2 class="section__title">3. Your details</h2>
        <p class="page__subtitle" style="margin: 0 0 1rem">
          {{ selectedEventType?.title }} · {{ formatDateTime(selectedSlot.startAt) }}
        </p>
        <form class="form" @submit.prevent="submitBooking">
          <div class="field">
            <label for="guest-name">Name</label>
            <InputText id="guest-name" v-model="guestName" autocomplete="name" />
          </div>
          <div class="field">
            <label for="guest-email">Email</label>
            <InputText
              id="guest-email"
              v-model="guestEmail"
              type="email"
              autocomplete="email"
            />
          </div>

          <Message v-if="bookingError" severity="error" :closable="false">
            {{ bookingError }}
            <ul
              v-if="bookingDetails.length"
              style="margin: 0.5rem 0 0; padding-left: 1.1rem"
            >
              <li v-for="detail in bookingDetails" :key="detail">{{ detail }}</li>
            </ul>
          </Message>

          <div class="form__actions">
            <Button
              type="submit"
              label="Confirm booking"
              icon="pi pi-check"
              :loading="submitting"
            />
            <Button
              label="Cancel"
              severity="secondary"
              text
              type="button"
              @click="selectedSlot = null"
            />
          </div>
        </form>
      </div>
    </template>
  </section>
</template>
