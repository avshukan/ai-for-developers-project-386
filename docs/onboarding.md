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

* A host to publish available 30-minute slots
* A guest to view available slots
* A guest to book a slot with name and email
* A host to view upcoming bookings

The project intentionally avoids advanced features.

Out of scope:

* Authentication
* User accounts
* External calendars
* Payments
* Notifications
* Rescheduling
* Cancellation
* Multiple hosts
* Teams

## Main Workflow

The project follows Design First.

Work order:

```text
Product docs
→ Domain model
→ TypeSpec API contract
→ Generated OpenAPI
→ Frontend and backend implementation
→ Tests
→ Docker
→ Deployment
```

Do not start implementation before the relevant contract and documentation exist.

## Source of Truth

* Domain terminology: `docs/glossary.md`
* Product behavior: `docs/product.md`
* Business rules: `docs/domain.md`
* API contract: `api/main.tsp`
* Generated OpenAPI: generated artifact only
* Agent workflow: `AGENTS.md`

## Current Stack

* Frontend: Vue 3 + Vite
* Backend: Go
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

If a command is missing, add it or mark it as `TBD`.

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

## MVP User Scenarios

The first complete version must support:

1. Guest opens booking page.
2. Guest sees available 30-minute slots.
3. Guest selects a slot.
4. Guest enters name and email.
5. Guest confirms booking.
6. Slot becomes unavailable.
7. Host opens bookings page.
8. Host sees upcoming bookings.

## Notes for Agents

Before changing code:

* Read `AGENTS.md`
* Check the glossary
* Check whether docs or TypeSpec must change first
* Keep changes small
* Use `make`
* Report assumptions and remaining `TBD` items
