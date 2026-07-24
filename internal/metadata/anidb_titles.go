package metadata

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	anidbTitleDumpURL             = "https://anidb.net/api/anime-titles.xml.gz"
	anidbTitleDumpCacheTTL        = 24 * time.Hour
	anidbTitleDumpTimeout         = 15 * time.Second
	anidbTitleDumpMaxCompressed   = 16 << 20
	anidbTitleDumpMaxUncompressed = 64 << 20
)

type anidbTitleCache struct {
	anime     []anidbTitleAnime
	expiresAt time.Time
}

type anidbTitleDump struct {
	Anime []anidbTitleAnime `xml:"anime"`
}

type anidbTitleAnime struct {
	ID     int          `xml:"aid,attr"`
	Titles []anidbTitle `xml:"title"`
}

func (p anidbProvider) findAnimeID(ctx context.Context, query string) (int, error) {
	query = normalizedAniDBTitle(query)
	if query == "" {
		return 0, nil
	}
	anime, err := p.animeTitles(ctx)
	if err != nil {
		return 0, err
	}
	partial := 0
	for _, item := range anime {
		for _, title := range item.Titles {
			candidate := normalizedAniDBTitle(title.Value)
			if candidate == query {
				return item.ID, nil
			}
			if partial == 0 && strings.Contains(candidate, query) {
				partial = item.ID
			}
		}
	}
	return partial, nil
}

func (p anidbProvider) animeTitles(ctx context.Context) ([]anidbTitleAnime, error) {
	endpoint := p.titleDumpEndpoint()
	p.service.mu.Lock()
	cached, ok := p.service.anidbTitleCaches[endpoint]
	p.service.mu.Unlock()
	if ok && time.Now().Before(cached.expiresAt) {
		return cached.anime, nil
	}

	anime, err := p.downloadAnimeTitles(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	p.service.mu.Lock()
	if p.service.anidbTitleCaches == nil {
		p.service.anidbTitleCaches = map[string]anidbTitleCache{}
	}
	p.service.anidbTitleCaches[endpoint] = anidbTitleCache{anime: anime, expiresAt: time.Now().Add(anidbTitleDumpCacheTTL)}
	p.service.mu.Unlock()
	return anime, nil
}

func (p anidbProvider) downloadAnimeTitles(ctx context.Context, endpoint string) ([]anidbTitleAnime, error) {
	ctx, cancel := context.WithTimeout(ctx, anidbTitleDumpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	client, version, err := anidbClientIdentity(p.config)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/gzip, application/xml")
	req.Header.Set("User-Agent", client+"/"+version)
	resp, err := p.service.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ProviderHTTPError{StatusCode: resp.StatusCode}
	}
	compressed, err := io.ReadAll(io.LimitReader(resp.Body, anidbTitleDumpMaxCompressed+1))
	if err != nil {
		return nil, err
	}
	if len(compressed) > anidbTitleDumpMaxCompressed {
		return nil, fmt.Errorf("AniDB anime title dump exceeds compressed size limit")
	}
	archive, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("open AniDB anime title dump: %w", err)
	}
	defer archive.Close()
	contents, err := io.ReadAll(io.LimitReader(archive, anidbTitleDumpMaxUncompressed+1))
	if err != nil {
		return nil, err
	}
	if len(contents) > anidbTitleDumpMaxUncompressed {
		return nil, fmt.Errorf("AniDB anime title dump exceeds uncompressed size limit")
	}
	var dump anidbTitleDump
	if err := xml.Unmarshal(contents, &dump); err != nil {
		return nil, fmt.Errorf("decode AniDB anime title dump: %w", err)
	}
	return dump.Anime, nil
}

func (p anidbProvider) titleDumpEndpoint() string {
	base := strings.TrimRight(strings.TrimSpace(p.config.BaseURL), "/")
	if base == "" || base == "http://api.anidb.net:9001" || base == "https://api.anidb.net:9001" {
		return anidbTitleDumpURL
	}
	return base + "/api/anime-titles.xml.gz"
}

func normalizedAniDBTitle(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
