create schema if not exists app;

create table if not exists app.users (
    id uuid primary key,
    username text not null unique,
    password_hash text not null,
    role text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

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
    year integer,
    monitored boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_media_items_type_title
    on app.media_items (media_type, title);

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
