package jobs

import (
	"context"
	"fmt"

	"media-manager/internal/events"
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
)

type audioTranscodePlan struct {
	args     FulfillmentActionArgs
	track    storage.MediaFileTrackFact
	target   storage.MediaProfileAudioTarget
	decision AudioConversionDecision
}

func enqueueAudioTranscodeTrackJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) (int, error) {
	if enqueue == nil {
		return 0, fmt.Errorf("audio transcode enqueue function is not configured")
	}
	if args.TargetType != "" && args.TargetType != "audio" {
		return 0, nil
	}
	policy, err := audioTranscodePolicy(ctx, settings, item, false)
	if err != nil {
		return 0, err
	}
	plans := audioTranscodePlansForPolicy(policy, item, args)
	for _, plan := range plans {
		jobID, err := enqueue(ctx, "audio_transcode", plan.args)
		if err != nil {
			return 0, err
		}
		initializeQueuedAudioTranscodeProgress(ctx, settings, eventBroker, jobID, item, plan.track)
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode track queued", map[string]any{
			"mediaItemId":       item.ID.String(),
			"title":             item.Title,
			"filePath":          plan.args.FilePath,
			"trackId":           plan.args.TrackID,
			"languageId":        plan.args.LanguageID,
			"jobId":             jobID,
			"decisionReason":    plan.decision.Reason,
			"targetCodec":       plan.decision.TargetCodec,
			"targetChannels":    plan.decision.TargetChannels,
			"targetBitrateKbps": plan.decision.TargetBitrateKbps,
		})
	}
	return len(plans), nil
}

func initializeQueuedAudioTranscodeProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	jobID int64,
	item storage.MediaItem,
	track storage.MediaFileTrackFact,
) {
	if settings == nil {
		return
	}
	zero := int32(0)
	progress := normalizedProgressData(&zero, "Waiting to transcode audio", audioTranscodeProgressData(item, track, mediaFactDurationMs(item, track)))
	execution, err := settings.UpdateSystemJobExecutionProgressData(ctx, jobID, &zero, "Waiting to transcode audio", progress)
	if err == nil {
		publishJobExecutionUpdated(eventBroker, execution)
	}
}

func audioTranscodePlansForPolicy(
	policy string,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) []audioTranscodePlan {
	plans := []audioTranscodePlan{}
	for _, fact := range item.FileFacts {
		if args.FilePath != "" && fact.FilePath != args.FilePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.TrackType != "audio" {
				continue
			}
			target, ok := audioTargetForTrack(item, args, track)
			if !ok || !fulfillmentLanguageInScope(args, target.LanguageID) {
				continue
			}
			decision := DecideAudioConversion(audioConversionInputForTrack(policy, target, track))
			if !decision.Needed || !decision.Allowed {
				continue
			}
			plans = append(plans, audioTranscodePlan{
				args: FulfillmentActionArgs{
					MediaItemID: item.ID.String(),
					FilePath:    track.FilePath,
					TargetType:  "audio",
					LanguageID:  target.LanguageID,
					TrackID:     track.ID.String(),
				},
				track:    track,
				target:   target,
				decision: decision,
			})
		}
	}
	return plans
}

func fulfillmentLanguageInScope(args FulfillmentActionArgs, languageID string) bool {
	return args.LanguageID == "" || satisfaction.LanguageMatches(languageID, args.LanguageID)
}
