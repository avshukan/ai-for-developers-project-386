# Architecture

> Stub. This document will describe the technical architecture of the Call
> Booking application. Sections marked `TBD` are intentionally unresolved at this
> stage of the project and will be filled in as decisions are made.

## Overview

Call Booking is a small meeting scheduling application inspired by Cal.com. It is
split into two independently implemented parts that communicate only through the
documented API contract:

- **Frontend** — Vue 3 + Vite SPA (`frontend/`).
- **Backend** — Go API service (`backend/`).

The API contract is the source of truth. It is authored in TypeSpec
(`api/main.tsp`) and generated into OpenAPI (`openapi/openapi.yaml`). Neither side
imports the other's code.

## Components

### Frontend

Vue 3 + Vite SPA. Talks to the backend only through the public API contract via
`src/api/`. See `frontend/README.md` for details.

### Backend

Go API service (`backend/`) built on the standard library `net/http` — no web
framework. It implements the contract generated from TypeSpec and enforces the
core business rules (including the no-double-booking invariant). Layout:

- `cmd/server` — entrypoint, configuration (env vars) and demo seed data.
- `internal/domain` — entities (`EventType`, `Availability`, `Slot`, `Booking`).
- `internal/slots` — pure slot-generation rules (30-minute slots, 14-day window,
  past/booked exclusion).
- `internal/store` — in-memory state and the booking invariants.
- `internal/httpapi` — routing, input validation, contract error mapping, CORS.

See `backend/README.md` for run instructions and
`docs/adr/0002-backend-stack-and-storage.md` for the stack/storage decision.

### Storage

In-memory store guarded by a mutex (`backend/internal/store`). Data resets on
restart, which the MVP permits. The store owns the no-double-booking invariant:
the existence, slot-availability, conflict and insert steps run under one lock.
See `docs/adr/0002-backend-stack-and-storage.md`.

## Boundaries

- Frontend and backend are implemented and deployed independently.
- The only coupling between them is the API contract (`openapi/openapi.yaml`,
  generated from `api/main.tsp`).
- Business rules and input validation live in the backend.

## Cross-cutting decisions

- Date-time values use UTC ISO 8601.
- Slot duration is fixed to 30 minutes in the MVP.
- Architecture decisions are recorded as ADRs in `docs/adr/`.

## Deployment

`TBD` — the application must be runnable through Docker. Final deployment target
is `TBD`.

## Related documents

- `docs/product.md` — product scope and user scenarios
- `docs/domain.md` — domain model and business rules
- `docs/glossary.md` — domain terminology
- `api/main.tsp` — API contract source of truth
- `docs/adr/` — architecture decision records
