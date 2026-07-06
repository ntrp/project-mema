package storage

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
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
	normalized.PreferredProtocol = normalizePreferredProtocol(input.PreferredProtocol)
	normalized.SeriesPackPreference = normalizeSeriesPackPreference(input.SeriesPackPreference)
	normalized.QualityIDs = qualityIDs
	normalized.TargetLanguageScores = normalizeTargetLanguageScores(input)
	normalized.TargetLanguages = languageIDsFromScores(normalized.TargetLanguageScores)
	normalized.SubtitleLanguages = normalizeSubtitleLanguages(input.SubtitleLanguages)
	componentTargets, err := normalizeComponentTargets(input.ComponentTargets)
	if err != nil {
		return MediaProfileInput{}, err
	}
	normalized.ComponentTargets = componentTargets
	normalized.CustomFormatScores = normalizeCustomFormatScores(input.CustomFormatScores)
	return normalized, nil
}

func normalizePreferredProtocol(value string) string {
	switch strings.TrimSpace(value) {
	case "torrent", "usenet":
		return strings.TrimSpace(value)
	default:
		return "any"
	}
}

func normalizeSeriesPackPreference(value string) string {
	switch strings.TrimSpace(value) {
	case "preferPacks", "preferEpisodes":
		return strings.TrimSpace(value)
	default:
		return "auto"
	}
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

func normalizeSubtitleLanguages(values []MediaProfileSubtitleLanguage) []MediaProfileSubtitleLanguage {
	seen := map[string]struct{}{}
	languages := []MediaProfileSubtitleLanguage{}
	for _, value := range values {
		language := strings.ToLower(strings.Join(strings.Fields(value.LanguageID), "-"))
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		seen[language] = struct{}{}
		languages = append(languages, MediaProfileSubtitleLanguage{
			LanguageID:   language,
			Score:        value.Score,
			Required:     value.Required,
			SubtitleType: normalizeSubtitleType(value.SubtitleType),
		})
	}
	return languages
}

func normalizeSubtitleType(value string) string {
	switch strings.TrimSpace(value) {
	case "embedded", "external":
		return strings.TrimSpace(value)
	default:
		return "any"
	}
}

func normalizeComponentTargets(values []MediaProfileComponentTarget) ([]MediaProfileComponentTarget, error) {
	targets := []MediaProfileComponentTarget{}
	for _, value := range values {
		target, ok, err := normalizeComponentTarget(value)
		if err != nil {
			return nil, err
		}
		if ok {
			targets = append(targets, target)
		}
	}
	return targets, nil
}

func normalizeComponentTarget(value MediaProfileComponentTarget) (MediaProfileComponentTarget, bool, error) {
	target := MediaProfileComponentTarget{
		ID:               value.ID,
		ComponentType:    strings.TrimSpace(value.ComponentType),
		Required:         value.Required,
		LanguageID:       normalizedTextPtr(value.LanguageID, "-"),
		Codec:            normalizedTextPtr(value.Codec),
		Channels:         normalizedTextPtr(value.Channels),
		Source:           normalizeComponentSource(value.Source),
		FallbackBehavior: normalizeComponentFallback(value.FallbackBehavior),
	}
	if target.ID == uuid.Nil {
		target.ID = uuid.New()
	}
	switch target.ComponentType {
	case "video":
		target.LanguageID = nil
		target.Channels = nil
	case "audio":
		if target.LanguageID == nil && target.Codec == nil && target.Channels == nil {
			return MediaProfileComponentTarget{}, false, nil
		}
	case "subtitle":
		target.Channels = nil
		if target.Source == "release" {
			target.Source = "subtitleProvider"
		}
		if target.LanguageID == nil && target.Codec == nil {
			return MediaProfileComponentTarget{}, false, nil
		}
	default:
		return MediaProfileComponentTarget{}, false, ErrInvalidInput
	}
	if target.Source == "" || target.FallbackBehavior == "" {
		return MediaProfileComponentTarget{}, false, ErrInvalidInput
	}
	return target, true, nil
}

func normalizeComponentSource(value string) string {
	switch strings.TrimSpace(value) {
	case "release", "subtitleProvider", "existing":
		return strings.TrimSpace(value)
	default:
		return "release"
	}
}

func normalizeComponentFallback(value string) string {
	switch strings.TrimSpace(value) {
	case "strict", "preferExisting", "allowMissing":
		return strings.TrimSpace(value)
	default:
		return "strict"
	}
}

func normalizedTextPtr(value *string, emptyValues ...string) *string {
	if value == nil {
		return nil
	}
	text := strings.ToLower(strings.Join(strings.Fields(*value), "-"))
	for _, empty := range emptyValues {
		if text == empty {
			return nil
		}
	}
	if text == "" {
		return nil
	}
	return &text
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
		&profile.RemoveNonEnabledSubtitleLanguages,
		&profile.PreferredProtocol,
		&profile.SeriesPackPreference,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}
