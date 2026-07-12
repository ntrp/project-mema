package sitehtml

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type stubService struct {
	response func(*http.Request, bool) *http.Response
}

func (s stubService) DoProviderRequest(req *http.Request, _ string, download bool) (*http.Response, error) {
	return s.response(req, download), nil
}

func TestHTMLAndSelectionHelpers(t *testing.T) {
	svc := stubService{response: func(_ *http.Request, _ bool) *http.Response {
		return resp(200, `<main data-lang="eng"><h1> A  title </h1></main>`)
	}}
	req, _ := http.NewRequest(http.MethodGet, "https://example.test", nil)
	if err := Test(req, svc, "fixture"); err != nil {
		t.Fatal(err)
	}
	doc, body, err := DoHTML(svc, req, "fixture")
	if err != nil || len(body) == 0 {
		t.Fatal(err)
	}
	main := doc.Find("main")
	if Text(main, "h1") != "A title" || Attr(main, "data-lang") != "eng" {
		t.Fatal("selection helpers")
	}
	if !Supports("movie", "movie", "serie") || Supports("audio", "movie") || Unsupported("x", "audio") == nil {
		t.Fatal("support helpers")
	}
	if Resolve("https://example.test/a/page", "../sub") != "https://example.test/sub" {
		t.Fatal("resolve")
	}
}
func TestDownloadAndFilenameHelpers(t *testing.T) {
	svc := stubService{response: func(_ *http.Request, download bool) *http.Response {
		if download {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBody()))}
		}
		return resp(200, "ok")
	}}
	req, _ := http.NewRequest(http.MethodGet, "https://example.test/sub.zip", nil)
	dl, err := Download(req, svc, "fixture", true, "release")
	if err != nil {
		t.Fatal(err)
	}
	if string(dl.Content) != "subtitle" {
		t.Fatalf("%q", dl.Content)
	}
	if !LooksArchived("file.rar") || FilenameFor("https://example.test/a/file.zip", "") != "file.zip" || FilenameFor(":bad", "fallback") != "fallback" {
		t.Fatal("archive helpers")
	}
	bad := stubService{response: func(_ *http.Request, _ bool) *http.Response { return resp(429, "") }}
	if _, _, err = DoHTML(bad, req, "fixture"); err == nil {
		t.Fatal("expected rate error")
	}
}
func resp(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}
}
func zipBody() []byte {
	var b bytes.Buffer
	w := zip.NewWriter(&b)
	f, _ := w.Create("fixture.srt")
	_, _ = f.Write([]byte("subtitle"))
	_ = w.Close()
	return b.Bytes()
}

var _ providercore.Service = stubService{}
