insert into app.media_profiles (
    id,
    name,
    is_default,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    final_container,
    remove_unwanted_audio,
    remove_unwanted_subtitles,
    preferred_protocol,
    series_pack_preference
)
values
    ('low-quality-demo', 'Low Quality Demo', false, true, 'webdl-480p', 0, 100, 1, 'mkv', false, false, 'any', 'preferEpisodes')
on conflict (id) do nothing;

insert into app.media_profile_video_targets (profile_id)
values
    ('low-quality-demo')
on conflict (profile_id) do nothing;

insert into app.media_profile_audio_targets (profile_id, language_id, score, sort_order)
values
    ('low-quality-demo', 'english', 0, 0)
on conflict (profile_id, language_id) do nothing;

insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
values
    ('low-quality-demo', 'sdtv', 1),
    ('low-quality-demo', 'dvd', 2),
    ('low-quality-demo', 'webdl-480p', 3),
    ('low-quality-demo', 'webrip-480p', 3)
on conflict (profile_id, quality_id) do nothing;

insert into app.media_profile_custom_formats (profile_id, custom_format_id, score)
values
    ('low-quality-demo', '8df7fe4569063d39319f07e69707285a', 100),
    ('low-quality-demo', '3df5e6dfef4b09bb6002f732bed5b774', 5),
    ('low-quality-demo', '90a6f9a284dff5103f6346090e6280c8', -10000),
    ('low-quality-demo', 'e2315f990da2e2cbfc9fa5b7a6fcfe48', -10000),
    ('low-quality-demo', 'e1a997ddb54e3ecbfe06341ad323c458', -10000)
on conflict (profile_id, custom_format_id) do nothing;

insert into app.users (id, username, password_hash, role)
values (
    'f20c5f51-eeb2-4895-9537-46e1129a9757',
    'guest',
    'pbkdf2-sha256$210000$ZGV2LWd1ZXN0LXNlZWQhIQ$X2RpVhDzffETuSK0e7nwDc6J7ZZH7jVCTqB3or4lVwk',
    'user'
)
on conflict (username) do nothing;

insert into app.tags (id, name, created_at, updated_at)
values
    ('17110be5-44e8-48b5-ae1f-fc85c713a791', 'HD', '2026-07-04 14:20:00+00', '2026-07-04 14:20:00+00'),
    ('ef580810-4d65-44c3-a8ff-cbdb56cf8f1b', 'Kids', '2026-07-04 14:20:00+00', '2026-07-04 14:20:00+00'),
    ('d06f380f-4c62-4af8-92d2-944074ffce97', 'Anime', '2026-07-04 14:20:00+00', '2026-07-04 14:20:00+00')
on conflict (lower(name)) do update
set name = excluded.name,
    updated_at = excluded.updated_at;

insert into app.library_folders (id, path, kind)
values ('10000000-0000-4000-8000-000000000100', '.data/media/test-movie', 'movie')
on conflict (path) do update
set kind = excluded.kind,
    updated_at = now();

