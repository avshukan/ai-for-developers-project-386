.PHONY: install typespec openapi check frontend-install frontend-dev frontend-build frontend-typecheck frontend-test test api-mock

install:
	npm install

typespec:
	npm run typespec

openapi:
	npm run openapi

check:
	npm run check

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

# Aggregate test target (AGENTS.md expects `make test`).
test: frontend-test

# --- API mock (Prism) ---
# Serves the generated OpenAPI contract as a mock API on http://localhost:4010.
api-mock:
	npx --yes @stoplight/prism-cli mock openapi/openapi.yaml
