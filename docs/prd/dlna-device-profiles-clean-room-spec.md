# PRD: DLNA Device Profiles Clean-Room Specification

## Status

Draft

## Clean-Room Boundary

This document describes the renderer profile system Mema needs for broad DLNA
compatibility. It is based on public behavior categories visible in Universal
Media Server renderer configuration files, protocol behavior, and Mema's own
debugging observations. It must not copy UMS profile expressions, codec rule
lines, comments, or implementation code.

Allowed inputs:

- Public device family names and profile names.
- The existence of profile concepts such as matchers, direct-play capability
  rules, transcode targets, seeking flags, subtitle flags, and image quirks.
- Behavior observed from real devices and Mema tests.

Forbidden inputs:

- Verbatim UMS `.conf` rule lines.
- Verbatim UMS matching expressions.
- UMS Java implementation code.
- UMS profile comments copied as requirements text.

## Goal

Make DLNA renderer profiles seeded, editable, testable, and complete enough to
represent the full UMS-supported device catalog in Mema-owned data.

## UMS Profile Coverage Summary

The current UMS renderer catalog contains 217 renderer/device profiles. The
catalog spans televisions, Blu-ray players, game consoles, mobile clients,
audio receivers, set-top boxes, Android renderers, desktop players, and DLNA
control-point style clients.

UMS profiles commonly contain these concepts:

| Concept              | Meaning in Mema                                                                                       |
| -------------------- | ----------------------------------------------------------------------------------------------------- |
| Display identity     | Human-facing renderer name and icon.                                                                  |
| Match inputs         | User-Agent, additional headers, UPnP device details, model data, and manual selection.                |
| Match priority       | More specific profiles override broad family profiles.                                                |
| Direct-play rules    | Container, video codec, audio codec, subtitle, image, and MIME compatibility.                         |
| Transcode target     | Preferred output family when direct play fails.                                                       |
| Seek behavior        | Byte seek, time seek, exclusive time seek, or no reliable seek.                                       |
| Transfer behavior    | Chunked transfer, fixed content length, and renderer-specific response headers.                       |
| DLNA flags           | ProtocolInfo, profile name, operation flags, conversion flags, and MIME translation behavior.         |
| Subtitle behavior    | External subtitle formats, internal subtitle handling, subtitle-as-source/resource behavior.          |
| Artwork behavior     | Thumbnail format, padding, resource placement, album-art profile, and folder thumbnails.              |
| Metadata behavior    | Whether date, folder, audio, video, or subtitle metadata should be exposed.                           |
| Compatibility quirks | Device-specific hacks such as versioned object IDs, tree shape, resume disabling, or bitrate halving. |

## Mema Profile Data Model

Renderer profiles should be stored as rows seeded by defaults and editable by
the user. Seed data should be versioned so Mema can add new defaults without
overwriting user changes. Allow users to reset their profiles.

### Profile Identity

Fields:

- `id`: stable Mema-owned identifier, for example `lg-webos`.
- `name`: user-visible profile name.
- `vendor`: normalized vendor or client family.
- `deviceClass`: `tv`, `bluray`, `console`, `mobile`, `desktop`, `receiver`,
  `speaker`, `settop`, `cast`, `control_point`, or `generic`.
- `source`: `mema_seed`, `user`, or `imported_clean_room`.
- `sourceVersion`: seed version used to create the row.
- `enabled`: whether automatic matching may use this profile.
- `priority`: integer; higher value wins when several profiles match.
- `iconKey`: Mema-owned icon key, not copied asset.
- `notes`: user-editable compatibility notes.

### Match Rules

Fields:

- `userAgentRules`: user-editable contains/regex rules.
- `headerRules`: header-name plus contains/regex rules.
- `upnpDetailRules`: friendly name, manufacturer, model name, model number,
  device description, and vendor extension rules.
- `clientHints`: optional manually learned tokens from observed clients.
- `matchMode`: `any`, `all`, or weighted score.
- `minScore`: threshold for weighted matching.

