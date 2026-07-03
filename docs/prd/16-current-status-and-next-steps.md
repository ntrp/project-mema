# Current PRD Status And Next Steps

Status: Draft

Date: 2026-06-28

Codename `Mema` remains temporary until the final project name is chosen.

## Summary

The project is past the initial skeleton and has a first monitored-media/search/grab slice in place.

The current implementation can:

- Run a Go API and SvelteKit frontend with OpenAPI-generated Go and TypeScript artifacts.
- Serve the built frontend from Go for production.
- Run PostgreSQL through Docker Compose.
- Reset the development database/schema.
- Authenticate a local admin session.
- Configure multiple indexers and download clients from settings.
- Test Torznab, Newznab, RSS, Transmission, and SABnzbd integrations.
- Add, list, view, and remove monitored movie/series media items.
- Show movie and series cards with detail pages.
- Configure TMDB and TVDB metadata providers.
- Search metadata providers for movies and series and cache search results to reduce repeated provider queries.
- Search configured indexers for a monitored item through a River job.
- Store release candidates and search errors.
- Enqueue a release grab through a River job and send it to the first enabled download client.
- Show basic download activity.
- Detect local or container media tools.

The current implementation does not yet close the full MVP vertical slice because it does not yet have quality profiles, release scoring, download completion polling, import, stream analysis, subtitles, or track assembly.

## PRD Status Matrix

