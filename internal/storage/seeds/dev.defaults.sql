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

insert into app.indexers (
    id,
    definition_id,
    name,
    type,
    implementation,
    implementation_name,
    protocol,
    privacy,
    language,
    indexer_urls,
    legacy_urls,
    base_url,
    categories,
    media_type_scopes,
    tag_scopes,
    fields,
    capabilities,
    redirect,
    app_profile_id,
    prefer_magnet_url,
    supports_rss,
    supports_search,
    supports_redirect,
    supports_pagination,
    enabled,
    priority,
    created_at,
    updated_at
)
values
    (
        '34ef8f5b-347b-41c5-a4f8-8892365a4dd3',
        'generic-torznab',
        'Local Movie Torznab',
        'torznab',
        'Torznab',
        'Local Movie Torznab',
        'torrent',
        'private',
        'EN',
        '{http://localhost:18082/api}'::text[],
        '{}'::text[],
        'http://localhost:18082/api',
        '{2000,2040}'::integer[],
        '{movie}'::text[],
        '{HD}'::text[],
        '[]'::jsonb,
        '{"categories":[{"id":2000,"name":"Movies","children":[{"id":2040,"name":"HD","children":[]}]}],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}'::jsonb,
        true,
        'default',
        false,
        true,
        true,
        true,
        true,
        true,
        90,
        '2026-07-04 14:20:00+00',
        '2026-07-04 14:20:00+00'
    ),
    (
        '425c5553-5721-4e18-b701-90c84c30b410',
        'generic-torznab',
        'Local Series Torznab',
        'torznab',
        'Torznab',
        'Local Series Torznab',
        'torrent',
        'private',
        'EN',
        '{http://localhost:18083/api}'::text[],
        '{}'::text[],
        'http://localhost:18083/api',
        '{5000,5040,5070}'::integer[],
        '{serie,anime}'::text[],
        '{Kids,Anime}'::text[],
        '[]'::jsonb,
        '{"categories":[{"id":5000,"name":"TV","children":[{"id":5040,"name":"HD","children":[]},{"id":5070,"name":"Anime","children":[]}]}],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}'::jsonb,
        true,
        'default',
        false,
        true,
        true,
        true,
        true,
        true,
        100,
        '2026-07-04 14:20:00+00',
        '2026-07-04 14:20:00+00'
    )
on conflict (id) do update
set media_type_scopes = excluded.media_type_scopes,
    tag_scopes = excluded.tag_scopes,
    categories = excluded.categories,
    updated_at = excluded.updated_at;
