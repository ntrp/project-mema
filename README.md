# Media Manager Codename

This repository is the scaffold for a self-hosted media manager. The current
name is only a codename; keep generated public identifiers neutral until the
final project name is chosen.

The first product slice is video-first: movies, TV, anime-specific metadata and
season handling, subtitles, indexer search, torrent/NZB download clients,
quality profiles, and MKV track sourcing/assembly.

## Current Stack

- Backend: Go, chi, pgx, River, contract-first OpenAPI.
- Frontend: SvelteKit/Svelte 5 as a browser-only SPA.
- Database: PostgreSQL.
- Production serving: Go serves the built SvelteKit static files.
- API client: generated TypeScript types from `api/openapi.yaml` with
  `openapi-typescript` and `openapi-fetch`.
- Media tools: Docker image includes ffmpeg, ffprobe, mkvmerge, mkvextract, and
  mediainfo. Local development uses host-installed tools.

## License

This project is licensed under the GNU Affero General Public License v3.0 or
later. See [LICENSE](LICENSE).

## Repository Layout

- `api/openapi.yaml`: source-of-truth API contract.
- `cmd/server`: Go entrypoint.
- `internal`: private Go application packages.
- `web`: SvelteKit app, built as static files.
- `docs/website`: Astro/Starlight documentation website.
- `docs/adr` and `docs/prd`: architecture decisions and product docs.
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

Install documentation website dependencies:

```sh
make docs-install
```

Start PostgreSQL:

```sh
docker compose up -d postgres
```

PostgreSQL is exposed on host port `15432` by default so it can run alongside
other local PostgreSQL instances. If that port is already in use, choose another
host port:

```sh
POSTGRES_PORT=16432 docker compose up -d postgres
export DATABASE_URL=postgres://media_manager:media_manager@localhost:16432/media_manager?sslmode=disable
```

Optional Transmission over PIA OpenVPN setup:

```sh
cp docker/transmission-openvpn.env.example .data/secrets/transmission-openvpn.env
$EDITOR .data/secrets/transmission-openvpn.env
docker compose up -d transmission-openvpn
```

Transmission is exposed at `http://localhost:19091` by default. The service uses
PIA via `haugene/transmission-openvpn`, stores its ignored local config under
`.data/transmission-openvpn/config`, and stores downloads under
`.data/media/transmission`. Override the PIA region with
`TRANSMISSION_OPENVPN_CONFIG`, for example:

```sh
TRANSMISSION_OPENVPN_CONFIG=sweden docker compose up -d transmission-openvpn
```

Create the development schema:

```sh
make db-reset
```

Run the Go API:

```sh
make dev-api
```

The API listens on `18080` by default. If that port is already in use:

```sh
ADDR=:19080 make dev-api
```

Run the SvelteKit dev server in another terminal:

```sh
make dev-web
```

Vite listens on `15173` by default and proxies `/api` to
`http://127.0.0.1:18080`. When the Go API uses a non-default port, point the
Vite proxy at it:

```sh
VITE_API_PROXY_TARGET=http://127.0.0.1:19080 make dev-web
```

Open the frontend URL printed by Vite. The dev frontend calls the Go API at
`/api`; for full same-origin behavior, use the built app served by Go.

Run the documentation website:

```sh
make docs-dev
```

The documentation website listens on `15174` by default.

## Build And Check

Generate OpenAPI clients and server code:

```sh
make api-generate
```

Verify committed OpenAPI generated artifacts without rewriting them:

```sh
make verify-generated
```

Run backend tests, Svelte checks, linting, and formatting checks:

```sh
make check
```

Run the behavior-backed API/integration and E2E suites:

```sh
make test-api
make test-e2e
```

Behavior scenarios live in `features/behavior` and use stable IDs such as
`SCN-AUTH-001`. Gherkin scenarios drive API/integration and E2E tests; focused
Go and TypeScript unit tests remain normal code tests and reference the relevant
scenario ID when they cover cataloged behavior.

Generate coverage reports with the shared 60% target:

```sh
make coverage
```

Backend reports are written to `coverage/`. Frontend reports are written to
`web/coverage/`.

Build the static web app and Go server binary:

```sh
make build
```

Run the built server after `make build`:

```sh
WEB_DIR=web/build ./bin/server
```

Build the unified Docker image with the Go backend, built SvelteKit files, and
required media tools:

```sh
make docker-build
```

Run the app image with the compose-managed PostgreSQL service:

```sh
docker compose up --build app
```

For a direct container run, point `DATABASE_URL` at a reachable PostgreSQL
instance and mount media storage at `/data`:

```sh
docker run --rm -p 18080:18080 \
  -e DATABASE_URL=postgres://media_manager:media_manager@host.docker.internal:15432/media_manager?sslmode=disable \
  -v "$PWD/.data/media:/data" \
  project-mema:local
```

## Contract-First Workflow

1. Update `api/openapi.yaml`.
2. Run `make api-generate`.
3. Implement or refactor the Go backend to satisfy the generated server
   interface.
4. Update the Svelte UI against the generated TypeScript API types.
5. Run `make check` before continuing. It includes `make verify-generated`, so
   stale generated Go or TypeScript artifacts fail the check instead of being
   rewritten silently.

OpenAPI is currently pinned to 3.0.3 because the selected Go generator has
better support for OpenAPI 3.0 than 3.1.

## File Naming Templates

File naming settings are edited in the web UI at `/settings/library`.
Templates use tokens wrapped in braces, for example
`{movie_title} ({release_year}) {quality_full}`.

When typing inside a template token, the UI shows an autocomplete list
fuzzy-filtered by parameter name. Each suggestion displays the parameter and an
example value. Press Up or Down to move through the filtered suggestions; after
one second, the highlighted row shows the parameter description in a tooltip.
Press Enter or Tab to insert the highlighted parameter.

Season and episode numbers expose padded variants such as `{season:0}`,
`{season:00}`, `{episode:0}`, `{episode:00}`, and `{episode:000}`.

## Database Workflow

The application applies schema migrations and production defaults on startup.
Development-only database cleanup, reset, and local seed application live outside
the server command:

```sh
make db-clean
make db-reset
make db-seed-local
```

`make db-reset` drops and recreates the `app` schema, runs migrations, applies
tracked defaults, applies tracked development defaults, and then tries the
gitignored `scripts/seeds/dev.local.sql`. Missing or invalid local seed files
are skipped with a warning. Destructive reset commands only run against local
database hosts unless `ALLOW_REMOTE_DEV_DB_RESET=true` is set.

River schema changes are handled separately through:

```sh
make river-migrate
```

## Storage Access

Application storage is being converted from handwritten pgx to `sqlc` behind
the `internal/storage.SettingsStore` facade. Converted query files live under
`internal/storage/queries`, with generated code committed under
`internal/storage/generated`.

Storage changes should follow
[`docs/adr/0001-storage-access-strategy.md`](docs/adr/0001-storage-access-strategy.md):
keep existing pgx code stable, make transaction ownership explicit, and only add
generated storage code behind storage wrappers.

Regenerate sqlc storage artifacts after editing schema or storage queries:

```sh
make sqlc-generate
```

Verify committed sqlc artifacts without rewriting them:

```sh
make verify-sqlc-generated
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
