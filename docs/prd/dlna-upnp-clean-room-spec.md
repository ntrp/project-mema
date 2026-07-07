# PRD: DLNA and UPnP Clean-Room Specification

## Status

Draft

## Context

Mema currently has browser preview playback, direct HTTP file streaming, VLC
playlist generation, media probing, and HLS segment generation. Those pieces
are not a DLNA server. The existing `internal/playback` package contains a
device-profile decision model for browser playback, while `internal/httpapi`
owns ffprobe, ffmpeg HLS segmenting, range-capable file serving, preview
tokens, and preview-specific endpoints.

The goal is to implement full DLNA and UPnP AV capabilities without creating a
second streaming stack. Browser preview and DLNA must share one media delivery
system for direct serving, remuxing, transcoding, probing, thumbnails,
subtitles, and compatibility decisions.

## Research Inputs

The following public projects were cloned under `.research` for behavioral
analysis:

- `.research/anacrolix-dms`: UPnP DLNA media server behavior, SSDP, SOAP,
  ContentDirectory, ConnectionManager, raw resources, transcodes, thumbnails,
  subtitles, and compatibility workarounds.
- `.research/huin-goupnp`: UPnP client discovery, device XML, SCPD, SOAP, and
  generated MediaServer and MediaRenderer client contracts.
- `.research/koron-go-ssdp`: SSDP advertiser, monitor, search, alive, byebye,
  and interface-limited multicast behavior.

This specification is clean-room. It defines Mema's desired behavior from
public protocol concepts, observed feature categories, and current Mema
architecture. Implementation must be original and tested against protocol
contracts, not copied from the researched source code.

## Goals

- Provide a UPnP AV 1.0 and DLNA 1.5 compatible Digital Media Server.
- Expose the Mema library through ContentDirectory browse and search.
- Advertise the server through SSDP on configured network interfaces.
- Serve root device XML and service SCPD documents.
- Handle SOAP control for ContentDirectory and ConnectionManager.
- Support raw file resources, compatible direct streams, remuxed streams,
  transcoded streams, thumbnails, artwork, and subtitles.
- Share one media delivery implementation between browser preview, VLC
  streaming, and DLNA.
- Preserve all existing browser preview behavior and public API semantics.
- Support device profiles and compatibility overrides per renderer.
- Provide observability, diagnostics, and compatibility tests.
- Default to a safe LAN-only posture and require explicit enablement.

## Non-Goals

- Copying implementation code from researched projects.
- Replacing Mema's metadata providers or library model.
- Exposing the full web application session surface to DLNA clients.
- Requiring internet exposure, remote access, or router traversal.
- Guaranteeing every renderer-specific quirk in the first release.
- Implementing write-capable ContentDirectory mutations unless explicitly
  enabled by a later story.

## Current Mema Playback Capabilities

Current preview functionality must continue to work:

- `internal/playback` chooses direct play, direct stream, or transcode.
- `BrowserVideoProfile` supports direct MP4/M4V H.264/AAC playback and HLS
  MPEG-TS output with H.264/AAC.
- WebKit preview forces HLS video transcode when native HLS cannot use a copied
  video stream.
- HTTP streaming uses `http.ServeContent` and byte range support.
- VLC playback produces an M3U with a time-limited stream token.
- Preview HLS playlists use fixed or keyframe-aligned six second segments.
- Preview segment serving runs ffmpeg and streams MPEG-TS bytes.
- ffprobe reads streams, chapters, container format, duration, keyframes,
  bitrate, channel count, language, title, and codec metadata.
- Tool execution is restricted through `internal/tools` validation, timeouts,
  output limits, and safe path checks.

The DLNA implementation must consume the same underlying delivery planner and
tool runner instead of implementing separate ffprobe, ffmpeg, or file-serving
logic.

## Architecture

### Packages

The implementation should be split into small single-purpose packages. No Go
module should exceed the project line limits.

Proposed package boundaries:

- `internal/dlna/server`: lifecycle, settings, network interface selection,
  startup, shutdown, and integration with the app HTTP server.
