-- name: ListSystemJobs :many
select id,
    state::text as state,
    kind,
    queue,
    attempt::int as attempt,
    max_attempts::int as max_attempts,
    priority::int as priority,
    args::text as args,
    metadata::text as metadata,
    coalesce(array_to_json(errors), '[]'::json)::text as errors,
    coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text)::text as info_message,
    scheduled_at,
    created_at,
    attempted_at,
    finalized_at
from river_job
where (cardinality(sqlc.arg(states)::text[]) = 0 or state::text = any(sqlc.arg(states)::text[]))
    and (sqlc.arg(queue)::text = '' or queue = sqlc.arg(queue)::text)
    and (sqlc.arg(kind)::text = '' or kind = sqlc.arg(kind)::text)
    and (
        sqlc.arg(search_query)::text = ''
        or kind ilike '%' || sqlc.arg(search_query)::text || '%'
        or queue ilike '%' || sqlc.arg(search_query)::text || '%'
        or args::text ilike '%' || sqlc.arg(search_query)::text || '%'
        or errors::text ilike '%' || sqlc.arg(search_query)::text || '%'
    )
order by coalesce(finalized_at, attempted_at, scheduled_at, created_at) desc, id desc
limit sqlc.arg(row_limit);

-- name: GetSystemJob :one
select id,
    state::text as state,
    kind,
    queue,
    attempt::int as attempt,
    max_attempts::int as max_attempts,
    priority::int as priority,
    args::text as args,
    metadata::text as metadata,
    coalesce(array_to_json(errors), '[]'::json)::text as errors,
    coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text)::text as info_message,
    scheduled_at,
    created_at,
    attempted_at,
    finalized_at
from river_job
where id = $1;

-- name: UpsertSystemJobSchedule :one
insert into app.system_job_schedules (id, name, kind, queue, interval_seconds)
values ($1, $2, $3, $4, $5)
on conflict (id) do update set
    name = excluded.name,
    kind = excluded.kind,
    queue = excluded.queue,
    interval_seconds = excluded.interval_seconds,
    updated_at = now()
returning *;

-- name: ListSystemJobSchedules :many
select s.id,
    s.name,
    s.kind,
    s.queue,
    s.interval_seconds,
    s.paused,
    s.created_at,
    s.updated_at,
    coalesce(active.river_job_id, 0)::bigint as active_river_job_id,
    coalesce(active.status, '')::text as active_status,
    active.progress_percent as active_progress_percent,
    coalesce(active.progress_label, '')::text as active_progress_label,
    coalesce(active.info_message, '')::text as active_info_message,
    coalesce(last_run.river_job_id, 0)::bigint as last_river_job_id,
    coalesce(last_run.status, '')::text as last_status,
    coalesce(last_run.created_at, 'epoch'::timestamptz) as last_created_at,
    last_run.finalized_at as last_finalized_at
from app.system_job_schedules s
left join lateral (
    select river_job_id, status, progress_percent, progress_label, info_message
    from app.system_job_executions
    where schedule_id = s.id
        and status in ('available', 'scheduled', 'retryable', 'running')
    order by updated_at desc, river_job_id desc
    limit 1
) active on true
left join lateral (
    select river_job_id, status, created_at, finalized_at
    from app.system_job_executions
    where schedule_id = s.id
    order by coalesce(finalized_at, updated_at, created_at) desc, river_job_id desc
    limit 1
) last_run on true
order by s.name;

-- name: GetSystemJobSchedule :one
select *
from app.system_job_schedules
where id = $1;

-- name: UpdateSystemJobSchedulePaused :one
update app.system_job_schedules
set paused = $2,
    updated_at = now()
where id = $1
returning *;

