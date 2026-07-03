# PRD: Media Library Management

Status: Draft

## Summary

Mema must first manage video libraries across movies, TV shows, and anime-oriented movie/TV content. The library model must support monitoring, search, import, upgrade, rename, deletion tracking, and file provenance. Books, audiobooks, music artists, albums, and tracks are later milestones.

## Goals

- Add, monitor, unmonitor, refresh, search, import, upgrade, and delete media.
- Support separate libraries with independent root folders and default profiles.
- Support multiple folder schemes by media type.
- Track wanted, missing, downloaded, imported, rejected, superseded, and assembled states.
- Store technical metadata for imported files and each media stream.
- Preserve source provenance for final assembled media.
- Support manual and automatic import.

## Media Domains

- Movies
- TV shows
- Animated movies
- Animated shows
- Books, later
- Audiobooks, later
- Music artists, later
- Albums, later
- Tracks, later

## Functional Requirements

- Users can create one or more libraries.
- Each library has a media type, root path, naming convention, default quality profile, default language/track profile, and import behavior.
- Users can add media manually by search.
- Users can bulk import existing folders.
- Users can scan existing files and match them to metadata entities.
- Users can mark media as monitored or unmonitored.
- Users can monitor entire shows, seasons, specials, individual episodes, authors, series, artists, albums, or tracks depending on domain.
- Users can see missing items and queued searches.
- Users can see imported files and stream-level details.
- Users can see why an item is eligible for upgrade.
- Users can repair metadata matches.
- Users can rename and reorganize files according to templates.
- Users can define recycling/trash behavior for replaced files.

## Data Concepts

- Library
- Media item
- Collection or series
- Season
- Episode
- Edition
- Release
- Download
- Imported file
- Stream
- Desired component
- Acquired component
- Final artifact
- Provenance record

## User Stories

- As a user, I can create a movie library with a root folder and default profiles.
- As a user, I can create a TV library and monitor selected seasons.
- As a user, I can configure anime-specific movie/TV metadata, season behavior, and dual-audio defaults.
- As a user, I can import an existing movie folder and have Mema match files to metadata.
- As a user, I can inspect the video, audio, and subtitle streams inside a file.
- As a user, I can see whether a file satisfies my target profile.
- As a user, I can manually override a metadata match.
- As a user, I can rename files without changing their content.
- As a user, I can preview file moves before applying them.

## Acceptance Criteria

- Existing files can be scanned and represented in the library with stream details.
- Every managed item has a clear monitored state.
- Every imported file records where it came from.
- Replacing or upgrading a file retains history.
- The UI distinguishes missing media, downloaded media, imported media, and assembled media.
- Manual import supports selecting media item, edition, language, quality, and whether to copy, hardlink, or move.

## Open Questions

- What are the exact top-level media types in the UI?
- Which anime metadata providers should be required from the start?
- How should specials, OVAs, extras, deleted scenes, and bonus features be modeled?
- Should editions be first-class for movies: theatrical, director's cut, extended, remaster, IMAX, commentary?
- Should a single media item support multiple final files for different profiles?
- Should multiple users be able to maintain separate wanted lists?
- Should library paths support remote mounts, SMB/NFS paths, and Windows paths?
- Should Mema detect moved/deleted files automatically through scheduled scans?
- Should hardlinks be preferred by default for torrents?
- Should imports be atomic with rollback when post-processing fails?
- How much manual import UX is required before first release?
