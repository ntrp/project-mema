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
	qualityIDs := compactUnique(request.QualityIds)
	if len(qualityIDs) == 0 {
		writeError(w, http.StatusBadRequest, "quality_required", "Select at least one quality")
		return storage.MediaProfileInput{}, false
	}
	audioTargets := mediaProfileAudioTargets(request.AudioTargets)
	if len(audioTargets) == 0 {
		writeError(w, http.StatusBadRequest, "audio_required", "Add at least one audio language")
		return storage.MediaProfileInput{}, false
	}
	return storage.MediaProfileInput{
		Name:                              name,
		IsDefault:                         request.IsDefault,
		FinalContainer:                    string(request.FinalContainer),
		QualityIDs:                        qualityIDs,
		UpgradesAllowed:                   request.UpgradesAllowed,
		UpgradeUntilQualityID:             request.UpgradeUntilQualityId,
		MinimumCustomFormatScore:          request.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     request.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: request.MinimumCustomFormatScoreIncrement,
		RemoveUnwantedAudio:               request.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:         string(request.AudioLossyTranscodePolicy),
		RemoveUnwantedSubtitles:           request.RemoveUnwantedSubtitles,
		SubtitleMode:                      string(request.SubtitleMode),
		AllowSubtitleReleaseFallback:      request.AllowSubtitleReleaseFallback,
		PreferredProtocol:                 string(request.PreferredProtocol),
		SeriesPackPreference:              string(request.SeriesPackPreference),
		VideoTarget:                       mediaProfileVideoTarget(request.VideoTarget),
		AudioTargets:                      audioTargets,
		SubtitleTargets:                   mediaProfileSubtitleTargets(request.SubtitleTargets),
		CustomFormatScores:                mediaProfileCustomFormatScores(request.CustomFormatScores),
	}, true
}

func compactUnique(values []string) []string {
	result := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func mediaProfileVideoTarget(value MediaProfileVideoTarget) storage.MediaProfileVideoTarget {
	return storage.MediaProfileVideoTarget{
		Codecs:              compactUniquePtr(value.Codecs),
		CodecRequired:       boolValueOrFalse(value.CodecRequired),
		CodecScore:          int32ValueOrZero(value.CodecScore),
		HDRFormats:          compactUniquePtr(value.HdrFormats),
		HDRRequired:         boolValueOrFalse(value.HdrRequired),
		HDRScore:            int32ValueOrZero(value.HdrScore),
		PixelFormats:        compactUniquePtr(value.PixelFormats),
		PixelFormatRequired: boolValueOrFalse(value.PixelFormatRequired),
		PixelFormatScore:    int32ValueOrZero(value.PixelFormatScore),
	}
}

func compactUniquePtr(values *[]string) []string {
	if values == nil {
		return nil
	}
	return compactUnique(*values)
}

func mediaProfileAudioTargets(values []MediaProfileAudioTarget) []storage.MediaProfileAudioTarget {
	targets := make([]storage.MediaProfileAudioTarget, 0, len(values))
	for _, value := range values {
		targets = append(targets, storage.MediaProfileAudioTarget{
			LanguageID:           value.LanguageId,
			Score:                value.Score,
			TargetCodec:          value.TargetCodec,
			TargetChannels:       compactUniquePtr(value.TargetChannels),
			MinimumBitrateKbps:   value.MinimumBitrateKbps,
			PreferredBitrateKbps: value.PreferredBitrateKbps,
		})
	}
	return targets
}

func mediaProfileSubtitleTargets(values []MediaProfileSubtitleTarget) []storage.MediaProfileSubtitleTarget {
	targets := make([]storage.MediaProfileSubtitleTarget, 0, len(values))
	for _, value := range values {
		targets = append(targets, storage.MediaProfileSubtitleTarget{
			LanguageID: value.LanguageId,
			Score:      value.Score,
			Formats:    compactUniquePtr(value.Formats),
		})
	}
	return targets
}

func mediaProfileCustomFormatScores(values []MediaProfileCustomFormatScore) []storage.MediaProfileCustomFormatScore {
	scores := make([]storage.MediaProfileCustomFormatScore, 0, len(values))
	for _, value := range values {
		scores = append(scores, storage.MediaProfileCustomFormatScore{
			CustomFormatID: value.CustomFormatId,
			Score:          value.Score,
		})
	}
	return scores
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
		IsDefault:                         profile.IsDefault,
		FinalContainer:                    MediaProfileFinalContainer(profile.FinalContainer),
		QualityIds:                        profile.QualityIDs,
		UpgradesAllowed:                   profile.UpgradesAllowed,
		UpgradeUntilQualityId:             profile.UpgradeUntilQualityID,
		MinimumCustomFormatScore:          profile.MinimumCustomFormatScore,
		UpgradeUntilCustomFormatScore:     profile.UpgradeUntilCustomFormatScore,
		MinimumCustomFormatScoreIncrement: profile.MinimumCustomFormatScoreIncrement,
		RemoveUnwantedAudio:               profile.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:         MediaProfileLossyTranscodePolicy(profile.AudioLossyTranscodePolicy),
		RemoveUnwantedSubtitles:           profile.RemoveUnwantedSubtitles,
		SubtitleMode:                      mediaProfileSubtitleModeResponse(profile.SubtitleMode),
		AllowSubtitleReleaseFallback:      profile.AllowSubtitleReleaseFallback,
		PreferredProtocol:                 MediaProfilePreferredProtocol(profile.PreferredProtocol),
		SeriesPackPreference:              MediaProfileSeriesPackPreference(profile.SeriesPackPreference),
		VideoTarget:                       mediaProfileVideoTargetResponse(profile.VideoTarget),
		AudioTargets:                      mediaProfileAudioTargetResponses(profile.AudioTargets),
		SubtitleTargets:                   mediaProfileSubtitleTargetResponses(profile.SubtitleTargets),
		CustomFormatScores:                mediaProfileCustomFormatScoreResponses(profile.CustomFormatScores),
		CreatedAt:                         profile.CreatedAt,
		UpdatedAt:                         profile.UpdatedAt,
	}
}