- `internal/dlna/ssdp`: SSDP message construction, multicast listeners,
  M-SEARCH handling, alive announcements, byebye announcements, and interface
  binding.
- `internal/dlna/upnp`: root device documents, service descriptors, SOAP
  envelope handling, UPnP errors, action dispatch, and event subscription
  primitives.
- `internal/dlna/content`: ContentDirectory tree, object IDs, browse, search,
  pagination, sorting, system update IDs, and DIDL-Lite mapping.
- `internal/dlna/media`: DLNA resource descriptors, protocolInfo generation,
  content feature headers, seek and range handling, and delivery integration.
- `internal/dlna/profile`: renderer capability profiles, user-agent matching,
  DLNA profile names, direct play profiles, and transcode targets.
- `internal/playback`: shared media source, stream decision, media delivery
  plan, and device profile logic used by both preview and DLNA.
- `internal/media`: shared probing, keyframe discovery, file stat, MIME
  detection, thumbnail generation, and subtitle resource helpers, if extracting
  from `httpapi` is needed.
- `internal/httpapi`: thin endpoint adapters for preview, stream, VLC, and
  DLNA resource URLs.

### Shared Delivery Service

Create a `MediaDeliveryService` that owns the behavior currently split across
preview and stream handlers.

Responsibilities:

- Resolve a media item file path through the settings/library service.
- Validate absolute local file paths before tool execution.
- Probe streams, chapters, container metadata, duration, and keyframes.
- Convert probe output into `playback.MediaSource`.
- Select a `playback.DeviceProfile`.
- Build a delivery plan with direct file, remux, or transcode method.
- Serve direct files with byte range support and correct content type.
- Generate HLS playlists for preview and HLS-capable DLNA clients.
- Generate MPEG-TS or other DLNA transcode output for clients that do not
  support HLS.
- Produce DLNA resource metadata, including protocolInfo and content features.
- Generate or serve poster, backdrop, thumbnail, and album art resources.
- Serve subtitle resources or converted subtitle formats where supported.
- Emit structured logs and metrics for probe, decision, and tool execution.

Existing preview handlers should become adapters:

- Preview info asks the delivery service for a browser delivery plan.
- Preview endpoint serves the direct file or generated playlist from that plan.
- Preview segment endpoint asks the delivery service to execute the planned
  segment.
- VLC endpoint continues to mint stream tokens but serves through the same
  direct file delivery path.

DLNA handlers should use the same service:

- ContentDirectory asks for available resources for each library item.
- Resource endpoints execute the selected delivery plan.
- Renderer profiles choose direct, remux, transcode, subtitle, and artwork
  resources without changing preview behavior.

## UPnP and DLNA Server Surface

### SSDP

The server must:

- Join SSDP multicast on configured interfaces.
- Support IPv4 multicast at `239.255.255.250:1900`.
- Support IPv6 link-local and site-local SSDP where the platform allows it.
- Send periodic `ssdp:alive` NOTIFY packets.
- Send `ssdp:byebye` NOTIFY packets during clean shutdown.
- Respond to `M-SEARCH * HTTP/1.1` with `MAN: "ssdp:discover"`.
- Honor `ST` values for `ssdp:all`, `upnp:rootdevice`, the server UUID, the
  MediaServer device type, and advertised service types.
- Clamp or validate `MX` and randomize response delay within the allowed
  search window.
- Include `CACHE-CONTROL`, `EXT`, `LOCATION`, `SERVER`, `ST`, and `USN` in
  search responses.
- Include `HOST`, `NT`, `NTS`, `USN`, `LOCATION`, `SERVER`, and
  `CACHE-CONTROL` in alive notifications.
- Build `LOCATION` from the interface address that can reach the client.
- Allow interface allowlist configuration.
- Ignore loopback, down, non-multicast, and link-local addresses that cannot
  produce usable device description URLs.
- Tolerate imperfect SSDP packets in tests where safe, including missing final
  header terminators.

### Root Device Description

The server must expose a root device XML document with:

