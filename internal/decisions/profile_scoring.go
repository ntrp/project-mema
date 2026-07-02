package decisions

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func profileQualityScore(qualityID string, profile *storage.MediaProfile) int32 {
	if profile == nil {
		return releaseQualityScore(qualityID)
	}
	for index, value := range profile.QualityIDs {
		if value == qualityID {
			return int32(index+1) * 1000
		}
	}
	return 0
}

func customFormatScore(
	parsed ParsedRelease,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
) (int32, []ReleaseScoreContributor) {
	if profile == nil || len(profile.CustomFormatScores) == 0 || len(formats) == 0 {
		return 0, nil
	}
	scores := map[uuid.UUID]int32{}
	for _, score := range profile.CustomFormatScores {
		scores[score.CustomFormatID] = score.Score
	}
	var total int32
	contributors := []ReleaseScoreContributor{}
	for _, match := range MatchCustomFormats(parsed, formats) {
		formatID, err := uuid.Parse(match.ID)
		if err != nil {
			continue
		}
		score, ok := scores[formatID]
		if !ok {
			continue
		}
		total += score
		contributors = append(contributors, ReleaseScoreContributor{Label: match.Name, Score: score})
	}
	return total, contributors
}

func languageScore(
	parsed ParsedRelease,
	profile *storage.MediaProfile,
) (int32, []ReleaseScoreContributor, string) {
	if profile == nil || len(profile.TargetLanguageScores) == 0 {
		return 0, nil, ""
	}
	releaseLanguages := normalizedLanguages(parsed.Languages)
	var total int32
	contributors := []ReleaseScoreContributor{}
	for _, target := range profile.TargetLanguageScores {
		if _, ok := releaseLanguages[target.LanguageID]; !ok {
			if target.Required {
				return 0, nil, fmt.Sprintf("Required language %s is missing.", target.LanguageID)
			}
			continue
		}
		total += target.Score
		contributors = append(contributors, ReleaseScoreContributor{
			Label: fmt.Sprintf("Language: %s", target.LanguageID),
			Score: target.Score,
		})
	}
	if profile.RemoveNonEnabledLanguages {
		targets := map[string]struct{}{}
		for _, target := range profile.TargetLanguageScores {
			targets[target.LanguageID] = struct{}{}
		}
		for language := range releaseLanguages {
			if language == "multiple" {
				continue
			}
			if _, ok := targets[language]; !ok {
				return 0, nil, fmt.Sprintf("Language %s is not enabled in the profile.", language)
			}
		}
	}
	return total, contributors, ""
}

func normalizedLanguages(values []string) map[string]struct{} {
	languages := map[string]struct{}{}
	for _, value := range values {
		language := strings.ToLower(strings.Join(strings.Fields(value), "-"))
		if language != "" {
			languages[language] = struct{}{}
		}
	}
	return languages
}

func rankContributors(
	parsed ParsedRelease,
	qualityScore int32,
	customScore int32,
	languageScore int32,
	meta releaseMeta,
) []ReleaseScoreContributor {
	contributors := []ReleaseScoreContributor{
		{Label: fmt.Sprintf("Quality rank: %s", parsed.Quality), Score: qualityScore},
		{Label: "Custom formats", Score: customScore},
		{Label: "Languages", Score: languageScore},
		{Label: fmt.Sprintf("Protocol: %s", meta.IndexerType), Score: 0},
		{Label: "Size bytes", Score: int32(clampInt64(meta.SizeBytes))},
	}
	if meta.Seeders != nil {
		contributors = append(contributors, ReleaseScoreContributor{Label: "Seeders", Score: *meta.Seeders})
	}
	if meta.Peers != nil {
		contributors = append(contributors, ReleaseScoreContributor{Label: "Peers", Score: *meta.Peers})
	}
	return contributors
}

func clampInt64(value int64) int {
	if value > 2147483647 {
		return 2147483647
	}
	if value < -2147483648 {
		return -2147483648
	}
	return int(value)
}
