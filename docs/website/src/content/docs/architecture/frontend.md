---
title: Frontend Architecture
description: Frontend boundaries, state ownership, and real-time event delivery.
---

## Boundaries

SvelteKit owns routing in `src/routes`. Handwritten application code uses three
pragmatic boundaries under `src/lib`:

- `app` owns application-lifetime providers, session/bootstrap, navigation,
  notifications, and the real-time transport.
- `features` owns domain behavior and exposes a small public entry point per
  feature. Feature code must not import another feature's internals.
- `shared` owns domain-neutral UI, API helpers, and utilities.

Generated Orval bindings remain read-only. Feature API modules wrap them with
stable query keys, response mapping, error normalization, and invalidation rules.
New handwritten frontend modules and components must stay below 200 lines.

Settings routes load data by active section. A route requests only the resources
rendered by that section instead of calling the legacy aggregate settings loader.
Shared dependencies are explicit: indexers also load indexer-search configuration
and tags; library loads folders, path mappings, profiles, and metadata providers;
and profile editing loads profiles, custom formats, and languages. Sections whose
content is local or loaded by a nested component make no aggregate request.

## State Ownership

TanStack Query owns remote resources, including remote search results. Query
keys belong to the feature that understands the resource. Mutations invalidate
or update every affected key; optimistic updates are reserved for deterministic,
reversible operations.

Svelte runes own temporary browser state such as forms, selections, open modal
state, and navigation UI. A remote resource must not be mirrored into `$state`.

## Real-time Events

There is exactly one application notification SSE connection per authenticated
browser tab. The root session lifecycle starts it after authentication and stops
it during logout, session invalidation, or application teardown. Calls to start
the transport are idempotent, so route navigation and component remounts cannot
open duplicate `/api/events` streams.

The transport in `src/lib/app/realtime/appEventSource.ts` parses and dispatches
all `/api/events` messages. Consumers subscribe to its dispatcher or observe
TanStack Query and notification state; components never open `/api/events`
directly. Replayed event IDs are processed once, and a single reconnect strategy
prevents concurrent reconnect attempts.

ESLint enforces this ownership by rejecting `new EventSource(...)` outside the
transport. Two endpoint-specific exceptions are documented in configuration:
release-search progress and the administrator's live system-log diagnostic.
They are bounded interactive streams, not duplicate application notification
connections. Cross-tab connection sharing is intentionally out of scope.

## Migration Inventory

Keep this table current as each vertical slice moves out of the legacy shell.

| Slice             | Owner                                | Operations                                             | Query/cache policy                                            | SSE interaction                                         | Status      |
| ----------------- | ------------------------------------ | ------------------------------------------------------ | ------------------------------------------------------------- | ------------------------------------------------------- | ----------- |
| Activity          | `features/activity`                  | queue/blocklist reads; cancel, delete, import commands | Stable activity keys; commands invalidate affected lists      | Application events reconcile queue and blocklist caches | Complete    |
| Library and media | `features/library`, `features/media` | lists, requests, details, files, releases, and commands | Collection/detail keys; commands reconcile affected caches    | Media and job events reconcile affected resources       | Complete |
| Discovery         | `features/discovery`, `features/search` | sections, pagination, blacklist, and keyed searches  | Inputs and section pages are part of stable query keys         | Blacklist mutations filter or invalidate discovery data | Complete |
| Settings          | focused settings features            | section reads, scans, caches, profile, and admin commands | Authenticated domain keys; section-specific loading only    | System events update cache-owned inspection resources   | Complete |
| Session           | `app/session`                        | session lookup, login, logout, and cache lifecycle     | Logout clears every privileged feature namespace              | Owns application SSE start and stop                     | Complete |

For every migrated operation, record its generated API operation, semantic type
(query or mutation), exact key factory, invalidation/update behavior, and event
types in the owning feature's tests or adjacent documentation.

## Verification

Each slice adds characterization tests before migration and query/mutation tests
afterward. Tests cover key stability, invalidation, rollback, event deduplication,
connection lifecycle, and cache reconciliation. Handwritten files require at
least 60% coverage; generated declarations and non-business data objects are the
only intended long-term exclusions.
