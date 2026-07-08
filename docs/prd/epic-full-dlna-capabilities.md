# Epic: Full DLNA and UPnP Capabilities

## Status

Draft

## Epic Goal

Implement full DLNA and UPnP AV capabilities in Mema while preserving current
browser preview, VLC streaming, media probing, and transcoding behavior through
one shared media delivery system.

## Success Criteria

- Mema advertises as a DLNA Digital Media Server on configured LAN interfaces.
- UPnP clients can discover the server, read device descriptions, inspect SCPD
  documents, and call SOAP actions.
- ContentDirectory browse and search expose Mema library media with useful
  metadata, artwork, resources, and subtitles.
- ConnectionManager reports protocols that Mema can actually serve.
- Direct, remuxed, transcoded, thumbnail, artwork, and subtitle resources use
  the same delivery stack as browser preview.
- Existing preview endpoints and VLC playlist behavior continue to pass their
  current tests.
- DLNA is disabled by default, LAN-scoped when enabled, observable, and
  rate-limited.

## Story DLNA-001: Research Baseline and Protocol Contracts

Goal: Capture the protocol contract and clean-room implementation boundaries.

Scope:

- Keep cloned sources under `.research`.
- Document behavior categories from `anacrolix/dms`, `huin/goupnp`, and
  `koron/go-ssdp`.
- Define which behavior is protocol-derived and which behavior is a
  renderer-specific compatibility target.
- Create golden XML examples for device, SCPD, SOAP, and DIDL-Lite outputs.

Acceptance criteria:

- `.research/` remains git ignored.
- No copied source code is introduced.
- The DLNA clean-room spec exists under `docs/prd`.
- The device profile clean-room spec exists under `docs/prd`.
- The UMS feature parity plan exists under `docs/prd`.
- Golden protocol examples are generated from Mema-owned fixtures.

Tests:

- Documentation link check or file existence check.
- Golden fixture generation test once protocol packages exist.

Dependencies: None.

## Story DLNA-002: Extract Shared Media Delivery Service

Goal: Move current preview and stream behavior behind one reusable service.

Scope:

- Extract media path resolution adapters from HTTP handlers.
- Extract direct file serving with range support.
- Extract ffprobe stream, chapter, container, duration, and keyframe probing.
- Extract HLS playlist and MPEG-TS segment generation.
- Keep existing browser and WebKit preview decisions.
- Keep existing VLC stream token and M3U behavior.
- Return delivery plans that can be consumed by preview and DLNA.

Acceptance criteria:

- Preview info returns the same fields as before.
- Preview direct playback still uses direct file serving when eligible.
- Preview HLS playlists have the same segment behavior as before.
- Preview segment ffmpeg args are unchanged unless tests are updated for a
  deliberate improvement.
- VLC playlist streaming still works.
- No DLNA package runs its own ffprobe or ffmpeg execution path.

Tests:

- Existing preview tests.
- New delivery service unit tests for direct, remux, and transcode plans.
- New direct file range tests.

Dependencies: DLNA-001.

## Story DLNA-003: DLNA Settings and Lifecycle

Goal: Add disabled-by-default server configuration and lifecycle management.

Scope:

- Add settings for enabled state, friendly name, interfaces, allowed CIDRs,
  announce interval, transcode enablement, thumbnails, subtitles, and renderer
  default profile.
- Start and stop the DLNA server with the application.
- Support clean shutdown and restart when settings change.
- Show server status and last error in settings.

Acceptance criteria:

- DLNA is disabled by default.
- Enabling DLNA starts the server without affecting the app HTTP API.
- Disabling DLNA sends shutdown notifications when SSDP exists.
- Invalid interface or CIDR settings are rejected with clear errors.

Tests:

- Settings validation tests.
- Lifecycle start, stop, and restart tests with fake network services.

Dependencies: DLNA-002.

## Story DLNA-004: SSDP Advertiser and Discovery Responder

Goal: Implement SSDP multicast advertisement and M-SEARCH responses.

Scope:

- Join multicast on configured interfaces.
- Send alive notifications for root device, UUID, MediaServer, and services.
- Send byebye notifications on shutdown.
- Respond to `ssdp:all`, `upnp:rootdevice`, UUID, device type, and service
  type searches.
- Generate interface-correct `LOCATION` URLs.
- Handle IPv4 and supported IPv6 SSDP.
- Add compatibility parsing for safe malformed packet cases.