Matching order:

1. Manual device override by UUID.
2. Manual device override by IP.
3. Existing connected renderer association by IP.
4. UPnP UUID and device details.
5. HTTP headers.
6. Default profile.

Mema should keep UMS-like sticky IP behavior: once a renderer is recognized at
an IP, weak later media requests from the same IP reuse that renderer unless a
higher-priority profile is positively matched.

### Direct-Play Capability Rules

Capability rules should be structured records, not free-form UMS syntax.

Fields:

- `mediaKind`: `video`, `audio`, `image`, or `subtitle`.
- `containers`: allowed containers.
- `videoCodecs`: allowed video codecs.
- `audioCodecs`: allowed audio codecs.
- `subtitleFormats`: allowed subtitle formats.
- `mimeType`: served MIME type.
- `maxWidth`, `maxHeight`, `maxBitrateMbps`, `maxLevel`, `bitDepths`.
- `hdrModes`: SDR, HDR10, HDR10+, HLG, Dolby Vision.
- `audioChannels`: optional channel limits.
- `constraints`: structured flags such as no GMC, no QPEL, no unsupported
  packed bitstream, or no odd-dimension muxing.

Rules are evaluated against Mema probe metadata. A file is direct-playable only
when at least one enabled rule matches container, codecs, dimensions, bitrate,
and required constraints.

### Delivery Preferences

Fields:

- `preferredProtocol`: `direct`, `remux`, `transcode`, or `hls`.
- `avoidHLS`: hides HLS resources from clients that mis-handle them.
- `preferHLS`: orders HLS first for HLS-first clients.
- `transcodeVideoTarget`: structured output target such as MPEG-TS/H.264/AC3.
- `transcodeAudioTarget`: structured output target such as AAC stereo or AC3.
- `remuxTargets`: ordered container remux options.
- `forceTranscodeExtensions`: extensions that should not direct play.
- `streamExtensions`: extensions that may direct stream without transcode.
- `allowAudioOnlyTranscode`: permits copy-video plus transcode-audio.
- `allowVideoTranscode`: permits full video transcode.

### DLNA HTTP and Protocol Flags

Fields:

- `seekMode`: `byte`, `time`, `both`, `time_exclusive`, or `none`.
- `chunkedTransfer`: allow chunked resource responses.
- `sendContentLength`: force content length where known.
- `sendDlnaOrgFlags`: include DLNA.ORG flags in protocolInfo.
- `dlnaOrgPnMode`: none, coarse, or accurate.
- `contentFeaturesHeaderMode`: none, direct-only, or all-resources.
- `transferModeHeader`: optional `Streaming`, `Interactive`, or omitted.
- `mimeOverrides`: Mema-owned MIME translations.
- `profileNameOverrides`: Mema-owned DLNA profile-name translations.

### Subtitle Rules

Fields:

- `externalSubtitleFormats`.
- `internalSubtitleFormats`.
- `streamSubtitlesForTranscodedVideo`.
- `offerSubtitlesAsResource`.
- `offerSubtitlesByProtocolInfo`.
- `subtitleHttpHeaderMode`.
- `removeUnsupportedSubtitleTags`.

### Artwork and Image Rules

Fields:

- `thumbnailAsResource`.
- `albumArtProfile`.
- `forceJpegThumbnails`.
- `thumbnailPadding`.
- `sendFolderThumbnails`.
- `imageMimeOverrides`.

### Metadata and Tree Rules

Fields:

- `contentTreeMode`: standard, fast, flat, or folder-biased.
- `needVersionedObjectId`.
- `disableServerResume`.
- `pushMetadata`.
- `sendDateMetadata`.
- `sendAudioMetadata`.
- `sendVideoMetadata`.
- `sendSubtitleMetadata`.
- `prependTrackNumbers`.
- `limitFolders`.

## Supported Device Catalog

The seeded catalog should cover these UMS-supported profile families using
Mema-owned capability records.

