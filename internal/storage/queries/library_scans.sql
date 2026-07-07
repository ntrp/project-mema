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
select distinct on (lsi.path)
    lsi.path,
    mi.id as media_item_id,
    mi.title as matched_title,
    mi.year as matched_year
from app.library_scan_items lsi
join app.library_scans ls on ls.id = lsi.scan_id
join app.media_items mi on mi.id = lsi.media_item_id
where ls.library_folder_id = $1
    and lsi.media_item_id is not null
    and lsi.status in ('auto_added', 'manually_added', 'restored')
order by lsi.path, lsi.updated_at desc;

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

-- name: ResetLibraryScanItemImport :one
with current as (
    select id,
        media_item_id as previous_media_item_id,
        match_source as previous_match_source
    from app.library_scan_items
    where app.library_scan_items.scan_id = sqlc.arg(scan_id)
        and app.library_scan_items.id = sqlc.arg(id)
        and app.library_scan_items.media_item_id is not null
),
reset as (
    update app.library_scan_items as lsi
    set status = 'pending',
        imported = false,
        matched_title = null,
        matched_year = null,
        matched_media_kind = null,
        matched_external_provider = null,
        matched_external_id = null,
        match_source = null,
        selected_metadata_provider_id = null,
        media_item_id = null,
        season_id = null,
        episode_id = null,
        updated_at = now()
    from current
    where lsi.id = current.id
        and lsi.scan_id = sqlc.arg(scan_id)
    returning lsi.id,
        lsi.scan_id,
        lsi.path,
        lsi.file_name,
        lsi.size_bytes,
        lsi.detected_title,
        lsi.detected_year,
        lsi.detected_media_kind,
        lsi.season_number,
        lsi.episode_number,
        lsi.status,
        lsi.imported,
        lsi.matched_title,
        lsi.matched_year,
        lsi.matched_media_kind,
        lsi.matched_external_provider,
        lsi.matched_external_id,
        lsi.match_source,
        lsi.selected_metadata_provider_id,
        lsi.duplicate_group_id,
        lsi.duplicate_removal_allowed,
        lsi.media_item_id,
        lsi.created_at,
        lsi.updated_at,
        current.previous_media_item_id,
        current.previous_match_source
)
select reset.*
from reset;

-- name: ResetLibraryScanItemsForMediaItem :exec
update app.library_scan_items
set status = 'pending',
    imported = false,
    matched_title = null,
    matched_year = null,
    matched_media_kind = null,
    matched_external_provider = null,
    matched_external_id = null,
    match_source = null,
    selected_metadata_provider_id = null,
    media_item_id = null,
    season_id = null,
    episode_id = null,
    updated_at = now()
where media_item_id = $1;

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
