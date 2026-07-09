package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/events"
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
	mediatools "media-manager/internal/tools"
)

func executeTrackFulfillmentOperation(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	operation targets.OperationType,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
) error {
	if operation == targets.OperationVideoTranscode {
		return executeVideoTranscodeTrack(ctx, settings, eventBroker, item, args, track)
	}
	if operation == targets.OperationAudioTranscode {
		return executeAudioTranscodeTrack(ctx, settings, eventBroker, item, args, track)
	}
	if operation == targets.OperationContainerRemux {
		return executeContainerRemuxFile(ctx, settings, eventBroker, item, args)
	}
	if operation == targets.OperationSubtitleExtraction || operation == targets.OperationSubtitleConversion {
		return executeSubtitleFulfillmentOperation(ctx, settings, eventBroker, operation, item, args)
	}
	err := fmt.Errorf("fulfillment operation %s is not executable for a selected track yet", operation)
	publishFulfillmentExecutionError(ctx, settings, eventBroker, operation, item, args, err)
	return err
}

func executeTargetFulfillmentOperation(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	operation targets.OperationType,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	filePath string,
) error {
	if operation == targets.OperationSubtitleEmbed ||
		operation == targets.OperationSubtitleExtraction ||
		operation == targets.OperationSubtitleConversion {
		if args.FilePath == "" {
			args.FilePath = filePath
		}
		return executeSubtitleFulfillmentOperation(ctx, settings, eventBroker, operation, item, args)
	}
	err := fmt.Errorf("fulfillment operation %s needs a direct track or file id before execution", operation)
	details := fulfillmentActionDetails(operation, args)
	details["mediaItemId"] = item.ID.String()
	details["title"] = item.Title
	details["filePath"] = filePath
	details["error"] = err.Error()
	publishSystemEvent(ctx, settings, eventBroker, jobEventError, "media", "Fulfillment execution blocked", details)
	return err
}

func executeAudioTranscodeTrack(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
) error {
	target, ok := audioTargetForTrack(item, args, track)
	if !ok {
		err := fmt.Errorf("no audio target matches selected track language")
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationAudioTranscode, item, args, err)
		return err
	}
	policy, err := audioTranscodePolicy(ctx, settings, item, args.Manual && args.TrackID != "")
	if err != nil {
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationAudioTranscode, item, args, err)
		return err
	}
	decision := DecideAudioConversion(audioConversionInputForTrack(policy, target, track))
	details := audioTranscodeDetails(item, args, track, target, decision)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode decision", details)
	if !decision.Needed {
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode skipped", details)
		return nil
	}
	if !decision.Allowed {
		err := fmt.Errorf("audio transcode blocked by policy: %s", decision.Reason)
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationAudioTranscode, item, args, err)
		return err
	}
	if !audioConversionHasExecutableWork(decision) {
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode skipped", details)
		return nil
	}
	outputPath, cleanup, err := tempOutputPath(track.FilePath)
	if err != nil {
		return err
	}
	defer cleanup()
	argsList, err := audioTrackTranscodeArgs(track.FilePath, outputPath, audioOrdinal(item, track), decision)
	if err != nil {
		return err
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode started", details)
	if err := runAudioTranscodeCommand(ctx, settings, eventBroker, item, track, argsList); err != nil {
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationAudioTranscode, item, args, err)
		return err
	}
	if err := replaceMediaFile(outputPath, track.FilePath); err != nil {
		return err
	}
	if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
		return fmt.Errorf("rescan media after audio transcode: %w", err)
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Audio transcode finished", details)
	return nil
}

func audioTranscodePolicy(
	ctx context.Context,
	settings *storage.SettingsStore,
	item storage.MediaItem,
	manualTrackAction bool,
) (string, error) {
	if manualTrackAction {
		return "manual", nil
	}
	if item.QualityProfileID == nil || strings.TrimSpace(*item.QualityProfileID) == "" {
		return "disabled", nil
	}
	profile, err := settings.GetMediaProfile(ctx, *item.QualityProfileID)
	if err != nil {
		return "", fmt.Errorf("load media profile for audio transcode: %w", err)
	}
	return profile.AudioLossyTranscodePolicy, nil
}

