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
* Tests: Vitest + Testing Library + MSW (frontend) — see `docs/adr/0001-frontend-testing-strategy.md`; Go standard `testing` + `httptest` (backend) — see `docs/adr/0002-backend-stack-and-storage.md`
* Deployment: TBD

OpenAPI is generated from TypeSpec.
Generated OpenAPI files must not be edited manually.

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
