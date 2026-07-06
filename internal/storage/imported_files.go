package storage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

var importedEpisodePattern = regexp.MustCompile(`(?i)s(\d{1,2})e(\d{1,3})`)

func (s *SettingsStore) RecordImportedMediaFile(ctx context.Context, item MediaItem, filePath string) error {
	return s.RecordImportedMediaFileWithHistory(ctx, item, "", filePath, "")
}

func (s *SettingsStore) RecordImportedMediaFileWithHistory(
	ctx context.Context,
	item MediaItem,
	sourcePath string,
	filePath string,
	importMode string,
) error {
	if item.LibraryFolderID == nil {
		return ErrNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := storagegen.New(s.pool).WithTx(tx)
	mediaItemID := item.ID
	exists, err := queries.ImportedMediaFileExists(ctx, storagegen.ImportedMediaFileExistsParams{
		MediaItemID: &mediaItemID,
		Path:        filePath,
	})
	if err != nil {
		return err
	}
	if exists {
		if err := recordImportedFileProvenance(ctx, tx, mediaItemID, sourcePath, filePath, "diskImport"); err != nil {
			return err
		}
		seasonID, episodeID, err := importedFileEpisodeReference(ctx, queries, mediaItemID, filePath)
		if err != nil {
			return err
		}
		if err := recordImportedFileSidecars(ctx, tx, mediaItemID, filePath, seasonID, episodeID, item.SubtitlePreferredMode); err != nil {
			return err
		}
		return tx.Commit(ctx)
	}

	scanID := uuid.New()
	if err := queries.CreateImportedFileLibraryScan(ctx, storagegen.CreateImportedFileLibraryScanParams{
		ID:              scanID,
		LibraryFolderID: *item.LibraryFolderID,
	}); err != nil {
		return err
	}

	kind, err := mediaItemKind(item.Type)
	if err != nil {
		return err
	}
	seasonID, episodeID, err := importedEpisodeReference(ctx, tx, item, filePath)
	if err != nil {
		return err
	}
	if err := queries.CreateImportedFileLibraryScanItem(ctx, storagegen.CreateImportedFileLibraryScanItemParams{
		ID:                uuid.New(),
		ScanID:            scanID,
		Path:              filePath,
		FileName:          filepath.Base(filePath),
		DetectedTitle:     item.Title,
		DetectedYear:      int4Value(item.Year),
		DetectedMediaKind: kind,
		MediaItemID:       &mediaItemID,
		SeasonID:          seasonID,
		EpisodeID:         episodeID,
	}); err != nil {
		return err
	}
	if _, err := createMediaFileHistory(ctx, tx, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        filePath,
		SourcePath:      optionalHistoryString(sourcePath),
		DestinationPath: optionalHistoryString(filePath),
		Operation:       "imported",
		Status:          "succeeded",
		ActorType:       "system",
		Details:         map[string]any{"importMode": importMode},
	}); err != nil {
		return err
	}
	if err := recordImportedFileProvenance(ctx, tx, mediaItemID, sourcePath, filePath, "diskImport"); err != nil {
		return err
	}
	if err := recordImportedFileSidecars(ctx, tx, mediaItemID, filePath, seasonID, episodeID, item.SubtitlePreferredMode); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func importedFileEpisodeReference(
	ctx context.Context,
	q *storagegen.Queries,
	mediaItemID uuid.UUID,
	filePath string,
) (*uuid.UUID, *uuid.UUID, error) {
	row, err := q.GetImportedFileEpisodeReference(ctx, storagegen.GetImportedFileEpisodeReferenceParams{
		MediaItemID: &mediaItemID,
		Path:        filePath,
	})
	if err != nil {
		return nil, nil, err
	}
	return row.SeasonID, row.EpisodeID, nil
}

func (s *SettingsStore) ImportedMediaFileEpisodeID(ctx context.Context, mediaItemID uuid.UUID, filePath string) (*uuid.UUID, error) {
	row, err := storagegen.New(s.pool).GetImportedFileEpisodeReference(ctx, storagegen.GetImportedFileEpisodeReferenceParams{
		MediaItemID: &mediaItemID,
		Path:        filePath,
	})
	if err != nil {
		return nil, err
	}
	return row.EpisodeID, nil
}

func importedEpisodeReference(
	ctx context.Context,
	q mediaItemQuerier,
	item MediaItem,
	filePath string,
) (*uuid.UUID, *uuid.UUID, error) {
	if item.Type != "serie" {
		return nil, nil, nil
	}
	seasons, err := listMediaSeriesSeasons(ctx, q, item.ID)
	if err != nil || !mediaSeriesHasEpisodes(seasons) {
		return nil, nil, err
	}
	seasonNumber, episodeNumber, ok := importedEpisodeNumbers(filePath)
	if !ok {
		return nil, nil, fmt.Errorf("episode import target could not be resolved from %s", filepath.Base(filePath))
	}
	for _, season := range seasons {
		if season.SeasonNumber != seasonNumber {
			continue
		}
		for _, episode := range season.Episodes {
			if episode.EpisodeNumber == episodeNumber {
				return &season.ID, &episode.ID, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("episode import target S%02dE%02d is not known for %s", seasonNumber, episodeNumber, item.Title)
}

func mediaSeriesHasEpisodes(seasons []MediaSeriesSeason) bool {
	for _, season := range seasons {
		if len(season.Episodes) > 0 {
			return true
		}
	}
	return false
}

func importedEpisodeNumbers(filePath string) (int32, int32, bool) {
	matches := importedEpisodePattern.FindStringSubmatch(filepath.Base(filePath))
	if len(matches) != 3 {
		return 0, 0, false
	}
	season, seasonErr := strconv.ParseInt(matches[1], 10, 32)
	episode, episodeErr := strconv.ParseInt(matches[2], 10, 32)
	if seasonErr != nil || episodeErr != nil {
		return 0, 0, false
	}
	return int32(season), int32(episode), true
}

func mediaItemKind(mediaType string) (string, error) {
	switch mediaType {
	case "movie":
		return "movie", nil
	case "serie":
		return "series", nil
	default:
		return "", errors.New("unsupported media type")
	}
}
