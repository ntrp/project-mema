# PRD: Books and Music Management

Status: Draft

## Summary

Mema should eventually cover book, audiobook, music artist, album, and track workflows. These domains share acquisition and library workflows with video, but metadata, completeness, and file formats differ. They are explicitly deferred until after the video-first MVP is proven.

## Goals

- Manage authors, books, editions, ebook files, and audiobook files.
- Manage artists, albums, releases, discs, and tracks.
- Support metadata enrichment and tagging.
- Support quality profiles for ebook, audiobook, and music formats.
- Support download, import, rename, and upgrade workflows.

## Books Functional Requirements

- Users can add authors and books.
- Users can monitor all books, selected books, or future books by author.
- Users can manage editions and formats.
- Mema can import ebooks and audiobooks.
- Mema can identify ISBNs and provider IDs.
- Mema can rename and organize by author, series, book, edition, and format.
- Mema can track missing formats separately.

## Music Functional Requirements

- Users can add artists.
- Users can monitor all albums, selected albums, or future albums by artist.
- Users can manage release groups, album versions, discs, and tracks.
- Mema can import FLAC, MP3, AAC, Opus, and other audio files.
- Mema can inspect audio codec, bitrate, sample rate, bit depth, and channels.
- Mema can tag files using metadata.
- Mema can organize by artist, album, disc, and track.
- Mema can distinguish albums, singles, EPs, live releases, compilations, and soundtracks.

## Acceptance Criteria

- A user can add an author and monitor selected books.
- A user can add a music artist and monitor selected albums.
- Mema can import a completed ebook or music download.
- Mema can score ebook/music formats against profiles.
- Mema can show missing books, albums, and tracks.

## Open Questions

- Books and music are planned after video domains.
- Should ebooks and audiobooks be separate libraries?
- Which ebook formats should be supported: EPUB, PDF, MOBI, AZW3, CBZ, CBR?
- Should comics/manga be included?
- Should audiobook chapters be parsed and tagged?
- Should music tagging modify files automatically?
- Should MusicBrainz be required for music?
- Should Discogs integration be included?
- Should Mema support multiple album releases and remasters?
- Should lyrics be managed?
- Should music videos be part of music or video domains?
