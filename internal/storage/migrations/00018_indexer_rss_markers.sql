-- +goose Up
alter table app.indexers
    add column if not exists rss_marker_published_at timestamptz,
    add column if not exists rss_marker_guid text,
    add column if not exists rss_marker_download_url text;

create index if not exists idx_indexers_rss_enabled
    on app.indexers (enabled, supports_rss, health_status, next_check_at);
