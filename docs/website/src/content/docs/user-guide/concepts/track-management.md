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
change can queue container remuxing. Automatic background schedules for these
operations exist in System > Jobs but start disabled. Video transcoding can fix
supported codec and pixel-format mismatches. HDR-only or resolution-only
mismatches are shown as profile problems, but they are not queued for automatic
transcoding until a safe tool path exists.

Scheduled video transcoding scans the library and creates one current one-shot
job for each eligible video track. Manual video transcode buttons create the
same track-scoped one-shot job for the selected row. Progress appears in System
> Jobs while the media tool reports transcoding progress, then the file is
rescanned when replacement finishes.

## Chapters

Chapter rows can be deleted individually. The chapter summary row deletes all
chapters in the file.
