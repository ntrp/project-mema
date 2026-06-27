# PRD: Engineering Architecture and Setup Standards

Status: Draft

## Summary

When implementation begins, the application should be set up with current, conservative Go and SvelteKit practices. The architecture should support a self-hosted application with a Go backend, SvelteKit frontend, background workers, media-processing jobs, and durable configuration.

`Mema` is only the current codename. Do not bake it into user-facing branding, package names, binary names, API names, Docker image names, or generated identifiers until the final project name is chosen.

This PRD intentionally documents engineering requirements without creating code yet.

## Current Guidance Baseline

- Go setup should follow official Go module layout guidance and avoid unnecessary framework or directory conventions.
- New Go code should target the current stable Go release line at setup time. As of 2026-06-27, the official Go release history lists Go 1.26.4 as the latest patch release.
- New SvelteKit code should use Svelte 5 and SvelteKit's current `sv` CLI workflow.
- New Svelte components should use runes-mode patterns.
- Tooling should be installed through official SvelteKit add-ons where available.

References checked while drafting:

- Official Go release history and module layout documentation.
- Official Svelte/SvelteKit MCP documentation for project creation, project structure, project types, server-only modules, TypeScript, accessibility, performance, `sv create`, `sv add`, `sv check`, ESLint, Prettier, Vitest, and Playwright.

## Backend Requirements

- Use Go modules.
- Prefer a simple, explicit Go project layout over cargo-culted `pkg/` and deep framework structures.
- Keep application entry points under `cmd/`.
- Keep internal application packages under `internal/`.
- Keep domain logic independent from HTTP handlers and persistence adapters.
- Use `context.Context` for request, worker, downloader, and external-process cancellation.
- Use structured logging.
- Use explicit configuration loading and validation at startup.
- Treat all filesystem paths as untrusted until resolved under configured roots.
- Invoke external tools like `ffmpeg`, `ffprobe`, and `mkvmerge` through argument arrays, not shell command strings.
- Redact secrets from logs and error output.
- Expose health and readiness endpoints.
- Write table-driven tests for parsers, scoring, import decisions, and profile evaluation.
- Use integration tests for database, download-client adapters, indexer adapters, and media-tool wrappers.
- Run `go test ./...`, `go vet`, formatting, and vulnerability checks in CI once code exists.

## Frontend Requirements

- Use SvelteKit with TypeScript.
- Use `npx sv create` for initial setup unless there is a strong reason not to.
- Add official SvelteKit tooling through `sv add` where available: ESLint, Prettier, Vitest, Playwright, and adapter configuration.
- Use Svelte 5 runes for new component state: `$state`, `$derived`, `$props`, and `$effect` only where side effects are genuinely needed.
- Avoid legacy Svelte patterns for new code: `export let`, `$:` reactive statements, `on:` event syntax, slot APIs, and legacy stores unless maintaining old code.
- Keep secrets and server-only frontend code in SvelteKit server-only modules when SvelteKit server code is used.
- Prefer progressive enhancement for forms that benefit from server validation.
- Keep accessibility checks in the normal quality gate through `sv check` and ESLint.
- Use Playwright for critical UI flows: adding media, editing profiles, interactive search, import review, and assembly review.
- Use Vitest for isolated component and utility tests.
- Design route data loading so large queues, histories, and search results can paginate or stream instead of blocking the UI.

## Architecture Direction

Mema should start as a modular monolith:

- Go backend owns the API, persistence, workers, media processing, indexer adapters, and download-client adapters.
- PostgreSQL is the first database.
- SvelteKit frontend owns the management UI.
- During development, run the Go API and SvelteKit dev server as separate processes.
- For production, the Go application should serve the built frontend assets so deployment is one app container plus PostgreSQL.
- Do not run a SvelteKit server in production.
- SvelteKit browser code calls the Go API directly using the generated OpenAPI client.
- Authentication is mandatory from version 0 using a local admin user and secure session cookie. OAuth/OIDC is deferred.
- PostgreSQL schema initialization should be resettable during early development because migration tooling is deferred until the project is live.
- Server-Sent Events are the first realtime mechanism for queues, jobs, search progress, and assembly progress.
- Use chi for routing, pgx for PostgreSQL access, and sqlc for generated SQL-first data access.
- Use River for background job management.
- Use OpenAPI from the start with a contract-first workflow. API behavior changes update the spec first, then generated Go/TypeScript artifacts, then backend and web implementation.
- Use `oapi-codegen` for Go OpenAPI generation and `openapi-typescript`/`openapi-fetch` for the web client.
- Use a top-level Go module with SvelteKit under `web/`.
- Use a Makefile for local commands.
- Bundle required media tools into the production Docker image and use locally installed tools for local runs.

These decisions affect authentication, API routing, CSRF, deployment, and Docker layout.

## Backend Package Candidates

Initial package boundaries should be boring and domain-driven:

- `cmd/<project-name>`
- `internal/config`
- `internal/httpapi`
- `internal/auth`
- `internal/library`
- `internal/metadata`
- `internal/indexers`
- `internal/downloads`
- `internal/imports`
- `internal/profiles`
- `internal/scoring`
- `internal/media`
- `internal/assembly`
- `internal/subtitles`
- `internal/jobs`
- `internal/jobs/river`
- `internal/storage`
- `internal/notifications`

These are candidates, not a final commitment.

## Frontend Area Candidates

Initial SvelteKit route areas:

- Dashboard
- Libraries
- Movies
- TV
- Anime
- Books
- Music
- Search
- Queue
- History
- Profiles
- Indexers
- Download clients
- Assembly jobs
- Settings
- System health

## Quality Gates

Before a code change is considered done:

- Go code is formatted with `gofmt`.
- Go tests pass for touched backend packages.
- Svelte code passes `sv check`.
- Frontend lint and formatting checks pass.
- Critical UI flows touched by a change have Playwright coverage or a documented reason for deferral.
- Media-processing changes include representative fixture tests where legally and practically possible.
- Profile/scoring changes include exact score-breakdown tests.

## External Tool Requirements

The backend should treat media tools as versioned capabilities:

- Detect installed tool versions.
- Show missing capabilities in system health.
- Store tool paths in configuration.
- Refuse jobs that require unavailable tools.
- Capture command arguments, exit codes, stderr summaries, and output artifact metadata.
- Avoid logging full paths or filenames if configured privacy mode is enabled.

## Open Questions

- Should background jobs be stored in the main database or a separate queue backend?
- How should the resettable development schema system be invoked and protected?