func audioTargetForTrack(
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
) (storage.MediaProfileAudioTarget, bool) {
	language := stringPtrValue(track.LanguageID)
	if language == "" {
		language = args.LanguageID
	}
	for _, target := range item.AudioTargets {
		if satisfaction.LanguageMatches(target.LanguageID, language) {
			return target, true
		}
	}
	return storage.MediaProfileAudioTarget{}, false
}

func audioTrackTranscodeArgs(
	inputPath string,
	outputPath string,
	audioIndex int,
	decision AudioConversionDecision,
) ([]string, error) {
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	if decision.TargetCodec == "" {
		return nil, fmt.Errorf("target audio codec is required")
	}
	stream := "a:" + strconv.Itoa(audioIndex)
	args := []string{
		"-hide_banner", "-loglevel", "error", "-y", "-i", inputPath,
		"-map", "0", "-c", "copy", "-c:" + stream, ffmpegAudioCodec(decision.TargetCodec),
	}
	if channels := ffmpegChannelCount(decision.TargetChannels); channels != "" {
		args = append(args, "-ac:"+stream, channels)
	}
	if decision.TargetBitrateKbps > 0 {
		args = append(args, "-b:"+stream, strconv.Itoa(int(decision.TargetBitrateKbps))+"k")
	}
	return append(args, outputPath), nil
}

func audioOrdinal(item storage.MediaItem, selected storage.MediaFileTrackFact) int {
	for _, fact := range item.FileFacts {
		if fact.FilePath != selected.FilePath {
			continue
		}
		audioIndex := 0
		for _, track := range fact.Tracks {
			if track.ID == selected.ID {
				return audioIndex
			}
			if track.TrackType == "audio" {
				audioIndex++
			}
		}
	}
	return 0
}

func tempOutputPath(inputPath string) (string, func(), error) {
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return "", func() {}, err
	}
	pattern := "." + strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + ".transcode-*" + filepath.Ext(inputPath)
	file, err := os.CreateTemp(filepath.Dir(inputPath), pattern)
	if err != nil {
		return "", func() {}, err
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", func() {}, err
	}
	return path, func() { _ = os.Remove(path) }, nil
}

func replaceMediaFile(sourcePath string, targetPath string) error {
	info, err := os.Stat(targetPath)
	if err == nil {
		_ = os.Chmod(sourcePath, info.Mode())
	}
	if err := os.Rename(sourcePath, targetPath); err != nil {
		return fmt.Errorf("replace media file: %w", err)
	}
	return nil
}

func audioTranscodeDetails(
	item storage.MediaItem,
	args FulfillmentActionArgs,
	track storage.MediaFileTrackFact,
	target storage.MediaProfileAudioTarget,
	decision AudioConversionDecision,
) map[string]any {
	details := fulfillmentActionDetails(targets.OperationAudioTranscode, args)
	details["mediaItemId"] = item.ID.String()
	details["title"] = item.Title
	details["trackId"] = track.ID.String()
	details["streamIndex"] = track.StreamIndex
	details["sourceCodec"] = stringPtrValue(track.Codec)
	details["sourceChannels"] = stringPtrValue(track.Channels)
	details["sourceBitrateKbps"] = int32PtrValue(track.BitrateKbps)
	details["targetLanguageId"] = target.LanguageID
	details["decisionStatus"] = decision.Status
	details["decisionReason"] = decision.Reason
	details["targetCodec"] = decision.TargetCodec
	details["targetChannels"] = decision.TargetChannels
	details["targetBitrateKbps"] = decision.TargetBitrateKbps
	return details
}

func publishFulfillmentExecutionError(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	operation targets.OperationType,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	err error,
) {
	details := fulfillmentActionDetails(operation, args)
	details["mediaItemId"] = item.ID.String()
	details["title"] = item.Title
	details["error"] = err.Error()
	publishSystemEvent(ctx, settings, eventBroker, jobEventError, "media", "Fulfillment execution failed", details)
}
