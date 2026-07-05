-- +goose Up
alter table app.media_release_candidates
    add column if not exists sources jsonb not null default '[]'::jsonb;

update app.media_release_candidates
set sources = jsonb_build_array(jsonb_strip_nulls(jsonb_build_object(
    'indexerId', indexer_id,
    'indexerName', indexer_name,
    'indexerProtocol', indexer_protocol,
    'title', title,
    'downloadUrl', download_url,
    'infoUrl', info_url,
    'guid', guid
)))
where sources = '[]'::jsonb;
