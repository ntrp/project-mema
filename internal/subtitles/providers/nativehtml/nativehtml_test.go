package nativehtml

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
)

type fakeService struct {
	handler func(*http.Request, string, bool) (*http.Response, error)
}

func (f fakeService) DoProviderRequest(req *http.Request, key string, download bool) (*http.Response, error) {
	return f.handler(req, key, download)
}

func resp(status int, body []byte) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(body))}
}

func TestAdapterSearchAuthAndDownloadArchive(t *testing.T) {
	spec := Spec{Key: "demo", BaseURL: "https://example.test", NeedsCookie: true, ArchiveDownload: true, SearchPath: func(providercore.SearchRequest) string { return "/find" }, Query: func(q url.Values, sr providercore.SearchRequest) { q.Set("q", QueryText(sr)) }, Candidate: func(sel *goquery.Selection, pageURL, lang string) (providercore.Candidate, bool) {
		return LinkCandidate(sel, pageURL, lang, "a[href]")
	}}
	adapter := New(spec)
	year := int32(2020)
	service := fakeService{handler: func(req *http.Request, key string, download bool) (*http.Response, error) {
		if key != "demo" || req.Header.Get("Cookie") != "sid=1" {
			t.Fatalf("bad key/cookie %s %q", key, req.Header.Get("Cookie"))
		}
		if !download {
			if req.URL.Path != "/find" || !strings.Contains(req.URL.Query().Get("q"), "Movie") {
				t.Fatalf("bad search url %s", req.URL.String())
			}
			return resp(200, []byte(`<table><tr data-language="en"><td class="title">Movie.Release</td><td><a href="/download/sub.zip">download</a></td></tr></table>`)), nil
		}
		return resp(200, zipBytes(t, "sub.srt", []byte("subtitle"))), nil
	}}
	candidates, err := adapter.Search(context.Background(), service, providercore.Config{SecretSettings: map[string]string{"cookies": "sid=1"}}, providercore.SearchRequest{Title: "Movie", Year: &year, LanguageID: "en"})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search()=%#v err=%v", candidates, err)
	}
	download, err := adapter.Download(context.Background(), service, providercore.Config{SecretSettings: map[string]string{"cookies": "sid=1"}}, candidates[0])
	if err != nil || string(download.Content) != "subtitle" {
		t.Fatalf("Download()=%q err=%v", download.Content, err)
	}
}

func TestAdapterRequiresCookies(t *testing.T) {
	_, err := New(Spec{Key: "private", NeedsCookie: true}).Search(context.Background(), fakeService{}, providercore.Config{}, providercore.SearchRequest{})
	if err == nil || !strings.Contains(err.Error(), providercore.ErrPrivateMembershipRequired.Error()) {
		t.Fatalf("expected private membership error, got %v", err)
	}
}

func zipBytes(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var b bytes.Buffer
	w := zip.NewWriter(&b)
	f, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
