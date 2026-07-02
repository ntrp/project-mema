insert into app.media_profiles (
    id,
    name,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_non_enabled_languages,
    preferred_protocol,
    series_pack_preference
)
values
    ('low-quality-demo', 'Low Quality Demo', true, 'webdl-480p', 0, 100, 1, false, 'any', 'preferEpisodes')
on conflict (id) do nothing;

insert into app.media_profile_languages (profile_id, language_id, score)
values
    ('low-quality-demo', 'english', 0)
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
