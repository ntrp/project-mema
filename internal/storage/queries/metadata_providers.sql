-- name: EnsureMetadataProvider :exec
insert into app.metadata_providers (
    id, name, type, base_url, enabled, priority
)
select $1, $2, $3, $4, $5, $6
where not exists (
    select 1 from app.metadata_providers where type = $3
);

-- name: ListMetadataProviders :many
select id, name, type, base_url, api_key, pin, access_token, session_token,
    session_token_expires_at, enabled, priority, created_at, updated_at
from app.metadata_providers
order by priority asc, name asc;

-- name: ListEnabledMetadataProviders :many
select id, name, type, base_url, api_key, pin, access_token, session_token,
    session_token_expires_at, enabled, priority, created_at, updated_at
from app.metadata_providers
where enabled = true
    and ((sqlc.arg(media_type) = 'movie' and type in ('tmdb', 'tvdb'))
        or (sqlc.arg(media_type) = 'serie' and type in ('tmdb', 'tvdb')))
order by priority asc, name asc;

-- name: GetMetadataProvider :one
select id, name, type, base_url, api_key, pin, access_token, session_token,
    session_token_expires_at, enabled, priority, created_at, updated_at
from app.metadata_providers
where id = $1;

-- name: CreateMetadataProvider :one
insert into app.metadata_providers (
    id, name, type, base_url, api_key, pin, access_token, enabled, priority
)
values (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(type),
    sqlc.arg(base_url),
    sqlc.narg(api_key),
    sqlc.narg(pin),
    sqlc.narg(access_token),
    sqlc.arg(enabled),
    sqlc.arg(priority)
)
returning id, name, type, base_url, api_key, pin, access_token, session_token,
    session_token_expires_at, enabled, priority, created_at, updated_at;

-- name: UpdateMetadataProvider :one
update app.metadata_providers
set name = $2,
    type = $3,
    base_url = $4,
    api_key = $5,
    pin = $6,
    access_token = $7,
    session_token = null,
    session_token_expires_at = null,
    enabled = $8,
    priority = $9,
    updated_at = now()
where id = $1
returning id, name, type, base_url, api_key, pin, access_token, session_token,
    session_token_expires_at, enabled, priority, created_at, updated_at;

-- name: DeleteMetadataProvider :execrows
delete from app.metadata_providers
where id = $1;

-- name: UpdateMetadataProviderSessionToken :exec
update app.metadata_providers
set session_token = $2,
    session_token_expires_at = $3,
    updated_at = now()
where id = $1;
