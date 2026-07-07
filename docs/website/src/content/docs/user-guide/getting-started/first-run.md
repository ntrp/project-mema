---
title: First Run
description: Start the local database, reset the development schema, and run the app.
---

## Prerequisites

- Go 1.26 or newer.
- Node.js 24 or newer with pnpm through Corepack.
- Docker or a local PostgreSQL 17 instance.
- Optional local media tools: `ffmpeg`, `ffprobe`, `mkvmerge`, `mkvextract`,
  and `mediainfo`.

## Install Frontend Dependencies

```sh
make web-install
```

## Start PostgreSQL

```sh
docker compose up -d postgres
```

The default local database URL is:

```txt
postgres://media_manager:media_manager@localhost:15432/media_manager?sslmode=disable
```

## Reset The Development Database

```sh
make db-reset
```

Development database cleanup, reset, and local seed application are external
tools. They are not part of the production server command.

## Run The App

Start the Go API:

```sh
make dev-api
```

Start the browser app in another terminal:

```sh
make dev-web
```

The API listens on port `18080` and the frontend listens on port `15173` by
default.
