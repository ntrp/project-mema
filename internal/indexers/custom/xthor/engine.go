package xthor

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://api.xthor.tk"

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
		return engine.FailedResult("Invalid Xthor request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "Xthor")
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
	if err := checkAPIState(decoded.Error); err != nil {
		return nil, err
	}
	baseURL := detailsBaseURL(config)
	releases := make([]engine.Release, 0, len(decoded.Torrents))
	for _, item := range decoded.Torrents {
		if item.Category == 106 {
			continue
		}
		release := item.release(config, baseURL)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL))
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("passkey", passkey(config))
	if search := strings.TrimSpace(query); search != "" {
		values.Set("search", strings.ReplaceAll(search, "'", ""))
	}
	if len(config.Categories) > 0 {
		categories := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, strconv.FormatInt(int64(category), 10))
		}
		values.Set("category", strings.Join(categories, "+"))
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		values.Set("freeleech", "1")
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func passkey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "passkey", "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func detailsBaseURL(config engine.Config) string {
	value := common.BaseURL(config, defaultBaseURL)
	value = strings.Replace(value, "://api.", "://", 1)
	return strings.TrimRight(value, "/")
}

type response struct {
	Error    apiError  `json:"Error"`
	Torrents []torrent `json:"Torrents"`
}

type apiError struct {
	Code  int    `json:"Code"`
	Descr string `json:"Descr"`
}

type torrent struct {
	ID             int    `json:"Id"`
	Category       int    `json:"Category"`
	Seeders        int32  `json:"Seeders"`
	Leechers       int32  `json:"Leechers"`
	Name           string `json:"Name"`
	TimesCompleted int32  `json:"Times_completed"`
	Size           int64  `json:"Size"`
	Added          int64  `json:"Added"`
	Freeleech      int    `json:"Freeleech"`
	NumFiles       int32  `json:"Numfiles"`
	DownloadLink   string `json:"Download_link"`
	TMDBID         int    `json:"Tmdb_id"`
}

func (t torrent) release(config engine.Config, baseURL string) engine.Release {
	id := strconv.Itoa(t.ID)
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": id})
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Name),
		DownloadURL:     strings.TrimSpace(t.DownloadLink),
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     common.UnixTime(t.Added),
	}
}

func checkAPIState(state apiError) error {
	switch state.Code {
	case 0, 2, 3:
		return nil
	case 1:
		return fmt.Errorf("passkey not found in tracker database")
	case 4:
		return fmt.Errorf("tracker is under DDOS attack, API disabled")
	case 8:
		return fmt.Errorf("triggered antispam protection")
	default:
		return fmt.Errorf("unknown xthor api state %d", state.Code)
	}
}