| PRD | Status | Implemented | Not Yet Implemented |
| --- | --- | --- | --- |
| `00-product-vision-and-scope.md` | Partial | Main architecture and video-first direction are reflected. ARR-suite upstream integration is not used. Books and music are deferred. | Sonarr-grade profiles, metadata providers, import, subtitle workflow, and independent component assembly are not implemented. |
| `01-media-library-management.md` | Partial | Movie/series monitored items exist. Items can be added, listed, opened in detail view, removed, and discovered from configured library folders with auto/manual matching. | TV seasons/episodes, anime season behavior, imported files, stream metadata, rename rules, states, history, manual import, and provenance are missing. |
| `02-metadata-and-identifiers.md` | Partial | TMDB and TVDB providers can be configured, tested, searched, and cached. Media items can store external provider IDs, overview, and poster paths. | AniList, AniDB mapping, aliases, metadata refresh jobs, richer canonical provider precedence, and ambiguous import handling are missing. |
| `03-indexer-search-release-discovery.md` | Partial | Multiple indexers can be configured, tested, searched, and enabled/disabled. Torznab, Newznab, RSS, caps parsing, and release candidate storage exist. | RSS sync, per-library indexer selection, tags, rate limits, health history, dedupe rules, rich normalization, rejection reasons, and score explanations are missing. |
| `04-download-clients-import-pipeline.md` | Partial | Transmission and SABnzbd can be configured, tested, and used for initial grab requests through jobs. Activity records queued/grabbed/failed states. | Download polling, completion detection, remote path mappings, import decisions, copy/move/hardlink, failed-download blocklist, component retention, file analysis, rollback, and manual review are missing. |
| `05-quality-profiles-release-scoring.md` | Not started | Release candidates are sorted roughly by seeders and size. | Quality profiles, custom formats, scoring rules, cutoffs, allowed/preferred/rejected rules, size limits, upgrade logic, score breakdowns, and independent video/audio/subtitle requirements are missing. |
| `06-track-sourcing-sync-muxing.md` | Not started | Media tool detection and River job infrastructure exist. Docker bundles ffmpeg, mediainfo, and mkvtoolnix. | Component profiles, stream extraction, sync confidence, mux jobs, provenance, manual review, artifact retention, and assembly status UI are missing. |
| `07-subtitles-captions-localization.md` | Not started | No dedicated subtitle implementation yet. | OpenSubtitles, wanted subtitle profiles, embedded subtitle detection, subtitle scoring, download, conversion, sync, muxing, external subtitle preservation, and subtitle provenance are missing. |
| `08-books-music-management.md` | Deferred | Deferred by product decision. | Books, audiobooks, artists, albums, tracks, ebook/audio import, and music tagging are intentionally out of V0. |
| `09-admin-security-configuration.md` | Partial | Local admin auth, session cookie, settings pages, Docker Compose, Docker image, dev reset, configurable root folders/providers, and tool detection exist. | Persistent DB-backed sessions, credential encryption/redaction hardening, settings for profiles/tool paths, health dashboard, queues dashboard, backups, restore, API keys, webhooks, retention, and schedule controls are missing. |
| `10-open-questions-backlog.md` | Partial | Core stack and V0 defaults have been resolved and mostly applied: PostgreSQL, chi, pgx, River, OpenAPI, direct browser API calls, Transmission, SABnzbd, Torznab/Newznab/RSS, Docker media tools. | Remaining decisions include hardlink defaults, remote path mappings, scoring model details, assembly tolerances, subtitle provider details, canonical metadata precedence, naming rules, public API/plugin/script surface, and mobile/public exposure. |
| `11-engineering-architecture-standards.md` | Partial | Repo layout, `cmd`/`internal`, SvelteKit in `web`, OpenAPI generation, Makefile, Dockerfile, Docker Compose, River, pgx, chi, and focused unit tests exist. | `sqlc` is not yet in use, CI is not present, integration/e2e coverage is limited, readiness checks are minimal, and deeper redaction/path-safety/tool-exec test coverage is missing. |
| `12-guided-delivery-process.md` | Partial | Phase 1 decisions are complete. Phase 6 architecture choices are mostly implemented. Phase 7 skeleton is mostly complete. Phase 8 has started with add/search/grab. | Phase 2 workflow definitions, Phase 3 domain model, Phase 4 state machines, Phase 5 scoring rules, and the rest of Phase 8 import/stream analysis are incomplete. |
| `13-mvp-scope-and-milestones.md` | Partial | Milestone 1 is mostly complete. Parts of Milestones 4 and 5 are started: indexer config/search, client config, grab initiation. | Milestone 2 metadata/libraries, Milestone 3 profiles/scoring, Milestone 5 monitoring/import/streams, Milestone 6 subtitles, Milestone 7 assembly, and Milestone 8 hardening remain. |
| `14-architecture-decisions.md` | Partial | Most accepted decisions are reflected: separate dev servers, Go-served static frontend, local admin auth, PostgreSQL, reset-dev, SSE endpoint, OpenAPI contract-first, chi, pgx, River, Docker-bundled tools, repo layout, Makefile. | `sqlc` is still missing, persistent sessions are not wired to the DB-backed `app.sessions` table, and the SvelteKit static adapter decision should be closed. |
| `15-openapi-contract-workflow.md` | Partial | OpenAPI contract exists, Go server/types are generated, TypeScript schema is generated, browser API client uses `openapi-fetch`, and Make targets generate/check API artifacts. | The contract has expanded beyond the starter surface, but API error coverage and drift enforcement could be stricter. OpenAPI-generated TypeScript currently generates schema types, not a fully generated operation wrapper layer. |

## Current MVP Boundary

The current app is an early MVP foundation, not a complete media manager MVP.

The strongest implemented path is:

1. Configure indexers and download clients.
2. Add a movie or series as monitored.
3. Run a background release search.
4. Review stored release candidates in the detail page.
5. Enqueue a release grab.
6. See activity move to grabbed or failed.

The missing path is:

1. Add a real movie or series with provider IDs and library/root folder.
3. Apply a profile and score releases.
4. Monitor the download until completion.
5. Import the completed file.
6. Analyze streams.
7. Show the imported file and stream state.
8. Decide if upgrades, subtitles, or component assembly are still needed.

## Recommended Build Order

### 1. Close The Video Vertical Slice

This should be the next milestone because it turns the current search/grab foundation into a usable MVP loop.

- Expand metadata beyond the current TMDB/TVDB search cache to refresh jobs, poster/backdrop handling, runtime, and aliases.
- Add the minimal movie import model: imported file, path, size, container, runtime, video/audio/subtitle streams.
- Add download polling jobs for Transmission and SABnzbd.
- Add remote path mapping support before import.
- Implement completed-download import with a conservative largest-video-file rule.
- Run `ffprobe` or MediaInfo after import and persist stream metadata.
- Show imported file and streams on the media detail page.

