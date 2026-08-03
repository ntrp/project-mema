package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"media-manager/internal/ratelimit"
)

func (s *Service) doJSON(ctx context.Context, config Config, method string, endpoint string, body any, target any) error {
	if err := s.wait(ctx, config); err != nil {
		return err
	}
	err := s.doJSONOnce(ctx, config, method, endpoint, body, target)
	if retry, wait := retryAfter(err); retry {
		if wait > 5*time.Second {
			return err
		}
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		if waitErr := s.wait(ctx, config); waitErr != nil {
			return waitErr
		}
		return s.doJSONOnce(ctx, config, method, endpoint, body, target)
	}
	return err
}

func (s *Service) doJSONOnce(ctx context.Context, config Config, method string, endpoint string, body any, target any) error {
	if config.Type == "tmdb" &&
		(config.APIKey == nil || strings.TrimSpace(*config.APIKey) == "") &&
		(config.AccessToken == nil || strings.TrimSpace(*config.AccessToken) == "") {
		return ErrCredentialsRequired
	}

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if config.AccessToken != nil && strings.TrimSpace(*config.AccessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(*config.AccessToken))
	} else if config.Type == "tmdb" && config.APIKey != nil && strings.TrimSpace(*config.APIKey) != "" {
		values := req.URL.Query()
		values.Set("api_key", strings.TrimSpace(*config.APIKey))
		req.URL.RawQuery = values.Encode()
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return rateLimitedError{retryAfter: ratelimit.DelayFromHeaders(resp.Header)}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ProviderHTTPError{StatusCode: resp.StatusCode}
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func (s *Service) wait(ctx context.Context, config Config) error {
	interval := 500 * time.Millisecond
	key := config.ID.String()
	switch config.Type {
	case "tvdb":
		interval = time.Second
	case "anidb":
		// AniDB requires at least two seconds between all HTTP API requests.
		interval = s.anidbInterval
		key = "anidb"
	}

	s.mu.Lock()
	now := time.Now()
	if config.Type == "anidb" && now.Before(s.anidbCooldownUntil) {
		wait := time.Until(s.anidbCooldownUntil)
		s.mu.Unlock()
		return rateLimitedError{retryAfter: wait}
	}
	next := s.lastByID[key].Add(interval)
	if next.Before(now) {
		next = now
	}
	s.lastByID[key] = next
	s.mu.Unlock()
	wait := time.Until(next)
	if wait <= 0 {
		return nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (s *Service) blockAniDB(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	until := time.Now().Add(duration)
	if until.After(s.anidbCooldownUntil) {
		s.anidbCooldownUntil = until
	}
}

type rateLimitedError struct {
	retryAfter time.Duration
}

func (e rateLimitedError) Error() string {
	return "metadata provider rate limit reached"
}

func retryAfter(err error) (bool, time.Duration) {
	var rateErr rateLimitedError
	if errors.As(err, &rateErr) {
		if rateErr.retryAfter <= 0 {
			return true, time.Second
		}
		return true, rateErr.retryAfter
	}
	return false, 0
}

func IsRateLimited(err error) bool {
	var rateErr rateLimitedError
	return errors.As(err, &rateErr)
}

func ProviderStatusCode(err error) (int, bool) {
	var providerErr ProviderHTTPError
	if errors.As(err, &providerErr) {
		return providerErr.StatusCode, true
	}
	return 0, false
}
