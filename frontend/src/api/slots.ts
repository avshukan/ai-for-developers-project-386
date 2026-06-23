/** Slot endpoints (guest views available slots for an event type). */
import { apiClient } from './client';
import type { Slot } from './types';

/**
 * GET /event-types/{eventTypeId}/slots — available 30-minute slots for the
 * next 14 days. Slots are global; `eventTypeId` validates the event type exists
 * and scopes the guest flow, it does not filter the slots.
 */
export function listSlots(eventTypeId: string): Promise<Slot[]> {
  return apiClient.get<Slot[]>(
    `/event-types/${encodeURIComponent(eventTypeId)}/slots`,
  );
}
