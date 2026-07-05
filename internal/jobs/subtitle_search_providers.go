package jobs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

func bestSubtitleCandidate(
	ctx context.Context,
	settings *storage.SettingsStore,
	service *subtitles.Service,
	request subtitles.SearchRequest,
) (subtitles.Candidate, storage.SubtitleProvider, error) {
	providers, err := settings.ListSubtitleProviders(ctx)
	if err != nil {
		return subtitles.Candidate{}, storage.SubtitleProvider{}, err
	}
	var failures []string
	for _, provider := range providers {
		if !provider.Enabled {
			continue
		}
		candidates, err := service.Search(ctx, subtitleConfig(provider), request)
		if err != nil {
			failures = append(failures, provider.Name+": "+err.Error())
			continue
		}
		if candidate, ok := chooseSubtitleCandidate(candidates, request.LanguageID); ok {
			return candidate, provider, nil
		}
	}
	if len(failures) > 0 {
		return subtitles.Candidate{}, storage.SubtitleProvider{}, errors.New(strings.Join(failures, "; "))
	}
	return subtitles.Candidate{}, storage.SubtitleProvider{}, errors.New("no subtitle candidate found")
}

func chooseSubtitleCandidate(candidates []subtitles.Candidate, languageID string) (subtitles.Candidate, bool) {
	var best subtitles.Candidate
	found := false
	for _, candidate := range candidates {
		if languageMatchKey(candidate.LanguageID) != languageMatchKey(languageID) {
			continue
		}
		if !found || candidate.DownloadCount > best.DownloadCount {
			best = candidate
			found = true
		}
	}
	return best, found
}

func subtitleConfig(provider storage.SubtitleProvider) subtitles.Config {
	return subtitles.Config{
		Name:     provider.Name,
		Type:     provider.Type,
		BaseURL:  provider.BaseURL,
		Username: provider.Username,
		Password: provider.Password,
		APIKey:   provider.APIKey,
	}
}

func subtitleRecord(
	item storage.MediaItem,
	provider storage.SubtitleProvider,
	candidate subtitles.Candidate,
	request subtitles.SearchRequest,
	path string,
	sourceURL string,
) storage.MediaItemSubtitleInput {
	return storage.MediaItemSubtitleInput{
		MediaItemID:  item.ID,
		ProviderID:   &provider.ID,
		ProviderName: provider.Name,
		LanguageID:   request.LanguageID,
		FilePath:     path,
		SourceURL:    optionalSubtitleString(firstNonEmpty(sourceURL, candidate.SourceURL)),
		ReleaseName:  optionalSubtitleString(candidate.ReleaseName),
	}
}

func optionalSubtitleString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func languageMatchKey(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "en", "eng", "english":
		return "english"
	case "de", "deu", "ger", "german":
		return "german"
	case "fr", "fra", "fre", "french":
		return "french"
	case "es", "spa", "spanish":
		return "spanish"
	case "ja", "jpn", "japanese":
		return "japanese"
	default:
		return strings.ReplaceAll(normalized, " ", "-")
	}
}

func parseSmallInt(value string) int32 {
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0
	}
	return int32(parsed)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
