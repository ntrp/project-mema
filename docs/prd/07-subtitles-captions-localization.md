# PRD: Subtitles, Captions, and Localization

Status: Draft

## Summary

Mema must provide Bazarr-like subtitle management, extended into the component-profile model. Subtitles can be embedded, external, downloaded, synchronized, converted, and muxed.

## Goals

- Search and download subtitles by media item and language.
- Manage embedded and external subtitles.
- Distinguish subtitle types.
- Synchronize subtitles to video.
- Integrate subtitles into final muxed files when configured.
- Preserve external subtitle files when configured.

## Subtitle Types

- Full
- Forced
- SDH/CC
- Signs and songs
- Commentary
- Lyrics, for music videos later

## Functional Requirements

- Users can define wanted subtitle languages per profile.
- Users can define subtitle type preferences.
- Users can prefer embedded or external subtitles.
- Mema can detect embedded subtitle streams.
- Mema can search subtitle providers.
- Mema can download subtitle files.
- Mema can score subtitle candidates.
- Mema can synchronize subtitle files to video.
- Mema can convert subtitle formats where safe.
- Mema can mux subtitles into final output.
- Mema can keep external `.srt`, `.ass`, `.ssa`, `.vtt`, or `.sub` files where configured.
- Mema can mark subtitles as hearing impaired, forced, default, or commentary.

## Provider Decisions

- OpenSubtitles is the first external subtitle provider.
- Embedded subtitle detection is required from the first subtitle workflow.

## Acceptance Criteria

- A user can define Italian subtitles as wanted for animated movies.
- Mema can detect that a file lacks Italian subtitles.
- Mema can search configured subtitle providers.
- Mema can download and sync an `.srt` subtitle to the selected video.
- Mema can mux the subtitle with correct language metadata.
- The UI distinguishes missing, found, synced, muxed, and failed subtitle states.

## Open Questions

- Which subtitle providers should be added after OpenSubtitles?
- Should subtitle provider credentials be stored encrypted?
- Should `.ass` anime subtitles be preserved without conversion?
- Should OCR of image subtitles be supported?
- Should Mema support subtitle translation?
- Should subtitle sync be automatic by default?
- Should users be able to prefer forced subtitles over full subtitles?
- Should subtitles be stored externally, muxed, or both?
- Should Mema support multiple subtitle files per language?
- Should bad subtitle reports affect future scoring?
