-- +goose Up
alter table app.indexers
    add column if not exists media_type_scopes text[] not null default '{movie,serie,anime,audio,book}',
    add column if not exists tag_scopes text[] not null default '{}';

update app.indexers
set media_type_scopes = array_replace(media_type_scopes, 'series', 'serie');

alter table app.indexers drop constraint if exists indexers_media_type_scopes_check;
alter table app.indexers
    add constraint indexers_media_type_scopes_check
    check (media_type_scopes <@ array['movie', 'serie', 'anime', 'audio', 'book']::text[]);

update app.media_items set media_type = 'serie' where media_type = 'series';
update app.media_requests set media_type = 'serie' where media_type = 'series';
update app.discover_blacklist set media_type = 'serie' where media_type = 'series';
update app.indexer_search_cache set media_type = 'serie' where media_type = 'series';
update app.indexer_search_history set media_type = 'serie' where media_type = 'series';
update app.metadata_search_cache set media_type = 'serie' where media_type = 'series';
update app.metadata_search_history set media_type = 'serie' where media_type = 'series';
update app.media_release_candidates set search_kind = 'serie' where search_kind = 'series';

alter table app.media_items drop constraint if exists media_items_media_type_check;
alter table app.media_items
    add constraint media_items_media_type_check check (media_type in ('movie', 'serie'));

alter table app.media_requests drop constraint if exists media_requests_media_type_check;
alter table app.media_requests
    add constraint media_requests_media_type_check check (media_type in ('movie', 'serie'));

alter table app.discover_blacklist drop constraint if exists discover_blacklist_media_type_check;
alter table app.discover_blacklist
    add constraint discover_blacklist_media_type_check check (media_type in ('movie', 'serie'));

alter table app.indexer_search_cache drop constraint if exists indexer_search_cache_media_type_check;
alter table app.indexer_search_cache
    add constraint indexer_search_cache_media_type_check check (media_type in ('movie', 'serie', 'mixed'));

alter table app.indexer_search_history drop constraint if exists indexer_search_history_media_type_check;
alter table app.indexer_search_history
    add constraint indexer_search_history_media_type_check check (media_type in ('movie', 'serie', 'mixed'));

alter table app.metadata_search_cache drop constraint if exists metadata_search_cache_media_type_check;
alter table app.metadata_search_cache
    add constraint metadata_search_cache_media_type_check check (media_type in ('movie', 'serie', 'mixed', 'person'));

alter table app.metadata_search_history drop constraint if exists metadata_search_history_media_type_check;
alter table app.metadata_search_history
    add constraint metadata_search_history_media_type_check check (media_type in ('movie', 'serie', 'mixed', 'person'));
