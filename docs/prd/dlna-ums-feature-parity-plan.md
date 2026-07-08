# PRD: DLNA UMS Feature Parity Plan

## Status

Draft

## Clean-Room Boundary

This plan targets feature parity with the behavior categories of Universal
Media Server, not code parity. Mema must implement original code, original
profile records, and Mema-owned tests. UMS source files can be used only to
identify public behavior categories and device coverage gaps.

## Desired UMS-Like Functionality

### Renderer Recognition

Mema needs a renderer registry that can recognize devices from:

- client IP
- UPnP UUID
- friendly name
- manufacturer
- model name
- model number
- User-Agent
- extra HTTP headers
- manually assigned override

The registry should keep UMS-like sticky IP behavior. A recognized renderer
stays associated with its IP so later media resource requests with weak headers
still use the right profile. A more specific higher-priority match may upgrade
the renderer profile.

### Device-Specific Configuration

Mema needs per-device configuration layered above the seeded profile:

1. Device override by UUID.
2. Device override by IP.
3. Seeded or cloned renderer profile.
4. Generic renderer fallback.

Device overrides should support profile assignment, display name, allowed
state, delivery policy overrides, and notes.

### UPnP and DLNA Server

Feature parity requires:

- SSDP alive, byebye, and M-SEARCH response.
- Interface-aware descriptor URLs.
- Stable root device descriptor.
- ContentDirectory service.
- ConnectionManager service.
- MediaReceiverRegistrar compatibility service.
- SCPD documents for every advertised service.
- SOAP action parser, dispatcher, responses, and faults.
- Event subscriptions with SID, timeout, sequence, initial notify, unsubscribe,
  expiry, and profile-based disablement.

### Content Directory

Feature parity requires:

- Stable opaque object IDs.
- Root containers for movies, series, collections, genres, years, recently
  added, recently updated, folders, playlists, and search results where useful.
- Browse metadata and direct children.
- Search with useful title, class, date, genre, and creator fields.
- Pagination with requested count and starting index.
- Sort support where clients request supported fields.
- Child counts when safe for the renderer.
- Ability to omit or simplify fields for fragile renderers.

### DIDL-Lite Metadata

Feature parity requires:

- Standard DIDL namespaces.
- Containers and items with UPnP classes.
- Titles, dates, genres, actors, directors, albums, artists, track numbers, and
  descriptions where Mema has metadata.
- Artwork and thumbnail resources.
- Subtitle resources.
- Media resources with protocolInfo, size, duration, bitrate, sample frequency,
  channel count, resolution, and color/HDR hints where available.
- Renderer-specific metadata trimming.

### Delivery Engine

Feature parity requires one shared delivery engine for browser preview, VLC,
and DLNA:

- Direct file serving with byte ranges.
- Direct stream with MIME and DLNA headers.
- Container remux.
- Audio-only transcode.
- Full video transcode.
- HLS for clients that prefer it.
- MPEG-TS style transcode for television profiles.
- Matroska remux/transcode where compatible.
- Fast-start path for initial play.
- Seekable cache path for later range requests.
- Transcode process lifecycle, cancellation, and cleanup.
- Delivery diagnostics for chosen profile, rule, and ffmpeg mode.

### Seek Behavior

Renderer profiles need explicit seek semantics:

- byte-range seek
- time seek
- byte and time seek
- transcode time seek only
- no seek
- initial range treated as play, not seek

Mema should distinguish `Range: bytes=0-` from real seek ranges. Initial play
should stream immediately. Nonzero seek may use a cached remux/transcode file.

### Transcoding and Compatibility

Feature parity requires:

- Structured direct-play rule evaluation.
- Target output profiles per renderer.
- Optional bitrate caps.
- Optional resolution caps.
- Video level limits.
- Bit-depth and HDR constraints.
- Audio codec/channel conversion.
- Container-specific mux policies.
- MIME and DLNA profile-name translation.
- Per-renderer ffmpeg option hooks represented as safe structured settings,
  not arbitrary shell strings.

### Subtitles

Feature parity requires:

