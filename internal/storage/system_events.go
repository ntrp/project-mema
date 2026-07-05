package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

const DefaultSystemEventRetentionDays = 7

type SystemEvent struct {
	ID        uuid.UUID
	Severity  string
	Category  string
	Message   string
	Data      map[string]any
	CreatedAt time.Time
}

type SystemEventInput struct {
	Severity string
	Category string
	Message  string
	Data     map[string]any
}

type SystemEventSettings struct {
	RetentionDays int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SystemEventSettingsInput struct {
	RetentionDays int32
}

func (s *SettingsStore) CreateSystemEvent(ctx context.Context, input SystemEventInput) (SystemEvent, error) {
	if input.Severity == "" || input.Category == "" || input.Message == "" {
		return SystemEvent{}, ErrInvalidInput
	}
	if err := s.PruneSystemEvents(ctx); err != nil {
		return SystemEvent{}, err
	}
	data := input.Data
	if data == nil {
		data = map[string]any{}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return SystemEvent{}, err
	}
	row, err := storagegen.New(s.pool).CreateSystemEvent(ctx, storagegen.CreateSystemEventParams{
		ID:       uuid.New(),
		Severity: input.Severity,
		Category: input.Category,
		Message:  input.Message,
		Data:     payload,
	})
	if err != nil {
		return SystemEvent{}, err
	}
	return systemEventFromRow(row)
}

func (s *SettingsStore) ListSystemEvents(ctx context.Context, limit int, before *time.Time) ([]SystemEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	if err := s.PruneSystemEvents(ctx); err != nil {
		return nil, err
	}
	rows, err := storagegen.New(s.pool).ListSystemEvents(ctx, storagegen.ListSystemEventsParams{
		Limit:  int32(limit),
		Before: before,
	})
	if err != nil {
		return nil, err
	}
	events := make([]SystemEvent, 0, len(rows))
	for _, row := range rows {
		event, err := systemEventFromRow(row)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *SettingsStore) DeleteSystemEvent(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := storagegen.New(s.pool).DeleteSystemEvent(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ClearSystemEvents(ctx context.Context) error {
	return storagegen.New(s.pool).ClearSystemEvents(ctx)
}

func (s *SettingsStore) GetSystemEventSettings(ctx context.Context) (SystemEventSettings, error) {
	row, err := storagegen.New(s.pool).GetSystemEventSettings(ctx, DefaultSystemEventRetentionDays)
	if err != nil {
		return SystemEventSettings{}, err
	}
	return systemEventSettingsFromGetRow(row), nil
}

func (s *SettingsStore) UpdateSystemEventSettings(
	ctx context.Context,
	input SystemEventSettingsInput,
) (SystemEventSettings, error) {
	if input.RetentionDays == 0 {
		input.RetentionDays = DefaultSystemEventRetentionDays
	}
	if input.RetentionDays < 1 || input.RetentionDays > 365 {
		return SystemEventSettings{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).UpdateSystemEventSettings(ctx, input.RetentionDays)
	if err != nil {
		return SystemEventSettings{}, err
	}
	settings := systemEventSettingsFromUpdateRow(row)
	return settings, s.PruneSystemEvents(ctx)
}

func (s *SettingsStore) PruneSystemEvents(ctx context.Context) error {
	settings, err := s.GetSystemEventSettings(ctx)
	if err != nil {
		return err
	}
	return storagegen.New(s.pool).PruneSystemEvents(ctx, settings.RetentionDays)
}

func systemEventFromRow(row storagegen.AppSystemEvent) (SystemEvent, error) {
	event := SystemEvent{
		ID:        row.ID,
		Severity:  row.Severity,
		Category:  row.Category,
		Message:   row.Message,
		CreatedAt: row.CreatedAt,
	}
	if len(row.Data) > 0 {
		_ = json.Unmarshal(row.Data, &event.Data)
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}
	return event, nil
}

func systemEventSettingsFromGetRow(row storagegen.GetSystemEventSettingsRow) SystemEventSettings {
	return SystemEventSettings{
		RetentionDays: row.RetentionDays,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func systemEventSettingsFromUpdateRow(row storagegen.UpdateSystemEventSettingsRow) SystemEventSettings {
	return SystemEventSettings{
		RetentionDays: row.RetentionDays,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
