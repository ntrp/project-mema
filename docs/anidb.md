# AniDB metadata provider

MEMA uses AniDB's official anime-titles dump for title-to-AID search and AniDB's public HTTP API for AID details, aliases, artwork, dates, and episode facts.

## Configuration

Create or enable the AniDB metadata provider and set:

- **Registered client name** to the application client identifier registered with AniDB. It must match `[a-z0-9]+`: lowercase ASCII letters and digits only, with no spaces.
- **Client version** to the positive integer version registered for that client.

These values identify the MEMA application client. AniDB user credentials are not requested or stored. The default AniDB provider is disabled until this identity is configured. Saving an invalid client name keeps the provider disabled; an invalid name reaching the request boundary returns a configuration error before any AniDB request is sent.

AniDB remains an auxiliary authority for AniDB IDs and episode facts. Enrichment does not replace an existing TMDB/TVDB primary display identity. AniDB artwork and start dates are retained as sourced auxiliary facts and in the AniDB provider-mapping provenance. Absolute numbering is attached only when an AniDB episode ID is already resolved to a local episode through a provider mapping. Specials and other non-regular AniDB episode types remain reviewable facts and are not assigned a standard season. Direct AniDB additions are marked as anime and use the `anidb_absolute` strategy.

## Contract and traffic policy

The authoritative [AniDB HTTP API Definition](https://wiki.anidb.net/HTTP_API_Definition) defines `request=anime` as an AID lookup. MEMA never sends the unsupported `aname` parameter. Search resolves a title to an AID locally from AniDB's [official anime-titles XML/GZip dump](https://anidb.net/api/anime-titles.xml.gz), documented under [AniDB API: Anime Titles](https://wiki.anidb.net/API#Anime_Titles), and then requests details with `aid`.

The title dump download has a 15-second deadline, a 16 MiB compressed limit, a 64 MiB uncompressed limit, and a 24-hour in-memory cache. Provider search/details also use MEMA's existing 24-hour persistent metadata cache and history. HTTP API calls send `client`, `clientver`, and `protover=1`, are throttled to at least two seconds apart, and retry HTTP 429 once only when the requested delay is at most five seconds.
