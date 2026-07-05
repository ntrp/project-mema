-- name: ListQualitySizeSettings :many
select
    quality_id,
    minimum_size_mb_per_minute,
    preferred_size_mb_per_minute,
    maximum_size_mb_per_minute,
    created_at,
    updated_at
from app.quality_size_settings;

-- name: EnsureQualitySizeSetting :exec
insert into app.quality_size_settings (
    quality_id,
    minimum_size_mb_per_minute,
    preferred_size_mb_per_minute,
    maximum_size_mb_per_minute
)
values (
    sqlc.arg(quality_id),
    sqlc.arg(minimum_size_mb_per_minute),
    sqlc.narg(preferred_size_mb_per_minute),
    sqlc.narg(maximum_size_mb_per_minute)
)
on conflict (quality_id) do nothing;

-- name: UpsertQualitySizeSetting :exec
insert into app.quality_size_settings (
    quality_id,
    minimum_size_mb_per_minute,
    preferred_size_mb_per_minute,
    maximum_size_mb_per_minute
)
values (
    sqlc.arg(quality_id),
    sqlc.arg(minimum_size_mb_per_minute),
    sqlc.narg(preferred_size_mb_per_minute),
    sqlc.narg(maximum_size_mb_per_minute)
)
on conflict (quality_id) do update
set minimum_size_mb_per_minute = excluded.minimum_size_mb_per_minute,
    preferred_size_mb_per_minute = excluded.preferred_size_mb_per_minute,
    maximum_size_mb_per_minute = excluded.maximum_size_mb_per_minute,
    updated_at = now();
