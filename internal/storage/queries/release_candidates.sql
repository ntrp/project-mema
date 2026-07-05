-- name: ClearReleaseCandidatesForMedia :exec
delete from app.media_release_candidates
where media_item_id = $1;

-- name: ClearReleaseSearchErrorsForMedia :exec
delete from app.media_release_search_errors
where media_item_id = $1;

-- name: AddReleaseCandidate :exec
insert into app.media_release_candidates (
    id,
    media_item_id,
    season_id,
    episode_id,
    indexer_id,
    indexer_name,
    indexer_protocol,
    title,
    download_url,
    info_url,
    guid,
    size_bytes,
    seeders,
    peers,
    published_at,
    search_kind,
    requested_season,
    requested_episode,
    sources
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id),
    sqlc.narg(indexer_id),
    sqlc.arg(indexer_name),
    sqlc.arg(indexer_protocol),
    sqlc.arg(title),
    sqlc.arg(download_url),
    sqlc.narg(info_url),
    sqlc.narg(guid),
    sqlc.arg(size_bytes),
    sqlc.narg(seeders),
    sqlc.narg(peers),
    sqlc.narg(published_at),
    sqlc.arg(search_kind),
    sqlc.narg(requested_season),
    sqlc.narg(requested_episode),
    sqlc.arg(sources)
);

-- name: AddReleaseSearchError :exec
insert into app.media_release_search_errors (id, media_item_id, message)
values (sqlc.arg(id), sqlc.arg(media_item_id), sqlc.arg(message));

-- name: GetReleaseCandidate :one
select id,
    media_item_id,
    season_id,
    episode_id,
    indexer_id,
    indexer_name,
    indexer_type,
    indexer_protocol,
    title,
    download_url,
    info_url,
    guid,
    size_bytes,
    seeders,
    peers,
    published_at,
    search_kind,
    requested_season,
    requested_episode,
    sources,
    created_at,
    updated_at
from app.media_release_candidates
where id = sqlc.arg(id)
    and media_item_id = sqlc.arg(media_item_id);

-- name: ListReleaseCandidates :many
select id,
    media_item_id,
    season_id,
    episode_id,
    indexer_id,
    indexer_name,
    indexer_type,
    indexer_protocol,
    title,
    download_url,
    info_url,
    guid,
    size_bytes,
    seeders,
    peers,
    published_at,
    search_kind,
    requested_season,
    requested_episode,
    sources,
    created_at,
    updated_at
from app.media_release_candidates
where media_item_id = $1
order by coalesce(seeders, -1) desc, size_bytes desc, created_at desc;

-- name: ListReleaseSearchErrors :many
select message
from app.media_release_search_errors
where media_item_id = $1
order by created_at asc;
