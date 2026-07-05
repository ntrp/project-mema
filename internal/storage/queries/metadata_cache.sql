-- name: GetMetadataSearchCacheResults :one
select results
from app.metadata_search_cache
where provider_id = sqlc.arg(provider_id)
    and media_type = sqlc.arg(media_type)
    and query = sqlc.arg(query)
    and year = sqlc.arg(year)
    and expires_at > now();

-- name: SetMetadataSearchCache :exec
insert into app.metadata_search_cache (
    provider_id,
    media_type,
    query,
    year,
    results,
    expires_at
)
values (
    sqlc.arg(provider_id),
    sqlc.arg(media_type),
    sqlc.arg(query),
    sqlc.arg(year),
    sqlc.arg(results),
    sqlc.arg(expires_at)
)
on conflict (provider_id, media_type, query, year) do update
set results = excluded.results,
    expires_at = excluded.expires_at,
    updated_at = now();

-- name: RecordMetadataSearchHistory :one
insert into app.metadata_search_history (
    id,
    provider_id,
    provider_name,
    provider_type,
    media_type,
    query,
    year,
    cache_hit,
    success,
    item_count,
    error,
    response
)
values (
    sqlc.arg(id),
    sqlc.arg(provider_id),
    sqlc.arg(provider_name),
    sqlc.arg(provider_type),
    sqlc.arg(media_type),
    sqlc.arg(query),
    sqlc.arg(year),
    sqlc.arg(cache_hit),
    sqlc.arg(success),
    sqlc.arg(item_count),
    sqlc.narg(error),
    sqlc.arg(response)
)
returning provider_name,
    provider_type,
    media_type,
    query,
    year,
    cache_hit,
    success,
    item_count,
    error,
    response::text as response,
    created_at;

-- name: MetadataCacheStats :one
select count(*)::int as total_entries,
    count(*) filter (where expires_at > now())::int as active_entries,
    count(*) filter (where expires_at <= now())::int as expired_entries,
    count(distinct provider_id)::int as provider_count
from app.metadata_search_cache;

-- name: ListMetadataSearchHistoryEntries :many
select provider_name,
    provider_type,
    media_type,
    query,
    year,
    cache_hit,
    success,
    item_count,
    error,
    response::text as response,
    created_at
from app.metadata_search_history
order by created_at desc
limit sqlc.arg(row_limit);

-- name: ListMetadataCacheEntries :many
select p.id as provider_id,
    p.name as provider_name,
    p.type as provider_type,
    c.media_type,
    c.query,
    c.year,
    case
        when jsonb_typeof(c.results) = 'array' then jsonb_array_length(c.results)
        else 1
    end::int as item_count,
    c.expires_at,
    c.created_at,
    c.updated_at,
    (c.expires_at <= now())::bool as expired
from app.metadata_search_cache c
join app.metadata_providers p on p.id = c.provider_id
order by c.updated_at desc
limit sqlc.arg(row_limit);

-- name: MetadataSearchHistoryCount :one
select count(*)::int
from app.metadata_search_history;

-- name: MetadataSearchHistoryStats :one
select count(*)::int as total_entries,
    count(*) filter (where cache_hit)::int as cache_hits,
    count(*) filter (where not cache_hit)::int as cache_misses,
    count(*) filter (where not success)::int as failures
from app.metadata_search_history;

-- name: ClearMetadataCache :execrows
delete from app.metadata_search_cache;

-- name: ClearMetadataCacheByPattern :execrows
delete from app.metadata_search_cache
where query ~* sqlc.arg(pattern);

-- name: DeleteMetadataCacheEntry :execrows
delete from app.metadata_search_cache
where provider_id = sqlc.arg(provider_id)
    and media_type = sqlc.arg(media_type)
    and query = sqlc.arg(query)
    and year = sqlc.arg(year);

-- name: ClearMetadataSearchHistory :execrows
delete from app.metadata_search_history;
