-- name: ListDownloadClients :many
select id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
from app.download_clients
order by priority asc, name asc;

-- name: ListEnabledDownloadClients :many
select id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
from app.download_clients
where enabled = true
order by priority asc, name asc;

-- name: GetDownloadClient :one
select id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
from app.download_clients
where id = $1;

-- name: CreateDownloadClient :one
insert into app.download_clients (
    id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority
)
values (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(type),
    sqlc.arg(protocol),
    sqlc.arg(base_url),
    sqlc.narg(username),
    sqlc.narg(password),
    sqlc.narg(api_key),
    sqlc.narg(category),
    sqlc.arg(enabled),
    sqlc.arg(priority)
)
returning id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at;

-- name: UpdateDownloadClient :one
update app.download_clients
set name = $2,
    type = $3,
    protocol = $4,
    base_url = $5,
    username = $6,
    password = $7,
    api_key = $8,
    category = $9,
    enabled = $10,
    priority = $11,
    updated_at = now()
where id = $1
returning id, name, type, protocol, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at;

-- name: DeleteDownloadClient :execrows
delete from app.download_clients
where id = $1;
