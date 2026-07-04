package myanonamouse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://www.myanonamouse.net/"

var sanitizeQuery = regexp.MustCompile(`[^\w]+`)

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
		return engine.FailedResult("Invalid MyAnonamouse request", "error", err.Error())
	}
	return engine.SuccessResult("MyAnonamouse indexer reachable")
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
	if cookie := mamCookie(config); cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("mam_id expired or invalid")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	var decoded response
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	if decoded.Error != "" && !strings.HasPrefix(strings.ToLower(decoded.Error), "nothing returned") {
		return nil, fmt.Errorf("myanonamouse api error: %s", decoded.Error)
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(decoded.Data))
	for _, item := range decoded.Data {
		release := item.release(config, baseURL)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	term := strings.TrimSpace(sanitizeQuery.ReplaceAllString(query, " "))
	values := map[string]string{
		"tor[text]":             term,
		"tor[searchType]":       searchType(config),
		"tor[srchIn][title]":    "true",
		"tor[srchIn][author]":   "true",
		"tor[srchIn][narrator]": "true",
		"tor[searchIn]":         "torrents",
		"tor[sortType]":         "default",
		"tor[perpage]":          "100",
		"tor[startNumber]":      "0",
		"thumbnails":            "1",
		"description":           "1",
	}
	if common.FieldBool(config, "searchInDescription") {
		values["tor[srchIn][description]"] = "true"
	}
	if common.FieldBool(config, "searchInSeries") {
		values["tor[srchIn][series]"] = "true"
	}
	if common.FieldBool(config, "searchInFilenames") {
		values["tor[srchIn][filenames]"] = "true"
	}
	if len(config.Categories) == 0 {
		values["tor[cat][]"] = "0"
	} else {
		for index, category := range config.Categories {
			values[fmt.Sprintf("tor[cat][%d]", index)] = strconv.FormatInt(int64(category), 10)
		}
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/tor/js/loadSearchJSONbasic.php", values)
}

func searchType(config engine.Config) string {
	switch int(common.FieldFloat(config, "searchType")) {
	case 1:
		return "active"
	case 2:
		return "fl"
	case 3:
		return "fl-VIP"
	case 4:
		return "VIP"
	case 5:
		return "nVIP"
	default:
		return "all"
	}
}

func mamCookie(config engine.Config) string {
	if cookie := common.FieldString(config, "cookie"); cookie != "" {
		return cookie
	}
	if mamID := common.FieldString(config, "mamId", "mam_id"); mamID != "" {
		return "mam_id=" + url.QueryEscape(mamID)
	}
	return ""
}

type response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    []item `json:"data"`
}

type item struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	AuthorInfo        string `json:"author_info"`
	LanguageCode      string `json:"lang_code"`
	Filetype          string `json:"filetype"`
	Vip               bool   `json:"vip"`
	Free              bool   `json:"free"`
	PersonalFreeLeech bool   `json:"personal_freeleech"`
	Category          string `json:"category"`
	Added             string `json:"added"`
	Grabs             int    `json:"times_completed"`
	Seeders           int    `json:"seeders"`
	Leechers          int    `json:"leechers"`
	NumFiles          int32  `json:"numfiles"`
	Size              string `json:"size"`
}

func (i item) release(config engine.Config, baseURL string) engine.Release {
	id := strconv.Itoa(i.ID)
	infoURL := common.ResolveURL(baseURL, "t/"+id)
	downloadURL, _ := common.URLWithQuery(baseURL, "/tor/download.php", map[string]string{"tid": id})
	if common.FieldBool(config, "useFreeleechWedge") && !i.Free && !i.PersonalFreeLeech {
		downloadURL += "&fl=1"
	}
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           i.title(),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       common.ParseSizeBytes(i.Size),
		Seeders:         common.Int32Ptr(i.Seeders),
		Peers:           common.Int32Ptr(i.Seeders + i.Leechers),
		PublishedAt:     common.ParseFlexibleTime(i.Added),
	}
}

func (i item) title() string {
	title := strings.TrimSpace(i.Title)
	if i.AuthorInfo != "" {
		var authors map[string]string
		if json.Unmarshal([]byte(i.AuthorInfo), &authors) == nil && len(authors) > 0 {
			values := []string{}
			for _, author := range authors {
				if strings.TrimSpace(author) != "" {
					values = append(values, strings.TrimSpace(author))
				}
			}
			if len(values) > 0 {
				title += " by " + strings.Join(values, ", ")
			}
		}
	}
	flags := []string{}
	if i.LanguageCode != "" {
		flags = append(flags, i.LanguageCode)
	}
	if i.Filetype != "" {
		flags = append(flags, strings.ToUpper(i.Filetype))
	}
	if len(flags) > 0 {
		title += " [" + strings.Join(flags, " / ") + "]"
	}
	if i.Vip {
		title += " [VIP]"
	}
	return title
}
