# MVP Scope and Milestones

Status: Draft

Codename `Mema` remains temporary until the final project name is chosen.

## Summary

Version 0 is video first. It covers movies, TV shows, anime-specific movie/TV handling, subtitles, indexers, torrent/NZB download clients, quality profiles, independent component profiles, and a constrained track-assembly workflow.

Books and music are product goals, but they are not part of version 0.

This is a full vanilla project. It does not integrate with existing ARR-suite applications as upstream systems. It implements the required capabilities directly while integrating with indexers, torrent clients, NZB clients, metadata providers, subtitle providers, and media-processing tools.

## In Scope For Version 0

- Movie libraries.
- TV libraries.
- Anime-specific behavior layered onto movies and TV.
- Anime metadata providers, season formats, title aliases, episode numbering, and profile presets.
- Subtitle discovery and management.
- Native indexer integration.
- Transmission torrent client integration.
- SABnzbd NZB client integration.
- Download monitoring.
- Import pipeline.
- Media stream analysis.
- Quality profiles.
- Release scoring.
- Independent video, audio, and subtitle component profiles.
- Track assembly as a constrained MVP prototype.
- MKV-focused final assembly.
- Manual review for low-confidence sync or compatibility cases.
- Provenance tracking for imported and assembled files.

## Out Of Scope For Version 0

- Books.
- Audiobooks.
- Music artists, albums, and tracks.
- ARR-suite application integration.
- Prowlarr as an upstream dependency.
- Plex/Jellyfin/Emby/Kodi replacement playback.
- DRM circumvention.
- AI translation or dubbing.
- Automatic handling of different cuts or editions in track assembly.
- Broad container support beyond MKV-focused assembly.

## Anime Model

Anime is not a separate top-level media domain in version 0.

Anime support is implemented through:

- Movie entities for anime movies.
- TV entities for anime series.
- Anime-specific metadata providers.
- Anime-specific title and alias matching.
- Anime-specific season and episode numbering behavior.
- Anime-specific subtitle handling, including signs/songs where possible.
- Anime-specific scoring/profile presets, especially for dual-audio and subtitles.

## Track Assembly MVP

The track-assembly MVP should be conservative.

Initial constraints:

- Prefer MKV output.
- Start with exact or near-exact runtime matching.
- Use manual review for low-confidence cases.
- Store source provenance per final stream.
- Reject or hold jobs where source components likely come from different cuts, censored versions, or incompatible editions.
- Treat automatic audio synchronization as a design decision still requiring threshold definition.

Example target:

- Video: 1080p
- Audio: German and English
- Subtitles: Italian
- Result: select a base video, detect missing target components, source compatible audio/subtitle components, sync where confidence allows, mux final MKV, and record provenance.

## Milestones

### Milestone 1: Skeleton

- Go backend starts.
- SvelteKit frontend starts.
- Health endpoint works.
- Basic configuration exists.
- Docker Compose boots the application.
- PostgreSQL is available through Docker Compose.
- Development reset flow can reinitialize app data/schema while migrations are deferred.
- Local admin authentication exists.
- Go API can serve built frontend assets for production.
- SvelteKit dev server can run separately during development.
- No SvelteKit server runs in production.
- SvelteKit browser code calls the Go API directly.
- OpenAPI spec generation exists from the first backend route.
- River is wired for background jobs.
- Repo uses a top-level Go module and `web/` for SvelteKit.
- Makefile exposes common local commands.
- Docker image bundles required media tools.
- Local runs detect locally installed media tools.
- CI quality gates run.

### Milestone 2: Library And Metadata

- Create movie and TV libraries.
- Add movie and TV items.
- Support anime-specific metadata choices for movie/TV items.
- Store external IDs and aliases.
- Scan existing files.

### Milestone 3: Profiles And Scoring

- Create quality profiles.
- Create release scoring rules.
- Create independent video/audio/subtitle component profiles.
- Show score breakdowns and rejection reasons.

### Milestone 4: Indexers And Search

- Configure Torznab, Newznab, and RSS indexers.
- Run interactive search.
- Run automatic search.
- Normalize releases.
- Score candidates.

### Milestone 5: Download And Import

- Configure Transmission and SABnzbd.
- Send releases to clients.
- Monitor downloads.
- Import completed files.
- Analyze media streams.
- Show imported streams in the UI.

### Milestone 6: Subtitle Workflow

- Detect missing subtitles.
- Search subtitle providers.
- Download subtitles.
- Sync and mux or store external subtitles based on profile.

### Milestone 7: Track Assembly Prototype

- Detect missing audio/subtitle components.
- Search and retain component sources.
- Extract selected streams.
- Mux final MKV.
- Store provenance.
- Require manual review for uncertain compatibility.

### Milestone 8: Video Hardening

- Upgrade decisions.
- Failed download handling.
- Job retries and cancellation.
- Path mappings.
- Better anime matching.
- Import rollback.
- Backup and restore.

### Milestone 9: Later Domains

- Books.
- Audiobooks.
- Music.

## Remaining Phase 1 Questions

- Should SvelteKit use adapter-static or a minimal custom static build strategy for Go asset serving?

## Resolved Integration Defaults

- First torrent client: Transmission.
- First NZB client: SABnzbd.
- First indexer protocols: Torznab, Newznab, and RSS.
- First movie metadata provider: TMDB.
- First TV metadata providers: TVDB and TMDB.
- First anime metadata approach: AniList plus AniDB mapping investigation.
- First subtitle provider: OpenSubtitles.
- Embedded subtitle detection is required.
- Database: PostgreSQL from the start.
- Migration tooling is deferred until the project is live.
- Development reset system reinitializes app data/schema while migrations are deferred.
- Development runs Go API and SvelteKit dev server separately.
- Production uses one Go app serving built frontend assets plus PostgreSQL.
- Authentication is mandatory in version 0 with a local admin user and session cookie.
- Realtime updates start with Server-Sent Events.
- Backend stack: chi, pgx, and sqlc.
- API contract: OpenAPI from the start.
- API workflow: update OpenAPI contract first for API behavior changes, generate Go/TypeScript artifacts, then update backend implementation and web usage.
- OpenAPI toolchain: `oapi-codegen`, `openapi-typescript`, and `openapi-fetch`.
- Job management: River.
- Repo layout: top-level Go module with SvelteKit in `web/`.
- Local task runner: Makefile.
- Dev reset command: `make db-reset` calling gated `<binary> reset-dev`.
- Media tools: bundled in Docker, installed locally for local runs.
- SvelteKit API access: direct browser-to-Go API calls through generated OpenAPI client.
