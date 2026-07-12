package gestdown

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

func TestGestdownSearchAndDownload(t *testing.T) {
	season, episode := int32(1), int32(4)
	svc := stubService(func(r *http.Request, provider string, download bool) (*http.Response, error) {
		if provider != key {
			t.Fatalf("provider = %s", provider)
		}
		if r.Header.Get("User-Agent") != "Bazarr" {
			t.Fatalf("missing Bazarr user agent")
		}
		switch r.URL.Path {
		case "/shows/external/tvdb/321":
			return jsonResp(`{"shows":[{"id":55}]}`), nil
		case "/subtitles/get/55/1/4/en":
			return jsonResp(`{"matchingSubtitles":[{"subtitleId":9,"downloadUri":"/download/9","version":"WEB-GRP, HDTV-OTHER","qualities":["1080p"],"completed":true,"hearingImpaired":false},{"subtitleId":10,"downloadUri":"/download/10","version":"draft","completed":false,"hearingImpaired":false}]}`), nil
		case "/download/9":
			if !download {
				t.Fatalf("download flag not set")
			}
			return textResp("1\r\ngestdown\r\n"), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		return nil, nil
	})
	cands, err := (adapter{}).Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", LanguageID: "eng", SeasonNumber: &season, EpisodeNumber: &episode, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"tvdb": "321"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 || cands[0].FileID != 9 || cands[0].ReleaseName != "WEB-GRP\nHDTV-OTHER" || cands[0].SourceURL != "https://api.gestdown.info/download/9" {
		t.Fatalf("unexpected candidates: %#v", cands)
	}
	dl, err := (adapter{}).Download(context.Background(), svc, providercore.Config{}, cands[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(dl.Content) != "1\ngestdown\n" {
		t.Fatalf("content = %q", dl.Content)
	}
}

func TestGestdownMissingTVDB(t *testing.T) {
	season, episode := int32(1), int32(1)
	_, err := (adapter{}).Search(context.Background(), nil, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", SeasonNumber: &season, EpisodeNumber: &episode})
	if err == nil {
		t.Fatal("expected missing tvdb error")
	}
}

func jsonResp(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func textResp(body string) *http.Response {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}
}
