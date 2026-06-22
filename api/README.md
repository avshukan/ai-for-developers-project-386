# API Contract

This directory contains the TypeSpec API contract for the Call Booking project.

The API contract is the source of truth for frontend and backend implementation.

## Files

* `main.tsp` — TypeSpec API contract
* `tspconfig.yaml` — TypeSpec compiler and OpenAPI emitter configuration

## Generated Output

OpenAPI is generated from TypeSpec into:

```text
openapi/openapi.yaml
```

Do not edit generated OpenAPI manually.

Change `api/main.tsp` first, then regenerate OpenAPI.

## Commands

From the repository root:

```bash
make install
make typespec
make openapi
make check
```

## API Scope

The contract covers the MVP scenarios:

* guest views event types
* host creates event types
* host creates availability ranges
* guest views available slots
* guest creates a booking
* host views upcoming bookings

## Rules

* The API follows `docs/product.md` and `docs/domain.md`.
* Domain names must match `docs/glossary.md`.
* Date-time values use UTC ISO 8601.
* Slot duration is fixed to 30 minutes in the MVP.
* Double booking must be rejected by the backend.
* Booking conflict must return `409 Conflict`.
