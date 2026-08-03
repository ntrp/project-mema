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
const (
	anidbProtocolVersion = "1"
	anidbBanCooldown     = 24 * time.Hour
)

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
	posterPath := p.cacheArtwork(ctx, anime.Picture)
	return []SearchResult{anidbSearchResult(anime, posterPath)}, nil
}

func (p anidbProvider) Details(ctx context.Context, request DetailsRequest) (MediaResource, error) {
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
	posterPath := p.cacheArtwork(ctx, anime.Picture)
	resource := anidbDetails(anime, posterPath)
	p.cacheCastArtwork(ctx, anime.Characters, resource.Cast)
	return resource, nil
}

func (p anidbProvider) cacheArtwork(ctx context.Context, picture string) *string {
	client, version, err := anidbClientIdentity(p.config)
	if err != nil {
		return nil
	}
	return p.service.cacheAniDBArtwork(ctx, picture, client+"/"+version)
}

func (p anidbProvider) cacheCastArtwork(ctx context.Context, characters []anidbCharacter, cast []Person) {
	pictures := map[string]string{}
	for _, character := range characters {
		for _, actor := range character.Seiyuu {
			if actor.ID > 0 && strings.TrimSpace(actor.Picture) != "" {
				pictures[strconv.Itoa(actor.ID)] = actor.Picture
			}
		}
	}
	profiles := map[string]*string{}
	for id, picture := range pictures {
		profiles[id] = p.cacheArtwork(ctx, picture)
	}
	for index := range cast {
		if cast[index].ExternalID != nil {
			cast[index].ProfilePath = profiles[*cast[index].ExternalID]
		}
	}
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

	return p.requestOnce(ctx, u.String(), target)
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
	req.Header.Set("User-Agent", req.URL.Query().Get("client")+"/"+req.URL.Query().Get("clientver"))
	resp, err := p.service.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		wait := ratelimit.DelayFromHeaders(resp.Header)
		if wait <= 0 {
			wait = anidbBanCooldown
		}
		p.service.blockAniDB(wait)
		return rateLimitedError{retryAfter: wait}
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
