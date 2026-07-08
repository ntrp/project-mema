---
title: Release Decisions And Track Management
description: How release scoring, upgrades, and media file satisfaction are calculated.
---

Release selection and file satisfaction are related but separate systems. Release
selection predicts whether an indexer result is worth grabbing from its release
name and provider metadata. File satisfaction evaluates an imported or existing
file from probed video, audio, subtitle, chapter, and sidecar data.

## Release Decision Flow

Release searches pass candidates through `internal/decisions.Engine`. Each
candidate is parsed from its release title, then evaluated against the requested
media item, selected media profile, custom formats, and language catalog.

The evaluator rejects candidates when they cannot safely satisfy the search:

- The parsed title does not match the requested movie or series.
- The parsed year does not match the requested movie or series year.
- The parsed season or episode does not match the requested resource.
- The parsed quality is unknown or not enabled in the profile.
- A required video target rejects the parsed video codec, HDR format, or pixel
  format.
- Parsed audio bitrate is below a configured minimum.
- The release language set contains a language that is not enabled while
  unwanted audio removal is enabled.
- The custom format score is below the profile minimum.

Non-rejected candidates receive a match severity and score contributors:

- Quality rank comes from profile quality order. Later qualities in the profile
  list score higher.
- Custom format score sums matched custom format weights from the profile.
- Language score adds target scores for matched audio and subtitle languages and
  adds penalties for missing target languages.
- Profile target score adds configured video and audio target scores when parsed
  release attributes match optional targets.

After matching, the chooser ranks surviving candidates by quality score, custom
format score, language score, preferred protocol, season-pack preference,
seeders, publish time, size, and title. Hard errors are never chosen. Warnings
can still win when they are the best acceptable candidate.

Search results persist matched custom format names and score so grab and import
workflows can keep using those acquisition facts. Custom formats remain scoring
and upgrade rules by default; they do not create video, audio, or subtitle target
rows.

## Upgrade Decisions

Upgrade checks compare a candidate against the current file state before the
candidate can replace existing media.

When no current file is present, a candidate can be selected without upgrade
comparison. When a current file exists, the app parses the current file path to
derive current quality and custom format score.

With no profile, only higher quality releases count as upgrades. With a profile,
the candidate is blocked when:

- Upgrades are disabled.
- The current file already reached the configured quality upgrade target.
- The current file already reached the configured custom format upgrade target.
- The candidate quality is below the current file quality.
- The candidate quality is equal and custom format score does not improve.
- The custom format score improvement is below the profile minimum increment.

Quality improvement can upgrade even when custom format score does not improve.
Equal quality requires a positive custom format improvement that meets the
minimum increment.

## File Satisfaction

Media detail pages build `MediaFileRow` values from API media item data and the
selected profile. Rows combine probed file metadata with profile expectations:

- `tracks`, `chapters`, and `otherFiles` come from file probing.
- `subtitleSatisfaction` comes from backend subtitle state.
- Expected audio targets, subtitle targets, unwanted-track flags, and upgrade
  settings come from the selected media profile.

The file row shows separate compact states for audio, subtitles, quality,
custom score, and status. These states are inspection aids; they do not rewrite
the backend release decision after import.

## Target Satisfaction State Matrix

The current media item status is coarse: `downloaded` means the media item has
at least one library scan row, or at least one completed download activity. It
does not currently require readable files or satisfied profile targets. Target
state should come from video, audio, and subtitle satisfaction instead.

Each profile target is evaluated only from data already available in the
database: stored file records, probed track records, and stored external subtitle
records. A target is satisfied when at least one database-backed candidate fully
matches the target. Manually created files become eligible only after a manual or
automatic rescan discovers them and persists the updated file, track, or subtitle
data. The media item state is then the rollup of every required video,
audio, and subtitle target for the file.

Parsed release data is used to decide whether to fetch a new file. After a file
is downloaded or imported, video and audio target evaluation uses probed data
from that file as the source of truth. Parsed release fields remain relevant only
for attributes that cannot be reliably probed, such as release group.

| Target type | Candidate inputs | Target satisfied when | Target partial when | Target missing when |
| --- | --- | --- | --- | --- |
| Video target | DB file record, DB probed video tracks, probed container/file fields, and stored non-probable release attributes such as release group. | At least one video track plus file-level data matches required codec, HDR, pixel format, resolution, quality, and other configured fields. | A video track exists and matches identity-level fields, but fails one or more required target fields. | No file record exists, no video track record exists, quality is unknown, or no probed/stored file value can match required fields. |
| Audio target | DB probed audio tracks for the file. | At least one audio track matches target language, codec, channels, and minimum bitrate. | One or more audio tracks match the target language, but none fully satisfy codec, channels, or bitrate. Unwanted audio can also make the audio target set partial when removal is enabled. | No file record exists, no audio track record exists, or no audio track matches the target language. |
| Subtitle target | DB embedded subtitle tracks, stored external subtitle records, and subtitle mode. | At least one candidate matches target language and the configured subtitle mode: embedded track for embedded mode, external subtitle record for external mode, either for mixed mode. | A subtitle candidate exists for the language but is in the wrong place for the mode, such as an external subtitle that still needs embedding. | No candidate exists for the target language in any allowed location. |

Target states aggregate upward:

| Level | State rule |
| --- | --- |
| Track/candidate | Persisted probe row, or persisted external subtitle row for subtitle targets. It can match, partially match, be unwanted, or be unusable for a specific target. |
| Target | One of `missing`, `partial`, `pending`, `satisfied`, `upgradeable`, `blocked`, or `failed`. `available` and `unmanaged` are not target states; if a profile requirement is not configured, no target row exists. |
| Media file | `downloaded` only when the file exists and every required target is `satisfied` or acceptable under the profile. `partial` when at least one required target is `missing`, `partial`, `pending`, `blocked`, or `failed` but the file is usable. `missing` when no usable file exists. |

