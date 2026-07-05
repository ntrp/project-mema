package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

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
	q storagegen.DBTX,
	profileID string,
) ([]string, error) {
	return storagegen.New(q).ListMediaProfileQualities(ctx, profileID)
}

func loadMediaProfileLanguages(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]MediaProfileLanguageScore, error) {
	rows, err := storagegen.New(q).ListMediaProfileLanguages(ctx, profileID)
	if err != nil {
		return nil, err
	}
	scores := make([]MediaProfileLanguageScore, 0, len(rows))
	for _, row := range rows {
		scores = append(scores, MediaProfileLanguageScore{
			LanguageID: row.LanguageID,
			Score:      row.Score,
			Required:   row.Required,
		})
	}
	return scores, nil
}

func loadMediaProfileCustomFormats(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]MediaProfileCustomFormatScore, error) {
	rows, err := storagegen.New(q).ListMediaProfileCustomFormats(ctx, profileID)
	if err != nil {
		return nil, err
	}
	scores := make([]MediaProfileCustomFormatScore, 0, len(rows))
	for _, row := range rows {
		scores = append(scores, MediaProfileCustomFormatScore{
			CustomFormatID: row.CustomFormatID,
			Score:          row.Score,
		})
	}
	return scores, nil
}

func replaceMediaProfileCustomFormats(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	scores []MediaProfileCustomFormatScore,
) error {
	queries := storagegen.New(q)
	if err := queries.ClearMediaProfileCustomFormats(ctx, profileID); err != nil {
		return err
	}
	for _, score := range scores {
		if err := queries.AddMediaProfileCustomFormat(ctx, storagegen.AddMediaProfileCustomFormatParams{
			ProfileID:      profileID,
			CustomFormatID: score.CustomFormatID,
			Score:          score.Score,
		}); err != nil {
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
	queries := storagegen.New(q)
	if err := queries.ClearMediaProfileLanguages(ctx, profileID); err != nil {
		return err
	}
	for _, score := range scores {
		if err := queries.AddMediaProfileLanguage(ctx, storagegen.AddMediaProfileLanguageParams{
			ProfileID:  profileID,
			LanguageID: score.LanguageID,
			Score:      score.Score,
			Required:   score.Required,
		}); err != nil {
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
