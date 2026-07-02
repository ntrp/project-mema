-- +goose Up
create table if not exists app.metadata_search_history (
    id uuid primary key,
    provider_id uuid references app.metadata_providers(id) on delete set null,
    provider_name text not null,
    provider_type text not null,
    media_type text not null check (media_type in ('movie', 'series', 'mixed')),
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
