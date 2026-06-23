# Call Booking — Frontend

Vue 3 + Vite SPA for the Call Booking MVP.

It talks to the backend **only** through the public API contract
(`openapi/openapi.yaml`, generated from `api/main.tsp`). It does not import any
backend code. During development it runs against a [Prism](https://stoplight.io/open-source/prism)
mock of that contract.

## Stack

- Vue 3 + Vite + TypeScript
- Vue Router
- PrimeVue (Aura theme) for UI components
- Vitest + Testing Library + MSW for tests (see `docs/adr/0001-frontend-testing-strategy.md`)

The stack follows the repo's documented decision in `AGENTS.md` / `README.md`.

## Prerequisites

- Node `>=20` (the repo pins `22.22.2` in `.nvmrc`). With nvm: `nvm use`.

## Setup

```bash
cd frontend
npm install
```

Configure the API base URL (optional in dev — it defaults to the Prism port):

```bash
cp .env.example .env   # VITE_API_BASE_URL=http://localhost:4010
```

## Run

Start the contract mock (in a separate terminal, from the repo root):

```bash
make api-mock          # npx @stoplight/prism-cli mock openapi/openapi.yaml
```

Start the dev server:

```bash
npm run dev            # or, from repo root: make frontend-dev
```

Open http://localhost:5173.

## Scripts

| Command             | What it does                          |
| ------------------- | ------------------------------------- |
| `npm run dev`       | Vite dev server                       |
| `npm run build`     | Type-check (`vue-tsc`) + production build |
| `npm run typecheck` | Type-check only                       |
| `npm run preview`   | Preview the production build          |
| `npm test`          | Run the test suite once (`TZ=UTC`)    |
| `npm run test:watch`| Run tests in watch mode               |

## Structure

```txt
frontend/
  src/
    api/         API layer (contract-aligned types, HTTP client, resource modules)
    components/  Layout + reusable async-state wrapper
    lib/         Date/time helpers
    pages/       Route pages (guest booking + host pages)
    router/      Vue Router setup
    test/        Test setup, fixtures, MSW handlers
    App.vue
    main.ts
```

## Pages

- `/` — Guest booking flow: choose event type → pick a slot → enter details → confirm.
- `/host/event-types` — create event types and list existing ones.
- `/host/availability` — publish an availability range.
- `/host/bookings` — view upcoming bookings.

Each data view has loading / error / empty / success states.

## Testing

```bash
npm test                # or, from repo root: make frontend-test (alias: make test)
```

Vitest + Testing Library run in jsdom with MSW mocking the API at the network
boundary, so the real `src/api` client runs against contract-shaped responses.
Tests run with `TZ=UTC` for deterministic times. CI runs them in
`.github/workflows/frontend-ci.yml`. See `docs/adr/0001-frontend-testing-strategy.md`
for scope and rationale.

## Notes

- All HTTP goes through `src/api/`; components never build URLs or call `fetch`.
- The base URL comes from `VITE_API_BASE_URL`; there is no hardcoded backend URL
  in components.
- The contract is the source of truth. If the UI and the contract disagree, fix
  the contract first (`api/main.tsp`), then regenerate OpenAPI.
