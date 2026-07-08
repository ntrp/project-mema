---
title: Libraries And File Handling
description: Configure root folders, path mappings, imports, naming, and delete behavior.
---

The Library settings tell Media Manager where finished media should live and how
files should be matched, named, imported, and removed.

Start with root folders. Add one folder for Movies and one for Series if you
manage both. A root folder is the place where imported files end up, not
necessarily the place where your download client first saves them.

## Root Folders

When you add a root folder, choose the folder path and the folder type. The type
matters because a movie and a series episode are organized differently.

After a folder is added, the app can scan it. A scan finds media files and
offers import rows. Each row can be matched to an existing media item or to a
new metadata result. During import, you choose or confirm the metadata provider
and profile that should be used.

Imported rows disappear from the active list so you can focus on the files that
still need attention. When duplicates are detected, the import view shows them
as conflicts you can resolve instead of silently guessing which file to keep.

## Path Mappings

Path mappings are needed when a download client and Media Manager use different
paths for the same files.

For example, a download client might report `/downloads/movie.mkv`, while Media
Manager sees that same file as `/data/downloads/movie.mkv`. Add a mapping from
the client path to the app path so completed downloads can be imported.

If downloads finish but imports cannot find the file, path mappings are the
first thing to check.

## File Naming

File naming settings control the folder and file names used when media is
organized. Movies and series have separate templates. Series also has different
templates for normal episodes, daily episodes, anime episodes, season folders,
and specials.

Use templates that produce names your media players can scan reliably. A good
template usually includes title, year for movies, season and episode numbers for
series, and enough quality information to distinguish releases.

Custom formats can be included in rename templates when the format is marked for
that use. This is useful when traits such as `REMUX`, `HDR`, or a preferred
release group should remain visible in file names.

## Delete Policy

The file delete policy controls what happens when you delete files from the app.

`Delete permanently` removes files from disk. `Move to recycle folder` moves
files under a hidden folder inside the library root. `Keep files` records that
the delete was skipped and leaves the file in place.

For cautious setups, start with a recycle folder. The recycle folder must be a
hidden relative folder, such as `.recycle`, so deleted files stay inside the
library root instead of moving to an unexpected location.

## Existing Libraries

For an existing library, add the root folders first and scan them. Match a small
batch, import it, and confirm that the media detail page shows the expected
file, quality, audio, subtitle, and status information.

Once the first batch looks right, continue with larger imports.
