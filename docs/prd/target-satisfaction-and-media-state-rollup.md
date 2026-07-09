# PRD: Target Satisfaction And Media State Rollup

## Status

Draft

## Context

Mema currently has several overlapping ways to describe whether media is done:
release decisions decide what to fetch, file detail views inspect probed tracks,
subtitle state reports whether configured subtitle languages exist, and media
items expose a coarse `missing`, `downloading`, or `downloaded` status.

This PRD defines one system for evaluating desired video, audio, and subtitle
targets after files are imported or scanned. It also defines how target state
rolls up into media state, what work can move media toward done, and how wanted
and file detail UI should represent missing, partial, pending, and unwanted
tracks or subtitles.

## Goals

- Define release decision steps before download.
- Define target satisfaction for video, audio, and subtitle targets.
- Use only database-backed facts for post-import satisfaction.
- Separate target states from track and subtitle visual states.
- Define media file and media item rollup rules.
- Define fulfillment jobs that can move media toward done.
- Define wanted table behavior for missing media, missing targets, and custom
  format upgrades.
- Make settings changes immediately affect satisfaction, pending operations,
  unwanted indicators, and wanted rows.

## Non-Goals

- Implement the refactor in this PRD.
- Define final database schema.
- Require transcoding, remuxing, subtitle conversion, or component sourcing jobs
  to exist before target satisfaction can be calculated.
- Treat metadata, artwork, chapters, or arbitrary sidecars as done-state
  targets.

## Core Model

### Target

A target is a wanted requirement from the selected media profile.

Target types:

- **Video target**: quality, video codec, HDR format, pixel format, resolution,
  container/file attributes, and release/provenance-only attributes when
  relevant.
- **Audio target**: language, codec, channels, minimum bitrate, and preferred
  bitrate.
- **Subtitle target**: language, format, and placement requirement from the
  active subtitle mode.

One target can evaluate many persisted candidates. No target configured means no
target row exists for that scope. Targets never become `unwanted`; only actual
tracks or subtitle records can be unwanted.

### Track Or Candidate

A candidate is an actual database-known item that can be checked against a
target.

Candidate types:

- Persisted probed video tracks.
- Persisted probed audio tracks.
- Persisted probed embedded subtitle tracks.
- Stored external subtitle records.
- Persisted file/provenance facts used by video or custom format evaluation.

Manual files or subtitles do not affect satisfaction until a manual or automatic
rescan discovers them and persists updated file, track, or subtitle data.

## Source Of Truth

Release parsing is used to choose what to fetch. After a file is downloaded,
imported, or scanned, satisfaction uses database-backed facts from the file and
its provenance.

Probeable fields use probed file data as source of truth:

- video codec
- HDR format
- pixel format and bit depth
- resolution
- audio codec
- audio channels
- audio bitrate
- embedded subtitle tracks
- container
- duration
- chapters

Release/provenance-only fields remain usable after import when they cannot be
reliably probed:

- release group
- release title or release name
- release hash
- indexer or provider identity
- release id
- uploader or indexer flags
- protocol
- publish stats
- proper, repack, real, or version marker
- edition when it is not probeable from the file

If a field is both parsed from release and later probeable from the file, the
probed value wins for target satisfaction.

## Initial Release Decision

Release selection predicts whether an indexer result is worth grabbing. It does
not prove final satisfaction.

Initial release decision steps:

1. Parse release title and filename into media identity, year, season/episode,
   quality, source, resolution, video/audio hints, language hints, release group,
   release flags, and release type.
2. Match the parsed release to the requested movie or series.
3. Reject unsafe candidates:
   - title, year, season, or episode does not match
   - quality is unknown or not enabled in the profile
   - required target mismatch is known from reliable release metadata
   - audio bitrate is below a hard minimum when parseable
   - release language set conflicts with remove-unwanted-audio settings when
     reliable
   - custom format score is below the profile minimum
