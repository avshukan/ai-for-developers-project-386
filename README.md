### Hexlet tests and linter status:
[![Actions Status](https://github.com/avshukan/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/avshukan/ai-for-developers-project-386/actions)

# Call Booking

A simple meeting scheduling application inspired by Cal.com.

The project is built as part of a Design First workflow. The API contract is defined before implementation and serves as the source of truth for both frontend and backend development.

## Goals

* Learn AI-assisted software development
* Practice Design First approach
* Build a complete web application from specification to deployment
* Keep frontend and backend development independent through an explicit API contract

## Scope

The application allows:

* Publishing available 30-minute meeting slots
* Viewing available slots
* Booking a slot by providing name and email
* Viewing upcoming bookings

The project intentionally excludes authentication, external calendar integrations, payments, notifications, and other advanced scheduling features.

## Technology Stack

* Frontend: Vue 3 + Vite
* Backend: Go
* API Contract: OpenAPI
* Testing: TBD
* Deployment: Docker

## Documentation

* Product requirements: `docs/01-product.md`
* Domain model: `docs/02-domain.md`
* API contract: `docs/03-api-contract.md`
* Architecture: `docs/04-architecture.md`
* Iteration plan: `docs/05-iteration-plan.md`

## Development Process

1. Define requirements.
2. Model the domain.
3. Design the API contract.
4. Implement frontend and backend independently.
5. Cover key scenarios with tests.
6. Package the application in Docker.
7. Deploy the application.
