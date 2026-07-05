-- name: BulkUpdateIndexer :exec
update app.indexers
set enabled = coalesce(sqlc.narg(enabled)::boolean, enabled),
    app_profile_id = coalesce(sqlc.narg(app_profile_id)::text, app_profile_id),
    priority = coalesce(sqlc.narg(priority)::integer, priority),
    minimum_seeders = coalesce(sqlc.narg(minimum_seeders)::integer, minimum_seeders),
    seed_ratio = coalesce(sqlc.narg(seed_ratio)::numeric, seed_ratio),
    seed_time = coalesce(sqlc.narg(seed_time)::integer, seed_time),
    pack_seed_time = coalesce(sqlc.narg(pack_seed_time)::integer, pack_seed_time),
    prefer_magnet_url = coalesce(sqlc.narg(prefer_magnet_url)::boolean, prefer_magnet_url),
    updated_at = now()
where id = sqlc.arg(id);
