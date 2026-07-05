package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (s *SettingsStore) PreviewMediaItemRename(ctx context.Context, id uuid.UUID) (MediaRenamePreview, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaRenamePreview{}, err
	}
	settings, err := s.GetFileNamingSettings(ctx)
	if err != nil {
		return MediaRenamePreview{}, err
	}
	rows := make([]MediaRenamePreviewRow, 0, len(item.FilePaths))
	for _, path := range item.FilePaths {
		rows = append(rows, mediaRenamePreviewRow(item, settings, path))
	}
	return MediaRenamePreview{Rows: rows}, nil
}

func mediaRenamePreviewRow(
	item MediaItem,
	settings FileNamingSettings,
	currentPath string,
) MediaRenamePreviewRow {
	row := MediaRenamePreviewRow{CurrentPath: currentPath, Messages: []string{}}
	proposed, messages, ok := mediaRenameProposedPath(item, settings, currentPath)
	row.Messages = append(row.Messages, messages...)
	row.ProposedPath = proposed
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		row.Status = "missing"
		row.Messages = append(row.Messages, "Source file is missing.")
		return row
	}
	if !ok {
		row.Status = "blocked"
		return row
	}
	if proposed == currentPath {
		row.Status = "unchanged"
		return row
	}
	if _, err := os.Stat(proposed); err == nil {
		row.Status = "conflict"
		row.Messages = append(row.Messages, "Destination already exists.")
		return row
	}
	row.Status = "safe"
	return row
}

func mediaRenameProposedPath(
	item MediaItem,
	settings FileNamingSettings,
	currentPath string,
) (string, []string, bool) {
	if item.LibraryFolderPath == nil || strings.TrimSpace(*item.LibraryFolderPath) == "" {
		return "", []string{"Library root is missing."}, false
	}
	root := filepath.Clean(strings.TrimSpace(*item.LibraryFolderPath))
	if item.Type == "serie" {
		return seriesRenamePath(item, settings, root, currentPath)
	}
	return movieRenamePath(item, settings, root, currentPath)
}

func movieRenamePath(
	item MediaItem,
	settings FileNamingSettings,
	root string,
	currentPath string,
) (string, []string, bool) {
	folder := filepath.Join(root, sanitizePathSegment(renderMediaTemplate(settings.MovieFolderFormat, mediaItemRenameInput(item))))
	file := sanitizePathSegment(renderMediaTemplate(settings.MovieFileFormat, mediaItemRenameInput(item)))
	return checkedRenamePath(root, folder, file+strings.ToLower(filepath.Ext(currentPath)))
}

func seriesRenamePath(
	item MediaItem,
	settings FileNamingSettings,
	root string,
	currentPath string,
) (string, []string, bool) {
	season, episode, ok := importedEpisodeNumbers(currentPath)
	if !ok {
		return "", []string{"Season and episode could not be detected."}, false
	}
	input := mediaItemRenameInput(item)
	input.Seasons = []MediaSeason{{SeasonNumber: season, Episodes: []MediaEpisode{{EpisodeNumber: episode}}}}
	folder := filepath.Join(
		root,
		sanitizePathSegment(renderMediaTemplate(settings.SeriesFolderFormat, input)),
		sanitizePathSegment(renderSeriesTemplate(settings.SeasonFolderFormat, item, season, episode)),
	)
	file := sanitizePathSegment(renderSeriesTemplate(settings.SeriesEpisodeFormat, item, season, episode))
	return checkedRenamePath(root, folder, file+strings.ToLower(filepath.Ext(currentPath)))
}

func checkedRenamePath(root string, folder string, file string) (string, []string, bool) {
	proposed := filepath.Clean(filepath.Join(folder, file))
	if _, err := safePathUnderRoot(root, proposed, false); err != nil {
		return proposed, []string{"Destination is outside the library root."}, false
	}
	return proposed, nil, true
}

func mediaItemRenameInput(item MediaItem) MediaItemInput {
	return MediaItemInput{Type: item.Type, Title: item.Title, Year: item.Year}
}
