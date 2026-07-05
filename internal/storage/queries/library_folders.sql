-- name: ListLibraryFolders :many
select id, path, created_at, updated_at
from app.library_folders
order by path asc;

-- name: GetLibraryFolder :one
select id, path, created_at, updated_at
from app.library_folders
where id = $1;

-- name: UpsertLibraryFolder :one
insert into app.library_folders (id, path)
values ($1, $2)
on conflict (path) do update set updated_at = now()
returning id, path, created_at, updated_at;

-- name: DeleteLibraryFolder :execrows
delete from app.library_folders
where id = $1;

-- name: LibraryFolderExists :one
select exists(select 1 from app.library_folders where id = $1);
