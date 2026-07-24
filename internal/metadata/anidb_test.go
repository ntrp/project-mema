package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
)

const anidbFixture = `<anime id="123" restricted="false">
<type>TV Series</type><episodecount>2</episodecount><startdate>2024-01-02</startdate><enddate>2024-03-04</enddate>
<titles><title xml:lang="x-jat" type="main">Example Anime</title><title xml:lang="en" type="official">Example English</title><title xml:lang="ja" type="official">例</title></titles>
<description>Example overview</description><picture>456.jpg</picture><episodes>
<episode id="1001"><epno type="1">1</epno><airdate>2024-01-02</airdate><title xml:lang="en">Regular</title></episode>
<episode id="1002"><epno type="2">1</epno><airdate>2024-02-02</airdate><title xml:lang="en">Special</title></episode>
<episode id="1003"><epno type="3">1</epno><title xml:lang="en">Credits</title></episode>
</episodes></anime>`

func TestAniDBSearchAndDetailsUseRegisteredClientAndMapFacts(t *testing.T) {
	var requests []url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Query())
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(anidbFixture))
	}))
	defer server.Close()

	service := NewService(server.Client(), nil)
	config := anidbTestConfig(server.URL)
	results, err := service.Search(context.Background(), config, SearchRequest{Query: "Example Anime", MediaType: "serie"})
	if err != nil || len(results) != 1 {
		t.Fatalf("Search() results=%#v err=%v", results, err)
	}
	if results[0].ExternalID != "123" || results[0].Title != "Example Anime" || results[0].PosterPath == nil {
		t.Fatalf("search result = %#v", results[0])
	}
	if requests[0].Get("aname") != "Example Anime" || requests[0].Get("client") != "testclient" || requests[0].Get("clientver") != "2" || requests[0].Get("protover") != "1" {
		t.Fatalf("search query = %#v", requests[0])
	}

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
	if details.EpisodeFacts[1].OriginalType != "2" || details.EpisodeFacts[1].AbsoluteNumber != nil || len(details.Seasons) != 0 {
		t.Fatalf("special episode was assigned standard numbering: %#v", details.EpisodeFacts[1])
	}
	if requests[1].Get("aid") != "123" || requests[1].Get("aname") != "" {
		t.Fatalf("details query = %#v", requests[1])
	}
}

func TestAniDBClearErrorsAndBounded429Retry(t *testing.T) {
	t.Run("identity required", func(t *testing.T) {
		service := NewService(http.DefaultClient, nil)
		_, err := service.Search(context.Background(), Config{Type: "anidb"}, SearchRequest{Query: "x", MediaType: "serie"})
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

	t.Run("429 retries once", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if calls.Add(1) == 1 {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			_, _ = w.Write([]byte(anidbFixture))
		}))
		defer server.Close()
		_, err := NewService(server.Client(), nil).Details(context.Background(), anidbTestConfig(server.URL), DetailsRequest{ExternalID: "123"})
		if err != nil || calls.Load() != 2 {
			t.Fatalf("calls=%d err=%v", calls.Load(), err)
		}
	})
}

func TestAniDBThrottleIsContextBounded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(anidbFixture)) }))
	defer server.Close()
	service := NewService(server.Client(), nil)
	config := anidbTestConfig(server.URL)
	if _, err := service.Details(context.Background(), config, DetailsRequest{ExternalID: "123"}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := service.Details(ctx, config, DetailsRequest{ExternalID: "123"}); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v", err)
	}
}

func anidbTestConfig(baseURL string) Config {
	name, version := "testclient", "2"
	return Config{ID: uuid.New(), Type: "anidb", BaseURL: baseURL, ClientName: &name, ClientVersion: &version}
}
