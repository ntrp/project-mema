-- +goose Up
alter table app.media_profiles
    add column if not exists remove_non_enabled_languages boolean not null default false;

alter table app.media_profile_languages
    add column if not exists required boolean not null default false;

-- +goose Down
alter table app.media_profile_languages
    drop column if exists required;

alter table app.media_profiles
    drop column if exists remove_non_enabled_languages;
