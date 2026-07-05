-- name: CreateMediaFileHistory :one
insert into app.media_file_history (
    id,
    media_item_id,
    file_path,
    source_path,
    destination_path,
    operation,
    status,
    actor_type,
    actor_id,
    job_id,
    details,
    failure_details
)
values (
    sqlc.arg(id),
    sqlc.narg(media_item_id),
    sqlc.arg(file_path),
    sqlc.narg(source_path),
    sqlc.narg(destination_path),
    sqlc.arg(operation),
    sqlc.arg(status),
    sqlc.arg(actor_type),
    sqlc.narg(actor_id),
    sqlc.narg(job_id),
    sqlc.arg(details),
    sqlc.narg(failure_details)
)
returning *;

-- name: ListMediaFileHistory :many
select *
from app.media_file_history
where media_item_id = $1
order by created_at desc, id desc;
