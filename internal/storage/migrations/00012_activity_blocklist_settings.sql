-- +goose Up
alter table app.indexer_search_settings
    add column if not exists automatic_blocklist_expiry_days integer not null default 7;

alter table app.indexer_search_settings
    drop constraint if exists indexer_search_settings_automatic_blocklist_expiry_days_check;

alter table app.indexer_search_settings
    add constraint indexer_search_settings_automatic_blocklist_expiry_days_check
        check (automatic_blocklist_expiry_days between 1 and 365);

-- +goose Down
alter table app.indexer_search_settings
    drop column if exists automatic_blocklist_expiry_days;
