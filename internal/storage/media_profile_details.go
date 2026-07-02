package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type mediaProfileQueryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (s *SettingsStore) populateMediaProfile(ctx context.Context, profile *MediaProfile) error {
	qualities, err := loadMediaProfileQualities(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	languageScores, err := loadMediaProfileLanguages(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	scores, err := loadMediaProfileCustomFormats(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	profile.QualityIDs = qualities
	profile.TargetLanguageScores = languageScores
	profile.TargetLanguages = languageIDsFromScores(languageScores)
	profile.CustomFormatScores = scores
	return nil
}

func loadMediaProfileQualities(
	ctx context.Context,
	q mediaProfileQueryer,
	profileID string,
) ([]string, error) {
	rows, err := q.Query(ctx, `
		select quality_id
		from app.media_profile_qualities
		where profile_id = $1
		order by sort_order, quality_id
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qualityIDs := []string{}
	for rows.Next() {
		var qualityID string
		if err := rows.Scan(&qualityID); err != nil {
			return nil, err
		}
		qualityIDs = append(qualityIDs, qualityID)
	}
	return qualityIDs, rows.Err()
}

func loadMediaProfileLanguages(
	ctx context.Context,
	q mediaProfileQueryer,
	profileID string,
) ([]MediaProfileLanguageScore, error) {
	rows, err := q.Query(ctx, `
		select language_id, score, required
		from app.media_profile_languages
		where profile_id = $1
		order by language_id
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := []MediaProfileLanguageScore{}
	for rows.Next() {
		var score MediaProfileLanguageScore
		if err := rows.Scan(&score.LanguageID, &score.Score, &score.Required); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func loadMediaProfileCustomFormats(
	ctx context.Context,
	q mediaProfileQueryer,
	profileID string,
) ([]MediaProfileCustomFormatScore, error) {
	rows, err := q.Query(ctx, `
		select pcf.custom_format_id, pcf.score
		from app.media_profile_custom_formats pcf
		join app.custom_formats cf on cf.id = pcf.custom_format_id
		where pcf.profile_id = $1
		order by lower(cf.name), pcf.custom_format_id
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := []MediaProfileCustomFormatScore{}
	for rows.Next() {
		var score MediaProfileCustomFormatScore
		if err := rows.Scan(&score.CustomFormatID, &score.Score); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func replaceMediaProfileCustomFormats(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	scores []MediaProfileCustomFormatScore,
) error {
	if _, err := q.Exec(ctx, `delete from app.media_profile_custom_formats where profile_id = $1`, profileID); err != nil {
		return err
	}
	for _, score := range scores {
		if _, err := q.Exec(ctx, `
			insert into app.media_profile_custom_formats (profile_id, custom_format_id, score)
			values ($1, $2, $3)
		`, profileID, score.CustomFormatID, score.Score); err != nil {
			return normalizeMediaProfileWriteError(err)
		}
	}
	return nil
}

func replaceMediaProfileLanguages(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	scores []MediaProfileLanguageScore,
) error {
	if _, err := q.Exec(ctx, `delete from app.media_profile_languages where profile_id = $1`, profileID); err != nil {
		return err
	}
	for _, score := range scores {
		if _, err := q.Exec(ctx, `
			insert into app.media_profile_languages (profile_id, language_id, score, required)
			values ($1, $2, $3, $4)
		`, profileID, score.LanguageID, score.Score, score.Required); err != nil {
			return normalizeMediaProfileWriteError(err)
		}
	}
	return nil
}

func normalizeCustomFormatScores(values []MediaProfileCustomFormatScore) []MediaProfileCustomFormatScore {
	seen := map[uuid.UUID]int{}
	scores := []MediaProfileCustomFormatScore{}
	for _, value := range values {
		if value.CustomFormatID == uuid.Nil {
			continue
		}
		if index, ok := seen[value.CustomFormatID]; ok {
			scores[index].Score = value.Score
			continue
		}
		seen[value.CustomFormatID] = len(scores)
		scores = append(scores, value)
	}
	return scores
}
