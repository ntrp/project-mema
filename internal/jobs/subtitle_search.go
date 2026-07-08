package jobs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

type SubtitleSearchArgs struct {
	MediaItemID string `json:"media_item_id" river:"unique"`
	LanguageID  string `json:"language_id" river:"unique"`
	FilePath    string `json:"file_path,omitempty" river:"unique"`
}

func (SubtitleSearchArgs) Kind() string {
	return "media.subtitle_search"
}

type SubtitleRetryArgs struct{}

func (SubtitleRetryArgs) Kind() string {
	return "media.subtitle_retry"
}

type SubtitleSearchWorker struct {
	river.WorkerDefaults[SubtitleSearchArgs]

	settings  *storage.SettingsStore
	subtitles *subtitles.Service
	events    *events.Broker
}

func (w *SubtitleSearchWorker) Work(ctx context.Context, job *river.Job[SubtitleSearchArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()
	recordJobProgress(ctx, w.settings, w.events, nil, "Searching subtitles")
	itemID, err := uuid.Parse(job.Args.MediaItemID)
	if err != nil {
		return fmt.Errorf("parse media item id: %w", err)
	}
	item, err := w.settings.GetMediaItem(ctx, itemID)
	if err != nil {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "subtitles", "Subtitle search failed to load media", map[string]any{"mediaItemId": itemID.String(), "error": err.Error()})
		return fmt.Errorf("load media item: %w", err)
	}
	return subtitleSearchDownload(ctx, w.settings, w.subtitles, w.events, item, job.Args)
}

type SubtitleRetryWorker struct {
	river.WorkerDefaults[SubtitleRetryArgs]

	settings  *storage.SettingsStore
	subtitles *subtitles.Service
	events    *events.Broker
}

func (w *SubtitleRetryWorker) Work(ctx context.Context, job *river.Job[SubtitleRetryArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()
	recordJobProgress(ctx, w.settings, w.events, nil, "Retrying missing subtitles")
	items, err := w.settings.ListMediaItems(ctx)
	if err != nil {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "subtitles", "Subtitle retry failed to list media", map[string]any{"error": err.Error()})
		return fmt.Errorf("list media: %w", err)
	}
	queued := 0
	failures := []string{}
	for _, item := range items {
		for _, args := range subtitleSearchRequestsForItem(item) {
			if err := subtitleSearchDownload(ctx, w.settings, w.subtitles, w.events, item, args); err != nil {
				failures = append(failures, fmt.Sprintf("%s: %s", item.Title, err.Error()))
			} else {
				queued++
			}
		}
	}
	if len(failures) > 0 {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "subtitles", "Subtitle retry finished with failures", map[string]any{"failureCount": len(failures), "successCount": queued})
		return fmt.Errorf("subtitle retry failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "subtitles", "Subtitle retry finished", map[string]any{"successCount": queued})
	return nil
}

func subtitleSearchDownload(
	ctx context.Context,
	settings *storage.SettingsStore,
	service *subtitles.Service,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args SubtitleSearchArgs,
) error {
	if strings.TrimSpace(args.FilePath) != "" {
		path, err := settings.MediaItemFilePath(ctx, item.ID, args.FilePath)
		if err != nil {
			return err
		}
		args.FilePath = path
	}
	request, ok := subtitleSearchRequest(item, args)
	if !ok {
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "subtitles", "Subtitle search skipped", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title})
		return nil
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "subtitles", "Subtitle search started", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "languageId": request.LanguageID})
	candidate, provider, err := bestSubtitleCandidate(ctx, settings, service, request)
	if err != nil {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "subtitles", "Subtitle search failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "languageId": request.LanguageID, "error": err.Error()})
		return err
	}
	download, err := service.Download(ctx, subtitleConfig(provider), candidate)
	if err != nil {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "subtitles", "Subtitle download failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "languageId": request.LanguageID, "error": err.Error()})
		return err
	}
	targetFormat := firstNonEmpty(subtitleTargetFormat(item, request.LanguageID), candidate.Format, "subrip")
	artifact, err := writeSubtitleFile(request, download.Content, targetFormat)
	if err != nil {
		return err
	}
	_, err = settings.UpsertMediaItemSubtitle(ctx, subtitleRecord(item, provider, candidate, request, artifact, download.URL))
	if err != nil {
		return err
	}
	slog.Debug("subtitle downloaded", "mediaItemId", item.ID, "path", artifact.Path, "languageId", request.LanguageID)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "subtitles", "Subtitle downloaded", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "languageId": request.LanguageID, "path": artifact.Path})
	return nil
}

type subtitleArtifact struct {
	Path      string
	Format    string
	Checksum  string
	SizeBytes int64
}

func writeSubtitleFile(request subtitles.SearchRequest, content []byte, targetFormat string) (subtitleArtifact, error) {
	converted, format, err := convertSubtitleContent(content, targetFormat)
	if err != nil {
		return subtitleArtifact{}, err
	}
	target := subtitleSidecarPath(request.FilePath, request.LanguageID, format)
	artifact := subtitleArtifact{
		Path:      target,
		Format:    format,
		Checksum:  subtitleChecksum(converted),
		SizeBytes: int64(len(converted)),
	}
	if _, err := os.Stat(target); err == nil {
		return artifact, nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return subtitleArtifact{}, err
	}
	return artifact, os.WriteFile(target, converted, 0o644)
}

func subtitleChecksum(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}
