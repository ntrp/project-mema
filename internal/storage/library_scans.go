package storage

import "github.com/jackc/pgx/v5"

func scanLibraryFolder(row pgx.Row) (LibraryFolder, error) {
	var folder LibraryFolder
	err := row.Scan(&folder.ID, &folder.Path, &folder.CreatedAt, &folder.UpdatedAt)
	return folder, err
}

func scanLibraryScan(row pgx.Row) (LibraryScan, error) {
	var scan LibraryScan
	err := row.Scan(
		&scan.ID,
		&scan.FolderID,
		&scan.FolderPath,
		&scan.Status,
		&scan.TotalFiles,
		&scan.AutoMatchedCount,
		&scan.ManualCount,
		&scan.CreatedAt,
		&scan.CompletedAt,
	)
	return scan, err
}

func scanLibraryScanItem(row pgx.Row) (LibraryScanItem, error) {
	var item LibraryScanItem
	err := row.Scan(
		&item.ID,
		&item.ScanID,
		&item.Path,
		&item.FileName,
		&item.DetectedTitle,
		&item.DetectedYear,
		&item.DetectedMediaKind,
		&item.Status,
		&item.MatchedTitle,
		&item.MatchedYear,
		&item.MatchedMediaKind,
		&item.MediaItemID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}
