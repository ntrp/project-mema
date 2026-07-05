-- name: CreateLibraryScan :exec
insert into app.library_scans (
    id,
    library_folder_id,
    status,
    total_files,
    auto_matched_count,
    manual_count,
    completed_at
)
values (
    sqlc.arg(id),
    sqlc.arg(library_folder_id),
    'completed',
    sqlc.arg(total_files),
    0,
    0,
    now()
);

-- name: AddLibraryScanItem :exec
insert into app.library_scan_items (
    id,
    scan_id,
    path,
    file_name,
    detected_title,
    detected_year,
    detected_media_kind,
    status,
    matched_title,
    matched_year,
    matched_media_kind,
    media_item_id
)
values (
    sqlc.arg(id),
    sqlc.arg(scan_id),
    sqlc.arg(path),
    sqlc.arg(file_name),
    sqlc.arg(detected_title),
    sqlc.narg(detected_year),
    sqlc.arg(detected_media_kind),
    sqlc.arg(status),
    sqlc.narg(matched_title),
    sqlc.narg(matched_year),
    sqlc.narg(matched_media_kind),
    sqlc.narg(media_item_id)
);

-- name: UpdateLibraryScanCounts :exec
update app.library_scans
set auto_matched_count = sqlc.arg(auto_matched_count),
    manual_count = sqlc.arg(manual_count)
where id = sqlc.arg(id);

-- name: GetLibraryScan :one
select s.id,
    s.library_folder_id,
    f.path as folder_path,
    s.status,
    s.total_files,
    s.auto_matched_count,
    s.manual_count,
    s.created_at,
    s.completed_at
from app.library_scans s
join app.library_folders f on f.id = s.library_folder_id
where s.id = $1;

-- name: GetLibraryScanFolderID :one
select library_folder_id
from app.library_scans
where id = $1;

-- name: MatchLibraryScanItem :one
update app.library_scan_items
set status = 'manually_added',
    matched_title = sqlc.arg(matched_title),
    matched_year = sqlc.narg(matched_year),
    matched_media_kind = sqlc.arg(matched_media_kind),
    media_item_id = sqlc.arg(media_item_id),
    updated_at = now()
where scan_id = sqlc.arg(scan_id)
    and id = sqlc.arg(id)
    and status = 'pending'
returning id,
    scan_id,
    path,
    file_name,
    detected_title,
    detected_year,
    detected_media_kind,
    status,
    matched_title,
    matched_year,
    matched_media_kind,
    media_item_id,
    created_at,
    updated_at;

-- name: RefreshLibraryScanManualCount :exec
update app.library_scans
set manual_count = (
    select count(*)::integer
    from app.library_scan_items
    where scan_id = sqlc.arg(id) and status = 'pending'
)
where id = sqlc.arg(id);

-- name: ListLibraryScanItems :many
select id,
    scan_id,
    path,
    file_name,
    detected_title,
    detected_year,
    detected_media_kind,
    status,
    matched_title,
    matched_year,
    matched_media_kind,
    media_item_id,
    created_at,
    updated_at
from app.library_scan_items
where scan_id = $1
order by status desc, path asc;
