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
    file_path,
    source_url,
    release_name
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id),
    sqlc.narg(provider_id),
    sqlc.arg(provider_name),
    sqlc.arg(language_id),
    sqlc.arg(file_path),
    sqlc.narg(source_url),
    sqlc.narg(release_name)
)
on conflict (media_item_id, language_id, file_path) do update
set provider_id = excluded.provider_id,
    provider_name = excluded.provider_name,
    source_url = excluded.source_url,
    release_name = excluded.release_name,
    updated_at = now()
returning *;
