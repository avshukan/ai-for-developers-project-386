import type { Booking, EventType, Slot } from '../api/types';

/** Contract-shaped fixtures (mirror openapi/openapi.yaml). */

export const eventTypes: EventType[] = [
  {
    id: 'evt_intro',
    title: 'Intro call',
    description: 'A 30-minute introduction call.',
    durationMinutes: 30,
  },
  {
    id: 'evt_strategy',
    title: 'Strategy session',
    description: 'Discuss your project roadmap.',
    durationMinutes: 30,
  },
];

// Two distinct UTC days so day-grouping produces two groups (tests run TZ=UTC).
export const slots: Slot[] = [
  { startAt: '2026-06-24T09:00:00Z', endAt: '2026-06-24T09:30:00Z', status: 'available' },
  { startAt: '2026-06-24T09:30:00Z', endAt: '2026-06-24T10:00:00Z', status: 'available' },
  { startAt: '2026-06-25T14:00:00Z', endAt: '2026-06-25T14:30:00Z', status: 'available' },
];

export const booking: Booking = {
  id: 'bkg_1',
  eventTypeId: 'evt_intro',
  startAt: '2026-06-24T09:00:00Z',
  endAt: '2026-06-24T09:30:00Z',
  guestName: 'Jane Doe',
  guestEmail: 'jane@example.com',
  createdAt: '2026-06-23T08:30:00Z',
};

export const bookings: Booking[] = [booking];
