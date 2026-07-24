# Current Features

This checklist describes the product capabilities found in the local docs, API/backend, frontend routes, tests, and closed GitHub issues.

## Platform, accounts, and operations

- [x] Self-hosted Go API and SvelteKit/Svelte 5 web application backed by PostgreSQL.
- [x] Docker and local-development startup flows, static frontend serving, and a contract-first OpenAPI workflow.
- [x] Local username/password authentication with login, logout, current-session bootstrap, and profile editing.
- [x] Database-backed sessions that survive process restarts and can be invalidated ([#43](https://github.com/ntrp/project-mema/issues/43)).
- [x] Administrator and regular-user roles with administrator-managed users.
- [x] Credential redaction and tests covering secrets in API responses, logs, events, and UI-facing state ([#44](https://github.com/ntrp/project-mema/issues/44)).
- [x] Filesystem root/path traversal protections ([#47](https://github.com/ntrp/project-mema/issues/47)).
- [x] Validated external media-tool execution, bounded output, timeouts, and safe argument construction ([#48](https://github.com/ntrp/project-mema/issues/48)).
- [x] Background jobs with progress, structured logs, retry behavior, manual controls, and system job views.
- [x] Server-sent application events, notices, system events, status, file logs, and log streaming.
- [x] Administrator system views for status, indexing history/cache, metadata history/cache, jobs, logs, and events.

## Discovery and metadata

- [x] TMDB metadata integration for movie and series search, details, discovery, people, artwork, trailers, similar titles, and recommendations.
- [x] TVDB metadata integration for series details, episodes, people, artwork, and facts.
- [x] Global and advanced search with metadata-provider selection.
- [x] Discovery home, trending/section feeds, movie and series feeds, pagination, and add/request actions.
- [x] Anime/animated discovery filters while retaining `movie` and `serie` as the core media types.
- [x] Metadata detail pages, collections, cast, crew, person pages, similar titles, and recommendations.
- [x] Discovery blacklist management.
- [x] Metadata refresh plus administrator cache/history inspection and clearing.

## Libraries, requests, and media state

- [x] Separate movie and series libraries with list and detail surfaces.
- [x] First-class seasons and episodes, monitoring controls, and series-specific actions.
- [x] Anime classification, provider-mapping storage, provider-derived aliases, numbering strategies, and absolute-number matching foundation ([#30](https://github.com/ntrp/project-mema/issues/30)).
- [x] User media requests, request detail, and administrator approval.
- [x] Movie and series root folders, server-visible folder selection, folder creation, scans, and path mappings.
- [x] Library import, matching, manual import, reset, and rescan workflows.
- [x] Media add, update, monitor/unmonitor, metadata refresh, and delete actions.
- [x] Stored file history, probe facts, provenance, rename preview/apply, file deletion, track deletion, and rescanning.
- [x] Browser preview, HLS segments, range streaming, subtitle preview, and time-limited VLC playlists.

## Search, acquisition, and download activity

- [x] Indexer catalogue, indexer configuration/testing, health, rate limiting, search cache, and search history.
- [x] Manual release search, automatic search, fallback queries, candidate ranking, and release grabbing.
- [x] RSS synchronization and matching.
- [x] Transmission and SABnzbd download-client configuration, testing, categories, status, progress, cancellation, and imported-file handling.
- [x] Activity queue, history, blocklist, cancellation, deletion, and manual import controls.
- [x] Release decisions that combine quality, custom-format scores, target warnings, upgrade rules, and blocklist state.

## Profiles, desired state, and wanted work

- [x] Quality definitions and sizes, quality ranking, upgrade-until rules, and minimum upgrade thresholds.
- [x] Custom formats with parsing tests, score thresholds, upgrade-until scores, and score-increment rules.
- [x] Media profiles with video, audio, subtitle, container, conversion, retention, and upgrade policies.
- [x] Video targets for quality/resolution, codec, HDR, pixel format, and container-related requirements.
- [x] Audio targets for language, codec, channels, bitrate, scoring, retention, and lossy-conversion policy.
- [x] Subtitle targets for language, format, score, embedded/external/mixed placement, and retention.
- [x] Database-backed target satisfaction states: missing, partial, pending, satisfied, upgradeable, blocked, and failed.
- [x] Media/file rollups, target/candidate visual states, and custom-format upgrade evaluation.
- [x] Wanted rows for missing media, unsatisfied targets, failed/blocked work, and upgrades.
- [x] Settings-driven recalculation and matching manual actions for automatic fulfillment jobs ([#122](https://github.com/ntrp/project-mema/issues/122)).

## Media processing, tracks, and provenance

- [x] `ffprobe` extraction of streams, chapters, container, duration, bitrate, language, title, codec, and keyframes.
- [x] Shared delivery planning for direct play, direct stream/remux, and transcode.
- [x] Video fulfillment and conversion planning.
- [x] Audio fulfillment from current or alternate releases plus explicit audio conversion policy.
- [x] Container remuxing and profile-target container handling.
- [x] Component-source retention, compatibility evaluation/review, extraction, release, and assembly actions.
- [x] MKV mux/assembly jobs with per-stream provenance and media-detail review/status ([#41](https://github.com/ntrp/project-mema/issues/41), [#42](https://github.com/ntrp/project-mema/issues/42)).
- [x] Track and subtitle source provenance persisted through final assembly.
- [x] Removal of unwanted tracks when profile policy explicitly allows it.

## Subtitles

- [x] Generated subtitle-provider catalogue with configuration forms, credential handling, enable/disable, and provider tests.
- [x] Manual and automatic subtitle search, candidate display, grab/download, progress, and fulfillment.
- [x] Embedded, external, and mixed subtitle placement.
- [x] Text subtitle conversion among SubRip, WebVTT, ASS, and SSA.
- [x] PGS/SUP bitmap subtitle recognition and matching when an existing bitmap track satisfies a target.
- [x] Subtitle metadata editing and deletion.

## DLNA and UPnP

- [x] LAN-scoped, disabled-by-default DLNA Digital Media Server with interface/CIDR controls and bounded streams/transcodes.
- [x] SSDP advertisement and discovery responses.
- [x] Device description, SCPD documents, SOAP actions/faults, and service state.
- [x] ContentDirectory object model, browse, search, sort, pagination, and system update IDs.
- [x] DIDL-Lite metadata/resources for movies, series, episodes, collections, genres, years, and recent content.
- [x] ConnectionManager protocol reporting.
- [x] Shared direct, remux, transcode, seek, artwork, thumbnail, and subtitle resource delivery.
- [x] UPnP event subscriptions and lifecycle handling.
- [x] Seeded and editable renderer profiles, matching rules, per-device overrides, clone/import/export/reset/restore, and compatibility controls.
- [x] Recent-device, profile-match, and delivery-decision diagnostics in settings/system UI.
- [x] DLNA security limits, audit logging, compatibility fixtures, and a device matrix.
- [x] Authenticated DLNA control point with remote MediaRenderer discovery, Play on device, transport, seek, volume, and mute controls ([#139](https://github.com/ntrp/project-mema/issues/139)).

## Frontend product surfaces

- [x] Authenticated application shell with responsive navigation, top search, notices, modals, and profile/logout controls.
- [x] Home/discover, search, movies, series, wanted, requests, activity, history, blocklist, and blacklist routes.
- [x] Media, metadata, collection, cast, crew, people, similar, and recommendation routes.
- [x] Settings for general behavior, library, download clients, indexers, quality, profiles, custom formats, metadata, subtitles, DLNA, languages, tags, and users.
- [x] Administrator-only system routes for status, DLNA, indexing, metadata, jobs, logs, and events.
- [x] Shared modal-based confirmations and application tooltip primitives rather than browser-native confirmation/tooltip behavior.

## Engineering support already present

- [x] Backend unit, integration, acceptance, and behavior-oriented tests.
- [x] Frontend Svelte checks, Vitest tests, linting, formatting, and coverage commands.
- [x] OpenAPI-generated Go and TypeScript contracts.
- [x] Local `make check` workflow covering generation and core backend/frontend verification.
- [x] Public user, developer, and architecture documentation website sources under `docs/website/`.
