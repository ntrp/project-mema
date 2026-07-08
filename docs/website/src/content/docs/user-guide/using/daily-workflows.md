---
title: Daily Workflows
description: Add media, search releases, follow downloads, and maintain files.
---

Once setup is complete, most work happens from Discovery, media detail pages,
manual search, wanted items, activity, and library import.

The app is designed so you can move between automatic and manual control. You
can let profiles guide release decisions, but you can still open a manual search
and choose a specific release when you want to override the automatic path.

## Add Media

Use Discovery or Search to find a movie or series. Open the result, choose the
library folder and profile, then add it. If you want the app to keep looking for
files or upgrades, leave monitoring enabled.

For series, monitoring can apply at different levels. A whole series, season, or
episode can be treated differently depending on what you want the app to manage.

## Search Releases

Manual search shows release candidates from enabled indexers that match the
media type, tags, and current indexer health. The results include scoring and
status information so you can see why a release looks good or bad. Missing
wanted audio or subtitle languages are shown as warnings, not hard rejections.
Automatic search gives those releases lower priority than otherwise comparable
releases that satisfy the language targets.

Use manual search when automatic search does not find what you want, when you
want to compare indexers, or when you want to grab a specific release.

Wanted can show missing media, specific unsatisfied profile targets, and
custom-format upgrade rows. Target rows include the parent media and file
context, the language or target type when relevant, and the operation that can
move the item forward.
Media cards and file cards use the same rollup state so missing, partial,
downloaded, and upgradeable states read consistently across library and detail
views.

## Grab And Download

When you grab a release, the app sends it to the matching download client.
Activity shows the download and later import state.

If a download finishes but does not import, check activity first. Then check the
download client category and path mappings. The most common issue is that the
download client reports a path the app cannot read.

## Import Existing Files

Use Library scans for files that already exist on disk. Scan a root folder,
review the rows, match each file to a media item or metadata result, choose the
profile, then import.

Start with a small batch. Existing libraries often contain naming differences,
duplicate files, old samples, extras, or alternate cuts. A small first import
makes it easier to tune matching and naming before importing everything.

## Review File Status

After import, open the media detail page. The file row shows file, audio,
subtitles, size, quality, score, status, and actions. Status details are shown
through compact badges and hover details so the row stays readable.

Use the file detail area when the row says partial or missing. It will show
which tracks and sidecars were detected and what the profile still wants. Track
rows show whether they match the profile, partially match it, need a follow-up
operation, are unwanted by current settings, or stand in for a missing target.

## Refresh And Maintain

Use metadata refresh when titles, posters, seasons, episodes, or people look
stale. Use subtitle search when subtitles are missing. Use rename preview when
you want to review file names before applying naming templates.

Every automatic fulfillment path also has a manual route. You can manually
search or grab releases, retry import, search or grab subtitles, rescan files,
and use component actions for remuxing, embedding, extraction, and stream
sourcing where the current media context supports them. Turning off an automatic
schedule stops future background runs, but it does not remove the matching
manual action.

System > Jobs lists the manual fulfillment action catalog next to scheduled job
controls, including the API route and worker path used by each action.

Treat profiles as living rules. If every release you like is being penalized,
adjust the profile or custom formats. If bad releases score too well, tighten
quality sizes, custom formats, or language targets.
