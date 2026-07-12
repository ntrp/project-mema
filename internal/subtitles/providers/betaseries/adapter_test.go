package betaseries

import (
	"archive/zip"
	"bytes"
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

func TestBetaSeriesSearchWithSeriesTVDBAndDownloadZip(t *testing.T) {
	season, episode := int32(2), int32(3)
	svc := stubService(func(r *http.Request, provider string, download bool) (*http.Response, error) {
		if provider != key {
			t.Fatalf("provider = %s", provider)
		}
		switch r.URL.Path {
		case "/shows/episodes":
			q := r.URL.Query()
			if q.Get("key") != "tok" || q.Get("thetvdb_id") != "100" || q.Get("season") != "2" || q.Get("episode") != "3" || q.Get("subtitles") != "1" || q.Get("v") != "3" {
				t.Fatalf("bad query: %s", r.URL.RawQuery)
			}
			return jsonResp(`{"errors":[],"episodes":[{"subtitles":[{"id":7,"language":"vo","file":"Show.S02E03-GRP","url":"https://api.betaseries.com/subs/7.zip","source":"addic7ed"},{"id":8,"language":"vf","file":"bad","url":"https://api.betaseries.com/subs/8.zip","source":"seriessub"}]}]}`), nil
		case "/subs/7.zip":
			if !download {
				t.Fatalf("download flag not set")
			}
			return zipResp(map[string]string{".hidden.srt": "bad", "Show.S02E03-GRP.srt": "1\r\nzip\r\n"}), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		return nil, nil
	})
	cfg := providercore.Config{SecretSettings: map[string]string{"token": "tok"}}
	cands, err := (adapter{}).Search(context.Background(), svc, cfg, providercore.SearchRequest{MediaType: "serie", SeasonNumber: &season, EpisodeNumber: &episode, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"tvdb": "100"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 || cands[0].LanguageID != "eng" || cands[0].SourceRef != "addic7ed" {
		t.Fatalf("unexpected candidates: %#v", cands)
	}
	dl, err := (adapter{}).Download(context.Background(), svc, cfg, cands[0])
	if err != nil {
		t.Fatal(err)
	}
	if got := string(dl.Content); got != "1\nzip\n" {
		t.Fatalf("content = %q", got)
	}
}

func TestBetaSeriesEpisodeDisplayAndNoSeriesFound(t *testing.T) {
	svc := stubService(func(r *http.Request, _ string, _ bool) (*http.Response, error) {
		if r.URL.Path != "/episodes/display" || r.URL.Query().Get("thetvdb_id") != "200" {
			t.Fatalf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(`{"errors":[{"code":4001}]}`))}, nil
	})
	cands, err := (adapter{}).Search(context.Background(), svc, providercore.Config{SecretSettings: map[string]string{"token": "tok"}}, providercore.SearchRequest{MediaType: "serie", MediaContext: providercore.MediaContext{EpisodeExternalIDs: map[string]string{"tvdb": "200"}}})
	if err != nil || len(cands) != 0 {
		t.Fatalf("cands=%#v err=%v", cands, err)
	}
}

func jsonResp(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func zipResp(files map[string]string) *http.Response {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for name, content := range files {
		w, _ := zw.Create(name)
		_, _ = w.Write([]byte(content))
	}
	_ = zw.Close()
	return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(buf.Bytes()))}
}
