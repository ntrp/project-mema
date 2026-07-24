insert into app.library_folders (id, path, kind)
values ('20000000-0000-4000-8000-000000000100', '.data/media/test-series', 'series')
on conflict (path) do update
set kind = excluded.kind,
    updated_at = now();

insert into app.media_items (
    id, media_type, content_kind, title, year, monitored, external_provider,
    external_id, metadata_status, original_language, series_type,
    numbering_strategy, season_count, episode_count, quality_profile_id,
    library_folder_id, media_folder_path, monitor_mode, minimum_availability
)
values
    (
        '20000000-0000-4000-8000-000000001001', 'serie', 'standard',
        'Severance', 2022, false, 'tmdb', '95396', 'seeded', 'EN', 'standard',
        'tmdb_season_episode', 1, 2, '01-ok-embedded',
        '20000000-0000-4000-8000-000000000100',
        '.data/media/test-series/01-standard-season', 'none', 'released'
    ),
    (
        '20000000-0000-4000-8000-000000001002', 'serie', 'anime',
        'Cowboy Bebop', 1998, false, 'anidb', '23', 'seeded', 'JA', 'absolute',
        'anidb_absolute', 1, 2, '09-mixed-existing-external',
        '20000000-0000-4000-8000-000000000100',
        '.data/media/test-series/02-anime-absolute', 'none', 'released'
    )
on conflict (id) do update
set media_type = excluded.media_type,
    content_kind = excluded.content_kind,
    title = excluded.title,
    year = excluded.year,
    monitored = excluded.monitored,
    external_provider = excluded.external_provider,
    external_id = excluded.external_id,
    metadata_status = excluded.metadata_status,
    original_language = excluded.original_language,
    series_type = excluded.series_type,
    numbering_strategy = excluded.numbering_strategy,
    season_count = excluded.season_count,
    episode_count = excluded.episode_count,
    quality_profile_id = excluded.quality_profile_id,
    library_folder_id = excluded.library_folder_id,
    media_folder_path = excluded.media_folder_path,
    monitor_mode = excluded.monitor_mode,
    minimum_availability = excluded.minimum_availability,
    updated_at = now();

insert into app.media_seasons (
    id, media_item_id, external_provider, season_number, name,
    episode_count, monitored, source
)
values
    (
        '20000000-0000-4000-8000-000000002001',
        '20000000-0000-4000-8000-000000001001', 'tmdb', 1, 'Season 1', 2,
        false, '{"kind":"development_seed"}'::jsonb
    ),
    (
        '20000000-0000-4000-8000-000000002002',
        '20000000-0000-4000-8000-000000001002', 'anidb', 1, 'Season 1', 2,
        false, '{"kind":"development_seed"}'::jsonb
    )
on conflict (id) do update
set media_item_id = excluded.media_item_id,
    external_provider = excluded.external_provider,
    season_number = excluded.season_number,
    name = excluded.name,
    episode_count = excluded.episode_count,
    monitored = excluded.monitored,
    source = excluded.source,
    updated_at = now();

insert into app.media_episodes (
    id, season_id, media_item_id, external_provider, season_number,
    episode_number, name, runtime_minutes, monitored, source
)
values
    ('20000000-0000-4000-8000-000000003001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000001001', 'tmdb', 1, 1, 'Good News About Hell', 57, false, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000003002', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000001001', 'tmdb', 1, 2, 'Half Loop', 54, false, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000003003', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000001002', 'anidb', 1, 1, 'Asteroid Blues', 24, false, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000003004', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000001002', 'anidb', 1, 2, 'Stray Dog Strut', 24, false, '{"kind":"development_seed"}'::jsonb)
on conflict (id) do update
set season_id = excluded.season_id,
    media_item_id = excluded.media_item_id,
    external_provider = excluded.external_provider,
    season_number = excluded.season_number,
    episode_number = excluded.episode_number,
    name = excluded.name,
    runtime_minutes = excluded.runtime_minutes,
    monitored = excluded.monitored,
    source = excluded.source,
    updated_at = now();