-- name: UpsertSystemJobExecution :one
insert into app.system_job_executions (
    river_job_id,
    schedule_id,
    classification,
    status,
    kind,
    queue,
    attempt,
    max_attempts,
    priority,
    args,
    metadata,
    errors,
    info_message,
    scheduled_at,
    created_at,
    attempted_at,
    finalized_at
) values (
    $1,
    sqlc.narg(schedule_id),
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    sqlc.narg(attempted_at),
    sqlc.narg(finalized_at)
)
on conflict (river_job_id) do update set
    schedule_id = excluded.schedule_id,
    classification = excluded.classification,
    status = excluded.status,
    kind = excluded.kind,
    queue = excluded.queue,
    attempt = excluded.attempt,
    max_attempts = excluded.max_attempts,
    priority = excluded.priority,
    args = excluded.args,
    metadata = excluded.metadata,
    errors = excluded.errors,
    info_message = excluded.info_message,
    scheduled_at = excluded.scheduled_at,
    created_at = excluded.created_at,
    attempted_at = excluded.attempted_at,
    finalized_at = excluded.finalized_at,
    updated_at = now()
returning *;

-- name: UpdateSystemJobExecutionProgress :one
update app.system_job_executions
set progress_percent = sqlc.narg(progress_percent),
    progress_label = $2,
    info_message = case when $2 = '' then info_message else $2 end,
    updated_at = now()
where river_job_id = $1
returning *;

-- name: ListCurrentOneShotJobExecutions :many
select *
from app.system_job_executions
where classification = 'one_shot'
    and status in ('available', 'scheduled', 'retryable', 'running')
order by coalesce(attempted_at, scheduled_at, created_at) desc, river_job_id desc
limit sqlc.arg(row_limit);

-- name: ListSystemJobExecutions :many
select *
from app.system_job_executions
where (cardinality(sqlc.arg(states)::text[]) = 0 or status = any(sqlc.arg(states)::text[]))
    and (sqlc.arg(schedule_id)::text = '' or coalesce(schedule_id, '') = sqlc.arg(schedule_id)::text)
    and (sqlc.arg(kind)::text = '' or kind = sqlc.arg(kind)::text)
    and (sqlc.arg(queue)::text = '' or queue = sqlc.arg(queue)::text)
    and (sqlc.narg(before)::timestamptz is null or coalesce(finalized_at, updated_at, created_at) < sqlc.narg(before)::timestamptz)
    and (
        sqlc.arg(search_query)::text = ''
        or kind ilike '%' || sqlc.arg(search_query)::text || '%'
        or queue ilike '%' || sqlc.arg(search_query)::text || '%'
        or info_message ilike '%' || sqlc.arg(search_query)::text || '%'
        or args::text ilike '%' || sqlc.arg(search_query)::text || '%'
        or errors::text ilike '%' || sqlc.arg(search_query)::text || '%'
    )
order by coalesce(finalized_at, updated_at, created_at) desc, river_job_id desc
limit sqlc.arg(row_limit);

-- name: GetSystemJobExecution :one
select *
from app.system_job_executions
where river_job_id = $1;

-- name: CreateSystemJobExecutionLog :one
insert into app.system_job_execution_logs (river_job_id, severity, message, data)
values ($1, $2, $3, $4)
returning *;

-- name: ListSystemJobExecutionLogs :many
select *
from app.system_job_execution_logs
where river_job_id = $1
order by created_at, id
limit sqlc.arg(row_limit);

-- name: GetSystemJobHistorySettings :one
select coalesce(
    (select retention_days from app.system_job_history_settings where id),
    sqlc.arg(retention_days)::int
)::int as retention_days;

-- name: UpdateSystemJobHistorySettings :one
insert into app.system_job_history_settings (id, retention_days)
values (true, $1)
on conflict (id) do update set
    retention_days = excluded.retention_days,
    updated_at = now()
returning *;

-- name: PruneSystemJobExecutions :exec
delete from app.system_job_executions
where finalized_at is not null
    and finalized_at < now() - make_interval(days => sqlc.arg(retention_days)::int);
