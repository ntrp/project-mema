-- name: ListMediaRequests :many
select r.id,
    r.requested_by_user_id,
    u.username as requested_by_username,
    r.media_type,
    r.title,
    r.year,
    r.external_provider,
    r.external_id,
    r.overview,
    r.poster_path,
    r.status,
    r.monitor_mode,
    r.series_type,
    r.minimum_availability,
    r.quality_profile_id,
    r.library_folder_id,
    r.media_item_id,
    r.decided_at,
    coalesce(array(
        select t.name
        from app.media_request_tags mrt
        join app.tags t on t.id = mrt.tag_id
        where mrt.media_request_id = r.id
        order by lower(t.name)
    ), '{}')::text[] as tags,
    r.created_at,
    r.updated_at
from app.media_requests r
join app.users u on u.id = r.requested_by_user_id
where sqlc.arg(include_all)::boolean = true
    or r.requested_by_user_id = sqlc.arg(user_id)
order by
    case r.status when 'pending' then 0 else 1 end,
    r.created_at desc;

-- name: GetMediaRequestForUser :one
select r.id,
    r.requested_by_user_id,
    u.username as requested_by_username,
    r.media_type,
    r.title,
    r.year,
    r.external_provider,
    r.external_id,
    r.overview,
    r.poster_path,
    r.status,
    r.monitor_mode,
    r.series_type,
    r.minimum_availability,
    r.quality_profile_id,
    r.library_folder_id,
    r.media_item_id,
    r.decided_at,
    coalesce(array(
        select t.name
        from app.media_request_tags mrt
        join app.tags t on t.id = mrt.tag_id
        where mrt.media_request_id = r.id
        order by lower(t.name)
    ), '{}')::text[] as tags,
    r.created_at,
    r.updated_at
from app.media_requests r
join app.users u on u.id = r.requested_by_user_id
where r.id = sqlc.arg(id)
    and (sqlc.arg(include_all)::boolean = true or r.requested_by_user_id = sqlc.arg(user_id));

-- name: GetMediaRequest :one
select r.id,
    r.requested_by_user_id,
    u.username as requested_by_username,
    r.media_type,
    r.title,
    r.year,
    r.external_provider,
    r.external_id,
    r.overview,
    r.poster_path,
    r.status,
    r.monitor_mode,
    r.series_type,
    r.minimum_availability,
    r.quality_profile_id,
    r.library_folder_id,
    r.media_item_id,
    r.decided_at,
    coalesce(array(
        select t.name
        from app.media_request_tags mrt
        join app.tags t on t.id = mrt.tag_id
        where mrt.media_request_id = r.id
        order by lower(t.name)
    ), '{}')::text[] as tags,
    r.created_at,
    r.updated_at
from app.media_requests r
join app.users u on u.id = r.requested_by_user_id
where r.id = $1;

-- name: GetMediaRequestForUpdate :one
select r.id,
    r.requested_by_user_id,
    u.username as requested_by_username,
    r.media_type,
    r.title,
    r.year,
    r.external_provider,
    r.external_id,
    r.overview,
    r.poster_path,
    r.status,
    r.monitor_mode,
    r.series_type,
    r.minimum_availability,
    r.quality_profile_id,
    r.library_folder_id,
    r.media_item_id,
    r.decided_at,
    coalesce(array(
        select t.name
        from app.media_request_tags mrt
        join app.tags t on t.id = mrt.tag_id
        where mrt.media_request_id = r.id
        order by lower(t.name)
    ), '{}')::text[] as tags,
    r.created_at,
    r.updated_at
from app.media_requests r
join app.users u on u.id = r.requested_by_user_id
where r.id = $1
for update;

-- name: CreateMediaRequest :one
insert into app.media_requests (
    id,
    requested_by_user_id,
    media_type,
    title,
    year,
    external_provider,
    external_id,
    overview,
    poster_path,
    monitor_mode,
    series_type,
    minimum_availability
)
values (
    sqlc.arg(id),
    sqlc.arg(requested_by_user_id),
    sqlc.arg(media_type),
    sqlc.arg(title),
    sqlc.narg(year),
    sqlc.narg(external_provider),
    sqlc.narg(external_id),
    sqlc.narg(overview),
    sqlc.narg(poster_path),
    sqlc.arg(monitor_mode),
    sqlc.narg(series_type),
    sqlc.arg(minimum_availability)
)
returning id;

-- name: ApproveMediaRequest :one
update app.media_requests
set status = 'approved',
    quality_profile_id = sqlc.arg(quality_profile_id),
    library_folder_id = sqlc.arg(library_folder_id),
    media_item_id = sqlc.arg(media_item_id),
    decided_at = now(),
    updated_at = now()
where app.media_requests.id = sqlc.arg(id)
returning app.media_requests.id,
    app.media_requests.requested_by_user_id,
    (
        select username from app.users where app.users.id = app.media_requests.requested_by_user_id
    )::text as requested_by_username,
    app.media_requests.media_type,
    app.media_requests.title,
    app.media_requests.year,
    app.media_requests.external_provider,
    app.media_requests.external_id,
    app.media_requests.overview,
    app.media_requests.poster_path,
    app.media_requests.status,
    app.media_requests.monitor_mode,
    app.media_requests.series_type,
    app.media_requests.minimum_availability,
    app.media_requests.quality_profile_id,
    app.media_requests.library_folder_id,
    app.media_requests.media_item_id,
    app.media_requests.decided_at,
    coalesce(array(
        select t.name
        from app.media_request_tags mrt
        join app.tags t on t.id = mrt.tag_id
        where mrt.media_request_id = sqlc.arg(id)
        order by lower(t.name)
    ), '{}')::text[] as tags,
    app.media_requests.created_at,
    app.media_requests.updated_at;
