package broadcasthenet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://api.broadcasthe.net"

var protocolPrefix = regexp.MustCompile(`^https?:`)

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
		return engine.FailedResult("Invalid BroadcastheNet request", "error", err.Error())
	}
	return engine.SuccessResult("BroadcastheNet indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint := common.BaseURL(config, defaultBaseURL)
	payload := rpcPayload(config, query)
	var decoded rpcResponse
	if err := postRPC(ctx, e.client, endpoint, payload, &decoded); err != nil {
		return nil, err
	}
	if decoded.Error != nil {
		return nil, fmt.Errorf("broadcasthe.net api error: %v", decoded.Error)
	}
	releases := make([]engine.Release, 0, len(decoded.Result.Torrents))
	for _, item := range decoded.Result.Torrents {
		release := item.release(config, endpoint)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func rpcPayload(config engine.Config, query string) map[string]any {
	params := map[string]any{}
	if search := strings.TrimSpace(query); search != "" {
		params["Search"] = strings.ReplaceAll(search, " ", "%")
	}
	return map[string]any{
		"jsonrpc": "2.0",
		"method":  "getTorrents",
		"params":  []any{apiKey(config), params, 100, 0},
		"id":      "prowlarr",
	}
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func postRPC(ctx context.Context, client engine.HTTPDoer, endpoint string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json-rpc, application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.HTTPStatusError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

type rpcResponse struct {
	Result torrents `json:"result"`
	Error  any      `json:"error"`
}

type torrents struct {
	Torrents map[string]torrent `json:"torrents"`
	Results  int                `json:"results"`
}

type torrent struct {
	GroupID     int    `json:"GroupID"`
	TorrentID   int    `json:"TorrentID"`
	Category    string `json:"Category"`
	Snatched    int32  `json:"Snatched"`
	Seeders     int32  `json:"Seeders"`
	Leechers    int32  `json:"Leechers"`
	Source      string `json:"Source"`
	Container   string `json:"Container"`
	Codec       string `json:"Codec"`
	Resolution  string `json:"Resolution"`
	Origin      string `json:"Origin"`
	ReleaseName string `json:"ReleaseName"`
	Size        int64  `json:"Size"`
	Time        int64  `json:"Time"`
	InfoHash    string `json:"InfoHash"`
	DownloadURL string `json:"DownloadURL"`
}

func (t torrent) release(config engine.Config, endpoint string) engine.Release {
	protocol := "https:"
	if strings.HasPrefix(endpoint, "http:") {
		protocol = "http:"
	}
	downloadURL := protocolPrefix.ReplaceAllString(strings.TrimSpace(t.DownloadURL), protocol)
	infoURL := fmt.Sprintf("%s//broadcasthe.net/torrents.php?id=%d&torrentid=%d", protocol, t.GroupID, t.TorrentID)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.ReplaceAll(strings.TrimSpace(t.ReleaseName), `\`, ""),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            fmt.Sprintf("BTN-%d", t.TorrentID),
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     common.UnixTime(t.Time),
	}
}
