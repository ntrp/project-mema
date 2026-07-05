package subtitles

import (
	"net/http"
	"net/url"
	"strings"
)

type openSubtitlesSearchResponse struct {
	Data []openSubtitlesSearchItem `json:"data"`
}

type openSubtitlesSearchItem struct {
	Attributes openSubtitlesAttributes `json:"attributes"`
}

type openSubtitlesAttributes struct {
	Language      string               `json:"language"`
	DownloadCount int                  `json:"download_count"`
	URL           string               `json:"url"`
	Files         []openSubtitlesFile  `json:"files"`
	Feature       openSubtitlesFeature `json:"feature_details"`
}

type openSubtitlesFile struct {
	FileID   int64  `json:"file_id"`
	FileName string `json:"file_name"`
}

type openSubtitlesFeature struct {
	Title string `json:"title"`
	Year  *int32 `json:"year"`
}

type openSubtitlesDownloadResponse struct {
	Link string `json:"link"`
}

type openSubtitlesLoginResponse struct {
	Token string `json:"token"`
}

func openSubtitlesEndpoint(baseURL string, path string) (*url.URL, error) {
	base, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || base.Scheme == "" || base.Host == "" {
		return nil, errInvalidBaseURL()
	}
	return base.JoinPath("api", "v1", path), nil
}

func errInvalidBaseURL() error {
	return errText("subtitle provider base URL is invalid")
}

type errText string

func (e errText) Error() string {
	return string(e)
}

func openSubtitlesHeaders(req *http.Request, config Config) {
	if config.APIKey != nil {
		req.Header.Set("Api-Key", strings.TrimSpace(*config.APIKey))
	}
	req.Header.Set("User-Agent", "project-mema")
}

func openSubtitlesAuth(req *http.Request, token string) {
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	}
}

func openSubtitlesLanguage(language string) string {
	switch languageMatchKey(language) {
	case "english":
		return "en"
	case "german":
		return "de"
	case "french":
		return "fr"
	case "spanish":
		return "es"
	case "japanese":
		return "ja"
	default:
		return strings.TrimSpace(language)
	}
}

func openSubtitlesCandidates(
	providerName string,
	requestLanguage string,
	payload openSubtitlesSearchResponse,
) []Candidate {
	candidates := []Candidate{}
	for _, item := range payload.Data {
		language := languageMatchKey(item.Attributes.Language)
		if language == "" {
			language = languageMatchKey(requestLanguage)
		}
		for _, file := range item.Attributes.Files {
			if file.FileID == 0 {
				continue
			}
			candidates = append(candidates, Candidate{
				ProviderName:  providerName,
				LanguageID:    language,
				FileID:        file.FileID,
				ReleaseName:   firstNonEmpty(file.FileName, item.Attributes.Feature.Title),
				DownloadCount: item.Attributes.DownloadCount,
				SourceURL:     item.Attributes.URL,
			})
		}
	}
	return candidates
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func languageMatchKey(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.TrimSuffix(normalized, " language")
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
