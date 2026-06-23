import { http, HttpResponse } from 'msw';
import userEvent from '@testing-library/user-event';
import { screen, waitFor } from '@testing-library/vue';
import { server } from '../test/msw/server';
import { renderWithApp } from '../test/render';
import GuestBookingPage from './GuestBookingPage.vue';

describe('GuestBookingPage', () => {
  it('loads and shows available event types (success)', async () => {
    renderWithApp(GuestBookingPage);
    expect(await screen.findByText('Intro call')).toBeInTheDocument();
    expect(screen.getByText('Strategy session')).toBeInTheDocument();
  });

  it('shows an empty state when there are no event types', async () => {
    server.use(http.get('*/event-types', () => HttpResponse.json([])));
    renderWithApp(GuestBookingPage);
    expect(
      await screen.findByText(/No event types are available/i),
    ).toBeInTheDocument();
  });

  it('shows an error state and retries on demand', async () => {
    server.use(
      http.get('*/event-types', () =>
        HttpResponse.json({ code: 'x', message: 'boom' }, { status: 500 }),
      ),
    );
    renderWithApp(GuestBookingPage);
    const retry = await screen.findByRole('button', { name: /try again/i });

    server.use(http.get('*/event-types', () => HttpResponse.json([
      { id: 'evt_intro', title: 'Intro call', description: 'd', durationMinutes: 30 },
    ])));
    await userEvent.click(retry);
    expect(await screen.findByText('Intro call')).toBeInTheDocument();
  });

  it('completes the booking happy path', async () => {
    renderWithApp(GuestBookingPage);

    await userEvent.click(await screen.findByRole('button', { name: /Intro call/i }));
    await userEvent.click(await screen.findByRole('button', { name: /09:00/ }));

    await userEvent.type(await screen.findByLabelText('Name'), 'Jane Doe');
    await userEvent.type(screen.getByLabelText('Email'), 'jane@example.com');
    await userEvent.click(screen.getByRole('button', { name: /confirm booking/i }));

    expect(await screen.findByText('Booking confirmed')).toBeInTheDocument();
  });

  it('handles a 409 conflict: shows a notice and drops the selected slot', async () => {
    renderWithApp(GuestBookingPage);

    await userEvent.click(await screen.findByRole('button', { name: /Intro call/i }));
    await userEvent.click(await screen.findByRole('button', { name: /09:00/ }));
    await userEvent.type(await screen.findByLabelText('Name'), 'Jane Doe');
    await userEvent.type(screen.getByLabelText('Email'), 'jane@example.com');

    server.use(
      http.post('*/bookings', () =>
        HttpResponse.json(
          { code: 'booking_conflict', message: 'That time was just booked.' },
          { status: 409 },
        ),
      ),
    );
    await userEvent.click(screen.getByRole('button', { name: /confirm booking/i }));

    expect(await screen.findByText(/just booked/i)).toBeInTheDocument();
    await waitFor(() =>
      expect(
        screen.queryByRole('button', { name: /confirm booking/i }),
      ).toBeNull(),
    );
  });
});
