package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/events"
	"media-manager/internal/mediafacts"
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
	"media-manager/internal/targets"
)

type FulfillmentActionArgs struct {
	MediaItemID        string `json:"media_item_id,omitempty" river:"unique"`
	FilePath           string `json:"file_path,omitempty" river:"unique"`
	TargetType         string `json:"target_type,omitempty" river:"unique"`
	LanguageID         string `json:"language_id,omitempty" river:"unique"`
	TrackID            string `json:"track_id,omitempty" river:"unique"`
	OtherFileID        string `json:"other_file_id,omitempty" river:"unique"`
	ExternalSubtitleID string `json:"external_subtitle_id,omitempty" river:"unique"`
	Manual             bool   `json:"manual,omitempty"`
}

type fulfillmentEnqueueFunc func(context.Context, string, FulfillmentActionArgs) (int64, error)

type VideoTranscodeArgs struct{ FulfillmentActionArgs }
type AudioTranscodeArgs struct{ FulfillmentActionArgs }
type AudioSourceArgs struct{ FulfillmentActionArgs }
type ContainerRemuxArgs struct{ FulfillmentActionArgs }
type SubtitleDownloadArgs struct{ FulfillmentActionArgs }
type SubtitleEmbedArgs struct{ FulfillmentActionArgs }
type SubtitleExtractArgs struct{ FulfillmentActionArgs }
type SubtitleConvertArgs struct{ FulfillmentActionArgs }

func (VideoTranscodeArgs) Kind() string   { return "media.fulfillment.video_transcode" }
func (AudioTranscodeArgs) Kind() string   { return "media.fulfillment.audio_transcode" }
func (AudioSourceArgs) Kind() string      { return "media.fulfillment.audio_source" }
func (ContainerRemuxArgs) Kind() string   { return "media.fulfillment.container_remux" }
func (SubtitleDownloadArgs) Kind() string { return "media.fulfillment.subtitle_download" }
func (SubtitleEmbedArgs) Kind() string    { return "media.fulfillment.subtitle_embed" }
func (SubtitleExtractArgs) Kind() string  { return "media.fulfillment.subtitle_extract" }
func (SubtitleConvertArgs) Kind() string  { return "media.fulfillment.subtitle_convert" }

