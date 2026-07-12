package subsynchro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativeutil"
)

const baseURL = "https://www.subsynchro.com"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subsynchro", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "movie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	form := url.Values{"title": {sr.Title}}
	if sr.Year != nil {
		form.Set("year", strconv.Itoa(int(*sr.Year)))
	}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subsynchro", BaseURL: baseURL, Path: "/include/ajax/subMarin.php", Form: form, Headers: map[string]string{"Referer": nativeutil.Absolute(cfg, baseURL, "/")}})
	if err != nil {
		return nil, err
	}
	return parse(data), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return nativeutil.DownloadSubtitle(ctx, svc, cfg, "subsynchro", baseURL, cand)
}

type response struct {
	Status int      `json:"status"`
	Data   []result `json:"data"`
}

type result struct {
	Release  string `json:"release"`
	Filename string `json:"filename"`
	Download string `json:"telechargement"`
	File     string `json:"fichier"`
}

func parse(data []byte) []providercore.Candidate {
	var body response
	dec := json.NewDecoder(bytes.NewReader(data))
	if dec.Decode(&body) != nil || body.Status != 200 {
		return nil
	}
	out := make([]providercore.Candidate, 0, len(body.Data))
	for _, item := range body.Data {
		if strings.TrimSpace(item.Download) == "" {
			continue
		}
		name := item.Release
		if name == "" {
			name = item.Filename
		}
		out = append(out, providercore.Candidate{ProviderName: "subsynchro", LanguageID: "fre", Format: nativeutil.Format(item.File), ReleaseName: name, SourceURL: item.Download})
	}
	return out
}
