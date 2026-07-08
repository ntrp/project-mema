---
title: Subtitles, Audio, And Tracks
description: Understand profile status, subtitle search, and track actions.
---

Media Manager inspects media files after import so you can see what each file
contains. The file detail area shows video information, audio tracks, subtitle
tracks, chapters, sidecar files, and compact status badges.

Those badges are compared against the selected profile. If a profile wants
German audio and English subtitles, the file status is judged against those
targets rather than against a generic idea of what a “good” file is.

## Audio Status

Audio status checks the profile’s audio targets against the tracks in the file.
Language is the most important part. Codec, channels, and bitrate matter when
you configured them.

If audio is marked partial or missing, open the file detail area and compare the
detected tracks with the selected profile. The file may be perfectly playable
but still not satisfy the profile you chose. Missing required audio languages
appear as placeholder rows in the track table, next to the detected audio
tracks. Audio rows that miss codec, channel, or bitrate targets are marked
partial and list the failed requirement in their hover details.

## Subtitle Status

Subtitle status checks wanted subtitle languages and formats. It also respects
the subtitle mode from the profile: Embedded, External, or Mixed.

When subtitles are missing, use subtitle search from the media page. If a
subtitle provider is configured, the app can search candidates and attach the
selected subtitle to the media item.

The subtitle rows reflect subtitle mode immediately. An embedded subtitle can
show as pending extraction for External mode, while an external sidecar can show
as pending embedding for Embedded mode. Extra subtitle rows are marked unwanted
only when the profile settings or subtitle mode make them unwanted.

## Subtitle Providers

Open Settings, then Subtitles. OpenSubtitles is the provider to configure for
real subtitle searches. Add the required credentials, enable it, save, and test.

The mock provider is useful only for predictable local testing. It is not a
replacement for a real subtitle source.

## Track Actions

The media file detail page can show actions for audio tracks, subtitle tracks,
chapters, and sidecar files. Available actions depend on the file and installed
media tools.

Deleting a track changes the media file. Use the profile status and track list
to confirm you are removing the right thing. If you often remove the same
unwanted tracks, consider expressing that in the profile so the file status
makes the problem obvious before you act.

## Preview And Inspection

Preview helps confirm that the file plays and that the detected audio and
subtitle choices make sense. Inspection details are especially useful after
manual imports, subtitle downloads, or track edits.

When a status badge looks surprising, trust the detailed track list first. It
shows what the app detected and gives you the clues needed to adjust the file,
profile, or language aliases.
