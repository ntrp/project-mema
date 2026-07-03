-- +goose Up
alter table app.indexers
    add column if not exists definition_id text,
    add column if not exists implementation text,
    add column if not exists implementation_name text,
    add column if not exists protocol text,
    add column if not exists privacy text not null default 'private',
    add column if not exists language text not null default 'en-US',
    add column if not exists encoding text,
    add column if not exists description text,
    add column if not exists indexer_urls text[] not null default '{}',
    add column if not exists legacy_urls text[] not null default '{}',
    add column if not exists fields jsonb not null default '[]'::jsonb,
    add column if not exists capabilities jsonb not null default '{"categories":[],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}'::jsonb,
    add column if not exists redirect boolean not null default true,
    add column if not exists app_profile_id text not null default 'default',
    add column if not exists minimum_seeders integer check (minimum_seeders is null or minimum_seeders >= 0),
    add column if not exists seed_ratio numeric(10, 2) check (seed_ratio is null or seed_ratio >= 0),
    add column if not exists seed_time integer check (seed_time is null or seed_time >= 0),
    add column if not exists pack_seed_time integer check (pack_seed_time is null or pack_seed_time >= 0),
    add column if not exists prefer_magnet_url boolean not null default false,
    add column if not exists supports_rss boolean not null default true,
    add column if not exists supports_search boolean not null default true,
    add column if not exists supports_redirect boolean not null default true,
    add column if not exists supports_pagination boolean not null default true;

update app.indexers
set protocol = case
        when type = 'newznab' then 'usenet'
        else 'torrent'
    end,
    definition_id = coalesce(definition_id, case
        when type = 'newznab' then 'generic-newznab'
        else 'generic-torznab'
    end),
    implementation = coalesce(implementation, case
        when type = 'newznab' then 'Newznab'
        else 'Cardigann'
    end),
    implementation_name = coalesce(implementation_name, name)
where protocol is null or definition_id is null or implementation is null or implementation_name is null;

alter table app.indexers
    alter column definition_id set not null,
    alter column implementation set not null,
    alter column implementation_name set not null,
    alter column type set default 'torznab',
    alter column protocol set not null;

alter table app.indexers drop constraint if exists indexers_protocol_check;
alter table app.indexers add constraint indexers_protocol_check check (protocol in ('torrent', 'usenet'));

alter table app.indexers drop constraint if exists indexers_privacy_check;
alter table app.indexers add constraint indexers_privacy_check check (privacy in ('public', 'private', 'semiPrivate'));

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

alter table app.indexer_search_history
    add column if not exists indexer_protocol text;

update app.indexer_search_history
set indexer_protocol = case
        when indexer_type = 'newznab' then 'usenet'
        else 'torrent'
    end
where indexer_protocol is null;

alter table app.indexer_search_history
    alter column indexer_protocol set not null,
    alter column indexer_protocol set default 'torrent',
    alter column indexer_type set default 'torznab';

alter table app.media_release_candidates
    add column if not exists indexer_protocol text;

update app.media_release_candidates
set indexer_protocol = case
        when indexer_type = 'newznab' then 'usenet'
        else 'torrent'
    end
where indexer_protocol is null;

alter table app.media_release_candidates
    alter column indexer_protocol set not null,
    alter column indexer_protocol set default 'torrent',
    alter column indexer_type set default 'torznab';

alter table app.release_blocklist
    add column if not exists indexer_protocol text;

update app.release_blocklist
set indexer_protocol = case
        when indexer_type = 'newznab' then 'usenet'
        when indexer_type = '' then 'torrent'
        else 'torrent'
    end
where indexer_protocol is null;

alter table app.release_blocklist
    alter column indexer_protocol set not null,
    alter column indexer_protocol set default 'torrent';
