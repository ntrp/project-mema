-- name: CreateSystemEvent :one
insert into app.system_events (id, severity, category, message, data)
values ($1, $2, $3, $4, sqlc.arg(data)::jsonb)
returning id, severity, category, message, data, created_at;

-- name: ListSystemEvents :many
select id, severity, category, message, data, created_at
from app.system_events
where (sqlc.narg(before)::timestamptz is null or created_at < sqlc.narg(before)::timestamptz)
order by created_at desc
limit $1;

-- name: DeleteSystemEvent :execrows
delete from app.system_events
where id = $1;

-- name: ClearSystemEvents :exec
delete from app.system_events;

-- name: GetSystemEventSettings :one
insert into app.system_event_settings (id, retention_days)
values (true, $1)
on conflict (id) do update set id = excluded.id
returning retention_days, created_at, updated_at;

-- name: UpdateSystemEventSettings :one
insert into app.system_event_settings (id, retention_days)
values (true, $1)
on conflict (id) do update
set retention_days = excluded.retention_days,
    updated_at = now()
returning retention_days, created_at, updated_at;

-- name: PruneSystemEvents :exec
delete from app.system_events
where created_at < now() - (sqlc.arg(retention_days)::int * interval '1 day');
