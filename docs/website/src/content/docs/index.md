---
title: Media Manager Documentation
description: Documentation for the self-hosted video-first media manager.
template: splash
hero:
  tagline: A self-hosted media manager for movies, TV, anime metadata, subtitles, release search, downloads, and MKV track management.
  actions:
    - text: Start locally
      link: /user-guide/getting-started/first-run/
      icon: right-arrow
      variant: primary
    - text: Learn the model
      link: /user-guide/concepts/media-lifecycle/
      icon: open-book
---

## What This App Does

The application manages a video library from request through discovery, release
search, download import, profile compliance, subtitles, and track-level media
maintenance.

It is built as a Go API with a SvelteKit browser app, PostgreSQL storage, River
jobs, contract-first OpenAPI, and local media tools such as `ffmpeg`,
`ffprobe`, `mkvmerge`, and `mkvextract`.

## Documentation Sections

- **User Guide** covers local setup, first-run commands, configuration, media
  lifecycle, metadata providers, indexers, download clients, profiles, imports,
  subtitles, audio, and track management.
- **Dev Guide** captures the development workflow used in this repository.
- **Architecture** describes runtime structure, storage ownership, and API
  contract rules.

## Current Status

The project is still pre-release. Schema changes are applied directly to the
initial schema, and development database resets happen through external dev
tooling rather than through the final application.
