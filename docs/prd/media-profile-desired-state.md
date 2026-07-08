# PRD: Media Profile Desired State

## Status

Draft

## Context

Media profiles currently combine quality selection, target audio languages,
subtitle languages, custom format score settings, and experimental component
targets. Those settings are stored on the profile, but the UI still presents
them as separate rule groups instead of one desired media state.

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
- `atmos`

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

- `srt`
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

## Data Model Direction

The refactored profile model should make desired video, audio, subtitle, and
container targets the canonical shape. Existing profile fields can be removed or
reshaped as part of the refactor because the project is not released yet.

Suggested target concepts:

- `profile_video_targets`
- `profile_audio_targets`
- `profile_subtitle_targets`
- `profile_container_targets`

or one generalized target table with typed detail fields. The final schema
choice should be made during implementation after evaluating sqlc query shape
and API ergonomics.

The profile API should expose a desired-state object that is easy for the UI to
render:

```json
{
  "container": {
    "format": "mkv",
    "required": true
  },
  "video": {
    "qualityIds": ["webdl-1080p"],
    "codec": {
      "value": "h265",
      "required": false,
      "score": 20
    },
    "perks": [
      {
        "kind": "hdr",
        "value": "hdr10",
        "required": false,
        "score": 10
      }
    ]
  },
  "audio": [
    {
      "languageId": "english",
      "required": true,
      "score": 100,
      "codec": "eac3",
      "channels": "5.1",
      "minimumBitrateKbps": 384,
      "preferredBitrateKbps": 768,
      "lossyTranscode": "disabled"
    }
  ],
  "subtitles": [
    {
      "languageId": "english",
      "required": true,
      "score": 25,
      "source": "external",
      "format": "srt"
    }
  ]
}
```

## Implementation Plan

Phase 1: PRD and final model design

- Document current profile fields and runtime consumers only as input to the
  refactor.
- Decide the final target schema for video, audio, subtitles, and container.
- Remove old profile fields unless they are part of the new model.

Phase 2: Unified UI and new profile contract

- Replace the separate profile rule sections with one desired-state editor.
- Write the new desired-state profile contract directly.
- Replace target language UI with audio language targets.
- Add final container target selection.

Phase 3: Backend desired-state model

- Add or reshape storage for video/audio/subtitle/container targets.
- Update OpenAPI and generated frontend types.
- Enforce required video target and at least one audio language target during
  profile validation.
- Enforce required final container target during profile validation.
- Preserve zero-subtitle-target profiles as valid.

Phase 4: Runtime integration

- Follow `docs/prd/target-satisfaction-and-media-state-rollup.md` for target
  evaluation, release decisions, media rollup, wanted-table behavior, and
  fulfillment jobs.

## Open Questions

- Should quality stay as a top-level profile concept or become a video target?
- Should component targets be generalized enough to replace audio/subtitle
  language tables?
- Should final container be top-level profile state or a video/container target?
- Which video perks should be first-class versus custom key/value targets?
- Which codecs and subtitle formats should be normalized enums versus free text?
- Which audio channel targets should be normalized enums versus free text?
- Should bitrate targets be minimum-only, preferred-only, or both?
- How should retained component sources be selected when a release has the
  desired audio or subtitle track but the imported main file does not?

## Success Criteria

- A user can define the desired final media state from one profile screen.
- The profile defines the final container format for stored media.
- Audio language settings no longer feel disconnected from component targets.
- The implementation can be delivered incrementally without breaking existing
  profile behavior.
