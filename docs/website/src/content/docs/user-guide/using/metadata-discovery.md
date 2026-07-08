---
title: Metadata And Discovery
description: Configure metadata providers and use them to find media.
---

Metadata is what turns a file name or search term into a real media item with a
title, year, poster, overview, seasons, episodes, cast, crew, and related media.

Media Manager currently exposes TMDB and TVDB as the main metadata providers in
Settings. You can enable either one or both. Provider priority decides which
provider is preferred when more than one can answer a request.

## Choosing Providers

TMDB is a good default for movies, collections, people, posters, and general
discovery. TVDB is useful when you care about series and episode structure.

For most setups, enable both. Keep the base URLs at their defaults unless you
use a proxy or a compatible hosted endpoint. Add the credentials each provider
requires, save, and test from the provider card.

If a provider test fails, fix that before continuing. A broken metadata provider
can make search, discovery, metadata refresh, and library matching feel broken
even when the rest of the app is healthy.

## Discovery

Discovery is the browsing area for finding new media. It uses enabled metadata
providers to show results, details, cast, crew, related media, recommendations,
collections, seasons, and episodes.

When you add something from discovery, you choose the library destination,
profile, monitoring behavior, and other media settings. That choice creates the
local media item that later search, downloads, imports, and profile checks work
against.

## Refreshing Metadata

Metadata can change after an item is added. Posters may be updated, episode
lists may be corrected, or a title may receive better translations. Use metadata
refresh from the media item when the local details look stale.

Refreshing metadata updates the descriptive information. It does not replace
your local file by itself and it does not change your library folders or profile
choices unless you edit those media settings.

## Search And Matching

Search uses metadata providers to turn a title query into selectable media
results. Library import also uses metadata when it needs help matching a file to
a known title.

If searches return the wrong title, check the provider choice, year, media type,
and language aliases. If searches return no results, test the provider first,
then try a simpler title.
