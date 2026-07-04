package subsplease

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://subsplease.org/"

var (
	groupRegex      = regexp.MustCompile(`(?i)\[?SubsPlease\]?\s*`)
	resolutionRegex = regexp.MustCompile(`(?i)\d{3,4}p`)
	sizeRegex       = regexp.MustCompile(`(?i)[&?]xl=(\d+)`)
)

type Engine struct {
	client engine.HTTPDoer
}

func New(clients ...engine.HTTPDoer) *Engine {
	var client engine.HTTPDoer
	if len(clients) > 0 {
		client = clients[0]
	}
	return &Engine{client: client}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	baseURL := common.BaseURL(config, defaultBaseURL)
	return common.TestURL(ctx, e.client, baseURL+"/", "SubsPlease")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	searchTerm := groupRegex.ReplaceAllString(strings.TrimSpace(query), "")
	searchTerm = resolutionRegex.ReplaceAllString(searchTerm, "")
	endpoint, err := common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/api/", map[string]string{
		"tz": "UTC",
		"f":  "search",
		"s":  strings.TrimSpace(searchTerm),
	})
	if err != nil {
		return nil, err
	}
	payload := map[string]release{}
	if err := common.GetJSON(ctx, e.client, endpoint, &payload); err != nil {
		return nil, err
	}
	releases := []engine.Release{}
	for _, item := range payload {
		releases = append(releases, item.releases(config, common.BaseURL(config, defaultBaseURL))...)
	}
	return releases, nil
}

type release struct {
	ReleaseDate time.Time  `json:"release_date"`
	Show        string     `json:"show"`
	Episode     string     `json:"episode"`
	Downloads   []download `json:"downloads"`
	Page        string     `json:"page"`
}

type download struct {
	Resolution string `json:"res"`
	Magnet     string `json:"magnet"`
}

func (r release) releases(config engine.Config, baseURL string) []engine.Release {
	out := make([]engine.Release, 0, len(r.Downloads))
	for _, download := range r.Downloads {
		magnet := strings.TrimSpace(download.Magnet)
		if magnet == "" {
			continue
		}
		title := "[SubsPlease] " + strings.TrimSpace(r.Show) + " - " + strings.TrimSpace(r.Episode) + " (" + strings.TrimSpace(download.Resolution) + "p)"
		published := r.ReleaseDate.UTC()
		seeders := int32(1)
		peers := int32(2)
		out = append(out, engine.Release{
			IndexerID:       config.ID,
			IndexerName:     config.Name,
			IndexerProtocol: config.Protocol,
			Title:           title,
			DownloadURL:     magnet,
			InfoURL:         baseURL + "/shows/" + strings.Trim(strings.TrimSpace(r.Page), "/") + "/",
			GUID:            magnet,
			SizeBytes:       releaseSize(download),
			Seeders:         &seeders,
			Peers:           &peers,
			PublishedAt:     &published,
		})
	}
	return out
}

func releaseSize(download download) int64 {
	if match := sizeRegex.FindStringSubmatch(download.Magnet); len(match) == 2 {
		if parsed, err := strconv.ParseInt(match[1], 10, 64); err == nil && parsed > 0 {
			return parsed
		}
	}
	switch strings.TrimSpace(download.Resolution) {
	case "1080":
		return 1395864371
	case "720":
		return 734003200
	case "480":
		return 367001600
	default:
		return 1073741824
	}
}
