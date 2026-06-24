# ADR 5: Deployment as a single combined Docker image on Render

## Status

Accepted

## Context

The course step "Docker и деплой" requires:

- a `Dockerfile` in the repository that the checker builds,
- the application starting automatically when the container runs,
- the application listening on the port given by the `PORT` environment variable,
- a single public URL after deployment.

Until now the app has been two independently served parts (ADR 0002,
`docs/architecture.md`): a Go API (`backend/`, already reads `PORT`) and a Vue 3 +
Vite SPA (`frontend/`) served separately, talking to the API through
`VITE_API_BASE_URL`. The deployment target was `TBD`.

The checker's contract (one Dockerfile, one `PORT`, one public link) does not map
cleanly onto two separately deployed services, so we need to decide the runtime
topology and the hosting platform.

## Options

**Runtime topology**

1. **Combined single image** — build the SPA to static files and let the Go
   server serve both the SPA and the API on one `PORT`. One image, one URL, exact
   fit for the checker. Requires the backend to optionally serve a static
   directory and the SPA to use same-origin (relative) API paths.
2. **Two separate services** — deploy backend and frontend as two services with
   two URLs. Matches the existing "deployed independently" wording, but means two
   Dockerfiles / two public links, the SPA must bake the backend URL in at build
   time, and CORS must be pinned to the frontend origin. A poor fit for a
   single-link checker.

**Hosting platform**

1. **Render** — free Docker web service, injects `PORT`, gives a public
   `*.onrender.com` URL, deploys from the GitHub repo, automatable via the Render
   MCP server.
2. **Railway** — equivalent Docker/`PORT`/public-URL model; kept as a fallback if
   Render requires payment or is unavailable.
3. Other PaaS / self-hosting — more setup than this learning step needs.

## Decision

- **Topology:** a single combined Docker image. A multi-stage `Dockerfile` builds
  the SPA (`node:22-alpine`), builds the Go server (`golang:1.26-alpine`), and
  ships a small runtime image that runs the server and serves the built SPA from
  `STATIC_DIR`.
- **Backend change:** `httpapi.NewRouter` takes a `staticDir`; when set, a
  catch-all route serves the SPA with client-side-routing fallback
  (`backend/internal/web`). The specific API patterns keep precedence, so API
  routing is unchanged. When `staticDir` is empty the backend stays API-only
  (local dev, unit tests, and the Playwright e2e suite are unaffected).
- **SPA change:** the production build sets `VITE_API_BASE_URL=""` (a Docker build
  arg) so the SPA calls the API with same-origin relative paths.
- **Platform:** deploy the image to Render as a Docker web service on the free
  instance type; `PORT` is injected by Render. Railway is the documented fallback.

## Consequences

- One image, one container, one public URL — matches the course requirement and
  keeps the local run path (`make docker-build` / `make docker-run`) identical to
  production.
- Same-origin serving means the SPA no longer needs cross-origin requests; the
  default `CORS_ALLOWED_ORIGIN=*` remains harmless and no origin pinning is
  required.
- "Independent" now means **build-time** independence (the parts are still built
  and tested separately and share only the API contract); they are **combined at
  runtime** in the shipped image. ADR 0002's two-part model and ADR 0003's
  split-process e2e setup remain valid.
- Storage is still in-memory (ADR 0002). On Render's free tier the service spins
  down when idle and resets data on restart, but `SEED_DATA=true` re-seeds on
  every start, so the app stays demoable. Durable storage remains future work.
- If the parts ever need independent deployment or scaling, splitting back into
  two services (and re-pinning CORS) would be a new ADR.
