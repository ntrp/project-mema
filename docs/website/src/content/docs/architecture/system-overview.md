---
title: System Overview
description: Runtime architecture and repository layout.
---

## Runtime Shape

- The Go server exposes the `/api` HTTP API.
- The built SvelteKit SPA is served by the Go server for production.
- The production Docker image builds the SPA in a Node stage, builds the Go
  server in a Go stage, then copies both into a Debian runtime image with the
  required media tools.
- PostgreSQL stores application state.
- River runs background jobs.
- System job schedules, execution snapshots, and per-run structured logs are
  stored in PostgreSQL so the System > Jobs view can show fixed schedules,
  active one-shot work, and finished execution history independently from
  River's queue table.
- Metadata, indexer, subtitle, and download-client services integrate with
  external providers.
- Media tools inspect and modify local files.

## Repository Layout

| Path               | Purpose                                         |
| ------------------ | ----------------------------------------------- |
| `api/openapi.yaml` | API contract source of truth.                   |
| `cmd/server`       | Final application entrypoint.                   |
| `cmd/devdb`        | External development database tooling.          |
| `internal`         | Go application packages.                        |
| `web`              | SvelteKit browser application.                  |
| `docs`             | ADRs, PRDs, and this documentation website.     |
| `features`         | Behavior specifications and planning artifacts. |
| `scripts`          | Verification and development helper scripts.    |

## Development Boundary

Development reset, cleanup, and local seed logic is intentionally outside the
final app. The server starts the app; developer tooling prepares local database
state.

## Production Image

The root `Dockerfile` is the supported unified image build. It copies the
SvelteKit static output to `/app/web`, sets `WEB_DIR=/app/web`, and runs the Go
server on `ADDR=:18080`. API and DLNA routes are handled by backend routers;
other paths fall back to the SPA shell.

## Frontend Route Ownership

The SvelteKit shell owns authentication, navigation, global search, notices, and
shared modals. Route components own route-specific surfaces and should request
only the backend data needed by the active route. The app shell controller uses
`routeData.ts` to load route-scoped data instead of preloading settings,
library, activity, and system data for every page. New route work should add a
dedicated route component under `web/src/lib/features` and avoid routing through
generic section switches when the route has its own page.

## Background Jobs

River remains the execution engine for background work. The app mirrors job
lifecycle changes into `app.system_job_executions` and records structured
execution logs in `app.system_job_execution_logs`. Execution rows keep
structured progress data alongside the progress label and percent. Progress data
can identify media item, media title, file path, target, phase, unit counts,
timestamps, and pending provider or tool operation when a worker knows those
fields. Jobs are currently inserted with one maximum attempt, so worker failures
are surfaced directly instead of River retrying them. Fixed scheduled jobs are
registered from the application catalog and synchronized into
`app.system_job_schedules`, where category, description, automatic/manual flags,
pause state, and interval settings are persisted.

Schedules can be disabled by pausing them. Media Refresh is enabled by default
and refreshes file metadata for every media item using the same rescan path as
the media-detail Refresh file metadata action. Media Fulfillment is registered
disabled by default and scans for needed media operations when enabled or run
manually. Disabled automatic schedules do not enqueue periodic work, but
media-detail manual actions can still enqueue the matching operation when the
current row has enough context. Schedules can also be marked configurable.
Configurable schedules keep the catalog minimum as their River tick, but the
persisted interval decides whether a run is due. This lets administrators raise
or lower automatic fulfillment intervals without a server restart.

Routine schedules, such as download client activity sync, are marked with a
separate history policy. Routine successful runs are hidden from the default
history view and use shorter retention, while failures remain visible so regular
health checks do not bury meaningful background work. All River jobs are
non-retryable. A failed one-shot execution is finalized as discarded, while a
failed recurring execution is finalized and can be enqueued again by its
schedule after the configured interval. Manual schedule runs enqueue the same
fixed job definition with application schedule metadata, so they update the
fixed schedule's active run, history, and next run calculation instead of
appearing as unrelated one-shot work.

The `/api/events` stream publishes both `system.job.updated` for River row
changes and `system.job.execution.updated` for dashboard execution/progress
updates.
