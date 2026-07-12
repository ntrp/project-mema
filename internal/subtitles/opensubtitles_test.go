package subtitles

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestOpenSubtitlesComAliasSearchAndDownload(t *testing.T) {
	apiKey := "subtitle-key"
	seenPaths := []string{}
	service := NewService(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		seenPaths = append(seenPaths, req.URL.Host+req.URL.Path)
		if req.Header.Get("Api-Key") != apiKey {
			t.Fatalf("missing api key header: %#v", req.Header)
		}
		switch req.URL.Host + req.URL.Path {
		case "api.opensubtitles.com/api/v1/subtitles":
			return stringResponse(200, `{"data":[{"attributes":{"language":"en","download_count":8,"url":"https://www.opensubtitles.com/subtitle/1","feature_details":{"title":"Scenario Movie"},"files":[{"file_id":42,"file_name":"Scenario.Movie.srt"}]}}]}`), nil
		case "api.opensubtitles.com/api/v1/download":
			return stringResponse(200, `{"link":"https://dl.opensubtitles.com/download/subtitle.srt"}`), nil
		case "dl.opensubtitles.com/download/subtitle.srt":
			return stringResponse(200, "1\n00:00:00,000 --> 00:00:01,000\nScenario\n"), nil
		default:
			t.Fatalf("unexpected request %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})})
	config := Config{Type: "opensubtitles", Name: "OpenSubtitles", BaseURL: "https://api.opensubtitles.com", APIKey: &apiKey}

	candidates, err := service.Search(context.Background(), config, SearchRequest{Title: "Scenario Movie", LanguageID: "english"})
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0].FileID != 42 || candidates[0].LanguageID != "english" {
		t.Fatalf("candidates = %#v", candidates)
	}
	download, err := service.Download(context.Background(), config, candidates[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(download.Content), "Scenario") || download.URL == "" {
		t.Fatalf("download = %#v", download)
	}
	if len(seenPaths) != 3 {
		t.Fatalf("seenPaths = %#v", seenPaths)
	}
}

func TestAddic7edRequiresCaptchaSessionBeforeRuntimeSearch(t *testing.T) {
	service := NewService(nil)
	_, err := service.Search(context.Background(), Config{Type: "addic7ed"}, SearchRequest{Title: "Scenario", LanguageID: "english"})
	if !errors.Is(err, providercore.ErrCaptchaRequired) {
		t.Fatalf("expected captcha prerequisite error, got %v", err)
	}
}

func TestOpenSubtitlesBlocksUnexpectedDownloadHost(t *testing.T) {
	apiKey := "subtitle-key"
	service := NewService(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host+req.URL.Path == "api.opensubtitles.com/api/v1/download" {
			return stringResponse(200, `{"link":"https://evil.example/subtitle.srt"}`), nil
		}
		t.Fatalf("unexpected request %s", req.URL.String())
		return nil, nil
	})})
	_, err := service.Download(context.Background(), Config{Type: "opensubtitlescom", BaseURL: "https://api.opensubtitles.com", APIKey: &apiKey}, Candidate{FileID: 7})
	if err == nil || !strings.Contains(err.Error(), "outside provider allowlist") {
		t.Fatalf("err = %v", err)
	}
}

func TestOpenSubtitlesRejectsOversizedSubtitleDownload(t *testing.T) {
	apiKey := "subtitle-key"
	service := NewService(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Host + req.URL.Path {
		case "api.opensubtitles.com/api/v1/download":
			return stringResponse(200, `{"link":"https://dl.opensubtitles.com/download/huge.srt"}`), nil
		case "dl.opensubtitles.com/download/huge.srt":
			return stringResponse(200, strings.Repeat("x", (10<<20)+1)), nil
		default:
			t.Fatalf("unexpected request %s", req.URL.String())
			return nil, nil
		}
	})})
	_, err := service.Download(context.Background(), Config{Type: "opensubtitlescom", BaseURL: "https://api.opensubtitles.com", APIKey: &apiKey}, Candidate{FileID: 7})
	if err == nil || !strings.Contains(err.Error(), "too large") {
		t.Fatalf("err = %v", err)
	}
}

func stringResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
