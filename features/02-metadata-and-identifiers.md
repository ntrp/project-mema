# PRD: Metadata and Identifiers

Status: Draft

## Summary

Mema needs reliable metadata and identifier mapping across media domains. It must separate canonical metadata identity from local file identity and release identity.

## Goals

- Integrate metadata providers per media type.
- Maintain stable internal IDs.
- Store external IDs from providers.
- Resolve ambiguous matches.
- Support metadata refresh and drift.
- Support alternate titles and aliases for search.

## Candidate Providers

- Movies: TMDB first; IMDb and TVDB where applicable later
- TV: TVDB and TMDB first; IMDb where applicable later
- Anime: AniList and AniDB mapping investigation first; MyAnimeList and TVDB mapping later as needed
- Books: Google Books, OpenLibrary, Goodreads alternatives, ISBNdb-like providers
- Music: MusicBrainz, ListenBrainz, Discogs where applicable
- Subtitles: OpenSubtitles, Subscene alternatives, embedded metadata, release names

## Functional Requirements

- Users can search external metadata when adding media.
- Mema stores provider IDs and aliases.
- Mema supports alternate titles by language and region.
- Mema supports release dates by country or region where available.
- Mema supports original language and preferred metadata language.
- Mema supports metadata refresh schedules.
- Mema records when provider metadata changes.
- Mema allows manual correction of title, year, edition, season mapping, episode mapping, and aliases.

## Identifier Types

- Internal UUID or ULID
- Provider IDs
- File hash
- Perceptual/audio fingerprint
- Release GUID from indexer
- Download client ID
- Torrent info hash where applicable
- NZB ID where applicable
- Import batch ID
- Final artifact ID

## Acceptance Criteria

- A media item can be matched to multiple external provider IDs.
- Search uses aliases, alternate titles, and year disambiguation.
- Metadata refresh does not overwrite user-locked fields.
- The UI shows metadata source and last refresh time.
- Ambiguous import matches require manual confirmation.

## Open Questions

- Which metadata provider should be canonical per media type after the first provider integrations are proven?
- Should users be able to choose provider priority?
- Should metadata be cached indefinitely?
- Should metadata provider API keys be optional or required?
- How should region-specific titles be selected?
- Should anime use absolute numbering, aired numbering, DVD numbering, or user-selectable mappings?
- Should music use MusicBrainz release groups or specific releases as the main album entity?
- Should books track editions independently from works?
- Should audiobooks and ebooks share the same book/work identity?
- Should local NFO files be read and written?
- Should Mema generate metadata sidecar files for media servers?
