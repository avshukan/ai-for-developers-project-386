.PHONY: install typespec openapi check \
	frontend-install frontend-dev frontend-build frontend-typecheck frontend-test \
	backend-build backend-run backend-test backend-vet backend-tidy \
	test api-mock

install:
	npm install

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
test: frontend-test backend-test

# --- API mock (Prism) ---
# Serves the generated OpenAPI contract as a mock API on http://localhost:4010.
api-mock:
	npx --yes @stoplight/prism-cli mock openapi/openapi.yaml
