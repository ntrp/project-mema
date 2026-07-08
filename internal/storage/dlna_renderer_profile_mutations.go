package storage

import (
	"context"
	"errors"
	"regexp"
	"strings"

	storagegen "media-manager/internal/storage/generated"
)

var dlnaRendererProfileIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,79}$`)

func (s *SettingsStore) CreateDLNARendererProfile(
	ctx context.Context,
	id string,
	input DLNARendererProfileInput,
) (DLNARendererProfile, error) {
	id = strings.TrimSpace(id)
	if !dlnaRendererProfileIDPattern.MatchString(id) {
		return DLNARendererProfile{}, ErrInvalidInput
	}
	normalized, err := normalizeDLNARendererProfileInput(input)
	if err != nil {
		return DLNARendererProfile{}, err
	}
	row, err := storagegen.New(s.pool).CreateDLNARendererProfile(ctx, dlnaRendererProfileCreateParams(id, normalized))
	if err != nil {
		return DLNARendererProfile{}, normalizeMediaProfileWriteError(err)
	}
	return dlnaRendererProfileFromRow(row), nil
}

func (s *SettingsStore) DeleteDLNARendererProfile(ctx context.Context, id string) error {
	rows, err := storagegen.New(s.pool).DeleteDLNARendererProfile(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) CloneDLNARendererProfile(
	ctx context.Context,
	sourceID string,
	targetID string,
	name string,
) (DLNARendererProfile, error) {
	source, err := s.GetDLNARendererProfile(ctx, sourceID)
	if err != nil {
		return DLNARendererProfile{}, err
	}
	input := dlnaRendererProfileInputFromProfile(source)
	input.Name = strings.TrimSpace(name)
	if input.Name == "" {
		input.Name = source.Name + " Copy"
	}
	return s.CreateDLNARendererProfile(ctx, targetID, input)
}

func (s *SettingsStore) ImportDLNARendererProfile(
	ctx context.Context,
	id string,
	input DLNARendererProfileInput,
) (DLNARendererProfile, error) {
	_, err := s.GetDLNARendererProfile(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return s.CreateDLNARendererProfile(ctx, id, input)
	}
	if err != nil {
		return DLNARendererProfile{}, err
	}
	return s.UpdateDLNARendererProfile(ctx, id, input)
}

func dlnaRendererProfileCreateParams(
	id string,
	input DLNARendererProfileInput,
) storagegen.CreateDLNARendererProfileParams {
	return storagegen.CreateDLNARendererProfileParams{
		ID:               id,
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

func dlnaRendererProfileInputFromProfile(profile DLNARendererProfile) DLNARendererProfileInput {
	return DLNARendererProfileInput{
		Name:             profile.Name,
		Vendor:           profile.Vendor,
		DeviceClass:      profile.DeviceClass,
		Enabled:          profile.Enabled,
		Priority:         profile.Priority,
		IconKey:          profile.IconKey,
		Notes:            profile.Notes,
		MatchRules:       profile.MatchRules,
		CapabilityRules:  profile.CapabilityRules,
		DeliverySettings: profile.DeliverySettings,
		DLNAFlags:        profile.DLNAFlags,
		SubtitleRules:    profile.SubtitleRules,
		ArtworkRules:     profile.ArtworkRules,
		MetadataRules:    profile.MetadataRules,
		Quirks:           profile.Quirks,
	}
}