4. Score surviving candidates:
   - quality rank
   - custom format score
   - audio and subtitle language target score
   - video, audio, and subtitle target hints when parseable
   - preferred protocol
   - series pack preference
   - seeders
   - publish time
   - size
5. Choose the best acceptable candidate. Warnings can win when no hard rejection
   applies.

The decision engine must avoid false rejection when release metadata is not
reliable enough to prove a required target mismatch. Unknown target state should
be resolved after import and probing.

## Target Satisfaction

Target satisfaction is recalculated from persisted file, track, subtitle, and
provenance facts whenever the relevant media item, file, profile, or settings
change.

### Target States

| State | Meaning |
| --- | --- |
| `missing` | No database-backed candidate can satisfy the target. |
| `partial` | At least one related candidate exists, but every related candidate fails one or more requirements. |
| `pending` | A candidate exists and can satisfy the target after a known operation. |
| `satisfied` | At least one candidate fully satisfies the target. |
| `upgradeable` | The current state is usable, but profile upgrade rules want a better candidate. |
| `blocked` | The target cannot be satisfied without user action, settings change, or missing tool support. |
| `failed` | The last fulfillment attempt for this target failed. |

`available` and `unmanaged` are not target states. If no target exists, no target
row is created.

### Video Targets

Video targets evaluate persisted file and video track data first. Release or
provenance data is used only for attributes that cannot be reliably probed.

Video target satisfaction:

- `satisfied`: at least one persisted video track plus file/provenance facts
  matches required quality, codec, HDR, pixel format, resolution, container, and
  other configured video fields.
- `partial`: a video track exists and is related to the target, but fails one or
  more required fields.
- `missing`: no file record exists, no video track record exists, quality is
  unknown, or no video track/file fact can satisfy required fields.
- `pending`: a known operation such as remux or video transcode can satisfy the
  target and is required by current settings.
- `upgradeable`: the file is usable, but quality, custom format, or video target
  upgrade rules prefer a better release or transform.
- `blocked`: the configured video target needs unsupported tooling or impossible
  conversion.
- `failed`: the last video fulfillment job failed.

### Audio Targets

Each audio target is evaluated against all persisted audio tracks for the media
file.

Audio target satisfaction:

- `satisfied`: at least one audio track matches language, codec, channel target,
  and bitrate requirements.
- `partial`: one or more tracks match the target language, but none fully satisfy
  codec, channels, or bitrate.
- `missing`: no file record exists, no audio track record exists, or no audio
  track matches the target language.
- `pending`: a known audio operation can satisfy the target, such as audio
  transcoding or sourcing a track from another release.
- `upgradeable`: a usable track exists, but profile upgrade rules prefer a
  better audio candidate.
- `blocked`: target cannot be fulfilled because configured transcoding or source
  policies do not allow it, required tools are missing, or no known source can
  provide it.
- `failed`: the last audio fulfillment job failed.

### Subtitle Targets

Each subtitle target is evaluated against persisted embedded subtitle tracks and
stored external subtitle records.

Subtitle mode changes evaluation:

- `embedded`: only embedded subtitle tracks satisfy the target.
- `external`: stored external subtitle records satisfy the target.
- `mixed`: embedded subtitle tracks or stored external subtitle records satisfy
  the target.

Subtitle target satisfaction:

- `satisfied`: at least one candidate matches target language, format when
  required, and active subtitle mode.
- `partial`: a candidate exists for the language, but fails format or another
  target requirement.
- `pending`: a candidate exists in the wrong placement for the active mode, such
  as an external subtitle that must be embedded.
- `missing`: no candidate exists for the target language in any allowed location.
- `upgradeable`: a usable subtitle exists, but profile rules prefer a better
  format or source.
- `blocked`: target cannot be fulfilled because provider, conversion, extraction,
  or embedding support is unavailable under current settings.
- `failed`: the last subtitle fulfillment job failed.

## Settings-Driven Recalculation

