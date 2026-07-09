package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/events"
	"media-manager/internal/mediafacts"
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
	mediatools "media-manager/internal/tools"
)

func executeSubtitleFulfillmentOperation(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	operation targets.OperationType,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) error {
	switch operation {
	case targets.OperationSubtitleExtraction:
		return executeSubtitleExtraction(ctx, settings, eventBroker, item, args)
	case targets.OperationSubtitleConversion:
		return executeSubtitleConversion(ctx, settings, eventBroker, item, args)
	case targets.OperationSubtitleEmbed:
		return executeSubtitleEmbed(ctx, settings, eventBroker, item, args)
	default:
		return fmt.Errorf("unsupported subtitle operation %s", operation)
	}
}

func executeSubtitleExtraction(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) error {
	track, ok := subtitleTrackCandidate(item, args)
	if !ok {
		return fmt.Errorf("subtitle extraction target not found")
	}
	format := firstNonEmpty(subtitleTargetFormat(item, stringPtrValue(track.LanguageID)), stringPtrValue(track.Format), "subrip")
	outputPath := subtitleSidecarPath(track.FilePath, stringPtrValue(track.LanguageID), format)
	if err := runSubtitleExtractCommand(ctx, track.FilePath, subtitleOrdinal(item, track), outputPath); err != nil {
		return err
	}
	if err := recordGeneratedSubtitle(ctx, settings, item, outputPath, stringPtrValue(track.LanguageID), format, storage.SubtitleRetentionExternal); err != nil {
		return err
	}
	recordJobProgress(ctx, settings, eventBroker, progressInt32(100), "Subtitle extraction complete")
	return nil
}

func executeSubtitleConversion(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) error {
	if subtitle, ok := externalSubtitleCandidate(item, args); ok {
		target := firstNonEmpty(subtitleTargetFormat(item, subtitle.LanguageID), subtitle.Format)
		return convertExternalSubtitle(ctx, settings, eventBroker, item, subtitle, args.FilePath, target)
	}
	track, ok := subtitleTrackCandidate(item, args)
	if !ok {
		return fmt.Errorf("subtitle conversion target not found")
	}
	format := firstNonEmpty(subtitleTargetFormat(item, stringPtrValue(track.LanguageID)), stringPtrValue(track.Format), "subrip")
	outputPath := subtitleSidecarPath(track.FilePath, stringPtrValue(track.LanguageID), format)
	if err := runSubtitleExtractCommand(ctx, track.FilePath, subtitleOrdinal(item, track), outputPath); err != nil {
		return err
	}
	if err := recordGeneratedSubtitle(ctx, settings, item, outputPath, stringPtrValue(track.LanguageID), format, subtitleRetentionModeForProfile(item.SubtitleMode)); err != nil {
		return err
	}
	recordJobProgress(ctx, settings, eventBroker, progressInt32(100), "Subtitle conversion complete")
	return nil
}

func executeSubtitleEmbed(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) error {
	subtitle, ok := externalSubtitleCandidate(item, args)
	if !ok {
		return fmt.Errorf("subtitle merge target not found")
	}
	outputPath, cleanup, err := tempOutputPath(args.FilePath)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := runSubtitleEmbedCommand(ctx, args.FilePath, subtitle.FilePath, outputPath); err != nil {
		return err
	}
	if err := replaceMediaFile(outputPath, args.FilePath); err != nil {
		return err
	}
	if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
		return fmt.Errorf("rescan media after subtitle merge: %w", err)
	}
	recordJobProgress(ctx, settings, eventBroker, progressInt32(100), "Subtitle merge complete")
	return nil
}

func convertExternalSubtitle(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	subtitle storage.MediaItemSubtitle,
	mediaPath string,
	targetFormat string,
) error {
	content, err := os.ReadFile(subtitle.FilePath)
	if err != nil {
		return err
	}
	converted, format, err := convertSubtitleContent(content, targetFormat)
	if err != nil {
		return err
	}
	outputPath := subtitleSidecarPath(firstNonEmpty(mediaPath, subtitle.FilePath), subtitle.LanguageID, format)
	if err := os.WriteFile(outputPath, converted, 0o644); err != nil {
		return err
	}
	if err := recordGeneratedSubtitle(ctx, settings, item, outputPath, subtitle.LanguageID, format, subtitle.RetentionMode); err != nil {
		return err
	}
	recordJobProgress(ctx, settings, eventBroker, progressInt32(100), "Subtitle conversion complete")
	return nil
}

