/** Booking endpoints (guest creates a booking; host lists upcoming bookings). */
import { apiClient } from './client';
import type { Booking, CreateBookingRequest } from './types';

/** POST /bookings — create a booking for a selected slot. */
export function createBooking(
  request: CreateBookingRequest,
): Promise<Booking> {
  return apiClient.post<Booking>('/bookings', request);
}

/** GET /host/bookings — upcoming bookings across all event types. */
export function listHostBookings(): Promise<Booking[]> {
  return apiClient.get<Booking[]>('/host/bookings');
}