Acceptance criteria:

- A local SSDP search can discover the server.
- Responses include valid CACHE-CONTROL, EXT, LOCATION, SERVER, ST, and USN.
- Alive and byebye messages are emitted for all advertised targets.
- Interface allowlist limits announcement and response surfaces.

Tests:

- Packet unit tests.
- Loopback or fake UDP integration tests.
- Search response target matching tests.

Dependencies: DLNA-003.

## Story DLNA-005: Device Description and SCPD Documents

Goal: Serve UPnP root device and service descriptor XML.

Scope:

- Generate stable root device XML.
- Serve ContentDirectory v1 SCPD.
- Serve ConnectionManager v1 SCPD.
- Optionally serve MediaReceiverRegistrar compatibility SCPD.
- Include icons and presentation URL.
- Add XML golden tests.

Acceptance criteria:

- `GET /dlna/rootDesc.xml` returns valid XML.
- SCPD URLs from root description return valid XML.
- Service control and event URLs are stable.
- Device UDN remains stable across restarts for the same server identity.

Tests:

- XML golden tests.
- URL resolution tests.
- UPnP client parse tests.

Dependencies: DLNA-004.

## Story DLNA-006: SOAP Engine and UPnP Errors

Goal: Implement generic UPnP SOAP request handling.

Scope:

- Parse SOAPACTION.
- Parse SOAP envelopes and action arguments.
- Dispatch to registered services.
- Marshal action responses.
- Marshal UPnP SOAP faults.
- Normalize standard error codes.
- Log action metadata.

Acceptance criteria:

- Unknown actions return UPnP invalid-action faults.
- Invalid arguments return argument errors.
- Service handlers do not manually build SOAP envelopes.
- SOAP tests cover success and fault cases.

Tests:

- SOAP parser unit tests.
- Fault XML golden tests.
- Action dispatch tests.

Dependencies: DLNA-005.

## Story DLNA-007: ContentDirectory Object Model

Goal: Map the Mema library to stable UPnP containers and objects.

Scope:

- Define opaque object IDs.
- Implement root containers for movies, TV shows, collections, recently added,
  recently updated, genres, and years.
- Map movies, shows, seasons, episodes, and media files.
- Include parent IDs and child counts.
- Hide missing or inaccessible files.

Acceptance criteria:

- Object IDs are stable and do not expose filesystem paths.
- Root browse returns expected containers.
- Media item browse returns expected child objects.
- Missing files are skipped or exposed as unavailable according to spec.

Tests:

- Object ID encode/decode tests.
- Library fixture browse tests.
- Child count tests.

Dependencies: DLNA-006.

## Story DLNA-008: DIDL-Lite Metadata and Resource Mapping

Goal: Serialize ContentDirectory objects into DLNA-compatible DIDL-Lite.

Scope:

- Add DIDL-Lite structs and serializers.
- Map title, class, date, genre, artist, album, artwork, size, duration,
  bitrate, resolution, channels, and sample frequency.
- Generate resources from the shared delivery service.
- Escape all XML and URL fields.

Acceptance criteria:

- Containers and items serialize with required namespaces.
- Media items include direct and compatible alternate resources.
- Resource attributes reflect actual probe metadata when available.
- DIDL output is deterministic for tests.

Tests:

- DIDL golden tests.
- XML escaping tests.
- Resource attribute tests.

Dependencies: DLNA-007, DLNA-002.

## Story DLNA-009: ContentDirectory Browse

Goal: Implement `BrowseMetadata` and `BrowseDirectChildren`.

Scope:

- Decode browse arguments.
- Support starting index and requested count.
- Respect filter where practical.
- Support sort criteria for known fields.
- Return NumberReturned, TotalMatches, and UpdateID.
- Return correct UPnP errors for invalid objects.

Acceptance criteria:

- UPnP clients can browse root and media containers.
- Pagination works without changing total count.
- Browse metadata returns exactly one object for valid IDs.
- Invalid IDs return 701 no such object.

Tests:

- SOAP Browse tests.
- Pagination tests.
- Invalid object tests.

Dependencies: DLNA-008.

## Story DLNA-010: ContentDirectory Search

Goal: Implement first-class search over the Mema library.

Scope:

- Return supported search capabilities.
- Parse supported criteria.
- Search titles, classes, genre, year/date, and creator fields where present.
- Support pagination and sorting.
- Reject unsupported criteria clearly.

