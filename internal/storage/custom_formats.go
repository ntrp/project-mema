package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

type CustomFormat struct {
	ID                      uuid.UUID
	Name                    string
	IncludeInRenameTemplate bool
	IncludeSpecs            []CustomFormatSpec
	ExcludeSpecs            []CustomFormatSpec
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CustomFormatSpec struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

type CustomFormatInput struct {
	ID                      uuid.UUID
	Name                    string
	IncludeInRenameTemplate bool
	IncludeSpecs            []CustomFormatSpec
	ExcludeSpecs            []CustomFormatSpec
}

func (s *SettingsStore) ListCustomFormats(ctx context.Context) ([]CustomFormat, error) {
	rows, err := storagegen.New(s.pool).ListCustomFormats(ctx)
	if err != nil {
		return nil, err
	}

	formats := make([]CustomFormat, 0, len(rows))
	for _, row := range rows {
		format, err := customFormatFromRow(row)
		if err != nil {
			return nil, err
		}
		formats = append(formats, format)
	}
	return formats, nil
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
	row, err := storagegen.New(s.pool).CreateCustomFormat(ctx, storagegen.CreateCustomFormatParams{
		ID:                      id,
		Name:                    input.Name,
		IncludeInRenameTemplate: input.IncludeInRenameTemplate,
		IncludeSpecs:            includeSpecs,
		ExcludeSpecs:            excludeSpecs,
	})
	if err != nil {
		return CustomFormat{}, err
	}
	return customFormatFromRow(row)
}

func (s *SettingsStore) UpdateCustomFormat(ctx context.Context, id uuid.UUID, input CustomFormatInput) (CustomFormat, error) {
	includeSpecs, excludeSpecs, err := marshalCustomFormatSpecs(input)
	if err != nil {
		return CustomFormat{}, err
	}
	row, err := storagegen.New(s.pool).UpdateCustomFormat(ctx, storagegen.UpdateCustomFormatParams{
		ID:                      id,
		Name:                    input.Name,
		IncludeInRenameTemplate: input.IncludeInRenameTemplate,
		IncludeSpecs:            includeSpecs,
		ExcludeSpecs:            excludeSpecs,
	})
	if err != nil {
		return CustomFormat{}, normalizeCustomFormatWriteError(err)
	}
	return customFormatFromRow(row)
}

func (s *SettingsStore) DeleteCustomFormat(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := storagegen.New(s.pool).DeleteCustomFormat(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
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

func normalizeCustomFormatWriteError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func customFormatFromRow(row storagegen.AppCustomFormat) (CustomFormat, error) {
	format := CustomFormat{
		ID:                      row.ID,
		Name:                    row.Name,
		IncludeInRenameTemplate: row.IncludeInRenameTemplate,
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               row.UpdatedAt,
	}
	if err := json.Unmarshal(row.IncludeSpecs, &format.IncludeSpecs); err != nil {
		return CustomFormat{}, err
	}
	if err := json.Unmarshal(row.ExcludeSpecs, &format.ExcludeSpecs); err != nil {
		return CustomFormat{}, err
	}
	return format, nil
}
