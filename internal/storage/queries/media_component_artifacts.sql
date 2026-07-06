-- name: ListMediaComponentArtifactsForSource :many
select id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at
from app.media_component_artifacts
where source_id = sqlc.arg(source_id)
order by created_at desc, id;

-- name: GetMediaComponentArtifact :one
select id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at
from app.media_component_artifacts
where id = sqlc.arg(id);

-- name: CreateMediaComponentArtifact :one
insert into app.media_component_artifacts (
    id,
    media_item_id,
    source_id,
    stream_id,
    stream_type,
    language,
    output_path,
    job_id
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(source_id),
    sqlc.arg(stream_id),
    sqlc.arg(stream_type),
    sqlc.arg(language),
    sqlc.arg(output_path),
    sqlc.arg(job_id)
)
returning id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: AssignMediaComponentArtifactJob :one
update app.media_component_artifacts
set job_id = sqlc.arg(job_id),
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: StartMediaComponentArtifact :one
update app.media_component_artifacts
set status = 'running',
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: CompleteMediaComponentArtifact :one
update app.media_component_artifacts
set status = 'succeeded',
    tool_summary = sqlc.arg(tool_summary),
    error_message = null,
    size_bytes = sqlc.arg(size_bytes),
    updated_at = now(),
    completed_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: FailMediaComponentArtifact :one
update app.media_component_artifacts
set status = 'failed',
    tool_summary = sqlc.arg(tool_summary),
    error_message = sqlc.arg(error_message),
    updated_at = now(),
    completed_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, source_id, stream_id, stream_type, language, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;
