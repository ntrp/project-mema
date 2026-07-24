package metadata

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/anidb"
	"media-manager/internal/ratelimit"
)

// AniDB HTTP API contract: https://wiki.anidb.net/HTTP_API_Definition
// The registered application client name/version are required; user login is not.
const anidbProtocolVersion = "1"

var (
	ErrAniDBClientIdentityRequired = errors.New("AniDB application client name and version are required")
	ErrAniDBClientNameInvalid      = errors.New("AniDB application client name must contain only lowercase ASCII letters and digits")
)

type anidbProvider struct {
	service *Service
	config  Config
}

func (p anidbProvider) Search(ctx context.Context, request SearchRequest) ([]SearchResult, error) {
	if request.MediaType != "serie" {
		return nil, nil
	}
	aid, err := p.findAnimeID(ctx, request.Query)
	if err != nil {
		return nil, err
	}
	if aid == 0 {
		return []SearchResult{}, nil
	}
	var anime anidbAnime
	if err := p.request(ctx, map[string]string{"request": "anime", "aid": strconv.Itoa(aid)}, &anime); err != nil {
		return nil, err
	}
	if anime.ID == 0 || anidbMainTitle(anime) == "" {
		return nil, fmt.Errorf("AniDB returned incomplete anime search data")
	}
	return []SearchResult{anidbSearchResult(anime)}, nil
}

func (p anidbProvider) Details(ctx context.Context, request DetailsRequest) (Details, error) {
	id, err := strconv.Atoi(strings.TrimSpace(request.ExternalID))
	if err != nil || id <= 0 {
		return Details{}, fmt.Errorf("invalid AniDB id")
	}
	var anime anidbAnime
	if err := p.request(ctx, map[string]string{"request": "anime", "aid": strconv.Itoa(id)}, &anime); err != nil {
		return Details{}, err
	}
	if anime.ID == 0 || anidbMainTitle(anime) == "" {
		return Details{}, fmt.Errorf("AniDB returned incomplete anime details")
	}
	return anidbDetails(anime), nil
}

func (p anidbProvider) request(ctx context.Context, params map[string]string, target any) error {
	client, version, err := anidbClientIdentity(p.config)
	if err != nil {
		return err
	}
	base := strings.TrimRight(p.config.BaseURL, "/")
	if base == "" {
		base = "http://api.anidb.net:9001"
	}
	u, err := url.Parse(base + "/httpapi")
	if err != nil {
		return err
	}
	query := u.Query()
	query.Set("client", client)
	query.Set("clientver", version)
	query.Set("protover", anidbProtocolVersion)
	for key, value := range params {
		query.Set(key, strings.TrimSpace(value))
	}
	u.RawQuery = query.Encode()

	err = p.requestOnce(ctx, u.String(), target)
	if retry, wait := retryAfter(err); retry && wait <= 5*time.Second {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		return p.requestOnce(ctx, u.String(), target)
	}
	return err
}

func (p anidbProvider) requestOnce(ctx context.Context, endpoint string, target any) error {
	if err := p.service.wait(ctx, p.config); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/xml")
	resp, err := p.service.httpClient.Do(req)
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
	if err := xml.NewDecoder(io.LimitReader(resp.Body, 16<<20)).Decode(target); err != nil {
		return fmt.Errorf("decode AniDB XML response: %w", err)
	}
	return nil
}

func anidbClientIdentity(config Config) (string, string, error) {
	if config.ClientName == nil || config.ClientVersion == nil {
		return "", "", ErrAniDBClientIdentityRequired
	}
	name := *config.ClientName
	version := strings.TrimSpace(*config.ClientVersion)
	if name == "" {
		return "", "", ErrAniDBClientIdentityRequired
	}
	if !anidb.ValidClientName(name) {
		return "", "", ErrAniDBClientNameInvalid
	}
	parsedVersion, err := strconv.Atoi(version)
	if err != nil || parsedVersion <= 0 {
		return "", "", ErrAniDBClientIdentityRequired
	}
	return name, strconv.Itoa(parsedVersion), nil
}
