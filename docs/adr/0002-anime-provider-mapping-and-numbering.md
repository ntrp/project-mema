# ADR 0002: Anime Provider Mapping And Numbering

## Status

Accepted

## Context

Series and episode rows are now first-class library state. TMDB/TVDB metadata is enough for normal TV, but anime matching needs multiple provider IDs, title aliases, and a selected numbering strategy. Anime remains a classification of existing media types: `movie` or `serie`.

## Decision

Use `media_items.content_kind` to distinguish `standard` from `anime`. Do not add an `anime` top-level media type.

Keep `media_items.external_provider` and `external_id` as the primary display/search default. Store every known provider identifier in `media_provider_mappings` for media items, seasons, and episodes. Supported provider values include placeholders for `anilist` and `anidb`; clients for those providers are out of scope here.

Store provider-derived aliases in `media_item_aliases`. Aliases have a normalized value, optional language, kind, provider provenance, and source payload. User-editable aliases are a follow-up.

Store the selected series numbering strategy on `media_items.numbering_strategy`. Normal series default to `tmdb_season_episode`; anime series default to `anidb_absolute`. Exceptional or provider-specific episode rows can be represented in `media_episode_numbering`.

TMDB remains the default primary display provider when available. Conflicts are represented as data: provider mappings and aliases are persisted, while review UI and user override workflows are deferred.

## Matching Rules

Release search expands queries conservatively with canonical, romaji, english, synonym, and release-title aliases. Release matching accepts the same aliases when scoring a candidate.

Anime imports/searches can treat absolute episode numbers as first-class input when the selected strategy is `anidb_absolute`. Season/episode matching remains the default for normal series.

## Examples

`Frieren: Beyond Journey's End` is a `serie` with `content_kind = anime`, TMDB as primary display provider when available, AniList/AniDB provider mappings when known, aliases such as `Sousou no Frieren`, and default `anidb_absolute` numbering.

`Spirited Away` is a `movie` with `content_kind = anime`; it can store provider mappings and aliases but has no numbering strategy.

## Consequences

Search/import code can use aliases and absolute numbering without reopening the TV/episode schema. Later AniList/AniDB provider clients, user-editable aliases, conflict review, and manual numbering overrides can write into the same tables.
