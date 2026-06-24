# ADR 4: Release automation with Conventional Commits and release-please

## Status

Accepted

## Context

The project has CI for checks but no release process: no version history, no
changelog, no tags. The course step asks us to automate releases so that version
bumps and the changelog are derived from commit history rather than maintained by
hand — which fits an agent-driven workflow, where commits are produced
mechanically and should carry enough structure to drive tooling.

Two decisions are coupled:

- **Commit message format** — tooling needs a machine-readable convention to know
  whether a change is a feature, a fix, or breaking.
- **Release tooling** — something that reads that history and proposes a version
  and changelog.

The repository contains several units (TypeSpec contract at the root, `frontend/`
SPA, `backend/` Go service), but it is a single learning application released as
one thing. There is no independent distribution of the parts.

## Options

**Commit convention**

1. **Conventional Commits** (`feat:`, `fix:`, `chore:`, `docs:`, …, `feat!:` /
   `BREAKING CHANGE:` for breaking) — the de-facto standard, directly consumable
   by release-please and SemVer-mapping.
2. Free-form messages — no tooling can derive versions; rejected.

**Release tooling**

1. **release-please** (`googleapis/release-please-action`) — reads Conventional
   Commits on the default branch and maintains a long-lived "release PR" with the
   computed next version and generated `CHANGELOG.md`; merging that PR tags and
   creates the GitHub release. No local release step, no manual version edits.
2. semantic-release — releases immediately on merge to the default branch. Powerful
   but publishes without a review gate, and is oriented toward publishing npm
   packages, which this repo does not do.
3. Manual tagging + hand-written changelog — defeats the purpose of the step.

**release-please topology**

1. **Single root release (`release-type: simple`, manifest mode)** — one version
   and one `CHANGELOG.md` for the whole application; version tracked in
   `.release-please-manifest.json`, language-agnostic so it does not have to bump
   any one of the three `package.json`/`go.mod` units.
2. Monorepo mode with a component per package (frontend/backend/contract) —
   separate versions and changelogs. More moving parts than a single learning app
   needs and implies the parts are released independently, which they are not.

## Decision

- **Commit convention:** Conventional Commits for all commits, including those
  authored by AI agents. The rule and the allowed types are documented in
  `AGENTS.md` ("Commit Convention") so agents follow it by default.
- **Tooling:** release-please via `googleapis/release-please-action@v4`, run on
  every push to `main` in `.github/workflows/release-please.yml`.
- **Topology:** a single root release using `release-type: simple` in manifest
  mode (`release-please-config.json` + `.release-please-manifest.json`),
  producing one `CHANGELOG.md` and one set of tags/GitHub releases for the whole
  application. Bootstrap version is `0.1.0` (matching the current package
  versions; the repo has no tags yet).
- **Flow:** merges to `main` make release-please open or update a release PR with
  the next SemVer version and changelog; merging that PR creates the tag and
  GitHub release. No one edits versions or the changelog by hand.

## Consequences

- Versioning and changelog generation become a by-product of well-formed commits;
  the discipline shifts to writing good Conventional Commit messages (now an
  explicit agent rule).
- SemVer mapping: `fix:` → patch, `feat:` → minor, `feat!:`/`BREAKING CHANGE:` →
  major (still `0.x`, so breaking changes bump the minor until a `1.0.0` is cut
  deliberately). `docs:`/`chore:`/`refactor:`/`test:`/`ci:` do not bump the
  version but can appear in the changelog.
- The workflow uses the default `GITHUB_TOKEN`. For it to open the release PR, the
  repository setting **Settings → Actions → General → Workflow permissions →
  "Allow GitHub Actions to create and approve pull requests"** must be enabled. If
  org policy forbids that, a Personal Access Token can be supplied as the action's
  `token:` instead. This is the one manual prerequisite.
- Choosing a single root release means the three sub-packages are not versioned
  independently. If a part ever needs its own release cadence, switching to
  release-please monorepo mode would be a new ADR.
- The auto-generated `hexlet-check.yml` workflow is unaffected; release-please is
  an additional, independent workflow.
