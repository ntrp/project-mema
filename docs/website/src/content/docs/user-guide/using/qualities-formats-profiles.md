---
title: Qualities, Formats, And Profiles
description: Define what a finished file should look like.
---

Profiles are the heart of Media Manager. A profile describes the finished state
you want for a media item: video quality, upgrade behavior, audio languages,
subtitle languages, custom format scores, and final container.

Qualities and custom formats support profiles. Qualities describe the broad
release level, such as SD, 1080p, 2160p, web, Blu-ray, or remux-style releases.
Custom formats describe finer traits such as codec, release group, edition,
source, indexer flags, language, or text found in the release title.

## Quality Sizes

Quality sizes define the expected size range for each quality. The app uses
these values while judging releases. A release that is much smaller than
expected may be low quality or mislabeled. A release that is much larger than
expected may be wasteful for the selected profile.

The default quality sizes are estimated from resolution, source, and a
high-quality H.264 encode baseline. Lower-quality sources stay permissive, while
Blu-ray, remux, and UHD qualities get progressively larger preferred and maximum
ranges. Tune them only after you see real releases being accepted or rejected for
the wrong size reason.

## Custom Formats

Custom formats let you reward or reject specific release traits. A format can
look for conditions such as a release title pattern, source, resolution,
quality, video codec, audio codec, release group, edition, indexer flag, or
language.

Conditions say what should match. Negated conditions say what should not match.
The `Required` switch controls whether a condition must be present or simply
contributes when it appears.

Keep custom formats focused. One format should describe one idea, such as a
preferred release group, a dual-audio anime release, a codec preference, or a
retag you usually want to avoid.

## Profile Basics

Every profile has a name, a final container, selected qualities, and at least
one audio target. One profile can be marked as the default. The default profile
is used when the app needs a starting choice for new media.

The final container decides the intended output container, such as MKV or MP4.
This is a desired state setting. Whether the app can already transform every
file into that container depends on the available media tools and current media
workflow.

## Upgrades

When upgrades are allowed, the app can consider replacing an existing file with
a better release. The `Upgrade until` setting stops that process once a selected
quality is reached.

Custom format score settings add another upgrade layer. A minimum score can
reject weak releases. An upgrade-until score can define when the custom format
part is good enough. A minimum score increment prevents tiny score improvements
from replacing an existing file too aggressively.

## Video Targets

Video targets let you express preferences for codec, HDR format, and pixel
format. Leaving a field empty means any value is acceptable for that part of the
profile.

Use these controls when quality alone is not specific enough. For example, you
might allow 2160p but prefer a particular HDR format or codec.

## Audio Targets

Audio targets describe wanted language tracks. Each target has a language and a
score. It can also include a target codec, target channel layout, and minimum
bitrate.

When a release is missing a wanted audio language, release search marks it with
a warning instead of rejecting it outright. Automatic search still prefers
otherwise comparable releases that include the wanted language.

A profile can remove audio tracks that are not wanted. Treat that setting
carefully: it is useful for clean final files, but it should only be enabled
when your profile accurately describes every language you want to keep.

Lossy audio conversion controls how willing the app should be to convert audio
when a target codec is requested. Keep conversion disabled unless you know you
want the app to create different audio tracks.

## Subtitle Targets

Subtitle targets work like audio targets. Choose the wanted language, optional
formats, and score. The subtitle mode controls where subtitles should live.

When a release is missing a wanted subtitle language, release search marks it
with a warning instead of rejecting it outright. Automatic search gives those
releases lower priority than otherwise comparable subtitle-complete releases.

`Embedded` means wanted subtitles should be inside the media file. `External`
means they should be sidecar files. `Mixed` allows both, keeping existing
external subtitles external while allowing downloaded subtitles to be embedded.

The profile can remove subtitle tracks that are not wanted and can also allow
subtitle searches in other releases when the current release does not provide a
good match.

## A Practical First Profile

For a first profile, choose a small set of qualities, one required audio
language, a subtitle mode of Mixed, and no strict custom format minimum. Mark it
as the default.

After you have imported a few real files, add stricter audio, subtitle, video,
and custom format rules. Profiles are easier to tune from real examples than
from theory.
