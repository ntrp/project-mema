package imports

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

type Service struct {
	settings *storage.SettingsStore
}

type ManualImportInput struct {
	SourcePath     string
	TargetFileName string
	ImportMode     ImportMode
	MovieTitle     string
	Year           *int32
	SeasonNumber   *int32
	EpisodeNumber  *int32
	EpisodeTitle   string
	ReleaseGroup   string
	Edition        string
	Quality        string
	Languages      []string
}

func NewService(settings *storage.SettingsStore) *Service {
	return &Service{settings: settings}
}

func (s *Service) ImportCompletedDownload(ctx context.Context, activity storage.DownloadActivity, files []downloadclients.StatusFile) error {
	slog.Debug("import completed download started", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "reportedFileCount", len(files))
	item, err := s.settings.GetMediaItem(ctx, activity.MediaItemID)
	if err != nil {
		slog.Error("import completed download media item load failed", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "error", err)
		return fmt.Errorf("load media item: %w", err)
	}
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		slog.Error("import completed download missing media folder", "activityId", activity.ID, "mediaItemId", item.ID)
		return fmt.Errorf("media folder is not set")
	}
	mappings, err := s.settings.ListPathMappings(ctx)
	if err != nil {
		slog.Error("import completed download path mapping load failed", "activityId", activity.ID, "error", err)
		return fmt.Errorf("load path mappings: %w", err)
	}
	selection, err := selectCompletedDownloadCandidates(files, mappings)
	if err != nil {
		slog.Error("import completed download source discovery failed", "activityId", activity.ID, "error", err)
		return err
	}
	if len(selection.SelectedSources) == 0 {
		slog.Error("import completed download had no valid video candidates", "activityId", activity.ID, "reportedFileCount", len(files), "rejectedCandidates", selection.RejectedCandidates)
		return fmt.Errorf("download client did not report valid import candidates%s", rejectedCandidateSummary(selection.RejectedCandidates))
	}

	if err := os.MkdirAll(*item.MediaFolderPath, 0o755); err != nil {
		slog.Error("import completed download media folder create failed", "activityId", activity.ID, "mediaFolderPath", *item.MediaFolderPath, "error", err)
		return fmt.Errorf("create media folder: %w", err)
	}
	for _, source := range selection.SelectedSources {
		target := filepath.Join(*item.MediaFolderPath, filepath.Base(source))
		slog.Debug("linking completed download file", "activityId", activity.ID, "source", source, "target", target)
		if err := importFile(source, target, ImportModeHardlink); err != nil {
			slog.Error("completed download file link failed", "activityId", activity.ID, "source", source, "target", target, "error", err)
			return err
		}
		if err := s.settings.RecordImportedMediaFile(ctx, item, target); err != nil {
			slog.Error("completed download imported file record failed", "activityId", activity.ID, "target", target, "error", err)
			return fmt.Errorf("record imported file: %w", err)
		}
	}
	slog.Debug("import completed download finished", "activityId", activity.ID, "mediaItemId", item.ID, "linkedFileCount", len(selection.SelectedSources), "rejectedFileCount", len(selection.RejectedCandidates))
	return nil
}

func (s *Service) ImportManualDownload(ctx context.Context, activity storage.DownloadActivity, input ManualImportInput) error {
	item, err := s.settings.GetMediaItem(ctx, activity.MediaItemID)
	if err != nil {
		return fmt.Errorf("load media item: %w", err)
	}
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return fmt.Errorf("media folder is not set")
	}
	mappings, err := s.settings.ListPathMappings(ctx)
	if err != nil {
		return fmt.Errorf("load path mappings: %w", err)
	}
	source := mapPath(input.SourcePath, mappings)
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("source file is not visible to the app: %s", source)
	}
	if info.IsDir() {
		return fmt.Errorf("source path must be a file")
	}
	if !isVideoFile(source) {
		return fmt.Errorf("source path is not a video file")
	}
	targetName, err := manualTargetFileName(item, input, source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*item.MediaFolderPath, 0o755); err != nil {
		return fmt.Errorf("create media folder: %w", err)
	}
	target := filepath.Join(*item.MediaFolderPath, targetName)
	if err := importFile(source, target, input.ImportMode); err != nil {
		return err
	}
	if err := s.settings.RecordImportedMediaFile(ctx, item, target); err != nil {
		return fmt.Errorf("record imported file: %w", err)
	}
	return nil
}

func mapPath(source string, mappings []storage.PathMapping) string {
	source = filepath.Clean(source)
	sortPathMappings(mappings)
	for _, mapping := range mappings {
		clientPath := filepath.Clean(mapping.ClientPath)
		if source == clientPath || strings.HasPrefix(source, clientPath+string(os.PathSeparator)) {
			relative := strings.TrimPrefix(source, clientPath)
			relative = strings.TrimLeft(relative, string(os.PathSeparator))
			return filepath.Join(mapping.AppPath, relative)
		}
	}
	return source
}

func sortPathMappings(mappings []storage.PathMapping) {
	sort.SliceStable(mappings, func(i, j int) bool {
		return len(mappings[i].ClientPath) > len(mappings[j].ClientPath)
	})
}

func manualTargetFileName(item storage.MediaItem, input ManualImportInput, source string) (string, error) {
	if name := cleanFileName(input.TargetFileName); name != "" {
		if filepath.Ext(name) == "" {
			name += strings.ToLower(filepath.Ext(source))
		}
		return name, nil
	}
	ext := strings.ToLower(filepath.Ext(source))
	if ext == "" {
		ext = ".mkv"
	}
	title := strings.TrimSpace(input.MovieTitle)
	if title == "" {
		title = item.Title
	}
	parts := []string{}
	if item.Type == "serie" {
		if input.SeasonNumber == nil || input.EpisodeNumber == nil {
			return "", fmt.Errorf("season and episode are required for series imports")
		}
		parts = append(parts, title, fmt.Sprintf("S%02dE%02d", *input.SeasonNumber, *input.EpisodeNumber))
		if episodeTitle := strings.TrimSpace(input.EpisodeTitle); episodeTitle != "" {
			parts = append(parts, episodeTitle)
		}
	} else {
		year := input.Year
		if year == nil {
			year = item.Year
		}
		if year != nil {
			title = fmt.Sprintf("%s (%d)", title, *year)
		}
		parts = append(parts, title)
	}
	parts = appendNonEmpty(parts, input.Edition, strings.Join(cleanStrings(input.Languages), " "), input.Quality)
	name := cleanFileName(strings.Join(parts, " - "))
	if group := cleanFileName(input.ReleaseGroup); group != "" {
		name = strings.TrimSpace(name + " - " + group)
	}
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	}
	return name + ext, nil
}

func appendNonEmpty(values []string, candidates ...string) []string {
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" {
			values = append(values, candidate)
		}
	}
	return values
}

func cleanStrings(values []string) []string {
	cleaned := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}

func cleanFileName(value string) string {
	value = strings.TrimSpace(value)
	replacer := strings.NewReplacer("/", " ", "\\", " ", ":", " ", "*", " ", "?", " ", "\"", "", "<", " ", ">", " ", "|", " ")
	return strings.Join(strings.Fields(replacer.Replace(value)), " ")
}

func isVideoFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".avi", ".m4v", ".mkv", ".mov", ".mp4", ".mpeg", ".mpg", ".ts", ".webm", ".wmv":
		return true
	default:
		return false
	}
}
