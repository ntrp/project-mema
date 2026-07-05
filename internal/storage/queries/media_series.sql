-- name: ListMediaSeasonRows :many
select *
from app.media_seasons
where media_item_id = $1
order by season_number asc;

-- name: ListMediaEpisodeRows :many
select *
from app.media_episodes
where media_item_id = $1
order by season_number asc, episode_number asc;

-- name: UpsertMediaSeasonRow :one
insert into app.media_seasons (
    id,
    media_item_id,
    external_provider,
    external_id,
    season_number,
    name,
    overview,
    air_date,
    poster_path,
    episode_count,
    monitored,
    source
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(external_provider),
    sqlc.narg(external_id),
    sqlc.arg(season_number),
    sqlc.arg(name),
    sqlc.narg(overview),
    sqlc.narg(air_date),
    sqlc.narg(poster_path),
    sqlc.narg(episode_count),
    sqlc.arg(monitored),
    sqlc.arg(source)::jsonb
)
on conflict (media_item_id, season_number) do update
set external_provider = excluded.external_provider,
    external_id = excluded.external_id,
    name = excluded.name,
    overview = excluded.overview,
    air_date = excluded.air_date,
    poster_path = excluded.poster_path,
    episode_count = excluded.episode_count,
    monitored = excluded.monitored,
    source = excluded.source,
    updated_at = now()
returning *;

-- name: UpsertMediaEpisodeRow :one
insert into app.media_episodes (
    id,
    season_id,
    media_item_id,
    external_provider,
    external_id,
    season_number,
    episode_number,
    name,
    overview,
    air_date,
    still_path,
    runtime_minutes,
    monitored,
    source
) values (
    sqlc.arg(id),
    sqlc.arg(season_id),
    sqlc.arg(media_item_id),
    sqlc.narg(external_provider),
    sqlc.narg(external_id),
    sqlc.arg(season_number),
    sqlc.arg(episode_number),
    sqlc.arg(name),
    sqlc.narg(overview),
    sqlc.narg(air_date),
    sqlc.narg(still_path),
    sqlc.narg(runtime_minutes),
    sqlc.arg(monitored),
    sqlc.arg(source)::jsonb
)
on conflict (media_item_id, season_number, episode_number) do update
set season_id = excluded.season_id,
    external_provider = excluded.external_provider,
    external_id = excluded.external_id,
    name = excluded.name,
    overview = excluded.overview,
    air_date = excluded.air_date,
    still_path = excluded.still_path,
    runtime_minutes = excluded.runtime_minutes,
    monitored = excluded.monitored,
    source = excluded.source,
    updated_at = now()
returning *;

-- name: UpdateMediaSeasonMonitoredRow :one
update app.media_seasons
set monitored = $2,
    updated_at = now()
where id = $1
returning *;

-- name: UpdateMediaSeasonEpisodesMonitored :exec
update app.media_episodes
set monitored = $2,
    updated_at = now()
where season_id = $1;

-- name: UpdateMediaEpisodeMonitoredRow :one
update app.media_episodes
set monitored = $2,
    updated_at = now()
where id = $1
returning *;

-- name: SyncMediaSeasonMonitoredFromEpisodes :one
update app.media_seasons
set monitored = exists (
        select 1
        from app.media_episodes
        where season_id = sqlc.arg(season_id)
            and monitored
    ),
    updated_at = now()
where id = sqlc.arg(season_id)
returning *;
