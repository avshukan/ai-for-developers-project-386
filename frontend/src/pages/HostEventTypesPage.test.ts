import { http, HttpResponse } from 'msw';
import userEvent from '@testing-library/user-event';
import { screen } from '@testing-library/vue';
import { server } from '../test/msw/server';
import { renderWithApp } from '../test/render';
import HostEventTypesPage from './HostEventTypesPage.vue';

describe('HostEventTypesPage', () => {
  it('lists existing event types', async () => {
    renderWithApp(HostEventTypesPage);
    expect(await screen.findByText('Intro call')).toBeInTheDocument();
  });

  it('shows an empty state when there are none', async () => {
    server.use(http.get('*/event-types', () => HttpResponse.json([])));
    renderWithApp(HostEventTypesPage);
    expect(await screen.findByText(/No event types yet/i)).toBeInTheDocument();
  });

  it('creates an event type and adds it to the list', async () => {
    server.use(http.get('*/event-types', () => HttpResponse.json([])));
    renderWithApp(HostEventTypesPage);
    await screen.findByText(/No event types yet/i);

    await userEvent.type(screen.getByLabelText('Title'), 'New call');
    await userEvent.type(screen.getByLabelText('Description'), 'Desc');
    await userEvent.click(
      screen.getByRole('button', { name: /create event type/i }),
    );

    expect(await screen.findByText(/Created "New call"/)).toBeInTheDocument();
    expect(screen.getByText('New call')).toBeInTheDocument();
  });
});
