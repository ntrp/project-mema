-- name: ListMediaItems :many
select m.id,
    m.media_type,
    m.title,
    m.year,
    m.monitored,
    m.external_provider,
    m.external_id,
    m.overview,
    m.poster_path,
    m.collection_id,
    m.collection_name,
    m.backdrop_path,
    m.metadata_status,
    m.original_language,
    m.series_type,
    m.release_date,
    m.first_air_date,
    m.runtime_minutes,
    m.season_count,
    m.episode_count,
    m.vote_average,
    m.genres,
    m.keywords,
    m.facts,
    m.seasons,
    m.cast_members,
    m.crew_members,
    m.recommendations,
    m.similar_media,
    m.monitor_mode,
    m.minimum_availability,
    m.quality_profile_id,
    mp.name as quality_profile_name,
    case
        when exists (
            select 1
            from app.library_scan_items status_lsi
            where status_lsi.media_item_id = m.id
        ) then 'downloaded'
        when exists (
            select 1
            from app.download_activity status_activity
            where status_activity.media_item_id = m.id
                and status_activity.status in ('queued', 'grabbed', 'downloading')
        ) then 'downloading'
        when exists (
            select 1
            from app.download_activity status_activity
            where status_activity.media_item_id = m.id
                and status_activity.status = 'completed'
        ) then 'downloaded'
        else 'missing'
    end as status,
    m.library_folder_id,
    m.media_folder_path,
    coalesce(lf.path, (
        select lf2.path
        from app.library_scan_items lsi2
        join app.library_scans ls2 on ls2.id = lsi2.scan_id
        join app.library_folders lf2 on lf2.id = ls2.library_folder_id
        where lsi2.media_item_id = m.id
        order by lsi2.updated_at desc
        limit 1
    ), '')::text as library_folder_path,
    array(
        select distinct lsi.path
        from app.library_scan_items lsi
        where lsi.media_item_id = m.id
        order by lsi.path
    )::text[] as file_paths,
    coalesce(array(
        select t.name
        from app.media_item_tags mit
        join app.tags t on t.id = mit.tag_id
        where mit.media_item_id = m.id
        order by lower(t.name)
    ), '{}')::text[] as tags,
    m.created_at,
    m.updated_at
from app.media_items m
left join app.media_profiles mp on mp.id = m.quality_profile_id
left join app.library_folders lf on lf.id = m.library_folder_id
order by m.created_at desc, m.title asc;

-- name: SearchMediaItems :many
select m.id,
    m.media_type,
    m.title,
    m.year,
    m.monitored,
    m.external_provider,
    m.external_id,
    m.overview,
    m.poster_path,
    m.collection_id,
    m.collection_name,
    m.backdrop_path,
    m.metadata_status,
    m.original_language,
    m.series_type,
    m.release_date,
    m.first_air_date,
    m.runtime_minutes,
    m.season_count,
    m.episode_count,
    m.vote_average,
    m.genres,
    m.keywords,
    m.facts,
    m.seasons,
    m.cast_members,
    m.crew_members,
    m.recommendations,
    m.similar_media,
    m.monitor_mode,
    m.minimum_availability,
    m.quality_profile_id,
    mp.name as quality_profile_name,
    case
        when exists (select 1 from app.library_scan_items status_lsi where status_lsi.media_item_id = m.id) then 'downloaded'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status in ('queued', 'grabbed', 'downloading')) then 'downloading'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status = 'completed') then 'downloaded'
        else 'missing'
    end as status,
    m.library_folder_id,
    m.media_folder_path,
    coalesce(lf.path, (
        select lf2.path
        from app.library_scan_items lsi2
        join app.library_scans ls2 on ls2.id = lsi2.scan_id
        join app.library_folders lf2 on lf2.id = ls2.library_folder_id
        where lsi2.media_item_id = m.id
        order by lsi2.updated_at desc
        limit 1
    ), '')::text as library_folder_path,
    array(select distinct lsi.path from app.library_scan_items lsi where lsi.media_item_id = m.id order by lsi.path)::text[] as file_paths,
    coalesce(array(select t.name from app.media_item_tags mit join app.tags t on t.id = mit.tag_id where mit.media_item_id = m.id order by lower(t.name)), '{}')::text[] as tags,
    m.created_at,
    m.updated_at
