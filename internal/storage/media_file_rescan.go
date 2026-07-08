package storage

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

var mediaFileExtensions = map[string]struct{}{
	".avi":  {},
	".m4v":  {},
	".mkv":  {},
	".mov":  {},
	".mp4":  {},
	".mpeg": {},
	".mpg":  {},
	".ts":   {},
	".webm": {},
	".wmv":  {},
}

func (s *SettingsStore) RescanMediaItemFiles(ctx context.Context, id uuid.UUID) (MediaItem, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaItem{}, err
	}
	if item.LibraryFolderID == nil || item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return MediaItem{}, ErrInvalidInput
	}
	root := filepath.Clean(strings.TrimSpace(*item.MediaFolderPath))
	files, err := mediaFilesInRoot(root)
	if err != nil {
		return MediaItem{}, ErrInvalidInput
	}
	kind, err := mediaItemKind(item.Type)
	if err != nil {
		return MediaItem{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := storagegen.New(s.pool).WithTx(tx)
	mediaItemID := item.ID
	existing, err := queries.ListMediaFileRecordsForItem(ctx, &mediaItemID)
	if err != nil {
		return MediaItem{}, err
	}
	scanID := uuid.New()
	if err := queries.CreateMediaFileRescanLibraryScan(ctx, storagegen.CreateMediaFileRescanLibraryScanParams{
		ID:              scanID,
		LibraryFolderID: *item.LibraryFolderID,
		TotalFiles:      int32(len(files)),
	}); err != nil {
		return MediaItem{}, err
	}
	for index := range existing {
		existing[index].Path = absoluteCleanPathOrClean(existing[index].Path)
		existing[index].FileName = filepath.Base(existing[index].Path)
	}
	currentPaths := mediaFilePathSet(files)
	knownPaths := map[string]struct{}{}
	usedMovedSources := map[uuid.UUID]struct{}{}
	for _, record := range existing {
		if _, ok := currentPaths[record.Path]; ok {
			knownPaths[record.Path] = struct{}{}
			if missingMediaFileStatus(record.Status) {
				if err := updateMediaFileRecordStatus(ctx, queries, record, "restored"); err != nil {
					return MediaItem{}, err
				}
				if err := recordMediaFileReconciliation(ctx, tx, mediaItemID, record.Path, "restored", "succeeded", nil, nil); err != nil {
					return MediaItem{}, err
				}
			}
			if err := recordMediaFileFactFromPath(ctx, tx, mediaItemID, record.SeasonID, record.EpisodeID, record.Path, "rescan"); err != nil {
				return MediaItem{}, err
			}
			if err := recordImportedFileSidecars(ctx, tx, mediaItemID, record.Path, record.SeasonID, record.EpisodeID, item.SubtitleMode); err != nil {
				return MediaItem{}, err
			}
			continue
		}
		if activeMediaFileStatus(record.Status) {
			if err := updateMediaFileRecordStatus(ctx, queries, record, "missing"); err != nil {
				return MediaItem{}, err
			}
			if err := recordMediaFileReconciliation(ctx, tx, mediaItemID, record.Path, "missing", "succeeded", nil, nil); err != nil {
				return MediaItem{}, err
			}
		}
	}
	for _, path := range files {
		if _, ok := knownPaths[path]; ok {
			continue
		}
		status := "auto_added"
		var sourcePath *string
		if source := movedSourceCandidate(existing, path, currentPaths, usedMovedSources); source != nil {
			status = "moved_candidate"
			sourcePath = &source.Path
			usedMovedSources[source.ID] = struct{}{}
		}
		if err := queries.CreateMediaFileRescanLibraryScanItem(ctx, storagegen.CreateMediaFileRescanLibraryScanItemParams{
			ID:                uuid.New(),
			ScanID:            scanID,
			Path:              path,
			FileName:          filepath.Base(path),
			DetectedTitle:     item.Title,
			DetectedYear:      int4Value(item.Year),
			DetectedMediaKind: kind,
			Status:            status,
			MediaItemID:       &mediaItemID,
			SeasonID:          nil,
			EpisodeID:         nil,
		}); err != nil {
			return MediaItem{}, err
		}
		if status == "moved_candidate" {
			if err := recordMediaFileReconciliation(ctx, tx, mediaItemID, path, "moved_candidate", "skipped", sourcePath, &path); err != nil {
				return MediaItem{}, err
			}
		}
		if err := recordMediaFileFactFromPath(ctx, tx, mediaItemID, nil, nil, path, "rescan"); err != nil {
			return MediaItem{}, err
		}
		if err := recordImportedFileSidecars(ctx, tx, mediaItemID, path, nil, nil, item.SubtitleMode); err != nil {
			return MediaItem{}, err
		}
	}
	if err := queries.TouchMediaItem(ctx, item.ID); err != nil {
		return MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, id)
}

func mediaFilePathSet(paths []string) map[string]struct{} {
	set := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		set[path] = struct{}{}
	}
	return set
}

func activeMediaFileStatus(status string) bool {
	return status == "auto_added" || status == "manually_added" || status == "restored"
}

func missingMediaFileStatus(status string) bool {
	return status == "missing" || status == "moved_candidate"
}

func updateMediaFileRecordStatus(ctx context.Context, q *storagegen.Queries, record storagegen.AppLibraryScanItem, status string) error {
	return q.UpdateLibraryScanItemStatus(ctx, storagegen.UpdateLibraryScanItemStatusParams{
		ID:     record.ID,
		Status: status,
	})
}

func movedSourceCandidate(
	records []storagegen.AppLibraryScanItem,
	path string,
	currentPaths map[string]struct{},
	used map[uuid.UUID]struct{},
) *storagegen.AppLibraryScanItem {
	name := filepath.Base(path)
	for index := range records {
		record := &records[index]
		if record.FileName != name || record.Path == path {
			continue
		}
		if _, exists := currentPaths[record.Path]; exists {
			continue
		}
		if record.Status != "missing" && !activeMediaFileStatus(record.Status) {
			continue
		}
		if _, ok := used[record.ID]; ok {
			continue
		}
		return record
	}
	return nil
}

func recordMediaFileReconciliation(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	filePath string,
	operation string,
	status string,
	sourcePath *string,
	destinationPath *string,
) error {
	_, err := createMediaFileHistory(ctx, q, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        filePath,
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		Operation:       operation,
		Status:          status,
		ActorType:       "system",
	})
	return err
}

func mediaFilesInRoot(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrInvalidInput
	}
	files := []string{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if _, ok := mediaFileExtensions[strings.ToLower(filepath.Ext(entry.Name()))]; !ok {
			return nil
		}
		files = append(files, path)
		return nil
	})
	sort.Strings(files)
	return files, err
}
