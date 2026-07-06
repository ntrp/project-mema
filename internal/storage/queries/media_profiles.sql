-- name: ListMediaProfiles :many
select id,
    name,
    is_default,
    final_container,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_unwanted_audio,
    audio_lossy_transcode_policy,
    remove_unwanted_subtitles,
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
    final_container,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_unwanted_audio,
    audio_lossy_transcode_policy,
    remove_unwanted_subtitles,
    preferred_protocol,
    series_pack_preference
)
values (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(is_default),
    sqlc.arg(final_container),
    sqlc.arg(upgrades_allowed),
    sqlc.narg(upgrade_until_quality_id),
    sqlc.arg(minimum_custom_format_score),
    sqlc.arg(upgrade_until_custom_format_score),
    sqlc.arg(minimum_custom_format_score_increment),
    sqlc.arg(remove_unwanted_audio),
    sqlc.arg(audio_lossy_transcode_policy),
    sqlc.arg(remove_unwanted_subtitles),
    sqlc.arg(preferred_protocol),
    sqlc.arg(series_pack_preference)
);

-- name: UpdateMediaProfile :execrows
update app.media_profiles
set name = sqlc.arg(name),
    is_default = sqlc.arg(is_default),
    final_container = sqlc.arg(final_container),
    upgrades_allowed = sqlc.arg(upgrades_allowed),
    upgrade_until_quality_id = sqlc.narg(upgrade_until_quality_id),
    minimum_custom_format_score = sqlc.arg(minimum_custom_format_score),
    upgrade_until_custom_format_score = sqlc.arg(upgrade_until_custom_format_score),
    minimum_custom_format_score_increment = sqlc.arg(minimum_custom_format_score_increment),
    remove_unwanted_audio = sqlc.arg(remove_unwanted_audio),
    audio_lossy_transcode_policy = sqlc.arg(audio_lossy_transcode_policy),
    remove_unwanted_subtitles = sqlc.arg(remove_unwanted_subtitles),
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
    final_container,
    upgrades_allowed,
    upgrade_until_quality_id,
    minimum_custom_format_score,
    upgrade_until_custom_format_score,
    minimum_custom_format_score_increment,
    remove_unwanted_audio,
    audio_lossy_transcode_policy,
    remove_unwanted_subtitles,
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

-- name: GetMediaProfileVideoTarget :one
select codecs,
    codec_required,
    codec_score,
    hdr_formats,
    hdr_required,
    hdr_score,
    pixel_formats,
    pixel_format_required,
    pixel_format_score
from app.media_profile_video_targets
where profile_id = $1
limit 1;

-- name: ListMediaProfileAudioTargets :many
select language_id,
    score,
    target_codec,
    target_channels,
    minimum_bitrate_kbps,
    preferred_bitrate_kbps
from app.media_profile_audio_targets
where profile_id = $1
order by sort_order, language_id;

-- name: ListMediaProfileSubtitleTargets :many
select language_id,
    score,
    source,
    formats
from app.media_profile_subtitle_targets
where profile_id = $1
order by sort_order, language_id;

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

-- name: UpsertMediaProfileVideoTarget :exec
insert into app.media_profile_video_targets (
    profile_id,
    codecs,
    codec_required,
    codec_score,
    hdr_formats,
    hdr_required,
    hdr_score,
    pixel_formats,
    pixel_format_required,
    pixel_format_score
)
values (
    sqlc.arg(profile_id),
    sqlc.arg(codecs),
    sqlc.arg(codec_required),
    sqlc.arg(codec_score),
    sqlc.arg(hdr_formats),
    sqlc.arg(hdr_required),
    sqlc.arg(hdr_score),
    sqlc.arg(pixel_formats),
    sqlc.arg(pixel_format_required),
    sqlc.arg(pixel_format_score)
)
on conflict (profile_id) do update
set codecs = excluded.codecs,
    codec_required = excluded.codec_required,
    codec_score = excluded.codec_score,
    hdr_formats = excluded.hdr_formats,
    hdr_required = excluded.hdr_required,
    hdr_score = excluded.hdr_score,
    pixel_formats = excluded.pixel_formats,
    pixel_format_required = excluded.pixel_format_required,
    pixel_format_score = excluded.pixel_format_score;

-- name: ClearMediaProfileAudioTargets :exec
delete from app.media_profile_audio_targets
where profile_id = $1;

-- name: AddMediaProfileAudioTarget :exec
insert into app.media_profile_audio_targets (
    profile_id,
    language_id,
    score,
    target_codec,
    target_channels,
    minimum_bitrate_kbps,
    preferred_bitrate_kbps,
    sort_order
)
values (
    sqlc.arg(profile_id),
    sqlc.arg(language_id),
    sqlc.arg(score),
    sqlc.narg(target_codec),
    sqlc.arg(target_channels),
    sqlc.narg(minimum_bitrate_kbps),
    sqlc.narg(preferred_bitrate_kbps),
    sqlc.arg(sort_order)
);

-- name: ClearMediaProfileSubtitleTargets :exec
delete from app.media_profile_subtitle_targets
where profile_id = $1;

-- name: AddMediaProfileSubtitleTarget :exec
insert into app.media_profile_subtitle_targets (
    profile_id,
    language_id,
    score,
    source,
    formats,
    sort_order
)
values (
    sqlc.arg(profile_id),
    sqlc.arg(language_id),
    sqlc.arg(score),
    sqlc.arg(source),
    sqlc.arg(formats),
    sqlc.arg(sort_order)
);

-- name: ClearMediaProfileCustomFormats :exec
delete from app.media_profile_custom_formats
where profile_id = $1;

-- name: AddMediaProfileCustomFormat :exec
insert into app.media_profile_custom_formats (profile_id, custom_format_id, score)
values (sqlc.arg(profile_id), sqlc.arg(custom_format_id), sqlc.arg(score));
