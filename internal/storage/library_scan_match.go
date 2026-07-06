package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) MatchLibraryScanItem(ctx context.Context, scanID uuid.UUID, itemID uuid.UUID, input LibraryMatchInput) (LibraryScanItem, MediaItem, error) {
	return s.matchLibraryScanItem(ctx, scanID, itemID, input, false)
}

func (s *SettingsStore) ImportLibraryScanItem(ctx context.Context, scanID uuid.UUID, itemID uuid.UUID, input LibraryMatchInput) (LibraryScanItem, MediaItem, error) {
	return s.matchLibraryScanItem(ctx, scanID, itemID, input, true)
}

func (s *SettingsStore) matchLibraryScanItem(ctx context.Context, scanID uuid.UUID, itemID uuid.UUID, input LibraryMatchInput, imported bool) (LibraryScanItem, MediaItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mediaType, contentKind, ok := mediaKindToMediaTypeAndContent(input.MediaKind)
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
	item, err := matchedLibraryMediaItem(ctx, tx, folderID, mediaType, contentKind, input)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	itemPath, err := storagegen.New(tx).GetLibraryScanItemPath(ctx, storagegen.GetLibraryScanItemPathParams{
		ScanID: scanID,
		ID:     itemID,
	})
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	seasonID, episodeID, err := importedEpisodeReference(ctx, tx, item, itemPath)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	matchSource := libraryMatchSource(input)
	row, err := storagegen.New(tx).MatchLibraryScanItem(ctx, storagegen.MatchLibraryScanItemParams{
		MatchedTitle:               textValue(&input.Title),
		MatchedYear:                int4Value(input.Year),
		MatchedMediaKind:           textValue(&input.MediaKind),
		MatchedExternalProvider:    textValue(input.ExternalProvider),
		MatchedExternalID:          textValue(input.ExternalID),
		MatchSource:                textValue(&matchSource),
		Imported:                   imported,
		SelectedMetadataProviderID: input.MetadataProviderID,
		MediaItemID:                &item.ID,
		SeasonID:                   seasonID,
		EpisodeID:                  episodeID,
		ScanID:                     scanID,
		ID:                         itemID,
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
	if imported {
		if err := recordImportedFileProvenance(ctx, tx, item.ID, "", itemPath, "libraryScanImport"); err != nil {
			return LibraryScanItem{}, MediaItem{}, err
		}
		if err := recordImportedFileSidecars(ctx, tx, item.ID, itemPath, seasonID, episodeID, item.SubtitlePreferredMode); err != nil {
			return LibraryScanItem{}, MediaItem{}, err
		}
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

func matchedLibraryMediaItem(
	ctx context.Context,
	q mediaItemQuerier,
	folderID uuid.UUID,
	mediaType string,
	contentKind string,
	input LibraryMatchInput,
) (MediaItem, error) {
	if input.MediaItemID != nil {
		item, err := getMediaItem(ctx, q, *input.MediaItemID)
		if err != nil {
			return MediaItem{}, err
		}
		if item.Type != mediaType {
			return MediaItem{}, ErrInvalidInput
		}
		return item, nil
	}
	return createMediaItemIfMissing(ctx, q, MediaItemInput{
		Type:                  mediaType,
		ContentKind:           contentKind,
		Title:                 input.Title,
		Year:                  input.Year,
		Monitored:             input.Monitored,
		ExternalProvider:      input.ExternalProvider,
		ExternalID:            input.ExternalID,
		Overview:              input.Overview,
		PosterPath:            input.PosterPath,
		MediaMetadataSnapshot: input.MediaMetadataSnapshot,
		MonitorMode:           input.MonitorMode,
		SeriesType:            input.SeriesType,
		MinimumAvailability:   input.MinimumAvailability,
		QualityProfileID:      &input.QualityProfileID,
		LibraryFolderID:       &folderID,
	})
}

func libraryMatchSource(input LibraryMatchInput) string {
	if input.MediaItemID != nil {
		return "library"
	}
	if input.ExternalProvider != nil || input.ExternalID != nil {
		return "provider"
	}
	return "manual"
}
