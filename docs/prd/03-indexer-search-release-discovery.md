# PRD: Indexer, Search, and Release Discovery

Status: Draft

## Summary

Mema must aggregate indexer and tracker search across media domains through a native integration layer. It should support automatic search, interactive search, RSS sync, indexer health, category mapping, and release normalization.

## Goals

- Support multiple indexers.
- Support multiple protocols.
- Normalize release data across indexers.
- Score releases against profiles.
- Explain rejection reasons.
- Route accepted releases to appropriate download clients.
- Support indexer tags and per-library indexer selection.

## Functional Requirements

- Users can add, test, edit, disable, and delete indexers.
- Users can assign indexers to media types and libraries.
- Users can configure indexer categories.
- Users can run interactive search for a media item.
- Users can run automatic search for missing or upgradeable items.
- Mema periodically checks RSS feeds for wanted releases.
- Mema deduplicates releases found on multiple indexers.
- Mema normalizes release title, size, seeders, leechers, grabs, publish date, age, protocol, source, group, edition, quality, language hints, and category.
- Mema shows why each release was accepted or rejected.
- Mema supports indexer-specific rate limits and health checks.

## Candidate Protocols

- Torznab, first release
- Newznab, first release
- RSS, first release
- Native APIs for high-value providers, later

## Release Decision Inputs

- Metadata match confidence
- Media type
- Quality profile
- Release profile
- Custom format score
- Language/track profile
- Size limits
- Source type
- Codec
- HDR/DV flags
- Audio hints
- Subtitle hints
- Edition
- Release group
- Repack/proper status
- Freeleech or seed requirements
- Download client availability

## Acceptance Criteria

- A user can configure at least one Torznab/Newznab indexer.
- A user can run interactive search and inspect scored results.
- Rejected results show concrete reasons.
- RSS sync can discover wanted media.
- Indexer failures surface in a health page.
- Search results preserve raw release data for debugging.

## Open Questions

- Torznab, Newznab, and RSS are required for MVP.
- Should Mema expose its own indexer proxy API for external tools later?
- Should indexers be global, per library, per media type, or all three?
- Should private tracker ratio/freeleech rules affect scoring?
- Should release-group preferences be global or profile-specific?
- Should interactive search allow selecting video/audio/subtitle components separately?
- How should Mema detect that two indexer results are the same release?
- Should searches run immediately, queued, or both?
- Should users be able to write custom parser rules?
- Should Mema expose Torznab/Newznab endpoints for external tools?
