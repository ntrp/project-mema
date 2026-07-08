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
insert into app.system_job_schedules (
    id,
    name,
    category,
    description,
    kind,
    queue,
    interval_seconds,
    interval_configurable,
    history_policy,
    automatic,
    manual_action_available
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
on conflict (id) do update set
    name = excluded.name,
    category = excluded.category,
    description = excluded.description,
    kind = excluded.kind,
    queue = excluded.queue,
    interval_seconds = case
        when excluded.interval_configurable then app.system_job_schedules.interval_seconds
        else excluded.interval_seconds
    end,
    interval_configurable = excluded.interval_configurable,
    history_policy = excluded.history_policy,
    automatic = excluded.automatic,
    manual_action_available = excluded.manual_action_available,
    updated_at = now()
returning *;

-- name: ListSystemJobSchedules :many
select s.id,
    s.name,
    s.category,
    s.description,
    s.kind,
    s.queue,
    s.interval_seconds,
    s.interval_configurable,
    s.history_policy,
    s.automatic,
    s.manual_action_available,
    s.paused,
    s.created_at,
    s.updated_at,
    coalesce(active.river_job_id, 0)::bigint as active_river_job_id,
    coalesce(active.status, '')::text as active_status,
    active.progress_percent as active_progress_percent,
    coalesce(active.progress_label, '')::text as active_progress_label,
    coalesce(active.progress_data, '{}'::jsonb) as active_progress_data,
    coalesce(active.info_message, '')::text as active_info_message,
    coalesce(last_run.river_job_id, 0)::bigint as last_river_job_id,
    coalesce(last_run.status, '')::text as last_status,
    coalesce(last_run.created_at, 'epoch'::timestamptz) as last_created_at,
    last_run.finalized_at as last_finalized_at
from app.system_job_schedules s
left join lateral (
    select river_job_id, status, progress_percent, progress_label, progress_data, info_message
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

-- name: SystemJobScheduleReady :one
select exists(
    select 1
    from app.system_job_schedules s
    where s.id = $1
        and not s.paused
        and not exists (
            select 1
            from app.system_job_executions active
            where active.schedule_id = s.id
                and active.status in ('available', 'scheduled', 'retryable', 'running')
        )
        and coalesce(
            (
                select max(last_run.created_at) + s.interval_seconds * interval '1 second' <= now()
                from app.system_job_executions last_run
                where last_run.schedule_id = s.id
            ),
            true
        )
);

-- name: UpdateSystemJobSchedulePaused :one
update app.system_job_schedules
set paused = $2,
    updated_at = now()
where id = $1
returning *;

-- name: UpdateSystemJobScheduleInterval :one
update app.system_job_schedules
set interval_seconds = $2,
    updated_at = now()
where id = $1
    and interval_configurable
returning *;

-- name: UpsertSystemJobExecution :one
insert into app.system_job_executions (
    river_job_id,
    schedule_id,
    classification,
    history_policy,
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
    coalesce((select history_policy from app.system_job_schedules where id = sqlc.narg(schedule_id)), 'standard'),
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
    history_policy = excluded.history_policy,
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
    progress_data = coalesce(sqlc.narg(progress_data), '{}'::jsonb),
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
    and (
        sqlc.arg(include_routine)::bool
        or history_policy <> 'routine'
        or status in ('retryable', 'cancelled', 'discarded')
    )
    and (sqlc.narg(before)::timestamptz is null or updated_at < sqlc.narg(before)::timestamptz)
    and (
        sqlc.arg(search_query)::text = ''
        or kind ilike '%' || sqlc.arg(search_query)::text || '%'
        or queue ilike '%' || sqlc.arg(search_query)::text || '%'
        or info_message ilike '%' || sqlc.arg(search_query)::text || '%'
        or args::text ilike '%' || sqlc.arg(search_query)::text || '%'
        or errors::text ilike '%' || sqlc.arg(search_query)::text || '%'
    )
order by updated_at desc, river_job_id desc
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
    )::int as retention_days,
    coalesce(
        (select routine_retention_hours from app.system_job_history_settings where id),
        sqlc.arg(routine_retention_hours)::int
    )::int as routine_retention_hours;

-- name: UpdateSystemJobHistorySettings :one
insert into app.system_job_history_settings (id, retention_days, routine_retention_hours)
values (true, $1, $2)
on conflict (id) do update set
    retention_days = excluded.retention_days,
    routine_retention_hours = excluded.routine_retention_hours,
    updated_at = now()
returning *;

-- name: PruneSystemJobExecutions :exec
delete from app.system_job_executions
where finalized_at is not null
    and (
        (
            history_policy = 'routine'
            and status = 'completed'
            and finalized_at < now() - make_interval(hours => sqlc.arg(routine_retention_hours)::int)
        )
        or (
            (history_policy <> 'routine' or status <> 'completed')
            and finalized_at < now() - make_interval(days => sqlc.arg(retention_days)::int)
        )
    );
