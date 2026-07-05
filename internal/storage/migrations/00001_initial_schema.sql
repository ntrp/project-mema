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
    created_at timestamptz not null default now()
);

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
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint indexers_media_type_scopes_check check (media_type_scopes <@ array['movie', 'serie', 'anime', 'audio', 'book']::text[])
);

create index if not exists idx_indexers_priority
    on app.indexers (priority, name);

create index if not exists idx_indexers_health
    on app.indexers (enabled, health_status, next_check_at);

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
    upgrades_allowed boolean not null default true,
    upgrade_until_quality_id text,
    minimum_custom_format_score integer not null default 0,
    upgrade_until_custom_format_score integer not null default 0,
    minimum_custom_format_score_increment integer not null default 1 check (minimum_custom_format_score_increment >= 0),
    remove_non_enabled_languages boolean not null default false,
    preferred_protocol text not null default 'any' check (preferred_protocol in ('any', 'torrent', 'usenet')),
    series_pack_preference text not null default 'auto' check (series_pack_preference in ('auto', 'preferPacks', 'preferEpisodes')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_media_profiles_name_lower
    on app.media_profiles (lower(name));

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

create table if not exists app.media_profile_languages (
    profile_id text not null references app.media_profiles(id) on delete cascade,
    language_id text not null,
    score integer not null default 0,
    required boolean not null default false,
    primary key (profile_id, language_id)
);

create index if not exists idx_media_profile_languages_profile
    on app.media_profile_languages (profile_id, language_id);

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

create table if not exists app.media_items (
    id uuid primary key,
    media_type text not null check (media_type in ('movie', 'serie')),
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
    type text not null check (type in ('tmdb', 'tvdb')),
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

create table if not exists app.library_folders (
    id uuid primary key,
    path text not null unique,
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
    alter table app.media_items drop constraint if exists media_items_keywords_check;
    alter table app.media_items drop constraint if exists media_items_recommendations_check;
    alter table app.media_items drop constraint if exists media_items_similar_media_check;
    alter table app.media_items
        add constraint media_items_monitor_mode_check check (monitor_mode in ('none', 'only_media', 'collection', 'all_episodes', 'future_episodes', 'missing_episodes', 'existing_episodes', 'no_specials'));
    alter table app.media_items
        add constraint media_items_series_type_check check (series_type is null or series_type in ('standard', 'daily', 'absolute'));
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
    detected_title text not null,
    detected_year integer,
    detected_media_kind text not null check (detected_media_kind in ('movie', 'series', 'anime_movie', 'anime_series', 'unknown')),
    status text not null check (status in ('pending', 'auto_added', 'manually_added')),
    matched_title text,
    matched_year integer,
    matched_media_kind text check (matched_media_kind in ('movie', 'series', 'anime_movie', 'anime_series', 'unknown')),
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

create table if not exists app.media_release_candidates (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    season_id uuid references app.media_seasons(id) on delete set null,
    episode_id uuid references app.media_episodes(id) on delete set null,
    indexer_id uuid,
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
