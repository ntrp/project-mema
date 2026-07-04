package animebytes

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://animebytes.tv/"

var trailingNumber = regexp.MustCompile(`\W(?:S\d\d?E)?\d?\d$|\W\d+$`)

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
	endpoint, err := searchURL(config, "test", "series")
	if err != nil {
		return engine.FailedResult("Invalid AnimeBytes request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "AnimeBytes")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query, mediaType)
	if err != nil {
		return nil, err
	}
	var decoded response
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
		return nil, err
	}
	if decoded.Error != "" {
		return nil, fmt.Errorf("animebytes api error: %s", decoded.Error)
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := []engine.Release{}
	for _, group := range decoded.Groups {
		for _, item := range group.Torrents {
			if common.FieldBool(config, "freeleechOnly") && item.RawDownMultiplier != 0 {
				continue
			}
			release := item.release(config, baseURL, group)
			if release.Title != "" && release.DownloadURL != "" {
				releases = append(releases, release)
			}
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string, mediaType string) (string, error) {
	searchType := "anime"
	if strings.EqualFold(mediaType, "music") {
		searchType = "music"
	}
	values := map[string]string{
		"username":     common.FieldString(config, "username"),
		"torrent_pass": common.FieldString(config, "passkey", "torrentPass", "torrent_pass"),
		"sort":         "grouptime",
		"way":          "desc",
		"type":         searchType,
		"searchstr":    cleanSearch(query),
		"limit":        "50",
	}
	for _, category := range config.Categories {
		values[strconv.FormatInt(int64(category), 10)] = "1"
	}
	if common.FieldBool(config, "freeleechOnly") {
		values["freeleech"] = "1"
	}
	if common.FieldBool(config, "excludeHentai") && searchType == "anime" {
		values["hentai"] = "0"
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/scrape.php", values)
}

func cleanSearch(query string) string {
	query = trailingNumber.ReplaceAllString(strings.TrimSpace(query), "")
	query = regexp.MustCompile(`(?i)\bThe Movie$`).ReplaceAllString(query, "")
	return strings.TrimSpace(query)
}

type response struct {
	Matches int     `json:"Matches"`
	Error   string  `json:"Error"`
	Groups  []group `json:"Groups"`
}

type group struct {
	ID           int64             `json:"ID"`
	CategoryName string            `json:"CategoryName"`
	FullName     string            `json:"FullName"`
	GroupName    string            `json:"GroupName"`
	SeriesName   string            `json:"SeriesName"`
	Year         int               `json:"Year,string"`
	Synonyms     map[string]string `json:"SynonymnsV2"`
	Torrents     []torrent         `json:"Torrents"`
}

type torrent struct {
	ID                int64   `json:"ID"`
	RawDownMultiplier float64 `json:"RawDownMultiplier"`
	Link              string  `json:"Link"`
	Property          string  `json:"Property"`
	Snatched          int32   `json:"Snatched"`
	Seeders           int     `json:"Seeders"`
	Leechers          int     `json:"Leechers"`
	Size              int64   `json:"Size"`
	FileCount         int32   `json:"FileCount"`
	UploadTime        string  `json:"UploadTime"`
}

func (t torrent) release(config engine.Config, baseURL string, g group) engine.Release {
	title := g.title()
	if t.Property != "" {
		title += " [" + strings.TrimSpace(t.Property) + "]"
	}
	infoURL := common.ResolveURL(baseURL, "torrent/"+strconv.FormatInt(t.ID, 10)+"/group")
	published := parseAnimeBytesTime(t.UploadTime)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     t.Link,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(t.Seeders),
		Peers:           common.Int32Ptr(t.Seeders + t.Leechers),
		PublishedAt:     published,
	}
}

func (g group) title() string {
	title := engine.FirstNonEmpty(g.SeriesName, g.FullName, g.GroupName)
	if g.Year > 0 && !strings.Contains(title, strconv.Itoa(g.Year)) {
		title += " (" + strconv.Itoa(g.Year) + ")"
	}
	return strings.TrimSpace(title)
}

func parseAnimeBytesTime(value string) *time.Time {
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(value), time.UTC)
	if err != nil {
		return common.ParseFlexibleTime(value)
	}
	return &parsed
}
