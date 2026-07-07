-- name: DeleteMediaItemRecord :execrows
delete from app.media_items
where id = $1;

-- name: DeleteMediaItemsForLibraryFolder :execrows
with folder_root as (
    select id,
        case when path = '/' then '/' else rtrim(path, '/') end as path
    from app.library_folders
    where app.library_folders.id = $1
),
media_to_delete as (
    select m.id
    from app.media_items m
    join folder_root f on true
    where m.library_folder_id = f.id
        or (
            m.media_folder_path is not null
            and (
                m.media_folder_path = f.path
                or (f.path = '/' and m.media_folder_path like '/%')
                or (f.path <> '/' and m.media_folder_path like f.path || '/%')
            )
        )
),
deleted_import_attempts as (
    delete from app.import_attempts
    where media_item_id in (select id from media_to_delete)
),
deleted_file_history as (
    delete from app.media_file_history
    where media_item_id in (select id from media_to_delete)
)
delete from app.media_items
where id in (select id from media_to_delete);

-- name: TouchMediaItem :exec
update app.media_items
set updated_at = now()
where id = $1;

-- name: UpdateMediaItemMetadataRecord :exec
update app.media_items
set
    media_type = sqlc.arg(media_type),
    content_kind = sqlc.arg(content_kind),
    title = sqlc.arg(title),
    year = sqlc.narg(year),
    external_provider = sqlc.narg(external_provider),
    external_id = sqlc.narg(external_id),
    overview = sqlc.narg(overview),
    poster_path = sqlc.narg(poster_path),
    collection_id = sqlc.narg(collection_id),
    collection_name = sqlc.narg(collection_name),
    backdrop_path = sqlc.narg(backdrop_path),
    metadata_status = sqlc.narg(metadata_status),
    original_language = sqlc.narg(original_language),
    numbering_strategy = sqlc.narg(numbering_strategy),
    release_date = sqlc.narg(release_date),
    first_air_date = sqlc.narg(first_air_date),
    runtime_minutes = sqlc.narg(runtime_minutes),
    season_count = sqlc.narg(season_count),
    episode_count = sqlc.narg(episode_count),
    vote_average = sqlc.narg(vote_average),
    genres = sqlc.arg(genres)::jsonb,
    keywords = sqlc.arg(keywords)::jsonb,
    facts = sqlc.arg(facts)::jsonb,
    seasons = sqlc.arg(seasons)::jsonb,
    cast_members = sqlc.arg(cast_members)::jsonb,
    crew_members = sqlc.arg(crew_members)::jsonb,
    recommendations = sqlc.arg(recommendations)::jsonb,
    similar_media = sqlc.arg(similar_media)::jsonb,
    updated_at = now()
where id = sqlc.arg(id);

-- name: UpdateMediaItemOptionsRecord :execrows
update app.media_items
set quality_profile_id = coalesce(sqlc.narg(quality_profile_id)::text, quality_profile_id),
    minimum_availability = coalesce(sqlc.narg(minimum_availability)::text, minimum_availability),
    monitored = coalesce(sqlc.narg(monitored)::boolean, monitored),
    monitor_mode = coalesce(sqlc.narg(monitor_mode)::text, monitor_mode),
    seasons = case when sqlc.arg(update_seasons)::boolean then sqlc.arg(seasons)::jsonb else seasons end,
    library_folder_id = coalesce(sqlc.narg(library_folder_id)::uuid, library_folder_id),
    media_folder_path = coalesce(sqlc.narg(media_folder_path)::text, media_folder_path),
    updated_at = now()
where id = sqlc.arg(id);
