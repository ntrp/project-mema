package storage

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	storagegen "media-manager/internal/storage/generated"
)

type DLNARendererDeviceOverride struct {
	ID                      uuid.UUID
	RendererUUID            *string
	IPAddress               *string
	ProfileID               string
	DisplayName             string
	Allowed                 bool
	DeliveryPolicyOverrides json.RawMessage
	Notes                   string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type DLNARendererDeviceOverrideInput struct {
	ID                      uuid.UUID
	RendererUUID            *string
	IPAddress               *string
	ProfileID               string
	DisplayName             string
	Allowed                 bool
	DeliveryPolicyOverrides json.RawMessage
	Notes                   string
}

func (s *SettingsStore) ListDLNARendererDeviceOverrides(
	ctx context.Context,
) ([]DLNARendererDeviceOverride, error) {
	rows, err := storagegen.New(s.pool).ListDLNARendererDeviceOverrides(ctx)
	if err != nil {
		return nil, err
	}
	overrides := make([]DLNARendererDeviceOverride, 0, len(rows))
	for _, row := range rows {
		overrides = append(overrides, dlnaRendererDeviceOverrideFromRow(row))
	}
	return overrides, nil
}

func (s *SettingsStore) UpsertDLNARendererDeviceOverride(
	ctx context.Context,
	input DLNARendererDeviceOverrideInput,
) (DLNARendererDeviceOverride, error) {
	normalized, err := normalizeDLNARendererDeviceOverrideInput(input)
	if err != nil {
		return DLNARendererDeviceOverride{}, err
	}
	row, err := storagegen.New(s.pool).UpsertDLNARendererDeviceOverride(ctx, dlnaRendererDeviceOverrideParams(normalized))
	if err != nil {
		return DLNARendererDeviceOverride{}, err
	}
	return dlnaRendererDeviceOverrideFromRow(row), nil
}

func (s *SettingsStore) DeleteDLNARendererDeviceOverride(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteDLNARendererDeviceOverride(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func normalizeDLNARendererDeviceOverrideInput(
	input DLNARendererDeviceOverrideInput,
) (DLNARendererDeviceOverrideInput, error) {
	input.ProfileID = strings.TrimSpace(input.ProfileID)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Notes = strings.TrimSpace(input.Notes)
	input.RendererUUID = cleanOptionalString(input.RendererUUID)
	input.IPAddress = cleanOptionalString(input.IPAddress)
	if input.ID == uuid.Nil {
		input.ID = uuid.New()
	}
	if input.ProfileID == "" || (input.RendererUUID == nil && input.IPAddress == nil) {
		return DLNARendererDeviceOverrideInput{}, ErrInvalidInput
	}
	if input.IPAddress != nil && net.ParseIP(*input.IPAddress) == nil {
		return DLNARendererDeviceOverrideInput{}, ErrInvalidInput
	}
	if !rendererJSONObject(input.DeliveryPolicyOverrides) {
		return DLNARendererDeviceOverrideInput{}, ErrInvalidInput
	}
	return input, nil
}

func dlnaRendererDeviceOverrideParams(
	input DLNARendererDeviceOverrideInput,
) storagegen.UpsertDLNARendererDeviceOverrideParams {
	return storagegen.UpsertDLNARendererDeviceOverrideParams{
		ID:                      input.ID,
		RendererUuid:            dlnaOptionalText(input.RendererUUID),
		IpAddress:               dlnaOptionalText(input.IPAddress),
		ProfileID:               input.ProfileID,
		DisplayName:             input.DisplayName,
		Allowed:                 input.Allowed,
		DeliveryPolicyOverrides: input.DeliveryPolicyOverrides,
		Notes:                   input.Notes,
	}
}

func dlnaRendererDeviceOverrideFromRow(row storagegen.AppDlnaRendererDeviceOverride) DLNARendererDeviceOverride {
	return DLNARendererDeviceOverride{
		ID: row.ID, RendererUUID: textPtr(row.RendererUuid), IPAddress: textPtr(row.IpAddress),
		ProfileID: row.ProfileID, DisplayName: row.DisplayName, Allowed: row.Allowed,
		DeliveryPolicyOverrides: row.DeliveryPolicyOverrides, Notes: row.Notes,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cleaned := strings.TrimSpace(*value)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}

func dlnaOptionalText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}
