package providers

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"media-manager/internal/subtitles/providercore"
)

const bsplayerKey = "bsplayer"

func init() { Register(bsplayerKey, bsplayerAdapter{}) }

type bsplayerAdapter struct{}

func (bsplayerAdapter) Test(context.Context, providercore.Service, providercore.Config) error {
	return fmt.Errorf("%w: BSPlayer SOAP search is disabled in Bazarr because the upstream API is unreliable", providercore.ErrProviderBrokenUpstream)
}

func (bsplayerAdapter) Search(context.Context, providercore.Service, providercore.Config, providercore.SearchRequest) ([]providercore.Candidate, error) {
	return nil, fmt.Errorf("%w: BSPlayer SOAP search is disabled in Bazarr because the upstream API is unreliable", providercore.ErrProviderBrokenUpstream)
}

func (bsplayerAdapter) Download(ctx context.Context, service providercore.Service, _ providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	link := strings.TrimSpace(candidate.SourceURL)
	if link == "" {
		return providercore.Download{}, fmt.Errorf("%w: bsplayer candidate has no download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	request.Header.Set("User-Agent", "Mozilla/4.0 (compatible; Synapse)")
	response, err := service.DoProviderRequest(request, bsplayerKey, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("provider returned HTTP %d", response.StatusCode)
	}
	reader, err := gzip.NewReader(io.LimitReader(response.Body, providerReadLimit+1))
	if err != nil {
		return providercore.Download{}, err
	}
	defer reader.Close()
	content, err := io.ReadAll(io.LimitReader(reader, providerReadLimit+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(content) > providerReadLimit {
		return providercore.Download{}, fmt.Errorf("provider response size limit exceeded")
	}
	return providercore.Download{Content: content, URL: link}, nil
}
