-- +goose Up
create schema if not exists app;

create table if not exists app.users (
    id uuid primary key,
    username text not null unique,
    password_hash text not null,
    display_name text not null default '',
    picture_url text not null default '',
    role text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose StatementBegin
do $$
begin
    alter table app.users drop constraint if exists users_role_check;
    alter table app.users
        add constraint users_role_check check (role in ('admin', 'user'));
end $$;
-- +goose StatementEnd

create unique index if not exists idx_users_username_lower
    on app.users (lower(username));

create table if not exists app.sessions (
    id text primary key,
    user_id uuid not null references app.users(id) on delete cascade,
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

alter table app.sessions add column if not exists updated_at timestamptz not null default now();

create table if not exists app.download_clients (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('transmission', 'sabnzbd')),
    protocol text not null default 'torrent' check (protocol in ('torrent', 'usenet')),
    base_url text not null,
    username text,
    password text,
    api_key text,
    category text,
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose StatementBegin
do $$
begin
    alter table app.download_clients drop constraint if exists download_clients_type_protocol_check;
    alter table app.download_clients
        add constraint download_clients_type_protocol_check
        check (
            (type = 'transmission' and protocol = 'torrent')
            or (type = 'sabnzbd' and protocol = 'usenet')
        );
end $$;
-- +goose StatementEnd

create index if not exists idx_download_clients_priority
    on app.download_clients (priority, name);

create table if not exists app.indexers (
    id uuid primary key,
    definition_id text not null default 'generic-torznab',
    name text not null,
    type text not null default 'torznab' check (type in ('torznab', 'newznab', 'rss')),
    implementation text not null default 'Cardigann',
    implementation_name text not null default '',
    protocol text not null default 'torrent' check (protocol in ('torrent', 'usenet')),
    privacy text not null default 'private' check (privacy in ('public', 'private', 'semiPrivate')),
    language text not null default 'en-US',
    encoding text,
    description text,
    indexer_urls text[] not null default '{}',
    legacy_urls text[] not null default '{}',
    base_url text not null,
    api_key text,
    categories integer[] not null default '{}',
    media_type_scopes text[] not null default '{movie,serie,anime,audio,book}',
    tag_scopes text[] not null default '{}',
    fields jsonb not null default '[]'::jsonb,
    capabilities jsonb not null default '{"categories":[],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}'::jsonb,
    redirect boolean not null default true,
    app_profile_id text not null default 'default',
    minimum_seeders integer check (minimum_seeders is null or minimum_seeders >= 0),
    seed_ratio numeric(10, 2) check (seed_ratio is null or seed_ratio >= 0),
    seed_time integer check (seed_time is null or seed_time >= 0),
    pack_seed_time integer check (pack_seed_time is null or pack_seed_time >= 0),
    prefer_magnet_url boolean not null default false,
    supports_rss boolean not null default true,
    supports_search boolean not null default true,
    supports_redirect boolean not null default true,
    supports_pagination boolean not null default true,
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    health_status text not null default 'healthy' check (health_status in ('healthy', 'temporary_disabled', 'disabled')),
    last_query_at timestamptz,
    last_success_at timestamptz,
    last_failure_at timestamptz,
    next_check_at timestamptz,
    last_status_code integer,
    last_error text,
    failure_count integer not null default 0 check (failure_count >= 0),
    rss_marker_published_at timestamptz,
    rss_marker_guid text,
    rss_marker_download_url text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint indexers_media_type_scopes_check check (media_type_scopes <@ array['movie', 'serie', 'anime', 'audio', 'book']::text[])
);

create index if not exists idx_indexers_priority
    on app.indexers (priority, name);

create index if not exists idx_indexers_health
    on app.indexers (enabled, health_status, next_check_at);

create index if not exists idx_indexers_rss_enabled
    on app.indexers (enabled, supports_rss, health_status, next_check_at);

create table if not exists app.indexer_proxies (
    id uuid primary key,
    name text not null,
    implementation text not null,
    link text not null,
    enabled boolean not null default true,
    on_health_issue boolean not null default false,
    supports_on_health_issue boolean not null default true,
    include_health_warnings boolean not null default false,
    test_command text not null default 'test',
    fields jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_indexer_proxies_name
    on app.indexer_proxies (name);

create table if not exists app.indexer_search_settings (
    id boolean primary key default true check (id),
    cache_duration_minutes integer not null default 1440 check (cache_duration_minutes between 0 and 43200),
    history_retention_days integer not null default 7 check (history_retention_days between 1 and 365),
    automatic_blocklist_expiry_days integer not null default 7 check (automatic_blocklist_expiry_days between 1 and 365),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

insert into app.indexer_search_settings (
    id, cache_duration_minutes, history_retention_days, automatic_blocklist_expiry_days
)
values (true, 1440, 7, 7)
on conflict (id) do nothing;

create table if not exists app.indexer_search_cache (
    indexer_id uuid not null references app.indexers(id) on delete cascade,
    media_type text not null check (media_type in ('movie', 'serie', 'mixed')),
    query text not null,
    response jsonb not null,
    result_count integer not null default 0 check (result_count >= 0),
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    primary key (indexer_id, media_type, query)
);

create index if not exists idx_indexer_search_cache_expires
    on app.indexer_search_cache (expires_at);

create table if not exists app.indexer_search_history (
    id uuid primary key,
    indexer_id uuid references app.indexers(id) on delete set null,
    indexer_name text not null,
    indexer_type text not null default 'torznab',
    indexer_protocol text not null default 'torrent',
    media_type text not null check (media_type in ('movie', 'serie', 'mixed')),
    query text not null,
    cache_hit boolean not null,
    success boolean not null,
    result_count integer not null default 0 check (result_count >= 0),
    error text,
    response jsonb not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_indexer_search_history_created_at
    on app.indexer_search_history (created_at desc);

create table if not exists app.quality_size_settings (
    quality_id text primary key,
    minimum_size_mb_per_minute numeric(10, 2) not null default 0 check (minimum_size_mb_per_minute >= 0),
    preferred_size_mb_per_minute numeric(10, 2) check (preferred_size_mb_per_minute >= 0),
    maximum_size_mb_per_minute numeric(10, 2) check (maximum_size_mb_per_minute >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint quality_size_settings_order_check check (
        (preferred_size_mb_per_minute is null or preferred_size_mb_per_minute >= minimum_size_mb_per_minute)
        and (maximum_size_mb_per_minute is null or maximum_size_mb_per_minute >= minimum_size_mb_per_minute)
        and (
            preferred_size_mb_per_minute is null
            or maximum_size_mb_per_minute is null
            or preferred_size_mb_per_minute <= maximum_size_mb_per_minute
        )
    )
);

create table if not exists app.discover_blacklist (
    id uuid primary key,
    media_type text not null check (media_type in ('movie', 'serie')),
    title text not null,
    year integer,
    external_provider text,
    external_id text,
    overview text,
    poster_path text,
    created_at timestamptz not null default now()
);

create unique index if not exists idx_discover_blacklist_external
    on app.discover_blacklist (media_type, external_provider, external_id)
    where external_provider is not null and external_id is not null;

create unique index if not exists idx_discover_blacklist_title
    on app.discover_blacklist (media_type, lower(title), coalesce(year, 0))
    where external_provider is null or external_id is null;

create table if not exists app.media_profiles (
    id text primary key,
    name text not null,
    is_default boolean not null default false,
    final_container text not null default 'mkv' check (final_container in ('mkv', 'mp4')),
    upgrades_allowed boolean not null default true,
    upgrade_until_quality_id text,
    minimum_custom_format_score integer not null default 0,
    upgrade_until_custom_format_score integer not null default 0,
    minimum_custom_format_score_increment integer not null default 1 check (minimum_custom_format_score_increment >= 0),
    remove_unwanted_audio boolean not null default false,
    audio_lossy_transcode_policy text not null default 'disabled' check (audio_lossy_transcode_policy in ('disabled', 'losslessToLossy', 'lossyToLossy')),
    remove_unwanted_subtitles boolean not null default false,
    subtitle_preferred_mode text not null default 'mixed' check (subtitle_preferred_mode in ('mixed', 'embedded', 'external')),
    allow_subtitle_release_fallback boolean not null default false,
    preferred_protocol text not null default 'any' check (preferred_protocol in ('any', 'torrent', 'usenet')),
    series_pack_preference text not null default 'auto' check (series_pack_preference in ('auto', 'preferPacks', 'preferEpisodes')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_media_profiles_name_lower
    on app.media_profiles (lower(name));

create unique index if not exists idx_media_profiles_default
    on app.media_profiles (is_default)
    where is_default;

create table if not exists app.custom_formats (
    id uuid primary key,
    name text not null,
    include_in_rename_template boolean not null default false,
    include_specs jsonb not null default '[]'::jsonb,
    exclude_specs jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint custom_formats_include_specs_array_check check (jsonb_typeof(include_specs) = 'array'),
    constraint custom_formats_exclude_specs_array_check check (jsonb_typeof(exclude_specs) = 'array')
);

create index if not exists idx_custom_formats_name_lower
    on app.custom_formats (lower(name));

create table if not exists app.media_profile_custom_formats (
    profile_id text not null references app.media_profiles(id) on delete cascade,
    custom_format_id uuid not null references app.custom_formats(id) on delete cascade,
    score integer not null default 0,
    primary key (profile_id, custom_format_id)
);

create index if not exists idx_media_profile_custom_formats_profile
    on app.media_profile_custom_formats (profile_id, custom_format_id);

create table if not exists app.media_profile_video_targets (
    profile_id text primary key references app.media_profiles(id) on delete cascade,
    codecs text[] not null default '{}',
    codec_required boolean not null default false,
    codec_score integer not null default 0,
    hdr_formats text[] not null default '{}',
    hdr_required boolean not null default false,
    hdr_score integer not null default 0,
    pixel_formats text[] not null default '{}',
    pixel_format_required boolean not null default false,
    pixel_format_score integer not null default 0
);

create table if not exists app.media_profile_audio_targets (
    profile_id text not null references app.media_profiles(id) on delete cascade,
    language_id text not null,
    score integer not null default 0,
    target_codec text,
    target_channels text[] not null default '{}',
    minimum_bitrate_kbps integer check (minimum_bitrate_kbps is null or minimum_bitrate_kbps > 0),
    preferred_bitrate_kbps integer check (preferred_bitrate_kbps is null or preferred_bitrate_kbps > 0),
    sort_order integer not null default 0,
    primary key (profile_id, language_id)
);

create index if not exists idx_media_profile_audio_targets_profile
    on app.media_profile_audio_targets (profile_id, sort_order, language_id);

create table if not exists app.media_profile_subtitle_targets (
    profile_id text not null references app.media_profiles(id) on delete cascade,
    language_id text not null,
    score integer not null default 0,
    formats text[] not null default '{}',
    sort_order integer not null default 0,
    primary key (profile_id, language_id)
);

create index if not exists idx_media_profile_subtitle_targets_profile
    on app.media_profile_subtitle_targets (profile_id, sort_order, language_id);

create table if not exists app.languages (
    code text primary key,
    display_name text not null,
    aliases jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint languages_code_check check (code = upper(code) and length(code) between 2 and 8),
    constraint languages_display_name_check check (length(trim(display_name)) > 0),
    constraint languages_aliases_array_check check (jsonb_typeof(aliases) = 'array')
);

create unique index if not exists idx_languages_display_name_lower
    on app.languages (lower(display_name));

create table if not exists app.media_profile_qualities (
    profile_id text not null references app.media_profiles(id) on delete cascade,
    quality_id text not null,
    sort_order integer not null default 0,
    primary key (profile_id, quality_id)
);

create index if not exists idx_media_profile_qualities_profile_sort
    on app.media_profile_qualities (profile_id, sort_order, quality_id);

create table if not exists app.file_naming_settings (
    id integer primary key default 1 check (id = 1),
    movie_file_format text not null,
    movie_folder_format text not null,
    series_episode_format text not null,
    daily_episode_format text not null,
    anime_episode_format text not null,
    series_folder_format text not null,
    season_folder_format text not null,
    specials_folder_format text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists app.file_delete_settings (
    id boolean primary key default true check (id),
    mode text not null default 'permanent' check (mode in ('permanent', 'recycle', 'keep')),
    recycle_folder text not null default '.recycle',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

insert into app.file_delete_settings (id, mode, recycle_folder)
values (true, 'permanent', '.recycle')
on conflict (id) do nothing;

create table if not exists app.media_items (
    id uuid primary key,
    media_type text not null check (media_type in ('movie', 'serie')),
    content_kind text not null default 'standard' check (content_kind in ('standard', 'anime')),
    title text not null,
    year integer,
    monitored boolean not null default true,
    external_provider text,
    external_id text,
    overview text,
    poster_path text,
    collection_id text,
    collection_name text,
    backdrop_path text,
    metadata_status text,
    original_language text,
    series_type text check (series_type in ('standard', 'daily', 'absolute')),
    numbering_strategy text check (numbering_strategy in ('tmdb_season_episode', 'tvdb_season_episode', 'anidb_absolute', 'manual')),
    release_date text,
    first_air_date text,
    runtime_minutes integer,
    season_count integer,
    episode_count integer,
    vote_average double precision,
    genres jsonb not null default '[]'::jsonb check (jsonb_typeof(genres) = 'array'),
    keywords jsonb not null default '[]'::jsonb check (jsonb_typeof(keywords) = 'array'),
    facts jsonb not null default '[]'::jsonb check (jsonb_typeof(facts) = 'array'),
    seasons jsonb not null default '[]'::jsonb check (jsonb_typeof(seasons) = 'array'),
    cast_members jsonb not null default '[]'::jsonb check (jsonb_typeof(cast_members) = 'array'),
    crew_members jsonb not null default '[]'::jsonb check (jsonb_typeof(crew_members) = 'array'),
    recommendations jsonb not null default '[]'::jsonb check (jsonb_typeof(recommendations) = 'array'),
    similar_media jsonb not null default '[]'::jsonb check (jsonb_typeof(similar_media) = 'array'),
    quality_profile_id text,
    media_folder_path text,
    monitor_mode text not null default 'only_media' check (monitor_mode in ('none', 'only_media', 'collection', 'all_episodes', 'future_episodes', 'missing_episodes', 'existing_episodes', 'no_specials')),
    minimum_availability text not null default 'released' check (minimum_availability in ('announced', 'in_cinema', 'released')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_items_type_title
    on app.media_items (media_type, title);

create table if not exists app.media_seasons (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    external_provider text,
    external_id text,
    season_number integer not null check (season_number >= 0),
    name text not null default '',
    overview text,
    air_date text,
    poster_path text,
    episode_count integer,
    monitored boolean not null default true,
    source jsonb not null default '{}'::jsonb check (jsonb_typeof(source) = 'object'),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (media_item_id, season_number)
);

create index if not exists idx_media_seasons_media_item_id
    on app.media_seasons (media_item_id, season_number);

create table if not exists app.media_episodes (
    id uuid primary key,
    season_id uuid not null references app.media_seasons(id) on delete cascade,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    external_provider text,
    external_id text,
    season_number integer not null check (season_number >= 0),
    episode_number integer not null check (episode_number >= 0),
    name text not null default '',
    overview text,
    air_date text,
    still_path text,
    runtime_minutes integer,
    monitored boolean not null default true,
    source jsonb not null default '{}'::jsonb check (jsonb_typeof(source) = 'object'),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (media_item_id, season_number, episode_number)
);

create index if not exists idx_media_episodes_media_item_id
    on app.media_episodes (media_item_id, season_number, episode_number);

create index if not exists idx_media_episodes_season_id
    on app.media_episodes (season_id, episode_number);

create table if not exists app.media_provider_mappings (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    season_id uuid references app.media_seasons(id) on delete cascade,
    episode_id uuid references app.media_episodes(id) on delete cascade,
    entity_type text not null check (entity_type in ('media_item', 'season', 'episode')),
    provider_name text not null check (provider_name in ('tmdb', 'tvdb', 'anilist', 'anidb')),
    provider_entity_type text not null,
    external_id text not null,
    canonical boolean not null default false,
    confidence double precision,
    source jsonb not null default '{}'::jsonb check (jsonb_typeof(source) = 'object'),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint media_provider_mapping_entity_check check (
        (entity_type = 'media_item' and season_id is null and episode_id is null)
        or (entity_type = 'season' and season_id is not null and episode_id is null)
        or (entity_type = 'episode' and episode_id is not null)
    )
);

create index if not exists idx_media_provider_mappings_media_item
    on app.media_provider_mappings (media_item_id, provider_name, external_id);

create unique index if not exists idx_media_provider_mappings_unique_entity
    on app.media_provider_mappings (
        media_item_id,
        coalesce(season_id, '00000000-0000-0000-0000-000000000000'::uuid),
        coalesce(episode_id, '00000000-0000-0000-0000-000000000000'::uuid),
        provider_name,
        provider_entity_type,
        external_id
    );

create table if not exists app.media_item_aliases (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    alias text not null,
    normalized_alias text not null,
    language text,
    alias_kind text not null check (alias_kind in ('canonical', 'romaji', 'english', 'native', 'synonym', 'release_title')),
    provider_name text check (provider_name in ('tmdb', 'tvdb', 'anilist', 'anidb')),
    provider_mapping_id uuid references app.media_provider_mappings(id) on delete set null,
    source jsonb not null default '{}'::jsonb check (jsonb_typeof(source) = 'object'),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_item_aliases_media_item
    on app.media_item_aliases (media_item_id, alias_kind);

create unique index if not exists idx_media_item_aliases_unique_value
    on app.media_item_aliases (
        media_item_id,
        normalized_alias,
        alias_kind,
        coalesce(language, ''),
        coalesce(provider_name, '')
    );

create table if not exists app.media_episode_numbering (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    season_id uuid references app.media_seasons(id) on delete cascade,
    episode_id uuid not null references app.media_episodes(id) on delete cascade,
    provider_name text not null check (provider_name in ('tmdb', 'tvdb', 'anilist', 'anidb')),
    numbering_scheme text not null check (numbering_scheme in ('season_episode', 'absolute')),
    season_number integer check (season_number is null or season_number >= 0),
    episode_number integer check (episode_number is null or episode_number >= 0),
    absolute_number integer check (absolute_number is null or absolute_number >= 0),
    source jsonb not null default '{}'::jsonb check (jsonb_typeof(source) = 'object'),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint media_episode_numbering_value_check check (
        (numbering_scheme = 'season_episode' and season_number is not null and episode_number is not null)
        or (numbering_scheme = 'absolute' and absolute_number is not null)
    )
);

create unique index if not exists idx_media_episode_numbering_unique_scheme
    on app.media_episode_numbering (episode_id, provider_name, numbering_scheme);

create table if not exists app.tags (
    id uuid primary key,
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_tags_name_lower
    on app.tags (lower(name));

create table if not exists app.log_file_settings (
    id boolean primary key default true check (id),
    enabled boolean not null default false,
    directory text not null default '.data/logs',
    retention_days integer not null default 7 check (retention_days between 1 and 365),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

insert into app.log_file_settings (id, enabled, directory, retention_days)
values (true, false, '.data/logs', 7)
on conflict (id) do nothing;

create table if not exists app.system_events (
    id uuid primary key,
    severity text not null check (severity in ('info', 'warning', 'error')),
    category text not null,
    message text not null,
    data jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now()
);

create index if not exists idx_system_events_created_at
    on app.system_events (created_at desc);

create table if not exists app.system_event_settings (
    id boolean primary key default true check (id),
    retention_days integer not null default 7 check (retention_days between 1 and 365),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

insert into app.system_event_settings (id, retention_days)
values (true, 7)
on conflict (id) do nothing;

create table if not exists app.media_item_tags (
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    tag_id uuid not null references app.tags(id) on delete cascade,
    primary key (media_item_id, tag_id)
);

create table if not exists app.metadata_providers (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('tmdb', 'tvdb', 'anilist', 'anidb')),
    base_url text not null,
    api_key text,
    pin text,
    access_token text,
    session_token text,
    session_token_expires_at timestamptz,
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_metadata_providers_priority
    on app.metadata_providers (priority, name);

-- +goose StatementBegin
do $$
begin
    alter table app.metadata_providers drop constraint if exists metadata_providers_type_check;
    alter table app.metadata_providers
        add constraint metadata_providers_type_check check (type in ('tmdb', 'tvdb', 'anilist', 'anidb'));
end $$;
-- +goose StatementEnd

create table if not exists app.metadata_search_cache (
    provider_id uuid not null references app.metadata_providers(id) on delete cascade,
    media_type text not null check (media_type in ('movie', 'serie', 'mixed', 'person')),
    query text not null,
    year integer not null default 0,
    results jsonb not null,
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    primary key (provider_id, media_type, query, year)
);

alter table app.metadata_search_cache
    alter column year set default 0;

update app.metadata_search_cache
set year = 0
where year is null;

alter table app.metadata_search_cache
    alter column year set not null;

create index if not exists idx_metadata_search_cache_expires
    on app.metadata_search_cache (expires_at);

create table if not exists app.metadata_search_history (
    id uuid primary key,
    provider_id uuid references app.metadata_providers(id) on delete set null,
    provider_name text not null,
    provider_type text not null,
    media_type text not null check (media_type in ('movie', 'serie', 'mixed', 'person')),
    query text not null,
    year integer not null default 0,
    cache_hit boolean not null,
    success boolean not null,
    item_count integer not null default 0 check (item_count >= 0),
    error text,
    response jsonb not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_metadata_search_history_created_at
    on app.metadata_search_history (created_at desc);

create table if not exists app.subtitle_providers (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('opensubtitles')),
    base_url text not null,
    username text,
    password text,
    api_key text,
    enabled boolean not null default true,
    priority integer not null default 100 check (priority between 0 and 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_subtitle_providers_type_name_lower
    on app.subtitle_providers (type, lower(name));

create table if not exists app.library_folders (
    id uuid primary key,
    path text not null unique,
    kind text not null default 'movie' check (kind in ('movie', 'series')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists app.path_mappings (
    id uuid primary key,
    client_path text not null,
    app_path text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_path_mappings_client_path
    on app.path_mappings (client_path);

alter table app.media_items
    add column if not exists library_folder_id uuid references app.library_folders(id) on delete set null;

alter table app.media_items
    add column if not exists series_type text;

alter table app.media_items
    add column if not exists content_kind text not null default 'standard';

alter table app.media_items
    add column if not exists numbering_strategy text;

alter table app.media_items
    add column if not exists keywords jsonb not null default '[]'::jsonb;

alter table app.media_items
    add column if not exists recommendations jsonb not null default '[]'::jsonb;

alter table app.media_items
    add column if not exists similar_media jsonb not null default '[]'::jsonb;

-- +goose StatementBegin
do $$
begin
    alter table app.media_items drop constraint if exists media_items_monitor_mode_check;
    alter table app.media_items drop constraint if exists media_items_series_type_check;
    alter table app.media_items drop constraint if exists media_items_content_kind_check;
    alter table app.media_items drop constraint if exists media_items_numbering_strategy_check;
    alter table app.media_items drop constraint if exists media_items_keywords_check;
    alter table app.media_items drop constraint if exists media_items_recommendations_check;
    alter table app.media_items drop constraint if exists media_items_similar_media_check;
    alter table app.media_items
        add constraint media_items_monitor_mode_check check (monitor_mode in ('none', 'only_media', 'collection', 'all_episodes', 'future_episodes', 'missing_episodes', 'existing_episodes', 'no_specials'));
    alter table app.media_items
        add constraint media_items_series_type_check check (series_type is null or series_type in ('standard', 'daily', 'absolute'));
    alter table app.media_items
        add constraint media_items_content_kind_check check (content_kind in ('standard', 'anime'));
    alter table app.media_items
        add constraint media_items_numbering_strategy_check check (numbering_strategy is null or numbering_strategy in ('tmdb_season_episode', 'tvdb_season_episode', 'anidb_absolute', 'manual'));
    alter table app.media_items
        add constraint media_items_keywords_check check (jsonb_typeof(keywords) = 'array');
    alter table app.media_items
        add constraint media_items_recommendations_check check (jsonb_typeof(recommendations) = 'array');
    alter table app.media_items
        add constraint media_items_similar_media_check check (jsonb_typeof(similar_media) = 'array');
end $$;
-- +goose StatementEnd

create table if not exists app.media_requests (
    id uuid primary key,
    requested_by_user_id uuid not null references app.users(id) on delete cascade,
    media_type text not null check (media_type in ('movie', 'serie')),
    title text not null,
    year integer,
    external_provider text,
    external_id text,
    overview text,
    poster_path text,
    monitor_mode text not null default 'only_media',
    series_type text,
    minimum_availability text not null default 'released',
    status text not null default 'pending',
    quality_profile_id text,
    library_folder_id uuid references app.library_folders(id) on delete set null,
    media_item_id uuid references app.media_items(id) on delete set null,
    decided_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose StatementBegin
do $$
begin
    alter table app.media_requests drop constraint if exists media_requests_status_check;
    alter table app.media_requests add column if not exists monitor_mode text not null default 'only_media';
    alter table app.media_requests add column if not exists series_type text;
    alter table app.media_requests add column if not exists minimum_availability text not null default 'released';
    alter table app.media_requests drop constraint if exists media_requests_monitor_mode_check;
    alter table app.media_requests drop constraint if exists media_requests_series_type_check;
    alter table app.media_requests drop constraint if exists media_requests_minimum_availability_check;
    alter table app.media_requests
        add constraint media_requests_status_check check (status in ('pending', 'approved'));
    alter table app.media_requests
        add constraint media_requests_monitor_mode_check check (monitor_mode in ('none', 'only_media', 'collection', 'all_episodes', 'future_episodes', 'missing_episodes', 'existing_episodes', 'no_specials'));
    alter table app.media_requests
        add constraint media_requests_series_type_check check (series_type is null or series_type in ('standard', 'daily', 'absolute'));
    alter table app.media_requests
        add constraint media_requests_minimum_availability_check check (minimum_availability in ('announced', 'in_cinema', 'released'));
end $$;
-- +goose StatementEnd

create index if not exists idx_media_requests_status_created
    on app.media_requests (status, created_at desc);

create index if not exists idx_media_requests_requested_by
    on app.media_requests (requested_by_user_id, created_at desc);

create table if not exists app.media_request_tags (
    media_request_id uuid not null references app.media_requests(id) on delete cascade,
    tag_id uuid not null references app.tags(id) on delete cascade,
    primary key (media_request_id, tag_id)
);

create table if not exists app.library_scans (
    id uuid primary key,
    library_folder_id uuid not null references app.library_folders(id) on delete cascade,
    status text not null check (status in ('completed', 'failed')),
    total_files integer not null default 0 check (total_files >= 0),
    auto_matched_count integer not null default 0 check (auto_matched_count >= 0),
    manual_count integer not null default 0 check (manual_count >= 0),
    created_at timestamptz not null default now(),
    completed_at timestamptz
);

create index if not exists idx_library_scans_folder_created
    on app.library_scans (library_folder_id, created_at desc);

create table if not exists app.library_scan_items (
    id uuid primary key,
    scan_id uuid not null references app.library_scans(id) on delete cascade,
    path text not null,
    file_name text not null,
    size_bytes bigint not null default 0 check (size_bytes >= 0),
    detected_title text not null,
    detected_year integer,
    detected_media_kind text not null check (detected_media_kind in ('movie', 'series', 'anime_movie', 'anime_series', 'unknown')),
    season_number integer check (season_number is null or season_number >= 0),
    episode_number integer check (episode_number is null or episode_number >= 0),
    status text not null check (status in ('pending', 'auto_added', 'manually_added', 'missing', 'moved_candidate', 'restored')),
    imported boolean not null default false,
    matched_title text,
    matched_year integer,
    matched_media_kind text check (matched_media_kind in ('movie', 'series', 'anime_movie', 'anime_series', 'unknown')),
    matched_external_provider text,
    matched_external_id text,
    match_source text check (match_source is null or match_source in ('library', 'provider', 'manual')),
    selected_metadata_provider_id uuid references app.metadata_providers(id) on delete set null,
    duplicate_group_id text,
    duplicate_removal_allowed boolean not null default false,
    media_item_id uuid references app.media_items(id) on delete set null,
    season_id uuid references app.media_seasons(id) on delete set null,
    episode_id uuid references app.media_episodes(id) on delete set null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_library_scan_items_scan_status
    on app.library_scan_items (scan_id, status, path);

create index if not exists idx_library_scan_items_episode
    on app.library_scan_items (episode_id)
    where episode_id is not null;

create table if not exists app.download_activity (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    release_title text not null,
    indexer_name text not null,
    download_client_name text not null,
    download_id text,
    download_url text not null,
    status text not null check (status in ('queued', 'grabbed', 'downloading', 'completed', 'cancelled', 'failed')),
    progress_percent integer check (progress_percent is null or (progress_percent >= 0 and progress_percent <= 100)),
    error text,
    failure_type text check (failure_type is null or failure_type in ('download', 'import')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose StatementBegin
do $$
begin
    alter table app.download_activity add column if not exists download_id text;
    alter table app.download_activity add column if not exists progress_percent integer;
    alter table app.download_activity add column if not exists failure_type text;
    alter table app.download_activity drop constraint if exists download_activity_progress_percent_check;
    alter table app.download_activity drop constraint if exists download_activity_failure_type_check;
    alter table app.download_activity drop constraint if exists download_activity_status_check;
    alter table app.download_activity
        add constraint download_activity_status_check check (status in ('queued', 'grabbed', 'downloading', 'completed', 'cancelled', 'failed'));
    alter table app.download_activity
        add constraint download_activity_progress_percent_check check (progress_percent is null or (progress_percent >= 0 and progress_percent <= 100));
    alter table app.download_activity
        add constraint download_activity_failure_type_check check (failure_type is null or failure_type in ('download', 'import'));
end $$;
-- +goose StatementEnd

create index if not exists idx_download_activity_created
    on app.download_activity (created_at desc);

create table if not exists app.import_attempts (
    id uuid primary key,
    activity_id uuid not null,
    media_item_id uuid not null,
    source_path text,
    target_path text,
    import_mode text not null default 'hardlink' check (import_mode in ('hardlink', 'copy', 'move')),
    status text not null check (status in ('succeeded', 'failed')),
    failure_stage text check (failure_stage in ('load_media_item', 'load_path_mappings', 'select_source', 'create_media_folder', 'file_operation', 'record_media_file')),
    error_message text,
    created_targets text[] not null default '{}',
    inserted_media_file_paths text[] not null default '{}',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_import_attempts_activity_created
    on app.import_attempts (activity_id, created_at desc);

create index if not exists idx_import_attempts_media_item_created
    on app.import_attempts (media_item_id, created_at desc);

create table if not exists app.media_release_candidates (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    season_id uuid references app.media_seasons(id) on delete set null,
    episode_id uuid references app.media_episodes(id) on delete set null,
    indexer_id uuid references app.indexers(id) on delete set null,
    indexer_name text not null,
    indexer_type text not null default 'torznab',
    indexer_protocol text not null default 'torrent',
    title text not null,
    download_url text not null,
    info_url text,
    guid text,
    size_bytes bigint not null default 0,
    seeders integer,
    peers integer,
    published_at timestamptz,
    search_kind text not null default 'title',
    requested_season integer,
    requested_episode integer,
    sources jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_release_candidates_media_item
    on app.media_release_candidates (media_item_id, created_at desc);

create index if not exists idx_media_release_candidates_episode
    on app.media_release_candidates (episode_id)
    where episode_id is not null;

create table if not exists app.release_blocklist (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    release_title text not null,
    indexer_name text not null,
    indexer_type text not null default '',
    indexer_protocol text not null default 'torrent',
    download_client_name text not null default '',
    download_url text,
    info_url text,
    guid text,
    reason text not null,
    source text not null,
    temporary boolean not null default true,
    expires_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint release_blocklist_expiry_check check (
        (temporary = false and expires_at is null) or (temporary = true and expires_at is not null)
    )
);

alter table app.release_blocklist
    add column if not exists download_client_name text not null default '';

create index if not exists idx_release_blocklist_media
    on app.release_blocklist (media_item_id, created_at desc);

create index if not exists idx_release_blocklist_expiry
    on app.release_blocklist (temporary, expires_at);

create table if not exists app.media_release_search_errors (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    message text not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_media_release_search_errors_media_item
    on app.media_release_search_errors (media_item_id, created_at desc);

create table if not exists app.media_item_subtitles (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    season_id uuid references app.media_seasons(id) on delete set null,
    episode_id uuid references app.media_episodes(id) on delete set null,
    provider_id uuid references app.subtitle_providers(id) on delete set null,
    provider_name text not null,
    language_id text not null,
    format text not null default 'srt',
    file_path text not null,
    source_url text,
    source_reference text,
    release_name text,
    provider_subtitle_id text,
    checksum text,
    size_bytes bigint check (size_bytes is null or size_bytes >= 0),
    downloaded_at timestamptz not null default now(),
    selected boolean not null default true,
    retention_mode text not null default 'external'
        check (retention_mode in ('external', 'mux', 'ignore')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (media_item_id, language_id, file_path)
);

alter table app.media_item_subtitles
    add column if not exists format text not null default 'srt';

alter table app.media_item_subtitles
    add column if not exists source_reference text;

alter table app.media_item_subtitles
    add column if not exists provider_subtitle_id text;

alter table app.media_item_subtitles
    add column if not exists checksum text;

alter table app.media_item_subtitles
    add column if not exists size_bytes bigint;

alter table app.media_item_subtitles
    add column if not exists downloaded_at timestamptz not null default now();

alter table app.media_item_subtitles
    add column if not exists selected boolean not null default true;

alter table app.media_item_subtitles
    add column if not exists retention_mode text not null default 'external';

create index if not exists idx_media_item_subtitles_media_item
    on app.media_item_subtitles (media_item_id, language_id);

create index if not exists idx_media_item_subtitles_episode
    on app.media_item_subtitles (episode_id)
    where episode_id is not null;

create table if not exists app.media_component_sources (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    source_role text not null check (source_role in ('baseVideo', 'audio', 'subtitle', 'other')),
    source_file_path text not null,
    retained_path text not null,
    release_title text,
    release_group text,
    release_name text,
    release_id text,
    source_metadata text,
    stream_inventory text not null default '',
    checksum text,
    size_bytes bigint check (size_bytes is null or size_bytes >= 0),
    retention_state text not null default 'retained' check (retention_state in ('retained', 'released')),
    retained_at timestamptz not null default now(),
    released_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_component_sources_media
    on app.media_component_sources (media_item_id, retained_at desc);

create table if not exists app.media_component_provenance (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    component_type text not null check (component_type in ('container', 'video', 'audio', 'subtitle', 'chapter', 'attachment', 'sidecar')),
    component_key text not null,
    release_group text not null,
    release_name text not null,
    release_id text,
    source_provider text,
    source_file_path text,
    retained_source_id uuid references app.media_component_sources(id) on delete set null,
    source_stream_id integer,
    transformation_chain jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint media_component_provenance_chain_array_check check (jsonb_typeof(transformation_chain) = 'array')
);

create unique index if not exists idx_media_component_provenance_component
    on app.media_component_provenance (media_item_id, component_type, component_key);

create table if not exists app.media_component_artifacts (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    source_id uuid not null references app.media_component_sources(id) on delete cascade,
    stream_id integer not null check (stream_id >= 0),
    stream_type text not null check (stream_type in ('audio', 'subtitle')),
    language text,
    output_path text not null,
    status text not null default 'queued' check (status in ('queued', 'running', 'succeeded', 'failed')),
    tool_name text not null default 'mkvextract',
    tool_summary text not null default '',
    error_message text,
    job_id text,
    size_bytes bigint check (size_bytes is null or size_bytes >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    completed_at timestamptz
);

create index if not exists idx_media_component_artifacts_source
    on app.media_component_artifacts (source_id, created_at desc);

create table if not exists app.media_component_compatibility_decisions (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    base_source_id uuid not null references app.media_component_sources(id) on delete cascade,
    component_source_id uuid not null references app.media_component_sources(id) on delete cascade,
    confidence_state text not null check (confidence_state in ('exact', 'likely', 'uncertain', 'incompatible')),
    automation_state text not null check (automation_state in ('allowed', 'blocked')),
    review_state text not null check (review_state in ('notRequired', 'pending', 'approved', 'rejected')),
    reason text not null,
    runtime_delta_ms integer,
    evidence jsonb not null default '{}'::jsonb check (jsonb_typeof(evidence) = 'object'),
    review_reason text,
    reviewed_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (media_item_id, base_source_id, component_source_id)
);

create index if not exists idx_media_component_compatibility_component
    on app.media_component_compatibility_decisions (component_source_id, created_at desc);

create table if not exists app.media_component_assembly_runs (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    base_source_id uuid not null references app.media_component_sources(id) on delete cascade,
    output_path text not null,
    status text not null default 'queued' check (status in ('queued', 'running', 'succeeded', 'failed')),
    tool_name text not null default 'mkvmerge',
    tool_summary text not null default '',
    error_message text,
    job_id text,
    size_bytes bigint check (size_bytes is null or size_bytes >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    completed_at timestamptz
);

create table if not exists app.media_component_assembly_inputs (
    id uuid primary key,
    run_id uuid not null references app.media_component_assembly_runs(id) on delete cascade,
    source_id uuid references app.media_component_sources(id) on delete set null,
    artifact_id uuid references app.media_component_artifacts(id) on delete set null,
    stream_type text not null check (stream_type in ('video', 'audio', 'subtitle')),
    input_path text not null,
    provenance jsonb not null default '{}'::jsonb check (jsonb_typeof(provenance) = 'object'),
    created_at timestamptz not null default now()
);

create index if not exists idx_media_component_assembly_runs_media
    on app.media_component_assembly_runs (media_item_id, created_at desc);

create index if not exists idx_media_component_assembly_inputs_run
    on app.media_component_assembly_inputs (run_id, created_at);

create table if not exists app.media_file_history (
    id uuid primary key,
    media_item_id uuid references app.media_items(id) on delete set null,
    file_path text not null,
    source_path text,
    destination_path text,
    operation text not null check (operation in ('imported', 'renamed', 'moved', 'replaced', 'deleted', 'missing', 'moved_candidate', 'restored', 'superseded')),
    status text not null check (status in ('succeeded', 'failed', 'skipped')),
    actor_type text not null default 'system' check (actor_type in ('system', 'user', 'job')),
    actor_id text,
    job_id text,
    details jsonb not null default '{}'::jsonb check (jsonb_typeof(details) = 'object'),
    failure_details text,
    created_at timestamptz not null default now()
);

create index if not exists idx_media_file_history_media_item
    on app.media_file_history (media_item_id, created_at desc);

create index if not exists idx_media_file_history_file_path
    on app.media_file_history (file_path, created_at desc);

create table if not exists app.media_item_sidecars (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    media_file_path text not null,
    file_path text not null,
    sidecar_type text not null check (sidecar_type in ('metadata', 'subtitle')),
    language_id text,
    format text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (media_item_id, file_path)
);

create index if not exists idx_media_item_sidecars_media
    on app.media_item_sidecars (media_item_id, media_file_path, sidecar_type);