from app.media_items m
left join app.media_profiles mp on mp.id = m.quality_profile_id
left join app.library_folders lf on lf.id = m.library_folder_id
where m.title ilike '%' || sqlc.arg(query)::text || '%'
    and (sqlc.narg(media_type)::text is null or m.media_type = sqlc.narg(media_type)::text)
order by
    case when lower(m.title) = lower(sqlc.arg(query)::text) then 0 else 1 end,
    m.title asc
limit sqlc.arg(row_limit);

-- name: GetMediaItem :one
select m.id,
    m.media_type,
    m.title,
    m.year,
    m.monitored,
    m.external_provider,
    m.external_id,
    m.overview,
    m.poster_path,
    m.collection_id,
    m.collection_name,
    m.backdrop_path,
    m.metadata_status,
    m.original_language,
    m.series_type,
    m.release_date,
    m.first_air_date,
    m.runtime_minutes,
    m.season_count,
    m.episode_count,
    m.vote_average,
    m.genres,
    m.keywords,
    m.facts,
    m.seasons,
    m.cast_members,
    m.crew_members,
    m.recommendations,
    m.similar_media,
    m.monitor_mode,
    m.minimum_availability,
    m.quality_profile_id,
    mp.name as quality_profile_name,
    case
        when exists (select 1 from app.library_scan_items status_lsi where status_lsi.media_item_id = m.id) then 'downloaded'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status in ('queued', 'grabbed', 'downloading')) then 'downloading'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status = 'completed') then 'downloaded'
        else 'missing'
    end as status,
    m.library_folder_id,
    m.media_folder_path,
    coalesce(lf.path, (
        select lf2.path
        from app.library_scan_items lsi2
        join app.library_scans ls2 on ls2.id = lsi2.scan_id
        join app.library_folders lf2 on lf2.id = ls2.library_folder_id
        where lsi2.media_item_id = m.id
        order by lsi2.updated_at desc
        limit 1
    ), '')::text as library_folder_path,
    array(select distinct lsi.path from app.library_scan_items lsi where lsi.media_item_id = m.id order by lsi.path)::text[] as file_paths,
    coalesce(array(select t.name from app.media_item_tags mit join app.tags t on t.id = mit.tag_id where mit.media_item_id = m.id order by lower(t.name)), '{}')::text[] as tags,
    m.created_at,
    m.updated_at
from app.media_items m
left join app.media_profiles mp on mp.id = m.quality_profile_id
left join app.library_folders lf on lf.id = m.library_folder_id
where m.id = $1;

-- name: ListMissingMediaItems :many
select m.id,
    m.media_type,
    m.title,
    m.year,
    m.monitored,
    m.external_provider,
    m.external_id,
    m.overview,
    m.poster_path,
    m.collection_id,
    m.collection_name,
    m.backdrop_path,
    m.metadata_status,
    m.original_language,
    m.series_type,
    m.release_date,
    m.first_air_date,
    m.runtime_minutes,
    m.season_count,
    m.episode_count,
    m.vote_average,
    m.genres,
    m.keywords,
    m.facts,
    m.seasons,
    m.cast_members,
    m.crew_members,
    m.recommendations,
    m.similar_media,
    m.monitor_mode,
    m.minimum_availability,
    m.quality_profile_id,
    mp.name as quality_profile_name,
    case
        when exists (select 1 from app.library_scan_items status_lsi where status_lsi.media_item_id = m.id) then 'downloaded'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status in ('queued', 'grabbed', 'downloading')) then 'downloading'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status = 'completed') then 'downloaded'
        else 'missing'
    end as status,
    m.library_folder_id,
    m.media_folder_path,
    coalesce(lf.path, (
        select lf2.path
        from app.library_scan_items lsi2
        join app.library_scans ls2 on ls2.id = lsi2.scan_id
        join app.library_folders lf2 on lf2.id = ls2.library_folder_id
        where lsi2.media_item_id = m.id
        order by lsi2.updated_at desc
        limit 1
    ), '')::text as library_folder_path,
    array(select distinct lsi.path from app.library_scan_items lsi where lsi.media_item_id = m.id order by lsi.path)::text[] as file_paths,
    coalesce(array(select t.name from app.media_item_tags mit join app.tags t on t.id = mit.tag_id where mit.media_item_id = m.id order by lower(t.name)), '{}')::text[] as tags,
    m.created_at,
    m.updated_at
