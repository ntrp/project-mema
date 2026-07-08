---
title: Troubleshooting
description: Common setup and workflow problems.
---

Most problems come from one of four places: credentials, paths, scopes, or
profiles. Start with the area closest to the failure, then work outward.

## Provider Tests Fail

If a metadata, subtitle, indexer, or download client test fails, check the base
URL and credentials first. Make sure the service is reachable from the same
environment where Media Manager runs.

For metadata providers, confirm the API key, token, or PIN. For subtitle
providers, confirm the provider account details. For indexers, confirm the API
key and any extra catalog fields. For download clients, confirm the base URL,
username, password, or API key.

## Searches Return Nothing

When metadata search returns nothing, test the metadata provider and try a
simpler title. Include the year when a title is ambiguous.

When release search returns nothing, check that at least one enabled indexer has
the right media scope, matching tags, useful categories, and healthy status.
Also check whether the indexer supports direct search.

## Downloads Finish But Do Not Import

This usually means the app cannot find the completed file. Compare the path
shown by the download client with the path Media Manager can read. If they are
different, add a path mapping.

Also check the download category. If the client puts media in a different
category than expected, the app may not pick up the completed download the way
you expect.

## Background Jobs Look Stuck

Open System > Jobs to see fixed scheduled jobs, current one-shot jobs, and
execution history. Fixed schedules can be run now, paused, or resumed from that
page. Running or queued executions can be aborted from the row actions.

The download client activity sync is a routine fixed schedule with a
configurable interval. It can run as often as every 15 seconds; higher values
can be saved from the fixed schedules table.

Use the execution history logs when a job fails or appears stuck. Each run keeps
structured messages, progress updates, errors, and relevant IDs so you can see
which indexer, download client, media item, or path was involved.

Routine successful sync runs are hidden from the normal execution history so
they do not bury other background work. Turn on Include routine runs when you
need to inspect them, or use the routine retention setting to keep them for a
shorter or longer window.

## The Wrong Indexers Are Used

Review media scopes and tag scopes. Media without tags can use any enabled
indexer for the matching media type. Media with tags only uses enabled indexers
that share at least one tag.

If an indexer should never be used for a media type, remove that scope from the
indexer. If a tagged item should use a specific indexer, give both the media
item and the indexer the same tag.

## A File Is Marked Partial

Partial usually means the file exists but does not fully match the selected
profile. Open the file detail view and compare the detected audio tracks,
subtitle tracks, quality, and custom format score with the profile.

If the file is acceptable to you, adjust the profile. If the profile is correct,
search for subtitles, replace the release, or fix the file.

## Subtitle Search Finds Bad Matches

Check the media metadata first, especially title, year, season, and episode.
Then check subtitle provider settings and profile subtitle targets.

If the provider finds candidates in the wrong language, add or improve language
aliases. If the provider finds no candidates, try manual search with a simpler
query or allow searching subtitles in other releases from the profile.

## Quality Or Score Looks Wrong

Check quality sizes and custom formats. Quality sizes affect whether a release
is a reasonable size for the selected quality. Custom formats affect scoring for
specific traits.

If good releases score poorly, the profile may be too strict. If bad releases
score well, add or adjust custom formats, raise minimum scores, or tighten
quality size ranges.

## Imports Match The Wrong Title

Use the import row controls to choose the correct metadata result before
importing. If matching is often wrong for a library, check file naming, title
aliases in the source files, and whether the folder type is correct.

Movie folders should contain movies. Series folders should contain series. A
mixed root makes automatic matching harder.
