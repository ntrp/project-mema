package subtis

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
)

const defaultBaseURL = "https://api.subt.is/v1"

type adapter struct{}
type apiResponse struct {
	Subtitle struct{ Link string `json:"subtitle_link"` } `json:"subtitle"`
	Title    struct{ Name string `json:"title_name"` } `json:"title"`
}

func init() { providers.Register("subtis", adapter{}) }

func (adapter) Test(context.Context, providercore.Service, providercore.Config) error { return nil }

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if request.MediaType != "" && request.MediaType != "movie" { return nil, nil }
	filename := path.Base(firstNonEmpty(request.FilePath, request.MediaContext.File.Path, request.MediaContext.File.Name, request.Title))
	if filename == "." || filename == "/" || filename == "" { return nil, nil }
	steps := cascade(config, request, filename)
	for _, step := range steps {
		link, title, ok, err := fetch(ctx, service, step.url)
		if err != nil { return nil, err }
		if !ok { continue }
		if title == "" { title = "Unknown" }
		return []providercore.Candidate{{ProviderName: "subtis", LanguageID: firstNonEmpty(request.LanguageID, "es"), Format: "srt", ReleaseName: title + step.suffix, SourceURL: link, SourceRef: step.url}}, nil
	}
	return nil, nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, _ providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(candidate.SourceURL) == "" { return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing) }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, candidate.SourceURL, nil)
	if err != nil { return providercore.Download{}, err }
	resp, err := service.DoProviderRequest(req, "subtis", true)
	if err != nil { return providercore.Download{}, err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 { return providercore.Download{}, fmt.Errorf("Subt.is status %d", resp.StatusCode) }
	data, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil { return providercore.Download{}, err }
	if len(data) == 0 { return providercore.Download{}, fmt.Errorf("empty subtitle content") }
	return providercore.Download{Content: data, URL: candidate.SourceURL}, nil
}

type step struct{ url, suffix string }

func cascade(config providercore.Config, request providercore.SearchRequest, filename string) []step {
	base := providercore.NewConfig(config).BaseURL(defaultBaseURL)
	var out []step
	if h := hashFor(request); h != "" { out = append(out, step{base + "/subtitle/find/file/hash/" + h, ""}) }
	if request.MediaContext.File.SizeBytes > 0 { out = append(out, step{fmt.Sprintf("%s/subtitle/find/file/bytes/%d", base, request.MediaContext.File.SizeBytes), ""}) }
	out = append(out, step{base + "/subtitle/find/file/name/" + url.PathEscape(filename), ""})
	out = append(out, step{base + "/subtitle/file/alternative/" + url.PathEscape(filename), " [fuzzy match]"})
	return out
}

func fetch(ctx context.Context, service providercore.Service, raw string) (string, string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil { return "", "", false, err }
	req.Header.Set("Accept", "application/json"); req.Header.Set("User-Agent", "Bazarr/Subtis/0.9.2")
	resp, err := service.DoProviderRequest(req, "subtis", false)
	if err != nil { return "", "", false, err }
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound { return "", "", false, nil }
	if resp.StatusCode < 200 || resp.StatusCode > 299 { return "", "", false, fmt.Errorf("Subt.is status %d", resp.StatusCode) }
	var parsed apiResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(&parsed); err != nil { return "", "", false, nil }
	if parsed.Subtitle.Link == "" { return "", "", false, nil }
	return parsed.Subtitle.Link, parsed.Title.Name, true, nil
}

func hashFor(request providercore.SearchRequest) string {
	if h := request.MediaContext.File.Hashes["opensubtitles"]; h != "" { return h }
	file := firstNonEmpty(request.FilePath, request.MediaContext.File.Path)
	if file == "" { return "" }
	h, err := opensubtitlesHash(file)
	if err != nil { return "" }
	return h
}

func opensubtitlesHash(file string) (string, error) {
	f, err := os.Open(file); if err != nil { return "", err }
	defer f.Close()
	st, err := f.Stat(); if err != nil || st.Size() <= 0 { return "", err }
	size := st.Size(); chunk := int64(65536); if size < chunk { chunk = size }
	sum := uint64(size) + checksum(f, 0, chunk) + checksum(f, size-chunk, chunk)
	buf := make([]byte, 8); binary.BigEndian.PutUint64(buf, sum)
	return hex.EncodeToString(buf), nil
}

func checksum(f *os.File, offset, n int64) uint64 {
	buf := make([]byte, n); f.ReadAt(buf, offset)
	for len(buf)%8 != 0 { buf = append(buf, 0) }
	var sum uint64
	for i := 0; i < len(buf); i += 8 { sum += binary.LittleEndian.Uint64(buf[i:i+8]) }
	return sum
}

func firstNonEmpty(values ...string) string { for _, v := range values { if strings.TrimSpace(v) != "" { return strings.TrimSpace(v) } }; return "" }