Acceptance criteria:

- Search returns matching movies and episodes.
- Search results use the same DIDL mapping as browse.
- Unsupported criteria do not trigger unbounded scans.

Tests:

- Criteria parser tests.
- Search fixture tests.
- SOAP Search tests.

Dependencies: DLNA-009.

## Story DLNA-011: ConnectionManager

Goal: Report available source protocols and basic connection status.

Scope:

- Implement `GetProtocolInfo`.
- Implement `GetCurrentConnectionIDs`.
- Implement `GetCurrentConnectionInfo`.
- Generate Source protocol info from delivery profiles.
- Keep Sink empty for server-only mode.

Acceptance criteria:

- ProtocolInfo only advertises formats Mema can serve.
- ConnectionManager SOAP calls work through a UPnP client.
- Static connection responses are accepted by common renderers.

Tests:

- ProtocolInfo generation tests.
- SOAP ConnectionManager tests.

Dependencies: DLNA-006, DLNA-008.

## Story DLNA-012: DLNA Resource Serving

Goal: Serve direct, remuxed, and transcoded media resources to DLNA clients.

Scope:

- Add DLNA resource URL signing or opaque resource IDs.
- Serve direct files through shared delivery.
- Serve remux/transcode streams through shared delivery.
- Support HEAD requests without starting transcodes.
- Implement DLNA seek and content feature headers.
- Enforce concurrent stream and transcode limits.

Acceptance criteria:

- Direct resources support byte ranges.
- DLNA transcode resources stop when clients disconnect.
- HEAD on a transcode resource returns headers only.
- Browser preview still uses the same HLS behavior as before.

Tests:

- Resource GET and HEAD tests.
- Range tests.
- Cancellation tests.
- Existing preview tests.

Dependencies: DLNA-002, DLNA-011.

## Story DLNA-013: Artwork and Thumbnail Resources

Goal: Expose posters, backdrops, and generated thumbnails to DLNA clients.

Scope:

- Prefer metadata artwork already known by Mema.
- Generate thumbnails from media files when configured.
- Cache thumbnails by media file identity and modification time.
- Advertise JPEG and PNG thumbnail resources.
- Serve fallback icons.

Acceptance criteria:

- DIDL includes album art or icon URLs when available.
- Clients can fetch artwork resources without authentication cookies.
- Thumbnail generation respects tool limits and concurrency.

Tests:

- Artwork URL tests.
- Thumbnail cache tests.
- Fallback icon tests.

Dependencies: DLNA-008, DLNA-012.

## Story DLNA-014: Subtitle Exposure and Conversion

Goal: Expose compatible subtitle resources for embedded and external subtitles.

Scope:

- Map Mema subtitle tracks and external subtitle files to DLNA metadata.
- Serve external subtitles when compatible.
- Convert supported subtitle formats to configured text formats where allowed.
- Honor renderer profile subtitle capabilities.
- Avoid advertising unsupported subtitles for known-incompatible renderers.

Acceptance criteria:

- External SRT subtitles can be served to clients that support them.
- Unsupported subtitle formats are omitted or converted according to profile.
- Subtitle language metadata is preserved.

Tests:

- Subtitle DIDL mapping tests.
- Subtitle resource serving tests.
- Conversion planner tests.

Dependencies: DLNA-008, DLNA-012.

## Story DLNA-015: Eventing and System Update IDs

Goal: Notify subscribed clients when library-visible content changes.

Scope:

- Implement ContentDirectory event subscriptions.
- Send initial `SystemUpdateID`.
- Increment update ID on library changes.
- Notify active subscribers.
- Expire old subscribers and support unsubscribe.
- Add compatibility mode for problematic renderers.

Acceptance criteria:

- SUBSCRIBE returns SID and TIMEOUT.
- Initial NOTIFY is sent.
- Update notification is sent after a library change.
- UNSUBSCRIBE removes the subscriber.

Tests:

- Subscription lifecycle tests.
- Notify HTTP callback tests.
- Update ID tests.

Dependencies: DLNA-006, DLNA-009.

## Story DLNA-016: Renderer Profiles and Compatibility Overrides

Goal: Make renderer-specific behavior data-driven.

Scope:

- Add profile matching by user-agent, friendly name, headers, and client IP.
- Define initial generic, VLC, Kodi, Samsung, LG, Sony, and Chromecast-like
  profiles.
