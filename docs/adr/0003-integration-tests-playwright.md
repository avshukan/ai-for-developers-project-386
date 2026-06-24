# ADR 3: Integration tests with Playwright

## Status

Accepted

## Context

Unit/component tests already exist on both sides in isolation: the frontend uses
Vitest + Testing Library + MSW (see `docs/adr/0001-frontend-testing-strategy.md`)
and the backend uses Go `testing` + `httptest` (see
`docs/adr/0002-backend-stack-and-storage.md`). Both mock their counterpart — MSW
stubs the API for the frontend, and the backend tests call handlers directly.

Nothing yet verifies that the real frontend and the real backend work together.
The course step requires an integration layer that drives the **main booking
scenario** end to end in a real browser, against a real running backend.

`docs/product.md` already defines the user scenarios. The integration layer must
trace to those scenarios rather than redefine them; the primary one is
"Main Guest Scenario" (open booking page → pick event type → pick slot → enter
name/email → confirm → slot becomes unavailable).

## Options

**Tooling**

1. **Playwright (`@playwright/test`, TypeScript)** — drives a real browser,
   first-class TypeScript, built-in test runner, auto-waiting, can boot the app
   under test via its `webServer` option. Recommended by the course step and has
   an MCP integration for agent-driven runs.
2. Cypress — capable e2e tool, but a heavier model, weaker multi-server startup
   story, and TypeScript ergonomics that agents handle less cleanly.
3. Selenium / WebDriver — lower level, more boilerplate, no batteries-included
   runner.

**Where the tests live**

1. **Top-level `e2e/` package** — a separate npm package with its own
   `package.json`. Keeps the Playwright dependency and browser binaries out of
   the shipped frontend bundle and out of the API-contract root package, and
   makes "the tests span both apps" explicit.
2. Inside `frontend/` — reuses the frontend toolchain but couples a full-stack
   concern to the SPA package and pulls Playwright into frontend installs.

**How the app is started for a run**

1. **Playwright `webServer` boots both servers** — the Go backend
   (`go run ./cmd/server`, seeded in-memory data) and the Vite SPA, wired so the
   SPA's `VITE_API_BASE_URL` points at the local backend. One command runs
   everything, locally and in CI.
2. Require the developer/CI to start servers manually — more steps, easy to get
   wrong, no Docker image exists yet anyway (deployment is still `TBD`).

## Decision

- **Tool:** Playwright with `@playwright/test` in TypeScript.
- **Location:** a top-level `e2e/` package with its own `package.json` and
  `playwright.config.ts`.
- **App startup:** Playwright's `webServer` starts both the Go backend and the
  Vite SPA for the run. The backend uses its built-in demo seed
  (`SEED_DATA=true`) so event types and availability exist; the SPA is served
  with `VITE_API_BASE_URL` pointing at the local backend so the two communicate
  over the real HTTP contract. Runs use `TZ=UTC` to match the API's UTC times.
- **Coverage:** the main booking scenario from `docs/product.md` (happy path to
  confirmation) plus the core invariant that a booked slot becomes unavailable
  (no double booking). Tests assert on user-visible behaviour, not implementation
  details, per `AGENTS.md`.
- **CI:** a dedicated `.github/workflows/e2e.yml` workflow installs Go, Node, the
  frontend and e2e dependencies and the Chromium browser, then runs the suite;
  the HTML report is uploaded as an artifact.

## Consequences

- There is now a third test layer that exercises the frontend and backend
  together over the real API, complementing the existing unit/component layers.
- The e2e package adds a dev-only dependency (`@playwright/test`) and downloads a
  browser binary in CI. It ships with nothing — it is not part of the frontend
  bundle or the backend image.
- Because the backend store is in-memory and seeded on startup, a run starts from
  a known state; tests are written to not depend on each other's mutations of
  that shared state within a run.
- Until a Docker topology is defined (still `TBD`), the suite runs against
  `go run` + Vite directly. When Docker/compose lands, `webServer` can point at
  the composed services instead; this ADR's decision (Playwright, top-level
  package, scenario coverage) stays valid.
- Relevant `make` targets (`e2e-install`, `e2e`) are added so agents do not guess
  commands.