from app.media_items m
left join app.media_profiles mp on mp.id = m.quality_profile_id
left join app.library_folders lf on lf.id = m.library_folder_id
where m.monitored = true
    and not exists (
        select 1
        from app.library_scan_items lsi
        where lsi.media_item_id = m.id
    )
    and not exists (
        select 1
        from app.download_activity activity
        where activity.media_item_id = m.id
            and activity.status in ('queued', 'grabbed', 'downloading')
    )
order by m.created_at asc;

-- name: FindMonitoredMediaMatch :one
select m.id,
    m.media_type,
    m.title,
    m.year,
    m.monitored,
    m.external_provider,
    m.external_id,
    m.overview,
    m.poster_path,
    m.collection_id,
    m.collection_name,
    m.backdrop_path,
    m.metadata_status,
    m.original_language,
    m.series_type,
    m.release_date,
    m.first_air_date,
    m.runtime_minutes,
    m.season_count,
    m.episode_count,
    m.vote_average,
    m.genres,
    m.keywords,
    m.facts,
    m.seasons,
    m.cast_members,
    m.crew_members,
    m.recommendations,
    m.similar_media,
    m.monitor_mode,
    m.minimum_availability,
    m.quality_profile_id,
    mp.name as quality_profile_name,
    case
        when exists (select 1 from app.library_scan_items status_lsi where status_lsi.media_item_id = m.id) then 'downloaded'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status in ('queued', 'grabbed', 'downloading')) then 'downloading'
        when exists (select 1 from app.download_activity status_activity where status_activity.media_item_id = m.id and status_activity.status = 'completed') then 'downloaded'
        else 'missing'
    end as status,
    m.library_folder_id,
    m.media_folder_path,
    coalesce(lf.path, (
        select lf2.path
        from app.library_scan_items lsi2
        join app.library_scans ls2 on ls2.id = lsi2.scan_id
        join app.library_folders lf2 on lf2.id = ls2.library_folder_id
        where lsi2.media_item_id = m.id
        order by lsi2.updated_at desc
        limit 1
    ), '')::text as library_folder_path,
    array(select distinct lsi.path from app.library_scan_items lsi where lsi.media_item_id = m.id order by lsi.path)::text[] as file_paths,
    coalesce(array(select t.name from app.media_item_tags mit join app.tags t on t.id = mit.tag_id where mit.media_item_id = m.id order by lower(t.name)), '{}')::text[] as tags,
    m.created_at,
    m.updated_at
from app.media_items m
left join app.media_profiles mp on mp.id = m.quality_profile_id
left join app.library_folders lf on lf.id = m.library_folder_id
where m.monitored = true
    and m.quality_profile_id is not null
    and lower(trim(regexp_replace(m.title, '[^[:alnum:]]+', ' ', 'g'))) =
        lower(trim(regexp_replace(sqlc.arg(title)::text, '[^[:alnum:]]+', ' ', 'g')))
    and (sqlc.narg(year)::integer is null or m.year = sqlc.narg(year)::integer)
order by
    case when m.year = sqlc.narg(year)::integer then 0 when m.year is null then 1 else 2 end,
    m.updated_at desc
limit 1;

