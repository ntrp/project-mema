package storage

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CustomFormat struct {
	ID           uuid.UUID
	Name         string
	IncludeSpecs []CustomFormatSpec
	ExcludeSpecs []CustomFormatSpec
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CustomFormatSpec struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

type CustomFormatInput struct {
	ID           uuid.UUID
	Name         string
	IncludeSpecs []CustomFormatSpec
	ExcludeSpecs []CustomFormatSpec
}

//go:embed trash_custom_formats.json
var trashCustomFormatsJSON []byte

func (s *SettingsStore) ListCustomFormats(ctx context.Context) ([]CustomFormat, error) {
	if err := s.EnsureDefaultCustomFormats(ctx); err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		select id, name, include_specs, exclude_specs, created_at, updated_at
		from app.custom_formats
		order by lower(name)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formats := []CustomFormat{}
	for rows.Next() {
		format, err := scanCustomFormat(rows)
		if err != nil {
			return nil, err
		}
		formats = append(formats, format)
	}
	return formats, rows.Err()
}

func (s *SettingsStore) EnsureDefaultCustomFormats(ctx context.Context) error {
	if err := s.NormalizeDefaultCustomFormatNames(ctx); err != nil {
		return err
	}

	var count int
	if err := s.pool.QueryRow(ctx, `select count(*) from app.custom_formats`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	formats, err := defaultCustomFormats()
	if err != nil {
		return err
	}
	for _, format := range formats {
		if _, err := s.CreateCustomFormat(ctx, format); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingsStore) NormalizeDefaultCustomFormatNames(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		update app.custom_formats
		set name = regexp_replace(name, '^(Radarr|Sonarr) ', ''),
			updated_at = now()
		where name ~ '^(Radarr|Sonarr) '
	`)
	return err
}

func (s *SettingsStore) CreateCustomFormat(ctx context.Context, input CustomFormatInput) (CustomFormat, error) {
	id := input.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	includeSpecs, excludeSpecs, err := marshalCustomFormatSpecs(input)
	if err != nil {
		return CustomFormat{}, err
	}
	return scanCustomFormatRow(s.pool.QueryRow(ctx, `
		insert into app.custom_formats (id, name, include_specs, exclude_specs)
		values ($1, $2, $3::jsonb, $4::jsonb)
		returning id, name, include_specs, exclude_specs, created_at, updated_at
	`, id, input.Name, includeSpecs, excludeSpecs))
}

func (s *SettingsStore) UpdateCustomFormat(ctx context.Context, id uuid.UUID, input CustomFormatInput) (CustomFormat, error) {
	includeSpecs, excludeSpecs, err := marshalCustomFormatSpecs(input)
	if err != nil {
		return CustomFormat{}, err
	}
	return scanCustomFormatRow(s.pool.QueryRow(ctx, `
		update app.custom_formats
		set name = $2, include_specs = $3::jsonb, exclude_specs = $4::jsonb, updated_at = now()
		where id = $1
		returning id, name, include_specs, exclude_specs, created_at, updated_at
	`, id, input.Name, includeSpecs, excludeSpecs))
}

func (s *SettingsStore) DeleteCustomFormat(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.custom_formats where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func defaultCustomFormats() ([]CustomFormatInput, error) {
	var formats []CustomFormatInput
	err := json.Unmarshal(trashCustomFormatsJSON, &formats)
	return formats, err
}

func marshalCustomFormatSpecs(input CustomFormatInput) ([]byte, []byte, error) {
	includeSpecs, err := json.Marshal(input.IncludeSpecs)
	if err != nil {
		return nil, nil, err
	}
	excludeSpecs, err := json.Marshal(input.ExcludeSpecs)
	if err != nil {
		return nil, nil, err
	}
	return includeSpecs, excludeSpecs, nil
}

func scanCustomFormatRow(row pgx.Row) (CustomFormat, error) {
	format, err := scanCustomFormat(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return CustomFormat{}, ErrNotFound
	}
	return format, err
}

func scanCustomFormat(row pgx.Row) (CustomFormat, error) {
	var format CustomFormat
	var includeSpecs []byte
	var excludeSpecs []byte
	err := row.Scan(
		&format.ID,
		&format.Name,
		&includeSpecs,
		&excludeSpecs,
		&format.CreatedAt,
		&format.UpdatedAt,
	)
	if err != nil {
		return CustomFormat{}, err
	}
	if err := json.Unmarshal(includeSpecs, &format.IncludeSpecs); err != nil {
		return CustomFormat{}, err
	}
	if err := json.Unmarshal(excludeSpecs, &format.ExcludeSpecs); err != nil {
		return CustomFormat{}, err
	}
	return format, nil
}
