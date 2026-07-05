package storage

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"
)

const (
	FileDeleteModePermanent = "permanent"
	FileDeleteModeRecycle   = "recycle"
	FileDeleteModeKeep      = "keep"

	DefaultRecycleFolder = ".recycle"
)

type FileDeleteSettings struct {
	Mode          string
	RecycleFolder string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FileDeleteSettingsInput struct {
	Mode          string
	RecycleFolder string
}

func (s *SettingsStore) GetFileDeleteSettings(ctx context.Context) (FileDeleteSettings, error) {
	row, err := storagegen.New(s.pool).GetFileDeleteSettings(ctx, storagegen.GetFileDeleteSettingsParams{
		Mode:          FileDeleteModePermanent,
		RecycleFolder: DefaultRecycleFolder,
	})
	if err != nil {
		return FileDeleteSettings{}, err
	}
	return fileDeleteSettingsFromGetRow(row), nil
}

func (s *SettingsStore) UpdateFileDeleteSettings(
	ctx context.Context,
	input FileDeleteSettingsInput,
) (FileDeleteSettings, error) {
	input = normalizeFileDeleteSettings(input)
	if err := validateFileDeleteSettings(input); err != nil {
		return FileDeleteSettings{}, err
	}
	row, err := storagegen.New(s.pool).UpdateFileDeleteSettings(ctx, storagegen.UpdateFileDeleteSettingsParams{
		Mode:          input.Mode,
		RecycleFolder: input.RecycleFolder,
	})
	if err != nil {
		return FileDeleteSettings{}, err
	}
	return fileDeleteSettingsFromUpdateRow(row), nil
}

func normalizeFileDeleteSettings(input FileDeleteSettingsInput) FileDeleteSettingsInput {
	input.Mode = strings.TrimSpace(input.Mode)
	if input.Mode == "" {
		input.Mode = FileDeleteModePermanent
	}
	input.RecycleFolder = strings.TrimSpace(input.RecycleFolder)
	if input.RecycleFolder == "" {
		input.RecycleFolder = DefaultRecycleFolder
	}
	return input
}

func validateFileDeleteSettings(input FileDeleteSettingsInput) error {
	switch input.Mode {
	case FileDeleteModePermanent, FileDeleteModeRecycle, FileDeleteModeKeep:
	default:
		return ErrInvalidInput
	}
	if filepath.IsAbs(input.RecycleFolder) {
		return ErrInvalidInput
	}
	if _, err := safePathUnderRoot(string(filepath.Separator)+"library", input.RecycleFolder, false); err != nil {
		return err
	}
	firstSegment := strings.Split(filepath.ToSlash(filepath.Clean(input.RecycleFolder)), "/")[0]
	if !strings.HasPrefix(firstSegment, ".") {
		return ErrInvalidInput
	}
	return nil
}

func fileDeleteSettingsFromGetRow(row storagegen.GetFileDeleteSettingsRow) FileDeleteSettings {
	return FileDeleteSettings{
		Mode:          row.Mode,
		RecycleFolder: row.RecycleFolder,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func fileDeleteSettingsFromUpdateRow(row storagegen.UpdateFileDeleteSettingsRow) FileDeleteSettings {
	return FileDeleteSettings{
		Mode:          row.Mode,
		RecycleFolder: row.RecycleFolder,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
