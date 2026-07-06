# PRD: Media Profile Desired State

## Status

Draft

## Context

Media profiles currently combine quality selection, release scoring rules,
target audio languages, subtitle languages, custom formats, and experimental
component targets. Those settings are stored on the profile, but the UI still
presents them as separate rule groups instead of one desired media state.

The goal is to make a profile describe the final state the user wants for a
media item. Search, import, extraction, conversion, and future transcoding jobs
should then use that profile to either find a matching release or transform the
current file toward the desired state.

## Problem

Users should not need to think in separate terms like release scoring,
component targets, subtitle satisfaction, and conversion jobs. They should be
able to answer one question per profile:

What should my media look like when it is finished?

The system should translate that desired state into:

- initial release ranking
- continued upgrade search
- retained component source decisions
- subtitle search or conversion jobs when subtitle targets are configured
- audio/video transcode jobs where allowed
- final container selection and remux jobs where needed
- file status and missing/unwanted component indicators

## Goals

- Provide one unified profile UI for desired video, audio, and subtitle state.
- Keep existing quality profile behavior as the video quality baseline.
- Model video, audio, and subtitle targets as first-class profile requirements.
- Require every profile to define a video target.
- Require every profile to define at least one audio language target.
- Keep subtitle targets optional.
- Preserve release scoring as one use of profile targets, not the only use.
- Support both search-based fulfillment and job-based fulfillment.
- Make lossy transformations explicit and opt-in.
- Show current media file compliance against the selected profile.
- Keep the model extensible for future tools such as video transcoding, audio
  transcoding, subtitle conversion, remuxing, and source retention.

## Non-Goals

- Implement transcoding jobs in the first PR.
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
the implementation choice matters, the UI should expose it as fulfillment
policy.

Example policy labels:

- Find in release
- Convert when possible
- Prefer existing file
- Allow missing
- Strict

The media detail file summary should compare the selected profile to the actual
file and show satisfied, missing, transformable, and unwanted components.

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
- Existing quality profile ordering continues to rank releases.
- Existing upgrade-until-quality behavior continues to work.
- The desired-state UI presents quality as the first video target.

### Codec Target

The profile should allow an optional video codec target.

Examples:

- `h264`
- `h265`
- `av1`

The codec target can be fulfilled by:

- finding a release that already matches
- keeping search/upgrades active until a matching release appears
- transcoding the existing file through a job, when enabled

Acceptance criteria:

- A profile can define no codec target.
- A profile can define one preferred codec target.
- The profile can specify whether codec mismatch is strict, preferred, or
  transformable.
- Release scoring can favor matching codec targets.
- File status can show whether the current video codec is satisfied.

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
- File probing should map detected video metadata to these targets where
  available.

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
- Required audio languages must affect release acceptance and file
  satisfaction.
- Optional audio languages can contribute score without blocking acceptance.

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
- Transformable codec mismatch can enqueue a future audio transcode job.

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
- File probing should compare actual channel layout or channel count to the
  profile target when available.
- Release scoring may favor channel matches when release metadata provides
  enough information.
- Transformable channel mismatch can enqueue a future audio transcode job when
  the configured codec and lossy transcode policy allow it.

### Audio Bitrate

Each audio language may define a target bitrate.

Acceptance criteria:

- Bitrate can be omitted per language.
- Bitrate target should support minimum and preferred values.
- File probing should compare actual bitrate to profile target when available.
- Release scoring may favor bitrate matches when release metadata provides
  enough information.

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
- A profile with zero subtitle targets treats subtitle status as ignored.
- Each subtitle language can be independently required or optional.
- Existing `any`, `embedded`, and `external` behavior remains supported.
- Required external subtitles can continue to trigger subtitle search.

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
- Subtitle conversion jobs should be able to use the target format later.

## Final Container Requirements

The profile must define the final container the media should be stored in.

Examples:

- `mkv`
- `mp4`

Container selection is separate from the video, audio, and subtitle targets. A
release may satisfy every stream target but still require a remux job if the
container does not match the profile.

Acceptance criteria:

- A profile cannot be saved without a final container target.
- The final container target is shown in the video section or a dedicated
  container section of the desired-state editor.
- File evaluation reports whether the current media container satisfies the
  profile.
- Container mismatch can be marked transformable when remuxing is available.
- Remuxing must preserve all selected final components and their provenance.
- The model should leave room for container-level options such as chapter
  retention, attachment retention, and metadata retention.

## Component Provenance Requirements

Every component that makes up the final media must retain provenance.

Components include:

- final container
- video stream
- each audio stream
- each subtitle stream
- chapters
- attachments or fonts when retained
- sidecar subtitle files when subtitles are external

For every component, the system must store:

- source release group
- source release name
- source release id when available
- source indexer or provider when available
- source file path or retained source id
- source stream id when the component came from a media stream
- transformation chain when the component was extracted, converted,
  transcoded, or remuxed

If a component is generated from another component, it must keep the original
release provenance and record the transformation that produced the current
component.

Acceptance criteria:

- No imported, extracted, converted, transcoded, remuxed, or retained component
  can be persisted without provenance.
- Provenance survives remux, subtitle conversion, audio transcode, and video
  transcode jobs.
