---
title: Track Management
description: Audio, subtitle, chapter, and external file handling.
---

Media file details expose embedded tracks, chapters, external subtitle files,
and sidecar files.

## Audio

Audio status checks the selected profile requirements:

- target codec
- target channels
- minimum bitrate
- required language

Audio tracks can be deleted from the file detail view when the media toolchain
supports the operation.

When a row is missing or partial, action buttons can queue audio transcoding or
audio sourcing. These manual buttons use the same background workers as the
disabled-by-default automatic fulfillment jobs.

## Subtitles

Subtitle mode decides where wanted subtitles should live:

- **Embedded** means wanted subtitles should be internal tracks.
- **External** means wanted subtitles should be sidecar files.
- **Mixed** keeps preexisting external subtitles external, but new subtitle
  downloads are embedded.

External subtitles can be searched manually or automatically from subtitle track
rows and external file rows.

Subtitle rows and external subtitle files can also expose fulfillment actions:
download missing subtitles, embed external subtitles, extract embedded
subtitles, or convert text subtitle formats. Available actions depend on the
current subtitle mode, target format, and known file context.

## Video And Container

Partial video rows can queue video transcoding. Rows that need only a container
change can queue container remuxing. The automatic Media Fulfillment schedule in
System > Jobs scans for these operations and starts disabled. Video transcoding can fix
supported codec and pixel-format mismatches. HDR-only or resolution-only
mismatches are shown as profile problems, but they are not queued for automatic
transcoding until a safe tool path exists.

Container remuxing copies every stream into the profile final container without
re-encoding. Media Fulfillment scans files that need only a container change and
creates one current one-shot job per file. Manual remux buttons create the same
file-scoped one-shot job from the file overview row.

Media Fulfillment also creates one current one-shot job for each eligible video
or audio track, subtitle extraction, subtitle conversion, or subtitle merge.
Manual buttons create the same scoped one-shot jobs for selected rows. System >
Jobs shows Media Fulfillment scan progress as processed media entries out of the
total, and child job progress while each media tool runs. Files are rescanned
when replacement finishes.

## Chapters

Chapter rows can be deleted individually. The chapter summary row deletes all
chapters in the file.
