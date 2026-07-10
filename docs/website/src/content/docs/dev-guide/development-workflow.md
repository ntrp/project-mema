---
title: Development Workflow
description: Commands and repository rules for contributors.
---

## Start Local Services

```sh
docker compose up -d postgres
make db-reset
make dev-watch
```

`make dev-watch` starts the frontend through Vite dev server and runs the API
watcher. Backend source changes rebuild and restart the API. Changes to
`api/openapi.yaml` regenerate both Go and TypeScript API bindings before the API
restarts.

## Verification

Use the focused command that matches the change. Common commands are:

```sh
make check
go test ./...
pnpm -C web run check
pnpm -C web exec vitest run <test-file>
```

When running Go tests directly, use the workspace-local cache path if your
environment cannot write to the default Go cache:

```sh
GOCACHE=/Users/ntrp/_pws/project-mema/.cache/go-build go test ./...
```

## Code Ownership

- Keep backend modules under 300 lines unless generated.
- Keep frontend modules and Svelte components under 200 lines.
- Prefer small feature folders over large mixed files.
- Use modal confirmation dialogs rather than browser confirm dialogs.
- Use the tooltip component rather than browser title tooltips.
- Keep generated files committed when the repository expects them.
