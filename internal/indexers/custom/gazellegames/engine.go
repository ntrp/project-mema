package gazellegames

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://gazellegames.net/"

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
		return engine.FailedResult("Invalid GazelleGames request", "error", err.Error())
	}
	return engine.SuccessResult("GazelleGames indexer reachable")
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
	req.Header.Set("X-API-Key", apiKey(config))
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	var decoded response
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	if !strings.EqualFold(decoded.Status, "success") {
		return []engine.Release{}, nil
	}
	var groups map[string]group
	if err := json.Unmarshal(decoded.Response, &groups); err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := []engine.Release{}
	for groupID, group := range groups {
		for torrentID, item := range group.Torrents {
			item.ID, _ = strconv.Atoi(torrentID)
			if !strings.EqualFold(item.TorrentType, "TORRENT") {
				continue
			}
			if common.FieldBool(config, "freeleechOnly") && !item.free() {
				continue
			}
			release := item.release(config, baseURL, groupID, group)
			if release.Title != "" && release.DownloadURL != "" {
				releases = append(releases, release)
			}
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	values := map[string]string{
		"request":      "search",
		"search_type":  "torrents",
		"empty_groups": "filled",
		"order_by":     "time",
		"order_way":    "desc",
	}
	if strings.TrimSpace(query) != "" {
		if common.FieldBool(config, "searchGroupNames") {
			values["groupname"] = strings.ReplaceAll(query, ".", " ")
		} else {
			values["searchstr"] = strings.ReplaceAll(query, ".", " ")
		}
	}
	if common.FieldBool(config, "freeleechOnly") {
		values["freetorrent"] = "1"
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/api.php", values)
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

type response struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
}

type group struct {
	Artists  []artist           `json:"Artists"`
	Torrents map[string]torrent `json:"Torrents"`
	Year     int                `json:"Year"`
}

type artist struct {
	Name string `json:"Name"`
}

type torrent struct {
	ID            int       `json:"-"`
	CategoryID    int       `json:"CategoryID"`
	Format        string    `json:"Format"`
	Encoding      string    `json:"Encoding"`
	Language      string    `json:"Language"`
	Region        string    `json:"Region"`
	RemasterYear  string    `json:"RemasterYear"`
	RemasterTitle string    `json:"RemasterTitle"`
	ReleaseTitle  string    `json:"ReleaseTitle"`
	Miscellaneous string    `json:"Miscellaneous"`
	Scene         int       `json:"Scene"`
	Dupable       int       `json:"Dupable"`
	Time          time.Time `json:"Time"`
	TorrentType   string    `json:"TorrentType"`
	FileCount     int32     `json:"FileCount"`
	Size          string    `json:"Size"`
	Snatched      int32     `json:"Snatched"`
	Seeders       int       `json:"Seeders"`
	Leechers      int       `json:"Leechers"`
	FreeTorrent   string    `json:"FreeTorrent"`
	PersonalFL    bool      `json:"PersonalFL"`
	LowSeedFL     bool      `json:"LowSeedFL"`
	GameDoxType   string    `json:"GameDOXType"`
}

func (t torrent) release(config engine.Config, baseURL string, groupID string, g group) engine.Release {
	id := strconv.Itoa(t.ID)
	if t.ID == 0 {
		id = ""
	}
	infoURL, _ := common.URLWithQuery(baseURL, "/torrents.php", map[string]string{"id": groupID, "torrentid": id})
	downloadURL, _ := common.URLWithQuery(baseURL, "/torrents.php", map[string]string{
		"action":       "download",
		"id":           id,
		"authkey":      "prowlarr",
		"torrent_pass": common.FieldString(config, "passkey", "torrentPass", "torrent_pass"),
	})
	published := t.Time
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           t.title(g),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       parseSize(t.Size),
		Seeders:         common.Int32Ptr(t.Seeders),
		Peers:           common.Int32Ptr(t.Seeders + t.Leechers),
		PublishedAt:     &published,
	}
}

func (t torrent) title(g group) string {
	title := strings.TrimSpace(t.ReleaseTitle)
	if g.Year > 0 && title != "" && !strings.Contains(title, strconv.Itoa(g.Year)) {
		title += " (" + strconv.Itoa(g.Year) + ")"
	}
	if t.RemasterTitle != "" {
		title += " [" + strings.TrimSpace(t.RemasterTitle+" "+t.RemasterYear) + "]"
	}
	flags := []string{strings.TrimSpace(t.Format + " " + t.Encoding)}
	for _, artist := range g.Artists {
		if artist.Name != "" {
			flags = append(flags, artist.Name)
		}
	}
	for _, value := range []string{t.Language, t.Region, t.Miscellaneous, t.GameDoxType} {
		if strings.TrimSpace(value) != "" {
			flags = append(flags, strings.TrimSpace(value))
		}
	}
	out := []string{}
	for _, flag := range flags {
		if strings.TrimSpace(flag) != "" {
			out = append(out, strings.TrimSpace(flag))
		}
	}
	if len(out) > 0 {
		title += " [" + strings.Join(out, " / ") + "]"
	}
	return strings.TrimSpace(title)
}

func (t torrent) free() bool {
	return t.LowSeedFL || t.PersonalFL || strings.EqualFold(t.FreeTorrent, "FreeLeech") || strings.EqualFold(t.FreeTorrent, "Neutral")
}

func parseSize(value string) int64 {
	if parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64); err == nil {
		return parsed
	}
	return common.ParseSizeBytes(value)
}
