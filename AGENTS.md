# AGENTS.md

## Purpose

This file defines how AI agents must work in this repository.

It is an operational guide for agents, not a product specification, glossary, or architecture document.

Agents must optimize for:

* Design First workflow
* small changes
* explicit contracts
* documented decisions
* reproducible commands
* minimal hidden assumptions

## Project Stack

* Frontend: Vue 3 + Vite SPA
* Backend: Go API
* API design: TypeSpec
* Generated API format: OpenAPI
* Runtime/package format: Docker
* Tests: TBD
* Deployment: TBD

Do not replace the stack without an ADR.

## Course Context

This is a learning project.

The course requires:

* studying the Cal.com domain
* defining behavior before implementation
* defining the API contract before frontend/backend implementation
* using TypeSpec as the source of truth for the API
* generating OpenAPI from TypeSpec
* implementing frontend and backend independently
* testing core user scenarios
* building and running the application in Docker

Do not turn this into a production-grade scheduling system unless documentation explicitly changes the scope.

## Sources of Truth

Use these files as primary sources of truth:

* `AGENTS.md` — rules for AI agents
* `README.md` — short project overview
* `docs/onboarding.md` — project map and reading order
* `docs/glossary.md` — domain terminology
* `docs/product.md` — product scope and user scenarios
* `docs/domain.md` — domain model and business rules
* `docs/architecture.md` — technical architecture
* `docs/adr/` — architecture decision records
* `api/main.tsp` — source of truth for the API contract
* generated OpenAPI files — generated artifacts only

Do not duplicate detailed product rules, glossary definitions, or architecture decisions in this file.
Link to the authoritative document instead.

Some files may not exist yet.
Create a missing file only when it is required by the current workflow step.

## Design First Workflow

Follow this order:

1. Understand the task.
2. Read the relevant documentation.
3. Check the glossary.
4. Challenge the idea before implementation.
5. Update documentation if behavior or decisions change.
6. Update TypeSpec if the API changes.
7. Generate OpenAPI from TypeSpec.
8. Implement code.
9. Run checks through `make`.
10. Summarize changes and remaining risks.

Do not implement API behavior before the API contract is defined in TypeSpec.

Do not silently choose behavior for unresolved product or domain `TBD` items.
Critical `TBD` items that affect API behavior must be resolved in documentation before TypeSpec is written.

Every API behavior must trace to a documented user scenario or business rule.

## TypeSpec Rules

`api/main.tsp` is the source of truth for the API.

Rules:

* Change the API through TypeSpec first.
* Regenerate OpenAPI after TypeSpec changes.
* Keep endpoint names, models, and examples aligned with `docs/glossary.md`.
* Do not introduce API fields that are not documented in `docs/domain.md`.
* Do not manually edit generated OpenAPI files.
* If TypeSpec tooling is not configured yet, mark commands as `TBD` and propose the missing setup.

Generated OpenAPI location: `TBD`.

## Documentation Responsibility

Keep responsibilities separated:

* Product scope and user scenarios belong in `docs/product.md`.
* Domain terms belong in `docs/glossary.md`.
* Business rules and invariants belong in `docs/domain.md`.
* Architecture decisions belong in `docs/architecture.md` or `docs/adr/`.
* Agent workflow rules belong in `AGENTS.md`.

Avoid copying the same content across files.
If a rule belongs to another document, reference that document instead.

## Grill Before Build

Before making domain, API, or architecture changes, challenge the proposal.

Check:

* Does it fit the course task?
* Does it fit `docs/product.md`?
* Does it conflict with `docs/glossary.md`?
* Does it require an ADR?
* Does it change the API contract?
* Does it affect frontend/backend independence?
* Can it be tested as a user scenario?

If the idea is risky, document the concern before implementation.

## ADR Rules

Use ADRs for important technical or product decisions.

Create an ADR for significant changes to:

