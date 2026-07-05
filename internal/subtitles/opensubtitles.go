package subtitles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrCredentialsRequired = errors.New("subtitle provider API key is required")
	ErrUnsupportedProvider = errors.New("unsupported subtitle provider")
)

func (s *Service) testOpenSubtitles(ctx context.Context, config Config) error {
	if config.Type != "opensubtitles" {
		return ErrUnsupportedProvider
	}
	if config.APIKey == nil || strings.TrimSpace(*config.APIKey) == "" {
		return ErrCredentialsRequired
	}
	base, err := url.Parse(strings.TrimSpace(config.BaseURL))
	if err != nil || base.Scheme == "" || base.Host == "" {
		return errors.New("subtitle provider base URL is invalid")
	}
	endpoint := base.JoinPath("api", "v1", "infos", "languages")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Api-Key", strings.TrimSpace(*config.APIKey))
	req.Header.Set("User-Agent", "project-mema")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("subtitle provider returned HTTP %d", resp.StatusCode)
	}
	return nil
}