| Family                                                          | Profiles to seed                                                                                                                                                                                                                |
| --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Generic                                                         | DefaultRenderer                                                                                                                                                                                                                 |
| Amazon                                                          | Amazon Fire TV Stick Vimu Player                                                                                                                                                                                                |
| Android and Chromecast                                          | Android, BubbleUPnP, Chromecast, Chromecast Ultra                                                                                                                                                                               |
| Apple                                                           | iPad/iPhone, AirPlayer, VLC for iOS/macOS, VLC for older iOS, VLC for Apple TV, VLC for Apple TV 4K                                                                                                                             |
| Desktop and apps                                                | VLC for desktop, Kodi, XBMC, Movian, MediaPlayer, foobar2000 mobile, Bigscreen, Pigasus, Skybox VR Player                                                                                                                       |
| BlackBerry and Nokia                                            | KalemSoft Media Player, Nokia N900                                                                                                                                                                                              |
| Bush, DirecTV, Fetch, Freebox, Netgem, Telstra, Thomson         | Bush Freeview Set Top Box, DirecTV, Fetch TV, Freebox, Netgem N7700, Telstra T-Box, Thomson U3 Series                                                                                                                           |
| Cambridge, Denon, Onkyo, Yamaha                                 | Azur BD, Denon AVR-4311CI, Denon AVR-X4200W, Onkyo TX-NR7xx, Onkyo TX-NR8xx, Yamaha AV Receiver, Yamaha R-N303, Yamaha R-N500, Yamaha RX-A1010, Yamaha RX-A2050, Yamaha RX-V3900, Yamaha RX-V500D, Yamaha RX-V671               |
| D-Link, Netgear, Western Digital, Realtek                       | D-Link DSM-510, Netgear NeoTV, WD TV Live, Realtek                                                                                                                                                                              |
| Freecom, Hama, Linkplay, Linn, Logitech, Lumin, Naim, Technisat | Freecom MusicPal, Hama IR320, WiiM, Linn App, Logitech Squeezebox, Lumin, Lumin U1 Mini, Naim Mu-So, Naim Mu-So Qb, Technisat S1+                                                                                               |
| Hisense, Sharp, Telefunken, Vizio, VideoWeb                     | Hisense K680, Sharp Aquos, Telefunken TV, Vizio Smart TV, VideoWeb TV                                                                                                                                                           |
| LG Blu-ray and legacy                                           | LG Blu-ray BDP, LG Blu-ray BP, LG BP550, LG Smart TV Upgrader, LG EG910V, LG LA6200, LG LA644V, LG LB/LCD 2014, LG LM620, LG LM660, LG LS5700, LG UB820V, LG UH770                                                              |
| LG modern TV                                                    | LG LED LCD, LG LED LCD 2022+, LG WebOS TV, LG NANO TV, LG OLED, LG OLED 2020/2021, LG OLED 2022, LG TV 2023+, LG TV 2025+                                                                                                       |
| Microsoft                                                       | Windows Media Player, Xbox 360, Xbox One                                                                                                                                                                                        |
| Miracast and AnyCast                                            | AnyCast, Miracast M806, Mirascreen                                                                                                                                                                                              |
| OPPO and Pioneer                                                | OPPO BDP, OPPO BDP-83, Pioneer BDP, Pioneer Kuro                                                                                                                                                                                |
| Panasonic Blu-ray and recorder                                  | Panasonic BDT, BDT220, BDT360, DMR, Home Theater SC-BTT                                                                                                                                                                         |
| Panasonic TV                                                    | Panasonic TV, AS600, AS650, CX680, CX700, DX, E6, ET60, GT50, GX800B, HZ1500, S60, ST60, U30Z, TX-L32V10E, VT60                                                                                                                 |
| Philips                                                         | Philips Android TV, Philips Aurea, Philips TV, Philips PUS TV, Philips 6500 Series TV, Streamium                                                                                                                                |
| Popcorn Hour and Showtime                                       | Popcorn Hour, Showtime 3, Showtime 4                                                                                                                                                                                            |
| Roku                                                            | Roku DVP 10, Roku 3 NSP 3, Roku 3 NSP 5, Roku 3 NSP 6-7, Roku 4 NSP 6-7, Roku TV, Roku TV 4K, Roku TV NSP 8, Roku Ultra                                                                                                         |
| Samsung Blu-ray, home theater, and audio                        | Samsung BD-C6800, Samsung H6500, Samsung HT-E3, Samsung HT-F4, Samsung Soundbar, Samsung Soundbar MS750                                                                                                                         |
| Samsung mobile                                                  | Samsung Mobile, Galaxy S5, Galaxy S7, Note Tab                                                                                                                                                                                  |
| Samsung TV legacy                                               | Samsung WiseLink, C/D Series, C6600, D6400, D7000, 5300 Series, PL51E490, EH5300, EH6070, ES6100, ES6575, ES8000, ES8005, F5100, F5505, F5900, H4500, H6203, H6400, J55xx, J6200, E+ Series                                     |
| Samsung TV modern                                               | Samsung UHD, LED UHD, 8 Series, 9 Series, OLED, Q7 Series, Q9 Series, The Frame, SMT-G7400, 2018 QLED TV, UHD TV 2019+, 8K TV 2019+, 2021+ QLED TV, 2021 AU9/Q6/43Q7/50Q7, 2021 AU8/AU7/BEA/32Q6, 2021 Q5, 2021+ NEO QLED TV 8K |
| Sony Blu-ray, consoles, mobile, receiver                        | Sony Blu-ray, Sony Blu-ray 2013, BDP-S3700, UBP-X800M2, PlayStation 3, PlayStation 4, PlayStation Vita, Sony Home Theatre System, Sony SA-NS310, Sony SMP-N100, Sony STR-DN1080, Sony STR-DA5800ES, Sony Xperia, Xperia Z3      |
| Sony Bravia                                                     | Sony Bravia, 4500, 5500, AG, BX305, EX, EX620, EX725, HX, HX75, NX70x, NX800, W, X, X Series TV, XBR, XBR OLED, XD/XE/XF, XH, XR                                                                                                |

