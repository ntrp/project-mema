-- name: ListMediaItemSidecars :many
select *
from app.media_item_sidecars
where media_item_id = $1
order by media_file_path, sidecar_type, file_path;

-- name: UpsertMediaItemSidecar :one
insert into app.media_item_sidecars (
    id,
    media_item_id,
    media_file_path,
    file_path,
    sidecar_type,
    subtype,
    language_id,
    format
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(media_file_path),
    sqlc.arg(file_path),
    sqlc.arg(sidecar_type),
    sqlc.narg(subtype),
    sqlc.narg(language_id),
    sqlc.narg(format)
)
on conflict (media_item_id, file_path) do update
set media_file_path = excluded.media_file_path,
    sidecar_type = excluded.sidecar_type,
    subtype = excluded.subtype,
    language_id = excluded.language_id,
    format = excluded.format,
    updated_at = now()
returning *;
