package storage

import (
	"context"
	"errors"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *SettingsStore) ListMediaProfiles(ctx context.Context) ([]MediaProfile, error) {
	rows, err := storagegen.New(s.pool).ListMediaProfiles(ctx)
	if err != nil {
		return nil, err
	}

	profiles := make([]MediaProfile, 0, len(rows))
	for _, row := range rows {
		profile := mediaProfileFromRow(row)
		if err := s.populateMediaProfile(ctx, &profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (s *SettingsStore) MediaProfileExists(ctx context.Context, id string) (bool, error) {
	return storagegen.New(s.pool).MediaProfileExists(ctx, id)
}

func (s *SettingsStore) GetMediaProfile(ctx context.Context, id string) (MediaProfile, error) {
	return s.getMediaProfile(ctx, strings.TrimSpace(id))
}

func (s *SettingsStore) CreateMediaProfile(ctx context.Context, input MediaProfileInput) (MediaProfile, error) {
	name := normalizeMediaProfileName(input.Name)
	if name == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	qualityIDs, err := normalizeProfileQualityIDs(input.QualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	normalized, err := normalizeMediaProfileInput(input, qualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	id := mediaProfileSlug(name)
	if id == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	return s.saveMediaProfile(ctx, id, name, qualityIDs, normalized, true)
}

func (s *SettingsStore) UpdateMediaProfile(ctx context.Context, id string, input MediaProfileInput) (MediaProfile, error) {
	name := normalizeMediaProfileName(input.Name)
	if name == "" || strings.TrimSpace(id) == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	qualityIDs, err := normalizeProfileQualityIDs(input.QualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	normalized, err := normalizeMediaProfileInput(input, qualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	return s.saveMediaProfile(ctx, strings.TrimSpace(id), name, qualityIDs, normalized, false)
}

func (s *SettingsStore) DeleteMediaProfile(ctx context.Context, id string) error {
	rows, err := storagegen.New(s.pool).DeleteMediaProfile(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) saveMediaProfile(
	ctx context.Context,
	id string,
	name string,
	qualityIDs []string,
	input MediaProfileInput,
	create bool,
) (MediaProfile, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaProfile{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	q := storagegen.New(tx)

	if input.IsDefault {
		if err := q.ClearDefaultMediaProfiles(ctx); err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
	}
	if create {
		if err := q.CreateMediaProfile(ctx, mediaProfileParams(id, name, input)); err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
	} else {
		rows, err := q.UpdateMediaProfile(ctx, mediaProfileUpdateParams(id, name, input))
		if err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
		if rows == 0 {
			return MediaProfile{}, ErrNotFound
		}
	}

	if err := replaceMediaProfileQualities(ctx, tx, id, qualityIDs); err != nil {
		return MediaProfile{}, err
	}
	if err := replaceMediaProfileTargets(ctx, tx, id, input); err != nil {
		return MediaProfile{}, normalizeMediaProfileWriteError(err)
	}
	if err := replaceMediaProfileCustomFormats(ctx, tx, id, input.CustomFormatScores); err != nil {
		return MediaProfile{}, normalizeMediaProfileWriteError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaProfile{}, normalizeMediaProfileWriteError(err)
	}
	return s.getMediaProfile(ctx, id)
}

func (s *SettingsStore) getMediaProfile(ctx context.Context, id string) (MediaProfile, error) {
	row, err := storagegen.New(s.pool).GetMediaProfile(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaProfile{}, ErrNotFound
	}
	if err != nil {
		return MediaProfile{}, err
	}
	profile := mediaProfileFromRow(row)
	return profile, s.populateMediaProfile(ctx, &profile)
}

type mediaProfileQuerier = storagegen.DBTX

func replaceMediaProfileQualities(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	qualityIDs []string,
) error {
	queries := storagegen.New(q)
	if err := queries.ClearMediaProfileQualities(ctx, profileID); err != nil {
		return err
	}
	definitions := QualitySizeDefinitionMap()
	for index, qualityID := range qualityIDs {
		definition := definitions[qualityID]
		sortOrder := int32(index)
		if definition.ID != "" {
			sortOrder = definition.SortOrder
		}
		if err := queries.AddMediaProfileQuality(ctx, storagegen.AddMediaProfileQualityParams{
			ProfileID: profileID,
			QualityID: qualityID,
			SortOrder: sortOrder,
		}); err != nil {
			return err
		}
	}
	return nil
}

func mediaProfileParams(id string, name string, input MediaProfileInput) storagegen.CreateMediaProfileParams {
	return storagegen.CreateMediaProfileParams{
		ID:                                id,
		Name:                              name,
		IsDefault:                         input.IsDefault,
		FinalContainer:                    input.FinalContainer,
		UpgradesAllowed:                   input.UpgradesAllowed,
		UpgradeUntilQualityID:             textValue(input.UpgradeUntilQualityID),
		MinimumCustomFormatScore:          input.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     input.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: input.MinimumCustomFormatScoreIncrement,
		RemoveUnwantedAudio:               input.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:         input.AudioLossyTranscodePolicy,
		RemoveUnwantedSubtitles:           input.RemoveUnwantedSubtitles,
		SubtitlePreferredMode:             input.SubtitlePreferredMode,
		AllowSubtitleReleaseFallback:      input.AllowSubtitleReleaseFallback,
		PreferredProtocol:                 input.PreferredProtocol,
		SeriesPackPreference:              input.SeriesPackPreference,
	}
}

func mediaProfileUpdateParams(id string, name string, input MediaProfileInput) storagegen.UpdateMediaProfileParams {
	return storagegen.UpdateMediaProfileParams{
		ID:                                id,
		Name:                              name,
		IsDefault:                         input.IsDefault,
		FinalContainer:                    input.FinalContainer,
		UpgradesAllowed:                   input.UpgradesAllowed,
		UpgradeUntilQualityID:             textValue(input.UpgradeUntilQualityID),
		MinimumCustomFormatScore:          input.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     input.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: input.MinimumCustomFormatScoreIncrement,
		RemoveUnwantedAudio:               input.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:         input.AudioLossyTranscodePolicy,
		RemoveUnwantedSubtitles:           input.RemoveUnwantedSubtitles,
		SubtitlePreferredMode:             input.SubtitlePreferredMode,
		AllowSubtitleReleaseFallback:      input.AllowSubtitleReleaseFallback,
		PreferredProtocol:                 input.PreferredProtocol,
		SeriesPackPreference:              input.SeriesPackPreference,
	}
}

func mediaProfileFromRow(row storagegen.AppMediaProfile) MediaProfile {
	return MediaProfile{
		ID:                                row.ID,
		Name:                              row.Name,
		IsDefault:                         row.IsDefault,
		FinalContainer:                    row.FinalContainer,
		UpgradesAllowed:                   row.UpgradesAllowed,
		UpgradeUntilQualityID:             textPtr(row.UpgradeUntilQualityID),
		MinimumCustomFormatScore:          row.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     row.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: row.MinimumCustomFormatScoreIncrement,
		RemoveUnwantedAudio:               row.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:         row.AudioLossyTranscodePolicy,
		RemoveUnwantedSubtitles:           row.RemoveUnwantedSubtitles,
		SubtitlePreferredMode:             row.SubtitlePreferredMode,
		AllowSubtitleReleaseFallback:      row.AllowSubtitleReleaseFallback,
		PreferredProtocol:                 row.PreferredProtocol,
		SeriesPackPreference:              row.SeriesPackPreference,
		CreatedAt:                         row.CreatedAt,
		UpdatedAt:                         row.UpdatedAt,
	}
}

func normalizeMediaProfileWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23503") {
		return ErrInvalidInput
	}
	return err
}
