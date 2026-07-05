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
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LibraryScan struct {
	ID               uuid.UUID
	FolderID         uuid.UUID
	FolderPath       string
	Status           string
	TotalFiles       int32
	AutoMatchedCount int32
	ManualCount      int32
	Items            []LibraryScanItem
	CreatedAt        time.Time
	CompletedAt      *time.Time
}

type LibraryScanItem struct {
	ID                uuid.UUID
	ScanID            uuid.UUID
	Path              string
	FileName          string
	DetectedTitle     string
	DetectedYear      *int32
	DetectedMediaKind string
	Status            string
	MatchedTitle      *string
	MatchedYear       *int32
	MatchedMediaKind  *string
	MediaItemID       *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type LibraryScanItemInput struct {
	Path              string
	FileName          string
	DetectedTitle     string
	DetectedYear      *int32
	DetectedMediaKind string
	SafeMatch         bool
}

type LibraryMatchInput struct {
	MediaKind           string
	Title               string
	Year                *int32
	Monitored           bool
	QualityProfileID    string
	MonitorMode         string
	MinimumAvailability string
	ExternalProvider    *string
	ExternalID          *string
	Overview            *string
	PosterPath          *string
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

func (s *SettingsStore) CreateLibraryFolder(ctx context.Context, path string) (LibraryFolder, error) {
	row, err := storagegen.New(s.pool).UpsertLibraryFolder(ctx, storagegen.UpsertLibraryFolderParams{
		ID:   uuid.New(),
		Path: path,
	})
	if err != nil {
		return LibraryFolder{}, err
	}
	return libraryFolderFromRow(row), nil
}

func (s *SettingsStore) DeleteLibraryFolder(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := storagegen.New(s.pool).DeleteLibraryFolder(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
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
			ID:                uuid.New(),
			ScanID:            scanID,
			Path:              input.Path,
			FileName:          input.FileName,
			DetectedTitle:     input.DetectedTitle,
			DetectedYear:      int4Value(input.DetectedYear),
			DetectedMediaKind: input.DetectedMediaKind,
			Status:            "pending",
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

func (s *SettingsStore) MatchLibraryScanItem(ctx context.Context, scanID uuid.UUID, itemID uuid.UUID, input LibraryMatchInput) (LibraryScanItem, MediaItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mediaType, ok := mediaKindToMediaType(input.MediaKind)
	if !ok {
		return LibraryScanItem{}, MediaItem{}, ErrNotFound
	}
	folderID, err := storagegen.New(tx).GetLibraryScanFolderID(ctx, scanID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LibraryScanItem{}, MediaItem{}, ErrNotFound
		}
		return LibraryScanItem{}, MediaItem{}, err
	}
	item, err := createMediaItemIfMissing(ctx, tx, MediaItemInput{
		Type:                mediaType,
		Title:               input.Title,
		Year:                input.Year,
		Monitored:           input.Monitored,
		ExternalProvider:    input.ExternalProvider,
		ExternalID:          input.ExternalID,
		Overview:            input.Overview,
		PosterPath:          input.PosterPath,
		MonitorMode:         input.MonitorMode,
		MinimumAvailability: input.MinimumAvailability,
		QualityProfileID:    &input.QualityProfileID,
		LibraryFolderID:     &folderID,
	})
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	row, err := storagegen.New(tx).MatchLibraryScanItem(ctx, storagegen.MatchLibraryScanItemParams{
		MatchedTitle:     textValue(&input.Title),
		MatchedYear:      int4Value(input.Year),
		MatchedMediaKind: textValue(&input.MediaKind),
		MediaItemID:      &item.ID,
		ScanID:           scanID,
		ID:               itemID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryScanItem{}, MediaItem{}, ErrNotFound
	}
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	updated := libraryScanItemFromMatchRow(row)
	if err := storagegen.New(tx).RefreshLibraryScanManualCount(ctx, scanID); err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	item, err = getMediaItem(ctx, tx, item.ID)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	return updated, item, nil
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
