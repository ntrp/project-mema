---
title: System Overview
description: Runtime architecture and repository layout.
---

## Runtime Shape

- The Go server exposes the `/api` HTTP API.
- The built SvelteKit SPA is served by the Go server for production.
- PostgreSQL stores application state.
- River runs background jobs.
- Metadata, indexer, subtitle, and download-client services integrate with
  external providers.
- Media tools inspect and modify local files.

## Repository Layout

| Path | Purpose |
| --- | --- |
| `api/openapi.yaml` | API contract source of truth. |
| `cmd/server` | Final application entrypoint. |
| `cmd/devdb` | External development database tooling. |
| `internal` | Go application packages. |
| `web` | SvelteKit browser application. |
| `docs` | ADRs, PRDs, and this documentation website. |
| `features` | Behavior specifications and planning artifacts. |
| `scripts` | Verification and development helper scripts. |

## Development Boundary

Development reset, cleanup, and local seed logic is intentionally outside the
final app. The server starts the app; developer tooling prepares local database
state.