## LG Clean-Room Behavior

LG should not be one profile. Mema should seed a broad LG webOS fallback plus
model-era profiles.

Required LG behavior:

- Match broad webOS/LG tokens for older and unknown LG TVs.
- Prefer model-specific LG profile when UPnP details include model family.
- Treat 2023-era LG models separately from 2025-era LG models because DTS
  support differs by generation.
- Avoid HLS by default for LG DLNA playback unless a specific profile proves it.
- Prefer visible direct MKV resources when compatible.
- Use audio-only remux/transcode when video is compatible but audio is not.
- Preserve initial play speed by streaming first and building seek cache only
  when needed.
- Support SRT and WebVTT at minimum, then add richer subtitle formats when
  profile rules prove them.

## Settings UI Requirements

Add a new panel in Settings > DLNA named `Device profiles`.

Panel sections:

- `Seeded profiles`: searchable table with name, family, enabled state,
  priority, source version, and modified state.
- `Recent devices`: observed IP, UUID, friendly name, matched profile,
  last-seen time, and override selector.
- `Profile editor`: modal or route for identity, matching, direct play,
  delivery, DLNA flags, subtitles, artwork, metadata, and quirks.
- `Seed management`: reset one profile to seed, clone seeded profile, disable
  profile, and import/export JSON.
- `Diagnostics`: show why a request matched a profile and which rule accepted
  or rejected a media file.

User edits must never be overwritten by seed updates. Seed updates create a new
seed version and mark profiles as `customized` when user data differs.

## Acceptance Criteria

- Renderer profiles are persisted and editable.
- All seeded profiles can be disabled, cloned, or reset.
- Recent DLNA clients can be assigned a manual profile override.
- Matching explains profile choice in diagnostics.
- Capability evaluation explains direct/remux/transcode decision.
- Mema can seed a profile catalog covering every family listed above.
- No UMS `.conf` line is copied into seed data or docs.
