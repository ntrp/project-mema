-- +goose Up
alter table app.metadata_search_cache
    drop constraint if exists metadata_search_cache_media_type_check,
    add constraint metadata_search_cache_media_type_check
        check (media_type in ('movie', 'serie', 'mixed', 'person'));

alter table app.metadata_search_history
    drop constraint if exists metadata_search_history_media_type_check,
    add constraint metadata_search_history_media_type_check
        check (media_type in ('movie', 'serie', 'mixed', 'person'));

-- +goose Down
delete from app.metadata_search_cache
where media_type = 'person';

delete from app.metadata_search_history
where media_type = 'person';

alter table app.metadata_search_cache
    drop constraint if exists metadata_search_cache_media_type_check,
    add constraint metadata_search_cache_media_type_check
        check (media_type in ('movie', 'serie', 'mixed'));

alter table app.metadata_search_history
    drop constraint if exists metadata_search_history_media_type_check,
    add constraint metadata_search_history_media_type_check
        check (media_type in ('movie', 'serie', 'mixed'));