type VideoTranscodeWorker struct {
	river.WorkerDefaults[VideoTranscodeArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type AudioTranscodeWorker struct {
	river.WorkerDefaults[AudioTranscodeArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type AudioSourceWorker struct {
	river.WorkerDefaults[AudioSourceArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type ContainerRemuxWorker struct {
	river.WorkerDefaults[ContainerRemuxArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type SubtitleDownloadWorker struct {
	river.WorkerDefaults[SubtitleDownloadArgs]
	settings  *storage.SettingsStore
	subtitles *subtitles.Service
	events    *events.Broker
}

type SubtitleEmbedWorker struct {
	river.WorkerDefaults[SubtitleEmbedArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type SubtitleExtractWorker struct {
	river.WorkerDefaults[SubtitleExtractArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type SubtitleConvertWorker struct {
	river.WorkerDefaults[SubtitleConvertArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

func (w *VideoTranscodeWorker) Work(ctx context.Context, job *river.Job[VideoTranscodeArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, w.enqueueFulfillment, targets.OperationVideoTranscode, job.Args.FulfillmentActionArgs)
}

func (w *AudioTranscodeWorker) Work(ctx context.Context, job *river.Job[AudioTranscodeArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, w.enqueueFulfillment, targets.OperationAudioTranscode, job.Args.FulfillmentActionArgs)
}

func (w *AudioSourceWorker) Work(ctx context.Context, job *river.Job[AudioSourceArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, nil, targets.OperationAudioSourcing, job.Args.FulfillmentActionArgs)
}

func (w *ContainerRemuxWorker) Work(ctx context.Context, job *river.Job[ContainerRemuxArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, nil, targets.OperationContainerRemux, job.Args.FulfillmentActionArgs)
}

func (w *SubtitleDownloadWorker) Work(ctx context.Context, job *river.Job[SubtitleDownloadArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, w.subtitles, nil, targets.OperationSubtitleDownload, job.Args.FulfillmentActionArgs)
}

func (w *SubtitleEmbedWorker) Work(ctx context.Context, job *river.Job[SubtitleEmbedArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, nil, targets.OperationSubtitleEmbed, job.Args.FulfillmentActionArgs)
}

func (w *SubtitleExtractWorker) Work(ctx context.Context, job *river.Job[SubtitleExtractArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, nil, targets.OperationSubtitleExtraction, job.Args.FulfillmentActionArgs)
}

func (w *SubtitleConvertWorker) Work(ctx context.Context, job *river.Job[SubtitleConvertArgs]) (err error) {
	return runFulfillmentWorker(ctx, job.JobRow, w.settings, w.events, nil, nil, targets.OperationSubtitleConversion, job.Args.FulfillmentActionArgs)
}

func runFulfillmentWorker(
	ctx context.Context,
	job *rivertype.JobRow,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	subtitleService *subtitles.Service,
	enqueueFulfillment fulfillmentEnqueueFunc,
	operation targets.OperationType,
	args FulfillmentActionArgs,
) (err error) {
	ctx = withJobExecution(ctx, job.ID)
	recordJobUpdated(ctx, settings, eventBroker, job, "running")
	defer func() { recordJobFinished(ctx, settings, eventBroker, job, err) }()
	details := fulfillmentActionDetails(operation, args)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Fulfillment worker invoked", details)
	recordJobProgressData(ctx, settings, eventBroker, nil, "Finding fulfillment targets", details)
	items, err := fulfillmentItems(ctx, settings, args.MediaItemID)
	if err != nil {
		return err
	}
	count := 0
	evaluated := 0
	skipped := 0
	for _, item := range items {
		liveItem := mediafacts.WithLiveFileFacts(item, args.FilePath)
		item = liveItem
		activeArgs := args
		var scopedTrack *storage.MediaFileTrackFact
		item, activeArgs, scopedTrack, err = fulfillmentApplyTrackScope(item, activeArgs)
		if err != nil {
			return err
		}
		if scopedTrack != nil {
			publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Fulfillment track resolved", map[string]any{
				"mediaItemId": item.ID.String(),
				"title":       item.Title,
				"operation":   operation,
				"filePath":    activeArgs.FilePath,
				"trackId":     scopedTrack.ID.String(),
				"trackType":   scopedTrack.TrackType,
				"streamIndex": scopedTrack.StreamIndex,
				"languageId":  stringPtrValue(scopedTrack.LanguageID),
				"codec":       stringPtrValue(scopedTrack.Codec),
			})
			if err := executeTrackFulfillmentOperation(
				ctx,
				settings,
				eventBroker,
				operation,
				liveItem,
				activeArgs,
				*scopedTrack,
			); err != nil {
				return err
			}
			count++
			continue
		}
		if operation == targets.OperationVideoTranscode || operation == targets.OperationAudioTranscode {
			added, err := enqueueTrackTranscodeJobs(ctx, settings, eventBroker, enqueueFulfillment, operation, item, activeArgs)
			if err != nil {
				return err
			}
			count += added
			continue
		}
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Fulfillment media evaluated", map[string]any{
			"mediaItemId": item.ID.String(),
			"title":       item.Title,
			"operation":   operation,
			"filePath":    activeArgs.FilePath,
			"factCount":   len(item.FileFacts),
			"trackCount":  fulfillmentTrackCount(item),
		})
		if operation == targets.OperationSubtitleDownload {
			count += runSubtitleDownloads(ctx, settings, subtitleService, eventBroker, item, activeArgs)
			continue
		}
		for _, input := range satisfaction.WantedTargetInputsForItem(item) {
			target := input.Target
			if activeArgs.FilePath != "" && input.FilePath != activeArgs.FilePath {
				skipped++
				continue
			}
			if !fulfillmentTargetInRequestScope(activeArgs, target) {
				skipped++
				continue
			}
			evaluated++
			if fulfillmentTargetMatches(operation, activeArgs, target) {
				count++
				publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Fulfillment action identified", map[string]any{
					"mediaItemId": item.ID.String(),
					"title":       item.Title,
					"operation":   operation,
					"filePath":    input.FilePath,
					"targetType":  target.Type,
					"targetState": target.State,
					"languageId":  target.LanguageID,
					"reasons":     target.Reasons,
				})
				if err := executeTargetFulfillmentOperation(ctx, settings, eventBroker, operation, item, activeArgs, input.FilePath); err != nil {
					return err
				}
				continue
			}
			publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Fulfillment target skipped", map[string]any{
				"mediaItemId": item.ID.String(),
				"title":       item.Title,
				"operation":   operation,
				"filePath":    input.FilePath,
				"targetType":  target.Type,
				"targetState": target.State,
				"languageId":  target.LanguageID,
				"reason":      fulfillmentSkipReason(operation, args, target),
			})
		}
	}
	done := int32(100)
	summary := fulfillmentActionDetails(operation, args)
	summary["targetCount"] = count
	summary["evaluatedTargetCount"] = evaluated
	summary["skippedByScopeCount"] = skipped
	recordJobProgressData(ctx, settings, eventBroker, &done, "Fulfillment execution summary", summary)
	if count == 0 {
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "media", "No fulfillment targets found", summary)
		if fulfillmentActionScoped(args) {
			return fmtScopedFulfillmentNotFound(args)
		}
	}
	return nil
}

func fulfillmentTargetMatches(operation targets.OperationType, args FulfillmentActionArgs, target targets.Target) bool {
	if !fulfillmentTargetInRequestScope(args, target) {
		return false
	}
	if target.RequiredOperation != nil && target.RequiredOperation.Type == operation {
		return true
	}
	if operation == targets.OperationVideoTranscode && target.Type == targets.TypeVideo && target.State == targets.StatePartial {
		return true
	}
	if operation == targets.OperationAudioTranscode && target.Type == targets.TypeAudio && target.State == targets.StatePartial {
		return true
	}
	return operation == targets.OperationAudioSourcing && target.Type == targets.TypeAudio && target.State == targets.StateMissing
}

func fulfillmentTargetInRequestScope(args FulfillmentActionArgs, target targets.Target) bool {
	if args.TargetType != "" && string(target.Type) != args.TargetType {
		return false
	}
	return args.LanguageID == "" || satisfaction.LanguageMatches(target.LanguageID, args.LanguageID)
}
