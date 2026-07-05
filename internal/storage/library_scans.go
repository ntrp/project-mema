package storage

import storagegen "media-manager/internal/storage/generated"

func libraryFolderFromRow(row storagegen.AppLibraryFolder) LibraryFolder {
	return LibraryFolder{
		ID:        row.ID,
		Path:      row.Path,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func libraryScanFromRow(row storagegen.GetLibraryScanRow) LibraryScan {
	return LibraryScan{
		ID:               row.ID,
		FolderID:         row.LibraryFolderID,
		FolderPath:       row.FolderPath,
		Status:           row.Status,
		TotalFiles:       row.TotalFiles,
		AutoMatchedCount: row.AutoMatchedCount,
		ManualCount:      row.ManualCount,
		CreatedAt:        row.CreatedAt,
		CompletedAt:      row.CompletedAt,
	}
}

func libraryScanItemFromRow(row storagegen.AppLibraryScanItem) LibraryScanItem {
	return LibraryScanItem{
		ID:                row.ID,
		ScanID:            row.ScanID,
		Path:              row.Path,
		FileName:          row.FileName,
		DetectedTitle:     row.DetectedTitle,
		DetectedYear:      int4Ptr(row.DetectedYear),
		DetectedMediaKind: row.DetectedMediaKind,
		Status:            row.Status,
		MatchedTitle:      textPtr(row.MatchedTitle),
		MatchedYear:       int4Ptr(row.MatchedYear),
		MatchedMediaKind:  textPtr(row.MatchedMediaKind),
		MediaItemID:       row.MediaItemID,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