- External subtitle discovery.
- Internal subtitle exposure where supported.
- Subtitle format conversion where safe.
- SRT and WebVTT baseline support.
- ASS/SSA/MicroDVD/SAMI/TEXT support where profiles allow it.
- Subtitle streaming for transcoded video when the profile supports it.
- Option to offer subtitles as resource URLs or protocolInfo metadata.

### Images and Artwork

Feature parity requires:

- Server icons in root descriptor.
- Media artwork.
- Folder thumbnails.
- Album art profile behavior.
- JPEG fallback and padding options for fragile renderers.
- Thumbnail caching.

### Metadata and Library Behavior

Feature parity requires:

- Fast browse paths for large libraries.
- Cache invalidation when visible library content changes.
- SystemUpdateID increments.
- Optional resume metadata suppression for devices that mis-handle it.
- Optional date metadata suppression or year-only date metadata.
- Folder limits for restricted devices.

### Security and Access Control

Feature parity requires:

- DLNA disabled by default.
- Interface allowlist.
- CIDR allowlist.
- Optional per-device allow/deny.
- LAN-only defaults.
- No filesystem path disclosure.
- Safe ffmpeg argument construction.
- Bounded process execution and cleanup.

### Diagnostics

Feature parity requires:

- Recent renderer list.
- Last matched profile and match reason.
- Last SOAP action and object ID.
- Last delivery decision and ffmpeg mode.
- Active stream list.
- Event subscription list.
- SSDP announce status per interface.
- Descriptor URL list per interface.
- Profile decision trace for a selected file/client.

## Implementation Plan

### Phase 1: Profile Storage and Seed Model

Scope:

- Add renderer profile tables to the initial schema.
- Add renderer device override tables to the initial schema.
- Add seed version tables for profile defaults.
- Add JSON-capable structured columns for match rules, direct-play rules,
  delivery settings, DLNA flags, subtitles, artwork, metadata, and quirks.
- Seed generic, VLC, BubbleUPnP, Chromecast, Samsung, LG, Sony, and Kodi first.

Project rule: because Mema is not released, schema work must edit
`internal/storage/migrations/00001_initial_schema.sql` directly and reset the
development DB. Do not create follow-up migrations.

Acceptance:

- Profiles are seeded on DB reset.
- User edits survive app restart.
- Seed reset can restore one profile.
- No `dev.local.sql` changes.

### Phase 2: Profile Engine

Scope:

- Replace hardcoded `DefaultRendererProfiles()` with repository-backed
  profiles.
- Keep in-memory cache with settings reload.
- Add match scoring and priority.
- Add sticky IP and UUID association.
- Add profile upgrade when higher-priority match appears.
- Add match explanation output.

Acceptance:

- Weak LG media GET keeps LG profile after SOAP recognition.
- Manual override by IP/UUID wins.
- Diagnostics show chosen rule and fallback path.

### Phase 3: Capability Evaluator

Scope:

- Convert probe metadata into normalized codec/container fields.
- Evaluate direct-play rules.
- Decide direct, remux, audio transcode, video transcode, or HLS.
- Return reason codes for UI and tests.
- Add LG generation tests for DTS-supported and DTS-unsupported families.

Acceptance:

- Compatible files direct play.
- Unsupported audio can trigger audio-only transcode.
- Unsupported video can trigger full transcode when allowed.
- Resource advertisement matches actual delivery capability.

### Phase 4: Seeded Full Device Catalog

Scope:

- Build Mema-owned profile records for every family in
  `dlna-device-profiles-clean-room-spec.md`.
- Do not copy UMS rule lines or regexes.
- Use public device manuals, user-agent observations, and Mema test fixtures to
  create original rules.
- Mark low-confidence profiles as enabled only when matching confidence is
  high enough.

Acceptance:

- All catalog families exist as seeded profiles.
- Each seeded profile has identity, match rules, delivery settings, and at
  least one direct-play or fallback behavior.
- Confidence level is visible in the UI.

### Phase 5: Settings UI Device Profiles Panel

Scope:

- Add `Device profiles` panel to Settings > DLNA.
- Add searchable profile table.
- Add recent-device table.
- Add profile editor route or modal.
- Add clone, reset, enable/disable, import, and export actions.
- Add assignment control from recent device to profile override.
- Add decision trace viewer for selected device plus media file.

