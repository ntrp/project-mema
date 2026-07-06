package storage

import (
	"regexp"
	"strings"
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
	if len(qualityIDs) == 0 {
		return MediaProfileInput{}, ErrInvalidInput
	}
	normalized.FinalContainer = normalizeContainer(input.FinalContainer)
	normalized.PreferredProtocol = normalizePreferredProtocol(input.PreferredProtocol)
	normalized.SeriesPackPreference = normalizeSeriesPackPreference(input.SeriesPackPreference)
	normalized.QualityIDs = qualityIDs
	normalized.VideoTarget = normalizeVideoTarget(input.VideoTarget)
	normalized.AudioTargets = normalizeAudioTargets(input.AudioTargets)
	if len(normalized.AudioTargets) == 0 {
		return MediaProfileInput{}, ErrInvalidInput
	}
	normalized.SubtitleTargets = normalizeSubtitleTargets(input.SubtitleTargets)
	normalized.CustomFormatScores = normalizeCustomFormatScores(input.CustomFormatScores)
	return normalized, nil
}

func normalizeContainer(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "mp4":
		return "mp4"
	default:
		return "mkv"
	}
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

func normalizeVideoTarget(value MediaProfileVideoTarget) MediaProfileVideoTarget {
	return MediaProfileVideoTarget{
		Codecs:              normalizedTextList(value.Codecs),
		CodecRequired:       value.CodecRequired,
		CodecScore:          value.CodecScore,
		HDRFormats:          normalizedTextList(value.HDRFormats),
		HDRRequired:         value.HDRRequired,
		HDRScore:            value.HDRScore,
		PixelFormats:        normalizedTextList(value.PixelFormats),
		PixelFormatRequired: value.PixelFormatRequired,
		PixelFormatScore:    value.PixelFormatScore,
	}
}

func normalizeAudioTargets(values []MediaProfileAudioTarget) []MediaProfileAudioTarget {
	seen := map[string]struct{}{}
	targets := []MediaProfileAudioTarget{}
	for _, value := range values {
		language := normalizeLanguageID(value.LanguageID)
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		seen[language] = struct{}{}
		targets = append(targets, MediaProfileAudioTarget{
			LanguageID:           language,
			Score:                value.Score,
			Required:             value.Required,
			Codecs:               normalizedTextList(value.Codecs),
			Channels:             normalizedTextList(value.Channels),
			MinimumBitrateKbps:   positiveInt32Ptr(value.MinimumBitrateKbps),
			PreferredBitrateKbps: positiveInt32Ptr(value.PreferredBitrateKbps),
			LossyTranscodePolicy: normalizeLossyPolicy(value.LossyTranscodePolicy),
		})
	}
	return targets
}

func normalizeSubtitleTargets(values []MediaProfileSubtitleTarget) []MediaProfileSubtitleTarget {
	seen := map[string]struct{}{}
	targets := []MediaProfileSubtitleTarget{}
	for _, value := range values {
		language := normalizeLanguageID(value.LanguageID)
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		seen[language] = struct{}{}
		targets = append(targets, MediaProfileSubtitleTarget{
			LanguageID: language,
			Score:      value.Score,
			Required:   value.Required,
			Source:     normalizeSubtitleSource(value.Source),
			Formats:    normalizedTextList(value.Formats),
		})
	}
	return targets
}

func normalizeLanguageID(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), "-"))
}

func normalizeSubtitleSource(value string) string {
	switch strings.TrimSpace(value) {
	case "embedded", "external":
		return strings.TrimSpace(value)
	default:
		return "any"
	}
}

func normalizeLossyPolicy(value string) string {
	switch strings.TrimSpace(value) {
	case "losslessToLossy", "lossyToLossy":
		return strings.TrimSpace(value)
	default:
		return "disabled"
	}
}

func positiveInt32Ptr(value *int32) *int32 {
	if value == nil || *value <= 0 {
		return nil
	}
	return value
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

func normalizedTextList(values []string) []string {
	seen := map[string]struct{}{}
	result := []string{}
	for _, value := range values {
		text := strings.ToLower(strings.Join(strings.Fields(value), "-"))
		if text == "" {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		result = append(result, text)
	}
	return result
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