func runSubtitleExtractCommand(ctx context.Context, inputPath string, subtitleIndex int, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	_, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           "ffmpeg",
		Args:           []string{"-y", "-i", inputPath, "-map", fmt.Sprintf("0:s:%d", subtitleIndex), outputPath},
		Timeout:        30 * time.Minute,
		MaxOutputBytes: 16 * 1024,
		MaxStderrBytes: 64 * 1024,
	})
	return err
}

func runSubtitleEmbedCommand(ctx context.Context, mediaPath string, subtitlePath string, outputPath string) error {
	_, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           "ffmpeg",
		Args:           []string{"-y", "-i", mediaPath, "-i", subtitlePath, "-map", "0", "-map", "1:0", "-c", "copy", outputPath},
		Timeout:        2 * time.Hour,
		MaxOutputBytes: 16 * 1024,
		MaxStderrBytes: 128 * 1024,
	})
	return err
}

func recordGeneratedSubtitle(ctx context.Context, settings *storage.SettingsStore, item storage.MediaItem, path string, languageID string, format string, retention storage.SubtitleRetentionMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	_, err = settings.UpsertMediaItemSubtitle(ctx, storage.MediaItemSubtitleInput{
		MediaItemID:   item.ID,
		LanguageID:    languageID,
		Format:        format,
		FilePath:      path,
		ProviderName:  "local",
		SizeBytes:     progressInt64(info.Size()),
		RetentionMode: retention,
	})
	return err
}

func subtitleTrackCandidate(item storage.MediaItem, args FulfillmentActionArgs) (storage.MediaFileTrackFact, bool) {
	for _, fact := range item.FileFacts {
		if args.FilePath != "" && fact.FilePath != args.FilePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.TrackType != "subtitle" || !fulfillmentLanguageMatches(args, stringPtrValue(track.LanguageID)) {
				continue
			}
			if args.TrackID == "" || args.TrackID == track.ID.String() {
				return track, true
			}
		}
	}
	return storage.MediaFileTrackFact{}, false
}

func externalSubtitleCandidate(item storage.MediaItem, args FulfillmentActionArgs) (storage.MediaItemSubtitle, bool) {
	for _, subtitle := range item.ExternalSubtitles {
		if args.ExternalSubtitleID != "" && subtitle.ID.String() != args.ExternalSubtitleID {
			continue
		}
		if !fulfillmentLanguageMatches(args, subtitle.LanguageID) || !sameMediaBase(subtitle.FilePath, args.FilePath) {
			continue
		}
		if args.OtherFileID == "" || args.OtherFileID == subtitle.ID.String() {
			return subtitle, true
		}
	}
	for _, sidecar := range item.Sidecars {
		if !sidecarMatchesArgs(item.ID, sidecar, args) {
			continue
		}
		return storage.MediaItemSubtitle{MediaItemID: item.ID, LanguageID: stringPtrValue(sidecar.LanguageID), Format: stringPtrValue(sidecar.Format), FilePath: sidecar.FilePath, Selected: true, RetentionMode: subtitleRetentionModeForProfile(item.SubtitleMode)}, true
	}
	return storage.MediaItemSubtitle{}, false
}

func sidecarMatchesArgs(mediaItemID uuid.UUID, sidecar storage.MediaItemSidecar, args FulfillmentActionArgs) bool {
	if sidecar.SidecarType != storage.MediaSidecarSubtitle || sidecar.MediaFilePath != args.FilePath {
		return false
	}
	if !fulfillmentLanguageMatches(args, stringPtrValue(sidecar.LanguageID)) {
		return false
	}
	synthetic := mediafacts.OtherFileID(mediaItemID, args.FilePath, sidecar.FilePath, string(storage.MediaSidecarSubtitle)).String()
	return args.OtherFileID == "" || args.OtherFileID == sidecar.ID.String() || args.OtherFileID == synthetic
}

func fulfillmentLanguageMatches(args FulfillmentActionArgs, languageID string) bool {
	return args.LanguageID == "" || satisfaction.LanguageMatches(languageID, args.LanguageID)
}

func subtitleOrdinal(item storage.MediaItem, selected storage.MediaFileTrackFact) int {
	index := 0
	for _, fact := range item.FileFacts {
		if fact.FilePath != selected.FilePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.TrackType != "subtitle" {
				continue
			}
			if track.ID == selected.ID {
				return index
			}
			index++
		}
	}
	return 0
}

func progressInt32(value int32) *int32 { return &value }

func progressInt64(value int64) *int64 { return &value }
