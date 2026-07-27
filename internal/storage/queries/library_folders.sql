-- name: ListLibraryFolders :many
select id, path, kind, manual, created_at, updated_at
from app.library_folders
order by path asc;

-- name: GetLibraryFolder :one
select id, path, kind, manual, created_at, updated_at
from app.library_folders
where id = $1;

-- name: UpsertLibraryFolder :one
insert into app.library_folders (id, path, kind, manual)
values ($1, $2, $3, $4)
on conflict (path) do update set kind = excluded.kind, manual = excluded.manual, updated_at = now()
returning id, path, kind, manual, created_at, updated_at;

-- name: UpdateLibraryFolderManual :one
update app.library_folders
set manual = $2, updated_at = now()
where id = $1
returning id, path, kind, manual, created_at, updated_at;

-- name: DeleteLibraryFolder :execrows
delete from app.library_folders
where id = $1;

-- name: LibraryFolderExists :one
select exists(select 1 from app.library_folders where id = $1);
