package betaseries

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/security"
)

const (
	key        = "betaseries"
	defaultURL = "https://api.betaseries.com"
	maxBody    = 10 << 20
)

type adapter struct{}

type responseEnvelope struct {
	Errors   []betaError   `json:"errors"`
	Episode  betaEpisode   `json:"episode"`
	Episodes []betaEpisode `json:"episodes"`
}

type betaError struct {
	Code int `json:"code"`
}

type betaEpisode struct {
	Subtitles []betaSubtitle `json:"subtitles"`
}

type betaSubtitle struct {
	ID       int64  `json:"id"`
	Language string `json:"language"`
	File     string `json:"file"`
	URL      string `json:"url"`
	Source   any    `json:"source"`
}

func init() { providers.Register(key, adapter{}) }

func (adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	_, _, err := do(ctx, service, config, http.MethodGet, "/members/infos", keyQuery(config), false)
	return err
}

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if request.MediaType != "" && request.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	endpoint, q, err := searchRequest(config, request)
	if err != nil {
		return nil, err
	}
	data, _, err := do(ctx, service, config, http.MethodGet, endpoint, q, false)
	if err != nil {
		return nil, err
	}
	var result responseEnvelope
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return []providercore.Candidate{}, nil
	}
	subs := result.Episode.Subtitles
	if len(subs) == 0 && len(result.Episodes) > 0 {
		subs = result.Episodes[0].Subtitles
	}
	out := make([]providercore.Candidate, 0, len(subs))
	for _, sub := range subs {
		lang := language(sub.Language)
		if lang == "" || strings.EqualFold(fmt.Sprint(sub.Source), "seriessub") {
			continue
		}
		out = append(out, providercore.Candidate{ProviderName: key, LanguageID: lang, FileID: sub.ID, Format: format(sub.URL), ReleaseName: sub.File, SourceURL: sub.URL, SourceRef: fmt.Sprint(sub.Source)})
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(candidate.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, _, err := do(ctx, service, config, http.MethodGet, candidate.SourceURL, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	content := data
	if zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data))); err == nil {
		content, err = firstSubtitle(zr)
		if err != nil {
			return providercore.Download{}, err
		}
	} else if member, err := providercore.ExtractSubtitle(path.Base(candidate.SourceURL), data, security.ArchiveLimits{}); err == nil {
		content = member.Content
	}
	return providercore.Download{Content: bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n")), URL: candidate.SourceURL}, nil
}

func searchRequest(config providercore.Config, r providercore.SearchRequest) (string, url.Values, error) {
	q := keyQuery(config)
	q.Set("subtitles", "1")
	q.Set("v", "3")
	if id := r.MediaContext.EpisodeExternalIDs["tvdb"]; id != "" {
		q.Set("thetvdb_id", id)
		return "/episodes/display", q, nil
	}
	id := r.MediaContext.SeasonExternalIDs["tvdb"]
	if id == "" {
		id = r.MediaContext.ExternalIDs["tvdb"]
	}
	if id == "" || r.SeasonNumber == nil || r.EpisodeNumber == nil {
		return "", nil, fmt.Errorf("%w: tvdb season/episode ids are required", providercore.ErrProviderPrerequisiteMissing)
	}
	q.Set("thetvdb_id", id)
	q.Set("season", strconv.Itoa(int(*r.SeasonNumber)))
	q.Set("episode", strconv.Itoa(int(*r.EpisodeNumber)))
	return "/shows/episodes", q, nil
}

func keyQuery(config providercore.Config) url.Values {
	secret, ok := providercore.NewConfig(config).RequiredSecret("token")
	if !ok {
		return url.Values{}
	}
	return url.Values{"key": []string{secret}}
}

func language(code string) string {
	switch strings.ToLower(code) {
	case "vo":
		return "eng"
	case "vf":
		return "fra"
	default:
		return ""
	}
}

func format(raw string) string {
	ext := strings.TrimPrefix(path.Ext(raw), ".")
	if ext == "" || len(ext) > 4 {
		return "srt"
	}
	return ext
}

func firstSubtitle(zr *zip.Reader) ([]byte, error) {
	for _, f := range zr.File {
		name := strings.ToLower(path.Base(f.Name))
		if strings.HasPrefix(name, ".") || !(strings.HasSuffix(name, ".srt") || strings.HasSuffix(name, ".ass") || strings.HasSuffix(name, ".ssa") || strings.HasSuffix(name, ".vtt")) {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return io.ReadAll(io.LimitReader(r, maxBody))
	}
	return nil, fmt.Errorf("%w: archive contains no subtitle", providercore.ErrProviderBrokenUpstream)
}

func do(ctx context.Context, service providercore.Service, config providercore.Config, method, endpoint string, q url.Values, download bool) ([]byte, *http.Response, error) {
	if _, ok := providercore.NewConfig(config).RequiredSecret("token"); !ok && !download {
		return nil, nil, fmt.Errorf("%w: token is required", providercore.ErrProviderPrerequisiteMissing)
	}
	raw := endpoint
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		u, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultURL))
		u.Path = path.Join(u.Path, endpoint)
		u.RawQuery = q.Encode()
		raw = u.String()
	}
	req, err := http.NewRequestWithContext(ctx, method, raw, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "Sub-Zero/2")
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
	if resp.StatusCode == http.StatusBadRequest {
		var parsed responseEnvelope
		if json.Unmarshal(data, &parsed) == nil {
			for _, e := range parsed.Errors {
				if e.Code == 4001 {
					return []byte(`{"errors":[{"code":4001}]}`), resp, nil
				}
				if e.Code == 1001 {
					return nil, resp, fmt.Errorf("%w: invalid token", providercore.ErrProviderPrerequisiteMissing)
				}
			}
		}
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, resp, fmt.Errorf("%w: http status %d", providercore.ErrProviderPrerequisiteMissing, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return data, resp, nil
}
