-- name: CreateImportAttempt :one
insert into app.import_attempts (
    id,
    activity_id,
    media_item_id,
    source_path,
    target_path,
    import_mode,
    status,
    failure_stage,
    error_message,
    created_targets,
    inserted_media_file_paths
)
values (
    sqlc.arg(id),
    sqlc.arg(activity_id),
    sqlc.arg(media_item_id),
    sqlc.narg(source_path),
    sqlc.narg(target_path),
    sqlc.arg(import_mode),
    sqlc.arg(status),
    sqlc.narg(failure_stage),
    sqlc.narg(error_message),
    sqlc.arg(created_targets),
    sqlc.arg(inserted_media_file_paths)
)
returning *;

-- name: ListImportAttemptsForActivity :many
select *
from app.import_attempts
where activity_id = $1
order by created_at desc;
