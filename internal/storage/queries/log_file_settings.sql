-- name: GetLogFileSettings :one
insert into app.log_file_settings (id, enabled, directory, retention_days)
values (true, false, $1, $2)
on conflict (id) do update set id = excluded.id
returning enabled, directory, retention_days, created_at, updated_at;

-- name: UpdateLogFileSettings :one
insert into app.log_file_settings (id, enabled, directory, retention_days)
values (true, $1, $2, $3)
on conflict (id) do update
set enabled = excluded.enabled,
    directory = excluded.directory,
    retention_days = excluded.retention_days,
    updated_at = now()
returning enabled, directory, retention_days, created_at, updated_at;
