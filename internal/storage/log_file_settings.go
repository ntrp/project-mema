package storage

import (
	"context"
	"strings"
	"time"
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
	settings, err := scanLogFileSettings(s.pool.QueryRow(ctx, `
		insert into app.log_file_settings (id, enabled, directory, retention_days)
		values (true, false, $1, $2)
		on conflict (id) do update set id = excluded.id
		returning enabled, directory, retention_days, created_at, updated_at
	`, DefaultLogFileDirectory, DefaultLogRetentionDays))
	if err != nil {
		return LogFileSettings{}, err
	}
	return settings, nil
}

func (s *SettingsStore) UpdateLogFileSettings(
	ctx context.Context,
	input LogFileSettingsInput,
) (LogFileSettings, error) {
	input = normalizeLogFileSettings(input)
	if input.Directory == "" || input.RetentionDays < 1 || input.RetentionDays > 365 {
		return LogFileSettings{}, ErrInvalidInput
	}
	return scanLogFileSettings(s.pool.QueryRow(ctx, `
		insert into app.log_file_settings (id, enabled, directory, retention_days)
		values (true, $1, $2, $3)
		on conflict (id) do update
		set enabled = excluded.enabled,
			directory = excluded.directory,
			retention_days = excluded.retention_days,
			updated_at = now()
		returning enabled, directory, retention_days, created_at, updated_at
	`, input.Enabled, input.Directory, input.RetentionDays))
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

func scanLogFileSettings(row interface {
	Scan(dest ...any) error
}) (LogFileSettings, error) {
	var settings LogFileSettings
	err := row.Scan(
		&settings.Enabled,
		&settings.Directory,
		&settings.RetentionDays,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	return settings, err
}
