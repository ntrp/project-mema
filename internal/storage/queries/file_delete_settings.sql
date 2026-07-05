-- name: GetFileDeleteSettings :one
insert into app.file_delete_settings (id, mode, recycle_folder)
values (true, $1, $2)
on conflict (id) do update set id = excluded.id
returning mode, recycle_folder, created_at, updated_at;

-- name: UpdateFileDeleteSettings :one
insert into app.file_delete_settings (id, mode, recycle_folder)
values (true, $1, $2)
on conflict (id) do update
set mode = excluded.mode,
    recycle_folder = excluded.recycle_folder,
    updated_at = now()
returning mode, recycle_folder, created_at, updated_at;
