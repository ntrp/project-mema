package jobs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestSCNMedia014RSSSyncGrabsMatchingMovie(t *testing.T) {
	ctx, store := jobsTestStore(t)
	indexer := rssSyncIndexer(t, ctx, store, rssSyncFeed("Wanted.Accepted.2026.1080p.WEB-DL", "https://indexer.test/download/accepted"))
	downloadClient := rssSyncSAB(t, true)
	wanted := rssSyncMedia(t, ctx, store, "Wanted Accepted")
	rssSyncDownloadClient(t, ctx, store, downloadClient)

	if err := runRSSSyncWorker(ctx, nil, store, indexers.NewService(indexer.Client()), downloadclients.NewService(downloadClient.Client()), decisions.NewEngine(), nil); err != nil {
		t.Fatalf("rss sync: %v", err)
	}

	activities, err := store.ListDownloadActivity(ctx)
	if err != nil {
		t.Fatalf("list download activity: %v", err)
	}
	if len(activities) != 1 || activities[0].Status != "grabbed" || activities[0].MediaItemID != wanted.ID {
		t.Fatalf("download activity = %#v", activities)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 1 || snapshot.Releases[0].SearchKind != "rss" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func TestSCNMedia014RSSSyncIgnoresNonMatchingRelease(t *testing.T) {
	ctx, store := jobsTestStore(t)
	indexer := rssSyncIndexer(t, ctx, store, rssSyncFeed("Different.Movie.2026.1080p.WEB-DL", "https://indexer.test/download/rejected"))
	wanted := rssSyncMedia(t, ctx, store, "Wanted Rejected")
	rssSyncDownloadClient(t, ctx, store, rssSyncSAB(t, true))

	if err := runRSSSyncWorker(ctx, nil, store, indexers.NewService(indexer.Client()), downloadclients.NewService(indexer.Client()), decisions.NewEngine(), nil); err != nil {
		t.Fatalf("rss sync: %v", err)
	}

	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 0 || len(snapshot.Errors) != 0 {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func TestSCNMedia014RSSSyncSkipsHandledMedia(t *testing.T) {
	ctx, store := jobsTestStore(t)
	indexer := rssSyncIndexer(t, ctx, store, rssSyncFeed("Already.Downloading.2026.1080p.WEB-DL", "https://indexer.test/download/handled"))
	item := rssSyncMedia(t, ctx, store, "Already Downloading")
	rssSyncDownloadClient(t, ctx, store, rssSyncSAB(t, true))
	if _, err := store.CreateDownloadActivity(ctx, storage.DownloadActivityInput{
		MediaItemID:        item.ID,
		ReleaseTitle:       "Existing",
		IndexerName:        "Indexer",
		DownloadClientName: "Client",
		DownloadURL:        "https://indexer.test/download/existing",
		Status:             "downloading",
	}); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	if err := runRSSSyncWorker(ctx, nil, store, indexers.NewService(indexer.Client()), downloadclients.NewService(indexer.Client()), decisions.NewEngine(), nil); err != nil {
		t.Fatalf("rss sync: %v", err)
	}
	activities, err := store.ListDownloadActivity(ctx)
	if err != nil || len(activities) != 1 {
		t.Fatalf("download activity = %#v, %v", activities, err)
	}
}

func TestSCNMedia014RSSSyncExcludesNonRSSAndNonRSSProfileIndexers(t *testing.T) {
	ctx, store := jobsTestStore(t)
	disabledRSS := rssSyncIndexerWithInput(t, ctx, store, rssSyncFeed("Wanted.Excluded.2026.1080p.WEB-DL", "https://indexer.test/download/one"), storage.IndexerInput{
		SupportsRSS:    false,
		SupportsSearch: true,
		Enabled:        true,
		AppProfileID:   "default",
	})
	_ = rssSyncIndexerWithInput(t, ctx, store, rssSyncFeed("Wanted.Excluded.2026.1080p.WEB-DL", "https://indexer.test/download/two"), storage.IndexerInput{
		SupportsRSS:    true,
		SupportsSearch: true,
		Enabled:        true,
		AppProfileID:   "no-rss",
	})
	wanted := rssSyncMedia(t, ctx, store, "Wanted Excluded")
	rssSyncDownloadClient(t, ctx, store, rssSyncSAB(t, true))

	if err := runRSSSyncWorker(ctx, nil, store, indexers.NewService(disabledRSS.Client()), downloadclients.NewService(disabledRSS.Client()), decisions.NewEngine(), nil); err != nil {
		t.Fatalf("rss sync: %v", err)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 0 {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func TestSCNMedia014RSSSyncMarkerPreventsDuplicateProcessing(t *testing.T) {
	ctx, store := jobsTestStore(t)
	indexer := rssSyncIndexer(t, ctx, store, rssSyncFeed("Wanted.Marker.2026.1080p.WEB-DL", "https://indexer.test/download/marker"))
	wanted := rssSyncMedia(t, ctx, store, "Wanted Marker")
	downloadServer := rssSyncSAB(t, true)
	rssSyncDownloadClient(t, ctx, store, downloadServer)
	service := indexers.NewService(indexer.Client())
	clients := downloadclients.NewService(downloadServer.Client())

	if err := runRSSSyncWorker(ctx, nil, store, service, clients, decisions.NewEngine(), nil); err != nil {
		t.Fatalf("first rss sync: %v", err)
	}
	if err := runRSSSyncWorker(ctx, nil, store, service, clients, decisions.NewEngine(), nil); err != nil {
		t.Fatalf("second rss sync: %v", err)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 1 {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func jobsTestStore(t *testing.T) (context.Context, *storage.SettingsStore) {
	t.Helper()
	databaseURL := testdb.Create(t)
	ctx := context.Background()
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return ctx, storage.NewSettingsStore(pool)
}

func rssSyncIndexer(t *testing.T, ctx context.Context, store *storage.SettingsStore, feed string) *httptest.Server {
	t.Helper()
	return rssSyncIndexerWithInput(t, ctx, store, feed, storage.IndexerInput{
		SupportsRSS:    true,
		SupportsSearch: true,
		Enabled:        true,
		AppProfileID:   "default",
	})
}

func rssSyncIndexerWithInput(t *testing.T, ctx context.Context, store *storage.SettingsStore, feed string, input storage.IndexerInput) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "" {
			http.Error(w, "rss sync must not send q", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(feed))
	}))
	t.Cleanup(server.Close)
	input.DefinitionID = "generic-newznab"
	input.Name = "RSS Sync"
	input.Implementation = "Newznab"
	input.ImplementationName = "Newznab"
	input.Protocol = "usenet"
	input.BaseURL = server.URL
	input.APIKey = stringPtr("key")
	input.Categories = []int32{2000}
	input.Priority = 1
	if _, err := store.CreateIndexer(ctx, input); err != nil {
		t.Fatalf("create indexer: %v", err)
	}
	return server
}

func rssSyncSAB(t *testing.T, success bool) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("mode") != "addurl" {
			http.Error(w, "unexpected mode", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if success {
			_, _ = w.Write([]byte(`{"status":true,"nzo_ids":["nzo-wanted"]}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":false,"error":"download rejected"}`))
	}))
	t.Cleanup(server.Close)
	return server
}

func rssSyncMedia(t *testing.T, ctx context.Context, store *storage.SettingsStore, title string) storage.MediaItem {
	t.Helper()
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type:      "movie",
		Title:     title,
		Year:      int32Ptr(2026),
		Monitored: true,
	})
	if err != nil {
		t.Fatalf("create media item: %v", err)
	}
	return item
}

func rssSyncDownloadClient(t *testing.T, ctx context.Context, store *storage.SettingsStore, server *httptest.Server) {
	t.Helper()
	if _, err := store.CreateDownloadClient(ctx, storage.DownloadClientInput{
		Name:     "RSS SAB",
		Type:     "sabnzbd",
		Protocol: "usenet",
		BaseURL:  server.URL,
		APIKey:   stringPtr("key"),
		Enabled:  true,
		Priority: 1,
	}); err != nil {
		t.Fatalf("create download client: %v", err)
	}
}

func rssSyncFeed(title string, downloadURL string) string {
	return fmt.Sprintf(`<rss><channel><item>
		<title>%s</title>
		<link>%s</link>
		<guid>%s</guid>
		<pubDate>Fri, 03 Jul 2026 04:00:00 +0200</pubDate>
		<size>8589934592</size>
	</item></channel></rss>`, title, downloadURL, downloadURL)
}

func stringPtr(value string) *string {
	return &value
}
