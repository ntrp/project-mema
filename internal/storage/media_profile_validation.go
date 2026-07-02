package storage

import (
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

func normalizeProfileQualityIDs(values []string) ([]string, error) {
	definitions := QualitySizeDefinitionMap()
	seen := map[string]struct{}{}
	qualityIDs := []string{}
	for _, value := range values {
		qualityID := strings.TrimSpace(value)
		if qualityID == "" {
			continue
		}
		if _, ok := definitions[qualityID]; !ok {
			return nil, ErrInvalidInput
		}
		if _, ok := seen[qualityID]; ok {
			continue
		}
		seen[qualityID] = struct{}{}
		qualityIDs = append(qualityIDs, qualityID)
	}
	return qualityIDs, nil
}

func normalizeMediaProfileInput(
	input MediaProfileInput,
	qualityIDs []string,
) (MediaProfileInput, error) {
	normalized := input
	qualitySet := map[string]struct{}{}
	for _, qualityID := range qualityIDs {
		qualitySet[qualityID] = struct{}{}
	}
	if input.UpgradeUntilQualityID != nil {
		qualityID := strings.TrimSpace(*input.UpgradeUntilQualityID)
		if qualityID == "" {
			normalized.UpgradeUntilQualityID = nil
		} else {
			if _, ok := qualitySet[qualityID]; !ok {
				return MediaProfileInput{}, ErrInvalidInput
			}
			normalized.UpgradeUntilQualityID = &qualityID
		}
	}
	if input.MinimumCustomFormatScoreIncrement < 0 {
		return MediaProfileInput{}, ErrInvalidInput
	}
	normalized.QualityIDs = qualityIDs
	normalized.TargetLanguageScores = normalizeTargetLanguageScores(input)
	normalized.TargetLanguages = languageIDsFromScores(normalized.TargetLanguageScores)
	normalized.CustomFormatScores = normalizeCustomFormatScores(input.CustomFormatScores)
	return normalized, nil
}

func normalizeTargetLanguageScores(input MediaProfileInput) []MediaProfileLanguageScore {
	if len(input.TargetLanguageScores) > 0 {
		return normalizeLanguageScoreValues(input.TargetLanguageScores)
	}
	scores := make([]MediaProfileLanguageScore, 0, len(input.TargetLanguages))
	for _, languageID := range input.TargetLanguages {
		scores = append(scores, MediaProfileLanguageScore{LanguageID: languageID})
	}
	return normalizeLanguageScoreValues(scores)
}

func normalizeLanguageScoreValues(values []MediaProfileLanguageScore) []MediaProfileLanguageScore {
	seen := map[string]struct{}{}
	scores := []MediaProfileLanguageScore{}
	for _, value := range values {
		language := strings.ToLower(strings.Join(strings.Fields(value.LanguageID), "-"))
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		seen[language] = struct{}{}
		scores = append(scores, MediaProfileLanguageScore{
			LanguageID: language,
			Score:      value.Score,
			Required:   value.Required,
		})
	}
	return scores
}

func languageIDsFromScores(scores []MediaProfileLanguageScore) []string {
	languages := make([]string, 0, len(scores))
	for _, score := range scores {
		languages = append(languages, score.LanguageID)
	}
	return languages
}

func normalizeMediaProfileName(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

var nonProfileSlugCharacter = regexp.MustCompile(`[^a-z0-9]+`)

func mediaProfileSlug(name string) string {
	slug := strings.Trim(nonProfileSlugCharacter.ReplaceAllString(strings.ToLower(name), "-"), "-")
	if len(slug) > 80 {
		slug = strings.Trim(slug[:80], "-")
	}
	return slug
}

func scanMediaProfileBase(row pgx.Row) (MediaProfile, error) {
	var profile MediaProfile
	err := row.Scan(
		&profile.ID,
		&profile.Name,
		&profile.UpgradesAllowed,
		&profile.UpgradeUntilQualityID,
		&profile.MinimumCustomFormatScore,
		&profile.UpgradeUntilCustomFormatScore,
		&profile.MinimumCustomFormatScoreIncrement,
		&profile.RemoveNonEnabledLanguages,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}
