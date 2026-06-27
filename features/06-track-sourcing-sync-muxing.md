# PRD: Track Sourcing, Synchronization, and Muxing

Status: Draft

## Summary

Mema must support independently sourcing video, audio, and subtitle tracks from multiple releases, synchronizing them, and muxing a final media artifact. This is the most technically distinctive and risky feature.

## Goals

- Allow target video quality, audio tracks, and subtitle tracks to be defined independently.
- Search for missing components after a base video is selected.
- Extract usable tracks from downloaded files.
- Synchronize external audio and subtitle tracks to the target video.
- Mux final files with deterministic stream order, language tags, titles, and defaults.
- Preserve source provenance and confidence scores.
- Require manual review when sync confidence is low.

## Example Requirement

For animated movies:

- Video target: 1080p
- Audio target: German and English
- Subtitle target: Italian

Mema should be able to select a 1080p video release, source German audio from another release, source English audio from another release, source Italian subtitles from another release or subtitle provider, synchronize all tracks to the selected video, and produce one final file.

## Functional Requirements

- Users can define component profiles for video, audio, and subtitles.
- Mema can determine which target components are missing from an imported file.
- Mema can search for releases likely to contain missing components.
- Mema can download candidate component sources without final library import.
- Mema can inspect streams in each candidate file.
- Mema can extract selected streams.
- Mema can compare runtime, frame rate, chapter structure, and fingerprints to estimate compatibility.
- Mema can synchronize subtitle timing.
- Mema can synchronize audio timing where feasible.
- Mema can detect likely incompatible cuts or editions.
- Mema can mux selected streams into a final container.
- Mema can tag language, title, forced/default flags, disposition, and track order.
- Mema can store every source release used for every final track.
- Mema can show final assembly status and logs.
- Mema can fall back to manual review.

## Processing Stages

1. Select base video.
2. Analyze base video streams.
3. Determine missing target components.
4. Search candidate component sources.
5. Download candidate sources.
6. Analyze streams from candidate sources.
7. Select component tracks.
8. Estimate compatibility.
9. Extract tracks.
10. Synchronize tracks.
11. Mux final file.
12. Validate final file.
13. Import final file.
14. Retain or clean source components.

## Technical Capabilities Needed

- MediaInfo or ffprobe stream analysis.
- ffmpeg extraction and muxing.
- mkvmerge support, likely preferred for Matroska final assembly.
- Subtitle sync tools or internal subtitle timing algorithms.
- Audio fingerprinting for sync offset detection.
- Runtime and chapter comparison.
- Black-frame/silence or dialogue-anchor detection, if needed later.
- Worker queue with cancellation and retry.
- Persistent logs and artifacts.

## Confidence Model

Mema should compute confidence for each component:

- Exact same release: high
- Same runtime, frame rate, and chapter layout: high
- Same runtime but different source group: medium
- Small constant offset detected: medium
- Multiple drift points detected and corrected: medium or low
- Different runtime, different edition, or unknown frame rate: low
- Failed validation: reject

## Acceptance Criteria

- Mema can identify missing audio/subtitle targets for a file.
- Mema can extract audio and subtitle streams from a downloaded file.
- Mema can mux selected streams into a final MKV.
- Mema can set language metadata and default/forced flags.
- Mema can retain provenance per stream.
- Mema refuses automatic muxing when compatibility confidence is below threshold.
- Users can manually approve or reject component candidates.

## Open Questions

- Should the first version support muxing only MKV output?
- Should MP4 output be supported at all, given subtitle/audio limitations?
- Which external tools are acceptable hard dependencies?
- Should audio sync be automatic in MVP, or manual review only?
- What synchronization confidence threshold should allow automatic muxing?
- Should Mema support variable frame rate correction?
- Should Mema support PAL speedup and slowdown correction?
- Should Mema support matching different cuts, editions, censored versions, or only exact runtime matches?
- Should dubbed audio be searched from lower-quality video releases?
- Should Mema keep the donor video file after extracting audio?
- Should Mema support commentary tracks?
- Should Mema support forced subtitles separately from full subtitles?
- Should Mema support anime signs/songs subtitles as a special subtitle type?
- Should users be able to define stream order and defaults per profile?
- How should Mema validate final sync without human playback?
- Should manual review include web playback with offset controls?
- Should Mema write sidecar project files for reproducible remuxing?

