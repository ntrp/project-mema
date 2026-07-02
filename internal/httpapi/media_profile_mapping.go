package httpapi

import (
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaProfileInput(w http.ResponseWriter, request MediaProfileRequest) (storage.MediaProfileInput, bool) {
	name := strings.Join(strings.Fields(request.Name), " ")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.MediaProfileInput{}, false
	}

	qualityIDs := make([]string, 0, len(request.QualityIds))
	seen := map[string]struct{}{}
	for _, value := range request.QualityIds {
		qualityID := strings.TrimSpace(value)
		if qualityID == "" {
			continue
		}
		if _, ok := seen[qualityID]; ok {
			continue
		}
		seen[qualityID] = struct{}{}
		qualityIDs = append(qualityIDs, qualityID)
	}
	if len(qualityIDs) == 0 {
		writeError(w, http.StatusBadRequest, "quality_required", "Select at least one quality")
		return storage.MediaProfileInput{}, false
	}

	targetLanguages := make([]string, 0, len(request.TargetLanguages))
	for _, value := range request.TargetLanguages {
		language := strings.TrimSpace(value)
		if language != "" {
			targetLanguages = append(targetLanguages, language)
		}
	}
	targetLanguageScores := make([]storage.MediaProfileLanguageScore, 0, len(request.TargetLanguageScores))
	for _, value := range request.TargetLanguageScores {
		targetLanguageScores = append(targetLanguageScores, storage.MediaProfileLanguageScore{
			LanguageID: value.LanguageId,
			Score:      value.Score,
			Required:   value.Required,
		})
	}
	customFormatScores := make([]storage.MediaProfileCustomFormatScore, 0, len(request.CustomFormatScores))
	for _, value := range request.CustomFormatScores {
		customFormatScores = append(customFormatScores, storage.MediaProfileCustomFormatScore{
			CustomFormatID: value.CustomFormatId,
			Score:          value.Score,
		})
	}

	return storage.MediaProfileInput{
		Name:                              name,
		QualityIDs:                        qualityIDs,
		UpgradesAllowed:                   request.UpgradesAllowed,
		UpgradeUntilQualityID:             request.UpgradeUntilQualityId,
		MinimumCustomFormatScore:          request.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     request.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: request.MinimumCustomFormatScoreIncrement,
		RemoveNonEnabledLanguages:         request.RemoveNonEnabledLanguages,
		TargetLanguages:                   targetLanguages,
		TargetLanguageScores:              targetLanguageScores,
		CustomFormatScores:                customFormatScores,
	}, true
}

func mediaProfileListResponse(profiles []storage.MediaProfile) MediaProfileListResponse {
	response := MediaProfileListResponse{Profiles: make([]MediaProfile, 0, len(profiles))}
	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, mediaProfileResponse(profile))
	}
	return response
}

func mediaProfileResponse(profile storage.MediaProfile) MediaProfile {
	return MediaProfile{
		Id:                                profile.ID,
		Name:                              profile.Name,
		QualityIds:                        profile.QualityIDs,
		UpgradesAllowed:                   profile.UpgradesAllowed,
		UpgradeUntilQualityId:             profile.UpgradeUntilQualityID,
		MinimumCustomFormatScore:          profile.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     profile.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: profile.MinimumCustomFormatScoreIncrement,
		RemoveNonEnabledLanguages:         profile.RemoveNonEnabledLanguages,
		TargetLanguages:                   profile.TargetLanguages,
		TargetLanguageScores:              mediaProfileLanguageScoreResponses(profile.TargetLanguageScores),
		CustomFormatScores:                mediaProfileCustomFormatScoreResponses(profile.CustomFormatScores),
		CreatedAt:                         profile.CreatedAt,
		UpdatedAt:                         profile.UpdatedAt,
	}
}

func mediaProfileLanguageScoreResponses(scores []storage.MediaProfileLanguageScore) []MediaProfileLanguageScore {
	response := make([]MediaProfileLanguageScore, 0, len(scores))
	for _, score := range scores {
		response = append(response, MediaProfileLanguageScore{
			LanguageId: score.LanguageID,
			Score:      score.Score,
			Required:   score.Required,
		})
	}
	return response
}

func mediaProfileCustomFormatScoreResponses(
	scores []storage.MediaProfileCustomFormatScore,
) []MediaProfileCustomFormatScore {
	response := make([]MediaProfileCustomFormatScore, 0, len(scores))
	for _, score := range scores {
		response = append(response, MediaProfileCustomFormatScore{
			CustomFormatId: openapi_types.UUID(score.CustomFormatID),
			Score:          score.Score,
		})
	}
	return response
}
