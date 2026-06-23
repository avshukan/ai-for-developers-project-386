import { http, HttpResponse } from 'msw';
import { booking, bookings, eventTypes, slots } from '../fixtures';

/**
 * Default happy-path handlers. Paths use a `*` prefix so they match regardless
 * of the configured API base URL. Individual tests override with `server.use`.
 */
export const handlers = [
  http.get('*/event-types', () => HttpResponse.json(eventTypes)),
  http.get('*/event-types/:eventTypeId/slots', () => HttpResponse.json(slots)),
  http.get('*/host/bookings', () => HttpResponse.json(bookings)),

  http.post('*/host/event-types', async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ id: 'evt_new', ...body }, { status: 201 });
  }),
  http.post('*/host/availability', async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(
      { id: 'avl_new', createdAt: '2026-06-23T08:00:00Z', ...body },
      { status: 201 },
    );
  }),
  http.post('*/bookings', async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ ...booking, ...body }, { status: 201 });
  }),
];