- Support per-renderer admin override.
- Drive protocolInfo, DIDL resources, eventing behavior, and delivery plans
  from profiles.

Acceptance criteria:

- Generic profile works without overrides.
- Known clients receive tailored resources and headers.
- Admin can override a renderer profile without restarting the app.

Tests:

- Profile match tests.
- ProtocolInfo per profile tests.
- Delivery plan per profile tests.

Dependencies: DLNA-012, DLNA-015.

## Story DLNA-017: DLNA Diagnostics UI

Goal: Provide enough UI to operate and debug DLNA.

Scope:

- Add settings UI for DLNA enablement and server options.
- Show server status, bound interfaces, advertised URLs, and last SSDP event.
- Show recent clients, selected profile, last SOAP action, and last error.
- Show active streams and active transcodes.
- Add manual restart action.

Acceptance criteria:

- User can enable DLNA and see advertised URLs.
- User can identify which renderer profile was selected.
- User can see active streams and stop/restart the server if needed.
- Confirmation actions use modals and tooltips, not browser dialogs.

Tests:

- Svelte component tests for settings states.
- API tests for diagnostics payload.

Dependencies: DLNA-003, DLNA-016.

## Story DLNA-018: Security, Limits, and Audit Logging

Goal: Keep DLNA local, bounded, and auditable.

Scope:

- Enforce interface and CIDR allowlists.
- Add discovery, SOAP, browse, search, probe, thumbnail, and transcode limits.
- Add active stream limits.
- Add audit events for client, action, media, delivery method, and result.
- Ensure stream tokens and absolute paths are not logged or exposed.

Acceptance criteria:

- Requests outside allowed CIDRs are rejected.
- Limits return clear errors and do not start expensive work.
- Audit entries exist for each stream and SOAP action.
- Tests prove object IDs and URLs do not contain absolute paths.

Tests:

- Allowlist tests.
- Rate and concurrency limit tests.
- Audit log tests.
- Path leakage tests.

Dependencies: DLNA-012, DLNA-017.

## Story DLNA-019: Compatibility Test Matrix

Goal: Validate real clients and document known behavior.

Scope:

- Add a compatibility checklist for VLC, Kodi, Windows, Samsung, LG, Sony,
  BubbleUPnP, and iOS/tvOS clients.
- Record discovery, browse, search, artwork, direct play, remux, transcode,
  seeking, subtitles, stop, and disconnect behavior.
- Add reproducible test media fixtures for DLNA.
- Add troubleshooting docs.

Acceptance criteria:

- At least VLC and one UPnP client pass automated or reproducible manual tests.
- Known limitations are documented.
- Compatibility findings feed back into renderer profiles.

Tests:

- Automated integration tests where clients can run headless.
- Manual checklist stored in docs.

Dependencies: DLNA-016, DLNA-018.

## Story DLNA-020: Optional MediaRenderer and Control Point Capabilities

Goal: Add renderer/control capabilities only after server mode is stable.

Scope:

- Discover MediaRenderer devices as a control point.
- Support AVTransport actions such as SetAVTransportURI, Play, Pause, Stop,
  Seek, Next, and Previous.
- Support RenderingControl actions such as volume and mute.
- Optionally expose Mema as a MediaRenderer if a local playback target exists.

Acceptance criteria:

- Server-only DLNA is complete before this story starts.
- Control point actions are available only through explicit UI.
- Renderer discovery does not interfere with MediaServer advertisement.

Tests:

- UPnP client contract tests.
- Fake renderer SOAP tests.

Dependencies: DLNA-019.

## Implementation Order

1. DLNA-001
2. DLNA-002
3. DLNA-003
4. DLNA-004
5. DLNA-005
6. DLNA-006
7. DLNA-007
8. DLNA-008
9. DLNA-009
10. DLNA-011
11. DLNA-012
12. DLNA-010
13. DLNA-013
14. DLNA-014
15. DLNA-015
16. DLNA-016
17. DLNA-017
18. DLNA-018
19. DLNA-019
20. DLNA-020

## Current Preview Integration Constraint

No story may introduce a DLNA-only implementation of:

- ffprobe metadata extraction.
- keyframe discovery.
- file range serving.
- HLS playlist generation.
- ffmpeg segment execution.
- direct/remux/transcode decision logic.
- content type detection.

If DLNA needs a new output format, it must extend the shared delivery service
and then become available to any caller whose profile allows it.
