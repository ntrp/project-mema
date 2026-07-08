package storage

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

type DLNARendererProfile struct {
	ID               string
	Name             string
	Vendor           string
	DeviceClass      string
	Source           string
	SourceVersion    int32
	Customized       bool
	Enabled          bool
	Priority         int32
	IconKey          string
	Notes            string
	MatchRules       json.RawMessage
	CapabilityRules  json.RawMessage
	DeliverySettings json.RawMessage
	DLNAFlags        json.RawMessage
	SubtitleRules    json.RawMessage
	ArtworkRules     json.RawMessage
	MetadataRules    json.RawMessage
	Quirks           json.RawMessage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type DLNARendererProfileInput struct {
	Name             string
	Vendor           string
	DeviceClass      string
	Enabled          bool
	Priority         int32
	IconKey          string
	Notes            string
	MatchRules       json.RawMessage
	CapabilityRules  json.RawMessage
	DeliverySettings json.RawMessage
	DLNAFlags        json.RawMessage
	SubtitleRules    json.RawMessage
	ArtworkRules     json.RawMessage
	MetadataRules    json.RawMessage
	Quirks           json.RawMessage
}

func (s *SettingsStore) ListDLNARendererProfiles(ctx context.Context) ([]DLNARendererProfile, error) {
	rows, err := storagegen.New(s.pool).ListDLNARendererProfiles(ctx)
	if err != nil {
		return nil, err
	}
	profiles := make([]DLNARendererProfile, 0, len(rows))
	for _, row := range rows {
		profiles = append(profiles, dlnaRendererProfileFromRow(row))
	}
	return profiles, nil
}

func (s *SettingsStore) GetDLNARendererProfile(ctx context.Context, id string) (DLNARendererProfile, error) {
	row, err := storagegen.New(s.pool).GetDLNARendererProfile(ctx, strings.TrimSpace(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return DLNARendererProfile{}, ErrNotFound
	}
	if err != nil {
		return DLNARendererProfile{}, err
	}
	return dlnaRendererProfileFromRow(row), nil
}

func (s *SettingsStore) UpdateDLNARendererProfile(
	ctx context.Context,
	id string,
	input DLNARendererProfileInput,
) (DLNARendererProfile, error) {
	normalized, err := normalizeDLNARendererProfileInput(input)
	if err != nil {
		return DLNARendererProfile{}, err
	}
	row, err := storagegen.New(s.pool).UpdateDLNARendererProfile(ctx, dlnaRendererProfileParams(id, normalized))
	if errors.Is(err, pgx.ErrNoRows) {
		return DLNARendererProfile{}, ErrNotFound
	}
	if err != nil {
		return DLNARendererProfile{}, err
	}
	return dlnaRendererProfileFromRow(row), nil
}

func (s *SettingsStore) ResetDLNARendererProfile(ctx context.Context, id string) (DLNARendererProfile, error) {
	row, err := storagegen.New(s.pool).ResetDLNARendererProfile(ctx, strings.TrimSpace(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return DLNARendererProfile{}, ErrNotFound
	}
	if err != nil {
		return DLNARendererProfile{}, err
	}
	return dlnaRendererProfileFromRow(row), nil
}

func normalizeDLNARendererProfileInput(input DLNARendererProfileInput) (DLNARendererProfileInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Vendor = strings.TrimSpace(input.Vendor)
	input.DeviceClass = strings.TrimSpace(input.DeviceClass)
	input.IconKey = strings.TrimSpace(input.IconKey)
	input.Notes = strings.TrimSpace(input.Notes)
	if input.Name == "" || input.DeviceClass == "" || input.Priority < 0 || input.Priority > 1000 {
		return DLNARendererProfileInput{}, ErrInvalidInput
	}
	if err := validateRendererProfileJSON(input); err != nil {
		return DLNARendererProfileInput{}, err
	}
	return input, nil
}

func validateRendererProfileJSON(input DLNARendererProfileInput) error {
	values := []json.RawMessage{
		input.MatchRules, input.CapabilityRules, input.DeliverySettings, input.DLNAFlags,
		input.SubtitleRules, input.ArtworkRules, input.MetadataRules, input.Quirks,
	}
	for _, value := range values {
		if !rendererJSONObject(value) {
			return ErrInvalidInput
		}
	}
	return validateRendererMatchRegexes(input.MatchRules)
}

func dlnaRendererProfileParams(
	id string,
	input DLNARendererProfileInput,
) storagegen.UpdateDLNARendererProfileParams {
	return storagegen.UpdateDLNARendererProfileParams{
		ID:               strings.TrimSpace(id),
		Name:             input.Name,
		Vendor:           input.Vendor,
		DeviceClass:      input.DeviceClass,
		Enabled:          input.Enabled,
		Priority:         input.Priority,
		IconKey:          input.IconKey,
		Notes:            input.Notes,
		MatchRules:       input.MatchRules,
		CapabilityRules:  input.CapabilityRules,
		DeliverySettings: input.DeliverySettings,
		DlnaFlags:        input.DLNAFlags,
		SubtitleRules:    input.SubtitleRules,
		ArtworkRules:     input.ArtworkRules,
		MetadataRules:    input.MetadataRules,
		Quirks:           input.Quirks,
	}
}

func dlnaRendererProfileFromRow(row storagegen.AppDlnaRendererProfile) DLNARendererProfile {
	return DLNARendererProfile{
		ID: row.ID, Name: row.Name, Vendor: row.Vendor, DeviceClass: row.DeviceClass,
		Source: row.Source, SourceVersion: row.SourceVersion, Customized: row.Customized,
		Enabled: row.Enabled, Priority: row.Priority, IconKey: row.IconKey, Notes: row.Notes,
		MatchRules: row.MatchRules, CapabilityRules: row.CapabilityRules,
		DeliverySettings: row.DeliverySettings, DLNAFlags: row.DlnaFlags,
		SubtitleRules: row.SubtitleRules, ArtworkRules: row.ArtworkRules,
		MetadataRules: row.MetadataRules, Quirks: row.Quirks,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}
