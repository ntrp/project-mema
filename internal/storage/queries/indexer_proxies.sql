-- name: ListIndexerProxies :many
select
    id,
    name,
    implementation,
    link,
    enabled,
    on_health_issue,
    supports_on_health_issue,
    include_health_warnings,
    test_command,
    fields,
    created_at,
    updated_at
from app.indexer_proxies
order by name asc;

-- name: CreateIndexerProxy :one
insert into app.indexer_proxies (
    id,
    name,
    implementation,
    link,
    enabled,
    on_health_issue,
    supports_on_health_issue,
    include_health_warnings,
    test_command,
    fields
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning
    id,
    name,
    implementation,
    link,
    enabled,
    on_health_issue,
    supports_on_health_issue,
    include_health_warnings,
    test_command,
    fields,
    created_at,
    updated_at;

-- name: UpdateIndexerProxy :one
update app.indexer_proxies
set name = $2,
    implementation = $3,
    link = $4,
    enabled = $5,
    on_health_issue = $6,
    supports_on_health_issue = $7,
    include_health_warnings = $8,
    test_command = $9,
    fields = $10,
    updated_at = now()
where id = $1
returning
    id,
    name,
    implementation,
    link,
    enabled,
    on_health_issue,
    supports_on_health_issue,
    include_health_warnings,
    test_command,
    fields,
    created_at,
    updated_at;

-- name: GetIndexerProxy :one
select
    id,
    name,
    implementation,
    link,
    enabled,
    on_health_issue,
    supports_on_health_issue,
    include_health_warnings,
    test_command,
    fields,
    created_at,
    updated_at
from app.indexer_proxies
where id = $1;

-- name: DeleteIndexerProxy :execrows
delete from app.indexer_proxies
where id = $1;
