package htmltable

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

type Options struct {
	Name           string
	DefaultBaseURL string
	SearchPath     string
	QueryParam     string
	CategoryParam  string
	CategoryJoin   string
	ExtraParams    map[string]string
	LoginPath      string
	UsernameParam  string
	PasswordParam  string
	ExtraLogin     map[string]string
}

type Engine struct {
	client  engine.HTTPDoer
	options Options
}

var intRegex = regexp.MustCompile(`\d+`)

func New(options Options, clients ...engine.HTTPDoer) *Engine {
	var client engine.HTTPDoer
	if len(clients) > 0 {
		client = clients[0]
	}
	return &Engine{client: client, options: options}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	_, err := e.Search(ctx, config, "test", "")
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	return engine.SuccessResult(e.options.Name + " indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := e.searchURL(config, query)
	if err != nil {
		return nil, err
	}
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	if cookie, err := e.cookie(ctx, config); err != nil {
		return nil, err
	} else if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, e.options.DefaultBaseURL)
	releases := []engine.Release{}
	doc.Find("tr, div.browsePoster, div.torrent, li").Each(func(_ int, row *goquery.Selection) {
		release := releaseFromNode(config, baseURL, row)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	})
	return releases, nil
}

func (e *Engine) searchURL(config engine.Config, query string) (string, error) {
	path := e.options.SearchPath
	if path == "" {
		path = "/browse.php"
	}
	values := map[string]string{}
	for key, value := range e.options.ExtraParams {
		values[key] = value
	}
	param := e.options.QueryParam
	if param == "" {
		param = "search"
	}
	if query = strings.TrimSpace(query); query != "" && param != "-" {
		values[param] = strings.ReplaceAll(query, ".", " ")
	}
	if e.options.CategoryParam != "" && len(config.Categories) > 0 {
		parts := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			parts = append(parts, strconv.FormatInt(int64(category), 10))
		}
		join := e.options.CategoryJoin
		if join == "" {
			join = ";"
		}
		values[e.options.CategoryParam] = strings.Join(parts, join)
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		if _, ok := values["freeleech"]; !ok {
			values["freeleech"] = "1"
		}
	}
	return common.URLWithQuery(common.BaseURL(config, e.options.DefaultBaseURL), path, values)
}

func (e *Engine) cookie(ctx context.Context, config engine.Config) (string, error) {
	if cookie := common.FieldString(config, "cookie"); cookie != "" {
		return cookie, nil
	}
	if e.options.LoginPath == "" {
		return "", nil
	}
	username := common.FieldString(config, "username")
	password := common.FieldString(config, "password")
	if username == "" || password == "" {
		return "", nil
	}
	form := url.Values{}
	form.Set(defaultValue(e.options.UsernameParam, "username"), username)
	form.Set(defaultValue(e.options.PasswordParam, "password"), password)
	for key, value := range e.options.ExtraLogin {
		form.Set(key, value)
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, e.options.DefaultBaseURL), e.options.LoginPath, nil)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", engine.HTTPStatusError(resp)
	}
	cookies := make([]string, 0, len(resp.Cookies()))
	for _, cookie := range resp.Cookies() {
		cookies = append(cookies, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookies, "; "), nil
}

func releaseFromNode(config engine.Config, baseURL string, row *goquery.Selection) engine.Release {
	title, infoURL, downloadURL := links(baseURL, row)
	if downloadURL == "" {
		downloadURL = infoURL
	}
	seeders, leechers := seedCounts(row)
	release := engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            engine.FirstNonEmpty(infoURL, downloadURL, title),
		SizeBytes:       common.ParseSizeBytes(row.Text()),
		PublishedAt:     published(row),
	}
	if seeders >= 0 {
		release.Seeders = common.Int32Ptr(seeders)
	}
	if seeders >= 0 && leechers >= 0 {
		release.Peers = common.Int32Ptr(seeders + leechers)
	}
	return release
}

func links(baseURL string, row *goquery.Selection) (string, string, string) {
	title, infoURL, downloadURL := "", "", ""
	row.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		resolved := common.ResolveURL(baseURL, href)
		text := clean(link.Text())
		if text == "" {
			text, _ = link.Attr("title")
			text = clean(text)
		}
		lower := strings.ToLower(href)
		isDownload := strings.Contains(lower, "download") || strings.Contains(lower, ".torrent") || strings.HasPrefix(lower, "magnet:")
		if isDownload && downloadURL == "" {
			downloadURL = resolved
		}
		if !isDownload && infoURL == "" {
			infoURL = resolved
		}
		if !isDownload && title == "" && len(text) > 2 {
			title = text
		}
	})
	return title, infoURL, downloadURL
}

func seedCounts(row *goquery.Selection) (int, int) {
	seeders, leechers := -1, -1
	numbers := []int{}
	row.Find("td, span, a").Each(func(_ int, cell *goquery.Selection) {
		text := strings.ReplaceAll(clean(cell.Text()), ",", "")
		value, ok := firstInt(text)
		if ok {
			numbers = append(numbers, value)
		}
		class, _ := cell.Attr("class")
		label, _ := cell.Attr("title")
		name := strings.ToLower(class + " " + label)
		if strings.Contains(name, "seed") && !strings.Contains(name, "leech") && ok {
			seeders = value
		}
		if strings.Contains(name, "leech") && ok {
			leechers = value
		}
	})
	if seeders < 0 && len(numbers) >= 2 {
		seeders = numbers[len(numbers)-2]
		leechers = numbers[len(numbers)-1]
	}
	return seeders, leechers
}

func published(row *goquery.Selection) *time.Time {
	var value string
	row.Find("time").EachWithBreak(func(_ int, item *goquery.Selection) bool {
		value, _ = item.Attr("datetime")
		if value == "" {
			value, _ = item.Attr("title")
		}
		if value == "" {
			value = item.Text()
		}
		return value == ""
	})
	if value != "" {
		return common.ParseFlexibleTime(value)
	}
	row.Find("td").EachWithBreak(func(_ int, cell *goquery.Selection) bool {
		text := clean(cell.Text())
		if strings.ContainsAny(text, ":-/") && common.ParseFlexibleTime(text) != nil {
			value = text
			return false
		}
		return true
	})
	return common.ParseFlexibleTime(value)
}

func clean(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func firstInt(value string) (int, bool) {
	match := intRegex.FindString(value)
	if match == "" {
		return 0, false
	}
	parsed, err := strconv.Atoi(match)
	return parsed, err == nil
}

func defaultValue(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
