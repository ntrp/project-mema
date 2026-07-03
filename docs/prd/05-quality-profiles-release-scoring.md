# PRD: Quality Profiles and Release Scoring

Status: Draft

## Summary

Mema needs comprehensive quality and profile management comparable to Sonarr, but extended to independent video, audio, subtitle, source, edition, and domain-specific scoring.

## Goals

- Let users define exactly how releases are scored.
- Support quality profiles per media type and library.
- Support custom formats and preferred words.
- Support independent scoring for video, audio, subtitles, release source, release group, edition, codec, HDR, and size.
- Explain every score and rejection.
- Allow automatic upgrade decisions.

## Functional Requirements

- Users can create, clone, edit, and delete profiles.
- Profiles can be assigned globally, per library, per media item, or per collection/series.
- Profiles define allowed, preferred, and rejected qualities.
- Profiles define cutoff criteria.
- Profiles define size limits by runtime or fixed size.
- Profiles define custom formats with positive or negative scores.
- Profiles define required terms and forbidden terms.
- Profiles define preferred release groups.
- Profiles define edition preferences.
- Profiles define language preferences.
- Profiles define upgrade behavior.
- Profiles can be tested against a release title.
- Mema stores score breakdowns for every candidate.

## Profile Types

- Media quality profile
- Video component profile
- Audio component profile
- Subtitle component profile
- Release profile
- Custom format profile
- Naming profile
- Import profile

## Candidate Scoring Dimensions

- Resolution
- Source: BluRay, WEB-DL, WEBRip, HDTV, DVD, CAM, etc.
- Codec: H.264, H.265, AV1, VC-1, MPEG-2
- Bit depth
- HDR format: HDR10, HDR10+, Dolby Vision, HLG, SDR
- Audio codec: AAC, AC3, EAC3, TrueHD, DTS, DTS-HD MA, FLAC, Opus, MP3
- Audio channels
- Audio language
- Commentary tracks
- Subtitle language
- Subtitle type: full, forced, SDH, signs/songs, commentary
- Container
- Release group
- Repack/proper
- Edition
- Size
- Seeders/age
- Custom words

## Acceptance Criteria

- Users can create a profile that prefers 1080p WEB-DL over 720p but rejects CAM.
- Users can create custom format scores and see their effect.
- Users can set a cutoff quality.
- Users can define required German and English audio independently from video quality.
- Users can define wanted Italian subtitles independently from video and audio.
- Interactive search shows score breakdown and rejection reasons.
- Upgrade decisions are deterministic and explainable.

## Open Questions

- Should the scoring engine be rule-based, weighted, or both?
- Should custom formats follow Sonarr/Radarr semantics closely for familiarity?
- Should users be able to import/export profiles as YAML or JSON?
- Should profiles support inheritance?
- Should profile changes retroactively mark files as upgradeable?
- Should there be profile presets for movies, TV, anime, books, music, and audiobooks?
- Should audio/subtitle requirements be hard requirements or score-based preferences?
- How should Mema rank a release that has perfect video but missing target audio?
- Should component assembly allow lower-scored intermediate downloads if final output satisfies the profile?
- Should release title parsing and actual file analysis produce separate scores?
- Should users be able to create scripting hooks for scoring?
- Should scoring support regex, tokenized rules, or a safer expression language?

