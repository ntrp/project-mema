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
    size_bytes,
    detected_title,
    detected_year,
    detected_media_kind,
    season_number,
    episode_number,
    status,
    imported,
    matched_title,
    matched_year,
    matched_media_kind,
    matched_external_provider,
    matched_external_id,
    match_source,
    selected_metadata_provider_id,
    duplicate_group_id,
    duplicate_removal_allowed,
    media_item_id
)
values (
    sqlc.arg(id),
    sqlc.arg(scan_id),
    sqlc.arg(path),
    sqlc.arg(file_name),
    sqlc.arg(size_bytes),
    sqlc.arg(detected_title),
    sqlc.narg(detected_year),
    sqlc.arg(detected_media_kind),
    sqlc.narg(season_number),
    sqlc.narg(episode_number),
    sqlc.arg(status),
    sqlc.arg(imported),
    sqlc.narg(matched_title),
    sqlc.narg(matched_year),
    sqlc.narg(matched_media_kind),
    sqlc.narg(matched_external_provider),
    sqlc.narg(matched_external_id),
    sqlc.narg(match_source),
    sqlc.narg(selected_metadata_provider_id),
    sqlc.narg(duplicate_group_id),
    sqlc.arg(duplicate_removal_allowed),
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
    f.kind as folder_kind,
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

-- name: ListActiveImportedPathsForLibraryFolder :many
select distinct lsi.path
from app.library_scan_items lsi
join app.library_scans ls on ls.id = lsi.scan_id
where ls.library_folder_id = $1
    and lsi.media_item_id is not null
    and lsi.status in ('auto_added', 'manually_added', 'restored');

-- name: GetLibraryScanItemPath :one
select path
from app.library_scan_items
where scan_id = sqlc.arg(scan_id)
    and id = sqlc.arg(id);

-- name: MatchLibraryScanItem :one
update app.library_scan_items
set status = 'manually_added',
    imported = sqlc.arg(imported),
    matched_title = sqlc.arg(matched_title),
    matched_year = sqlc.narg(matched_year),
    matched_media_kind = sqlc.arg(matched_media_kind),
    matched_external_provider = sqlc.narg(matched_external_provider),
    matched_external_id = sqlc.narg(matched_external_id),
    match_source = sqlc.arg(match_source),
    selected_metadata_provider_id = sqlc.narg(selected_metadata_provider_id),
    media_item_id = sqlc.arg(media_item_id),
    season_id = sqlc.narg(season_id),
    episode_id = sqlc.narg(episode_id),
    updated_at = now()
where scan_id = sqlc.arg(scan_id)
    and id = sqlc.arg(id)
    and status = 'pending'
returning id,
    scan_id,
    path,
    file_name,
    size_bytes,
    detected_title,
    detected_year,
    detected_media_kind,
    season_number,
    episode_number,
    status,
    imported,
    matched_title,
    matched_year,
    matched_media_kind,
    matched_external_provider,
    matched_external_id,
    match_source,
    selected_metadata_provider_id,
    duplicate_group_id,
    duplicate_removal_allowed,
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
    size_bytes,
    detected_title,
    detected_year,
    detected_media_kind,
    season_number,
    episode_number,
    status,
    imported,
    matched_title,
    matched_year,
    matched_media_kind,
    matched_external_provider,
    matched_external_id,
    match_source,
    selected_metadata_provider_id,
    duplicate_group_id,
    duplicate_removal_allowed,
    media_item_id,
    created_at,
    updated_at
from app.library_scan_items
where scan_id = $1
order by status desc, path asc;
