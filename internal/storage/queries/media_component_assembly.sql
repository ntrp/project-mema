-- name: ListMediaComponentAssemblyRuns :many
select id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at
from app.media_component_assembly_runs
where media_item_id = sqlc.arg(media_item_id)
order by created_at desc, id;

-- name: ListMediaComponentAssemblyInputs :many
select id, run_id, source_id, artifact_id, stream_type, input_path, provenance, created_at
from app.media_component_assembly_inputs
where run_id = sqlc.arg(run_id)
order by created_at, id;

-- name: GetMediaComponentAssemblyRun :one
select id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at
from app.media_component_assembly_runs
where id = sqlc.arg(id);

-- name: CreateMediaComponentAssemblyRun :one
insert into app.media_component_assembly_runs (
    id,
    media_item_id,
    base_source_id,
    output_path,
    job_id
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(base_source_id),
    sqlc.arg(output_path),
    sqlc.arg(job_id)
)
returning id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: CreateMediaComponentAssemblyInput :one
insert into app.media_component_assembly_inputs (
    id,
    run_id,
    source_id,
    artifact_id,
    stream_type,
    input_path,
    provenance
) values (
    sqlc.arg(id),
    sqlc.arg(run_id),
    sqlc.arg(source_id),
    sqlc.arg(artifact_id),
    sqlc.arg(stream_type),
    sqlc.arg(input_path),
    sqlc.arg(provenance)
)
returning id, run_id, source_id, artifact_id, stream_type, input_path, provenance, created_at;

-- name: AssignMediaComponentAssemblyJob :one
update app.media_component_assembly_runs
set job_id = sqlc.arg(job_id),
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: StartMediaComponentAssemblyRun :one
update app.media_component_assembly_runs
set status = 'running',
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: CompleteMediaComponentAssemblyRun :one
update app.media_component_assembly_runs
set status = 'succeeded',
    tool_summary = sqlc.arg(tool_summary),
    error_message = null,
    size_bytes = sqlc.arg(size_bytes),
    updated_at = now(),
    completed_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;

-- name: FailMediaComponentAssemblyRun :one
update app.media_component_assembly_runs
set status = 'failed',
    tool_summary = sqlc.arg(tool_summary),
    error_message = sqlc.arg(error_message),
    updated_at = now(),
    completed_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, base_source_id, output_path, status, tool_name, tool_summary, error_message, job_id, size_bytes, created_at, updated_at, completed_at;
