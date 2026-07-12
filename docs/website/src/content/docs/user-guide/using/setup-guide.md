---
title: Setup Guide
description: Configure the app in the order that avoids most first-run problems.
---

The easiest setup path is to configure the services that identify media first,
then the services that find and download releases, and only then the profiles
that decide what a finished file should contain.

You can change any of these settings later. The order below simply gives the app
enough information to make good choices while you add your first media.

## Start With Metadata

Open Settings, then Metadata. Enable the metadata providers you want to use and
enter their credentials. TMDB is the normal starting point for movies and broad
discovery. TVDB is useful for series and episode information.

Save each provider and use the test action before moving on. A provider can be
enabled but still fail because the API key, PIN, token, or base URL is wrong.
Testing here saves time later when searches or refreshes return nothing.

## Add Languages

Open Settings, then Languages. Add the languages you care about before creating
profiles. Profiles use these languages for audio and subtitle targets.

Use aliases for names that often appear in file names or release titles. For
example, German might include aliases such as `DE`, `DEU`, `GER`, and `Deutsch`.
Aliases help the app understand releases and tracks that use short or regional
labels.

## Add Library Folders

Open Settings, then Library, and add at least one movie folder and one series
folder if you manage both kinds of media. The folder type matters because movies
and series use different matching, naming, and import behavior.

Use paths as the app sees them. In a container setup, that might not be the same
path your download client reports. Path mappings handle that difference.

## Add Path Mappings

Path mappings translate a path reported by a download client into the path Media
Manager can read.

For example, if SABnzbd reports completed files under `/downloads/complete` but
Media Manager sees the same files under `/data/downloads`, add a mapping from
`/downloads` to `/data/downloads`. Without this, downloads can finish
successfully but fail to import because the app cannot find the completed file.

## Add Download Clients

Open Settings, then Download Clients. Add Transmission for torrents, SABnzbd for
Usenet, or both.

Give each client a clear name, base URL, credentials, category, priority, and
enabled state. The category should match the category or label you expect the
download client to use for media downloads. Test the connection after saving.

Priority matters when more than one client can handle a release. Lower-priority
numbers are treated as more important in the app’s settings lists and selection
behavior.

## Add Indexers

Open Settings, then Indexers. Start from the catalog, choose the indexer you use,
then enter the base URL, API key, categories, and any extra fields the selected
indexer requires.

Set media scopes so the indexer is used for the right content. A movie-only
indexer should not be searched for series unless you intentionally enable that.
If you use tags, set tag scopes too. Media with tags only searches tagged
indexers that share at least one tag. Media without tags can use any matching
indexer for its media type.

Save and test each indexer. A healthy indexer is the foundation for manual
search, automatic search, and recent-feed workflows.

## Review Quality Sizes

Open Settings, then Quality. Quality sizes define the file-size ranges the app
uses when judging releases. They help reject releases that are far too small or
far too large for a quality.

The defaults are enough for a first run. Revisit them when you notice releases
being rejected or accepted for the wrong size reason.

## Create Custom Formats

Open Settings, then Custom Formats. Custom formats let you reward or reject
specific release traits such as a codec, source, release group, edition, indexer
flag, language, or text pattern in the release title.

Create custom formats before profiles when you want profile scores to use them.
A format can have conditions that must match and negated conditions that reject
a match. Keep names practical, because they appear later in profile scoring and
rename templates.

## Create Profiles

Open Settings, then Profiles. Create at least one default profile. Select the
qualities you want, the final container, upgrade behavior, audio targets,
subtitle targets, custom format scores, and score thresholds.

Profiles are the most important setup choice because every added media item uses
one. Start with a simple profile, test it on a few files, then add stricter
rules once you understand how your indexers name releases.

## Configure DLNA Device Profiles

Open Settings, then DLNA, then Device Profiles when a TV or player needs
renderer-specific compatibility rules. Search the seeded profile table, select
a profile to edit it in a modal, clone it into a custom profile, create a new
custom profile, or reset a seeded profile back to its defaults.

Recent DLNA devices show the IP, last seen time, matched profile, and override
selector. Use the override selector or manual override form to pin a known
device, such as an LG TV IP or renderer UUID, to a chosen profile. Use the
decision trace panel to compare a selected device and media file path against
the effective profile rules before changing playback settings. Run trace shows
why a renderer matched a profile and why the selected media will direct play,
remux, transcode, or use HLS; the output uses the file name instead of exposing
the full media path.

Renderer profile delivery settings can also control output container and seek
mode. Use byte seek for normal range requests, time or time-exclusive for
renderers that prefer DLNA time seeking, and none for devices that fail on seek.
Profile subtitle, artwork, and metadata rules control what appears in the
renderer browse response: subtitle resource formats, album-art URLs, thumbnail
behavior, dates, rich media fields, folder data, and child counts. Trim these
fields when a TV lists folders slowly, shows duplicate artwork, or rejects
subtitle resources it cannot load.
Mema also supports renderer search, sorted browse results, and UPnP event
subscriptions; disable eventing in a profile only when a client fails during
subscription setup.

When profile edits make playback worse, reset a seeded profile to restore the
current built-in default, or clone the profile before experimenting. If a seed
profile receives an upgrade, customized profiles keep their local edits until
you reset them. Use recent devices, profile match trace, and delivery decision
trace to troubleshoot common cases: a TV matching the wrong profile, LG audio
falling back to transcode, seeking failures, missing artwork, or subtitle
formats the renderer cannot load. DLNA cleanup runs on startup and removes stale
thumbnail, remux, and transcode cache files; active stream and transcode rows in
Settings > DLNA show current long-running work.

## Add Subtitle Providers

Open Settings, then Subtitles. The picker shows the full Bazarr-compatible
provider catalog. Every catalog entry has a native runtime and can be configured,
enabled, and tested. Provider warnings describe requirements such as API keys,
private-site membership, browser cookies for CAPTCHA-protected sites, local
Whisper services, or media identifiers.

Configure OpenSubtitles.com for a straightforward online subtitle source, or
choose the providers serving your languages and media. The mock provider is
useful only when you want predictable test data.

Subtitle providers work together with profile subtitle targets. The provider
finds candidates; the profile decides which languages and formats are wanted.

## Try A Small Workflow

After setup, add one movie or one series episode. Run a manual search, inspect
the release scores, grab a release, and watch the download activity. When the
download completes, confirm that import finds the file, attaches it to the media
item, and shows the expected file status.

If this small workflow works, the core setup is ready.
