package jobs

import (
	"context"
	"fmt"
	"strconv"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
	mediatools "media-manager/internal/tools"
)

func executeVideoTranscodeTrack(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
) error {
	decision := DecideVideoConversion(videoConversionInputForTrack(item.VideoTarget, track))
	details := videoTranscodeDetails(item, args, track, decision)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Video transcode decision", details)
	if !decision.Needed {
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Video transcode skipped", details)
		return nil
	}
	if !decision.Allowed {
		err := fmt.Errorf("video transcode blocked: %s", decision.Reason)
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationVideoTranscode, item, args, err)
		return err
	}
	outputPath, cleanup, err := tempOutputPath(track.FilePath)
	if err != nil {
		return err
	}
	defer cleanup()
	argsList, err := videoTrackTranscodeArgs(track.FilePath, outputPath, videoOrdinal(item, track), decision)
	if err != nil {
		return err
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Video transcode started", details)
	if err := runVideoTranscodeCommand(ctx, settings, eventBroker, item, track, argsList); err != nil {
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationVideoTranscode, item, args, err)
		return err
	}
	if err := replaceMediaFile(outputPath, track.FilePath); err != nil {
		return err
	}
	if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
		return fmt.Errorf("rescan media after video transcode: %w", err)
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Video transcode finished", details)
	return nil
}

func videoTrackTranscodeArgs(
	inputPath string,
	outputPath string,
	videoIndex int,
	decision VideoConversionDecision,
) ([]string, error) {
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	stream := "v:" + strconv.Itoa(videoIndex)
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-i", inputPath, "-map", "0", "-c", "copy"}
	if decision.TargetCodec != "" {
		args = append(args, "-c:"+stream, ffmpegVideoCodec(decision.TargetCodec))
	}
	if decision.TargetPixel != "" {
		args = append(args, "-pix_fmt:"+stream, decision.TargetPixel)
	}
	return append(args, outputPath), nil
}

func videoOrdinal(item storage.MediaItem, selected storage.MediaFileTrackFact) int {
	for _, fact := range item.FileFacts {
		if fact.FilePath != selected.FilePath {
			continue
		}
		videoIndex := 0
		for _, track := range fact.Tracks {
			if track.ID == selected.ID {
				return videoIndex
			}
			if track.TrackType == "video" {
				videoIndex++
			}
		}
	}
	return 0
}

func videoTranscodeDetails(
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
	decision VideoConversionDecision,
) map[string]any {
	details := fulfillmentActionDetails(targets.OperationVideoTranscode, args)
	details["mediaItemId"] = item.ID.String()
	details["title"] = item.Title
	details["trackId"] = track.ID.String()
	details["streamIndex"] = track.StreamIndex
	details["sourceCodec"] = stringPtrValue(track.Codec)
	details["sourcePixel"] = stringPtrValue(track.PixelFormat)
	details["decisionStatus"] = decision.Status
	details["decisionReason"] = decision.Reason
	details["targetCodec"] = decision.TargetCodec
	details["targetPixel"] = decision.TargetPixel
	return details
}
