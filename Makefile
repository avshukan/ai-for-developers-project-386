.PHONY: install typespec openapi check frontend-install frontend-dev frontend-build frontend-typecheck api-mock

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

# --- API mock (Prism) ---
# Serves the generated OpenAPI contract as a mock API on http://localhost:4010.
api-mock:
	npx --yes @stoplight/prism-cli mock openapi/openapi.yaml
