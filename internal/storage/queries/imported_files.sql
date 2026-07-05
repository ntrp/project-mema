-- name: ImportedMediaFileExists :one
select exists(
    select 1
    from app.library_scan_items
    where media_item_id = $1 and path = $2
);

-- name: CreateImportedFileLibraryScan :exec
insert into app.library_scans (
    id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
)
values ($1, $2, 'completed', 1, 1, 0, now());

-- name: CreateImportedFileLibraryScanItem :exec
insert into app.library_scan_items (
    id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
    status, matched_title, matched_year, matched_media_kind, media_item_id, season_id, episode_id
)
values (
    sqlc.arg(id),
    sqlc.arg(scan_id),
    sqlc.arg(path),
    sqlc.arg(file_name),
    sqlc.arg(detected_title),
    sqlc.narg(detected_year),
    sqlc.arg(detected_media_kind),
    'auto_added',
    sqlc.arg(detected_title),
    sqlc.narg(detected_year),
    sqlc.arg(detected_media_kind),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id)
);

-- name: GetImportedFileEpisodeReference :one
select season_id, episode_id
from app.library_scan_items
where media_item_id = $1
    and path = $2;

-- name: DeleteLibraryScanItemsForMediaItem :exec
delete from app.library_scan_items
where media_item_id = $1;

-- name: ListMediaFileRecordsForItem :many
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
    season_id,
    episode_id,
    created_at,
    updated_at
from app.library_scan_items
where media_item_id = $1
order by updated_at desc;

-- name: UpdateLibraryScanItemStatus :exec
update app.library_scan_items
set status = $2,
    updated_at = now()
where id = $1;

-- name: RenameMediaFileRecord :execrows
update app.library_scan_items
set path = sqlc.arg(destination_path),
    file_name = sqlc.arg(file_name),
    status = 'restored',
    updated_at = now()
where media_item_id = sqlc.arg(media_item_id)
    and path = sqlc.arg(source_path);

-- name: CreateMediaFileRescanLibraryScan :exec
insert into app.library_scans (
    id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
)
values ($1, $2, 'completed', $3, $3, 0, now());

-- name: CreateMediaFileRescanLibraryScanItem :exec
insert into app.library_scan_items (
    id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
    status, matched_title, matched_year, matched_media_kind, media_item_id, season_id, episode_id
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
    sqlc.arg(detected_title),
    sqlc.narg(detected_year),
    sqlc.arg(detected_media_kind),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id)
);