Settings are part of target satisfaction. Changing a profile or media item
setting must recalculate target states, candidate visual states, pending
operations, and wanted rows from the same persisted facts.

Examples:

- Subtitle mode changes from `embedded` to `mixed`: an external subtitle that was
  `pending` becomes `satisfied`.
- Subtitle mode changes from `mixed` to `embedded`: an external subtitle becomes
  `pending` until embedded.
- Remove-unwanted-audio enabled: audio tracks outside enabled audio target
  languages are visually unwanted.
- Remove-unwanted-subtitles enabled: subtitle tracks or external subtitle records
  outside enabled subtitle target languages are visually unwanted.
- Remove-unwanted setting disabled: extra tracks remain visible but are no longer
  unwanted.
- Profile target removed: the target row disappears, and related missing or
  pending wanted rows disappear.

## Candidate Visual States

Track and subtitle rows explain how actual candidates relate to targets.

| Visual state | Meaning |
| --- | --- |
| `matching` | Candidate satisfies at least one target. |
| `partial` | Candidate relates to a target but fails one or more requirements. |
| `unwanted` | Candidate conflicts with profile or settings. |
| `pending_operation` | Candidate can satisfy a target after a known operation. |
| `missing_placeholder` | Synthetic row shown because a target has no candidate. |

Visual rules:

- A track can be `matching` for one target and still be unwanted for a different
  policy only when UI can explain both roles without ambiguity.
- Missing placeholders are target rows rendered in the track table, not real
  tracks.
- Pending subtitle rows must name the required operation, such as embed external
  subtitle, extract embedded subtitle, or convert subtitle format.
- Unwanted rows must come from settings or subtitle mode, not from target state.

## Media Rollup

Media file state rolls up from required video, audio, and subtitle target states.

| Rollup | Condition |
| --- | --- |
| `missing` | No usable file exists, or required file-level target state is missing with no usable media. |
| `downloading` | Download, import, transform, extraction, source, or merge work is queued or running. |
| `partial` | A usable file exists, but at least one required target is `missing`, `partial`, `pending`, `blocked`, or `failed`. |
| `downloaded` | A usable file exists and every required target is `satisfied` or acceptable under the profile. |
| `upgradeable` | The media is usable and done enough to watch, but profile upgrade rules prefer a better release, custom format score, or target candidate. |

Media item state rolls up from expected media files:

- Movies require one media file.
- Series require monitored episodes according to monitor mode.
- A media item is `downloaded` only when every required file is downloaded under
  the media file rollup.
- A media item is `partial` when at least one required file is usable but not
  fully satisfied.
- A media item is `missing` when required files do not exist or no usable file
  exists for required media.
- A media item is `downloading` when active work exists for required media.
- A media item is `upgradeable` when all required files are usable but one or
  more upgrade rules remain unmet.

## Custom Format Handling

Custom formats are profile scoring and upgrade rules by default. They are not
video, audio, or subtitle targets unless a future profile model explicitly ties a
custom format to one of those target types.

Before download:

- Custom formats match parsed release data.
- Matched custom formats contribute profile score.
- Minimum custom format score can reject a release.
- Custom format score affects upgrade choice.

At grab or import:

- Matched custom formats and total custom format score should be persisted with
  release provenance.
- Persisted facts should include enough detail to explain which formats matched
  and why.

After import:

- Custom format upgrade state uses persisted provenance and custom-format facts.
- Reparse current filenames only when provenance/custom-format facts are missing.
- Probeable custom-format specs must not override probed target satisfaction.
- A custom format upgrade target can make media `upgradeable`.
- Custom-format upgrade rows can appear in wanted independently from
  video/audio/subtitle target rows.

## Fulfillment Jobs

Fulfillment jobs move persisted state toward `satisfied` or `upgradeable`.
Jobs must write durable facts, then satisfaction recalculates from the database.
Automatic jobs must be configurable so users can choose which operations Mema
handles. Each automatic job type must expose at least an enabled flag and a
schedule/interval when it runs periodically. Users may disable any automatic
operation when they prefer to handle that work manually.

