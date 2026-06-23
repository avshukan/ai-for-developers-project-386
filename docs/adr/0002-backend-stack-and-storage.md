# ADR 2: Backend stack and storage

## Status

Accepted

## Context

The backend implementation step requires a Go API service (see `AGENTS.md` and
`README.md`) that serves the contract in `api/main.tsp` and enforces the core
booking rules from `docs/domain.md`. Two items were left `TBD` and must be
resolved before implementation:

- **Database choice** — `docs/domain.md` and `docs/product.md` list it as `TBD`,
  but the course step explicitly allows in-memory storage that may reset on
  restart.
- **Backend testing stack** — `AGENTS.md` lists it as `TBD`.

`AGENTS.md` also asks us to prefer the standard library over new dependencies and
to record significant decisions (database, testing strategy) as ADRs.

## Options

**HTTP layer**

1. **Go standard library `net/http`** — `http.ServeMux` (Go 1.22+) supports
   method + path patterns (`GET /event-types/{eventTypeId}/slots`), which covers
   every route in the contract. Zero dependencies.
2. A third-party router/framework (chi, gin, echo) — more features, but extra
   dependencies for routing we do not need at this size.

**Storage**

1. **In-memory store guarded by a mutex** — matches the course step; the
   no-double-booking invariant is enforced by checking and inserting under one
   lock.
2. SQLite / Postgres — durable, supports a DB-level uniqueness constraint, but
   adds a dependency, migrations and a running service the MVP does not need.

**Testing**

1. **Go standard `testing` + `net/http/httptest`** — table tests for the pure
   slot logic and store, recorder-based tests for the HTTP handlers. No
   dependencies.
2. A third-party assertion/suite framework (testify, ginkgo) — nicer assertions,
   but an avoidable dependency.

## Decision

- **HTTP:** Go standard library `net/http` with `http.ServeMux` pattern routing.
  No web framework.
- **Storage:** an in-memory store (`internal/store`) guarded by a `sync.Mutex`.
  Data resets on restart, which the course step permits. The store owns the
  no-double-booking invariant: the event-type existence check, slot-availability
  check, conflict check and insert all run under one lock, so concurrent requests
  for the same time cannot both succeed.
- **Testing:** the Go standard `testing` package plus `net/http/httptest`. Slot
  generation and the store are unit-tested (including a concurrent
  double-booking test run under `-race`); the HTTP layer is tested through
  `httpapi.NewRouter` against the contract's status codes and error shapes.

## Consequences

- The backend has **no production or test dependencies** beyond the Go standard
  library, keeping the Docker image small and the build reproducible.
- In-memory storage is not durable and is single-process. If the app ever needs
  persistence or horizontal scaling, the atomic double-booking guarantee must be
  re-established at the database level (e.g. a uniqueness constraint or a
  transaction); that would be a new ADR.
- Backend layout: `cmd/server` (entrypoint, config, seed), `internal/domain`
  (entities), `internal/slots` (pure slot rules), `internal/store` (in-memory
  state + invariants), `internal/httpapi` (routing, validation, error mapping,
  CORS).
- CORS is permissive (single configurable allowed origin, default `*`) because
  host pages are public in the MVP and there are no credentials to protect.
- Resolves the `TBD` for database choice and backend testing stack for the MVP.
