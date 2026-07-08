---
title: Indexers And Download Clients
description: Configure release search and downloads.
---

Indexers and download clients work together, but they do different jobs. An
indexer finds releases. A download client downloads the release you choose.

Transmission is used for torrent downloads. SABnzbd is used for Usenet
downloads. Indexers can be torrent, Usenet, or another supported protocol
depending on the selected catalog entry.

## Download Clients

Open Settings, then Download Clients. Add one client for each download service
you want Media Manager to use.

Transmission needs a base URL and, when your instance requires it, a username
and password. SABnzbd needs a base URL and API key. Both client types can use a
category, priority, and enabled state.

Use categories to keep Media Manager downloads separate from other downloads.
For example, a category such as `movies` or `media` makes completed files easier
to recognize and import. Use the test action after saving a client.

## Indexers

Open Settings, then Indexers. New indexers start from the catalog. Pick the
site or service, then review the prefilled values before entering credentials.

The most important fields are the base URL, API key, categories, enabled state,
and priority. Some catalog entries include extra fields. Fill those exactly as
your indexer expects.

Categories tell the indexer which parts of the site to search. If categories are
too narrow, good releases may never appear. If they are too broad, unrelated
releases may appear in searches.

## Media Scopes

Media scopes decide which content types an indexer should handle. The app shows
Movies, Series, Anime, Audio, and Books.

Keep scopes honest. A movie tracker should usually have Movies enabled. A TV
tracker should usually have Series enabled. Anime can be separated when you use
specialized anime indexers or want different release naming behavior.

## Tag Scopes

Tags are optional, but they are useful when a media item should search only a
specific group of indexers. A common pattern is a tag for a private tracker, an
anime source, or a language-specific source.

The matching rule is intentionally simple. Media without tags can use any
enabled indexer that matches the media scope. Media with tags only uses enabled
indexers that share at least one of those tags.

This lets untagged media behave broadly while tagged media becomes deliberate.

## Health And Priority

Indexer health helps explain why an indexer is or is not being used. A healthy
indexer can be searched. A temporarily disabled indexer has recently failed and
may be skipped until it is checked again. A disabled indexer is off.

Priority controls ordering. Use it to prefer the services you trust most. If two
indexers are equally good, give the one you want searched first the better
priority.

## Search Support And Recent Feeds

Some indexers support direct search. Some support recent feeds. Some support
both. Direct search is used when you manually search for releases or when the
app searches for a specific wanted item. Recent feeds are used for workflows
that look at newly posted releases.

If an indexer is enabled but never appears useful, check its support icons,
categories, media scopes, tag scopes, and health.
