package jobs

import (
	"context"
	"fmt"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

type videoTranscodePlan struct {
	args     FulfillmentActionArgs
	track    storage.MediaFileTrackFact
	decision VideoConversionDecision
}

func enqueueTrackTranscodeJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
	operation targets.OperationType,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) (int, error) {
	switch operation {
	case targets.OperationVideoTranscode:
		return enqueueVideoTranscodeTrackJobs(ctx, settings, eventBroker, enqueue, item, args)
	case targets.OperationAudioTranscode:
		return enqueueAudioTranscodeTrackJobs(ctx, settings, eventBroker, enqueue, item, args)
	default:
		return 0, nil
	}
}

func enqueueVideoTranscodeTrackJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) (int, error) {
	if enqueue == nil {
		return 0, fmt.Errorf("video transcode enqueue function is not configured")
	}
	if args.TargetType != "" && args.TargetType != "video" {
		return 0, nil
	}
	plans := videoTranscodePlans(item, args)
	for _, plan := range plans {
		jobID, err := enqueue(ctx, "video_transcode", plan.args)
		if err != nil {
			return 0, err
		}
		initializeQueuedVideoTranscodeProgress(ctx, settings, eventBroker, jobID, item, plan.track)
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Video transcode track queued", map[string]any{
			"mediaItemId":    item.ID.String(),
			"title":          item.Title,
			"filePath":       plan.args.FilePath,
			"trackId":        plan.args.TrackID,
			"jobId":          jobID,
			"decisionReason": plan.decision.Reason,
			"targetCodec":    plan.decision.TargetCodec,
			"targetPixel":    plan.decision.TargetPixel,
		})
	}
	return len(plans), nil
}

func videoTranscodePlans(item storage.MediaItem, args FulfillmentActionArgs) []videoTranscodePlan {
	plans := []videoTranscodePlan{}
	for _, fact := range item.FileFacts {
		if args.FilePath != "" && fact.FilePath != args.FilePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.TrackType != "video" {
				continue
			}
			decision := DecideVideoConversion(videoConversionInputForTrack(item.VideoTarget, track))
			if !decision.Needed || !decision.Allowed {
				continue
			}
			plans = append(plans, videoTranscodePlan{
				args: FulfillmentActionArgs{
					MediaItemID: item.ID.String(),
					FilePath:    track.FilePath,
					TargetType:  "video",
					TrackID:     track.ID.String(),
				},
				track:    track,
				decision: decision,
			})
		}
	}
	return plans
}
