# Open Questions Backlog

Status: Draft

This file collects cross-cutting questions that should be answered before technical design. The questions are grouped by decision area.

## Product Scope

Resolved:

- MVP scope is video first: movies, TV, anime-specific movie/TV handling, subtitles, indexers, download clients, profiles, and constrained track assembly.
- Books and music are deferred until after the video foundation.
- Anime is modeled through movie and TV entities with anime-specific metadata providers, season formats, matching, subtitles, and profile presets.
- ARR-suite application integration is out of scope. This is a vanilla project, not a wrapper around existing ARR applications.

Open:

1. Is mobile management required for MVP?
2. Should multi-user roles be added after the local-admin version?
3. Should the app be safe to expose publicly, or assume reverse proxy/VPN only?

## Deployment

1. What is the primary target beyond Docker Compose: bare-metal binary, NAS, Kubernetes, or all later?
2. What platforms must be supported: Linux amd64, Linux arm64, macOS, Windows?

Resolved:

- PostgreSQL is required from the start.
- No migration tooling until the project is live.
- During development, the Go API and SvelteKit dev server run as separate processes.
- In production, Go serves the built frontend assets.
- Deployment target for the skeleton is one app container plus PostgreSQL.
- Early development uses a resettable app/schema initialization system instead of migrations.
- Required media tools are bundled into the Docker image.
- Local runs use locally installed host tools.

## Indexers

Resolved:

- MVP indexer protocols are Torznab, Newznab, and RSS.

Open:

1. Should Mema expose an indexer proxy API for external tools later?
2. Should indexers be tagged and filtered per library?
3. How should private tracker rules influence scoring?
4. Should Mema support manual cookie-based indexers?

## Download Clients

Resolved:

- First torrent client is Transmission.
- First NZB client is SABnzbd.

Open:

1. Should failed-download handling be automatic?
2. How should remote path mappings work?
3. Should seeding rules be managed?
4. Should hardlink be the default for torrents?
5. Should source components be retained after final assembly?

## Quality And Scoring

1. Should scoring mirror Sonarr/Radarr custom formats or define a new model?
2. Should profile rules be GUI-only or import/exportable as YAML/JSON?
3. Should profiles support inheritance?
4. Should component requirements be hard constraints or weighted preferences?
5. Should scoring happen at release-title time, file-analysis time, or both?
6. Should low-scored component downloads be allowed if the final assembled file satisfies the target profile?

## Track Assembly

Resolved:

- Track assembly is in MVP as a constrained prototype.
- MVP starts with MKV-focused output, exact or near-exact runtime matching first, and manual review for low-confidence cases.

Open:

1. Is automatic audio synchronization allowed for high-confidence MVP cases, or should every audio sync require manual review?
2. What exact tolerance defines near-exact runtime matching?
3. Which final container formats are supported: MKV only, MP4 later?
4. Which external tools are acceptable hard dependencies?
5. Should Mema support PAL speedup/slowdown correction?
6. Should Mema support different cuts or editions?
7. Should manual review include playback and offset controls?
8. How much confidence is needed before automatic muxing?
9. How should provenance be displayed to the user?
10. Should final muxing be reproducible from a saved project manifest?

## Subtitles

1. Which subtitle providers are mandatory?
2. Should `.ass` subtitles be preserved for anime?
3. Should Mema support OCR for image subtitles?
4. Should Mema download forced subtitles separately from full subtitles?
5. Should subtitles be stored externally, muxed, or both?
6. Should Mema support subtitle translation?

## Metadata

Resolved:

- Movies start with TMDB.
- TV starts with TVDB and TMDB.
- Anime starts with AniList plus AniDB mapping investigation.

Open:

1. Which provider is canonical per media type after first integrations are proven?
2. How should provider priority be configured?
3. Should local NFO files be read and written?
4. Should anime use AniDB, AniList, TVDB, or a mapping layer long-term?
5. Should music use MusicBrainz release groups or specific releases?
6. Should books track works and editions separately?

## Library And Files

1. What folder structures should be supported by default?
2. Should users be able to define custom naming templates with conditionals?
3. Should imports be atomic?
4. Should replaced files go to trash?
5. Should Mema detect manual file changes through scheduled scans?
6. Should multiple final artifacts be allowed for the same media item?

## UI

1. What are the primary screens for MVP?
2. Should queues and history be global or per library?
3. Should profile editing be wizard-based, table-based, or both?
4. Should interactive search support component-level selection?
5. Should assembly jobs have a detailed timeline view?
6. Should advanced users be able to edit raw profile JSON?

Resolved:

- Use Server-Sent Events first for one-way realtime status updates.

## API And Extensibility

1. Should Mema have a public REST API from day one?
2. Should there be a plugin system?
3. Should custom scripts run before import, after import, before mux, and after mux?
4. Should webhooks be compatible with Sonarr/Radarr event payloads?
5. Should SvelteKit use adapter-static or a minimal custom static build strategy for Go asset serving?

Resolved:

- Use chi for HTTP routing.
- Use pgx for PostgreSQL access.
- Use sqlc for SQL-first data access generation.
- Authentication is mandatory in version 0 with a local admin user and session cookie.
- Use OpenAPI from the start.
- Use `oapi-codegen`, `openapi-typescript`, and `openapi-fetch`.
- Use contract-first OpenAPI: API behavior changes update the spec first, generate Go/TypeScript artifacts next, then update backend and web implementation.
- Use River for background job management.
- SvelteKit browser code calls the Go API directly.
- Do not run a SvelteKit server in production.
