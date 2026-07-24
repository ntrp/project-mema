# Missing and Future Features

Unchecked items are not complete as of **2026-07-14**. The first section is the
authoritative open GitHub backlog. Later sections capture limitations and
product/documentation gaps found locally even when no issue exists.

## Open GitHub backlog

All 11 open issues are represented below.

### Release readiness and operations

- [ ] **Configuration and database backup/restore** — define a versioned backup
      format and secret policy; add API/jobs plus create, download, upload,
      preview, restore, and recovery UI
      ([#45](https://github.com/ntrp/project-mema/issues/45)).
- [ ] **GitHub Actions quality gates** — run Go tests, Svelte checks/tests,
      lint/format checks, and generated OpenAPI drift checks on pull requests
      ([#46](https://github.com/ntrp/project-mema/issues/46)).
- [ ] **Decide the advanced operations surface** — scope V0/V1 versus deferred
      capabilities and produce ADRs/follow-up issues
      ([#49](https://github.com/ntrp/project-mema/issues/49)).
  - [ ] API keys for automation, if accepted.
  - [ ] Webhook event types and delivery guarantees, if accepted.
  - [ ] Authenticated metrics and label policy, if accepted.
  - [ ] Reverse-proxy authentication boundaries, if accepted.
  - [ ] OAuth/OIDC boundaries, if accepted.
- [ ] **Close the operations hardening epic** after its remaining open children
      are resolved ([#7](https://github.com/ntrp/project-mema/issues/7)).
  - [x] Durable sessions are complete via
        [#43](https://github.com/ntrp/project-mema/issues/43).
  - [x] Secret redaction is complete via
        [#44](https://github.com/ntrp/project-mema/issues/44).
  - [x] Filesystem hardening is complete via
        [#47](https://github.com/ntrp/project-mema/issues/47).
  - [x] External media-tool hardening is complete via
        [#48](https://github.com/ntrp/project-mema/issues/48).
  - [ ] Backup/restore, CI, and advanced-hook scoping remain open.

### Anime metadata and correction workflows

The storage/matching foundation from
[#30](https://github.com/ntrp/project-mema/issues/30) is implemented; these
user-facing/provider follow-ups are not.

- [ ] **User-editable anime aliases** — add/edit/delete/disable user aliases
      without losing provider-derived values; define how aliases affect search
      and matching ([#68](https://github.com/ntrp/project-mema/issues/68)).
- [ ] **Anime provider-conflict review UI** — show TMDB/TVDB/AniList/AniDB
      disagreements and let users select canonical values while preserving
      source data ([#69](https://github.com/ntrp/project-mema/issues/69)).
- [ ] **AniList metadata provider** — search/details, aliases, artwork, dates,
      mappings, rate-limit/auth policy, and fixtures
      ([#70](https://github.com/ntrp/project-mema/issues/70)).
- [ ] **AniDB metadata provider** — search/details, absolute-number hints,
      mappings, specials/OVA policy, caching/credentials, and fixtures
      ([#71](https://github.com/ntrp/project-mema/issues/71)).
- [ ] **Manual anime numbering overrides** — correction workflow for provider
      mappings, split seasons, specials, OVAs, and complex series numbering
      ([#72](https://github.com/ntrp/project-mema/issues/72)).

### First-start experience

- [ ] **Mandatory first-start setup wizard** — create the first admin without
      default credentials; validate database readiness, TMDB, an indexer, a
      download client, movie/series root folders, path mappings, and required
      defaults before entering the normal app
      ([#73](https://github.com/ntrp/project-mema/issues/73)).
- [ ] **Skippable guided tour** — per-user bottom-left stepper covering core
      routes, with skip/finish/restart state and responsive/error behavior
      ([#74](https://github.com/ntrp/project-mema/issues/74)).

## Explicit current limitations without an open issue

These limitations are stated by local docs or exposed by the code audit, but the
GitHub tracker has no open implementation issue for them.

- [ ] **Bitmap-to-text subtitle conversion** — PGS/SUP tracks can be recognized
      and matched, but conversion between bitmap subtitles and
      SubRip/WebVTT/ASS/SSA is not implemented
      (`docs/website/src/content/docs/user-guide/using/qualities-formats-profiles.md`).
- [ ] **Decide the user-facing DLNA control-point surface** — backend
      AVTransport and RenderingControl clients exist, but there is no current
      route/API for renderer discovery or play/pause/seek/volume actions. Issue
      [#108](https://github.com/ntrp/project-mema/issues/108) closed the isolated
      contract foundation; open a dedicated follow-up only if the product should
      expose those controls.
- [ ] **Optional Mema MediaRenderer mode** — the server intentionally does not
      advertise MediaRenderer/AVTransport/RenderingControl services. A local
      playback target and renderer mode remain conditional future work from
      [#108](https://github.com/ntrp/project-mema/issues/108).

## Deferred or conditional product ideas

These are mentioned by draft PRDs but are not committed backlog items. They need
a product decision before implementation.

- [ ] **Generic profile rule engine** — the desired-state PRD explicitly defers
      a full rule engine until the profile model is stable
      (`docs/prd/media-profile-desired-state.md`).
- [ ] **Custom formats as first-class target rows** — the rollup PRD keeps
      custom formats as an upgrade layer unless a future model explicitly
      attaches them to targets
      (`docs/prd/target-satisfaction-and-media-state-rollup.md`).
- [ ] **Additional future profile fields and policies** — preserve
      extensibility, but only add fields backed by a concrete user need and
      fulfillment strategy.

## Documentation and planning gaps

- [x] **Update implemented PRD statuses** — reconciled all six PRDs with the
      application on 2026-07-14, including delivered scope, implementation
      evidence, and explicit remaining boundaries.
- [ ] **Reconcile DLNA “planned architecture” wording** —
      `docs/website/src/content/docs/architecture/dlna.md` still labels profile
      architecture as planned even though renderer profiles, overrides,
      diagnostics, and delivery behavior are implemented.
- [ ] **Refresh the operations epic body** —
      [#7](https://github.com/ntrp/project-mema/issues/7) lists durable sessions,
      secret redaction, filesystem safety, and media-tool safety as remaining
      even though their child issues are closed.
- [ ] **Refine open anime issues** —
      [#68](https://github.com/ntrp/project-mema/issues/68) through
      [#72](https://github.com/ntrp/project-mema/issues/72) are labeled
      `to refine` and contain placeholder acceptance criteria/open decisions.
- [ ] **Assign roadmap milestones/priorities** — none of the 11 open issues had
      a milestone during this audit, so sequencing beyond issue dependencies is
      not explicit.
- [ ] **Keep this PM inventory current** whenever an issue closes, a limitation
      is accepted/deferred, or a user-facing capability changes.

## Implemented work that should not be mistaken for missing work

The reconciled PRDs record these broad epics as implemented. Do not reopen them
without a new, specific gap.

- [x] Full DLNA/UPnP server epic and child stories #89–#108
      ([#88](https://github.com/ntrp/project-mema/issues/88)).
- [x] Renderer profiles, device overrides, diagnostics, compatibility tests, and
      control-point client contracts.
- [x] Target satisfaction, DB-backed facts, rollups, wanted rows, fulfillment
      jobs, API/UI, tests, and docs through child issues #123–#137
      ([#122](https://github.com/ntrp/project-mema/issues/122)).
- [x] Anime provider mapping/alias/numbering storage and matching foundation
      ([#30](https://github.com/ntrp/project-mema/issues/30)); only the five open
      follow-ups above remain.