func mediaProfileVideoTargetResponse(target storage.MediaProfileVideoTarget) MediaProfileVideoTarget {
	return MediaProfileVideoTarget{
		Codecs:              &target.Codecs,
		CodecRequired:       &target.CodecRequired,
		CodecScore:          &target.CodecScore,
		HdrFormats:          &target.HDRFormats,
		HdrRequired:         &target.HDRRequired,
		HdrScore:            &target.HDRScore,
		PixelFormats:        &target.PixelFormats,
		PixelFormatRequired: &target.PixelFormatRequired,
		PixelFormatScore:    &target.PixelFormatScore,
	}
}

func mediaProfileAudioTargetResponses(targets []storage.MediaProfileAudioTarget) []MediaProfileAudioTarget {
	response := make([]MediaProfileAudioTarget, 0, len(targets))
	for _, target := range targets {
		response = append(response, MediaProfileAudioTarget{
			LanguageId:           target.LanguageID,
			Score:                target.Score,
			TargetCodec:          target.TargetCodec,
			TargetChannels:       &target.TargetChannels,
			MinimumBitrateKbps:   target.MinimumBitrateKbps,
			PreferredBitrateKbps: target.PreferredBitrateKbps,
		})
	}
	return response
}

func mediaProfileSubtitleTargetResponses(targets []storage.MediaProfileSubtitleTarget) []MediaProfileSubtitleTarget {
	response := make([]MediaProfileSubtitleTarget, 0, len(targets))
	for _, target := range targets {
		response = append(response, MediaProfileSubtitleTarget{
			LanguageId: target.LanguageID,
			Score:      target.Score,
			Formats:    &target.Formats,
		})
	}
	return response
}

func mediaProfileSubtitleModeResponse(value string) MediaProfileSubtitleMode {
	switch value {
	case "embedded":
		return MediaProfileSubtitleModeEmbedded
	case "external":
		return MediaProfileSubtitleModeExternal
	default:
		return MediaProfileSubtitleModeMixed
	}
}

func mediaProfileCustomFormatScoreResponses(scores []storage.MediaProfileCustomFormatScore) []MediaProfileCustomFormatScore {
	response := make([]MediaProfileCustomFormatScore, 0, len(scores))
	for _, score := range scores {
		response = append(response, MediaProfileCustomFormatScore{
			CustomFormatId: openapi_types.UUID(score.CustomFormatID),
			Score:          score.Score,
		})
	}
	return response
}

func boolValueOrFalse(value *bool) bool {
	return value != nil && *value
}

func int32ValueOrZero(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}
