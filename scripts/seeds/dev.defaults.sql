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

drop table if exists dev_test_media_profiles;
create temporary table dev_test_media_profiles (
    profile_id text primary key,
    audio_language_id text not null,
    audio_codec text not null,
    audio_channels text[] not null,
    audio_minimum_kbps integer not null,
    audio_preferred_kbps integer not null,
    audio_lossy_transcode_policy text not null,
    subtitle_mode text not null,
    remove_unwanted_audio boolean not null,
    remove_unwanted_subtitles boolean not null
) on commit drop;

insert into dev_test_media_profiles (
    profile_id,
    audio_language_id,
    audio_codec,
    audio_channels,
    audio_minimum_kbps,
    audio_preferred_kbps,
    audio_lossy_transcode_policy,
    subtitle_mode,
    remove_unwanted_audio,
    remove_unwanted_subtitles
)
values
    ('01-ok-embedded', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('02-missing-italian-audio', 'italian', 'aac', array['2.0'], 160, 256, 'lossyToLossy', 'mixed', true, false),
    ('03-wrong-audio-codec', 'english', 'eac3', array['5.1'], 640, 768, 'lossyToLossy', 'mixed', true, false),
    ('04-wrong-audio-channels', 'japanese', 'aac', array['2.0'], 160, 256, 'lossyToLossy', 'mixed', true, false),
    ('05-low-audio-bitrate', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('06-unwanted-audio', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('07-embedded-subtitle-needed', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'embedded', true, true),
    ('08-external-subtitle-mode', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'external', true, true),
    ('09-mixed-existing-external', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, true),
    ('10-unwanted-subtitle', 'korean', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'external', true, true),
    ('11-chapter-delete-summary', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('12-other-files-actions', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'external', true, true),
    ('13-three-movies-one-folder-grand-budapest', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('13-three-movies-one-folder-interstellar', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('13-three-movies-one-folder-truman-show', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('14-subrip-subtitle', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('15-ass-subtitle', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('16-wrong-video-resolution', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('17-audio-conversion-disabled', 'english', 'aac', array['stereo'], 192, 256, 'disabled', 'mixed', true, false),
    ('18-audio-conversion-lossless', 'english', 'aac', array['stereo'], 192, 256, 'losslessToLossy', 'mixed', true, false),
    ('19-audio-conversion-lossy', 'english', 'aac', array['stereo'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('20-wrong-video-codec', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false),
    ('21-wrong-container', 'english', 'aac', array['2.0'], 192, 256, 'lossyToLossy', 'mixed', true, false);

alter table dev_test_media_profiles
add column final_container text not null default 'mkv';

update dev_test_media_profiles
set final_container = 'mp4'
where profile_id = '21-wrong-container';

drop table if exists dev_test_media_subtitles;
create temporary table dev_test_media_subtitles (
    profile_id text not null,
    language_id text not null,
    score integer not null,
    sort_order integer not null
) on commit drop;

insert into dev_test_media_subtitles (profile_id, language_id, score, sort_order)
values
    ('01-ok-embedded', 'english', 0, 0),
    ('02-missing-italian-audio', 'italian', 0, 0),
    ('03-wrong-audio-codec', 'english', 0, 0),
    ('04-wrong-audio-channels', 'english', 0, 0),
    ('05-low-audio-bitrate', 'english', 0, 0),
    ('06-unwanted-audio', 'english', 0, 0),
    ('07-embedded-subtitle-needed', 'english', 50, 0),
    ('08-external-subtitle-mode', 'english', 50, 0),
    ('09-mixed-existing-external', 'english', 50, 0),
    ('09-mixed-existing-external', 'italian', 50, 1),
    ('10-unwanted-subtitle', 'english', 50, 0),
    ('11-chapter-delete-summary', 'english', 0, 0),
    ('12-other-files-actions', 'english', 50, 0),
    ('13-three-movies-one-folder-grand-budapest', 'english', 0, 0),
    ('13-three-movies-one-folder-interstellar', 'english', 0, 0),
    ('13-three-movies-one-folder-truman-show', 'english', 0, 0),
    ('14-subrip-subtitle', 'english', 0, 0),
    ('15-ass-subtitle', 'english', 0, 0),
    ('16-wrong-video-resolution', 'english', 0, 0),
    ('17-audio-conversion-disabled', 'english', 0, 0),
    ('18-audio-conversion-lossless', 'english', 0, 0),
    ('19-audio-conversion-lossy', 'english', 0, 0),
    ('20-wrong-video-codec', 'english', 0, 0),
    ('21-wrong-container', 'english', 0, 0);

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
select
    profile_id,
    profile_id,
    false,
    final_container,
    true,
    'bluray-1080p',
    0,
    0,
    1,
    remove_unwanted_audio,
    audio_lossy_transcode_policy,
    remove_unwanted_subtitles,
    subtitle_mode,
    true,
    'any',
    'auto'
from dev_test_media_profiles
on conflict (id) do update
set name = excluded.name,
    is_default = excluded.is_default,
    final_container = excluded.final_container,
    upgrades_allowed = excluded.upgrades_allowed,
    upgrade_until_quality_id = excluded.upgrade_until_quality_id,
    minimum_custom_format_score = excluded.minimum_custom_format_score,
    upgrade_until_custom_format_score = excluded.upgrade_until_custom_format_score,
    minimum_custom_format_score_increment = excluded.minimum_custom_format_score_increment,
    remove_unwanted_audio = excluded.remove_unwanted_audio,
    audio_lossy_transcode_policy = excluded.audio_lossy_transcode_policy,
    remove_unwanted_subtitles = excluded.remove_unwanted_subtitles,
    subtitle_mode = excluded.subtitle_mode,
    allow_subtitle_release_fallback = excluded.allow_subtitle_release_fallback,
    preferred_protocol = excluded.preferred_protocol,
    series_pack_preference = excluded.series_pack_preference,
    updated_at = now();

delete from app.media_profile_audio_targets
where profile_id in (select profile_id from dev_test_media_profiles);

delete from app.media_profile_subtitle_targets
where profile_id in (select profile_id from dev_test_media_profiles);

delete from app.media_profile_qualities
where profile_id in (select profile_id from dev_test_media_profiles);

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
from dev_test_media_profiles
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
select
    profile_id,
    audio_language_id,
    100,
    audio_codec,
    audio_channels,
    audio_minimum_kbps,
    audio_preferred_kbps,
    0
from dev_test_media_profiles;

insert into app.media_profile_subtitle_targets (profile_id, language_id, score, formats, sort_order)
select profile_id, language_id, score, array['subrip'], sort_order
from dev_test_media_subtitles;

insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
select profile_id, quality_id, sort_order
from dev_test_media_profiles
cross join (values
    ('webdl-720p', 1),
    ('webrip-720p', 2),
    ('webdl-1080p', 3),
    ('webrip-1080p', 4),
    ('bluray-1080p', 5)
) as qualities(quality_id, sort_order)
on conflict (profile_id, quality_id) do update
set sort_order = excluded.sort_order;

drop table if exists dev_test_media_imports;
create temporary table dev_test_media_imports (
    sort_order integer primary key,
    profile_id text not null,
    media_item_id uuid not null,
    scan_item_id uuid not null,
    year integer not null,
    external_id text not null,
    original_language text not null,
    folder_path text not null,
    file_name text not null,
    detected_title text not null
) on commit drop;

insert into dev_test_media_imports (
    sort_order,
    profile_id,
    media_item_id,
    scan_item_id,
    year,
    external_id,
    original_language,
    folder_path,
    file_name,
    detected_title
)
values
    (1, '01-ok-embedded', '10000000-0000-4000-8000-000000001001', '10000000-0000-4000-8000-000000002001', 2003, '12', 'EN', '.data/media/test-movie/01-ok-embedded', 'Finding.Nemo.2003.tmdb-12.1080p.WEB-DL.AAC2.0.EN.mkv', 'Finding Nemo'),
    (2, '02-missing-italian-audio', '10000000-0000-4000-8000-000000001002', '10000000-0000-4000-8000-000000002002', 2001, '194', 'FR', '.data/media/test-movie/02-missing-italian-audio', 'Amelie.2001.tmdb-194.1080p.WEB-DL.AAC2.0.EN.mkv', 'Amelie'),
    (3, '03-wrong-audio-codec', '10000000-0000-4000-8000-000000001003', '10000000-0000-4000-8000-000000002003', 1999, '603', 'EN', '.data/media/test-movie/03-wrong-audio-codec', 'The.Matrix.1999.tmdb-603.1080p.WEB-DL.AC3.5.1.EN.mkv', 'The Matrix'),
    (4, '04-wrong-audio-channels', '10000000-0000-4000-8000-000000001004', '10000000-0000-4000-8000-000000002004', 2001, '129', 'JA', '.data/media/test-movie/04-wrong-audio-channels', 'Spirited.Away.2001.tmdb-129.1080p.WEB-DL.AAC1.0.JA.mkv', 'Spirited Away'),
    (5, '05-low-audio-bitrate', '10000000-0000-4000-8000-000000001005', '10000000-0000-4000-8000-000000002005', 2015, '76341', 'EN', '.data/media/test-movie/05-low-audio-bitrate', 'Mad.Max.Fury.Road.2015.tmdb-76341.1080p.WEB-DL.AAC2.0.EN.mkv', 'Mad Max Fury Road'),
    (6, '06-unwanted-audio', '10000000-0000-4000-8000-000000001006', '10000000-0000-4000-8000-000000002006', 2007, '2062', 'EN', '.data/media/test-movie/06-unwanted-audio', 'Ratatouille.2007.tmdb-2062.1080p.WEB-DL.AAC2.0.EN-ES.mkv', 'Ratatouille'),
    (7, '07-embedded-subtitle-needed', '10000000-0000-4000-8000-000000001007', '10000000-0000-4000-8000-000000002007', 2014, '116149', 'EN', '.data/media/test-movie/07-embedded-subtitle-needed', 'Paddington.2014.tmdb-116149.1080p.WEB-DL.AAC2.0.EN.mkv', 'Paddington'),
    (8, '08-external-subtitle-mode', '10000000-0000-4000-8000-000000001008', '10000000-0000-4000-8000-000000002008', 2016, '329865', 'EN', '.data/media/test-movie/08-external-subtitle-mode', 'Arrival.2016.tmdb-329865.1080p.WEB-DL.AAC2.0.EN.mkv', 'Arrival'),
    (9, '09-mixed-existing-external', '10000000-0000-4000-8000-000000001009', '10000000-0000-4000-8000-000000002009', 2015, '150540', 'EN', '.data/media/test-movie/09-mixed-existing-external', 'Inside.Out.2015.tmdb-150540.1080p.WEB-DL.AAC2.0.EN.mkv', 'Inside Out'),
    (10, '10-unwanted-subtitle', '10000000-0000-4000-8000-000000001010', '10000000-0000-4000-8000-000000002010', 2019, '496243', 'KO', '.data/media/test-movie/10-unwanted-subtitle', 'Parasite.2019.tmdb-496243.1080p.WEB-DL.AAC2.0.KO.mkv', 'Parasite'),
    (11, '11-chapter-delete-summary', '10000000-0000-4000-8000-000000001011', '10000000-0000-4000-8000-000000002011', 2010, '27205', 'EN', '.data/media/test-movie/11-chapter-delete-summary', 'Inception.2010.tmdb-27205.1080p.WEB-DL.AAC2.0.EN.Chapters.mkv', 'Inception'),
    (12, '12-other-files-actions', '10000000-0000-4000-8000-000000001012', '10000000-0000-4000-8000-000000002012', 2008, '10681', 'EN', '.data/media/test-movie/12-other-files-actions', 'WALL-E.2008.tmdb-10681.1080p.WEB-DL.AAC2.0.EN.mkv', 'WALL E'),
    (13, '13-three-movies-one-folder-grand-budapest', '10000000-0000-4000-8000-000000001013', '10000000-0000-4000-8000-000000002013', 2014, '120467', 'EN', '.data/media/test-movie/13-three-movies-one-folder', 'The.Grand.Budapest.Hotel.2014.tmdb-120467.1080p.WEB-DL.AAC2.0.EN.mkv', 'The Grand Budapest Hotel'),
    (14, '13-three-movies-one-folder-interstellar', '10000000-0000-4000-8000-000000001014', '10000000-0000-4000-8000-000000002014', 2014, '157336', 'EN', '.data/media/test-movie/13-three-movies-one-folder', 'Interstellar.2014.tmdb-157336.1080p.WEB-DL.AAC2.0.EN.mkv', 'Interstellar'),
    (15, '13-three-movies-one-folder-truman-show', '10000000-0000-4000-8000-000000001015', '10000000-0000-4000-8000-000000002015', 1998, '37165', 'EN', '.data/media/test-movie/13-three-movies-one-folder', 'The.Truman.Show.1998.tmdb-37165.1080p.WEB-DL.AAC2.0.EN.mkv', 'The Truman Show'),
    (16, '14-subrip-subtitle', '10000000-0000-4000-8000-000000001016', '10000000-0000-4000-8000-000000002016', 2001, '585', 'EN', '.data/media/test-movie/14-subrip-subtitle', 'Monsters.Inc.2001.tmdb-585.1080p.WEB-DL.AAC2.0.EN.SubRip.mkv', 'Monsters Inc'),
    (17, '15-ass-subtitle', '10000000-0000-4000-8000-000000001017', '10000000-0000-4000-8000-000000002017', 2009, '14160', 'EN', '.data/media/test-movie/15-ass-subtitle', 'Up.2009.tmdb-14160.1080p.WEB-DL.AAC2.0.EN.ASS.mkv', 'Up'),
    (18, '16-wrong-video-resolution', '10000000-0000-4000-8000-000000001018', '10000000-0000-4000-8000-000000002018', 2006, '920', 'EN', '.data/media/test-movie/16-wrong-video-resolution', 'Cars.2006.tmdb-920.1080p.WEB-DL.AAC2.0.EN.WrongResolution.mkv', 'Cars'),
    (19, '17-audio-conversion-disabled', '10000000-0000-4000-8000-000000001019', '10000000-0000-4000-8000-000000002019', 1995, '862', 'EN', '.data/media/test-movie/17-audio-conversion-disabled', 'Toy.Story.1995.tmdb-862.1080p.WEB-DL.FLAC2.0.EN.AudioConversionDisabled.mkv', 'Toy Story'),
    (20, '18-audio-conversion-lossless', '10000000-0000-4000-8000-000000001020', '10000000-0000-4000-8000-000000002020', 1998, '9487', 'EN', '.data/media/test-movie/18-audio-conversion-lossless', 'A.Bugs.Life.1998.tmdb-9487.1080p.WEB-DL.FLAC2.0.EN.AudioConversionLossless.mkv', 'A Bugs Life'),
    (21, '19-audio-conversion-lossy', '10000000-0000-4000-8000-000000001021', '10000000-0000-4000-8000-000000002021', 2012, '62177', 'EN', '.data/media/test-movie/19-audio-conversion-lossy', 'Brave.2012.tmdb-62177.1080p.WEB-DL.AC3.2.0.EN.AudioConversionLossy.mkv', 'Brave'),
    (22, '20-wrong-video-codec', '10000000-0000-4000-8000-000000001022', '10000000-0000-4000-8000-000000002022', 2017, '354912', 'EN', '.data/media/test-movie/20-wrong-video-codec', 'Coco.2017.tmdb-354912.1080p.WEB-DL.MPEG4.AAC2.0.EN.WrongVideoCodec.mkv', 'Coco'),
    (23, '21-wrong-container', '10000000-0000-4000-8000-000000001023', '10000000-0000-4000-8000-000000002023', 2021, '508943', 'EN', '.data/media/test-movie/21-wrong-container', 'Luca.2021.tmdb-508943.1080p.WEB-DL.AAC2.0.EN.WrongContainer.mkv', 'Luca');

insert into app.media_items (
    id,
    media_type,
    content_kind,
    title,
    year,
    monitored,
    external_provider,
    external_id,
    metadata_status,
    original_language,
    quality_profile_id,
    library_folder_id,
    media_folder_path,
    monitor_mode,
    minimum_availability
)
select
    media_item_id,
    'movie',
    'standard',
    profile_id,
    year,
    false,
    'tmdb',
    external_id,
    'seeded',
    original_language,
    profile_id,
    '10000000-0000-4000-8000-000000000100',
    folder_path,
    'none',
    'released'
from dev_test_media_imports
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
    quality_profile_id = excluded.quality_profile_id,
    library_folder_id = excluded.library_folder_id,
    media_folder_path = excluded.media_folder_path,
    monitor_mode = excluded.monitor_mode,
    minimum_availability = excluded.minimum_availability,
    updated_at = now();

insert into app.library_scans (
    id,
    library_folder_id,
    status,
    total_files,
    auto_matched_count,
    manual_count,
    completed_at
)
select
    '10000000-0000-4000-8000-000000003001',
    '10000000-0000-4000-8000-000000000100',
    'completed',
    count(*)::integer,
    count(*)::integer,
    0,
    now()
from dev_test_media_imports
on conflict (id) do update
set library_folder_id = excluded.library_folder_id,
    status = excluded.status,
    total_files = excluded.total_files,
    auto_matched_count = excluded.auto_matched_count,
    manual_count = excluded.manual_count,
    completed_at = excluded.completed_at;

delete from app.library_scan_items
where scan_id = '10000000-0000-4000-8000-000000003001';

insert into app.library_scan_items (
    id,
    scan_id,
    path,
    file_name,
    size_bytes,
    detected_title,
    detected_year,
    detected_media_kind,
    status,
    imported,
    matched_title,
    matched_year,
    matched_media_kind,
    matched_external_provider,
    matched_external_id,
    match_source,
    selected_metadata_provider_id,
    duplicate_removal_allowed,
    media_item_id
)
select
    scan_item_id,
    '10000000-0000-4000-8000-000000003001',
    folder_path || '/' || file_name,
    file_name,
    0,
    detected_title,
    year,
    'movie',
    'auto_added',
    true,
    profile_id,
    year,
    'movie',
    'tmdb',
    external_id,
    'manual',
    '00000000-0000-4000-8000-000000000101',
    false,
    media_item_id
from dev_test_media_imports
order by sort_order;
