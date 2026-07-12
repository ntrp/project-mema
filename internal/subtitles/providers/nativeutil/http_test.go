package nativeutil

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type service struct{ seen *http.Request }

func (s *service) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.seen = req
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Disposition": {`attachment; filename="movie.zip"`}}, Body: io.NopCloser(strings.NewReader("<a href='/x'>x</a>"))}, nil
}

func TestDoBuildsFormRequestsAndHelpers(t *testing.T) {
	svc := &service{}
	data, resp, err := Do(context.Background(), svc, providercore.Config{BaseURL: "https://override.test/root"}, RequestSpec{Provider: "p", BaseURL: "https://base.test", Method: http.MethodPost, Path: "/search", Form: mapValues("q", "movie"), Headers: map[string]string{"X-Test": "1"}})
	if err != nil || string(data) == "" || resp == nil {
		t.Fatalf("Do = %q, %v, %v", data, resp, err)
	}
	if svc.seen.URL.String() != "https://override.test/search" || svc.seen.Header.Get("Content-Type") != "application/x-www-form-urlencoded" || svc.seen.Header.Get("X-Test") != "1" {
		t.Fatalf("unexpected request: %s %#v", svc.seen.URL, svc.seen.Header)
	}
	if Absolute(providercore.Config{BaseURL: "https://override.test/root"}, "https://base.test", "file.srt") != "https://override.test/root/file.srt" {
		t.Fatal("relative URL was not resolved against configured base path")
	}
	if Lang("", "Español") != "spa" || Format("/a/b/movie.zip") != "zip" || DownloadName("/fallback.srt", resp) != "movie.zip" {
		t.Fatal("helper normalization failed")
	}
	doc, err := Document(data)
	if err != nil || Attr(doc.Selection, "a", "href") != "/x" || FirstText(doc.Selection, "a") != "x" {
		t.Fatalf("document helpers failed: %v", err)
	}
}

func TestGetFormAndDownloadSubtitle(t *testing.T) {
	svc := roundTrip(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/search" && !strings.Contains(req.URL.RawQuery, "q=movie") {
			t.Fatalf("missing query form: %s", req.URL.String())
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Disposition": {`attachment; filename="sub.zip"`}}, Body: io.NopCloser(bytes.NewReader(zipData()))}, nil
	})
	_, _, err := Do(context.Background(), svc, providercore.Config{}, RequestSpec{Provider: "p", BaseURL: "https://base.test", Path: "/search", Form: mapValues("q", "movie")})
	if err != nil {
		t.Fatal(err)
	}
	dl, err := DownloadSubtitle(context.Background(), svc, providercore.Config{}, "p", "https://base.test", providercore.Candidate{SourceURL: "/download/sub.zip"})
	if err != nil || !bytes.Contains(dl.Content, []byte("line")) || dl.URL != "https://base.test/download/sub.zip" {
		t.Fatalf("DownloadSubtitle = %#v, %v", dl, err)
	}
}

func TestDoRejectsBadStatus(t *testing.T) {
	svc := roundTrip(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 503, Body: io.NopCloser(strings.NewReader("down"))}, nil
	})
	_, _, err := Do(context.Background(), svc, providercore.Config{}, RequestSpec{Provider: "p", BaseURL: "https://base.test", Path: "/"})
	if err == nil {
		t.Fatal("expected HTTP status error")
	}
}

type roundTrip func(*http.Request) (*http.Response, error)

func (r roundTrip) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	return r(req)
}

func mapValues(key, value string) url.Values { return url.Values{key: {value}} }

func zipData() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("sub.srt")
	_, _ = w.Write([]byte("line"))
	_ = zw.Close()
	return buf.Bytes()
}
