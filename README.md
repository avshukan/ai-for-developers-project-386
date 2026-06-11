### Hexlet tests and linter status:
[![Actions Status](https://github.com/avshukan/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/avshukan/ai-for-developers-project-386/actions)

# Call Booking

A simple meeting scheduling application inspired by Cal.com.

The project is built as part of a Design First workflow. The API contract is defined before implementation and serves as the source of truth for independent frontend and backend development.

## Goals

* Learn AI-assisted software development
* Practice Design First and API-first workflow
* Build a complete web application from specification to deployment
* Keep frontend and backend development independent through an explicit API contract

## MVP Scope

The application allows:

* A host to publish available time ranges for meetings
* A guest to view available 30-minute slots
* A guest to book a slot by providing name and email
* A host to view upcoming bookings

The project intentionally excludes authentication, user accounts, external calendar integrations, payments, notifications, cancellation, rescheduling, teams, and recurring events.

Host pages are public in the MVP. This is an intentional learning simplification, not a production security model.

## Technology Stack

* Frontend: Vue 3 + Vite SPA
* Backend: Go API
* API source of truth: TypeSpec
* Generated API format: OpenAPI
* Runtime/package format: Docker
* Tests: TBD
* Deployment: TBD

Generated OpenAPI files must not be edited manually. API changes must be made through TypeSpec first.

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

1. Define product behavior.
2. Define domain language and domain model.
3. Design the API contract in TypeSpec.
4. Generate OpenAPI from TypeSpec.
5. Implement frontend and backend independently.
6. Cover key scenarios with tests.
7. Package and run the application with Docker.
8. Deploy the application.

Do not start implementation before the relevant documentation and API contract exist.
