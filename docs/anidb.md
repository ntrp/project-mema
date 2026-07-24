# AniDB metadata provider

MEMA uses AniDB's public HTTP API for anime search, details, aliases, artwork, dates, and episode facts.

## Configuration

Create or enable the AniDB metadata provider and set:

- **Registered client name** to the application client identifier registered with AniDB. It must match `[a-z0-9]+`: lowercase ASCII letters and digits only, with no spaces.
- **Client version** to the positive integer version registered for that client.

These values identify the MEMA application client. AniDB user credentials are not requested or stored. The default AniDB provider is disabled until this identity is configured. Saving an invalid client name keeps the provider disabled; an invalid name reaching the request boundary returns a configuration error before any AniDB request is sent.

AniDB remains an auxiliary authority for AniDB IDs and episode facts. Enrichment does not replace an existing TMDB/TVDB primary display identity. Absolute numbering is attached only when an AniDB episode ID is already resolved to a local episode through a provider mapping. Specials and other non-regular AniDB episode types remain reviewable facts and are not assigned a standard season.

## Contract and traffic policy

Implementation was verified against the authoritative [AniDB HTTP API Definition](https://wiki.anidb.net/HTTP_API_Definition): requests send `client`, `clientver`, and `protover=1`; anime lookups use `aid` or `aname`. Calls are throttled to at least two seconds apart. HTTP 429 responses are typed and retried once only when the requested delay is at most five seconds. Search/details continue to use MEMA's existing 24-hour metadata cache and history.
