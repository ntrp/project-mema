---
title: Media Lifecycle
description: How media moves from request to managed files.
---

The application is organized around a media item and the desired final state of
its files.

## Main Flow

1. A media item is added or requested.
2. Metadata providers hydrate title, year, poster, collection, cast, crew, and
   related media.
3. Search providers return release candidates for automatic or manual search.
4. A selected release is sent to a download client.
5. Completed downloads are imported into the configured library folder.
6. File probing reads video, audio, subtitle, chapter, and sidecar state.
7. The selected media profile decides whether the file is ok, partial, or
   missing required components.
8. Follow-up actions can search subtitles, import external subtitles, delete
   tracks, or prepare future processing jobs.

## Important Surfaces

- **Discovery** is for finding media and adding it to the library.
- **Manual search** is for selecting releases directly.
- **Library import** is for matching existing files to media items.
- **Media detail** is for inspecting files, tracks, subtitles, chapters, and
  fulfillment status.
- **Settings** owns providers, indexers, download clients, library folders,
  profiles, languages, and system diagnostics.
