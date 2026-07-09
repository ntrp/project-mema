package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
	mediatools "media-manager/internal/tools"
)

func executeContainerRemuxFile(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) error {
	fact, ok := mediaFactForPath(item, args.FilePath)
	if !ok {
		err := fmt.Errorf("media file fact not found for remux: %s", args.FilePath)
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationContainerRemux, item, args, err)
		return err
	}
	targetContainer := normalizedContainer(item.FinalContainer)
	details := containerRemuxDetails(item, args, fact, targetContainer)
	if targetContainer == "" {
		err := fmt.Errorf("profile final container is required")
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationContainerRemux, item, args, err)
		return err
	}
	if normalizedContainer(mediaFactContainer(fact)) == targetContainer {
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Container remux skipped", details)
		return nil
	}
	targetPath := remuxTargetPath(fact.FilePath, targetContainer)
	outputPath, cleanup, err := tempRemuxOutputPath(fact.FilePath, targetContainer)
	if err != nil {
		return err
	}
	defer cleanup()
	argsList, err := containerRemuxArgs(fact.FilePath, outputPath, targetContainer)
	if err != nil {
		return err
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Container remux started", details)
	if err := runContainerRemuxCommand(ctx, settings, eventBroker, item, fact, argsList); err != nil {
		publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationContainerRemux, item, args, err)
		return err
	}
	if err := replaceRemuxedMediaFile(outputPath, fact.FilePath, targetPath); err != nil {
		return err
	}
	if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
		return fmt.Errorf("rescan media after container remux: %w", err)
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Container remux finished", details)
	return nil
}

func containerRemuxArgs(inputPath string, outputPath string, targetContainer string) ([]string, error) {
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-i", inputPath, "-map", "0", "-c", "copy"}
	if normalizedContainer(targetContainer) == "mp4" {
		args = append(args, "-movflags", "+faststart")
	}
	return append(args, outputPath), nil
}

func mediaFactForPath(item storage.MediaItem, path string) (storage.MediaFileFact, bool) {
	for _, fact := range item.FileFacts {
		if fact.FilePath == path {
			return fact, true
		}
	}
	return storage.MediaFileFact{}, false
}

func mediaFactContainer(fact storage.MediaFileFact) string {
	if fact.ContainerFormat != nil && strings.TrimSpace(*fact.ContainerFormat) != "" {
		return *fact.ContainerFormat
	}
	return filepath.Ext(fact.FilePath)
}

func normalizedContainer(value string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".")
}

func remuxTargetPath(inputPath string, targetContainer string) string {
	ext := "." + normalizedContainer(targetContainer)
	return strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ext
}

func tempRemuxOutputPath(inputPath string, targetContainer string) (string, func(), error) {
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return "", func() {}, err
	}
	ext := "." + normalizedContainer(targetContainer)
	pattern := "." + strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + ".remux-*" + ext
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

func replaceRemuxedMediaFile(outputPath string, sourcePath string, targetPath string) error {
	if info, err := os.Stat(sourcePath); err == nil {
		_ = os.Chmod(outputPath, info.Mode())
	}
	if err := os.Rename(outputPath, targetPath); err != nil {
		return fmt.Errorf("replace remuxed media file: %w", err)
	}
	if sourcePath != targetPath {
		_ = os.Remove(sourcePath)
	}
	return nil
}

func containerRemuxDetails(
	item storage.MediaItem,
	args FulfillmentActionArgs,
	fact storage.MediaFileFact,
	targetContainer string,
) map[string]any {
	details := fulfillmentActionDetails(targets.OperationContainerRemux, args)
	details["mediaItemId"] = item.ID.String()
	details["title"] = item.Title
	details["filePath"] = fact.FilePath
	details["sourceContainer"] = normalizedContainer(mediaFactContainer(fact))
	details["targetContainer"] = targetContainer
	return details
}
