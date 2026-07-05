-- name: ListDiscoverBlacklist :many
select id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at
from app.discover_blacklist
order by created_at desc, lower(title);

-- name: SaveDiscoverBlacklistByExternalID :one
insert into app.discover_blacklist (
    id, media_type, title, year, external_provider, external_id, overview, poster_path
)
values (
    sqlc.arg(id),
    sqlc.arg(media_type),
    sqlc.arg(title),
    sqlc.narg(year),
    sqlc.narg(external_provider),
    sqlc.narg(external_id),
    sqlc.narg(overview),
    sqlc.narg(poster_path)
)
on conflict (media_type, external_provider, external_id)
    where external_provider is not null and external_id is not null
do update
set title = excluded.title,
    year = excluded.year,
    overview = excluded.overview,
    poster_path = excluded.poster_path
returning id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at;

-- name: SaveDiscoverBlacklistByTitle :one
insert into app.discover_blacklist (
    id, media_type, title, year, external_provider, external_id, overview, poster_path
)
values (
    sqlc.arg(id),
    sqlc.arg(media_type),
    sqlc.arg(title),
    sqlc.narg(year),
    null,
    null,
    sqlc.narg(overview),
    sqlc.narg(poster_path)
)
on conflict (media_type, lower(title), coalesce(year, 0))
    where external_provider is null or external_id is null
do update
set overview = excluded.overview,
    poster_path = excluded.poster_path
returning id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at;

-- name: DeleteDiscoverBlacklistItem :execrows
delete from app.discover_blacklist
where id = $1;
