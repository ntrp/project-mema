package subtitles

import (
	"net/http"

	"media-manager/internal/subtitles/security"
)

func validateProviderURL(providerKey string, rawURL string, download bool) error {
	return security.ValidateProviderURL(canonicalProviderKey(providerKey), rawURL, download)
}

func (s *Service) doProviderRequest(req *http.Request, providerKey string, download bool) (*http.Response, error) {
	if err := validateProviderURL(providerKey, req.URL.String(), download); err != nil {
		return nil, err
	}
	client := *s.client
	previous := s.client.CheckRedirect
	client.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		from := req.URL.String()
		if len(via) > 0 {
			from = via[len(via)-1].URL.String()
		}
		if err := security.ValidateRedirect(canonicalProviderKey(providerKey), from, next.URL.String(), download); err != nil {
			return err
		}
		if previous != nil {
			return previous(next, via)
		}
		return nil
	}
	return client.Do(req)
}
