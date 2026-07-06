-- name: ListMediaItemSubtitles :many
select *
from app.media_item_subtitles
where media_item_id = $1
order by language_id, created_at desc;

-- name: UpsertMediaItemSubtitle :one
insert into app.media_item_subtitles (
    id,
    media_item_id,
    season_id,
    episode_id,
    provider_id,
    provider_name,
    language_id,
    format,
    file_path,
    source_url,
    source_reference,
    release_name,
    provider_subtitle_id,
    checksum,
    size_bytes,
    downloaded_at,
    selected,
    retention_mode
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id),
    sqlc.narg(provider_id),
    sqlc.arg(provider_name),
    sqlc.arg(language_id),
    sqlc.arg(format),
    sqlc.arg(file_path),
    sqlc.narg(source_url),
    sqlc.narg(source_reference),
    sqlc.narg(release_name),
    sqlc.narg(provider_subtitle_id),
    sqlc.narg(checksum),
    sqlc.narg(size_bytes),
    sqlc.arg(downloaded_at),
    sqlc.arg(selected),
    sqlc.arg(retention_mode)
)
on conflict (media_item_id, language_id, file_path) do update
set provider_id = excluded.provider_id,
    provider_name = excluded.provider_name,
    format = excluded.format,
    source_url = excluded.source_url,
    source_reference = excluded.source_reference,
    release_name = excluded.release_name,
    provider_subtitle_id = excluded.provider_subtitle_id,
    checksum = excluded.checksum,
    size_bytes = excluded.size_bytes,
    downloaded_at = excluded.downloaded_at,
    updated_at = now()
returning *;

-- name: GetMediaItemSubtitle :one
select *
from app.media_item_subtitles
where media_item_id = $1 and id = $2;

-- name: ClearSelectedMediaItemSubtitles :exec
update app.media_item_subtitles
set selected = false,
    updated_at = now()
where media_item_id = $1
  and language_id = $2
  and file_path = $3
  and id <> $4;

-- name: UpdateMediaItemSubtitleSelection :one
update app.media_item_subtitles
set selected = $3,
    retention_mode = $4,
    updated_at = now()
where media_item_id = $1 and id = $2
returning *;

-- name: DeleteMediaItemSubtitle :execrows
delete from app.media_item_subtitles
where media_item_id = $1 and id = $2;
