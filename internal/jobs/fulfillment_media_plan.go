package jobs

import (
	"context"
	"fmt"

	"media-manager/internal/events"
	"media-manager/internal/mediafacts"
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func enqueueMediaFulfillmentJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
) error {
	if settings == nil {
		return fmt.Errorf("settings store is not configured")
	}
	if enqueue == nil {
		return fmt.Errorf("fulfillment enqueue function is not configured")
	}
	recordJobProgressData(ctx, settings, eventBroker, nil, "Scanning media fulfillment targets", nil)
	items, err := settings.ListMediaItems(ctx)
	if err != nil {
		return err
	}
	total := len(items)
	recordMediaFulfillmentProgress(ctx, settings, eventBroker, 0, total, 0)
	queued := 0
	for index, item := range items {
		liveItem := mediafacts.WithLiveFileFacts(item, "")
		added, err := enqueueVideoTranscodeTrackJobs(ctx, settings, eventBroker, enqueue, liveItem, FulfillmentActionArgs{})
		if err != nil {
			return err
		}
		queued += added
		added, err = enqueueAudioTranscodeTrackJobs(ctx, settings, eventBroker, enqueue, liveItem, FulfillmentActionArgs{})
		if err != nil {
			return err
		}
		queued += added
		added, err = enqueueContainerRemuxFileJobs(ctx, settings, eventBroker, enqueue, liveItem, FulfillmentActionArgs{})
		if err != nil {
			return err
		}
		queued += added
		added, err = enqueueSubtitleFulfillmentJobs(ctx, settings, eventBroker, enqueue, liveItem)
		if err != nil {
			return err
		}
		queued += added
		recordMediaFulfillmentProgress(ctx, settings, eventBroker, index+1, total, queued)
	}
	done := int32(100)
	recordJobProgressData(ctx, settings, eventBroker, &done, "Media fulfillment scan complete", map[string]any{
		"queuedJobCount":          queued,
		"mediaItemCount":          total,
		"processedMediaItemCount": total,
	})
	return nil
}

func recordMediaFulfillmentProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	processed int,
	total int,
	queued int,
) {
	percent := mediaFulfillmentProgressPercent(processed, total)
	recordJobProgressData(ctx, settings, eventBroker, &percent, fmt.Sprintf("Scanning media %d/%d", processed, total), map[string]any{
		"processedMediaItemCount": processed,
		"mediaItemCount":          total,
		"queuedJobCount":          queued,
	})
}

func mediaFulfillmentProgressPercent(processed int, total int) int32 {
	if total <= 0 {
		return 100
	}
	if processed <= 0 {
		return 0
	}
	if processed >= total {
		return 100
	}
	return int32(processed * 100 / total)
}

func enqueueSubtitleFulfillmentJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
	item storage.MediaItem,
) (int, error) {
	count := 0
	for _, input := range satisfaction.WantedTargetInputsForItem(item) {
		operation := input.Target.RequiredOperation
		if operation == nil || input.Target.Type != targets.TypeSubtitle {
			continue
		}
		jobID, err := enqueue(ctx, string(operation.Type), FulfillmentActionArgs{
			MediaItemID: item.ID.String(),
			FilePath:    input.FilePath,
			TargetType:  "subtitle",
			LanguageID:  input.Target.LanguageID,
		})
		if err != nil {
			return count, err
		}
		initializeQueuedFulfillmentProgress(ctx, settings, eventBroker, jobID, "Waiting for "+string(operation.Type), input.FilePath, input.Target.LanguageID)
		count++
	}
	return count, nil
}

func fulfillmentTaskKey(filePath string, languageID string) string {
	return filePath + "\x00" + languageMatchKey(languageID)
}

func initializeQueuedFulfillmentProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	jobID int64,
	label string,
	filePath string,
	languageID string,
) {
	if settings == nil {
		return
	}
	zero := int32(0)
	progress := normalizedProgressData(&zero, label, map[string]any{
		"filePath":   filePath,
		"languageId": languageID,
	})
	execution, err := settings.UpdateSystemJobExecutionProgressData(ctx, jobID, &zero, label, progress)
	if err == nil {
		publishJobExecutionUpdated(eventBroker, execution)
	}
}
