# Integration tests (Playwright)

End-to-end tests that drive the **real** Vue SPA against the **real** Go backend
over the public API contract. Nothing is mocked here — this is the layer that
proves the frontend and backend work together.

For the why (tool, location, startup) see
[`docs/adr/0003-integration-tests-playwright.md`](../docs/adr/0003-integration-tests-playwright.md).

## What is covered

Scenarios are defined in [`docs/product.md`](../docs/product.md); the tests trace
to them rather than redefine them:

| Test (`tests/booking.spec.ts`)                     | Scenario in `docs/product.md` |
| -------------------------------------------------- | ----------------------------- |
| guest books an available slot end to end           | Main Guest Scenario           |
| a booked slot becomes unavailable (no double book) | Slot / Booking rules          |
| host sees a new booking on the bookings page        | Main Host Bookings Scenario   |

## How it runs

Playwright's `webServer` starts both services for the run (see
`playwright.config.ts`):

- **Backend** — `go run ./cmd/server` on port `8080` with `SEED_DATA=true`, so
  demo event types and availability exist.
- **Frontend** — Vite dev server on port `5173` with
  `VITE_API_BASE_URL=http://localhost:8080`, so the SPA calls the local backend
  instead of the Prism mock.

The browser is pinned to `timezoneId: 'UTC'` so the UTC times from the API render
deterministically.

## Running locally

Requires Node (see `../.nvmrc`) and Go (see `../backend/go.mod`).

```bash
# from repo root
make e2e-install   # install deps + the Chromium browser
make e2e           # run the suite (boots backend + frontend automatically)
```

Or directly:

```bash
cd e2e
npm install
npm run install:browsers
npm test
npm run report      # open the HTML report from the last run
```

You do **not** need to start the servers yourself — Playwright does. If you
already have them running locally, `reuseExistingServer` picks them up outside CI.
