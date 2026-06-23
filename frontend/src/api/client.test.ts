import { http, HttpResponse } from 'msw';
import { server } from '../test/msw/server';
import { ApiError, apiClient } from './client';

describe('apiClient', () => {
  it('parses a successful JSON response', async () => {
    server.use(http.get('*/ping', () => HttpResponse.json({ ok: true })));
    await expect(apiClient.get('/ping')).resolves.toEqual({ ok: true });
  });

  it('throws ApiError with status/code/message from a contract error body', async () => {
    server.use(
      http.post('*/bookings', () =>
        HttpResponse.json(
          { code: 'booking_conflict', message: 'Slot taken' },
          { status: 409 },
        ),
      ),
    );
    const call = () => apiClient.post('/bookings', {});
    await expect(call()).rejects.toBeInstanceOf(ApiError);
    await expect(call()).rejects.toMatchObject({
      status: 409,
      code: 'booking_conflict',
      message: 'Slot taken',
    });
  });

  it('surfaces validation details', async () => {
    server.use(
      http.post('*/host/event-types', () =>
        HttpResponse.json(
          { code: 'validation_error', message: 'Invalid', details: ['title required'] },
          { status: 400 },
        ),
      ),
    );
    await expect(apiClient.post('/host/event-types', {})).rejects.toMatchObject({
      status: 400,
      details: ['title required'],
    });
  });

  it('maps a network failure to ApiError(status 0)', async () => {
    server.use(http.get('*/event-types', () => HttpResponse.error()));
    await expect(apiClient.get('/event-types')).rejects.toMatchObject({ status: 0 });
  });
});
