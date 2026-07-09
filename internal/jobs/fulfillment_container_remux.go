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
		return failContainerRemux(ctx, settings, eventBroker, item, args, err)
	}
	if err := validateRemuxOutput(outputPath); err != nil {
		return failContainerRemux(ctx, settings, eventBroker, item, args, err)
	}
	if err := replaceRemuxedMediaFile(outputPath, fact.FilePath, targetPath); err != nil {
		return failContainerRemux(ctx, settings, eventBroker, item, args, err)
	}
	if err := settings.RecordContainerRemuxedMediaFile(ctx, item.ID, fact.FilePath, targetPath); err != nil {
		if fact.FilePath != targetPath {
			_ = os.Remove(targetPath)
		}
		return failContainerRemux(ctx, settings, eventBroker, item, args, err)
	}
	if err := removeRemuxSourceFile(fact.FilePath, targetPath); err != nil {
		return failContainerRemux(ctx, settings, eventBroker, item, args, err)
	}
	if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
		return failContainerRemux(ctx, settings, eventBroker, item, args, fmt.Errorf("rescan media after container remux: %w", err))
	}
	finalizeContainerRemuxProgress(ctx, settings, eventBroker, item, fact)
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
		args = append(args, "-c:s", "mov_text", "-movflags", "+faststart")
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
	if sourcePath != targetPath {
		if _, err := os.Stat(targetPath); err == nil {
			return fmt.Errorf("remux target already exists: %s", targetPath)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check remux target: %w", err)
		}
	}
	if info, err := os.Stat(sourcePath); err == nil {
		_ = os.Chmod(outputPath, info.Mode())
	}
	if err := os.Rename(outputPath, targetPath); err != nil {
		return fmt.Errorf("replace remuxed media file: %w", err)
	}
	return nil
}

func validateRemuxOutput(outputPath string) error {
	info, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("validate remux output: %w", err)
	}
	if info.IsDir() || info.Size() == 0 {
		return fmt.Errorf("remux output is empty: %s", outputPath)
	}
	return nil
}

func removeRemuxSourceFile(sourcePath string, targetPath string) error {
	if sourcePath == targetPath {
		return nil
	}
	if err := os.Remove(sourcePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove remux source: %w", err)
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

func failContainerRemux(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
	err error,
) error {
	publishFulfillmentExecutionError(ctx, settings, eventBroker, targets.OperationContainerRemux, item, args, err)
	return err
}
