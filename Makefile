.PHONY: install dev lint format typespec openapi check \
	frontend-install frontend-dev frontend-build frontend-typecheck frontend-test \
	backend-build backend-run backend-test backend-vet backend-tidy \
	test e2e-install e2e e2e-report api-mock \
	docker-build docker-run

# Image tag for the combined app image (Go API + built SPA on one port).
IMAGE ?= call-booking

# Install all dependencies: API contract (root), frontend, and Go modules.
install:
	npm install
	cd frontend && npm install
	cd backend && go mod download

# Run the whole app locally: Go API on :8080 and the SPA dev server on :5173,
# with the SPA pointed at the local backend. Ctrl-C stops both.
dev:
	@trap 'kill 0' INT TERM; \
		(cd backend && go run ./cmd/server) & \
		(cd frontend && VITE_API_BASE_URL=http://localhost:8080 npm run dev) & \
		wait

# Static checks: Go vet + frontend type-check.
lint:
	cd backend && go vet ./...
	cd frontend && npm run typecheck

# Format Go sources. (The frontend has no formatter configured yet.)
format:
	cd backend && gofmt -w .

typespec:
	npm run typespec

openapi:
	npm run openapi

# Contract + backend static checks.
check:
	npm run check
	cd backend && go vet ./...

# --- Frontend (Vue 3 + Vite SPA) ---
# Requires Node >= 20 (see .nvmrc). Run `nvm use` first if your default differs.
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-typecheck:
	cd frontend && npm run typecheck

frontend-test:
	cd frontend && npm test

# --- Backend (Go API) ---
# Requires Go (see backend/go.mod). Storage is in memory and resets on restart.
backend-build:
	cd backend && go build -o bin/server ./cmd/server

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test -race ./...

backend-vet:
	cd backend && go vet ./...

backend-tidy:
	cd backend && go mod tidy

# Aggregate test target (AGENTS.md expects `make test`).
# Unit/component layers only; e2e is separate (needs a browser + both servers).
test: frontend-test backend-test

# --- Integration tests (Playwright; real SPA + real backend) ---
# See docs/adr/0003-integration-tests-playwright.md. Requires Node and Go.
# Playwright boots the backend and frontend itself; no manual server start needed.
e2e-install:
	cd e2e && npm install && npm run install:browsers

e2e:
	cd e2e && npm test

e2e-report:
	cd e2e && npm run report

# --- API mock (Prism) ---
# Serves the generated OpenAPI contract as a mock API on http://localhost:4010.
api-mock:
	npx --yes @stoplight/prism-cli mock openapi/openapi.yaml

# --- Docker (single combined image: Go API + built SPA on one port) ---
# See docs/adr/0005-deployment-combined-docker-render.md.
docker-build:
	docker build -t $(IMAGE) .

# Run locally on :8080. The platform (Render) injects PORT in production.
docker-run:
	docker run --rm -e PORT=8080 -p 8080:8080 $(IMAGE)
