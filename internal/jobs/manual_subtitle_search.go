package jobs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

type ManualSubtitleCandidate struct {
	ID        string
	Provider  storage.SubtitleProvider
	Candidate subtitles.Candidate
	Match     ManualSubtitleMatch
	Protocol  string
}

type ManualSubtitleMatch struct {
	Severity string
	Label    string
	Details  []string
}

type ManualSubtitleGrabInput struct {
	ProviderID uuid.UUID
	FilePath   string
	LanguageID string
	Title      string
	Format     string
	FileID     *int64
	SourceURL  *string
	SourceRef  *string
}

func SearchManualSubtitles(
	ctx context.Context,
	settings *storage.SettingsStore,
	service *subtitles.Service,
	item storage.MediaItem,
	query string,
	languageID string,
	filePath string,
) ([]ManualSubtitleCandidate, []string, error) {
	request, ok, err := manualSubtitleSearchRequest(ctx, settings, item, query, languageID, filePath)
	if err != nil || !ok {
		return nil, nil, err
	}
	providers, err := settings.ListSubtitleProviders(ctx)
	if err != nil {
		return nil, nil, err
	}
	results := []ManualSubtitleCandidate{}
	logs := []string{fmt.Sprintf("Searching subtitles for %s", request.Title)}
	for _, provider := range providers {
		if !provider.Enabled {
			continue
		}
		logs = append(logs, fmt.Sprintf("Searching %s", provider.Name))
		candidates, err := service.Search(ctx, subtitleConfig(provider), request)
		if err != nil {
			logs = append(logs, fmt.Sprintf("%s failed: %s", provider.Name, err.Error()))
			continue
		}
		logs = append(logs, fmt.Sprintf("%s returned %d subtitle(s)", provider.Name, len(candidates)))
		for _, candidate := range candidates {
			results = append(results, manualSubtitleCandidate(provider, candidate, request))
		}
	}
	return results, logs, nil
}

func GrabManualSubtitle(
	ctx context.Context,
	settings *storage.SettingsStore,
	service *subtitles.Service,
	item storage.MediaItem,
	input ManualSubtitleGrabInput,
) error {
	path, err := settings.MediaItemFilePath(ctx, item.ID, input.FilePath)
	if err != nil {
		return err
	}
	provider, err := settings.GetSubtitleProvider(ctx, input.ProviderID)
	if err != nil {
		return err
	}
	if !provider.Enabled {
		return fmt.Errorf("subtitle provider %s is disabled", provider.Name)
	}
	request := subtitles.SearchRequest{
		MediaType:  item.Type,
		Title:      item.Title,
		LanguageID: strings.TrimSpace(input.LanguageID),
		Year:       item.Year,
		FilePath:   path,
	}
	candidate := subtitles.Candidate{
		ProviderName: provider.Name,
		LanguageID:   strings.TrimSpace(input.LanguageID),
		Format:       strings.TrimSpace(input.Format),
		ReleaseName:  strings.TrimSpace(input.Title),
		SourceURL:    stringPtrValue(input.SourceURL),
		SourceRef:    stringPtrValue(input.SourceRef),
	}
	if input.FileID != nil {
		candidate.FileID = *input.FileID
	}
	download, err := service.Download(ctx, subtitleConfig(provider), candidate)
	if err != nil {
		return err
	}
	targetFormat := firstNonEmpty(subtitleTargetFormat(item, request.LanguageID), candidate.Format, "srt")
	artifact, err := writeSubtitleFile(request, download.Content, targetFormat)
	if err != nil {
		return err
	}
	_, err = settings.UpsertMediaItemSubtitle(ctx, subtitleRecord(item, provider, candidate, request, artifact, download.URL))
	return err
}

func manualSubtitleSearchRequest(
	ctx context.Context,
	settings *storage.SettingsStore,
	item storage.MediaItem,
	query string,
	languageID string,
	filePath string,
) (subtitles.SearchRequest, bool, error) {
	path, err := settings.MediaItemFilePath(ctx, item.ID, filePath)
	if err != nil {
		return subtitles.SearchRequest{}, false, err
	}
	request := subtitles.SearchRequest{
		MediaType:  item.Type,
		Title:      firstNonEmpty(query, item.Title),
		LanguageID: strings.TrimSpace(languageID),
		Year:       item.Year,
		FilePath:   path,
	}
	if request.LanguageID == "" || path == "" {
		return subtitles.SearchRequest{}, false, nil
	}
	season, episode, ok := subtitleEpisodeNumbers(path)
	if ok {
		request.SeasonNumber = &season
		request.EpisodeNumber = &episode
	}
	return request, true, nil
}

func manualSubtitleCandidate(
	provider storage.SubtitleProvider,
	candidate subtitles.Candidate,
	request subtitles.SearchRequest,
) ManualSubtitleCandidate {
	match := manualSubtitleMatch(candidate, request)
	return ManualSubtitleCandidate{
		ID:        subtitleCandidateID(provider.ID, candidate),
		Provider:  provider,
		Candidate: candidate,
		Match:     match,
		Protocol:  "HTTP",
	}
}

func manualSubtitleMatch(candidate subtitles.Candidate, request subtitles.SearchRequest) ManualSubtitleMatch {
	if languageMatchKey(candidate.LanguageID) != languageMatchKey(request.LanguageID) {
		return ManualSubtitleMatch{
			Severity: "warning",
			Label:    "Language differs",
			Details:  []string{fmt.Sprintf("Wanted %s, provider returned %s", request.LanguageID, candidate.LanguageID)},
		}
	}
	return ManualSubtitleMatch{
		Severity: "success",
		Label:    "Match",
		Details:  []string{fmt.Sprintf("Language matches %s", request.LanguageID)},
	}
}

func subtitleCandidateID(providerID uuid.UUID, candidate subtitles.Candidate) string {
	parts := []string{providerID.String(), candidate.LanguageID, candidate.Format, candidate.ReleaseName}
	if candidate.FileID != 0 {
		parts = append(parts, strconv.FormatInt(candidate.FileID, 10))
	}
	if candidate.SourceRef != "" {
		parts = append(parts, candidate.SourceRef)
	}
	return strings.Join(parts, ":")
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
