package metadata

import (
	"compress/gzip"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAniDBSearchAndDetailsUseRegisteredClientAndMapFacts(t *testing.T) {
	animeFixture := readAniDBFixture(t, "anidb-anime.xml")
	var requests []url.Values
	var titleCalls, artworkCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/anime-titles.xml.gz":
			titleCalls++
			writeGzip(t, w, readAniDBFixture(t, "anidb-titles.xml"))
			return
		case "/images/456.jpg", "/images/789.jpg":
			artworkCalls++
			writeTestJPEG(t, w)
			return
		}
		requests = append(requests, r.URL.Query())
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write(animeFixture)
	}))
	defer server.Close()

	service := NewService(server.Client(), nil)
	service.ConfigureAniDBArtworkCache(t.TempDir())
	service.anidbImageBaseURL = server.URL + "/images"
	service.anidbInterval = time.Millisecond
	config := anidbTestConfig(server.URL)
	results, err := service.Search(context.Background(), config, SearchRequest{Query: "Example Anime", MediaType: "serie"})
	if err != nil || len(results) != 1 {
		t.Fatalf("Search() results=%#v err=%v", results, err)
	}
	if results[0].ExternalID != "123" || results[0].Title != "Example Anime" || results[0].PosterPath == nil || *results[0].PosterPath != "/api/media/artwork/anidb/456.jpg" {
		t.Fatalf("search result = %#v", results[0])
	}
	if titleCalls != 1 || requests[0].Get("aid") != "123" || requests[0].Get("aname") != "" {
		t.Fatalf("titleCalls=%d search query=%#v", titleCalls, requests[0])
	}
	assertAniDBIdentity(t, requests[0])

	config.ID = uuid.New()
	details, err := service.Details(context.Background(), config, DetailsRequest{MediaType: "serie", ExternalID: "123"})
	if err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	if len(details.Aliases) != 3 || len(details.EpisodeFacts) != 3 {
		t.Fatalf("aliases=%#v episodes=%#v", details.Aliases, details.EpisodeFacts)
	}
	if details.EpisodeFacts[0].AbsoluteNumber == nil || *details.EpisodeFacts[0].AbsoluteNumber != 1 {
		t.Fatalf("regular episode = %#v", details.EpisodeFacts[0])
	}
	if details.EpisodeFacts[1].OriginalType != "2" || details.EpisodeFacts[1].AbsoluteNumber != nil || len(details.Seasons) != 2 {
		t.Fatalf("AniDB season mapping = %#v facts=%#v", details.Seasons, details.EpisodeFacts[1])
	}
	if len(details.Cast) != 1 || details.Cast[0].ProfilePath == nil || *details.Cast[0].ProfilePath != "/api/media/artwork/anidb/789.jpg" {
		t.Fatalf("AniDB cast = %#v", details.Cast)
	}
	if requests[1].Get("aid") != "123" || requests[1].Get("aname") != "" {
		t.Fatalf("details query = %#v", requests[1])
	}
	if artworkCalls != 2 {
		t.Fatalf("artwork calls = %d", artworkCalls)
	}
}

func TestAniDBTitleDumpIsCachedAndContextBounded(t *testing.T) {
	var titleCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/anime-titles.xml.gz" {
			titleCalls.Add(1)
			writeGzip(t, w, readAniDBFixture(t, "anidb-titles.xml"))
			return
		}
		_, _ = w.Write(readAniDBFixture(t, "anidb-anime.xml"))
	}))
	defer server.Close()
	service := NewService(server.Client(), nil)
	service.anidbInterval = time.Millisecond
	for range 2 {
		if _, err := service.Search(context.Background(), anidbTestConfig(server.URL), SearchRequest{Query: "Example", MediaType: "serie"}); err != nil {
			t.Fatal(err)
		}
	}
	if titleCalls.Load() != 1 {
		t.Fatalf("title dump calls = %d", titleCalls.Load())
	}

	blocked := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer blocked.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	_, err := NewService(blocked.Client(), nil).Search(ctx, anidbTestConfig(blocked.URL), SearchRequest{Query: "Example", MediaType: "serie"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("bounded title download error = %v", err)
	}
}

