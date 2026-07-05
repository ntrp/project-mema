-- name: ListMediaProfiles :many
select id,
    name,
    is_default,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_non_enabled_languages,
    remove_non_enabled_subtitle_languages,
    preferred_protocol,
    series_pack_preference,
    created_at,
    updated_at
from app.media_profiles
order by lower(name);

-- name: MediaProfileExists :one
select exists(select 1 from app.media_profiles where id = $1);

-- name: CreateMediaProfile :exec
insert into app.media_profiles (
    id,
    name,
    is_default,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_non_enabled_languages,
    remove_non_enabled_subtitle_languages,
    preferred_protocol,
    series_pack_preference
)
values (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(is_default),
    sqlc.arg(upgrades_allowed),
    sqlc.narg(upgrade_until_quality_id),
    sqlc.arg(minimum_custom_format_score),
    sqlc.arg(upgrade_until_custom_format_score),
    sqlc.arg(minimum_custom_format_score_increment),
    sqlc.arg(remove_non_enabled_languages),
    sqlc.arg(remove_non_enabled_subtitle_languages),
    sqlc.arg(preferred_protocol),
    sqlc.arg(series_pack_preference)
);

-- name: UpdateMediaProfile :execrows
update app.media_profiles
set name = sqlc.arg(name),
    is_default = sqlc.arg(is_default),
    upgrades_allowed = sqlc.arg(upgrades_allowed),
    upgrade_until_quality_id = sqlc.narg(upgrade_until_quality_id),
    minimum_custom_format_score = sqlc.arg(minimum_custom_format_score),
    upgrade_until_custom_format_score = sqlc.arg(upgrade_until_custom_format_score),
    minimum_custom_format_score_increment = sqlc.arg(minimum_custom_format_score_increment),
    remove_non_enabled_languages = sqlc.arg(remove_non_enabled_languages),
    remove_non_enabled_subtitle_languages = sqlc.arg(remove_non_enabled_subtitle_languages),
    preferred_protocol = sqlc.arg(preferred_protocol),
    series_pack_preference = sqlc.arg(series_pack_preference),
    updated_at = now()
where id = sqlc.arg(id);

-- name: DeleteMediaProfile :execrows
delete from app.media_profiles
where id = $1;

-- name: GetMediaProfile :one
select id,
    name,
    is_default,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_non_enabled_languages,
    remove_non_enabled_subtitle_languages,
    preferred_protocol,
    series_pack_preference,
    created_at,
    updated_at
from app.media_profiles
where id = $1;

-- name: ClearDefaultMediaProfiles :exec
update app.media_profiles
set is_default = false,
    updated_at = now()
where is_default;

-- name: ListMediaProfileQualities :many
select quality_id
from app.media_profile_qualities
where profile_id = $1
order by sort_order, quality_id;

-- name: ListMediaProfileLanguages :many
select language_id, score, required
from app.media_profile_languages
where profile_id = $1
order by language_id;

-- name: ListMediaProfileSubtitleLanguages :many
select language_id, score, required, subtitle_type
from app.media_profile_subtitle_languages
where profile_id = $1
order by language_id;

-- name: ListMediaProfileCustomFormats :many
select pcf.custom_format_id, pcf.score
from app.media_profile_custom_formats pcf
join app.custom_formats cf on cf.id = pcf.custom_format_id
where pcf.profile_id = $1
order by lower(cf.name), pcf.custom_format_id;

-- name: ClearMediaProfileQualities :exec
delete from app.media_profile_qualities
where profile_id = $1;

-- name: AddMediaProfileQuality :exec
insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
values (sqlc.arg(profile_id), sqlc.arg(quality_id), sqlc.arg(sort_order));

-- name: ClearMediaProfileLanguages :exec
delete from app.media_profile_languages
where profile_id = $1;

-- name: AddMediaProfileLanguage :exec
insert into app.media_profile_languages (profile_id, language_id, score, required)
values (sqlc.arg(profile_id), sqlc.arg(language_id), sqlc.arg(score), sqlc.arg(required));

-- name: ClearMediaProfileSubtitleLanguages :exec
delete from app.media_profile_subtitle_languages
where profile_id = $1;

-- name: AddMediaProfileSubtitleLanguage :exec
insert into app.media_profile_subtitle_languages (profile_id, language_id, score, required, subtitle_type)
values (sqlc.arg(profile_id), sqlc.arg(language_id), sqlc.arg(score), sqlc.arg(required), sqlc.arg(subtitle_type));

-- name: ClearMediaProfileCustomFormats :exec
delete from app.media_profile_custom_formats
where profile_id = $1;

-- name: AddMediaProfileCustomFormat :exec
insert into app.media_profile_custom_formats (profile_id, custom_format_id, score)
values (sqlc.arg(profile_id), sqlc.arg(custom_format_id), sqlc.arg(score));
