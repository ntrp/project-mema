-- +goose Up
update app.indexer_search_settings
set history_retention_days = 7, updated_at = now()
where id = true and history_retention_days = 30;
