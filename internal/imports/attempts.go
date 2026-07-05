package imports

import (
	"context"
	"log/slog"
	"os"

	"media-manager/internal/storage"
)

const (
	importStatusSucceeded = "succeeded"
	importStatusFailed    = "failed"
)

type importRun struct {
	source string
	target string
	mode   ImportMode
}

func (s *Service) importWithAttempt(
	ctx context.Context,
	activity storage.DownloadActivity,
	item storage.MediaItem,
	run importRun,
) error {
	mode, err := normalizeImportMode(run.mode)
	if err != nil {
		s.recordImportAttempt(ctx, activity, run, importStatusFailed, "file_operation", err, nil, nil)
		return err
	}
	run.mode = mode
	alreadyPresent := sameExistingFile(run.source, run.target)
	if err := importFile(run.source, run.target, mode); err != nil {
		s.recordImportAttempt(ctx, activity, run, importStatusFailed, "file_operation", err, nil, nil)
		return err
	}

	createdTargets := createdTargets(run.target, alreadyPresent)
	if err := s.settings.RecordImportedMediaFileWithHistory(ctx, item, run.source, run.target, string(mode)); err != nil {
		rollbackErr := rollbackImportedTarget(run)
		recordErr := err
		if rollbackErr != nil {
			slog.Error("import rollback failed", "activityId", activity.ID, "target", run.target, "error", rollbackErr)
		}
		s.recordImportAttempt(ctx, activity, run, importStatusFailed, "record_media_file", recordErr, createdTargets, nil)
		return err
	}
	s.recordImportAttempt(ctx, activity, run, importStatusSucceeded, "", nil, createdTargets, []string{run.target})
	return nil
}

func (s *Service) recordImportAttempt(
	ctx context.Context,
	activity storage.DownloadActivity,
	run importRun,
	status string,
	stage string,
	cause error,
	createdTargets []string,
	insertedMediaFiles []string,
) {
	input := storage.ImportAttemptInput{
		ActivityID:             activity.ID,
		MediaItemID:            activity.MediaItemID,
		SourcePath:             optionalImportString(run.source),
		TargetPath:             optionalImportString(run.target),
		ImportMode:             string(run.mode),
		Status:                 status,
		FailureStage:           optionalImportString(stage),
		CreatedTargets:         createdTargets,
		InsertedMediaFilePaths: insertedMediaFiles,
	}
	if cause != nil {
		message := cause.Error()
		input.ErrorMessage = &message
	}
	if _, err := s.settings.CreateImportAttempt(ctx, input); err != nil {
		slog.Error("record import attempt failed", "activityId", activity.ID, "status", status, "stage", stage, "error", err)
	}
}

func createdTargets(target string, alreadyPresent bool) []string {
	if target == "" || alreadyPresent {
		return nil
	}
	return []string{target}
}

func rollbackImportedTarget(run importRun) error {
	if run.target == "" {
		return nil
	}
	if run.mode == ImportModeMove {
		if _, err := os.Stat(run.source); os.IsNotExist(err) {
			if _, err := os.Stat(run.target); err == nil {
				return os.Rename(run.target, run.source)
			}
		}
	}
	return os.Remove(run.target)
}

func optionalImportString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
