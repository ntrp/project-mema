package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	storagegen "media-manager/internal/storage/generated"
)

func libraryFolderFromRow(row storagegen.AppLibraryFolder) LibraryFolder {
	return LibraryFolder{
		ID:        row.ID,
		Path:      row.Path,
		Kind:      row.Kind,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func libraryScanFromRow(row storagegen.GetLibraryScanRow) LibraryScan {
	return LibraryScan{
		ID:               row.ID,
		FolderID:         row.LibraryFolderID,
		FolderPath:       row.FolderPath,
		FolderKind:       row.FolderKind,
		Status:           row.Status,
		TotalFiles:       row.TotalFiles,
		AutoMatchedCount: row.AutoMatchedCount,
		ManualCount:      row.ManualCount,
		CreatedAt:        row.CreatedAt,
		CompletedAt:      row.CompletedAt,
	}
}

func libraryScanItemFromRow(row storagegen.AppLibraryScanItem) LibraryScanItem {
	return libraryScanItemFromFields(libraryScanItemFields{
		ID: row.ID, ScanID: row.ScanID, Path: row.Path, FileName: row.FileName,
		SizeBytes: row.SizeBytes, DetectedTitle: row.DetectedTitle,
		DetectedYear: row.DetectedYear, DetectedMediaKind: row.DetectedMediaKind,
		SeasonNumber: row.SeasonNumber, EpisodeNumber: row.EpisodeNumber,
		Status: row.Status, Imported: row.Imported, MatchedTitle: row.MatchedTitle,
		MatchedYear: row.MatchedYear, MatchedMediaKind: row.MatchedMediaKind,
		MatchedExternalProvider: row.MatchedExternalProvider,
		MatchedExternalID:       row.MatchedExternalID, MatchSource: row.MatchSource,
		SelectedMetadataProviderID: row.SelectedMetadataProviderID,
		DuplicateGroupID:           row.DuplicateGroupID,
		DuplicateRemovalAllowed:    row.DuplicateRemovalAllowed,
		MediaItemID:                row.MediaItemID, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	})
}

func libraryScanItemFromListRow(row storagegen.ListLibraryScanItemsRow) LibraryScanItem {
	return libraryScanItemFromFields(libraryScanItemFields{
		ID: row.ID, ScanID: row.ScanID, Path: row.Path, FileName: row.FileName,
		SizeBytes: row.SizeBytes, DetectedTitle: row.DetectedTitle,
		DetectedYear: row.DetectedYear, DetectedMediaKind: row.DetectedMediaKind,
		SeasonNumber: row.SeasonNumber, EpisodeNumber: row.EpisodeNumber,
		Status: row.Status, Imported: row.Imported, MatchedTitle: row.MatchedTitle,
		MatchedYear: row.MatchedYear, MatchedMediaKind: row.MatchedMediaKind,
		MatchedExternalProvider: row.MatchedExternalProvider,
		MatchedExternalID:       row.MatchedExternalID, MatchSource: row.MatchSource,
		SelectedMetadataProviderID: row.SelectedMetadataProviderID,
		DuplicateGroupID:           row.DuplicateGroupID,
		DuplicateRemovalAllowed:    row.DuplicateRemovalAllowed,
		MediaItemID:                row.MediaItemID, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	})
}

func libraryScanItemFromMatchRow(row storagegen.MatchLibraryScanItemRow) LibraryScanItem {
	return libraryScanItemFromFields(libraryScanItemFields{
		ID: row.ID, ScanID: row.ScanID, Path: row.Path, FileName: row.FileName,
		SizeBytes: row.SizeBytes, DetectedTitle: row.DetectedTitle,
		DetectedYear: row.DetectedYear, DetectedMediaKind: row.DetectedMediaKind,
		SeasonNumber: row.SeasonNumber, EpisodeNumber: row.EpisodeNumber,
		Status: row.Status, Imported: row.Imported, MatchedTitle: row.MatchedTitle,
		MatchedYear: row.MatchedYear, MatchedMediaKind: row.MatchedMediaKind,
		MatchedExternalProvider: row.MatchedExternalProvider,
		MatchedExternalID:       row.MatchedExternalID, MatchSource: row.MatchSource,
		SelectedMetadataProviderID: row.SelectedMetadataProviderID,
		DuplicateGroupID:           row.DuplicateGroupID,
		DuplicateRemovalAllowed:    row.DuplicateRemovalAllowed,
		MediaItemID:                row.MediaItemID, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	})
}

type libraryScanItemFields struct {
	ID                         uuid.UUID
	ScanID                     uuid.UUID
	Path                       string
	FileName                   string
	SizeBytes                  int64
	DetectedTitle              string
	DetectedYear               pgtype.Int4
	DetectedMediaKind          string
	SeasonNumber               pgtype.Int4
	EpisodeNumber              pgtype.Int4
	Status                     string
	Imported                   bool
	MatchedTitle               pgtype.Text
	MatchedYear                pgtype.Int4
	MatchedMediaKind           pgtype.Text
	MatchedExternalProvider    pgtype.Text
	MatchedExternalID          pgtype.Text
	MatchSource                pgtype.Text
	SelectedMetadataProviderID *uuid.UUID
	DuplicateGroupID           pgtype.Text
	DuplicateRemovalAllowed    bool
	MediaItemID                *uuid.UUID
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

func libraryScanItemFromFields(row libraryScanItemFields) LibraryScanItem {
	return LibraryScanItem{
		ID:                         row.ID,
		ScanID:                     row.ScanID,
		Path:                       row.Path,
		FileName:                   row.FileName,
		SizeBytes:                  row.SizeBytes,
		DetectedTitle:              row.DetectedTitle,
		DetectedYear:               int4Ptr(row.DetectedYear),
		DetectedMediaKind:          row.DetectedMediaKind,
		SeasonNumber:               int4Ptr(row.SeasonNumber),
		EpisodeNumber:              int4Ptr(row.EpisodeNumber),
		Status:                     row.Status,
		Imported:                   row.Imported,
		MatchedTitle:               textPtr(row.MatchedTitle),
		MatchedYear:                int4Ptr(row.MatchedYear),
		MatchedMediaKind:           textPtr(row.MatchedMediaKind),
		MatchedExternalProvider:    textPtr(row.MatchedExternalProvider),
		MatchedExternalID:          textPtr(row.MatchedExternalID),
		MatchSource:                textPtr(row.MatchSource),
		SelectedMetadataProviderID: row.SelectedMetadataProviderID,
		DuplicateGroupID:           textPtr(row.DuplicateGroupID),
		DuplicateRemovalAllowed:    row.DuplicateRemovalAllowed,
		MediaItemID:                row.MediaItemID,
		CreatedAt:                  row.CreatedAt,
		UpdatedAt:                  row.UpdatedAt,
	}
}
