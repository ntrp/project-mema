# Guided Delivery Process

Status: Draft

This document is the working process for taking the project from idea to implementation. It is intentionally phased so product decisions, technical design, and code setup happen in the right order.

Codename `Mema` remains temporary until the final project name is chosen.

## Phase 1: Product Boundaries

Goal: decide what version 0 is and what it is not.

Decisions:

- MVP media domains.
- Whether anime is a separate domain or a profile/category for movies and TV.
- Whether books and music are first-release features or later milestones.
- Whether the application integrates with existing ARR-suite apps or stays fully independent.
- Whether track assembly is an MVP feature or starts as a prototype behind manual review.

Resolved decisions:

- Version 0 is video first.
- Version 0 includes movies, TV, anime-specific movie/TV metadata and season behavior, subtitles, indexers, download clients, quality profiles, and constrained track assembly.
- Anime is modeled through movie and TV entities, not a separate top-level domain.
- Books and music are later milestones.
- The project is fully independent from ARR-suite applications. It still integrates with torrent clients, NZB clients, and indexers.

Outputs:

- Updated product vision PRD.
- MVP scope document.
- Non-goals for version 0.
- Initial milestone list.

Exit criteria:

- We can describe version 0 in one paragraph.
- Every major product capability is either in scope, out of scope, or explicitly deferred.

## Phase 2: User Workflows

Goal: define the critical user journeys before data models or APIs.

Workflows:

- Add a movie.
- Add a TV show.
- Add anime with dual-audio/subtitle expectations.
- Configure an indexer.
- Configure a download client.
- Create a quality profile.
- Create independent video/audio/subtitle component profiles.
- Run interactive search.
- Grab a release.
- Import a completed download.
- Detect missing tracks.
- Source missing tracks.
- Sync and mux final output.
- Review failed or low-confidence assembly.

Outputs:

- Workflow PRD sections with acceptance criteria.
- First UI route map.
- Manual review rules.

Exit criteria:

- Each critical workflow has clear start state, user action, system behavior, success state, and failure state.

## Phase 3: Domain Model

Goal: define durable concepts before database implementation.

Core concepts:

- Library.
- Media item.
- Collection, season, episode, edition, author, artist, album, track.
- Release.
- Download.
- Imported file.
- Media stream.
- Desired component.
- Candidate component.
- Assembly job.
- Final artifact.
- Profile.
- Score breakdown.
- Provenance record.

Outputs:

- Data model design.
- Entity lifecycle diagrams.
- State machines for downloads, imports, and assembly.

Exit criteria:

- We can represent a normal ARR-style import and a multi-source assembled file without special cases.

## Phase 4: Scoring And Profiles

Goal: define the scoring engine well enough to implement deterministic behavior.

Decisions:

- Hard requirements vs weighted preferences.
- Profile inheritance.
- Custom format compatibility with Sonarr/Radarr mental models.
- Score breakdown model.
- Cutoff and upgrade behavior.
- Separate release-title scoring from post-download file-analysis scoring.

Outputs:

- Scoring engine PRD.
- Profile schema draft.
- Example profiles.
- Test matrix for accepted/rejected releases.

Exit criteria:

- Given a release and profile, the system can produce a deterministic score and rejection list.

## Phase 5: Track Assembly Design

Goal: reduce the riskiest feature into staged, testable capabilities.

Decisions:

- MKV-only initial output or multiple containers.
- Required external tools.
- Automatic sync threshold.
- Manual review workflow.
- Provenance retention.
- Source component retention.
- Supported mismatch cases: exact runtime only, offsets, drift, PAL speedup, different cuts.

Outputs:

- Assembly pipeline design.
- Confidence model.
- Job state machine.
- Tool capability matrix.
- Manual review requirements.

Exit criteria:

- We know what the first assembly prototype will and will not attempt.

## Phase 6: Technical Architecture

Goal: choose implementation shape before scaffolding.

Decisions:

