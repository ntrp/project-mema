# PRD: Media Profile Desired State

## Status

Implemented

Last reconciled with the application: **2026-07-14**.

## Implementation Snapshot

- [x] One profile editor covers video, audio, subtitle, container, custom-format,
      upgrade, and delivery policy.
- [x] Profiles require at least one quality and one audio-language target.
- [x] Subtitle targets remain optional.
- [x] Video targets support codecs, HDR formats, and pixel formats with required
      flags and scores.
- [x] Audio targets support language, codec, channel layouts, minimum/preferred
      bitrate, scores, retention, and explicit lossy-conversion policy.
- [x] Subtitle targets support language, formats, scores, placement, retention,
      and release fallback.
- [x] Final container, preferred delivery protocol, and series-pack preference
      are profile-level policy.
- [x] Target evaluation, rollup, wanted rows, and fulfillment are implemented by
      the closed target-satisfaction epic
      ([#122](https://github.com/ntrp/project-mema/issues/122)).
- [ ] Bitmap-to-text subtitle conversion remains unsupported; PGS/SUP tracks can
      be recognized and matched but not converted to text formats.
- [ ] A generic profile rule engine remains intentionally deferred.

Primary implementation: `internal/storage/media_profile_*`,
`internal/targets/`, `internal/satisfaction/`, `internal/jobs/fulfillment_*`,
`internal/httpapi/media_profile_*`, and
`web/src/lib/components/settings/profiles/`.

## Context

Media profiles now combine quality selection, target audio languages, subtitle
languages, custom-format score settings, video requirements, final container,
and delivery policies in one profile contract and editor.

The goal is to make a profile describe the final state the user wants for a
media item. Target satisfaction, media rollup, release decisions, fulfillment
jobs, wanted-table behavior, and post-import source-of-truth rules are specified
separately in `docs/prd/target-satisfaction-and-media-state-rollup.md`.

## Problem

Users should not need to think in separate terms like quality profiles,
language targets, container requirements, and component targets. They should be
able to answer one question per profile:

What should my media look like when it is finished?

The system should translate that desired state into a stable profile contract
that other subsystems can evaluate.

## Goals

- Provide one unified profile UI for desired video, audio, and subtitle state.
- Keep existing quality profile behavior as the video quality baseline.
- Model video, audio, and subtitle targets as first-class profile requirements.
- Require every profile to define a video target.
- Require every profile to define at least one audio language target.
- Keep subtitle targets optional.
- Make lossy transformations explicit and opt-in.
- Keep the model extensible for future target fields and policy settings.

## Non-Goals

- Replace the existing release quality ranking system immediately.
- Automatically delete existing tracks without a separate retention/removal
  decision.
- Require every profile to define subtitle targets.
- Require every profile to define every possible codec, bitrate, perk, or
  format target.
- Build a full rule engine before the desired-state model is stable.

## User Experience

Profiles should expose one desired-state editor with three sections:

1. Video
2. Audio
3. Subtitles

Each section should focus on outcome targets, not implementation details. Where
the implementation choice matters, the UI should expose it as policy fields.
Runtime behavior for those policies is specified in
`docs/prd/target-satisfaction-and-media-state-rollup.md`.

## Video Requirements

Every profile must define a video target. The minimum valid video target is the
existing quality selection.

### Quality

The profile must keep the existing quality list and quality upgrade behavior.
Example: `WEBDL-1080p`.

Acceptance criteria:

- A profile cannot be saved without a video target.
- At least one quality target is required until an explicit alternative video
  target model exists.
- Existing upgrade-until-quality behavior continues to work.
- The desired-state UI presents quality as the first video target.

### Codec Target

The profile should allow an optional video codec target.

Examples:

- `h264`
- `h265`
- `av1`

Acceptance criteria:

- A profile can define no codec target.
- A profile can define one preferred codec target.
- The profile can specify whether codec mismatch is strict, preferred, or
  transformable.

### Video Perks

The profile should support optional video feature targets.

Initial target examples:

- HDR mode, such as HDR10, HDR10+, Dolby Vision, SDR
- pixel format, such as `yuv420p`, `yuv420p10le`
- bit depth
- resolution details when quality alone is not enough

Acceptance criteria:

- Video perks are optional.
- Each perk can be scored or required independently.
- Unknown or unsupported perks should not block profile saving.

## Audio Requirements

The profile should define a list of desired audio languages. Each language is a
target row with its own requirements and scoring.

At least one audio language target is mandatory for every profile.

For every audio language, the profile should support:

- language
- required flag
- release scoring value
- target codec
- target channels
- target bitrate
- lossy transcode policy
- fulfillment policy

### Audio Language

Audio language replaces the current mental model of generic target languages in
the UI. The refactor should move the profile model directly to audio language
targets instead of keeping generic target languages as a parallel legacy model.

Acceptance criteria:

- A profile cannot be saved without at least one audio language target.
- Each language can be independently required or optional.

### Audio Codec

Each audio language may define a target codec.

Examples:

- `aac`
- `ac3`
- `eac3`
- `dts`
- `truehd`
- `flac`

Acceptance criteria:

- Codec can be omitted per language.
- Codec mismatch can be scored, strict, or transformable.

### Audio Channels

Each audio language may define target channels.

Examples:

- `2.0`
- `5.1`
- `7.1`

Acceptance criteria:

- Target channels can be omitted per language.
- Channel mismatch can be scored, strict, or transformable.

### Audio Bitrate

Each audio language may define a target bitrate.

Acceptance criteria:

- Bitrate can be omitted per language.
- Bitrate target should support minimum and preferred values.

### Lossy Transcoding Policy

The profile must explicitly control whether lossy audio transcoding is allowed.

Acceptance criteria:

- Lossless-to-lossy conversion can be allowed per target.
- Lossy-to-lossy conversion must be separately allowed.
- Lossy-to-lossy should default to disabled.
- The UI must explain the chosen policy through labels, not browser tooltips.

## Subtitle Requirements

The profile should define a list of desired subtitle languages. Each language is
a target row with its own requirements and format target.

Subtitle targets are optional. A profile with no subtitle targets should not
search for subtitles, require subtitle tracks, or mark subtitle status as
missing.

For every subtitle language, the profile should support:

- language
- required flag
- release scoring value
- subtitle source preference
- target subtitle format
- conversion policy

### Subtitle Language

Subtitle languages continue to be separate from audio languages.

Acceptance criteria:

- A profile can define zero or more subtitle language targets.
- A profile with zero subtitle targets is valid.
- Each subtitle language can be independently required or optional.
- Existing `any`, `embedded`, and `external` behavior remains supported.

### Subtitle Format

Each subtitle language may define a target format.

Examples:

- `subrip` (SubRip, `.srt` sidecar files)
- `ass`
- `vtt`
- `pgs`

Acceptance criteria:

- Format can be omitted per language.
- Format mismatch can be marked transformable when conversion is available.

## Final Container Requirements

The profile must define the final container the media should be stored in.

Examples:

- `mkv`
- `mp4`

Container selection is separate from the video, audio, and subtitle targets.

Acceptance criteria:

- A profile cannot be saved without a final container target.
- The final container target is shown in the video section or a dedicated
  container section of the desired-state editor.
- The model should leave room for container-level options such as chapter
  retention, attachment retention, and metadata retention.

## Implemented Data Model

The application uses typed profile fields and target records rather than a
generic rule table:

- quality IDs remain top-level and define the video quality baseline;
- final container and delivery/conversion/retention settings are profile-level;
- one typed video target stores codec, HDR, and pixel-format preferences;
- audio and subtitle targets are typed lists keyed by language;
- custom formats remain a separate score and upgrade layer.

The profile API exposes the implemented desired-state shape to the UI:

```json
{
  "finalContainer": "mkv",
  "qualityIds": ["webdl-1080p"],
  "videoTarget": {
    "codecs": ["h265"],
    "codecRequired": false,
    "codecScore": 20,
    "hdrFormats": ["hdr10"],
    "hdrRequired": false,
    "hdrScore": 10,
    "pixelFormats": []
  },
  "audioTargets": [
    {
      "languageId": "english",
      "score": 100,
      "targetCodec": "eac3",
      "targetChannels": ["5.1"],
      "minimumBitrateKbps": 384,
      "preferredBitrateKbps": 768
    }
  ],
  "audioLossyTranscodePolicy": "disabled",
  "subtitleTargets": [
    {
      "languageId": "english",
      "score": 25,
      "formats": ["subrip"]
    }
  ],
  "subtitleMode": "external"
}
```

## Delivered Implementation

- [x] Final typed profile model and initial-schema storage.
- [x] Unified editor and OpenAPI-generated frontend contract.
- [x] Required quality and audio target validation with optional subtitles.
- [x] Final-container and delivery-policy selection.
- [x] Release decisions, target evaluation, media rollup, wanted rows, and
      fulfillment jobs.

## Resolved Design Decisions

- Quality remains a top-level ordered profile concept and is the resolution
  baseline for the video target.
- Audio and subtitle targets remain separate typed language lists; component
  source/provenance workflows complement them rather than replacing them.
- Final container remains top-level profile state.
- Codec, HDR, and pixel format are first-class video fields; custom formats
  remain the extensible scoring mechanism.
- Codec, subtitle-format, and channel values are normalized strings validated by
  the backend and constrained by frontend catalogues.
- Audio targets support both minimum and preferred bitrate.
- Retained alternate sources and assembled outputs use component-source and
  provenance records instead of adding source selection to the profile target.

## Success Criteria

- [x] A user can define the desired final media state from one profile screen.
- [x] The profile defines the final container format for stored media.
- [x] Audio-language requirements are first-class profile targets and participate
      in release decisions, satisfaction, wanted rows, and fulfillment.
- [x] Existing quality/custom-format upgrade behavior remains integrated with
      the desired-state model.
