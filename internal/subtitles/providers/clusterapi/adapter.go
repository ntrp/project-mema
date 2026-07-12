package clusterapi

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const maxBody = 10 << 20

type Spec struct {
	Key             string
	DefaultBaseURL  string
	SearchPath      string
	TestPath        string
	RequiredSecret  string
	SecretQueryName string
	SecretHeader    string
	SeriesOnly      bool
	MovieOnly       bool
	RequireIMDb     bool
	Local           bool
	CommandName     string
	CommandArgs     []string
}

type Adapter struct{ Spec Spec }

func (a Adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if err := a.validate(config, providercore.SearchRequest{}); err != nil && !strings.Contains(err.Error(), "unsupported media type") {
		return err
	}
	if a.Spec.CommandName != "" {
		runner := config.CommandRunner
		if runner == nil {
			return fmt.Errorf("%w: %s is required", providercore.ErrProviderPrerequisiteMissing, a.Spec.CommandName)
		}
		if _, err := runner(ctx, a.Spec.CommandName, a.Spec.CommandArgs...); err != nil {
			return fmt.Errorf("%w: %s check failed: %v", providercore.ErrProviderPrerequisiteMissing, a.Spec.CommandName, err)
		}
	}
	if a.Spec.TestPath == "" {
		return nil
	}
	_, _, err := a.do(ctx, service, config, http.MethodGet, a.Spec.TestPath, nil, false)
	return classifyHTTP(err)
}

func (a Adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if a.Spec.Local && a.Spec.Key == "whisperai" {
		if err := a.Test(ctx, service, config); err != nil {
			return nil, err
		}
		return []providercore.Candidate{{ProviderName: a.Spec.Key, LanguageID: request.LanguageID, Format: "srt", ReleaseName: request.Title, SourceURL: a.url(config, a.Spec.SearchPath, request).String()}}, nil
	}
	if err := a.validate(config, request); err != nil {
		return nil, err
	}
	data, _, err := a.do(ctx, service, config, http.MethodGet, a.url(config, a.Spec.SearchPath, request).String(), nil, false)
	if err != nil {
		return nil, classifyHTTP(err)
	}
	return parseCandidates(a.Spec.Key, request.LanguageID, data)
}

func (a Adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	raw := strings.TrimSpace(candidate.SourceURL)
	if raw == "" && candidate.FileID != 0 {
		raw = a.url(config, "/download/"+strconv.FormatInt(candidate.FileID, 10), providercore.SearchRequest{}).String()
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, resp, err := a.do(context.Background(), service, config, http.MethodGet, raw, nil, true)
	if err != nil {
		return providercore.Download{}, classifyHTTP(err)
	}
	name := downloadName(raw, resp)
	member, err := providercore.ExtractSubtitle(name, data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: raw}, nil
}

func (a Adapter) validate(config providercore.Config, request providercore.SearchRequest) error {
	view := providercore.NewConfig(config)
	if a.Spec.RequiredSecret != "" {
		if _, ok := view.RequiredSecret(a.Spec.RequiredSecret); !ok {
			return fmt.Errorf("%w: %s is required", providercore.ErrProviderPrerequisiteMissing, a.Spec.RequiredSecret)
		}
	}
	if a.Spec.MovieOnly && request.MediaType != "" && request.MediaType != "movie" {
		return fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	if a.Spec.SeriesOnly && request.MediaType != "" && request.MediaType != "serie" {
		return fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	if a.Spec.RequireIMDb && imdbID(request) == "" {
		return fmt.Errorf("%w: imdb id is required", providercore.ErrProviderPrerequisiteMissing)
	}
	return nil
}

func (a Adapter) do(ctx context.Context, service providercore.Service, config providercore.Config, method, raw string, body io.Reader, download bool) ([]byte, *http.Response, error) {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = a.absolute(config, raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, raw, body)
	if err != nil {
		return nil, nil, err
	}
	a.authorize(req, config)
	resp, err := service.DoProviderRequest(req, a.Spec.Key, download)
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("http status %d", resp.StatusCode)
	}
	return data, resp, nil
}

func (a Adapter) url(config providercore.Config, endpoint string, request providercore.SearchRequest) *url.URL {
	u, _ := url.Parse(a.absolute(config, endpoint))
	q := u.Query()
	if request.Title != "" {
		q.Set("q", request.Title)
		q.Set("title", request.Title)
	}
	if request.LanguageID != "" {
		q.Set("language", request.LanguageID)
	}
	if request.Year != nil {
		q.Set("year", strconv.Itoa(int(*request.Year)))
	}
	if request.SeasonNumber != nil {
		q.Set("season", strconv.Itoa(int(*request.SeasonNumber)))
	}
	if request.EpisodeNumber != nil {
		q.Set("episode", strconv.Itoa(int(*request.EpisodeNumber)))
	}
	if id := imdbID(request); id != "" {
		q.Set("imdb_id", id)
		q.Set("imdb", id)
	}
	if request.MediaContext.File.SizeBytes > 0 {
		q.Set("size", strconv.FormatInt(request.MediaContext.File.SizeBytes, 10))
	}
	if a.Spec.SecretQueryName != "" {
		if secret, ok := providercore.NewConfig(config).RequiredSecret(a.Spec.RequiredSecret); ok {
			q.Set(a.Spec.SecretQueryName, secret)
		}
	}
	u.RawQuery = q.Encode()
	return u
}

func (a Adapter) absolute(config providercore.Config, endpoint string) string {
	base := providercore.NewConfig(config).BaseURL(a.Spec.DefaultBaseURL)
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}
	u, _ := url.Parse(base)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (a Adapter) authorize(req *http.Request, config providercore.Config) {
	view := providercore.NewConfig(config)
	if a.Spec.SecretHeader != "" {
		if secret, ok := view.RequiredSecret(a.Spec.RequiredSecret); ok {
			req.Header.Set(a.Spec.SecretHeader, secret)
		}
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
}

func imdbID(r providercore.SearchRequest) string {
	if r.MediaContext.ExternalIDs != nil {
		return r.MediaContext.ExternalIDs["imdb"]
	}
	return ""
}
func downloadName(raw string, resp *http.Response) string {
	if resp != nil {
		_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if params["filename"] != "" {
			return params["filename"]
		}
	}
	u, _ := url.Parse(raw)
	if base := path.Base(u.Path); base != "." && base != "/" {
		return base
	}
	return "subtitle.srt"
}
func classifyHTTP(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(msg, "401") || strings.Contains(msg, "403") {
		return fmt.Errorf("%w: %v", providercore.ErrProviderPrerequisiteMissing, err)
	}
	if strings.Contains(msg, "429") {
		return fmt.Errorf("%w: %v", providercore.ErrProviderBrokenUpstream, err)
	}
	return err
}
