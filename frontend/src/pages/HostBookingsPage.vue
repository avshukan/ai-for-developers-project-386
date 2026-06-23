<script setup lang="ts">
import { onMounted, ref } from 'vue';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import AsyncSection from '../components/AsyncSection.vue';
import { ApiError, listHostBookings, type Booking } from '../api';
import { formatDateTime } from '../lib/datetime';

const bookings = ref<Booking[] | null>(null);
const loading = ref(false);
const loadError = ref<string | null>(null);

async function load() {
  loading.value = true;
  loadError.value = null;
  try {
    bookings.value = await listHostBookings();
  } catch (error) {
    loadError.value =
      error instanceof ApiError ? error.message : 'Failed to load bookings.';
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>

<template>
  <section>
    <div
      style="
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 1rem;
        flex-wrap: wrap;
      "
    >
      <div>
        <h1 class="page__title">Upcoming bookings</h1>
        <p class="page__subtitle">Bookings across all event types.</p>
      </div>
      <Button
        label="Refresh"
        icon="pi pi-refresh"
        severity="secondary"
        outlined
        :loading="loading"
        @click="load"
      />
    </div>

    <AsyncSection
      :loading="loading"
      :error="loadError"
      :empty="(bookings?.length ?? 0) === 0"
      empty-message="No upcoming bookings yet."
      @retry="load"
    >
      <DataTable :value="bookings ?? []" striped-rows responsive-layout="scroll">
        <Column field="eventTypeId" header="Event type" />
        <Column header="Start">
          <template #body="{ data }: { data: Booking }">
            {{ formatDateTime(data.startAt) }}
          </template>
        </Column>
        <Column header="End">
          <template #body="{ data }: { data: Booking }">
            {{ formatDateTime(data.endAt) }}
          </template>
        </Column>
        <Column field="guestName" header="Guest" />
        <Column field="guestEmail" header="Email" />
      </DataTable>
    </AsyncSection>
  </section>
</template>
