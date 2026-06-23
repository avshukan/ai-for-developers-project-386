# ADR 1: Frontend testing strategy

## Status

Accepted

## Context

The course requires testing core user scenarios, and the docs listed the testing
stack as `TBD`. We need a first test layer for the Vue 3 + Vite frontend that:

- verifies UI behavior against the API contract, not backend rules;
- keeps the frontend independent of the backend (no backend code, no real server);
- stays small and behavior-focused (this is an MVP).

## Options

1. **Vitest + Testing Library (jsdom) + MSW** — Vite-native runner; user-centric
   component/integration tests; API mocked at the network boundary so the real
   `src/api` client runs against contract-shaped responses.
2. **Stub the API layer (`vi.mock` on `src/api/*`)** — simplest, but bypasses the
   real client and its error/contract mapping.
3. **Tests against the Prism mock** — true integration, but needs a running server;
   slower and flakier, unsuitable for unit/component CI.
4. **Add E2E (Playwright)** — highest fidelity, but heavier and slower than an MVP
   needs now.

## Decision

Use **Vitest + @testing-library/vue + @testing-library/jest-dom + MSW**, running in
**jsdom**, for unit and component/integration tests.

- Unit tests cover pure logic: the API client error mapping (`src/api/client.ts`)
  and date helpers (`src/lib/datetime.ts`).
- Component/integration tests cover core user scenarios (the guest booking flow
  including the 409 conflict path and all four UI states; a representative host
  page) with MSW returning contract-shaped data.
- Tests run with `TZ=UTC` for deterministic day/time handling.
- E2E (Playwright) and backend tests are **out of scope** for now; E2E is the
  intended home for the DatePicker/DataTable host pages, which are awkward in jsdom.

## Consequences

- Tests are fast, run in CI without a browser or backend, and exercise the real
  API client against the contract (MSW), reinforcing "contract as source of truth".
- Adds dev dependencies: `vitest`, `@vue/test-utils`, `@testing-library/vue`,
  `@testing-library/jest-dom`, `@testing-library/user-event`, `jsdom`, `msw`.
- PrimeVue needs a few jsdom polyfills (see `frontend/src/test/setup.ts`); browser-
  heavy widgets are deferred to a future E2E layer.
- Run with `make test` (or `npm test` in `frontend/`). CI runs in
  `.github/workflows/frontend-ci.yml` (the Hexlet workflow is locked and untouched).
