---
title: Configuration
description: Runtime environment variables and local development paths.
---

## Runtime Variables

| Variable | Purpose | Default |
| --- | --- | --- |
| `ADDR` | Go API bind address | `:18080` |
| `APP_ENV` | Runtime environment name | `development`; Docker image uses `production` |
| `DATABASE_URL` | PostgreSQL connection string | local Docker Postgres |
| `MEDIA_DATA_DIR` | Media working directory | `.data/media`; Docker image uses `/data` through compose |
| `WEB_DIR` | Static web build served by Go | `web/build`; Docker image uses `/app/web` |
| `APP_VERSION` | Version shown in system status | `0.0.0-dev` |
| `APP_COMMIT` | Commit shown in system status | `dev` |
| `APP_SOURCE_URL` | Source repository URL shown in status | `Not configured` |
| `ADMIN_USERNAME` | Default admin username | `admin` |
| `ADMIN_PASSWORD` | Default admin password | `admin` |
| `SESSION_TTL` | Session cookie lifetime | `24h` |

## Development Database Commands

| Command | Purpose |
| --- | --- |
| `make db-clean` | Drop and recreate the `app` schema. |
| `make db-reset` | Clean, migrate, apply defaults, apply dev defaults, then try local dev seed. |
| `make db-seed-local` | Apply only the ignored local development seed. |

The local development seed path is `scripts/seeds/dev.local.sql`. It is ignored
by git and can contain machine-local providers, paths, or credentials.

## Media Tools

The application probes and processes media through installed command-line tools.
Local development should have these available when testing file details,
preview, track deletion, subtitles, or assembly flows:

- `ffmpeg`
- `ffprobe`
- `mkvmerge`
- `mkvextract`
- `mediainfo`

For in-app setup after the server is running, use the
[Setup Guide](/user-guide/using/setup-guide/) instead of editing files by hand.
