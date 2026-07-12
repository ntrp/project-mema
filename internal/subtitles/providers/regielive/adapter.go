package regielive

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providercore"
)

const providerKey = "regielive"
const defaultBaseURL = "https://api.regielive.ro/bazarr/search.php"
const apiHeader = "API-BAZARR-YTZ-SL"
const maxBody = 10 << 20

var adapter Adapter

type Adapter struct{}

type searchResponse struct {
	Rezultate map[string]struct {
		Subtitrari map[string]struct {
			Titlu  string `json:"titlu"`
			URL    string `json:"url"`
			Rating struct{ Nota int `json:"nota"` } `json:"rating"`
		} `json:"subtitrari"`
	} `json:"rezultate"`
}

func init() { providers.Register(providerKey, adapter) }

func (Adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	_, err := request(ctx, service, config, searchURL(config, url.Values{"nume": {"test"}}), false, nil)
	return classify(err)
}

func (Adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, requestData providercore.SearchRequest) ([]providercore.Candidate, error) {
	params := url.Values{"nume": {requestData.Title}}
	if requestData.Year != nil {
		params.Set("an", strconv.Itoa(int(*requestData.Year)))
	}
	if requestData.MediaType == "serie" {
		if requestData.SeasonNumber != nil { params.Set("sezon", strconv.Itoa(int(*requestData.SeasonNumber))) }
		if requestData.EpisodeNumber != nil { params.Set("episod", strconv.Itoa(int(*requestData.EpisodeNumber))) }
	}
	data, err := request(ctx, service, config, searchURL(config, params), false, nil)
	if err != nil { return nil, classify(err) }
	var parsed searchResponse
	if err := json.Unmarshal(data, &parsed); err != nil { return nil, err }
	out := []providercore.Candidate{}
	for _, film := range parsed.Rezultate {
		for _, sub := range film.Subtitrari {
			out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: "ro", Format: "zip", ReleaseName: sub.Titlu, DownloadCount: sub.Rating.Nota, SourceURL: sub.URL})
		}
	}
	return out, nil
}

func (Adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36",
		"Accept": "application/octet-stream, */*",
		"Accept-Language": "en-US,en;q=0.9",
		"Referer": "https://subtitrari.regielive.ro",
	}
	data, err := request(ctx, service, config, candidate.SourceURL, true, headers)
	if err != nil { return providercore.Download{}, classify(err) }
	if strings.TrimSpace(string(data)) == "500" {
		return providercore.Download{}, fmt.Errorf("regielive download failed: server returned HTTP 500")
	}
	content, err := subtitleFromZip(data)
	if err != nil { return providercore.Download{}, err }
	return providercore.Download{Content: content, URL: candidate.SourceURL}, nil
}

func request(ctx context.Context, service providercore.Service, config providercore.Config, raw string, download bool, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil { return nil, err }
	if !download { req.Header.Set("RL-API", apiHeader) }
	for k, v := range headers { req.Header.Set(k, v) }
	resp, err := service.DoProviderRequest(req, providerKey, download)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil { return nil, err }
	if len(data) > maxBody { return nil, fmt.Errorf("response size limit exceeded") }
	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 { return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, strings.TrimSpace(string(data))) }
	return data, nil
}

func searchURL(config providercore.Config, params url.Values) string {
	raw := providercore.NewConfig(config).BaseURL(defaultBaseURL)
	u, _ := url.Parse(raw)
	u.RawQuery = params.Encode()
	return u.String()
}

func subtitleFromZip(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil { return nil, fmt.Errorf("regielive download failed: provider returned an invalid archive payload") }
	for _, file := range zr.File {
		if strings.HasPrefix(filepath.Base(file.Name), ".") || !isSubtitle(file.Name) { continue }
		rc, err := file.Open()
		if err != nil { return nil, err }
		content, readErr := io.ReadAll(rc)
		closeErr := rc.Close()
		if readErr != nil { return nil, readErr }
		if closeErr != nil { return nil, closeErr }
		return content, nil
	}
	return nil, fmt.Errorf("regielive download failed: archive did not contain a subtitle file")
}

func isSubtitle(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".srt", ".sub", ".ssa", ".ass", ".vtt": return true
	default: return false
	}
}

func classify(err error) error {
	if err == nil { return nil }
	if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "301") {
		return fmt.Errorf("%w: %v", providercore.ErrProviderBrokenUpstream, err)
	}
	return err
}
