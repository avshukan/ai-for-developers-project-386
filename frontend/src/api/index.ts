/** Public surface of the API layer. Components import from here. */
export { ApiError, apiBaseUrl } from './client';
export { listEventTypes, createEventType } from './eventTypes';
export { createAvailability } from './availability';
export { listSlots } from './slots';
export { createBooking, listHostBookings } from './bookings';
export type {
  EventType,
  CreateEventTypeRequest,
  Availability,
  CreateAvailabilityRequest,
  Slot,
  SlotStatus,
  Booking,
  CreateBookingRequest,
  ApiErrorCode,
  ApiErrorBody,
} from './types';