Exit criteria:

- A configured indexer and client can find, grab, monitor, import, analyze, and display one movie.

### 2. Add Minimal Profiles And Scoring

This should come before expanding TV/anime too far, because search and import decisions need profile semantics.

- Add quality profile CRUD.
- Add a minimal release parser for resolution, source, codec, edition, language terms, and release group.
- Add scoring rules with allowed, rejected, preferred, cutoff, and size bounds.
- Store score breakdowns on release candidates.
- Show rejection reasons and score details in interactive search.
- Use the top accepted candidate for automatic search.

Exit criteria:

- A user can express a basic 1080p WEB-DL profile, reject CAM/TS, and understand why each release was accepted or rejected.

### 3. Add TV And Anime Foundations

TV should become real after movie import is working.

- Add series, season, and episode tables.
- Add TMDB/TVDB series lookup.
- Add season/episode monitored state.
- Add episode release matching.
- Add anime metadata fields and provider mapping investigation hooks.
- Add anime-specific aliases and numbering strategy fields.

Exit criteria:

- A user can add a series, monitor seasons/episodes, search for one episode, grab it, import it, and see streams.

### 4. Add Subtitle Workflow

Subtitle handling is lower risk than audio sync and unlocks part of the independent-component vision.

- Add subtitle wanted-profile fields.
- Detect embedded subtitle streams during import.
- Add OpenSubtitles provider settings.
- Search and download missing subtitles.
- Store external subtitle files and metadata.
- Mux or keep external subtitles according to profile.

Exit criteria:

- A movie imported without Italian subtitles can be detected as missing subtitles, source a matching subtitle, and expose the result in the detail view.

### 5. Add Constrained Track Assembly

Track assembly should stay conservative until enough stream metadata, profiles, and provenance are reliable.

- Add component profile requirements for video, audio, and subtitles.
- Add component-source release searches.
- Retain downloaded component files separately from final imports.
- Extract selected streams from compatible MKV sources.
- Start with exact or near-exact runtime matching.
- Mux final MKV with `mkvmerge`.
- Store per-stream provenance.
- Require manual review for low-confidence or different-cut cases.

Exit criteria:

- For one movie, the app can retain a base video, source a missing subtitle or audio stream from a compatible component source, mux a final MKV, and show provenance.

### 6. Harden Operations

- Add CI for Go tests, Svelte checks, lint, format, and OpenAPI generation drift.
- Add queue/job dashboard and cancellation/retry controls.
- Add backup and restore.
- Add better health checks for database, tools, indexers, and clients.
- Add path traversal protections and external-command tests.
- Add credential redaction and encrypted-at-rest secret storage.

## Immediate Next Sprint

The next sprint should avoid starting track assembly too early. The current blocker is the missing import and stream-analysis backbone.

Recommended tickets:

1. Add imported media file and media stream tables.
2. Add Transmission and SABnzbd polling jobs.
3. Add completed-download import job with path mapping and largest-media-file detection.
4. Add media tool execution wrapper for `ffprobe` with structured output and tests.
5. Show imported file and stream metadata on the detail page.
6. Add metadata refresh/background sync jobs for stored TMDB/TVDB IDs.
7. Add the first minimal quality profile schema after import works.
8. Add release parsing and score breakdowns for interactive search.
9. Add provider health/rate-limit tracking for metadata providers and indexers.
10. Add series/season/episode schema after movie import is working.

## Notes

- `sqlc` is an accepted architecture decision, but the current storage layer uses hand-written pgx queries. Adopt `sqlc` before the schema grows much further, otherwise later conversion will be more expensive.
- The `app.sessions` table exists, but the current HTTP session store is in-memory. Moving sessions to PostgreSQL should happen before treating auth as production-ready.
- Books and music should remain deferred until movie and TV import, profile scoring, and at least the subtitle part of component handling are proven.
