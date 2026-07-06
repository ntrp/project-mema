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

insert into app.media_profile_audio_targets (profile_id, language_id, score, required, sort_order)
values
    ('low-quality-demo', 'english', 0, true, 0)
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
