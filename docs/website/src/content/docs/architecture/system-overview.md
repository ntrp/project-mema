---
title: System Overview
description: Runtime architecture and repository layout.
---

## Runtime Shape

- The Go server exposes the `/api` HTTP API.
- The built SvelteKit SPA is served by the Go server for production.
- PostgreSQL stores application state.
- River runs background jobs.
- System job schedules, execution snapshots, and per-run structured logs are
  stored in PostgreSQL so the System > Jobs view can show fixed schedules,
  one-shot work, and execution history independently from River's queue table.
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

## Background Jobs

River remains the execution engine for background work. The app mirrors job
lifecycle changes into `app.system_job_executions` and records structured
execution logs in `app.system_job_execution_logs`. Fixed scheduled jobs are
registered from the application catalog and synchronized into
`app.system_job_schedules`, where pause state is persisted.

Schedules can be marked configurable. Configurable schedules keep the catalog
minimum as their River tick, but the persisted interval decides whether a run is
due. This lets the download client activity sync run as often as every 15
seconds while still allowing administrators to raise the interval without a
server restart.

Routine schedules, such as download client activity sync, are marked with a
separate history policy. Routine successful runs are hidden from the default
history view and use shorter retention, while failures and retryable executions
remain visible so regular health checks do not bury meaningful background work.

The `/api/events` stream publishes both `system.job.updated` for River row
changes and `system.job.execution.updated` for dashboard execution/progress
updates.