Wanted views should expose both media and target rows. A missing media row shows
that no usable file exists for the movie or episode. A target row shows a
specific missing or partial video, audio, or subtitle target for an existing
media file. Target rows must include enough parent context to make ownership
clear, such as media title, season and episode when applicable, file path or file
label, target type, target language when applicable, and target state.

## Canonical Target And Candidate States

Use a normalized state list for each video, audio, and subtitle target. Keep
candidate visual states separate from target states, and keep media item status
as a rollup instead of a single overloaded download flag.

| Target state | Meaning | Counts as done |
| --- | --- | --- |
| `missing` | No database-backed candidate can satisfy the target. | No |
| `partial` | At least one related candidate exists, but every related candidate fails one or more requirements. | No |
| `pending` | A candidate exists and can satisfy the target after a known operation. | No |
| `satisfied` | At least one candidate fully satisfies the target. | Yes |
| `upgradeable` | Usable, but profile upgrade rules prefer a better candidate. | Yes |
| `blocked` | Cannot be satisfied without user action, settings change, or missing tool support. | No |
| `failed` | Last fulfillment attempt failed. Retry or manual action may be possible. | No |

| Candidate visual state | Meaning |
| --- | --- |
| `matching` | Candidate satisfies at least one target. |
| `partial` | Candidate relates to a target but fails one or more requirements. |
| `unwanted` | Candidate conflicts with profile or settings. |
| `pending_operation` | Candidate can satisfy a target after a known operation. |
| `missing_placeholder` | Synthetic row shown because a target has no candidate. |

The media item rollup should be derived from child states:

| Rollup | Condition |
| --- | --- |
| `missing` | No usable required media file is present. |
| `downloading` | Work for required media or targets is queued or running, and no higher-priority failed or blocked state should be shown. |
| `partial` | A usable file exists, but required video, audio, or subtitle target state is `missing`, `partial`, `pending`, `failed`, or `blocked`. |
| `downloaded` | Required media file exists and required child states are `satisfied` or acceptable under the profile. |
| `upgradeable` | Same as downloaded, but one or more required child states can improve within profile upgrade rules. |

## Video Satisfaction

Release-time video satisfaction is handled in the release decision engine. The
engine parses video codec, HDR format, and pixel format from the release name.
Configured video target fields can either add score or reject the release when
the field is marked required.

File-time video state is calculated by `internal/satisfaction` from persisted
media file facts, persisted video track rows, profile quality order, final
container, and video target settings. Live filesystem discovery does not
participate in satisfaction; imports and rescans must persist file facts first.

No persisted file or no persisted video track returns `missing`. A known
quality, codec, HDR, pixel format, or container mismatch returns `partial` with
the failed requirement names. A container mismatch can return `pending` when the
known operation is a remux. A satisfied file can still return `upgradeable` when
the profile quality upgrade target is higher than the persisted quality.

## Audio Satisfaction

Audio satisfaction is calculated by `internal/satisfaction` from persisted audio
track facts and profile audio targets.

If no file exists, audio is missing. If the file has no audio tracks, audio is
missing. Otherwise each target is matched against detected audio tracks by
language first, then target details:

- target codec
- target channel list
- minimum bitrate

If a target language is absent, the file is missing that audio target. If the
language exists but codec, channels, or bitrate do not match, the file is
partial. If at least one target matches and another target fails, the file is
partial. If every target matches, audio is satisfied.

When a profile has no explicit audio target list, legacy required audio
languages are used as simple language-only targets.

If unwanted audio removal is enabled, audio tracks whose language is outside the
profile target languages become `unwanted` candidates. This does not change the
audio target state directly; target state comes from the configured target's own
matching, partial, pending, or missing candidates.

Missing audio targets are inserted into the track list as red placeholder audio
rows next to detected audio tracks. This mirrors missing embedded subtitle rows
and makes the detailed track list the source of truth for why the compact audio
badge is missing or partial.

## Subtitle Satisfaction

Subtitle satisfaction is calculated from backend subtitle state plus the
profile subtitle mode.

The profile declares wanted subtitle languages and mode:

- Embedded means wanted subtitles should become internal tracks.
- External means wanted subtitles should remain as sidecar files.
- Mixed allows existing sidecars and downloaded subtitles while still showing
  embedded-track expectations where relevant.

When subtitle state is ignored, the file shows ignored. When all wanted
languages are matched, the file shows satisfied. When some wanted languages are
missing, the file shows partial or missing depending on whether any wanted
subtitle language is already matched.

For embedded mode, available external subtitle files can change a missing state
into partial because the subtitle exists but still needs import. Missing embedded
subtitle targets are inserted into the track list as red placeholder subtitle
rows.

## Track Management

The track table is assembled from detected embedded tracks, external subtitles
that should be muxed, missing target placeholders, and chapter rows.

Rows can be marked:

- Missing, when a profile target has no satisfying file track.
- Unwanted, when a detected audio or subtitle track is outside enabled profile
  languages or conflicts with subtitle mode.
- Deletable, when the row maps to an embedded audio, subtitle, or chapter
  delete request.

Audio and subtitle delete actions mutate the media file. After deletion or file
rescan, probing refreshes the detected track list and the satisfaction states are
calculated again from the new row data.

External subtitle sidecars are managed separately from embedded subtitle tracks.
They can be searched, imported, moved out, embedded, ignored, or deleted based on
subtitle mode and retention state.