- The UI can show which release group and release name supplied every final
  component.
- Release id is nullable only when the upstream source did not provide a stable
  id.
- Manually added files must still get provenance with a manual source marker
  and user-provided or derived release name.

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
    "required": true,
    "fulfillment": "remux"
  },
  "video": {
    "qualityIds": ["webdl-1080p"],
    "codec": {
      "value": "h265",
      "required": false,
      "score": 20,
      "fulfillment": "searchOrTranscode"
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
      "lossyTranscode": "disabled",
      "fulfillment": "searchOrTranscode"
    }
  ],
  "subtitles": [
    {
      "languageId": "english",
      "required": true,
      "score": 25,
      "source": "external",
      "format": "srt",
      "fulfillment": "searchOrConvert"
    }
  ],
  "provenancePolicy": {
    "requireComponentProvenance": true
  }
}
```

## Fulfillment Policy

Every target should eventually be evaluated into one of these outcomes:

- satisfied
- missing
- mismatch
- transformable
- unsupported
- ignored

Fulfillment policy should decide what the system may do:

- `scoreOnly`: use for release ranking only
- `search`: keep searching or upgrading until matched
- `convert`: transform existing file when possible
- `searchOrConvert`: prefer search, then convert if allowed
- `remux`: remux into the target container when streams are otherwise
  acceptable
- `existing`: prefer existing retained source or current file
- `allowMissing`: do not block completion

Names can change during implementation, but the distinction must remain clear.

## Release Scoring

Release scoring should use profile targets consistently:

- Video quality keeps its current rank behavior.
- Video codec and perks can contribute score.
- Audio language score continues to contribute score.
- Audio codec, channels, and bitrate can contribute score when parseable.
- Configured subtitle language scores continue to contribute score.
- Configured subtitle formats can contribute score when parseable.
- Container format can contribute score when release metadata is reliable.
- Required targets can reject releases only when the release metadata is
  reliable enough to make that decision.

If metadata is not reliable, the system should avoid false rejection and mark
the target as unknown until file probing/import.

## File Evaluation

After import or scan, file probing should evaluate actual tracks against the
profile:

- video quality, codec, pixel format, HDR, resolution
- audio language, codec, channels, bitrate
- configured subtitle language, embedded/external state, format
- container format
- component provenance

The UI should show:

- what is satisfied
- what is missing
- what exists but is unwanted under removal settings
- what can be transformed by a configured job
- what cannot currently be evaluated
- which release group and release supplied each final component

## Implementation Plan

Phase 1: PRD and final model design

- Document current profile fields and runtime consumers only as input to the
  refactor.
- Decide the final target schema for video, audio, subtitles, and container.
- Decide the component provenance schema before implementing component
  extraction or remux behavior.
- Remove old profile fields unless they are part of the new model.

Phase 2: Unified UI and new profile contract

- Replace the separate profile rule sections with one desired-state editor.
- Write the new desired-state profile contract directly.
- Replace target language UI with audio language targets.
- Add final container target selection.

Phase 3: Backend desired-state model

- Add or reshape storage for video/audio/subtitle/container targets.
- Add component provenance storage.
- Update OpenAPI and generated frontend types.
- Enforce required video target and at least one audio language target during
  profile validation.
- Enforce required final container target during profile validation.
- Preserve zero-subtitle-target profiles as valid.

Phase 4: Evaluation and scoring

- Add profile target evaluation for probed files.
- Extend release scoring to video codec/perks, audio codec/channels/bitrate,
  and subtitle/container format where metadata supports it.
- Surface target outcomes in media detail.
- Surface component provenance in media detail.

Phase 5: Fulfillment jobs

- Add job planning for supported transformations.
- Start with subtitle format conversion or audio transcode before video
  transcode.
- Add remux planning for final container mismatches.
- Require explicit lossy conversion policies before enqueueing lossy jobs.
- Ensure every fulfillment job writes updated provenance with the original
  release lineage intact.

## Open Questions

- Should quality stay as a top-level profile concept or become a video target?
- Should component targets be generalized enough to replace audio/subtitle
  language tables?
- Should final container be top-level profile state or a video/container target?
- Which video perks should be first-class versus custom key/value targets?
- Which codecs and subtitle formats should be normalized enums versus free text?
- Which audio channel targets should be normalized enums versus free text?
- Should bitrate targets be minimum-only, preferred-only, or both?
- How should the UI distinguish required-for-download from required-final-state?
- Should search continue indefinitely for transformable mismatches, or should
  conversion become eligible after a time or quality threshold?
- How should retained component sources be selected when a release has the
  desired audio or subtitle track but the imported main file does not?
- What should count as a stable release id for torrents, usenet releases, and
  subtitle providers?

## Success Criteria

- A user can define the desired final media state from one profile screen.
- The same profile explains release choice, upgrade behavior, missing track
  status, and future transformation jobs.
- The profile defines the final container format for stored media.
- Every final media component retains source release group, release name, and
  release id when available.
- Audio language settings no longer feel disconnected from component targets.
- Subtitle language settings can drive both subtitle search and future format
  conversion.
- The implementation can be delivered incrementally without breaking existing
  profile behavior.