- UPnP device namespace and spec version.
- Device type `urn:schemas-upnp-org:device:MediaServer:1`.
- Stable UDN generated from server identity.
- Friendly name configured by the user.
- Manufacturer, model name, model number, serial number, and presentation URL.
- DLNA capability and DLNA document metadata where needed for compatibility.
- Icon list with one or more generated or bundled icons.
- Service list for ContentDirectory and ConnectionManager.
- Optional service list entries for compatibility services.

### Service Descriptors

The server must expose SCPD XML for:

- ContentDirectory v1.
- ConnectionManager v1.
- Optional Microsoft MediaReceiverRegistrar compatibility service.
- Optional MediaRenderer services if renderer/control capabilities are added.

Service descriptors must list action arguments, directions, related state
variables, data types, allowed values, and evented state variables.

### SOAP

The server must:

- Parse SOAP envelopes for UPnP control requests.
- Dispatch by `SOAPACTION` service URN and action name.
- Decode input arguments by action contract.
- Return XML response envelopes with action response elements.
- Return UPnP SOAP faults with standard error codes.
- Preserve HTTP status behavior expected by UPnP clients.
- Avoid chunked request bodies where client compatibility requires explicit
  content length in tests.
- Log action, renderer identity, object ID, and error class without logging
  sensitive tokens.

### Eventing

The server must:

- Support `SUBSCRIBE` and `UNSUBSCRIBE` for evented services.
- Parse callback URLs from `CALLBACK`.
- Honor and clamp `TIMEOUT`.
- Return `SID` and actual timeout.
- Send initial `NOTIFY` for ContentDirectory `SystemUpdateID`.
- Increment and publish `SystemUpdateID` when library-visible content changes.
- Remove expired subscriptions.
- Track sequence numbers and wrap according to UPnP rules.
- Provide a compatibility mode to disable or stall eventing for problematic
  renderers only when configured.

## ContentDirectory

### Tree Model

Mema should expose a stable logical library tree, not raw filesystem roots by
default.

Root containers:

- Movies
- Series
- Collections
- Recently Added
- Recently Updated
- Genres
- Years
- Profiles or Libraries, if useful for multi-root installations

Each media item should expose:

- A metadata object for the media item.
- One or more file objects when multiple files are attached.
- Optional containers for seasons and episodes.
- Optional containers for alternate versions or extras.

Object IDs must be:

- Stable across server restarts.
- Opaque to clients.
- Safe to embed in XML and URLs.
- Reversible by the server without exposing raw absolute filesystem paths.

### Browse

The server must implement:

- `BrowseDirectChildren`.
- `BrowseMetadata`.
- `StartingIndex`.
- `RequestedCount`.
- `Filter`, at least enough to omit expensive optional metadata when not
  requested.
- `SortCriteria` for supported fields.
- `NumberReturned`, `TotalMatches`, and `UpdateID`.
- UPnP 701 no-such-object errors for invalid object IDs.

### Search

The server must implement search as first-class behavior, not a stub:

- `GetSearchCapabilities` returns supported search fields.
- `Search` accepts a container, criteria, filter, starting index, requested
  count, and sort criteria.
- Supported fields should include title, class, creator/artist where present,
  genre, date/year, and Mema identifiers where useful.
- Unsupported criteria should return a protocol error, not silently scan the
  entire library.
- Search results must use the same DIDL-Lite object mapping as browse.

### Sort

The server must implement:

- `GetSortCapabilities`.
- Sorting by title.
- Sorting by date/year.
- Sorting by recently added or updated when exposed.
- Stable tie-breaking by media ID.

### DIDL-Lite Mapping

DIDL-Lite output must include:

- `DIDL-Lite` root with `dc`, `upnp`, `dlna`, and default DIDL namespaces.
- Containers with `id`, `parentID`, `restricted`, `childCount`, `searchable`,
  `dc:title`, and `upnp:class`.
- Items with `id`, `parentID`, `restricted`, `dc:title`, `upnp:class`, and one
  or more `res` resources.
- `dc:date` or year where available.
- `upnp:genre`, `upnp:artist`, `upnp:album`, `upnp:albumArtURI`, and
  `upnp:icon` when available.
