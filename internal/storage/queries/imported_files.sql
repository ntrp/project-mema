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
    status, matched_title, matched_year, matched_media_kind, media_item_id
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
    sqlc.arg(media_item_id)
);

-- name: DeleteLibraryScanItemsForMediaItem :exec
delete from app.library_scan_items
where media_item_id = $1;

-- name: CreateMediaFileRescanLibraryScan :exec
insert into app.library_scans (
    id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
)
values ($1, $2, 'completed', $3, $3, 0, now());
