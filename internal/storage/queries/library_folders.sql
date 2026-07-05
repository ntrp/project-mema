-- name: ListLibraryFolders :many
select id, path, kind, created_at, updated_at
from app.library_folders
order by path asc;

-- name: GetLibraryFolder :one
select id, path, kind, created_at, updated_at
from app.library_folders
where id = $1;

-- name: UpsertLibraryFolder :one
insert into app.library_folders (id, path, kind)
values ($1, $2, $3)
on conflict (path) do update set kind = excluded.kind, updated_at = now()
returning id, path, kind, created_at, updated_at;

-- name: DeleteLibraryFolder :execrows
delete from app.library_folders
where id = $1;

-- name: LibraryFolderExists :one
select exists(select 1 from app.library_folders where id = $1);
