package storage

import (
	"context"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"
)

const (
	DefaultLogFileDirectory = ".data/logs"
	DefaultLogRetentionDays = 7
)

type LogFileSettings struct {
	Enabled       bool
	Directory     string
	RetentionDays int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type LogFileSettingsInput struct {
	Enabled       bool
	Directory     string
	RetentionDays int32
}

func (s *SettingsStore) GetLogFileSettings(ctx context.Context) (LogFileSettings, error) {
	row, err := storagegen.New(s.pool).GetLogFileSettings(ctx, storagegen.GetLogFileSettingsParams{
		Directory:     DefaultLogFileDirectory,
		RetentionDays: DefaultLogRetentionDays,
	})
	if err != nil {
		return LogFileSettings{}, err
	}
	return logFileSettingsFromRow(row), nil
}

func (s *SettingsStore) UpdateLogFileSettings(
	ctx context.Context,
	input LogFileSettingsInput,
) (LogFileSettings, error) {
	input = normalizeLogFileSettings(input)
	if input.Directory == "" || input.RetentionDays < 1 || input.RetentionDays > 365 {
		return LogFileSettings{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).UpdateLogFileSettings(ctx, storagegen.UpdateLogFileSettingsParams{
		Enabled:       input.Enabled,
		Directory:     input.Directory,
		RetentionDays: input.RetentionDays,
	})
	if err != nil {
		return LogFileSettings{}, err
	}
	return logFileSettingsFromUpdateRow(row), nil
}

func normalizeLogFileSettings(input LogFileSettingsInput) LogFileSettingsInput {
	input.Directory = strings.TrimSpace(input.Directory)
	if input.Directory == "" {
		input.Directory = DefaultLogFileDirectory
	}
	if input.RetentionDays == 0 {
		input.RetentionDays = DefaultLogRetentionDays
	}
	return input
}

func logFileSettingsFromRow(row storagegen.GetLogFileSettingsRow) LogFileSettings {
	return LogFileSettings{
		Enabled:       row.Enabled,
		Directory:     row.Directory,
		RetentionDays: row.RetentionDays,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func logFileSettingsFromUpdateRow(row storagegen.UpdateLogFileSettingsRow) LogFileSettings {
	return LogFileSettings{
		Enabled:       row.Enabled,
		Directory:     row.Directory,
		RetentionDays: row.RetentionDays,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
