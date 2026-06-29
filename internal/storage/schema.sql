create schema if not exists app;

create table if not exists app.users (
    id uuid primary key,
    username text not null unique,
    password_hash text not null,
    role text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

do $$
begin
    alter table app.users drop constraint if exists users_role_check;
    alter table app.users
        add constraint users_role_check check (role in ('admin', 'user'));
end $$;

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

create index if not exists idx_download_clients_priority
    on app.download_clients (priority, name);

create table if not exists app.indexers (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('torznab', 'newznab', 'rss')),
    base_url text not null,
    api_key text,
    categories integer[] not null default '{}',
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_indexers_priority
    on app.indexers (priority, name);

create table if not exists app.media_items (
    id uuid primary key,
    media_type text not null check (media_type in ('movie', 'series')),
    title text not null,
    year integer not null default 0,
    monitored boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

alter table app.media_items
    add column if not exists external_provider text;

alter table app.media_items
    add column if not exists external_id text;

alter table app.media_items
    add column if not exists overview text;

alter table app.media_items
    add column if not exists poster_path text;

alter table app.media_items
    add column if not exists quality_profile_id text;

alter table app.media_items
    alter column year drop not null;

create index if not exists idx_media_items_type_title
    on app.media_items (media_type, title);

create table if not exists app.tags (
    id uuid primary key,
    name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index if not exists idx_tags_name_lower
    on app.tags (lower(name));

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
    media_type text not null check (media_type in ('movie', 'series')),
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

create table if not exists app.library_folders (
    id uuid primary key,
    path text not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

alter table app.media_items
    add column if not exists library_folder_id uuid references app.library_folders(id) on delete set null;

create table if not exists app.media_requests (
    id uuid primary key,
    requested_by_user_id uuid not null references app.users(id) on delete cascade,
    media_type text not null check (media_type in ('movie', 'series')),
    title text not null,
    year integer,
    external_provider text,
    external_id text,
    overview text,
    poster_path text,
    status text not null default 'pending',
    quality_profile_id text,
    library_folder_id uuid references app.library_folders(id) on delete set null,
    media_item_id uuid references app.media_items(id) on delete set null,
    decided_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

do $$
begin
    alter table app.media_requests drop constraint if exists media_requests_status_check;
    alter table app.media_requests
        add constraint media_requests_status_check check (status in ('pending', 'approved'));
end $$;

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
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_library_scan_items_scan_status
    on app.library_scan_items (scan_id, status, path);

create table if not exists app.download_activity (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    release_title text not null,
    indexer_name text not null,
    download_client_name text not null,
    download_url text not null,
    status text not null check (status in ('queued', 'grabbed', 'failed')),
    error text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

do $$
begin
    alter table app.download_activity drop constraint if exists download_activity_status_check;
    alter table app.download_activity
        add constraint download_activity_status_check check (status in ('queued', 'grabbed', 'failed'));
end $$;

create index if not exists idx_download_activity_created
    on app.download_activity (created_at desc);

create table if not exists app.media_release_candidates (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    indexer_id uuid,
    indexer_name text not null,
    indexer_type text not null,
    title text not null,
    download_url text not null,
    info_url text,
    guid text,
    size_bytes bigint not null default 0,
    seeders integer,
    peers integer,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_release_candidates_media_item
    on app.media_release_candidates (media_item_id, created_at desc);

create table if not exists app.media_release_search_errors (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    message text not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_media_release_search_errors_media_item
    on app.media_release_search_errors (media_item_id, created_at desc);
