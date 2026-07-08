package storage

import (
	"context"
	"encoding/json"
	"strings"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) SyncSystemJobSchedules(ctx context.Context, definitions []SystemJobScheduleDefinition) error {
	queries := storagegen.New(s.pool)
	for _, definition := range definitions {
		if strings.TrimSpace(definition.ID) == "" || definition.IntervalSeconds < MinSystemJobScheduleIntervalSeconds {
			return ErrInvalidInput
		}
		_, err := queries.UpsertSystemJobSchedule(ctx, storagegen.UpsertSystemJobScheduleParams{
			ID:                   definition.ID,
			Name:                 definition.Name,
			Kind:                 definition.Kind,
			Queue:                definition.Queue,
			IntervalSeconds:      definition.IntervalSeconds,
			IntervalConfigurable: definition.IntervalConfigurable,
			HistoryPolicy:        systemJobHistoryPolicy(definition.HistoryPolicy),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingsStore) SetSystemJobScheduleInterval(ctx context.Context, id string, intervalSeconds int32) (SystemJobSchedule, error) {
	if intervalSeconds < MinSystemJobScheduleIntervalSeconds {
		return SystemJobSchedule{}, ErrInvalidInput
	}
	_, err := storagegen.New(s.pool).UpdateSystemJobScheduleInterval(ctx, storagegen.UpdateSystemJobScheduleIntervalParams{
		ID:              strings.TrimSpace(id),
		IntervalSeconds: intervalSeconds,
	})
	if err != nil {
		return SystemJobSchedule{}, err
	}
	schedules, err := s.ListSystemJobSchedules(ctx)
	if err != nil {
		return SystemJobSchedule{}, err
	}
	for _, schedule := range schedules {
		if schedule.ID == strings.TrimSpace(id) {
			return schedule, nil
		}
	}
	return SystemJobSchedule{}, ErrNotFound
}

func (s *SettingsStore) ListSystemJobSchedules(ctx context.Context) ([]SystemJobSchedule, error) {
	rows, err := storagegen.New(s.pool).ListSystemJobSchedules(ctx)
	if err != nil {
		return nil, err
	}
	schedules := make([]SystemJobSchedule, 0, len(rows))
	for _, row := range rows {
		schedules = append(schedules, systemJobScheduleFromRow(row))
	}
	return schedules, nil
}

func (s *SettingsStore) SetSystemJobSchedulePaused(ctx context.Context, id string, paused bool) (SystemJobSchedule, error) {
	_, err := storagegen.New(s.pool).UpdateSystemJobSchedulePaused(ctx, storagegen.UpdateSystemJobSchedulePausedParams{
		ID:     strings.TrimSpace(id),
		Paused: paused,
	})
	if err != nil {
		return SystemJobSchedule{}, err
	}
	schedules, err := s.ListSystemJobSchedules(ctx)
	if err != nil {
		return SystemJobSchedule{}, err
	}
	for _, schedule := range schedules {
		if schedule.ID == strings.TrimSpace(id) {
			return schedule, nil
		}
	}
	return SystemJobSchedule{}, ErrNotFound
}

func (s *SettingsStore) SystemJobSchedulePaused(ctx context.Context, id string) bool {
	schedule, err := storagegen.New(s.pool).GetSystemJobSchedule(ctx, strings.TrimSpace(id))
	return err == nil && schedule.Paused
}

func (s *SettingsStore) SystemJobScheduleReady(ctx context.Context, id string) bool {
	ready, err := storagegen.New(s.pool).SystemJobScheduleReady(ctx, strings.TrimSpace(id))
	return err == nil && ready
}

func (s *SettingsStore) UpsertSystemJobExecution(ctx context.Context, input SystemJobExecutionInput) (SystemJobExecution, error) {
	row, err := storagegen.New(s.pool).UpsertSystemJobExecution(ctx, storagegen.UpsertSystemJobExecutionParams{
		RiverJobID:     input.RiverJobID,
		ScheduleID:     nullableText(input.ScheduleID),
		Classification: input.Classification,
		Status:         input.Status,
		Kind:           input.Kind,
		Queue:          input.Queue,
		Attempt:        input.Attempt,
		MaxAttempts:    input.MaxAttempts,
		Priority:       input.Priority,
		Args:           jsonBytes(input.Args, []byte("{}")),
		Metadata:       jsonBytes(input.Metadata, []byte("{}")),
		Errors:         jsonBytes(input.Errors, []byte("[]")),
		InfoMessage:    input.InfoMessage,
		ScheduledAt:    input.ScheduledAt,
		CreatedAt:      input.CreatedAt,
		AttemptedAt:    input.AttemptedAt,
		FinalizedAt:    input.FinalizedAt,
	})
	return systemJobExecutionFromRow(row), err
}

func (s *SettingsStore) UpdateSystemJobExecutionProgress(ctx context.Context, riverJobID int64, progressPercent *int32, label string) (SystemJobExecution, error) {
	row, err := storagegen.New(s.pool).UpdateSystemJobExecutionProgress(ctx, storagegen.UpdateSystemJobExecutionProgressParams{
		RiverJobID:      riverJobID,
		ProgressLabel:   strings.TrimSpace(label),
		ProgressPercent: nullableInt4(progressPercent),
	})
	return systemJobExecutionFromRow(row), err
}

func (s *SettingsStore) ListCurrentOneShotJobExecutions(ctx context.Context, limit int32) ([]SystemJobExecution, error) {
	rows, err := storagegen.New(s.pool).ListCurrentOneShotJobExecutions(ctx, systemJobLimit(limit))
	return systemJobExecutionsFromRows(rows, err)
}

func (s *SettingsStore) ListSystemJobExecutions(ctx context.Context, filters SystemJobExecutionFilters) ([]SystemJobExecution, error) {
	rows, err := storagegen.New(s.pool).ListSystemJobExecutions(ctx, storagegen.ListSystemJobExecutionsParams{
		States:         filters.States,
		ScheduleID:     strings.TrimSpace(filters.ScheduleID),
		Kind:           strings.TrimSpace(filters.Kind),
		Queue:          strings.TrimSpace(filters.Queue),
		IncludeRoutine: filters.IncludeRoutine,
		Before:         filters.Before,
		SearchQuery:    strings.TrimSpace(filters.Query),
		RowLimit:       systemJobLimit(filters.Limit),
	})
	return systemJobExecutionsFromRows(rows, err)
}

func (s *SettingsStore) GetSystemJobExecution(ctx context.Context, riverJobID int64) (SystemJobExecution, error) {
	row, err := storagegen.New(s.pool).GetSystemJobExecution(ctx, riverJobID)
	return systemJobExecutionFromRow(row), err
}

func (s *SettingsStore) CreateSystemJobExecutionLog(ctx context.Context, riverJobID int64, severity string, message string, data map[string]any) (SystemJobExecutionLog, error) {
	payload, err := json.Marshal(nonNilMap(data))
	if err != nil {
		return SystemJobExecutionLog{}, err
	}
	row, err := storagegen.New(s.pool).CreateSystemJobExecutionLog(ctx, storagegen.CreateSystemJobExecutionLogParams{
		RiverJobID: riverJobID,
		Severity:   severity,
		Message:    strings.TrimSpace(message),
		Data:       payload,
	})
	return systemJobExecutionLogFromRow(row), err
}

func (s *SettingsStore) ListSystemJobExecutionLogs(ctx context.Context, riverJobID int64, limit int32) ([]SystemJobExecutionLog, error) {
	rows, err := storagegen.New(s.pool).ListSystemJobExecutionLogs(ctx, storagegen.ListSystemJobExecutionLogsParams{
		RiverJobID: riverJobID,
		RowLimit:   systemJobLimit(limit),
	})
	if err != nil {
		return nil, err
	}
	logs := make([]SystemJobExecutionLog, 0, len(rows))
	for _, row := range rows {
		logs = append(logs, systemJobExecutionLogFromRow(row))
	}
	return logs, nil
}

func (s *SettingsStore) GetSystemJobHistorySettings(ctx context.Context) (SystemJobHistorySettings, error) {
	row, err := storagegen.New(s.pool).GetSystemJobHistorySettings(ctx, storagegen.GetSystemJobHistorySettingsParams{
		RetentionDays:         DefaultSystemJobHistoryRetentionDays,
		RoutineRetentionHours: DefaultRoutineSystemJobRetentionHours,
	})
	return SystemJobHistorySettings{RetentionDays: row.RetentionDays, RoutineRetentionHours: row.RoutineRetentionHours}, err
}

func (s *SettingsStore) UpdateSystemJobHistorySettings(ctx context.Context, settings SystemJobHistorySettings) (SystemJobHistorySettings, error) {
	if settings.RetentionDays < 1 || settings.RetentionDays > 365 ||
		settings.RoutineRetentionHours < 1 || settings.RoutineRetentionHours > 168 {
		return SystemJobHistorySettings{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).UpdateSystemJobHistorySettings(ctx, storagegen.UpdateSystemJobHistorySettingsParams{
		RetentionDays:         settings.RetentionDays,
		RoutineRetentionHours: settings.RoutineRetentionHours,
	})
	if err != nil {
		return SystemJobHistorySettings{}, err
	}
	return SystemJobHistorySettings{RetentionDays: row.RetentionDays, RoutineRetentionHours: row.RoutineRetentionHours}, s.PruneSystemJobExecutions(ctx)
}

func (s *SettingsStore) PruneSystemJobExecutions(ctx context.Context) error {
	settings, err := s.GetSystemJobHistorySettings(ctx)
	if err != nil {
		return err
	}
	return storagegen.New(s.pool).PruneSystemJobExecutions(ctx, storagegen.PruneSystemJobExecutionsParams{
		RetentionDays:         settings.RetentionDays,
		RoutineRetentionHours: settings.RoutineRetentionHours,
	})
}

func systemJobHistoryPolicy(value string) string {
	if strings.TrimSpace(value) == "routine" {
		return "routine"
	}
	return "standard"
}

func jsonBytes(value []byte, fallback []byte) []byte {
	if len(value) == 0 || !json.Valid(value) {
		return fallback
	}
	return value
}

func nonNilMap(value map[string]any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return value
}
