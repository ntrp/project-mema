---
title: Library Import
description: Matching existing files to media items.
---

Library import scans configured folders, shows discovered files, and lets the
user match each file to a media item.

## Matching

Matching compares detected title, year, and media kind against local library
items and metadata providers. Title matching is forgiving: punctuation and
special characters are normalized so names such as `Amelie` and `Amélie` or
`Walle`, `Wall-e`, and `Wall e` can match the same title.

## Import Behavior

When importing a list, rows are imported one at a time. Each row can show its
own spinner while the import is active, and imported rows are removed from the
visible list after completion.

## Duplicate Handling

The scan can group duplicate paths and provide removal actions for duplicate
library files when an import is attached to an existing media item.
