---
title: Profiles
description: Desired media state for video, audio, subtitles, and file status.
---

Profiles describe what a finished media file should look like.

## Profile Areas

- **Video** defines the target quality and optional video properties.
- **Audio** defines target languages and requirements such as codec, channels,
  and minimum bitrate.
- **Subtitles** define wanted subtitle languages, formats, and placement mode.
- **Release scoring** helps rank search results before import or processing.

## File Status

File overview badges summarize whether the current file meets the selected
profile:

- **Ok** means every requirement is satisfied.
- **Partial** means the media is available but some requirements still need work.
- **Missing** means a required component is not present in a usable form.

Audio status is evaluated against audio requirements only. Subtitle status is
evaluated against subtitle targets and the configured subtitle mode.

## More Detail

For the full setup flow, see
[Qualities, Formats, And Profiles](/user-guide/using/qualities-formats-profiles/).
