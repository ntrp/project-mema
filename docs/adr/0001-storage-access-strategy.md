# ADR 0001: Storage Access Strategy

## Status

Accepted

## Context

The application uses PostgreSQL through `pgxpool` and a single
`storage.SettingsStore` facade. Most storage methods are handwritten SQL in
small files under `internal/storage`, with row scanners kept near the queries
that use them. The repository also has a `sqlc.yaml` placeholder that points at
`internal/storage/queries` and `internal/storage/generated`, but there are no
real query files or generated storage packages yet.

The current code has a broad handwritten surface: dozens of storage files,
integration tests for important flows, and transaction-heavy methods for media,
profiles, indexers, library scans, release candidates, and request approval.
Moving all of that to `sqlc` at once would create broad churn while the schema
and workflow boundaries are still changing.

At the same time, letting every new query grow as ad hoc handwritten pgx would
make transaction boundaries, scanner behavior, nullable fields, and generated
API contracts harder to keep consistent as the storage layer expands.

## Decision

Use a staged storage strategy:

1. Keep existing handwritten pgx code in place unless a bounded refactor is
   already needed for the feature being changed.
2. Use the next storage implementation issue (#51) to pilot `sqlc` in one small,
   representative storage area before adopting it more broadly.
3. Treat the pilot as the decision point for wider adoption. If it improves
   maintainability without forcing awkward transaction or JSON handling, future
   storage areas may move to `sqlc` incrementally. If it adds more friction than
   value, keep handwritten pgx and codify that pattern instead.

This is not a blanket migration. New work should not introduce generated
storage code outside the pilot area until the pilot has proven the convention.

## Rules For New Storage Code

- Keep storage entrypoints on `SettingsStore` unless a separate store has a
  clear ownership boundary.
- Keep each storage file focused on one resource or workflow.
- Put reusable scanner helpers near the resource that owns the row shape.
- Translate `pgx.ErrNoRows` to `storage.ErrNotFound` at public storage
  boundaries.
- Prefer typed input structs for writes rather than passing loosely related
  scalar arguments through call chains.
- Keep SQL ordered and explicit; avoid broad `select *` queries in application
  code.
- Add unit or integration coverage for new write paths, validation rules, and
  transaction behavior.

## Transaction Rules

- Start transactions at the storage method that owns the whole write workflow.
- Pass a narrow querier interface into helpers that must work with both
  `pgxpool.Pool` and `pgx.Tx`.
- Always defer rollback immediately after a successful `Begin`, ignoring the
  rollback error after commit.
- Return the committed entity from the same transaction when the caller depends
  on derived rows, tags, or related records.
- Do not start nested transactions inside helper functions.
- Keep River job transactions separate from application storage transactions
  unless a workflow explicitly requires one atomic boundary.

## `sqlc` Pilot Rules

The pilot in #51 should:

- use one bounded storage area with representative reads, writes, and tests;
- place source queries under `internal/storage/queries/<area>.sql`;
- place generated code under `internal/storage/generated`;
- keep generated files out of handwritten application logic;
- expose a small wrapper from handwritten storage code when application errors,
  validation, or transaction orchestration are needed;
- document the exact generation command and whether generated code is committed.

The pilot should avoid media item or indexer-wide rewrites. Good candidates are
small settings-style resources with clear CRUD and limited JSON payload
handling.

## Consequences

This keeps the current storage layer stable while the project is still evolving.
It also gives `sqlc` a real path to adoption without turning the existing pgx
surface into an unfinished partial migration.

The cost is temporary dual readiness: contributors must understand the current
handwritten pgx rules and the pilot constraints. That cost is bounded by keeping
generated code isolated until the pilot is accepted.

## Pilot Outcome

#51 applies this ADR to quality-size settings as the first `sqlc` pilot. That
area keeps domain validation and transaction ownership in handwritten storage
code while moving list, ensure, and upsert SQL into
`internal/storage/queries/quality_sizes.sql`. Generated code is committed under
`internal/storage/generated` and refreshed with:

```sh
make storage-generate
```

Future generated storage areas should copy that wrapper shape instead of
calling generated queries directly from HTTP handlers or jobs.

## Follow-Up

- Add a storage-generation drift check before moving additional storage areas.
- Migrate another bounded storage area only after the quality-size pilot remains
  stable through normal feature work.