- Resource attributes for protocolInfo, size, duration, bitrate, resolution,
  sample frequency, bits per sample, channel count, and color depth where
  available.
- XML escaping for all text and URL fields.

## ConnectionManager

The server must implement:

- `GetProtocolInfo` with all source protocolInfo values Mema can actually
  serve.
- `GetCurrentConnectionIDs`.
- `GetCurrentConnectionInfo`.
- Static output connections for simple clients.
- Future per-session connection tracking for active DLNA streams.

Protocol info must be generated from the same delivery profiles used by the
resource planner. Do not advertise formats that the server cannot serve.

## Media Resources

### Direct Resources

Direct resources must:

- Serve original files through the shared delivery service.
- Preserve byte range support.
- Set content type by extension and probe fallback.
- Set `Accept-Ranges`.
- Set DLNA content feature headers when requested through
  `getContentFeatures.dlna.org`.
- Avoid exposing absolute file paths in URLs.
- Support HEAD requests without starting expensive tool work.

### Remux and Transcode Resources

Remuxed and transcoded resources must:

- Be generated from `playback.StreamInfo` and device profiles.
- Support browser HLS exactly as preview does today.
- Support DLNA-friendly MPEG-TS or MP4 output where renderer profiles require
  it.
- Support time seek range headers where the output mode can honor them.
- Return `TimeSeekRange.dlna.org`, `contentFeatures.dlna.org`, and
  `transferMode.dlna.org` as applicable.
- Stop ffmpeg when clients disconnect.
- Treat broken pipes and cancelled contexts as normal client disconnects.
- Use bounded stderr capture or configured per-stream logs.
- Reuse `internal/tools` validation and execution limits.

### Thumbnails and Artwork

Artwork resources must:

- Prefer existing poster/backdrop metadata from Mema.
- Generate media thumbnails only when configured and tools are available.
- Cache generated thumbnails by media file identity and modification time.
- Advertise DLNA thumbnail profiles for JPEG and PNG variants.
- Provide fallback device icons when media thumbnails are unavailable.

### Subtitles

Subtitle resources must:

- Expose embedded and external subtitle availability in DIDL metadata where
  clients support it.
- Serve external subtitle files when compatible.
- Convert subtitles to a client-supported text format where needed and allowed.
- Reuse the existing subtitle track model and language requirements.
- Avoid advertising subtitles as playable video resources.
- Account for renderer differences because subtitle support varies heavily.

### Dynamic Streams

Dynamic streams are useful but risky. If implemented, they must:

- Be disabled by default.
- Use explicit admin-created stream definitions, not arbitrary files from the
  media directory.
- Run through allowlisted tools and argument templates.
- Use the same delivery service and logging controls.
- Never execute commands from untrusted library content.

## Device Profiles

Profiles should describe renderer capabilities:

- Friendly name or user-agent matchers.
- Direct play containers, video codecs, audio codecs, image formats, subtitle
  formats, and maximum bitrate.
- Remux targets.
- Transcode targets.
- DLNA profile names.
- Required compatibility headers.
- Eventing behavior.
- Search, sort, and thumbnail quirks.

Initial profiles:

- Generic DLNA 1.5.
- VLC and desktop clients.
- Kodi.
- Samsung TV family.
- LG TV family.
- Sony TV family.
- Chromecast-like profile when discovered through UPnP AV behavior.
- Browser profile retained for preview.

Profiles must be data-driven enough that adding a renderer does not require
rewriting server logic.

## Security

DLNA must be explicitly enabled. Defaults:

- Disabled unless the user turns it on.
- LAN-only binding.
- Interface allowlist.
- Client CIDR allowlist.
- No web session cookies accepted on DLNA endpoints.
- No admin API exposure through DLNA paths.
- No absolute paths in XML, object IDs, logs shown to clients, or URLs.
- Rate limits for discovery, SOAP actions, browse, search, thumbnails, and
  transcodes.
