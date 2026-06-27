# Media Manager Codename

This repository is the scaffold for a self-hosted media manager. The current
name is only a codename; keep generated public identifiers neutral until the
final project name is chosen.

The first product slice is video-first: movies, TV, anime-specific metadata and
season handling, subtitles, indexer search, torrent/NZB download clients,
quality profiles, and MKV track sourcing/assembly.

## Current Stack

- Backend: Go, chi, pgx, sqlc, River, contract-first OpenAPI.
- Frontend: SvelteKit/Svelte 5 as a browser-only SPA.
- Database: PostgreSQL.
- Production serving: Go serves the built SvelteKit static files.
- API client: generated TypeScript types from `api/openapi.yaml` with
  `openapi-typescript` and `openapi-fetch`.
- Media tools: Docker image includes ffmpeg, ffprobe, mkvmerge, mkvextract, and
  mediainfo. Local development uses host-installed tools.

## Repository Layout

- `api/openapi.yaml`: source-of-truth API contract.
- `cmd/server`: Go entrypoint.
- `internal`: private Go application packages.
- `web`: SvelteKit app, built as static files.
- `features`: PRDs, requirements, architecture decisions, and delivery notes.
- `tools/tools.go`: pinned code generation and CLI tool dependencies.

## Local Prerequisites

- Go 1.26 or newer.
- Node.js 24 or newer with pnpm through Corepack.
- Docker or a local PostgreSQL 17 instance.
- Optional local media tools: `ffmpeg`, `ffprobe`, `mkvmerge`, `mkvextract`,
  `mediainfo`.

## First Run

Install frontend dependencies:

```sh
make web-install
```

Start PostgreSQL:

```sh
docker compose up -d postgres
```

If port `5432` is already in use, choose another host port:

```sh
POSTGRES_PORT=55432 docker compose up -d postgres
export DATABASE_URL=postgres://media_manager:media_manager@localhost:55432/media_manager?sslmode=disable
```

Create the development schema:

```sh
ALLOW_DEV_RESET=true make db-reset
```

Run the Go API:

```sh
make dev-api
```

If port `8080` is already in use:

```sh
ADDR=:18080 make dev-api
```

Run the SvelteKit dev server in another terminal:

```sh
make dev-web
```

When the Go API uses a non-default port, point the Vite proxy at it:

```sh
VITE_API_PROXY_TARGET=http://127.0.0.1:18080 make dev-web
```

Open the frontend URL printed by Vite. The dev frontend calls the Go API at
`/api`; for full same-origin behavior, use the built app served by Go.

## Build And Check

Generate OpenAPI clients and server code:

```sh
make api-generate
```

Run backend tests, Svelte checks, linting, and formatting checks:

```sh
make check
```

Build the static web app and Go server binary:

```sh
make build
```

Run the built server after `make build`:

```sh
WEB_DIR=web/build ./bin/server
```

## Contract-First Workflow

1. Update `api/openapi.yaml`.
2. Run `make api-generate`.
3. Implement or refactor the Go backend to satisfy the generated server
   interface.
4. Update the Svelte UI against the generated TypeScript API types.
5. Run `make check` before continuing.

OpenAPI is currently pinned to 3.0.3 because the selected Go generator has
better support for OpenAPI 3.0 than 3.1.

## Database Workflow

The project starts with a simple protected reset path instead of migrations.
`make db-reset` only works when both conditions are true:

- `APP_ENV=development`
- `ALLOW_DEV_RESET=true`

This is intentionally destructive and should only be used before the project is
live. River schema changes are handled separately through:

```sh
make river-migrate
```

## Quality Baseline

- `gofmt` is the Go formatter.
- Go code stays in `internal` until there is a proven public package boundary.
- Svelte code uses Svelte 5 syntax and `svelte-check`.
- Generated files are committed only after regenerating from the contract.
- Backend API changes start at OpenAPI, then regenerate Go and TypeScript.

Useful development assistance found during setup:

- Official Svelte docs MCP: use for SvelteKit/Svelte 5 docs and autofixes.
- Browser/Playwright skills: use once the app has meaningful UI flows.
- GitHub app tools: useful after the repository is pushed and PRs/checks exist.
- Candidate Go skills to consider installing later:
  `samber/cc-skills-golang@golang-code-style`,
  `samber/cc-skills-golang@golang-error-handling`,
  `samber/cc-skills-golang@golang-testing`,
  `samber/cc-skills-golang@golang-security`, and
  `samber/cc-skills-golang@golang-performance`.

## Next Feature Slice

Start with persisted local authentication and app settings, then add indexer and
download-client configuration. That gives the UI a real protected workflow and
sets up the first end-to-end search/download/import path.
