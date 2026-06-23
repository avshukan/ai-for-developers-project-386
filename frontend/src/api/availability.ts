/** Availability endpoints (host publishes available time). */
import { apiClient } from './client';
import type { Availability, CreateAvailabilityRequest } from './types';

/** POST /host/availability — publish a host availability range. */
export function createAvailability(
  request: CreateAvailabilityRequest,
): Promise<Availability> {
  return apiClient.post<Availability>('/host/availability', request);
}
