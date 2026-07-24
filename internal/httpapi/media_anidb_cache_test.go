package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func TestAniDBSearchUsesExistingMetadataCacheAndHistory(t *testing.T) {
	store := testSettingsStore(t)
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(`<anime id="123"><titles><title xml:lang="x-jat" type="main">Cached Anime</title></titles></anime>`))
	}))
	defer upstream.Close()
	clientName, clientVersion := "testclient", "1"
	provider, err := store.CreateMetadataProvider(context.Background(), storage.MetadataProviderInput{
		Name: "AniDB cache test", Type: "anidb", BaseURL: upstream.URL,
		APIKey: &clientName, PIN: &clientVersion, Enabled: true, Priority: 120,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.DeleteMetadataProvider(context.Background(), provider.ID) })
	server := &Server{settings: store, metadata: metadata.NewService(upstream.Client(), nil), now: time.Now}
	request := metadata.SearchRequest{Query: "Cached Anime", MediaType: "serie"}
	for range 2 {
		results, searchErr := server.searchMetadataProvider(context.Background(), provider, request)
		if searchErr != nil || len(results) != 1 {
			t.Fatalf("results=%#v err=%v", results, searchErr)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("upstream calls = %d, want one cache miss", calls.Load())
	}
}
