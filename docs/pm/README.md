# Product Management Status

Last audited: **2026-07-14**

This folder is the product-level inventory for Mema. It reconciles the local documentation, the implemented backend and frontend surfaces, and the GitHub issue tracker.

## Documents

- [Current features](current-features.md) — capabilities supported by code and/or closed implementation evidence.
- [Missing and future features](missing-and-future.md) — open backlog, known limitations, conditional future ideas, and documentation debt.

## Status legend

- `[x]` — implemented. This means code or a reachable product surface exists and the issue evidence does not say the work is still open.
- `[ ]` — not implemented, incomplete, unverified, or awaiting a product decision.
- A nested checked item under an unchecked item means the foundation exists but the user-facing or follow-up work remains.

## Audit rules

1. An open GitHub issue is treated as not complete, even when related foundations exist.
2. A closed epic is treated as implemented only when its completion note and local code agree.
3. A draft PRD is not automatically a missing feature. Several PRDs were implemented by closed epics but were never moved out of `Draft`.
4. Explicit limitations in the user guide are tracked even when no GitHub issue exists.
5. Generated `Unimplemented` methods in `internal/httpapi/openapi.gen.go` are interface fallbacks, not proof that the real server route is missing.
6. This is a product-capability catalogue, not a list of every OpenAPI operation or UI component.

## Snapshot

- The product has an end-to-end movie, series, anime-aware, subtitle, acquisition, profile-fulfillment, and DLNA foundation.
- The target-satisfaction epic ([#122](https://github.com/ntrp/project-mema/issues/122)) and DLNA epic ([#88](https://github.com/ntrp/project-mema/issues/88)) are closed with implementation evidence.
- **11 GitHub issues remain open** as of the audit date.
- The most important release-readiness gaps are first-start onboarding, backup/restore, CI enforcement, and decisions about advanced operational integrations.
- The most important domain expansion gaps are AniList, AniDB, editable anime aliases, conflict review, and manual numbering overrides.

## Primary evidence

- Product and user docs: `README.md`, `docs/website/src/content/docs/`
- Product specifications: `docs/prd/`
- Architecture decisions: `docs/adr/`
- API contract: `api/openapi.yaml`
- Backend implementation: `internal/`, especially `internal/httpapi/`, `internal/jobs/`, `internal/dlna/`, `internal/storage/`, and `internal/satisfaction/`
- Frontend implementation: `web/src/routes/`, `web/src/lib/features/`, and `web/src/lib/components/`
- Behavior catalogue: `features/behavior/`
- GitHub issues: <https://github.com/ntrp/project-mema/issues>