insert into app.media_profiles (
    id,
    name,
    is_default,
    final_container,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_unwanted_audio,
    audio_lossy_transcode_policy,
    remove_unwanted_subtitles,
    subtitle_mode,
    allow_subtitle_release_fallback,
    preferred_protocol,
    series_pack_preference
)
values
    ('test-audio-aac-stereo', 'Test - Audio AAC stereo EN', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', false, 'mixed', true, 'any', 'auto'),
    ('test-audio-eac3-51', 'Test - Audio EAC3 5.1 EN', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', false, 'mixed', true, 'any', 'auto'),
    ('test-audio-italian-aac', 'Test - Audio AAC stereo IT', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', false, 'mixed', true, 'any', 'auto'),
    ('test-audio-japanese-aac', 'Test - Audio AAC stereo JA', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', false, 'mixed', true, 'any', 'auto'),
    ('test-subtitle-embedded-en-it', 'Test - Subtitles embedded EN IT', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', true, 'embedded', true, 'any', 'auto'),
    ('test-subtitle-external-en', 'Test - Subtitles external EN', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', true, 'external', true, 'any', 'auto'),
    ('test-subtitle-mixed-en-it', 'Test - Subtitles mixed EN IT', false, 'mkv', true, 'bluray-1080p', 0, 0, 1, true, 'lossyToLossy', true, 'mixed', true, 'any', 'auto')
on conflict (id) do update
set name = excluded.name,
    final_container = excluded.final_container,
    upgrades_allowed = excluded.upgrades_allowed,
    upgrade_until_quality_id = excluded.upgrade_until_quality_id,
    remove_unwanted_audio = excluded.remove_unwanted_audio,
    audio_lossy_transcode_policy = excluded.audio_lossy_transcode_policy,
    remove_unwanted_subtitles = excluded.remove_unwanted_subtitles,
    subtitle_mode = excluded.subtitle_mode,
    allow_subtitle_release_fallback = excluded.allow_subtitle_release_fallback,
    preferred_protocol = excluded.preferred_protocol,
    series_pack_preference = excluded.series_pack_preference,
    updated_at = now();

insert into app.media_profile_video_targets (
    profile_id,
    codecs,
    codec_score,
    hdr_formats,
    hdr_score,
    pixel_formats,
    pixel_format_score
)
select profile_id, array['h264'], 0, array['sdr'], 0, array['yuv420p'], 0
from (values
    ('test-audio-aac-stereo'),
    ('test-audio-eac3-51'),
    ('test-audio-italian-aac'),
    ('test-audio-japanese-aac'),
    ('test-subtitle-embedded-en-it'),
    ('test-subtitle-external-en'),
    ('test-subtitle-mixed-en-it')
) as profiles(profile_id)
on conflict (profile_id) do update
set codecs = excluded.codecs,
    codec_score = excluded.codec_score,
    hdr_formats = excluded.hdr_formats,
    hdr_score = excluded.hdr_score,
    pixel_formats = excluded.pixel_formats,
    pixel_format_score = excluded.pixel_format_score;

insert into app.media_profile_audio_targets (
    profile_id,
    language_id,
    score,
    target_codec,
    target_channels,
    minimum_bitrate_kbps,
    preferred_bitrate_kbps,
    sort_order
)
values
    ('test-audio-aac-stereo', 'english', 100, 'aac', array['2.0'], 192, 256, 0),
    ('test-audio-eac3-51', 'english', 100, 'eac3', array['5.1'], 640, 768, 0),
    ('test-audio-italian-aac', 'italian', 100, 'aac', array['2.0'], 160, 256, 0),
    ('test-audio-japanese-aac', 'japanese', 100, 'aac', array['2.0'], 160, 256, 0),
    ('test-subtitle-embedded-en-it', 'english', 100, 'aac', array['2.0'], 192, 256, 0),
    ('test-subtitle-external-en', 'english', 100, 'aac', array['2.0'], 192, 256, 0),
    ('test-subtitle-mixed-en-it', 'english', 100, 'aac', array['2.0'], 192, 256, 0)
on conflict (profile_id, language_id) do update
set score = excluded.score,
    target_codec = excluded.target_codec,
    target_channels = excluded.target_channels,
    minimum_bitrate_kbps = excluded.minimum_bitrate_kbps,
    preferred_bitrate_kbps = excluded.preferred_bitrate_kbps,
    sort_order = excluded.sort_order;

insert into app.media_profile_subtitle_targets (profile_id, language_id, score, formats, sort_order)
values
    ('test-audio-aac-stereo', 'english', 0, array['srt'], 0),
    ('test-audio-eac3-51', 'english', 0, array['srt'], 0),
    ('test-audio-italian-aac', 'italian', 0, array['srt'], 0),
    ('test-audio-japanese-aac', 'english', 0, array['srt'], 0),
    ('test-subtitle-embedded-en-it', 'english', 50, array['srt'], 0),
    ('test-subtitle-embedded-en-it', 'italian', 50, array['srt'], 1),
    ('test-subtitle-external-en', 'english', 50, array['srt'], 0),
    ('test-subtitle-mixed-en-it', 'english', 50, array['srt'], 0),
    ('test-subtitle-mixed-en-it', 'italian', 50, array['srt'], 1)
on conflict (profile_id, language_id) do update
set score = excluded.score,
    formats = excluded.formats,
    sort_order = excluded.sort_order;

insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
select profile_id, quality_id, sort_order
from (values
    ('test-audio-aac-stereo'),
    ('test-audio-eac3-51'),
    ('test-audio-italian-aac'),
    ('test-audio-japanese-aac'),
    ('test-subtitle-embedded-en-it'),
    ('test-subtitle-external-en'),
    ('test-subtitle-mixed-en-it')
) as profiles(profile_id)
cross join (values
    ('webdl-720p', 1),
    ('webrip-720p', 2),
    ('webdl-1080p', 3),
    ('webrip-1080p', 4),
    ('bluray-1080p', 5)
) as qualities(quality_id, sort_order)
on conflict (profile_id, quality_id) do nothing;