* Stack
* API design approach
* Database choice
* Deployment model
* Testing strategy
* Authentication decision
* Domain boundaries
* Major architecture boundaries

Do not create ADRs for every small domain-model edit.

ADR format:

```text
# ADR N: Title

## Status

Accepted | Proposed | Superseded

## Context

What problem are we solving?

## Options

What alternatives were considered?

## Decision

What did we choose?

## Consequences

What trade-offs follow from this decision?
```

Do not rewrite old ADRs to hide previous decisions.
Create a new ADR if a decision changes.

## Makefile Rules

Agents must use `make` targets instead of guessing commands.

Expected future targets:

```bash
make install
make dev
make test
make lint
make format
make typespec
make openapi
make docker-build
make docker-run
make check
```

If a target does not exist yet:

* Do not invent a different command silently.
* Add the missing target only when it is relevant to the current task.
* Otherwise mark it as `TBD`.

Prefer creating a Makefile early in the project.

## Testing Rules

Tests must focus on behavior, not implementation details.

Test scenarios must come from `docs/product.md` and `docs/domain.md`.

Testing stack is `TBD`.

Expected future test layers:

* backend API tests
* frontend component or integration tests
* end-to-end tests for core user scenarios

## Dependency Rules

Do not add production dependencies casually.

Before adding a dependency, explain:

* why it is needed
* why standard library or existing dependencies are not enough
* whether it affects Docker image size or deployment
* whether it adds maintenance risk

Prefer simple, boring dependencies.

## Generated Files

Generated files must be clearly marked.

Rules:

* Do not manually edit generated files.
* Regenerate them from the source file.
* If generated output changes unexpectedly, explain why.
* Keep generated artifacts reproducible through `make`.

## Frontend Rules

Frontend is a SPA.

Rules:

* Use Vue 3 + Vite.
* Do not use server-side rendering.
* Communicate with the backend only through the documented API.
* Do not hardcode backend behavior that is not in the API contract.
* Keep user scenarios simple and visible.

## Backend Rules

Backend is a Go API service.

Rules:

* Follow the API contract generated from TypeSpec.
* Keep business rules explicit.
* Validate input on the backend.
* Keep storage choice documented.
* Do not introduce authentication unless the product scope changes through ADR.

Database choice is `TBD`.

## Docker Rules

The application must be runnable through Docker.

Rules:

* Provide a Docker-based local run path.
* Keep Docker commands available through `make`.
* Prefer reproducible builds.
* Do not rely on local machine state.

Final deployment target is `TBD`.

## Onboarding Rule

A new agent must be able to understand the project by reading:

1. `AGENTS.md`
2. `README.md`
3. `docs/onboarding.md`
4. `docs/glossary.md`
5. Relevant task-specific docs

If `docs/onboarding.md` is missing, create it when onboarding becomes part of the workflow.

## Change Size

Prefer small changes.

A good agent step should usually change one of these:

* one documentation topic
* one API contract area
* one backend feature
* one frontend feature
* one test scenario
* one infrastructure concern

Avoid large mixed changes.

## Git Safety

Do not overwrite user work.

Rules:

* Do not delete files unless the task requires it.
* Do not rewrite unrelated files.
* Do not perform broad formatting unless requested.
* Do not change generated files without changing their source.
* Do not use destructive git commands.
* Do not commit unless explicitly asked.

## Completion Checklist

Before finishing a task, report:

* what changed
* which files changed
* which commands were run
* which checks passed or failed
* what remains `TBD`
* any risks or assumptions

For implementation tasks, prefer running:

```bash
make check
```

If `make check` does not exist, propose adding it.

## Improving This File

This file should evolve during the project.

If an agent notices:

* a repeated mistake
* a missing command
* an unclear convention
* a useful workflow rule
* a recurring source of confusion

then the agent should propose an update to `AGENTS.md`.

Do not update `AGENTS.md` silently unless the task explicitly asks for it.
