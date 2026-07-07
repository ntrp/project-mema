# DLNA compatibility matrix

This checklist records renderer behavior for the DLNA server and the media fixture cases in `internal/dlna/testdata/compatibility_media.json`.

## Fixture set

| Fixture | Purpose | Expected delivery |
| --- | --- | --- |
| baseline-mp4-h264-aac | Common MP4 direct play path | Direct file |
| matroska-remux-target | MKV container compatibility | Remuxed stream |
| webkit-hls-transcode-target | HLS-only clients | Transcoded HLS |
| external-subtitle-sidecar | External subtitle exposure | SRT or VTT subtitle resource |

## Client checklist

| Client | Profile | Discovery | Browse | Search | Artwork | Direct | Remux | Transcode | Seeking | Subtitles | Stop/disconnect | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| VLC | vlc | Pass | Pass | Pass | Pass | Pass | Manual | Manual | Manual | Pass | Manual | Automated profile test plus manual playback checklist |
| BubbleUPnP | bubbleupnp | Pass | Pass | Pass | Pass | Pass | Manual | Manual | Manual | Pass | Manual | Automated profile test plus manual playback checklist |
| Kodi | kodi | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Known client profile covered |
| Windows Media Player | generic | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Generic profile fallback |
| Samsung TV | samsung | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | SRT | Manual | Samsung profile covered |
| LG TV | lg | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | SRT | Manual | LG profile covered |
| Sony TV | sony | Manual | Manual | Manual | Manual | Manual | Manual | Manual | Manual | SRT | Manual | Sony profile covered |
| iOS/tvOS cast target | chromecast | Manual | Manual | Manual | Manual | HLS | HLS | HLS | Manual | VTT | Manual | HLS-first, eventing disabled |

## Manual run steps

1. Enable DLNA in Settings > DLNA and confirm the advertised URL appears.
2. Start the client on the same allowed LAN or add its subnet to Allowed CIDRs.
3. Confirm the Mema server appears in discovery within one announce interval.
4. Browse Movies, Series, and a media item with artwork and subtitles.
5. Play each fixture case and record the selected profile in Settings > DLNA.
6. Seek forward, stop playback, disconnect the client, and confirm active streams return to zero.
7. Check System > Events for DLNA SOAP and stream audit entries.

## Known limitations

- Remux and transcode results depend on the host `ffmpeg` build.
- Eventing is disabled for Chromecast-like clients because they do not need ContentDirectory event subscriptions.
- Some televisions cache stale SSDP records until their network stack is restarted.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| Client cannot discover Mema | Confirm DLNA is enabled, the server is bound to the expected interface, multicast is allowed, and the client IP is inside Allowed CIDRs. |
| Browse works but playback fails | Compare the selected renderer profile, protocolInfo order, and the fixture delivery mode. |
| Artwork or subtitles are missing | Confirm the media item has artwork or sidecar subtitles and that the client supports the profile subtitle format. |
| Streams stay active after stopping | Refresh Settings > DLNA and check audit events for disconnect or transcode errors. |
| Absolute paths appear in output | Treat as a bug; object IDs, resource URLs, and audit entries must stay opaque. |
