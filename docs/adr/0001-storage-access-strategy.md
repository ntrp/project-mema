# ADR 0001: Storage Access Strategy

## Status

Accepted

## Context

The application uses PostgreSQL through `pgxpool` and a single
`storage.SettingsStore` facade. Most storage methods are handwritten SQL in
small files under `internal/storage`, with row scanners kept near the queries
that use them. The repository also has a `sqlc.yaml` configuration that points
at `internal/storage/queries` and `internal/storage/generated`.

The current code has a broad handwritten surface: dozens of storage files,
integration tests for important flows, and transaction-heavy methods for media,
profiles, indexers, library scans, release candidates, and request approval.
Moving all of that to `sqlc` at once would create broad churn while the schema
and workflow boundaries are still changing.

At the same time, letting every new query grow as ad hoc handwritten pgx would
make transaction boundaries, scanner behavior, nullable fields, and generated
API contracts harder to keep consistent as the storage layer expands.

## Decision

Use an incremental `sqlc` conversion strategy:

1. Keep storage entrypoints on `SettingsStore` so application code does not
   couple directly to generated query structs.
2. Convert bounded storage areas to `sqlc` behind wrappers, preserving existing
   validation, transaction ownership, and public error behavior.
3. Prefer simple read/write resources first, then move to transaction-heavy
   workflows only when their generated query shape is covered by focused tests.

This is not a request to bypass transaction boundaries or wrappers. Generated
queries should replace handwritten SQL while the storage package continues to
own application behavior.

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

## `sqlc` Conversion Rules

- Use one query file per bounded storage area under
  `internal/storage/queries/<area>.sql`.
- Commit generated code under `internal/storage/generated`.
- Keep generated files out of handwritten application logic outside
  `internal/storage`.
- Expose small wrappers from handwritten storage code when application errors,
  validation, JSON mapping, or transaction orchestration are needed.
- Add focused storage coverage for each converted write path.
- Run `make sqlc-generate` after editing query files or schema.
- Run `make verify-sqlc-generated` to check committed generated artifacts.

## Consequences

This keeps the current storage layer stable while the project is still evolving.
It also gives `sqlc` a real path to adoption without turning the existing pgx
surface into an unfinished partial migration.

The cost is temporary dual readiness: contributors must understand the current
handwritten pgx rules and the pilot constraints. That cost is bounded by keeping
generated code isolated until the pilot is accepted.

## Current Converted Areas

Converted query files:

- `internal/storage/queries/custom_formats.sql`
- `internal/storage/queries/database_status.sql`
- `internal/storage/queries/discover_blacklist.sql`
- `internal/storage/queries/download_activity.sql`
- `internal/storage/queries/download_clients.sql`
- `internal/storage/queries/file_naming.sql`
- `internal/storage/queries/imported_files.sql`
- `internal/storage/queries/indexer_bulk.sql`
- `internal/storage/queries/indexer_proxies.sql`
- `internal/storage/queries/indexer_search.sql`
- `internal/storage/queries/indexers.sql`
- `internal/storage/queries/languages.sql`
- `internal/storage/queries/library_folders.sql`
- `internal/storage/queries/library_scans.sql`
- `internal/storage/queries/log_file_settings.sql`
- `internal/storage/queries/media_items.sql`
- `internal/storage/queries/media_items_mutations.sql`
- `internal/storage/queries/media_profiles.sql`
- `internal/storage/queries/media_requests.sql`
- `internal/storage/queries/metadata_cache.sql`
- `internal/storage/queries/metadata_providers.sql`
- `internal/storage/queries/path_mappings.sql`
- `internal/storage/queries/quality_sizes.sql`
- `internal/storage/queries/release_blocklist.sql`
- `internal/storage/queries/release_candidates.sql`
- `internal/storage/queries/system_jobs.sql`
- `internal/storage/queries/system_events.sql`
- `internal/storage/queries/tags.sql`
- `internal/storage/queries/user_profile.sql`
- `internal/storage/queries/users.sql`

`SettingsStore` remains the application-facing wrapper for validation,
transaction ownership, JSON mapping, and error normalization.

## Follow-Up

- Continue converting handwritten storage SQL to `sqlc` query files.
- Move broad media, indexer, or transaction-heavy workflows only with focused
  query files and matching tests for transaction behavior.
