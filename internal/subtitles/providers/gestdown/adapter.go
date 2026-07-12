package gestdown

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
)

const (
	key        = "gestdown"
	defaultURL = "https://api.gestdown.info"
	maxBody    = 10 << 20
)

type adapter struct{}

type showsResponse struct {
	Shows []struct {
		ID int64 `json:"id"`
	} `json:"shows"`
}

type subtitlesResponse struct {
	MatchingSubtitles []gestSubtitle `json:"matchingSubtitles"`
}

type gestSubtitle struct {
	SubtitleID      int64    `json:"subtitleId"`
	DownloadURI     string   `json:"downloadUri"`
	Version         string   `json:"version"`
	Qualities       []string `json:"qualities"`
	Completed       bool     `json:"completed"`
	HearingImpaired bool     `json:"hearingImpaired"`
}

func init() { providers.Register(key, adapter{}) }

func (adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	_, _, err := do(ctx, service, config, http.MethodGet, "/shows/external/tvdb/1", false)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil
	}
	return err
}

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if request.MediaType != "" && request.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	if request.SeasonNumber == nil || request.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode are required", providercore.ErrProviderPrerequisiteMissing)
	}
	seriesID := request.MediaContext.ExternalIDs["tvdb"]
	if seriesID == "" {
		seriesID = request.MediaContext.SeasonExternalIDs["tvdb"]
	}
	if seriesID == "" {
		return nil, fmt.Errorf("%w: tvdb id is required", providercore.ErrProviderPrerequisiteMissing)
	}
	shows, err := searchShow(ctx, service, config, seriesID)
	if err != nil || len(shows) == 0 {
		return nil, err
	}
	lang := language(request.LanguageID)
	out := []providercore.Candidate{}
	for _, showID := range shows {
		subs, err := searchSubtitles(ctx, service, config, showID, *request.SeasonNumber, *request.EpisodeNumber, lang)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return []providercore.Candidate{}, nil
			}
			return nil, err
		}
		for _, sub := range subs {
			if !sub.Completed {
				continue
			}
			out = append(out, providercore.Candidate{ProviderName: key, LanguageID: request.LanguageID, FileID: sub.SubtitleID, Format: "srt", ReleaseName: strings.Join(releases(sub.Version), "\n"), SourceURL: absolute(config, sub.DownloadURI), SourceRef: strings.Join(sub.Qualities, ",")})
		}
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(candidate.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, _, err := do(ctx, service, config, http.MethodGet, candidate.SourceURL, true)
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")), URL: candidate.SourceURL}, nil
}

func searchShow(ctx context.Context, service providercore.Service, config providercore.Config, seriesID string) ([]int64, error) {
	data, _, err := do(ctx, service, config, http.MethodGet, "/shows/external/tvdb/"+url.PathEscape(seriesID), false)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}
	var result showsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(result.Shows))
	for _, show := range result.Shows {
		ids = append(ids, show.ID)
	}
	return ids, nil
}

func searchSubtitles(ctx context.Context, service providercore.Service, config providercore.Config, showID int64, season int32, episode int32, lang string) ([]gestSubtitle, error) {
	endpoint := fmt.Sprintf("/subtitles/get/%d/%d/%d/%s", showID, season, episode, url.PathEscape(lang))
	data, _, err := do(ctx, service, config, http.MethodGet, endpoint, false)
	if err != nil {
		return nil, err
	}
	var result subtitlesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.MatchingSubtitles, nil
}

func language(alpha3 string) string {
	code := strings.ToLower(alpha3)
	if len(code) == 2 || strings.Contains(code, "-") {
		return code
	}
	if converted, ok := addic7edCodes[code]; ok {
		return converted
	}
	return code
}

var addic7edCodes = map[string]string{
	"ara": "ar", "aze": "az", "ben": "bn", "bos": "bs", "bul": "bg", "cat": "ca", "ces": "cs", "cze": "cs",
	"dan": "da", "deu": "de", "ger": "de", "ell": "el", "gre": "el", "eng": "en", "eus": "eu", "baq": "eu",
	"fas": "fa", "per": "fa", "fin": "fi", "fra": "fr", "fre": "fr", "glg": "gl", "heb": "he", "hrv": "hr",
	"hun": "hu", "hye": "hy", "arm": "hy", "ind": "id", "ita": "it", "jpn": "ja", "kor": "ko", "mkd": "mk",
	"msa": "ms", "may": "ms", "nld": "nl", "dut": "nl", "nor": "no", "pol": "pl", "por": "pt", "pob": "pt-BR",
	"ron": "ro", "rum": "ro", "rus": "ru", "slk": "sk", "slo": "sk", "slv": "sl", "spa": "es", "sqi": "sq",
	"alb": "sq", "srp": "sr", "swe": "sv", "tha": "th", "tur": "tr", "ukr": "uk", "vie": "vi", "zho": "zh",
}

func releases(version string) []string {
	items := strings.Split(version, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func absolute(config providercore.Config, endpoint string) string {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}
	u, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultURL))
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func do(ctx context.Context, service providercore.Service, config providercore.Config, method, endpoint string, download bool) ([]byte, *http.Response, error) {
	raw := absolute(config, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, raw, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "Bazarr")
	resp, err := service.DoProviderRequest(req, key, download)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, resp, err
	}
	if len(data) > maxBody {
		return nil, resp, fmt.Errorf("response size limit exceeded")
	}
	if resp.StatusCode == http.StatusLocked {
		return nil, resp, fmt.Errorf("%w: http status 423", providercore.ErrProviderBrokenUpstream)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return data, resp, nil
}
