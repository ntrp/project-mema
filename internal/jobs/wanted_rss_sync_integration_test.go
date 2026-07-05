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

func TestSCNMedia014WantedRSSSyncNoCandidateFlow(t *testing.T) {
	ctx, store := jobsTestStore(t)
	indexer := wantedSyncIndexer(t, ctx, store, `<rss><channel></channel></rss>`)
	wanted := wantedSyncMedia(t, ctx, store, "Wanted Empty")
	wantedSyncDownloadClient(t, ctx, store, wantedSyncSAB(t, true))

	if err := runWantedRSSSyncWorker(
		ctx,
		nil,
		store,
		indexers.NewService(indexer.Client()),
		downloadclients.NewService(indexer.Client()),
		decisions.NewEngine(),
		nil,
	); err != nil {
		t.Fatalf("wanted rss sync: %v", err)
	}

	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 0 || len(snapshot.Errors) != 1 || snapshot.Errors[0] != "No releases found" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func TestSCNMedia014WantedRSSSyncRejectedCandidateFlow(t *testing.T) {
	ctx, store := jobsTestStore(t)
	feed := wantedSyncFeed("Different.Movie.2026.1080p.WEB-DL", "https://indexer.test/download/rejected")
	indexer := wantedSyncIndexer(t, ctx, store, feed)
	wanted := wantedSyncMedia(t, ctx, store, "Wanted Rejected")
	wantedSyncDownloadClient(t, ctx, store, wantedSyncSAB(t, true))

	if err := runWantedRSSSyncWorker(
		ctx,
		nil,
		store,
		indexers.NewService(indexer.Client()),
		downloadclients.NewService(indexer.Client()),
		decisions.NewEngine(),
		nil,
	); err != nil {
		t.Fatalf("wanted rss sync: %v", err)
	}

	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 1 {
		t.Fatalf("releases = %#v", snapshot.Releases)
	}
	match := decisions.EvaluateReleaseMatch(wanted, snapshot.Releases[0])
	if match.Severity != "error" || len(match.Details) == 0 {
		t.Fatalf("match = %#v", match)
	}
	activities, err := store.ListDownloadActivity(ctx)
	if err != nil || len(activities) != 0 {
		t.Fatalf("download activity = %#v, %v", activities, err)
	}
}

func TestSCNMedia014WantedRSSSyncAcceptedCandidateFlow(t *testing.T) {
	ctx, store := jobsTestStore(t)
	feed := wantedSyncFeed("Wanted.Accepted.2026.1080p.WEB-DL", "https://indexer.test/download/accepted")
	indexer := wantedSyncIndexer(t, ctx, store, feed)
	downloadClient := wantedSyncSAB(t, true)
	wanted := wantedSyncMedia(t, ctx, store, "Wanted Accepted")
	wantedSyncDownloadClient(t, ctx, store, downloadClient)

	if err := runWantedRSSSyncWorker(
		ctx,
		nil,
		store,
		indexers.NewService(indexer.Client()),
		downloadclients.NewService(downloadClient.Client()),
		decisions.NewEngine(),
		nil,
	); err != nil {
		t.Fatalf("wanted rss sync: %v", err)
	}

	activities, err := store.ListDownloadActivity(ctx)
	if err != nil {
		t.Fatalf("list download activity: %v", err)
	}
	if len(activities) != 1 || activities[0].Status != "grabbed" || activities[0].DownloadID == nil || *activities[0].DownloadID != "nzo-wanted" {
		t.Fatalf("download activity = %#v", activities)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("list release results: %v", err)
	}
	if len(snapshot.Releases) != 1 || snapshot.Releases[0].Title != "Wanted.Accepted.2026.1080p.WEB-DL" {
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

func wantedSyncIndexer(
	t *testing.T,
	ctx context.Context,
	store *storage.SettingsStore,
	feed string,
) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(feed))
	}))
	t.Cleanup(server.Close)
	if _, err := store.CreateIndexer(ctx, storage.IndexerInput{
		DefinitionID:       "generic-newznab",
		Name:               "Wanted RSS",
		Implementation:     "Newznab",
		ImplementationName: "Newznab",
		Protocol:           "usenet",
		BaseURL:            server.URL,
		APIKey:             stringPtr("key"),
		Categories:         []int32{2000},
		SupportsRSS:        true,
		SupportsSearch:     true,
		Enabled:            true,
		Priority:           1,
	}); err != nil {
		t.Fatalf("create indexer: %v", err)
	}
	return server
}

func wantedSyncSAB(t *testing.T, success bool) *httptest.Server {
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

func wantedSyncMedia(
	t *testing.T,
	ctx context.Context,
	store *storage.SettingsStore,
	title string,
) storage.MediaItem {
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

func wantedSyncDownloadClient(
	t *testing.T,
	ctx context.Context,
	store *storage.SettingsStore,
	server *httptest.Server,
) {
	t.Helper()
	if _, err := store.CreateDownloadClient(ctx, storage.DownloadClientInput{
		Name:     "Wanted SAB",
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

func wantedSyncFeed(title string, downloadURL string) string {
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
