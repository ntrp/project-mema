package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

type LibraryFolder struct {
	ID        uuid.UUID
	Path      string
	Kind      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LibraryScan struct {
	ID               uuid.UUID
	FolderID         uuid.UUID
	FolderPath       string
	FolderKind       string
	Status           string
	TotalFiles       int32
	AutoMatchedCount int32
	ManualCount      int32
	Items            []LibraryScanItem
	CreatedAt        time.Time
	CompletedAt      *time.Time
}

type LibraryScanItem struct {
	ID                         uuid.UUID
	ScanID                     uuid.UUID
	Path                       string
	FileName                   string
	SizeBytes                  int64
	DetectedTitle              string
	DetectedYear               *int32
	DetectedMediaKind          string
	SeasonNumber               *int32
	EpisodeNumber              *int32
	Status                     string
	Imported                   bool
	MatchedTitle               *string
	MatchedYear                *int32
	MatchedMediaKind           *string
	MatchedExternalProvider    *string
	MatchedExternalID          *string
	MatchSource                *string
	SelectedMetadataProviderID *uuid.UUID
	DuplicateGroupID           *string
	DuplicateRemovalAllowed    bool
	MediaItemID                *uuid.UUID
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type LibraryScanItemInput struct {
	Path                       string
	FileName                   string
	SizeBytes                  int64
	DetectedTitle              string
	DetectedYear               *int32
	DetectedMediaKind          string
	SeasonNumber               *int32
	EpisodeNumber              *int32
	SafeMatch                  bool
	Imported                   bool
	MatchedTitle               *string
	MatchedYear                *int32
	MatchedMediaKind           *string
	MatchedExternalProvider    *string
	MatchedExternalID          *string
	MatchSource                *string
	SelectedMetadataProviderID *uuid.UUID
	DuplicateGroupID           *string
	DuplicateRemovalAllowed    bool
	MediaItemID                *uuid.UUID
}

type LibraryMatchInput struct {
	MediaKind           string
	Title               string
	Year                *int32
	Monitored           bool
	QualityProfileID    string
	MonitorMode         string
	MinimumAvailability string
	SeriesType          *string
	MetadataProviderID  *uuid.UUID
	MediaItemID         *uuid.UUID
	ExternalProvider    *string
	ExternalID          *string
	Overview            *string
	PosterPath          *string
	MediaMetadataSnapshot
}

type LibraryScanItemResetResult struct {
	Item               LibraryScanItem
	RemovedMediaItemID *uuid.UUID
}

type ActiveImportedPath struct {
	MediaItemID   uuid.UUID
	MatchedTitle  string
	MatchedYear   *int32
	MatchedSource string
}

func (s *SettingsStore) ListLibraryFolders(ctx context.Context) ([]LibraryFolder, error) {
	rows, err := storagegen.New(s.pool).ListLibraryFolders(ctx)
	if err != nil {
		return nil, err
	}

	folders := make([]LibraryFolder, 0, len(rows))
	for _, row := range rows {
		folders = append(folders, libraryFolderFromRow(row))
	}
	return folders, nil
}

func (s *SettingsStore) CreateLibraryFolder(ctx context.Context, path string, kind string) (LibraryFolder, error) {
	path, err := absoluteCleanPath(path)
	if err != nil {
		return LibraryFolder{}, err
	}
	row, err := storagegen.New(s.pool).UpsertLibraryFolder(ctx, storagegen.UpsertLibraryFolderParams{
		ID:   uuid.New(),
		Path: path,
		Kind: kind,
	})
	if err != nil {
		return LibraryFolder{}, err
	}
	return libraryFolderFromRow(row), nil
}

func (s *SettingsStore) DeleteLibraryFolder(ctx context.Context, id uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(tx)
	if _, err := queries.DeleteMediaItemsForLibraryFolder(ctx, id); err != nil {
		return err
	}
	rowsAffected, err := queries.DeleteLibraryFolder(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return tx.Commit(ctx)
}

func (s *SettingsStore) LibraryFolderExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return storagegen.New(s.pool).LibraryFolderExists(ctx, id)
}

func (s *SettingsStore) CreateLibraryScan(ctx context.Context, folder LibraryFolder, inputs []LibraryScanItemInput) (LibraryScan, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return LibraryScan{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	scanID := uuid.New()
	manualCount := int32(0)
	q := storagegen.New(tx)
	if err := q.CreateLibraryScan(ctx, storagegen.CreateLibraryScanParams{
		ID:              scanID,
		LibraryFolderID: folder.ID,
		TotalFiles:      int32(len(inputs)),
	}); err != nil {
		return LibraryScan{}, err
	}

	for _, input := range inputs {
		manualCount++
		if err := q.AddLibraryScanItem(ctx, storagegen.AddLibraryScanItemParams{
			ID:                         uuid.New(),
			ScanID:                     scanID,
			Path:                       input.Path,
			FileName:                   input.FileName,
			SizeBytes:                  input.SizeBytes,
			DetectedTitle:              input.DetectedTitle,
			DetectedYear:               int4Value(input.DetectedYear),
			DetectedMediaKind:          input.DetectedMediaKind,
			SeasonNumber:               int4Value(input.SeasonNumber),
			EpisodeNumber:              int4Value(input.EpisodeNumber),
			Status:                     "pending",
			Imported:                   input.Imported,
			MatchedTitle:               textValue(input.MatchedTitle),
			MatchedYear:                int4Value(input.MatchedYear),
			MatchedMediaKind:           textValue(input.MatchedMediaKind),
			MatchedExternalProvider:    textValue(input.MatchedExternalProvider),
			MatchedExternalID:          textValue(input.MatchedExternalID),
			MatchSource:                textValue(input.MatchSource),
			SelectedMetadataProviderID: input.SelectedMetadataProviderID,
			DuplicateGroupID:           textValue(input.DuplicateGroupID),
			DuplicateRemovalAllowed:    input.DuplicateRemovalAllowed,
			MediaItemID:                input.MediaItemID,
		}); err != nil {
			return LibraryScan{}, err
		}
	}

	if err := q.UpdateLibraryScanCounts(ctx, storagegen.UpdateLibraryScanCountsParams{
		AutoMatchedCount: 0,
		ManualCount:      manualCount,
		ID:               scanID,
	}); err != nil {
		return LibraryScan{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LibraryScan{}, err
	}
	return s.GetLibraryScan(ctx, scanID)
}

func (s *SettingsStore) GetLibraryScan(ctx context.Context, id uuid.UUID) (LibraryScan, error) {
	row, err := storagegen.New(s.pool).GetLibraryScan(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryScan{}, ErrNotFound
	}
	if err != nil {
		return LibraryScan{}, err
	}
	scan := libraryScanFromRow(row)
	items, err := s.listLibraryScanItems(ctx, id)
	if err != nil {
		return LibraryScan{}, err
	}
	scan.Items = items
	return scan, nil
}

func (s *SettingsStore) ActiveImportedPathsForLibraryFolder(ctx context.Context, folderID uuid.UUID) (map[string]ActiveImportedPath, error) {
	rows, err := storagegen.New(s.pool).ListActiveImportedPathsForLibraryFolder(ctx, folderID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]ActiveImportedPath, len(rows))
	for _, row := range rows {
		result[row.Path] = ActiveImportedPath{
			MediaItemID:   row.MediaItemID,
			MatchedTitle:  row.MatchedTitle,
			MatchedYear:   int4Ptr(row.MatchedYear),
			MatchedSource: "library",
		}
	}
	return result, nil
}

func (s *SettingsStore) listLibraryScanItems(ctx context.Context, scanID uuid.UUID) ([]LibraryScanItem, error) {
	rows, err := storagegen.New(s.pool).ListLibraryScanItems(ctx, scanID)
	if err != nil {
		return nil, err
	}

	items := make([]LibraryScanItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, libraryScanItemFromListRow(row))
	}
	return items, nil
}
