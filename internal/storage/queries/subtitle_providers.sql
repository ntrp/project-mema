-- name: ListSubtitleProviders :many
select *
from app.subtitle_providers
order by priority asc, lower(name) asc;

-- name: GetSubtitleProvider :one
select *
from app.subtitle_providers
where id = $1;

-- name: CreateSubtitleProvider :one
insert into app.subtitle_providers (
    id,
    name,
    type,
    base_url,
    username,
    password,
    api_key,
    enabled,
    priority
) values (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(type),
    sqlc.arg(base_url),
    sqlc.narg(username),
    sqlc.narg(password),
    sqlc.narg(api_key),
    sqlc.arg(enabled),
    sqlc.arg(priority)
)
returning *;

-- name: UpdateSubtitleProvider :one
update app.subtitle_providers
set name = sqlc.arg(name),
    type = sqlc.arg(type),
    base_url = sqlc.arg(base_url),
    username = sqlc.narg(username),
    password = sqlc.narg(password),
    api_key = sqlc.narg(api_key),
    enabled = sqlc.arg(enabled),
    priority = sqlc.arg(priority),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: DeleteSubtitleProvider :execrows
delete from app.subtitle_providers
where id = $1;

-- name: ListMockSubtitleProviderRows :many
select *
from app.mock_subtitle_provider_rows
where provider_id = $1
order by sort_order asc, lower(title) asc, language_id asc, format asc;

-- name: CreateMockSubtitleProviderRow :one
insert into app.mock_subtitle_provider_rows (
    id,
    provider_id,
    title,
    language_id,
    format,
    sort_order
) values (
    sqlc.arg(id),
    sqlc.arg(provider_id),
    sqlc.arg(title),
    sqlc.arg(language_id),
    sqlc.arg(format),
    sqlc.arg(sort_order)
)
returning *;

-- name: DeleteMockSubtitleProviderRows :exec
delete from app.mock_subtitle_provider_rows
where provider_id = $1;
