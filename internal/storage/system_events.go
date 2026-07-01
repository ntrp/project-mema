package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	return scanSystemEvent(s.pool.QueryRow(ctx, `
		insert into app.system_events (id, severity, category, message, data)
		values ($1, $2, $3, $4, $5)
		returning id, severity, category, message, data, created_at
	`, uuid.New(), input.Severity, input.Category, input.Message, payload))
}

func (s *SettingsStore) ListSystemEvents(ctx context.Context, limit int, before *time.Time) ([]SystemEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	if err := s.PruneSystemEvents(ctx); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		select id, severity, category, message, data, created_at
		from app.system_events
		where ($2::timestamptz is null or created_at < $2)
		order by created_at desc
		limit $1
	`, limit, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []SystemEvent{}
	for rows.Next() {
		event, err := scanSystemEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (s *SettingsStore) DeleteSystemEvent(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.system_events where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ClearSystemEvents(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `delete from app.system_events`)
	return err
}

func (s *SettingsStore) GetSystemEventSettings(ctx context.Context) (SystemEventSettings, error) {
	return scanSystemEventSettings(s.pool.QueryRow(ctx, `
		insert into app.system_event_settings (id, retention_days)
		values (true, $1)
		on conflict (id) do update set id = excluded.id
		returning retention_days, created_at, updated_at
	`, DefaultSystemEventRetentionDays))
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
	settings, err := scanSystemEventSettings(s.pool.QueryRow(ctx, `
		insert into app.system_event_settings (id, retention_days)
		values (true, $1)
		on conflict (id) do update
		set retention_days = excluded.retention_days,
			updated_at = now()
		returning retention_days, created_at, updated_at
	`, input.RetentionDays))
	if err != nil {
		return SystemEventSettings{}, err
	}
	return settings, s.PruneSystemEvents(ctx)
}

func (s *SettingsStore) PruneSystemEvents(ctx context.Context) error {
	settings, err := s.GetSystemEventSettings(ctx)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		delete from app.system_events
		where created_at < now() - ($1::int * interval '1 day')
	`, settings.RetentionDays)
	return err
}

func scanSystemEvent(row interface {
	Scan(dest ...any) error
}) (SystemEvent, error) {
	var event SystemEvent
	var payload []byte
	err := row.Scan(&event.ID, &event.Severity, &event.Category, &event.Message, &payload, &event.CreatedAt)
	if err != nil {
		return SystemEvent{}, err
	}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &event.Data)
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}
	return event, nil
}

func scanSystemEventSettings(row interface {
	Scan(dest ...any) error
}) (SystemEventSettings, error) {
	var settings SystemEventSettings
	err := row.Scan(&settings.RetentionDays, &settings.CreatedAt, &settings.UpdatedAt)
	return settings, err
}
