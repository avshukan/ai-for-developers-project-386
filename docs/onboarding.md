# Onboarding

This document helps a new AI agent or developer understand the project quickly.

## Reading Order

Start with these files:

1. `AGENTS.md`
2. `README.md`
3. `docs/glossary.md`
4. `docs/product.md`
5. `docs/domain.md`
6. `docs/architecture.md`
7. `api/main.tsp`

Some files may not exist yet. Missing files are expected at the beginning of the project.

## Project Summary

This is a small meeting scheduling application inspired by Cal.com.

The application allows:

* a predefined host to create event types
* a predefined host to publish available time
* a guest to view available event types
* a guest to choose an event type
* a guest to see available slots for the next 14 days
* a guest to book an available slot with name and email
* a host to view upcoming bookings

The project intentionally avoids advanced production features.

Out of scope:

* registration
* authentication
* personal accounts
* external calendars
* payments
* notifications
* cancellation
* rescheduling
* multiple hosts
* teams
* recurring events

## Course Workflow

The project follows Design First.

Work order:

```text
study Cal.com domain
→ define product behavior
→ define shared terminology
→ define domain model
→ create TypeSpec API contract
→ generate OpenAPI
→ implement frontend and backend independently
→ test core scenarios
→ build and run with Docker
→ prepare deployment
```

Do not start implementation before the relevant contract and documentation exist.

## Source of Truth

* Agent workflow: `AGENTS.md`
* Project overview: `README.md`
* Domain terminology: `docs/glossary.md`
* Product behavior: `docs/product.md`
* Business rules: `docs/domain.md`
* API contract source: `api/main.tsp`
* Generated OpenAPI: generated artifact only

OpenAPI must be generated from TypeSpec.
Generated OpenAPI files must not be edited manually.

## Current Stack

* Frontend: Vue 3 + Vite SPA
* Backend: Go API
* API design: TypeSpec
* API output: OpenAPI
* Runtime: Docker
* Tests: TBD
* Database: TBD
* Deployment: TBD

## Expected Commands

Use `make` commands only.

Expected future commands:

```bash
make install
make dev
make test
make lint
make format
make typespec
make openapi
make docker-build
make docker-run
make check
```

If a command is missing, add it only when it is relevant to the current workflow step. Otherwise mark it as `TBD`.

## First Project Milestones

1. Create base documentation.
2. Define glossary.
3. Define product scope.
4. Define domain model.
5. Create TypeSpec API contract.
6. Generate OpenAPI.
7. Scaffold frontend and backend.
8. Implement core booking flow.
9. Add tests.
10. Build Docker image.
11. Prepare deployment.

## MVP User Scenarios

The first complete version must support:

1. Host creates an event type.
2. Host publishes available time.
3. Guest opens the booking page.
4. Guest sees available event types.
5. Guest selects an event type.
6. Guest sees available slots for the next 14 days.
7. Guest selects an available slot.
8. Guest enters name and email.
9. Guest confirms booking.
10. Selected time becomes unavailable.
11. Host opens bookings page.
12. Host sees upcoming bookings across all event types.

## Core Constraint

Two bookings cannot exist for the same time, even for different event types.

This rule must be enforced by the backend.

## Notes for Agents

Before changing code:

* read `AGENTS.md`
* check the glossary
* check product and domain docs
* check whether TypeSpec must change first
* keep changes small
* use `make`
* report assumptions and remaining `TBD` items
