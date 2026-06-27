# PRD: Product Vision and Scope

Status: Draft

## Summary

Codename `Mema` refers to the project temporarily. The final product name is undecided.

Mema is a self-hosted media manager built with Go and SvelteKit. It combines library management, metadata enrichment, indexer aggregation, download-client orchestration, quality/profile scoring, subtitle management, and post-processing. The first version is video-first: movies, TV shows, anime-specific movie/TV handling, subtitles, indexers, download clients, profiles, and constrained track assembly.

Mema should be a full vanilla project. It should not integrate with Sonarr, Radarr, Readarr, Lidarr, Prowlarr, or Bazarr as upstream applications. It should implement equivalent product capabilities directly while integrating with torrent clients, NZB clients, indexers, metadata providers, subtitle providers, and media-processing tools.

Books and music remain long-term product goals, but they are deferred until the video engine is proven.

## Phase 1 Decisions

- Version 0 focuses on video first.
- Version 0 includes movies, TV, anime-specific movie/TV metadata and season behavior, subtitles, indexers, download clients, profiles, and track assembly.
- Anime is not a separate top-level domain in version 0. Anime uses movie and TV entities with anime-specific metadata providers, title matching, season formats, episode numbering, subtitle behavior, and profile presets.
- Track assembly is in MVP as a constrained prototype: MKV-focused output, exact or near-exact runtime matching first, and manual review for low-confidence sync cases.
- ARR-suite app integration is out of scope. The project does not rely on or connect to existing ARR applications.

## Goals

- Manage movies, TV shows, and anime-specific movie/TV variants in version 0.
- Support books, audiobooks, ebooks, music artists, albums, and tracks in later milestones.
- Search indexers and trackers through a native indexer integration layer.
- Integrate with multiple download clients.
- Provide Sonarr-grade quality profiles, release profiles, preferred words, custom formats, and scoring.
- Provide Bazarr-like subtitle discovery and management.
- Provide independent video, audio, and subtitle target profiles.
- Source video, audio, and subtitles from multiple releases when needed.
- Synchronize external audio/subtitle tracks to the selected video.
- Mux final files using deterministic, inspectable post-processing.
- Preserve provenance for every imported file and assembled track.
- Offer a modern SvelteKit UI with clear queues, decisions, conflicts, and manual overrides.
- Expose a documented local API for automation.

## Non-Goals For Initial MVP

- Hosting or streaming media playback as a Plex/Jellyfin replacement.
- Building a torrent or Usenet client.
- Circumventing DRM or decrypting protected commercial media.
- Cloud-hosted multi-tenant SaaS.
- Automated account creation for private trackers.
- AI-based translation or dubbing as a core first version.
- Integration with existing ARR-suite applications as upstream dependencies.
- Books and music in version 0.

## Product Principles

- One library model, multiple media domains.
- Transparent automation: every search, grab, rejection, import, and mux decision must be explainable.
- Profiles are composable: quality, language, source, release group, codec, edition, audio, subtitle, and custom rules should not be collapsed into a single hard-coded score.
- Manual override is first-class: users can inspect candidates and choose components.
- Automation must be conservative when synchronization confidence is low.
- Self-hosted reliability matters more than novelty.

## Core Personas

- Home media archivist: wants precise control over naming, quality, languages, editions, and folder layout.
- Multilingual household admin: wants different language and subtitle combinations per library, profile, series, or movie.
- Anime collector: wants release-group control, dual audio, subtitles, signs/songs handling, specials, and absolute/season episode mapping.
- Music collector: wants album and track completeness, release type handling, and metadata correctness in a later milestone.
- Book collector: wants ebook and audiobook acquisition with author/series organization in a later milestone.

## Main Workflows

1. Add wanted media.
2. Choose or inherit profiles.
3. Discover candidates from indexers.
4. Score and reject releases.
5. Send selected releases to download clients.
6. Monitor download state.
7. Import completed files.
8. Analyze media streams and metadata.
9. Decide whether the file is complete or needs extra tracks.
10. Search for missing audio/subtitle/video components.
11. Synchronize and mux final output when needed.
12. Rename, move, and update library state.
13. Continue monitoring for upgrades.

## Success Metrics

- A user can manage movie, TV, and anime-oriented video libraries from one UI.
- A user can express Sonarr-like quality scoring without editing config files.
- A user can define a profile such as `1080p video + German audio + English audio + Italian subtitles`.
- Mema can explain why a release was accepted or rejected.
- Mema can import a completed download and identify all embedded streams.
- Mema can assemble a final video file from separately sourced video, audio, and subtitle tracks when the inputs are compatible.
- Mema can detect low-confidence sync cases and require manual review.

## Major Risks

- Track synchronization is technically difficult and may require human review.
- Release naming is inconsistent across media types and indexers.
- Anime absolute episode numbering and movie editions can create matching ambiguity.
- Music and book metadata quality differs substantially from movie and TV metadata and is deferred until after the video foundation.
- Legal expectations vary by region; Mema must stay source-agnostic and user-operated.
- Combining all domains can create a large configuration surface.

## Open Questions

- Should animated movies and shows have special default scoring/profile presets?
- Should books and audiobooks be handled by one domain or separate domains?
- Should Mema support both torrents and Usenet in MVP?
- Which download clients are mandatory for first release?
- Should there be a plugin system for indexers, download clients, metadata providers, and post-processing?
- Should external tools like ffmpeg, mkvmerge, MediaInfo, fpcalc, or autosubsync be hard dependencies or optional capabilities?
- What is the expected deployment target: Docker Compose, single binary, Kubernetes, NAS packages?
- Should user authentication be optional for LAN deployments?
- Is multi-user role-based access required?
- Is mobile UI a first-class requirement?
- Should Mema integrate with Plex, Jellyfin, Emby, Kodi, or only manage files?
- Should Mema actively trigger library scans in media servers?