Possible jobs:

- **Release search/grab/import**: finds and imports a release for missing media or
  upgrades.
- **Video upgrade search**: searches for a better release when video/custom
  format upgrade rules remain unmet.
- **Video transcoding**: transforms video codec, resolution, HDR, pixel format, or
  related fields when policy and tooling allow.
- **Audio transcoding**: transforms audio codec, channels, or bitrate when policy
  allows.
- **Audio sourcing from another release**: fetches another release to extract or
  merge a desired audio track.
- **Container remuxing**: moves selected streams into the target container.
- **Subtitle download**: fetches stored external subtitle records for missing
  subtitle targets.
- **Subtitle merge/embed**: embeds an external subtitle when mode requires
  embedded subtitles.
- **Subtitle extraction**: extracts embedded subtitles when external subtitles are
  required or preferred.
- **Subtitle transformation/conversion**: converts subtitle format when target
  format requires it and tooling supports it.

V1 implements only these eight operations. Each operation has its own River
worker and fixed system schedule, inserted disabled by default until an
administrator enables periodic execution.

Job states should influence target state:

- queued/running job for a target can make target `pending` or media
  `downloading`
- failed job can make target `failed`
- missing tool, disabled policy, or impossible operation can make target
  `blocked`

While running, jobs must report precise progress and structured logs. Progress
should identify the media item, file or target being processed, current phase,
percent or unit progress when known, started time, last update time, and any
pending external tool command or provider request. Logs should be visible from
the system job history and from media-specific activity surfaces.

Every operation that an automatic job can perform must also be available as a
manual action from the relevant media row, track row, or external subtitle row.
Manual actions enqueue the operation-specific worker with media/file/track or
subtitle scope, persist the same facts, emit the same progress/log events, and
recalculate satisfaction with the same rules as automatic jobs.

## Wanted Table

Wanted views must expose both media rows and target rows.

Media rows:

- Show when no usable media file exists for a movie or monitored episode.
- Include media title and season/episode context when applicable.
- Drive normal release search/grab workflows.
- Use row kind `media` so existing no-file wanted cases remain distinct from
  file-level target repair work.

Target rows:

- Show when an existing media file has a required target in `missing`,
  `partial`, `pending`, `blocked`, or `failed`.
- Include parent context: media title, season/episode when applicable, file label
  or path, target type, target language when applicable, target state, and
  required operation when pending.
- Recalculate immediately when profile settings change.
- Use row kind `target`; satisfied, upgradeable, and removed profile targets do
  not produce target wanted rows.

Custom-format upgrade rows:

- Show when all required media is usable but the custom format upgrade target is
  unmet.
- Display separately from video, audio, and subtitle target rows.
- Include current score, target score, and parent media/file context.
- Use row kind `custom_format_upgrade`; these rows are driven by custom-format
  score delta instead of video, audio, or subtitle target state.

Wanted table rows should not be produced for targets that no longer exist after
profile changes.

## Acceptance Criteria

- PRD clearly distinguishes target state from track/subtitle visual state.
- PRD states that post-import satisfaction is DB-only and probe data wins for
  probeable fields.
- PRD covers initial release decision, target satisfaction, media rollup, custom
  format handling, fulfillment jobs, visual states, and wanted table behavior.
- PRD includes setting-dependent examples, including subtitle mode changing
  external subtitles from pending to satisfied.
- PRD does not introduce metadata, artwork, chapter, or arbitrary sidecar targets.

## Implementation Notes

- The existing `docs/prd/media-profile-desired-state.md` remains the broader
  profile model PRD. This PRD is the target satisfaction and rollup contract.
- The architecture docs should later link to this PRD when implementation starts.
- Because the project is unreleased, future schema changes should update the
  initial schema directly rather than adding migrations.
