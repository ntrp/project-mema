-- +goose Up
create table if not exists app.indexer_search_settings (
    id boolean primary key default true check (id),
    cache_duration_minutes integer not null default 1440 check (cache_duration_minutes between 0 and 43200),
    history_retention_days integer not null default 7 check (history_retention_days between 1 and 365),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

insert into app.indexer_search_settings (id, cache_duration_minutes, history_retention_days)
values (true, 1440, 7)
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
    indexer_type text not null,
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
