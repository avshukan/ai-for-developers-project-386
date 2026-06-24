# Changelog

## [0.3.0](https://github.com/avshukan/ai-for-developers-project-386/compare/v0.2.1...v0.3.0) (2026-06-24)


### Features

* **backend:** serve built SPA via STATIC_DIR ([8fa8f81](https://github.com/avshukan/ai-for-developers-project-386/commit/8fa8f8146d914adabde427a6320c652f5ce0d993))


### Documentation

* **adr:** add ADR 0005 for combined Docker deploy on Render ([93042d1](https://github.com/avshukan/ai-for-developers-project-386/commit/93042d12294cfdcc2a5de24df0b31981c768e8b3))


### Build System

* **docker:** add multi-stage Dockerfile and .dockerignore ([753686e](https://github.com/avshukan/ai-for-developers-project-386/commit/753686e6e4e6bfbeeacc2a439f1aee096df36f26))
* **make:** add docker-build and docker-run targets ([46a7922](https://github.com/avshukan/ai-for-developers-project-386/commit/46a7922c3a4fd57179e89e1c3f8c8fcb620ed597))

## [0.2.1](https://github.com/avshukan/ai-for-developers-project-386/compare/v0.2.0...v0.2.1) (2026-06-24)


### CI

* bump workflow actions to node24 majors ([b5b7f02](https://github.com/avshukan/ai-for-developers-project-386/commit/b5b7f0281f1e019f182a8ee4e47d6991e4df6d5d))

## [0.2.0](https://github.com/avshukan/ai-for-developers-project-386/compare/v0.1.0...v0.2.0) (2026-06-24)


### Features

* **backend:** implement Go API with in-memory store per the contract ([7655eae](https://github.com/avshukan/ai-for-developers-project-386/commit/7655eae0ade3936a20c6ac7add54c5bd8b67e39a))
* **frontend:** add Vue 3 + Vite SPA for the Call Booking MVP ([ee26c76](https://github.com/avshukan/ai-for-developers-project-386/commit/ee26c766a8244754146258ebee11d9d095bb2893))


### Bug Fixes

* emit OpenAPI to repo-root openapi/ and clarify slots contract ([854a089](https://github.com/avshukan/ai-for-developers-project-386/commit/854a0899db41d7ef3187d6a0d0dc00bec2db5652))


### Documentation

* add architecture.md stub and document frontend tests ([a4b5745](https://github.com/avshukan/ai-for-developers-project-386/commit/a4b57459438fdd15fe4ad21f460740c5942fb866))
* align domain validation rules with the API contract ([36a4324](https://github.com/avshukan/ai-for-developers-project-386/commit/36a4324cd6346cae7cb2aea65e1707e915172834))
* **api:** add request/response examples to the contract ([983e4c3](https://github.com/avshukan/ai-for-developers-project-386/commit/983e4c31ae3d08d75461c23daaff9016e9b50b97))
* **onboarding:** mark backend storage and testing as decided ([e32f1e9](https://github.com/avshukan/ai-for-developers-project-386/commit/e32f1e95bffdfa8f180771028dd7062d09b6f9e1))
* record e2e testing and release automation decisions ([5d4eee4](https://github.com/avshukan/ai-for-developers-project-386/commit/5d4eee465ed1891f3aa5ff3f4319e5a33da4fb5b))


### Tests

* **e2e:** add Playwright integration tests for the booking scenario ([e350a5b](https://github.com/avshukan/ai-for-developers-project-386/commit/e350a5b6c803a987f5ef8430a2b0c65c6082a7b2))
* **frontend:** add Vitest + Testing Library + MSW test suite ([0744e75](https://github.com/avshukan/ai-for-developers-project-386/commit/0744e75e1f19e1ac513c24ac2f6a4972ae31c23a))


### CI

* run integration tests in CI and automate releases with release-please ([c955662](https://github.com/avshukan/ai-for-developers-project-386/commit/c955662d47bfd936cc79dc93bf69092c04270324))
