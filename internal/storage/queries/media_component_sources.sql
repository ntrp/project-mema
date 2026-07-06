-- name: ListMediaComponentSources :many
select *
from app.media_component_sources
where media_item_id = $1
order by retained_at desc, source_role;

-- name: GetMediaComponentSource :one
select *
from app.media_component_sources
where media_item_id = $1 and id = $2;

-- name: CreateMediaComponentSource :one
insert into app.media_component_sources (
    id,
    media_item_id,
    source_role,
    source_file_path,
    retained_path,
    release_title,
    release_group,
    release_name,
    release_id,
    source_metadata,
    stream_inventory,
    checksum,
    size_bytes
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(source_role),
    sqlc.arg(source_file_path),
    sqlc.arg(retained_path),
    sqlc.narg(release_title),
    sqlc.narg(release_group),
    sqlc.narg(release_name),
    sqlc.narg(release_id),
    sqlc.narg(source_metadata),
    sqlc.arg(stream_inventory),
    sqlc.narg(checksum),
    sqlc.narg(size_bytes)
)
returning *;

-- name: ReleaseMediaComponentSource :one
update app.media_component_sources
set retention_state = 'released',
    released_at = now(),
    updated_at = now()
where media_item_id = $1 and id = $2
returning *;
