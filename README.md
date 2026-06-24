### Hexlet tests and linter status:
[![Actions Status](https://github.com/avshukan/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/avshukan/ai-for-developers-project-386/actions)

# Call Booking

A small meeting scheduling web application inspired by Cal.com.

This is a learning project focused on Design First, API-first development, AI-assisted implementation, Docker delivery, and testing of core user scenarios.

## Project Goal

The goal is to build a complete small web application:

1. Understand the scheduling domain.
2. Define product behavior.
3. Define the API contract.
4. Implement frontend and backend independently.
5. Cover core scenarios with tests.
6. Package and run the application with Docker.
7. Prepare the application for deployment.

## Course Context

The project is developed with AI agents as the main implementation tool.

The human role is to:

* define tasks
* review results
* improve documentation and contracts
* guide the project step by step

The project follows a Design First workflow.

## MVP Scope

The application allows:

* a predefined host to create event types
* a predefined host to publish available time
* a guest to view available event types
* a guest to choose an event type
* a guest to see available slots for the next 14 days
* a guest to book an available slot
* the host to view upcoming bookings

The MVP has no registration, authentication, personal accounts, or external calendar integrations.

## Core Rule

Two bookings cannot exist for the same time, even for different event types.

## Technology Stack

* Frontend: Vue 3 + Vite SPA
* Backend: Go API
* API source of truth: TypeSpec
* Generated API format: OpenAPI
* Runtime/package format: Docker
* Tests: Vitest + Testing Library + MSW (frontend) — see `docs/adr/0001-frontend-testing-strategy.md`; Go standard `testing` + `httptest` (backend) — see `docs/adr/0002-backend-stack-and-storage.md`; Playwright integration tests in `e2e/` — see `docs/adr/0003-integration-tests-playwright.md`
* Releases: release-please + Conventional Commits — see `docs/adr/0004-release-automation-release-please.md`
* Deployment: single combined Docker image (Go API + built SPA on one port) on Render — see `docs/adr/0005-deployment-combined-docker-render.md`

OpenAPI is generated from TypeSpec.
Generated OpenAPI files must not be edited manually.

## Testing

* `make test` — frontend (Vitest) + backend (Go) unit/component tests.
* `make e2e` — Playwright integration tests: the real SPA driven against the
  real backend in a browser, covering the main booking scenario. Run
  `make e2e-install` once first (installs deps + a browser). See `e2e/README.md`.

## Running with Docker

The whole app ships as one image: the Go server serves both the API and the built
SPA on a single port.

```bash
make docker-build              # build the image (tag: call-booking)
make docker-run                # run on http://localhost:8080
```

The container listens on the `PORT` environment variable (default `8080`), so it
runs unchanged on platforms that inject `PORT`. See
`docs/adr/0005-deployment-combined-docker-render.md`.

## Deployment

**Live app: https://call-booking-66ez.onrender.com**

Deployed as a single Docker web service on Render (Railway as fallback); `PORT` is
injected by the platform. Render auto-deploys the `main` branch from the
`Dockerfile`. On the free instance the service spins down when idle, so the first
request after a pause is slow (cold start) and in-memory data resets — demo data
re-seeds automatically on start. See `docs/architecture.md` (Deployment) and
`docs/adr/0005-deployment-combined-docker-render.md`.

## Commits and releases

Commits follow [Conventional Commits](https://www.conventionalcommits.org/)
(`feat:`, `fix:`, `docs:`, …), including agent-authored commits — see the
"Commit Convention" section of `AGENTS.md`. On merge to `main`, release-please
opens/updates a release PR with the next version and `CHANGELOG.md`. See
`docs/adr/0004-release-automation-release-please.md`.

## Documentation

* Agent rules: `AGENTS.md`
* Onboarding: `docs/onboarding.md`
* Glossary: `docs/glossary.md`
* Product requirements: `docs/product.md`
* Domain model: `docs/domain.md`
* Architecture: `docs/architecture.md`
* Architecture decisions: `docs/adr/`
* API contract source: `api/main.tsp`

Some files may be missing at the beginning of the project and can be created when needed.

## Development Process

```text
product behavior
→ domain language
→ domain model
→ TypeSpec API contract
→ generated OpenAPI
→ frontend implementation
→ backend implementation
→ tests
→ Docker
→ deployment
```

Do not start implementation before the relevant documentation and API contract exist.
