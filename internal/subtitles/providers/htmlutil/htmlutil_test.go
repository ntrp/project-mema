package htmlutil

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

type stubService struct {
	response func(*http.Request, bool) *http.Response
}

func (s stubService) DoProviderRequest(req *http.Request, _ string, download bool) (*http.Response, error) {
	return s.response(req, download), nil
}

func TestURLAndResponseHelpers(t *testing.T) {
	if got := BaseURL(providercore.Config{BaseURL: "https://custom.test"}, "https://fallback"); got != "https://custom.test" {
		t.Fatal(got)
	}
	if got := Resolve("https://example.test/path/page", "../file.srt"); got != "https://example.test/file.srt" {
		t.Fatal(got)
	}
	got, err := WithQuery("https://example.test/base", "/search", map[string]string{"q": "title", "empty": ""})
	if err != nil || got != "https://example.test/base/search?q=title" {
		t.Fatalf("%s %v", got, err)
	}
	body, err := ReadResponse(response(200, "ok"), 10, "test")
	if err != nil || string(body) != "ok" {
		t.Fatalf("%q %v", body, err)
	}
	if _, err = ReadResponse(response(429, ""), 10, "test"); err == nil {
		t.Fatal("expected rate error")
	}
	if _, err = ReadResponse(response(200, "too long"), 2, "test"); err == nil {
		t.Fatal("expected size error")
	}
}

func TestRequestTestAndDownload(t *testing.T) {
	svc := stubService{response: func(req *http.Request, download bool) *http.Response {
		if download {
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/zip"}}, Body: io.NopCloser(bytes.NewReader(zipBody()))}
		}
		return response(200, "ok")
	}}
	if err := Test(context.Background(), svc, providercore.Config{BaseURL: "https://example.test"}, "fixture", ""); err != nil {
		t.Fatal(err)
	}
	dl, err := Download(context.Background(), svc, "fixture", "https://example.test/sub.zip", "release", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dl.Content), "subtitle") {
		t.Fatalf("%q", dl.Content)
	}
	if !LooksArchive("application/gzip") || Filename("https://example.test/path/file.zip", "") != "file.zip" || Filename(":bad", "fallback") != "fallback" {
		t.Fatal("helper mismatch")
	}
}
func response(status int, body string) *http.Response {
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