insert into app.media_provider_mappings (
    id, media_item_id, entity_type, provider_name, provider_entity_type,
    external_id, canonical, confidence, source
)
values
    ('20000000-0000-4000-8000-000000004001', '20000000-0000-4000-8000-000000001001', 'media_item', 'tmdb', 'tv', '95396', true, 1, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000004002', '20000000-0000-4000-8000-000000001002', 'media_item', 'anidb', 'anime', '23', true, 1, '{"kind":"development_seed"}'::jsonb)
on conflict (id) do update
set media_item_id = excluded.media_item_id,
    provider_name = excluded.provider_name,
    provider_entity_type = excluded.provider_entity_type,
    external_id = excluded.external_id,
    canonical = excluded.canonical,
    confidence = excluded.confidence,
    source = excluded.source,
    updated_at = now();

insert into app.media_episode_numbering (
    id, media_item_id, season_id, episode_id, provider_name,
    numbering_scheme, season_number, episode_number, absolute_number, source
)
values
    ('20000000-0000-4000-8000-000000005001', '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003001', 'tmdb', 'season_episode', 1, 1, null, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000005002', '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003002', 'tmdb', 'season_episode', 1, 2, null, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000005003', '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003003', 'anidb', 'absolute', null, null, 1, '{"kind":"development_seed"}'::jsonb),
    ('20000000-0000-4000-8000-000000005004', '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003004', 'anidb', 'absolute', null, null, 2, '{"kind":"development_seed"}'::jsonb)
on conflict (episode_id, provider_name, numbering_scheme) do update
set season_number = excluded.season_number,
    episode_number = excluded.episode_number,
    absolute_number = excluded.absolute_number,
    source = excluded.source,
    updated_at = now();

insert into app.library_scans (
    id, library_folder_id, status, total_files, auto_matched_count,
    manual_count, completed_at
)
values (
    '20000000-0000-4000-8000-000000006001',
    '20000000-0000-4000-8000-000000000100', 'completed', 4, 4, 0, now()
)
on conflict (id) do update
set status = excluded.status,
    total_files = excluded.total_files,
    auto_matched_count = excluded.auto_matched_count,
    manual_count = excluded.manual_count,
    completed_at = excluded.completed_at;

delete from app.library_scan_items
where scan_id = '20000000-0000-4000-8000-000000006001';

insert into app.library_scan_items (
    id, scan_id, path, file_name, size_bytes, detected_title, detected_year,
    detected_media_kind, season_number, episode_number, status, imported,
    matched_title, matched_year, matched_media_kind, matched_external_provider,
    matched_external_id, match_source, duplicate_removal_allowed,
    media_item_id, season_id, episode_id
)
values
    ('20000000-0000-4000-8000-000000007001', '20000000-0000-4000-8000-000000006001', '.data/media/test-series/01-standard-season/Season 01/Severance.2022.S01E01.480p.WEB-DL.AAC2.0.EN.mkv', 'Severance.2022.S01E01.480p.WEB-DL.AAC2.0.EN.mkv', 0, 'Severance', 2022, 'series', 1, 1, 'auto_added', true, 'Severance', 2022, 'series', 'tmdb', '95396', 'manual', false, '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003001'),
    ('20000000-0000-4000-8000-000000007002', '20000000-0000-4000-8000-000000006001', '.data/media/test-series/01-standard-season/Season 01/Severance.2022.S01E02.480p.WEB-DL.AAC2.0.EN.mkv', 'Severance.2022.S01E02.480p.WEB-DL.AAC2.0.EN.mkv', 0, 'Severance', 2022, 'series', 1, 2, 'auto_added', true, 'Severance', 2022, 'series', 'tmdb', '95396', 'manual', false, '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003002'),
    ('20000000-0000-4000-8000-000000007003', '20000000-0000-4000-8000-000000006001', '.data/media/test-series/02-anime-absolute/Season 01/Cowboy.Bebop.1998.S01E01.480p.WEB-DL.AAC2.0.JA.mkv', 'Cowboy.Bebop.1998.S01E01.480p.WEB-DL.AAC2.0.JA.mkv', 0, 'Cowboy Bebop', 1998, 'anime_series', 1, 1, 'auto_added', true, 'Cowboy Bebop', 1998, 'anime_series', 'anidb', '23', 'manual', false, '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003003'),
    ('20000000-0000-4000-8000-000000007004', '20000000-0000-4000-8000-000000006001', '.data/media/test-series/02-anime-absolute/Season 01/Cowboy.Bebop.1998.S01E02.480p.WEB-DL.AAC2.0.JA.mkv', 'Cowboy.Bebop.1998.S01E02.480p.WEB-DL.AAC2.0.JA.mkv', 0, 'Cowboy Bebop', 1998, 'anime_series', 1, 2, 'auto_added', true, 'Cowboy Bebop', 1998, 'anime_series', 'anidb', '23', 'manual', false, '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003004');

insert into app.media_file_facts (
    id, media_item_id, season_id, episode_id, file_path, quality_id,
    container_format, duration_ms, size_bytes, source_kind
)
values
    ('20000000-0000-4000-8000-000000008001', '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003001', '.data/media/test-series/01-standard-season/Season 01/Severance.2022.S01E01.480p.WEB-DL.AAC2.0.EN.mkv', 'webdl-480p', 'mkv', 3000, 0, 'import'),
    ('20000000-0000-4000-8000-000000008002', '20000000-0000-4000-8000-000000001001', '20000000-0000-4000-8000-000000002001', '20000000-0000-4000-8000-000000003002', '.data/media/test-series/01-standard-season/Season 01/Severance.2022.S01E02.480p.WEB-DL.AAC2.0.EN.mkv', 'webdl-480p', 'mkv', 3000, 0, 'import'),
    ('20000000-0000-4000-8000-000000008003', '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003003', '.data/media/test-series/02-anime-absolute/Season 01/Cowboy.Bebop.1998.S01E01.480p.WEB-DL.AAC2.0.JA.mkv', 'webdl-480p', 'mkv', 3000, 0, 'import'),
    ('20000000-0000-4000-8000-000000008004', '20000000-0000-4000-8000-000000001002', '20000000-0000-4000-8000-000000002002', '20000000-0000-4000-8000-000000003004', '.data/media/test-series/02-anime-absolute/Season 01/Cowboy.Bebop.1998.S01E02.480p.WEB-DL.AAC2.0.JA.mkv', 'webdl-480p', 'mkv', 3000, 0, 'import')
on conflict (media_item_id, file_path) do update
set season_id = excluded.season_id,
    episode_id = excluded.episode_id,
    quality_id = excluded.quality_id,
    container_format = excluded.container_format,
    duration_ms = excluded.duration_ms,
    size_bytes = excluded.size_bytes,
    source_kind = excluded.source_kind,
    updated_at = now();
