package shizaproject

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://shiza-project.com/"

var trailingEpisode = regexp.MustCompile(`(?:[SsEe]?\d{1,4}){1,2}$`)

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
	_, err := e.Search(ctx, config, "test", "")
	if err != nil {
		return engine.FailedResult("Invalid Shizaproject request", "error", err.Error())
	}
	return engine.SuccessResult("Shizaproject indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	var decoded response
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
		return nil, err
	}
	if len(decoded.Errors) > 0 {
		return nil, fmt.Errorf("shizaproject api error: %s", decoded.Errors[0].Message)
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := []engine.Release{}
	for _, edge := range decoded.Data.Releases.Edges {
		node := edge.Node
		for _, item := range node.Torrents {
			release := item.release(config, baseURL, node)
			if release.Title != "" && release.DownloadURL != "" {
				releases = append(releases, release)
			}
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	variables, err := json.Marshal(map[string]any{
		"first": 50,
		"query": strings.TrimSpace(trailingEpisode.ReplaceAllString(strings.TrimSpace(query), "")),
	})
	if err != nil {
		return "", err
	}
	queryText := `query fetchReleases($first: Int, $query: String) { releases(first: $first, query: $query) { edges { node { name type originalName alternativeNames publishedAt slug torrents { synopsis downloaded seeders leechers size magnetUri updatedAt file { url } videoQualities } } } } }`
	values := map[string]string{"query": queryText, "variables": string(variables)}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/graphql", values)
	if err != nil {
		return "", err
	}
	return endpoint, nil
}

type response struct {
	Data struct {
		Releases struct {
			Edges []edge `json:"edges"`
		} `json:"releases"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type edge struct {
	Node node `json:"node"`
}

type node struct {
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	OriginalName     string     `json:"originalName"`
	AlternativeNames []string   `json:"alternativeNames"`
	PublishedAt      *time.Time `json:"publishedAt"`
	Slug             string     `json:"slug"`
	Torrents         []torrent  `json:"torrents"`
}

type torrent struct {
	Synopsis       string    `json:"synopsis"`
	Downloaded     int       `json:"downloaded"`
	Seeders        int       `json:"seeders"`
	Leechers       int       `json:"leechers"`
	Size           int64     `json:"size"`
	MagnetURI      string    `json:"magnetUri"`
	UpdatedAt      time.Time `json:"updatedAt"`
	VideoQualities []string  `json:"videoQualities"`
	File           struct {
		URL string `json:"url"`
	} `json:"file"`
}

func (t torrent) release(config engine.Config, baseURL string, n node) engine.Release {
	title := composeTitle(n, t)
	infoURL := common.ResolveURL(baseURL, "/releases/"+url.PathEscape(n.Slug)+"/")
	published := t.UpdatedAt
	if n.PublishedAt != nil && n.PublishedAt.After(published) {
		published = *n.PublishedAt
	}
	downloadURL := engine.FirstNonEmpty(t.File.URL, t.MagnetURI)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            downloadURL,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(t.Seeders),
		Peers:           common.Int32Ptr(t.Seeders + t.Leechers),
		PublishedAt:     &published,
	}
}

func composeTitle(n node, t torrent) string {
	names := []string{n.Name, n.OriginalName}
	names = append(names, n.AlternativeNames...)
	parts := []string{}
	seen := map[string]bool{}
	for _, name := range names {
		name = strings.TrimSpace(name)
		key := strings.ToLower(name)
		if name != "" && !seen[key] {
			seen[key] = true
			parts = append(parts, name)
		}
	}
	title := strings.TrimSpace(strings.Join(parts, " / ") + " " + strings.TrimSpace(t.Synopsis))
	if len(t.VideoQualities) > 0 {
		title += " [" + strings.Join(t.VideoQualities, " ") + "]"
	}
	return strings.TrimSpace(title)
}
