---
title: How Media Manager Works
description: The main parts of the app and how they fit together.
---

Media Manager is built around one question: what should this movie or series
look like when it is finished?

You add media, choose where it belongs, assign a profile, and let the app keep
track of the rest. Metadata providers describe the title. Indexers find possible
releases. Download clients fetch the selected release. Library folders decide
where files live. Profiles decide whether the file is good enough or still needs
work.

## The Main Pieces

Metadata providers are the source of titles, posters, overviews, seasons,
episodes, cast, crew, collections, recommendations, and related media. The app
can use more than one provider, so a movie and a series can come from the source
that has the best match.

Library folders are the places where finished media is stored. A folder is
marked as either Movies or Series so imports can use the right naming and folder
rules.

Indexers are the places the app searches for releases. Some indexers support
manual search, some support RSS-style recent feeds, and some support both. Their
media scopes and tag scopes decide when they should be used.

Download clients do the actual downloading. Transmission handles torrent
downloads, and SABnzbd handles Usenet downloads. Media Manager tracks the
download and later imports the completed files.

Profiles describe the desired final file. They include qualities, upgrade
rules, audio languages, subtitle languages, custom format scores, and the final
container. A profile is what turns a file from “downloaded” into “meets the
requirements.”

## The Usual Flow

A typical item starts in Discovery or Search. You add it, choose a profile and a
library folder, and decide whether it should be monitored. When the app searches,
it asks the matching indexers for releases and scores the results against the
profile. A selected release is sent to the matching download client.

When the download finishes, the app tries to import the file. During import it
matches the file to the media item, moves or attaches it to the library, reads
its video, audio, subtitle, and sidecar information, and updates the media page.
The file then shows compact status badges for audio, subtitles, quality, score,
and overall state.

If the file is incomplete, follow-up actions can search for subtitles, inspect
tracks, delete unwanted tracks, rename files, or replace the release with a
better one when upgrades are allowed.

## What “Ok”, “Partial”, And “Missing” Mean

`Ok` means the current file satisfies the selected profile. It has an acceptable
quality and meets the required audio and subtitle targets.

`Partial` means there is a usable file, but at least one requested part still
needs attention. A common example is a movie with the right video quality but a
missing subtitle language.

`Missing` means the app does not have a usable file or does not have a required
component.

These labels are meant to guide the next action. They do not mean the file is
unplayable; they mean the file does or does not match the profile you selected.
