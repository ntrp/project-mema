-- name: ListPathMappings :many
select id, client_path, app_path, created_at, updated_at
from app.path_mappings
order by client_path asc;

-- name: UpsertPathMapping :one
insert into app.path_mappings (id, client_path, app_path)
values ($1, $2, $3)
on conflict (client_path) do update
set app_path = excluded.app_path, updated_at = now()
returning id, client_path, app_path, created_at, updated_at;

-- name: DeletePathMapping :execrows
delete from app.path_mappings
where id = $1;
