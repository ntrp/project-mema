# Features PRD Index

This folder contains the product requirements for a self-hosted Go + SvelteKit media manager that combines the main capabilities of Sonarr, Radarr, Readarr, Lidarr, Prowlarr, and Bazarr.

Codename: Mema. The final product name is undecided, and any reference to `Mema` in these PRDs should be treated as a temporary placeholder.

## Product Intent

Mema manages movies, TV shows, animated movies and shows, books, music, indexers, download clients, subtitles, audio tracks, quality profiles, imports, renames, and post-processing from one application.

The distinguishing requirement is independent media-component targeting: users must be able to define target video quality, wanted audio tracks, and wanted subtitles separately, then allow Mema to source components from multiple releases and assemble a final file.

Example:

- Content type: animated movie
- Target video: 1080p
- Target audio: German and English
- Target subtitles: Italian
- Result: source the best matching video, audio, and subtitle components from one or more downloads, synchronize tracks, mux them into one final file, and retain provenance.

## PRDs

- [Product Vision and Scope](./00-product-vision-and-scope.md)
- [Media Library Management](./01-media-library-management.md)
- [Metadata and Identifiers](./02-metadata-and-identifiers.md)
- [Indexer, Search, and Release Discovery](./03-indexer-search-release-discovery.md)
- [Download Clients and Import Pipeline](./04-download-clients-import-pipeline.md)
- [Quality Profiles and Release Scoring](./05-quality-profiles-release-scoring.md)
- [Track Sourcing, Synchronization, and Muxing](./06-track-sourcing-sync-muxing.md)
- [Subtitles, Captions, and Localization](./07-subtitles-captions-localization.md)
- [Books and Music Management](./08-books-music-management.md)
- [Administration, Security, and Configuration](./09-admin-security-configuration.md)
- [Open Questions Backlog](./10-open-questions-backlog.md)
- [Engineering Architecture and Setup Standards](./11-engineering-architecture-standards.md)
- [Guided Delivery Process](./12-guided-delivery-process.md)
- [MVP Scope and Milestones](./13-mvp-scope-and-milestones.md)
- [Architecture Decisions](./14-architecture-decisions.md)
- [OpenAPI Contract Workflow](./15-openapi-contract-workflow.md)

## Requirement Status

All PRDs are currently `Draft`. They are meant to capture initial requirements and open questions before implementation planning begins.

Status values:

- `Draft`: Requirements are incomplete and need product decisions.
- `Ready for Design`: User stories and acceptance criteria are stable enough for technical design.
- `Ready for Build`: Data model, API, UX, and integration behavior are sufficiently specified.
- `Implemented`: Feature has shipped.
