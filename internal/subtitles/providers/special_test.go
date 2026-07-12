package providers

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ulikunitz/xz"

	"media-manager/internal/subtitles/providercore"
)

type fakeProviderService struct {
	handler func(*http.Request, string, bool) (*http.Response, error)
}

func (f fakeProviderService) DoProviderRequest(request *http.Request, providerType string, download bool) (*http.Response, error) {
	return f.handler(request, providerType, download)
}

func response(status int, body []byte) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(bytes.NewReader(body)), Header: http.Header{}}
}

func TestAnimeToshoTestRequestsFeed(t *testing.T) {
	called := false
	service := fakeProviderService{handler: func(request *http.Request, providerType string, download bool) (*http.Response, error) {
		called = providerType == animeToshoKey && !download && request.URL.Query().Get("show") == "torrent"
		return response(200, []byte(`{"files":[]}`)), nil
	}}
	if err := (animeToshoAdapter{}).Test(context.Background(), service, providercore.Config{}); err != nil || !called {
		t.Fatalf("Test() called=%v err=%v", called, err)
	}
}

func TestAnimeToshoSearchRequiresAniDBEpisodeID(t *testing.T) {
	_, err := (animeToshoAdapter{}).Search(context.Background(), fakeProviderService{}, providercore.Config{}, providercore.SearchRequest{LanguageID: "english"})
	if err == nil || !strings.Contains(err.Error(), providercore.ErrProviderPrerequisiteMissing.Error()) {
		t.Fatalf("expected prerequisite error, got %v", err)
	}
}

func TestAnimeToshoSearchAndDownloadXZ(t *testing.T) {
	service := fakeProviderService{handler: func(request *http.Request, providerType string, download bool) (*http.Response, error) {
		if providerType != animeToshoKey {
			t.Fatalf("provider = %s", providerType)
		}
		query := request.URL.Query()
		switch {
		case query.Get("eid") == "123":
			return response(200, []byte(`[{"id":77,"status":"complete","timestamp":2,"title":"Newest"},{"id":76,"status":"pending","timestamp":3}]`)), nil
		case query.Get("show") == "torrent" && query.Get("id") == "77":
			return response(200, []byte(`{"files":[{"filename":"Release.Name.mkv","attachments":[{"id":42,"type":"subtitle","info":{"lang":"eng","name":"English"}}]}]}`)), nil
		case download && strings.Contains(request.URL.Path, "/0000002a/42.xz"):
			return response(200, xzBytes(t, []byte("1\n00:00:01,000 --> 00:00:02,000\nHi\n"))), nil
		default:
			return response(404, nil), nil
		}
	}}
	candidates, err := (animeToshoAdapter{}).Search(context.Background(), service, providercore.Config{Name: "Anime"}, providercore.SearchRequest{LanguageID: "english", MediaContext: providercore.MediaContext{EpisodeExternalIDs: map[string]string{"anidb_episode_id": "123"}}})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search() candidates=%#v err=%v", candidates, err)
	}
	download, err := (animeToshoAdapter{}).Download(context.Background(), service, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download() err=%v", err)
	}
	if !bytes.Contains(download.Content, []byte("Hi")) {
		t.Fatalf("unexpected download content %q", download.Content)
	}
}

func TestNapiProjektTestRequestsKnownHash(t *testing.T) {
	called := false
	service := fakeProviderService{handler: func(request *http.Request, providerType string, download bool) (*http.Response, error) {
		called = providerType == napiprojektKey && !download && request.URL.Query().Get("f") == "00000000000000000000000000000000"
		return response(200, []byte("NPc0")), nil
	}}
	if err := (napiprojektAdapter{}).Test(context.Background(), service, providercore.Config{}); err != nil || !called {
		t.Fatalf("Test() called=%v err=%v", called, err)
	}
}

func TestNapiProjektHashSearchAndDownload(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "video-*.mkv")
	if err != nil {
		t.Fatal(err)
	}
	content := []byte("napiprojekt fixture")
	if _, err := file.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	wantHash := "f807ada25a5e1c463a1e4a45f3b4f541"
	var sawHash bool
	service := fakeProviderService{handler: func(request *http.Request, providerType string, download bool) (*http.Response, error) {
		if providerType != napiprojektKey || !download {
			t.Fatalf("provider=%s download=%v", providerType, download)
		}
		values, _ := url.ParseQuery(request.URL.RawQuery)
		if values.Get("f") == wantHash && values.Get("t") == napiprojektSubhash(wantHash) {
			sawHash = true
		}
		return response(200, []byte("subtitle body")), nil
	}}
	adapter := napiprojektAdapter{}
	candidates, err := adapter.Search(context.Background(), service, providercore.Config{Name: "Napi"}, providercore.SearchRequest{LanguageID: "polish", FilePath: file.Name()})
	if err != nil || len(candidates) != 1 || !sawHash {
		t.Fatalf("Search() candidates=%#v sawHash=%v err=%v", candidates, sawHash, err)
	}
	download, err := adapter.Download(context.Background(), service, providercore.Config{}, candidates[0])
	if err != nil || string(download.Content) != "subtitle body" {
		t.Fatalf("Download()=%q err=%v", download.Content, err)
	}
}

func TestBSPlayerBrokenSearchAndGzipDownload(t *testing.T) {
	adapter := bsplayerAdapter{}
	if err := adapter.Test(context.Background(), fakeProviderService{}, providercore.Config{}); err == nil || !strings.Contains(err.Error(), providercore.ErrProviderBrokenUpstream.Error()) {
		t.Fatalf("expected broken upstream test error, got %v", err)
	}
	if _, err := adapter.Search(context.Background(), fakeProviderService{}, providercore.Config{}, providercore.SearchRequest{}); err == nil || !strings.Contains(err.Error(), providercore.ErrProviderBrokenUpstream.Error()) {
		t.Fatalf("expected broken upstream error, got %v", err)
	}
	service := fakeProviderService{handler: func(request *http.Request, providerType string, download bool) (*http.Response, error) {
		if providerType != bsplayerKey || !download || request.Header.Get("User-Agent") == "" {
			t.Fatalf("bad request provider=%s download=%v ua=%q", providerType, download, request.Header.Get("User-Agent"))
		}
		return response(200, gzipBytes(t, []byte("bsplayer subtitle"))), nil
	}}
	download, err := adapter.Download(context.Background(), service, providercore.Config{}, providercore.Candidate{SourceURL: "https://s1.api.bsplayer-subtitles.com/sub.gz"})
	if err != nil || string(download.Content) != "bsplayer subtitle" {
		t.Fatalf("Download()=%q err=%v", download.Content, err)
	}
}

func gzipBytes(t *testing.T, content []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func xzBytes(t *testing.T, content []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer, err := xz.NewWriter(&buffer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