- Limits on concurrent transcodes and probes.
- Audit log for client IP, renderer identity, action, media item, and delivery
  method.

## Settings and UI

Settings should include:

- Enable DLNA server.
- Friendly name.
- Advertised network interfaces.
- Allowed client CIDRs.
- HTTP bind address or reuse application HTTP listener.
- Announce interval and max age.
- Enable search.
- Enable thumbnails.
- Enable subtitles.
- Enable transcodes.
- Maximum concurrent transcodes.
- Default renderer profile.
- Per-renderer profile override.
- Eventing compatibility mode.
- Diagnostics panel with discovered clients, active streams, last SOAP action,
  and last error.

## Observability

The implementation must emit:

- SSDP lifecycle events.
- SOAP action counts and errors.
- Browse and search timings.
- DIDL object counts.
- Active stream counts.
- Probe, thumbnail, remux, and transcode timings.
- Tool errors and cancellation reasons.
- Renderer profile chosen for each request.
- Delivery plan chosen for each resource request.

Logs must not leak stream tokens or absolute paths to clients.

## Testing

### Unit Tests

- SSDP packet construction and parsing.
- M-SEARCH target matching and MX delay bounds.
- Device XML generation.
- SCPD generation.
- SOAP request parsing and response envelopes.
- UPnP fault generation.
- Event subscription parsing, renewal, expiry, and sequence numbers.
- Object ID encoding and decoding.
- DIDL-Lite XML for containers, movies, episodes, artwork, subtitles, and
  multiple resources.
- Browse pagination and total counts.
- Search criteria parsing and unsupported criteria errors.
- Sort ordering.
- ProtocolInfo generation.
- DLNA content feature headers.
- Delivery planner decisions for browser and DLNA profiles.
- HLS playlist generation remains compatible with existing preview tests.

### Integration Tests

- Start the server on loopback and fetch root device XML.
- Use a UPnP client library to discover services by URL and validate SCPD.
- Call ContentDirectory `Browse` and `Search` through SOAP.
- Call ConnectionManager `GetProtocolInfo`.
- Fetch direct resources with byte ranges.
- Fetch HEAD and GET for transcode resources without starting work on HEAD.
- Cancel a transcode request and assert the tool process exits.
- Trigger a library update and assert `SystemUpdateID` changes.
- Subscribe to ContentDirectory events and receive an initial NOTIFY.

### Compatibility Matrix

Manual or automated compatibility should track:

- VLC desktop.
- VLC mobile.
- Kodi.
- Windows Media Player or Windows media sharing client.
- Samsung TV.
- LG TV.
- Sony TV.
- BubbleUPnP.
- Infuse or a similar iOS/tvOS client.

For each client, capture:

- Discovery.
- Browse.
- Search.
- Artwork.
- Direct play.
- Remux.
- Transcode.
- Seeking.
- Subtitles.
- Stop and disconnect behavior.

## Acceptance Criteria

- DLNA can be enabled from settings and appears on the local network.
- SSDP discovery works for root device and MediaServer service targets.
- Root device XML and SCPD XML validate against the expected UPnP shape.
- ContentDirectory browse exposes movies, TV shows, seasons, episodes, and
  files from the Mema library.
- ContentDirectory search returns correct paginated results.
- ConnectionManager reports only supported source protocols.
- Direct file resources can be played by at least VLC and one UPnP client.
- At least one DLNA transcode profile can be played by a renderer that cannot
  direct-play the source.
- Artwork is visible for clients that request it.
- External subtitles are advertised or served according to renderer capability.
- Browser preview still behaves exactly as before.
- VLC playlist streaming still behaves exactly as before.
- There is one shared delivery implementation for preview and DLNA resources.
- The implementation has unit and integration tests for each protocol layer.

## Rollout

1. Extract shared media delivery without changing user-visible behavior.
2. Add protocol-only DLNA server behind a disabled setting.
3. Add ContentDirectory browse over the Mema library.
4. Add resource serving through shared delivery.
5. Add search, eventing, artwork, subtitle, and renderer profile depth.
6. Expand compatibility testing and ship the setting as beta.
