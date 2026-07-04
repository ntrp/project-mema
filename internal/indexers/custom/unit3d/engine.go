package unit3d

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://unit3d.local"

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
	endpoint, err := searchURL(config, "test")
	if err != nil {
		return engine.FailedResult("Invalid UNIT3D request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "UNIT3D")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	body, err := engine.ReadLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseResponse(config, body)
}

func searchURL(config engine.Config, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/api/torrents/filter")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	if apiKey := common.FieldString(config, "apiKey", "apikey"); apiKey != "" {
		values.Set("api_token", apiKey)
	} else {
		values.Set("api_token", engine.StringValue(config.APIKey))
	}
	if strings.TrimSpace(query) != "" {
		values.Set("name", strings.TrimSpace(query))
	}
	for _, category := range config.Categories {
		values.Add("categories[]", strconv.FormatInt(int64(category), 10))
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func parseResponse(config engine.Config, body []byte) ([]engine.Release, error) {
	var decoded response
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(decoded.Data))
	for _, item := range decoded.Data {
		release := item.toRelease(config)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

type response struct {
	Data []torrent `json:"data"`
}

type torrent struct {
	ID         string     `json:"id"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	Size         int64  `json:"size"`
	Files        int32  `json:"num_file"`
	Grabs        int32  `json:"times_completed"`
	Seeders      int32  `json:"seeders"`
	Leechers     int32  `json:"leechers"`
	CreatedAt    string `json:"created_at"`
	DownloadLink string `json:"download_link"`
	DetailsLink  string `json:"details_link"`
	InfoHash     string `json:"info_hash"`
	Freeleech    bool   `json:"freeleech"`
	DoubleUpload bool   `json:"double_upload"`
	IMDBID       string `json:"imdb_id"`
	TMDBID       string `json:"tmdb_id"`
	TVDBID       string `json:"tvdb_id"`
	ReleaseYear  int    `json:"release_year"`
	Encode       string `json:"encode"`
	Resolution   string `json:"resolution"`
	Uploader     string `json:"uploader"`
}

func (t torrent) toRelease(config engine.Config) engine.Release {
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Attributes.Name),
		DownloadURL:     strings.TrimSpace(t.Attributes.DownloadLink),
		InfoURL:         strings.TrimSpace(t.Attributes.DetailsLink),
		GUID:            engine.FirstNonEmpty(t.Attributes.DetailsLink, t.ID),
		SizeBytes:       t.Attributes.Size,
		Seeders:         common.Int32Ptr(int(t.Attributes.Seeders)),
		Peers:           common.Int32Ptr(int(t.Attributes.Seeders + t.Attributes.Leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.Attributes.CreatedAt),
	}
}
