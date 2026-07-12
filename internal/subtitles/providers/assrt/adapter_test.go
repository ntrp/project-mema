package assrt

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type stubService func(*http.Request, string, bool) (*http.Response, error)

func (f stubService) DoProviderRequest(r *http.Request, provider string, download bool) (*http.Response, error) {
	return f(r, provider, download)
}

func TestAssrtSearchAndDownload(t *testing.T) {
	season, episode := int32(7), int32(10)
	var sawSearch, sawDetail bool
	svc := stubService(func(r *http.Request, provider string, download bool) (*http.Response, error) {
		if provider != key {
			t.Fatalf("provider = %s", provider)
		}
		if r.Header.Get("User-Agent") != "Sub-Zero/2" {
			t.Fatalf("user agent not set")
		}
		switch r.URL.Path {
		case "/v1/sub/search":
			sawSearch = true
			if got := r.URL.Query().Get("q"); got != "Rick and Morty S07E10" {
				t.Fatalf("q = %q", got)
			}
			if r.URL.Query().Get("token") != "tok" || r.URL.Query().Get("is_file") != "1" {
				t.Fatalf("missing search params: %s", r.URL.RawQuery)
			}
			return jsonResp(`{"sub":{"subs":[{"id":42,"videoname":"不知道","native_name":["Rick.and.Morty.S07E10-GRP"],"lang":{"langlist":{"langeng":"English","langchs":"Chinese"}}}]}}`), nil
		case "/v1/sub/detail":
			sawDetail = true
			if r.URL.Query().Get("id") != "42" {
				t.Fatalf("detail id = %s", r.URL.Query().Get("id"))
			}
			return jsonResp(`{"sub":{"subs":[{"filelist":[{"f":"Rick.S07E10.chs.srt","url":"https://api.assrt.net/download/chs"},{"f":"Rick.S07E10.eng.srt","url":"https://api.assrt.net/download/eng"}]}]}}`), nil
		case "/download/eng":
			if !download {
				t.Fatalf("download flag not set")
			}
			return textResp("1\r\n00:00:01,000 --> 00:00:02,000\r\nhello\r\n"), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		return nil, nil
	})
	cfg := providercore.Config{SecretSettings: map[string]string{"token": "tok"}}
	cands, err := (adapter{}).Search(context.Background(), svc, cfg, providercore.SearchRequest{MediaType: "serie", Title: "Rick and Morty", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 2 || cands[0].FileID != 42 || cands[0].ReleaseName != "Rick.and.Morty.S07E10-GRP" {
		t.Fatalf("unexpected candidates: %#v", cands)
	}
	cand := cands[0]
	for _, c := range cands {
		if c.LanguageID == "eng" {
			cand = c
		}
	}
	dl, err := (adapter{}).Download(context.Background(), svc, cfg, cand)
	if err != nil {
		t.Fatal(err)
	}
	if !sawSearch || !sawDetail || !strings.Contains(string(dl.Content), "hello\n") || strings.Contains(string(dl.Content), "\r") {
		t.Fatalf("download/content mismatch: %q", dl.Content)
	}
}

func TestAssrtTestUsesQuota(t *testing.T) {
	svc := stubService(func(r *http.Request, _ string, _ bool) (*http.Response, error) {
		if r.URL.Path != "/v1/user/quota" || r.URL.Query().Get("token") != "tok" {
			t.Fatalf("unexpected quota request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		return jsonResp(`{"user":{"quota":60}}`), nil
	})
	if err := (adapter{}).Test(context.Background(), svc, providercore.Config{SecretSettings: map[string]string{"token": "tok"}}); err != nil {
		t.Fatal(err)
	}
}

func jsonResp(body string) *http.Response { return response(body, "application/json") }
func textResp(body string) *http.Response { return response(body, "text/plain") }
func response(body, ct string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{ct}}, Body: io.NopCloser(strings.NewReader(body))}
}