- Go API and SvelteKit run as separate processes during development.
- Go serves built frontend assets in production.
- No SvelteKit server in production.
- SvelteKit browser code calls the Go API directly.
- Database: PostgreSQL from day one.
- Go router and persistence stack: chi, pgx, and sqlc.
- OpenAPI contract-first from the start.
- Background job persistence through River.
- Realtime updates: Server-Sent Events first.
- Docker layout.
- Authentication: local admin user and session cookie from version 0.
- Resettable app/schema initialization before migration tooling exists.
- Repository layout: top-level Go module with SvelteKit under `web/`.
- Local commands through Makefile.
- Media tools bundled in Docker and installed locally for local runs.
- OpenAPI toolchain: `oapi-codegen`, `openapi-typescript`, and `openapi-fetch`.

Outputs:

- Architecture document.
- API design direction.
- Repository layout.
- Development environment plan.

Exit criteria:

- We can scaffold code without expecting immediate structural rework.

## Phase 7: Scaffold

Goal: create a minimal working application skeleton.

Work:

- Initialize Go module.
- Create backend entry point.
- Add health endpoint.
- Add configuration loading.
- Add structured logging.
- Add PostgreSQL connection and schema initialization. Do not add migration tooling until the project is live.
- Add a development-only reset path for app data/schema changes.
- Add `make db-reset` that calls the gated `<binary> reset-dev` flow.
- Add River job management and its required job-table setup as part of bootstrap/reset.
- Add OpenAPI spec, Go server/types generation, and TypeScript web client generation.
- Add Docker image setup for bundled media tools.
- Add system health checks for local/bundled media tool availability.
- Add local admin authentication and session cookie plumbing.
- Add SSE endpoint foundation for job/search/queue progress.
- Create SvelteKit TypeScript frontend.
- Add ESLint, Prettier, Vitest, Playwright, and `sv check`.
- Add Docker Compose.
- Add CI quality gates.

Exit criteria:

- Backend starts.
- Frontend starts.
- Health endpoint works.
- Tests and checks run.
- Docker Compose can boot the skeleton.
- Development reset flow is documented and protected from production use.
- OpenAPI generation is part of the development workflow.
- API behavior changes start with the OpenAPI contract.
- River can enqueue and process a smoke-test job.
- Makefile targets exist for `dev`, `dev-api`, `dev-web`, `db-reset`, `test`, and `check`.
- Docker image includes required media tools.
- Local startup reports missing host media tools clearly.

## Phase 8: Vertical Slice 1

Goal: build the first useful end-to-end path.

Recommended slice:

- Create a library.
- Configure one indexer protocol.
- Configure one download client.
- Add one movie.
- Run interactive search.
- Score releases.
- Send one release to download client.
- Poll completion.
- Import file.
- Analyze streams with ffprobe.
- Show imported file and streams in UI.

Exit criteria:

- One media item can move from wanted to imported with visible scoring and stream metadata.

## Phase 9: Vertical Slice 2

Goal: prove the differentiating track-profile workflow.

Recommended slice:

- Define required audio/subtitle components.
- Detect missing components.
- Search for candidate component sources.
- Download or manually register component source.
- Extract selected stream.
- Mux into MKV.
- Store provenance.
- Show assembly job result.

Exit criteria:

- One final MKV can be assembled from base video plus at least one external subtitle or audio track.

## Phase 10: Expand Domains

Goal: add broader media-domain coverage after the core video engine works.

Order:

1. Movies.
2. TV.
3. Anime-specific matching, metadata providers, season behavior, and profile presets.
4. Subtitles.
5. Books.
6. Music.

Exit criteria:

- Each domain has metadata, wanted state, search, download, import, profiles, and UI support.

## Phase 11: Hardening

Goal: make the system reliable for real self-hosted use.

Work:

- Backup and restore.
- Authentication.
- Secrets handling.
- Health checks.
- Job retries and cancellation.
- Better path mapping.
- Import rollback.
- External tool diagnostics.
- E2E tests for critical flows.
- Migration tests.
- Performance testing with large libraries.

Exit criteria:

- The application can survive restarts, failed jobs, bad downloads, and config errors without corrupting library state.

## Phase 12: Release Preparation

Goal: prepare for a usable public alpha.

Work:

- Final project name.
- Branding cleanup.
- Docker image naming.
- Documentation.
- Example profiles.
- Example Docker Compose.
- Migration/import documentation for existing media libraries.
- Known limitations.
- Security review.

Exit criteria:

- A technical user can install, configure, and test the application from documentation.
