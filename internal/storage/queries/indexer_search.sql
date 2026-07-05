-- name: EnsureIndexerSearchSettings :exec
insert into app.indexer_search_settings (
    id,
    cache_duration_minutes,
    history_retention_days,
    automatic_blocklist_expiry_days
)
values (true, 1440, 7, 7)
on conflict (id) do nothing;

-- name: GetIndexerSearchSettings :one
select cache_duration_minutes,
    history_retention_days,
    automatic_blocklist_expiry_days
from app.indexer_search_settings
where id = true;

-- name: SaveIndexerSearchSettings :exec
insert into app.indexer_search_settings (
    id,
    cache_duration_minutes,
    history_retention_days,
    automatic_blocklist_expiry_days
)
values (
    true,
    sqlc.arg(cache_duration_minutes),
    sqlc.arg(history_retention_days),
    sqlc.arg(automatic_blocklist_expiry_days)
)
on conflict (id) do update
set cache_duration_minutes = excluded.cache_duration_minutes,
    history_retention_days = excluded.history_retention_days,
    automatic_blocklist_expiry_days = excluded.automatic_blocklist_expiry_days,
    updated_at = now();

-- name: GetIndexerSearchCacheResponse :one
select response
from app.indexer_search_cache
where indexer_id = sqlc.arg(indexer_id)
    and media_type = sqlc.arg(media_type)
    and query = sqlc.arg(query)
    and expires_at > now();

-- name: SetIndexerSearchCache :exec
insert into app.indexer_search_cache (
    indexer_id,
    media_type,
    query,
    response,
    result_count,
    expires_at
)
values (
    sqlc.arg(indexer_id),
    sqlc.arg(media_type),
    sqlc.arg(query),
    sqlc.arg(response),
    sqlc.arg(result_count),
    sqlc.arg(expires_at)
)
on conflict (indexer_id, media_type, query) do update
set response = excluded.response,
    result_count = excluded.result_count,
    expires_at = excluded.expires_at,
    updated_at = now();

-- name: RecordIndexerSearchHistory :one
insert into app.indexer_search_history (
    id,
    indexer_id,
    indexer_name,
    indexer_protocol,
    media_type,
    query,
    cache_hit,
    success,
    result_count,
    error,
    response
)
values (
    sqlc.arg(id),
    sqlc.arg(indexer_id),
    sqlc.arg(indexer_name),
    sqlc.arg(indexer_protocol),
    sqlc.arg(media_type),
    sqlc.arg(query),
    sqlc.arg(cache_hit),
    sqlc.arg(success),
    sqlc.arg(result_count),
    sqlc.narg(error),
    sqlc.arg(response)
)
returning indexer_name,
    indexer_protocol,
    media_type,
    query,
    cache_hit,
    success,
    result_count,
    error,
    response::text as response,
    created_at;

-- name: CleanupIndexerSearchHistory :execrows
delete from app.indexer_search_history
where created_at < now() - make_interval(days => sqlc.arg(retention_days)::int);

-- name: ClearIndexerSearchHistory :execrows
delete from app.indexer_search_history;

-- name: ClearIndexerSearchCache :execrows
delete from app.indexer_search_cache;

-- name: ClearIndexerSearchCacheByPattern :execrows
delete from app.indexer_search_cache
where query ~* sqlc.arg(pattern);

-- name: DeleteIndexerSearchCacheEntry :execrows
delete from app.indexer_search_cache
where indexer_id = sqlc.arg(indexer_id)
    and media_type = sqlc.arg(media_type)
    and query = sqlc.arg(query);

-- name: IndexerSearchCacheStats :one
select count(*)::int as total_entries,
    count(*) filter (where expires_at > now())::int as active_entries,
    count(*) filter (where expires_at <= now())::int as expired_entries,
    count(distinct indexer_id)::int as indexer_count
from app.indexer_search_cache;

-- name: ListIndexerSearchCacheEntries :many
select i.id as indexer_id,
    i.name as indexer_name,
    i.protocol as indexer_protocol,
    c.media_type,
    c.query,
    c.result_count,
    c.expires_at,
    c.created_at,
    c.updated_at,
    (c.expires_at <= now())::bool as expired
from app.indexer_search_cache c
join app.indexers i on i.id = c.indexer_id
order by c.updated_at desc
limit sqlc.arg(row_limit);

-- name: GetIndexerSearchCacheEntry :one
select i.id as indexer_id,
    i.name as indexer_name,
    i.protocol as indexer_protocol,
    c.media_type,
    c.query,
    c.result_count,
    c.expires_at,
    c.created_at,
    c.updated_at,
    (c.expires_at <= now())::bool as expired
from app.indexer_search_cache c
join app.indexers i on i.id = c.indexer_id
where c.indexer_id = sqlc.arg(indexer_id)
    and c.media_type = sqlc.arg(media_type)
    and c.query = sqlc.arg(query);

-- name: ListIndexerSearchHistoryEntries :many
select indexer_name,
    indexer_protocol,
    media_type,
    query,
    cache_hit,
    success,
    result_count,
    error,
    response::text as response,
    created_at
from app.indexer_search_history
order by created_at desc
limit sqlc.arg(row_limit);

-- name: IndexerSearchHistoryCount :one
select count(*)::int
from app.indexer_search_history;

-- name: IndexerSearchHistoryStats :one
select count(*)::int as total_entries,
    count(*) filter (where cache_hit)::int as cache_hits,
    count(*) filter (where not cache_hit)::int as cache_misses,
    count(*) filter (where not success)::int as failures
from app.indexer_search_history;
