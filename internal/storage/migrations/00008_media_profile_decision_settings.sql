-- +goose Up
alter table app.media_profiles
    add column if not exists preferred_protocol text not null default 'any',
    add column if not exists series_pack_preference text not null default 'auto';

alter table app.media_profiles
    drop constraint if exists media_profiles_preferred_protocol_check,
    add constraint media_profiles_preferred_protocol_check
        check (preferred_protocol in ('any', 'torrent', 'usenet'));

alter table app.media_profiles
    drop constraint if exists media_profiles_series_pack_preference_check,
    add constraint media_profiles_series_pack_preference_check
        check (series_pack_preference in ('auto', 'preferPacks', 'preferEpisodes'));

-- +goose Down
alter table app.media_profiles
    drop column if exists preferred_protocol,
    drop column if exists series_pack_preference;
