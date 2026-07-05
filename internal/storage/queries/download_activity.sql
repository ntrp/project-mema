-- name: CreateDownloadActivity :one
insert into app.download_activity (
    id, media_item_id, release_title, indexer_name, download_client_name,
    download_id, download_url, status, error, failure_type
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(release_title),
    sqlc.arg(indexer_name),
    sqlc.arg(download_client_name),
    sqlc.narg(download_id),
    sqlc.arg(download_url),
    sqlc.arg(status),
    sqlc.narg(error),
    sqlc.narg(failure_type)
)
returning id, media_item_id, release_title, indexer_name, download_client_name,
    download_id, download_url, status, progress_percent, error, failure_type,
    created_at, updated_at;

-- name: FailDownloadActivity :one
update app.download_activity
set status = 'failed',
    progress_percent = null,
    error = sqlc.narg(error),
    failure_type = sqlc.arg(failure_type),
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, release_title, indexer_name, download_client_name,
    download_id, download_url, status, progress_percent, error, failure_type,
    created_at, updated_at;

-- name: UpdateDownloadActivity :one
update app.download_activity
set status = sqlc.arg(status),
    download_id = coalesce(sqlc.narg(download_id)::text, download_id),
    progress_percent = sqlc.narg(progress_percent),
    error = sqlc.narg(error),
    failure_type = null,
    updated_at = now()
where id = sqlc.arg(id)
returning id, media_item_id, release_title, indexer_name, download_client_name,
    download_id, download_url, status, progress_percent, error, failure_type,
    created_at, updated_at;

-- name: ListDownloadActivity :many
select
    a.id,
    a.media_item_id,
    m.title as media_title,
    m.media_type,
    m.year as media_year,
    a.release_title,
    a.indexer_name,
    a.download_client_name,
    a.download_id,
    a.download_url,
    a.status,
    a.progress_percent,
    a.error,
    a.failure_type,
    a.created_at,
    a.updated_at
from app.download_activity a
join app.media_items m on m.id = a.media_item_id
order by a.created_at desc
limit 100;

-- name: GetDownloadActivity :one
select
    a.id,
    a.media_item_id,
    m.title as media_title,
    m.media_type,
    m.year as media_year,
    a.release_title,
    a.indexer_name,
    a.download_client_name,
    a.download_id,
    a.download_url,
    a.status,
    a.progress_percent,
    a.error,
    a.failure_type,
    a.created_at,
    a.updated_at
from app.download_activity a
join app.media_items m on m.id = a.media_item_id
where a.id = $1;

-- name: CancelDownloadActivity :one
update app.download_activity
set status = 'cancelled',
    progress_percent = null,
    error = null,
    failure_type = null,
    updated_at = now()
where id = $1
    and status in ('queued', 'grabbed', 'downloading')
returning id, media_item_id, release_title, indexer_name, download_client_name,
    download_id, download_url, status, progress_percent, error, failure_type,
    created_at, updated_at;

-- name: DeleteDownloadActivity :execrows
delete from app.download_activity
where id = $1
    and status in ('failed', 'cancelled');

-- name: ListActiveDownloadActivity :many
select
    a.id,
    a.media_item_id,
    m.title as media_title,
    m.media_type,
    m.year as media_year,
    a.release_title,
    a.indexer_name,
    a.download_client_name,
    a.download_id,
    a.download_url,
    a.status,
    a.progress_percent,
    a.error,
    a.failure_type,
    a.created_at,
    a.updated_at
from app.download_activity a
join app.media_items m on m.id = a.media_item_id
where a.status in ('queued', 'grabbed', 'downloading')
    and a.download_id is not null
order by a.updated_at asc;
