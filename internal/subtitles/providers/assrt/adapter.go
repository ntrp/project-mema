package assrt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/security"
)

const (
	key        = "assrt"
	defaultURL = "https://api.assrt.net/v1"
	maxBody    = 10 << 20
)

type adapter struct{}

type quotaResponse struct {
	Status string `json:"status"`
	ErrMsg string `json:"errmsg"`
	User   struct {
		Quota int `json:"quota"`
	} `json:"user"`
}

type searchResponse struct {
	Status string `json:"status"`
	ErrMsg string `json:"errmsg"`
	Sub    struct {
		Subs []searchSub `json:"subs"`
	} `json:"sub"`
}

type searchSub struct {
	ID         int64             `json:"id"`
	VideoName  string            `json:"videoname"`
	NativeName any               `json:"native_name"`
	Lang       map[string]any    `json:"lang"`
	LangList   map[string]string `json:"langlist"`
}

type detailResponse struct {
	Status string `json:"status"`
	ErrMsg string `json:"errmsg"`
	Sub    struct {
		Subs []detailSub `json:"subs"`
	} `json:"sub"`
}

type detailSub struct {
	URL      string       `json:"url"`
	FileList []detailFile `json:"filelist"`
}

type detailFile struct {
	Name string `json:"f"`
	URL  string `json:"url"`
}

func init() { providers.Register(key, adapter{}) }

func (adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	var out quotaResponse
	if err := doJSON(ctx, service, config, http.MethodGet, "/user/quota", tokenQuery(config), false, &out); err != nil {
		return err
	}
	if out.User.Quota <= 0 {
		return fmt.Errorf("%w: invalid assrt quota", providercore.ErrProviderBrokenUpstream)
	}
	return nil
}

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	q := tokenQuery(config)
	q.Set("q", query(request))
	q.Set("is_file", "1")
	var out searchResponse
	if err := doJSON(ctx, service, config, http.MethodGet, "/sub/search", q, false, &out); err != nil {
		return nil, err
	}
	candidates := make([]providercore.Candidate, 0, len(out.Sub.Subs))
	for _, sub := range out.Sub.Subs {
		for _, lang := range languages(sub) {
			name := videoName(sub)
			if name == "" || name == "不知道" {
				name = request.Title
			}
			candidates = append(candidates, providercore.Candidate{ProviderName: key, LanguageID: lang, FileID: sub.ID, Format: "srt", ReleaseName: name, SourceRef: lang})
		}
	}
	return candidates, nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	if candidate.FileID == 0 {
		return providercore.Download{}, fmt.Errorf("%w: assrt subtitle id is required", providercore.ErrProviderPrerequisiteMissing)
	}
	q := tokenQuery(config)
	q.Set("id", strconv.FormatInt(candidate.FileID, 10))
	var out detailResponse
	if err := doJSON(ctx, service, config, http.MethodGet, "/sub/detail", q, false, &out); err != nil {
		return providercore.Download{}, err
	}
	link := chooseDownload(out, candidate)
	if link == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing assrt download URL", providercore.ErrProviderBrokenUpstream)
	}
	data, _, err := do(ctx, service, config, http.MethodGet, link, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")), URL: link}, nil
}

func tokenQuery(config providercore.Config) url.Values {
	token, ok := providercore.NewConfig(config).RequiredSecret("token")
	if !ok {
		return url.Values{}
	}
	return url.Values{"token": []string{token}}
}

func query(r providercore.SearchRequest) string {
	parts := []string{r.Title}
	if r.MediaType == "movie" && r.Year != nil {
		parts = append(parts, strconv.Itoa(int(*r.Year)))
	}
	if r.MediaType == "serie" {
		if r.SeasonNumber != nil && r.EpisodeNumber != nil {
			parts = append(parts, fmt.Sprintf("S%02dE%02d", *r.SeasonNumber, *r.EpisodeNumber))
		} else if r.EpisodeNumber != nil {
			parts = append(parts, fmt.Sprintf("E%02d", *r.EpisodeNumber))
		}
	}
	return strings.Join(nonEmpty(parts), " ")
}

func nonEmpty(values []string) []string {
	out := values[:0]
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			out = append(out, strings.TrimSpace(v))
		}
	}
	return out
}

func languages(sub searchSub) []string {
	raw := sub.LangList
	if raw == nil && sub.Lang != nil {
		if nested, ok := sub.Lang["langlist"].(map[string]any); ok {
			raw = map[string]string{}
			for k := range nested {
				raw[k] = k
			}
		}
	}
	re := regexp.MustCompile(`^lang(\w+)$`)
	var out []string
	for k := range raw {
		if m := re.FindStringSubmatch(k); len(m) == 2 {
			out = append(out, assrtLanguage(m[1]))
		}
	}
	return out
}

func assrtLanguage(code string) string {
	switch strings.ToLower(code) {
	case "eng":
		return "eng"
	case "chs", "cht", "chi", "zho":
		return "zho"
	case "jpn":
		return "jpn"
	default:
		return strings.ToLower(code)
	}
}

func videoName(sub searchSub) string {
	if sub.VideoName != "" && sub.VideoName != "不知道" {
		return sub.VideoName
	}
	switch v := sub.NativeName.(type) {
	case string:
		return v
	case []any:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return sub.VideoName
}

func chooseDownload(out detailResponse, candidate providercore.Candidate) string {
	if len(out.Sub.Subs) == 0 {
		return ""
	}
	sub := out.Sub.Subs[0]
	if len(sub.FileList) == 0 {
		return sub.URL
	}
	for _, f := range sub.FileList {
		if strings.Contains(strings.ToLower(f.Name), candidate.SourceRef) {
			return f.URL
		}
	}
	return sub.FileList[0].URL
}

func doJSON(ctx context.Context, service providercore.Service, config providercore.Config, method, endpoint string, q url.Values, download bool, dst any) error {
	data, _, err := do(ctx, service, config, method, endpoint, q, download)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}
	if msg := statusMessage(dst); msg != "" {
		return fmt.Errorf("%w: %s", providercore.ErrProviderBrokenUpstream, msg)
	}
	return nil
}

func statusMessage(v any) string {
	switch x := v.(type) {
	case *quotaResponse:
		return x.ErrMsg
	case *searchResponse:
		return x.ErrMsg
	case *detailResponse:
		return x.ErrMsg
	}
	return ""
}

func do(ctx context.Context, service providercore.Service, config providercore.Config, method, endpoint string, q url.Values, download bool) ([]byte, *http.Response, error) {
	if _, ok := providercore.NewConfig(config).RequiredSecret("token"); !ok && !download {
		return nil, nil, fmt.Errorf("%w: token is required", providercore.ErrProviderPrerequisiteMissing)
	}
	raw := endpoint
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		u, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultURL))
		u.Path = path.Join(u.Path, endpoint)
		u.RawQuery = q.Encode()
		raw = u.String()
	}
	req, err := http.NewRequestWithContext(ctx, method, raw, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "Sub-Zero/2")
	resp, err := service.DoProviderRequest(req, key, download)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, resp, err
	}
	if len(data) > maxBody {
		return nil, resp, fmt.Errorf("response size limit exceeded")
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, resp, fmt.Errorf("%w: http status %d", providercore.ErrProviderPrerequisiteMissing, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	if download {
		if err := security.ValidateProviderURL(key, raw, true); err != nil {
			return nil, resp, err
		}
	}
	return data, resp, nil
}
