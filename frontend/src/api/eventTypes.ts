/** Event Type endpoints (guest list + host create). */
import { apiClient } from './client';
import type { CreateEventTypeRequest, EventType } from './types';

/** GET /event-types — all event types visible to guests. */
export function listEventTypes(): Promise<EventType[]> {
  return apiClient.get<EventType[]>('/event-types');
}

/** POST /host/event-types — create a new event type. */
export function createEventType(
  request: CreateEventTypeRequest,
): Promise<EventType> {
  return apiClient.post<EventType>('/host/event-types', request);
}
