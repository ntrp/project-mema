---
title: Storage
description: Database schema, sqlc, and storage ownership rules.
---

The application uses PostgreSQL. Application code reaches storage through the
`internal/storage.SettingsStore` facade.

## Schema

The project is pre-release, so schema changes are applied directly to:

```txt
internal/storage/migrations/00001_initial_schema.sql
```

Do not add follow-up migrations while the project remains unreleased.

## Generated Queries

Storage query files live under:

```txt
internal/storage/queries
```

Generated sqlc code is committed under:

```txt
internal/storage/generated
```

Run these after storage query or schema edits:

```sh
make sqlc-generate
make verify-sqlc-generated
```

## ADR

The storage access strategy is documented in:

[`docs/adr/0001-storage-access-strategy.md`](https://github.com/ntrp/project-mema/blob/main/docs/adr/0001-storage-access-strategy.md)
