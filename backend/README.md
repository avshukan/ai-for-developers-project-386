# Backend — Call Booking API

Go API service for the Call Booking MVP. It implements the contract in
`../api/main.tsp` (generated to `../openapi/openapi.yaml`) and enforces the
booking business rules from `../docs/domain.md`.

Built on the Go standard library `net/http` — no web framework, no third-party
dependencies. Storage is in memory and resets on restart. See
`../docs/adr/0002-backend-stack-and-storage.md` for the stack and storage
decision.

## Layout

```
cmd/server        entrypoint, configuration, demo seed data
internal/domain   entities: EventType, Availability, Slot, Booking
internal/slots    pure slot-generation rules (30-min, 14-day window, past/booked)
internal/store    in-memory state + the no-double-booking invariant
internal/httpapi  routing, input validation, contract error mapping, CORS
internal/web      serves the built SPA when STATIC_DIR is set (combined image)
```

## Run

From the repository root:

```bash
make backend-run        # go run ./cmd/server
make backend-build      # builds bin/server
make backend-test       # go test -race ./...
make backend-vet        # go vet ./...
```

Or directly from `backend/`:

```bash
go run ./cmd/server
```

The server listens on `http://localhost:8080` by default and seeds a couple of
event types plus two weeks of availability so the app is usable immediately.

### Configuration (environment variables)

| Variable              | Default | Purpose                                              |
| --------------------- | ------- | ---------------------------------------------------- |
| `PORT`                | `8080`  | TCP port to listen on                                |
| `CORS_ALLOWED_ORIGIN` | `*`     | `Access-Control-Allow-Origin` value                  |
| `SEED_DATA`           | `true`  | Seed demo data on startup; set to `false` to disable |
| `STATIC_DIR`          | (empty) | Directory with the built SPA to serve from the same origin; empty = API-only |

In the combined Docker image the server also serves the built SPA from
`STATIC_DIR`, so the SPA and API share one origin and CORS is not exercised.
`CORS_ALLOWED_ORIGIN` only matters when the SPA is served from a different origin
(local dev against the Prism mock, or a split deploy). See
`../docs/adr/0005-deployment-combined-docker-render.md`.

## Connecting the frontend

The SPA reads its base URL from `VITE_API_BASE_URL`. To point it at this backend
instead of the Prism mock, set in `frontend/.env`:

```
VITE_API_BASE_URL=http://localhost:8080
```

## Endpoints

| Method | Path                              | Purpose                          |
| ------ | --------------------------------- | -------------------------------- |
| GET    | `/event-types`                    | List event types (guest)         |
| POST   | `/host/event-types`               | Create an event type (host)      |
| POST   | `/host/availability`              | Publish an availability range    |
| GET    | `/event-types/{eventTypeId}/slots`| Available slots for next 14 days |
| POST   | `/bookings`                       | Create a booking (guest)         |
| GET    | `/host/bookings`                  | List upcoming bookings (host)    |

Errors follow the contract: `400 validation_error`, `404 not_found`,
`409 booking_conflict`.
