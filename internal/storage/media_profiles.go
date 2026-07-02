package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *SettingsStore) ListMediaProfiles(ctx context.Context) ([]MediaProfile, error) {
	rows, err := s.pool.Query(ctx, `
		select
			id,
			name,
			upgrades_allowed,
			upgrade_until_quality_id,
			minimum_custom_format_score,
			upgrade_until_custom_format_score,
			minimum_custom_format_score_increment,
			remove_non_enabled_languages,
			preferred_protocol,
			series_pack_preference,
			created_at,
			updated_at
		from app.media_profiles
		order by lower(name)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := []MediaProfile{}
	for rows.Next() {
		profile, err := scanMediaProfileBase(rows)
		if err != nil {
			return nil, err
		}
		if err := s.populateMediaProfile(ctx, &profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, rows.Err()
}

func (s *SettingsStore) MediaProfileExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `select exists(select 1 from app.media_profiles where id = $1)`, id).Scan(&exists)
	return exists, err
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
	tag, err := s.pool.Exec(ctx, `delete from app.media_profiles where id = $1`, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
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

	if create {
		if _, err := tx.Exec(ctx, `
			insert into app.media_profiles (
				id,
				name,
				upgrades_allowed,
				upgrade_until_quality_id,
				minimum_custom_format_score,
				upgrade_until_custom_format_score,
				minimum_custom_format_score_increment,
				remove_non_enabled_languages,
				preferred_protocol,
				series_pack_preference
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`,
			id,
			name,
			input.UpgradesAllowed,
			input.UpgradeUntilQualityID,
			input.MinimumCustomFormatScore,
			input.UpgradeUntilCustomFormatScore,
			input.MinimumCustomFormatScoreIncrement,
			input.RemoveNonEnabledLanguages,
			input.PreferredProtocol,
			input.SeriesPackPreference,
		); err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
	} else {
		tag, err := tx.Exec(ctx, `
			update app.media_profiles
			set name = $2,
				upgrades_allowed = $3,
				upgrade_until_quality_id = $4,
				minimum_custom_format_score = $5,
				upgrade_until_custom_format_score = $6,
				minimum_custom_format_score_increment = $7,
				remove_non_enabled_languages = $8,
				preferred_protocol = $9,
				series_pack_preference = $10,
				updated_at = now()
			where id = $1
		`,
			id,
			name,
			input.UpgradesAllowed,
			input.UpgradeUntilQualityID,
			input.MinimumCustomFormatScore,
			input.UpgradeUntilCustomFormatScore,
			input.MinimumCustomFormatScoreIncrement,
			input.RemoveNonEnabledLanguages,
			input.PreferredProtocol,
			input.SeriesPackPreference,
		)
		if err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
		if tag.RowsAffected() == 0 {
			return MediaProfile{}, ErrNotFound
		}
	}

	if err := replaceMediaProfileQualities(ctx, tx, id, qualityIDs); err != nil {
		return MediaProfile{}, err
	}
	if err := replaceMediaProfileLanguages(ctx, tx, id, input.TargetLanguageScores); err != nil {
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
	profile, err := scanMediaProfileBase(s.pool.QueryRow(ctx, `
		select
			id,
			name,
			upgrades_allowed,
			upgrade_until_quality_id,
			minimum_custom_format_score,
			upgrade_until_custom_format_score,
			minimum_custom_format_score_increment,
			remove_non_enabled_languages,
			preferred_protocol,
			series_pack_preference,
			created_at,
			updated_at
		from app.media_profiles
		where id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaProfile{}, ErrNotFound
	}
	if err != nil {
		return MediaProfile{}, err
	}
	return profile, s.populateMediaProfile(ctx, &profile)
}

type mediaProfileQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func replaceMediaProfileQualities(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	qualityIDs []string,
) error {
	if _, err := q.Exec(ctx, `delete from app.media_profile_qualities where profile_id = $1`, profileID); err != nil {
		return err
	}
	definitions := QualitySizeDefinitionMap()
	for index, qualityID := range qualityIDs {
		definition := definitions[qualityID]
		sortOrder := int32(index)
		if definition.ID != "" {
			sortOrder = definition.SortOrder
		}
		if _, err := q.Exec(ctx, `
			insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
			values ($1, $2, $3)
		`, profileID, qualityID, sortOrder); err != nil {
			return err
		}
	}
	return nil
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