-- name: CreateMediaItemRecord :one
insert into app.media_items (
    id,
    media_type,
    content_kind,
    title,
    year,
    monitored,
    external_provider,
    external_id,
    overview,
    poster_path,
    collection_id,
    collection_name,
    backdrop_path,
    metadata_status,
    original_language,
    series_type,
    numbering_strategy,
    release_date,
    first_air_date,
    runtime_minutes,
    season_count,
    episode_count,
    vote_average,
    genres,
    keywords,
    facts,
    seasons,
    cast_members,
    crew_members,
    recommendations,
    similar_media,
    monitor_mode,
    minimum_availability,
    quality_profile_id,
    library_folder_id,
    media_folder_path
)
values (
    sqlc.arg(id),
    sqlc.arg(media_type),
    sqlc.arg(content_kind),
    sqlc.arg(title),
    sqlc.narg(year),
    sqlc.arg(monitored),
    sqlc.narg(external_provider),
    sqlc.narg(external_id),
    sqlc.narg(overview),
    sqlc.narg(poster_path),
    sqlc.narg(collection_id),
    sqlc.narg(collection_name),
    sqlc.narg(backdrop_path),
    sqlc.narg(metadata_status),
    sqlc.narg(original_language),
    sqlc.narg(series_type),
    sqlc.narg(numbering_strategy),
    sqlc.narg(release_date),
    sqlc.narg(first_air_date),
    sqlc.narg(runtime_minutes),
    sqlc.narg(season_count),
    sqlc.narg(episode_count),
    sqlc.narg(vote_average),
    sqlc.arg(genres),
    sqlc.arg(keywords),
    sqlc.arg(facts),
    sqlc.arg(seasons),
    sqlc.arg(cast_members),
    sqlc.arg(crew_members),
    sqlc.arg(recommendations),
    sqlc.arg(similar_media),
    sqlc.arg(monitor_mode),
    sqlc.arg(minimum_availability),
    sqlc.narg(quality_profile_id),
    sqlc.narg(library_folder_id),
    sqlc.narg(media_folder_path)
)
returning id;

-- name: FindExistingMediaItemID :one
select id
from app.media_items
where lower(media_type) = lower(sqlc.arg(media_type)::text)
    and lower(title) = lower(sqlc.arg(title)::text)
    and ((sqlc.narg(year)::integer is null and year is null) or year = sqlc.narg(year)::integer)
order by created_at asc
limit 1;

-- name: UpdateExistingMediaItem :exec
update app.media_items
set quality_profile_id = coalesce(quality_profile_id, sqlc.narg(quality_profile_id)::text),
    library_folder_id = coalesce(library_folder_id, sqlc.narg(library_folder_id)::uuid),
    media_folder_path = coalesce(media_folder_path, sqlc.narg(media_folder_path)::text),
    monitor_mode = sqlc.arg(monitor_mode),
    minimum_availability = sqlc.arg(minimum_availability),
    monitored = sqlc.arg(monitored),
    series_type = coalesce(sqlc.narg(series_type)::text, series_type),
    content_kind = sqlc.arg(content_kind),
    numbering_strategy = coalesce(sqlc.narg(numbering_strategy)::text, numbering_strategy),
    updated_at = case
        when (quality_profile_id is null and sqlc.narg(quality_profile_id)::text is not null)
            or (library_folder_id is null and sqlc.narg(library_folder_id)::uuid is not null)
            or (media_folder_path is null and sqlc.narg(media_folder_path)::text is not null)
            or monitor_mode <> sqlc.arg(monitor_mode)
            or minimum_availability <> sqlc.arg(minimum_availability)
            or monitored <> sqlc.arg(monitored)
            or (sqlc.narg(series_type)::text is not null and series_type is distinct from sqlc.narg(series_type)::text)
            or content_kind <> sqlc.arg(content_kind)
            or (sqlc.narg(numbering_strategy)::text is not null and numbering_strategy is distinct from sqlc.narg(numbering_strategy)::text)
        then now()
        else updated_at
    end
where id = sqlc.arg(id);

-- name: GetMediaItemAnimeState :one
select content_kind, numbering_strategy
from app.media_items
where id = $1;
