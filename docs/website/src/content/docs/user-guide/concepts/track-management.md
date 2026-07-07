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

## Subtitles

Subtitle mode decides where wanted subtitles should live:

- **Embedded** means wanted subtitles should be internal tracks.
- **External** means wanted subtitles should be sidecar files.
- **Mixed** keeps preexisting external subtitles external, but new subtitle
  downloads are embedded.

External subtitles can be searched manually or automatically from subtitle track
rows and external file rows.

## Chapters

Chapter rows can be deleted individually. The chapter summary row deletes all
chapters in the file.