func TestAniDBTitleDumpSizeLimits(t *testing.T) {
	for _, test := range []struct {
		name, want string
		write      func(http.ResponseWriter)
	}{
		{name: "compressed", want: "compressed size limit", write: func(w http.ResponseWriter) {
			_, _ = io.CopyN(w, zeroReader{}, anidbTitleDumpMaxCompressed+1)
		}},
		{name: "uncompressed", want: "uncompressed size limit", write: func(w http.ResponseWriter) {
			archive := gzip.NewWriter(w)
			_, _ = io.CopyN(archive, zeroReader{}, anidbTitleDumpMaxUncompressed+1)
			_ = archive.Close()
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { test.write(w) }))
			defer server.Close()
			_, err := NewService(server.Client(), nil).Search(context.Background(), anidbTestConfig(server.URL), SearchRequest{Query: "Example", MediaType: "serie"})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

type zeroReader struct{}

func (zeroReader) Read(payload []byte) (int, error) {
	clear(payload)
	return len(payload), nil
}

func TestAniDBClearErrorsAndBanCooldown(t *testing.T) {
	t.Run("identity required", func(t *testing.T) {
		_, err := NewService(http.DefaultClient, nil).Search(context.Background(), Config{Type: "anidb"}, SearchRequest{Query: "x", MediaType: "serie"})
		if !errors.Is(err, ErrAniDBClientIdentityRequired) {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("invalid registered client name", func(t *testing.T) {
		version := "1"
		for _, name := range []string{"TestClient", "test client", "test_client", " testclient"} {
			config := Config{Type: "anidb", ClientName: &name, ClientVersion: &version}
			_, err := NewService(http.DefaultClient, nil).Search(context.Background(), config, SearchRequest{Query: "x", MediaType: "serie"})
			if !errors.Is(err, ErrAniDBClientNameInvalid) {
				t.Fatalf("name %q: error = %v", name, err)
			}
		}
	})
	t.Run("malformed XML", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("<anime")) }))
		defer server.Close()
		_, err := NewService(server.Client(), nil).Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "123"})
		if err == nil || !strings.Contains(err.Error(), "decode AniDB XML") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("upstream failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
		defer server.Close()
		_, err := NewService(server.Client(), nil).Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "123"})
		if status, ok := ProviderStatusCode(err); !ok || status != http.StatusBadGateway {
			t.Fatalf("error = %v status=%d", err, status)
		}
	})
	t.Run("long 429 is typed without retry", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer server.Close()
		_, err := NewService(server.Client(), nil).Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "123"})
		if !IsRateLimited(err) || calls.Load() != 1 {
			t.Fatalf("calls=%d err=%v", calls.Load(), err)
		}
	})
	t.Run("429 is not retried and starts a cooldown", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer server.Close()
		service := NewService(server.Client(), nil)
		_, firstErr := service.Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "123"})
		_, secondErr := service.Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "456"})
		if !IsRateLimited(firstErr) || !IsRateLimited(secondErr) || calls.Load() != 1 {
			t.Fatalf("calls=%d first=%v second=%v", calls.Load(), firstErr, secondErr)
		}
	})
}

func TestAniDBThrottleIsContextBounded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(readAniDBFixture(t, "anidb-anime.xml")) }))
	defer server.Close()
	service, config := NewService(server.Client(), nil), anidbTestConfig(server.URL)
	if _, err := service.Details(context.Background(), config, DetailsRequest{ExternalID: "123"}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	config.ID = uuid.New()
	if _, err := service.Details(ctx, config, DetailsRequest{ExternalID: "123"}); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v", err)
	}
}

func readAniDBFixture(t *testing.T, name string) []byte {
	t.Helper()
	contents, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func writeTestJPEG(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	w.Header().Set("Content-Type", "image/jpeg")
	picture := image.NewRGBA(image.Rect(0, 0, 1, 1))
	picture.Set(0, 0, color.White)
	if err := jpeg.Encode(w, picture, nil); err != nil {
		t.Fatal(err)
	}
}

func writeGzip(t *testing.T, w http.ResponseWriter, contents []byte) {
	t.Helper()
	archive := gzip.NewWriter(w)
	if _, err := archive.Write(contents); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
}

func assertAniDBIdentity(t *testing.T, query url.Values) {
	t.Helper()
	if query.Get("client") != "testclient" || query.Get("clientver") != "2" || query.Get("protover") != "1" || query.Get("request") != "anime" {
		t.Fatalf("AniDB query = %#v", query)
	}
}

func anidbTestConfig(baseURL string) Config {
	name, version := "testclient", "2"
	return Config{ID: uuid.New(), Type: "anidb", BaseURL: baseURL, ClientName: &name, ClientVersion: &version}
}
