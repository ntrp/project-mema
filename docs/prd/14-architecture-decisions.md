# Architecture Decisions

Status: Draft

Codename `Mema` remains temporary until the final project name is chosen.

## Summary

This document records accepted architecture decisions so scaffolding can proceed without re-litigating foundational choices.

## Accepted Decisions

### Frontend And Backend Runtime

- During development, run the Go API and SvelteKit dev server as separate processes.
- In production, the Go application serves the built frontend assets.
- Production deployment starts as one app container plus PostgreSQL.
- Do not run a SvelteKit server in production.
- The SvelteKit app is built as static web assets and served by Go.
- Browser code calls the Go API directly using the generated OpenAPI client.
- SvelteKit server routes and server actions are not part of the normal architecture.

Rationale:

- Separate dev servers keep SvelteKit fast and idiomatic during UI work.
- A single app container in production keeps self-hosted deployment simple.
- Go owns API, auth, workers, filesystem access, and media-processing jobs.

### Authentication

- Authentication is mandatory from version 0.
- Start with one local admin account and secure session cookies.
- OAuth, OIDC, LDAP, SSO, and multi-user roles are deferred.

Rationale:

- The app controls downloads, filesystem writes, credentials, and media-processing jobs, so unauthenticated operation is too risky even for early versions.
- Local admin auth is enough for a self-hosted first version.

### Database

- PostgreSQL is the first database.
- Do not add migration tooling until the project is live.
- During early development, provide a simple resettable app/schema initialization system.

Rationale:

- PostgreSQL avoids designing around SQLite limitations for jobs, locking, JSON, indexing, and future scale.
- Before real users and live data, resettable schema initialization is simpler than maintaining migration history.

### Development Reset System

The scaffold should include a protected development-only reset flow that can drop and recreate app-owned database objects and seed minimal data.

Requirements:

- Must be disabled by default in production.
- Must require an explicit development environment flag.
- Must clearly log destructive reset actions.
- Must reset only app-owned schema/data.
- Must be callable from a developer command, not from unauthenticated HTTP.
- Should support reseeding the local admin account.
- Should run River's required job-table setup as part of local bootstrap/reset.

Implementation decision:

- Add a Go CLI subcommand shaped as `<binary> reset-dev`.
- Gate it with `APP_ENV=development` and `ALLOW_DEV_RESET=true`.
- Expose it through `make db-reset`.

### Realtime Updates

- Use Server-Sent Events first.
- Use SSE for queues, jobs, search progress, import progress, and assembly progress.
- WebSockets are deferred until bidirectional realtime interaction is needed.

Rationale:

- Most realtime UI needs are server-to-client status streams.
- SSE is simpler operationally and fits progress updates well.

### OpenAPI

- Use OpenAPI from the start.
- Use a contract-first workflow: edit the OpenAPI spec as the API boundary, then generate backend and frontend artifacts.
- Use `oapi-codegen` for Go OpenAPI server/types generation.
- Use `openapi-typescript` and `openapi-fetch` for TypeScript client generation.
- Backend implementation must conform to the generated Go server interfaces/types.
- SvelteKit browser code must use the generated TypeScript client/types.
- Internal backend refactors that do not change API behavior do not require spec changes.
- API behavior changes should update the spec first, regenerate Go and TypeScript artifacts, then update backend implementation and web usage.

Rationale:

- The frontend and backend should not drift.
- API shape should be explicit from the first vertical slice.
- Refactoring order stays disciplined: API contract first, generated artifacts second, backend and UI implementation third.

### Backend HTTP Stack

- Use chi for routing.
- Use standard `net/http` primitives under chi.
- Keep handlers thin and delegate to application services.

Rationale:

- chi is small, idiomatic, and keeps the code close to standard Go.

### PostgreSQL Access

- Use pgx for PostgreSQL connectivity.
- Use sqlc for SQL-first generated data access.
- Avoid a heavy ORM while the domain model is still evolving.

Rationale:

- The project needs predictable queries, explicit transactions, and good PostgreSQL support.
- sqlc keeps SQL visible while giving typed Go accessors.

### Background Jobs

- Use River for background job management.
- Use River with PostgreSQL and pgx.
- River jobs should cover search, download polling, import, stream analysis, subtitle sync, muxing, cleanup, and long-running maintenance work.
- River's own required database schema setup is allowed as part of bootstrap/reset even though custom app migration tooling is deferred.

Rationale:

- The project already requires PostgreSQL.
- Media management needs durable, retryable, observable jobs.
- River fits the selected Go/PostgreSQL/pgx stack.

### Media Tool Distribution

- Bundle required media tools into the production Docker image.
- Local development runs use locally installed host tools.
- The app must detect and report tool versions and missing tools in system health.

Initial tools:

- `ffmpeg`
- `ffprobe`
- `mkvmerge`
- `mkvextract`
- `mediainfo`, if useful after implementation validation

Rationale:

- Docker deployments should work predictably without users manually building a media-tool environment.
- Local development stays transparent by using normal host-installed tools.

### Repository Layout

- Use a top-level Go module.
- Keep Go entry points under `cmd/<project-name>`.
- Keep backend implementation under `internal/`.
- Keep the SvelteKit app under `web/`.

Rationale:

- The backend is the application runtime and should own the root module.
- `web/` keeps the frontend clearly separated while still allowing production asset embedding/serving.

### Local Task Runner

- Use a `Makefile` for common local commands.
- Initial targets should include `dev`, `dev-api`, `dev-web`, `db-reset`, `test`, and `check`.

Rationale:

- Make is available broadly and keeps first-project commands obvious.
- More specialized tooling can be added later if the command surface grows.

## Still Open
- Whether SvelteKit should use adapter-static or a minimal custom static build strategy for Go asset serving.
