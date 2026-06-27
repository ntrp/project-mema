# PRD: Download Clients and Import Pipeline

Status: Draft

## Summary

Mema must integrate with download clients, monitor downloads, import completed files, analyze media streams, and route files into final library locations. The import pipeline is the foundation for both ordinary ARR-style imports and advanced component assembly.

## Goals

- Support multiple download clients.
- Route grabs by media type, profile, category, and library.
- Track download lifecycle.
- Import completed files safely.
- Analyze files before final placement.
- Support copy, move, and hardlink.
- Support failed download handling and blocklisting.
- Support partial import as components for later assembly.

## Candidate Download Clients

- Transmission, first torrent client
- qBittorrent, later
- Deluge
- SABnzbd, first NZB client
- NZBGet
- aria2, later

## Functional Requirements

- Users can add, test, edit, disable, and delete download clients.
- Users can assign clients to media types, libraries, and protocols.
- Mema can send releases to a selected client with category/tag metadata.
- Mema monitors active downloads.
- Mema detects completed downloads.
- Mema scans completed folders for importable files.
- Mema rejects samples, trailers, ads, tiny files, and unsupported formats.
- Mema analyzes video, audio, subtitle, ebook, audiobook, and music files.
- Mema imports files using copy, move, or hardlink.
- Mema supports manual import for unmatched files.
- Mema marks failed downloads and can search alternatives.
- Mema can retain downloaded files as source components for later muxing.

## Import States

- Grabbed
- Downloading
- Completed
- Scanning
- Import candidate
- Imported
- Component retained
- Assembling
- Assembled
- Failed
- Blocklisted
- Manual review required

## Acceptance Criteria

- Mema can send a release to a configured download client.
- Mema can monitor the release until completion.
- Mema can import the largest valid media file from a completed download.
- Mema records technical media info for imported files.
- Mema can reject a download and explain the reason.
- Mema can keep downloaded media as a component without final import.

## Open Questions

- First-release clients are Transmission for torrents and SABnzbd for NZBs.
- Should torrents default to hardlink import?
- Should Mema manage seeding goals, ratio limits, and post-import deletion?
- Should completed download handling be polling-based, webhook-based, or both?
- Should users be able to define per-client path mappings?
- How should remote path mappings work across Docker hosts and NAS mounts?
- Should failed download handling be automatic for all clients?
- Should imports run in worker queues with concurrency limits?
- Should archive extraction be built in?
- Should encrypted/packed releases be supported?
- Should Mema import extras and bonus files?
- Should partial component sources be retained forever, until assembled, or by retention policy?
