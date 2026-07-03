# OpenAPI Contract Workflow

Status: Draft

Codename `Mema` remains temporary until the final project name is chosen.

## Summary

The project uses a contract-first OpenAPI workflow. The OpenAPI document is the source of truth for public API behavior. Go backend code and SvelteKit frontend API clients are generated from the same contract.

## Decision

- Use OpenAPI from the start.
- Keep the OpenAPI contract in the repository.
- Generate Go server/types from the contract with `oapi-codegen`.
- Generate TypeScript types/client helpers with `openapi-typescript` and `openapi-fetch`.
- SvelteKit browser code calls the Go API directly through the generated TypeScript client.
- Go serves the built frontend assets in production.
- No SvelteKit server runtime is used in production.

## Proposed File Layout

```text
api/
  openapi.yaml
  generated/
    go/
    ts/
cmd/
  <project-name>/
internal/
  httpapi/
  ...
web/
  src/
    lib/
      api/
```

Initial contract file: [api/openapi.yaml](../api/openapi.yaml)

## Contract Rules

- API behavior changes start in `api/openapi.yaml`.
- Generated files must not be manually edited.
- Backend handlers must conform to generated Go interfaces/types.
- Frontend calls must use the generated TypeScript client.
- Internal backend refactors that do not change API behavior do not require contract changes.
- Breaking API changes must be intentional and reflected in the contract.
- Error response shapes must be defined in OpenAPI, not invented ad hoc in handlers.

## Change Workflow

1. Edit `api/openapi.yaml`.
2. Generate Go API artifacts.
3. Generate TypeScript API artifacts.
4. Implement or adjust Go handlers.
5. Update SvelteKit usage.
6. Run backend tests.
7. Run frontend checks.

## Initial API Surface

The scaffold should start with a small API contract:

- `GET /api/health`
- `GET /api/system/tools`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/session`
- `GET /api/events`

The OpenAPI server URL is `/api`, so path entries inside the contract omit the `/api` prefix.

## Initial Schemas

Required starter schemas:

- `HealthResponse`
- `ToolStatus`
- `ToolStatusResponse`
- `LoginRequest`
- `SessionResponse`
- `ErrorResponse`
- `EventEnvelope`

## Error Model

Use a consistent error shape:

```yaml
ErrorResponse:
  type: object
  required:
    - code
    - message
  properties:
    code:
      type: string
    message:
      type: string
    details:
      type: object
      additionalProperties: true
```

## Authentication Model

- The Go API owns authentication.
- Use local admin login in version 0.
- Use secure session cookies.
- Browser calls include cookies.
- Generated client should be configured to include credentials for same-origin API calls.
- The initial cookie name is `session`.
- The only initial user role is `admin`.

## SSE Model

`GET /api/events` is the first Server-Sent Events endpoint.

Initial event types:

- `system.heartbeat`
- `job.created`
- `job.updated`
- `job.completed`
- `job.failed`

The event schema should remain generic enough to support queue, search, import, subtitle, and assembly progress later.

The OpenAPI contract represents `text/event-stream` as a string response while documenting that each SSE `data` value is an `EventEnvelope` JSON object.

## Generation Commands

Exact commands will be finalized during scaffolding. The expected Makefile targets are:

```text
make api-generate
make api-check
```

`make check` should include API generation drift checks once the generator setup exists.

## Open Questions

- Should generated Go code live under `internal/httpapi/generated` instead of `api/generated/go`?
- Should generated TypeScript code live under `web/src/lib/api/generated` instead of `api/generated/ts`?
- Should the first generated Go path use strict server interfaces from `oapi-codegen`?
- Should the API include `/api/version` separately from `/api/health`?