Acceptance:

- User can edit seeded profiles.
- User can assign LG TV IP/UUID to a profile.
- User can reset profile to seed.
- UI shows customized state and source version.
- No browser confirm dialogs; confirmations use modal.
- Tooltips use tooltip component, not browser title.

### Phase 6: UPnP Service Parity

Scope:

- Complete eventing.
- Complete Search.
- Complete sort and pagination semantics.
- Add service state variables.
- Add richer SCPDs.
- Add MediaReceiverRegistrar behavior used by picky TVs.
- Add descriptor compatibility variants controlled by profile flags.

Acceptance:

- VLC, BubbleUPnP, LG, Samsung, Sony, Windows Media Player, Kodi, and
  Chromecast-like clients browse without malformed responses.
- SOAP faults are protocol-correct.
- Event subscriptions are observable and expire.

### Phase 7: Delivery Parity

Scope:

- Generalize remux cache beyond LG.
- Support profile-specific transcode targets.
- Support subtitle conversion.
- Support thumbnail cache.
- Support DLNA profile-name and MIME translation.
- Add seek-by-time support where renderer profile asks for it.

Acceptance:

- Initial play never waits for full cache unless profile explicitly requires
  it.
- Nonzero seek works when cached remux/transcode exists.
- Unsupported audio can still show visible media resource when delivery can
  fix it.
- Profiles cannot advertise resources Mema cannot serve.

### Phase 8: Test Harness and Compatibility Lab

Scope:

- Add profile fixture generator.
- Add probe fixture library covering MP4, MKV, AVI, MPEG-TS, subtitles, HDR,
  DTS, AC3, EAC3, AAC, HEVC, AVC, AV1, and image formats.
- Add golden DIDL/protocolInfo tests per profile family.
- Add simulated renderer requests for sticky IP and priority upgrades.
- Add manual matrix steps for TVs and apps.

Acceptance:

- Focused Go tests cover profile matching and delivery decisions.
- Compatibility matrix records manual device outcomes.
- Regressions show which profile/rule changed behavior.

### Phase 9: Operational Hardening

Scope:

- Add profile seed upgrade strategy.
- Add background cache eviction.
- Add stream cancellation on disconnect.
- Add rate limits for SOAP and resource endpoints.
- Add metrics/log events.
- Add docs for troubleshooting device profiles.

Acceptance:

- Long-running transcodes clean up.
- Cache has size and age bounds.
- DLNA settings page can explain current server state.
- User can recover from bad profile edits by reset.

## Dependency Order

1. Storage seed model.
2. Repository-backed profile engine.
3. Capability evaluator.
4. UI profile management.
5. Full catalog seeding.
6. UPnP parity.
7. Delivery parity.
8. Compatibility harness.
9. Operational hardening.

Storage and profile engine must come before full catalog work. UI can start
after the first repository-backed profiles exist. Delivery parity depends on
capability evaluation because resource advertisement must match what Mema can
actually serve.

## Risks

- UMS profile data cannot be copied directly; Mema needs original capability
  records.
- Some old devices require quirks that can only be verified with hardware.
- Too many editable fields can make the UI unusable unless grouped carefully.
- Profile changes can break visibility if advertised resources become stricter
  than delivery can support.
- Transcode seek parity is expensive because real seek often requires cached
  output or time-seek-aware transcoding.

## First Implementation Slice

Recommended first slice:

1. Persist renderer profiles and device overrides.
2. Seed `generic`, `vlc`, `bubbleupnp`, `chromecast`, `kodi`, `samsung-tv`,
   `sony-tv`, `lg-webos`, `lg-tv-2023-plus`, and `lg-tv-2025-plus`.
3. Move current hardcoded profile fields into seed data.
4. Add sticky IP and match explanation backed by persisted profiles.
5. Add Settings > DLNA > Device profiles panel with table, recent devices, and
   override assignment.
6. Add LG profile split so 2023-era DTS-capable and 2025-era DTS-incapable
   devices can make different audio decisions.
7. Verify LG TV, VLC, and BubbleUPnP manually.
