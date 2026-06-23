/**
 * API types.
 *
 * These mirror the API contract — the single source of truth — defined in
 * `api/main.tsp` and generated to `openapi/openapi.yaml` at the repo root.
 * They are hand-written (not imported from the backend) so the frontend stays
 * decoupled from the backend implementation. Keep them in sync with the
 * contract; do not add fields that are not in the contract.
 */

/** A bookable meeting type created by the host. */
export interface EventType {
  id: string;
  title: string;
  description: string;
  /** Fixed to 30 minutes in the MVP. */
  durationMinutes: number;
}

/** Request body for creating an event type. */
export interface CreateEventTypeRequest {
  title: string;
  description: string;
  /** Must be 30 in the MVP. */
  durationMinutes: number;
}

/** A time range published by the host. */
export interface Availability {
  id: string;
  /** Availability start time in UTC (ISO 8601). */
  startAt: string;
  /** Availability end time in UTC (ISO 8601). */
  endAt: string;
  /** Creation time in UTC (ISO 8601). */
  createdAt: string;
}

/** Request body for creating an availability range. */
export interface CreateAvailabilityRequest {
  startAt: string;
  endAt: string;
}

export type SlotStatus = 'available' | 'booked';

/** A 30-minute slot derived from host availability. */
export interface Slot {
  /** Slot start time in UTC (ISO 8601). */
  startAt: string;
  /** Slot end time in UTC (ISO 8601). */
  endAt: string;
  status: SlotStatus;
}

/** A confirmed booking created by a guest. */
export interface Booking {
  id: string;
  eventTypeId: string;
  /** Booking start time in UTC (ISO 8601). */
  startAt: string;
  /** Booking end time in UTC (ISO 8601). */
  endAt: string;
  guestName: string;
  guestEmail: string;
  /** Creation time in UTC (ISO 8601). */
  createdAt: string;
}

/** Request body for creating a booking. The backend derives endAt as startAt + 30 min. */
export interface CreateBookingRequest {
  eventTypeId: string;
  startAt: string;
  guestName: string;
  guestEmail: string;
}

/** Error `code` values returned by the contract's error responses. */
export type ApiErrorCode =
  | 'validation_error'
  | 'not_found'
  | 'booking_conflict';

/** Shape of the error bodies defined in the contract (400 / 404 / 409). */
export interface ApiErrorBody {
  code: ApiErrorCode;
  message: string;
  details?: string[];
}
