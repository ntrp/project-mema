package jobs

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestSCNMedia002AutoSearchSkipsAlreadyHandledMedia(t *testing.T) {
	for _, status := range []string{"downloaded", "downloading"} {
		item := storage.MediaItem{ID: uuid.New(), Type: "movie", Title: "Scenario Movie", Status: status}

		if err := autoSearchDownload(context.Background(), nil, nil, nil, decisions.NewEngine(), nil, item); err != nil {
			t.Fatalf("autoSearchDownload(%s) returned error: %v", status, err)
		}
	}
}

func TestSCNMedia002TopDecisionRejectionsReturnUniqueBlockingReasons(t *testing.T) {
	item := storage.MediaItem{Type: "movie", Title: "Scenario Movie", Year: int32Ptr(2026)}
	releases := []storage.ReleaseCandidateInput{
		{Title: "Different.Movie.2026.1080p.WEBDL"},
		{Title: "Different.Movie.2026.720p.WEBDL"},
		{Title: "Another.Movie.2025.1080p.WEBDL"},
	}

	reasons := topDecisionRejections(item, nil, nil, nil, releases)

	if len(reasons) == 0 || len(reasons) > 3 {
		t.Fatalf("reasons = %#v", reasons)
	}
	seen := map[string]struct{}{}
	for _, reason := range reasons {
		if strings.TrimSpace(reason) == "" {
			t.Fatalf("blank rejection reason in %#v", reasons)
		}
		if _, ok := seen[reason]; ok {
			t.Fatalf("duplicate rejection reason %q in %#v", reason, reasons)
		}
		seen[reason] = struct{}{}
	}
}

func TestAutoSearchExcludesBlocklistedReleaseCandidates(t *testing.T) {
	ctx, store := jobsTestStore(t)
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type:      "movie",
		Title:     "Scenario Movie " + uuid.NewString(),
		Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	guid := "blocked-" + uuid.NewString()
	blocked := storage.ReleaseCandidateInput{
		MediaItemID:     item.ID,
		IndexerName:     "Scenario Indexer",
		IndexerProtocol: "usenet",
		Title:           "Scenario.Movie.2026.1080p.WEB-DL",
		DownloadURL:     "https://indexer.test/blocked",
		GUID:            &guid,
		SizeBytes:       8_000_000_000,
	}
	allowed := storage.ReleaseCandidateInput{
		MediaItemID:     item.ID,
		IndexerName:     "Scenario Indexer",
		IndexerProtocol: "usenet",
		Title:           "Scenario.Movie.2026.720p.WEB-DL",
		DownloadURL:     "https://indexer.test/allowed",
		SizeBytes:       4_000_000_000,
	}
	expiresAt := time.Now().Add(time.Hour)
	if _, err := store.BlockReleaseCandidate(ctx, blocked, "client rejected", "download_failed", &expiresAt); err != nil {
		t.Fatal(err)
	}

	filtered, err := unblockedReleaseCandidates(ctx, store, []storage.ReleaseCandidateInput{blocked, allowed})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].Title != allowed.Title {
		t.Fatalf("filtered releases = %#v", filtered)
	}
}

func TestAutoSearchAlternativeRetryIsBounded(t *testing.T) {
	err := fmt.Errorf("%w: Scenario.Movie.2026.1080p", errRetryAlternativeRelease)
	if !shouldRetryAlternativeRelease(err, 1) {
		t.Fatal("first retry should be allowed")
	}
	if shouldRetryAlternativeRelease(err, maxAutomaticGrabAttempts) {
		t.Fatal("retry at limit should not be allowed")
	}
	if !automaticRetryLimitReached(err, maxAutomaticGrabAttempts) {
		t.Fatal("retry limit should be reached")
	}
	if automaticRetryLimitReached(context.Canceled, maxAutomaticGrabAttempts) {
		t.Fatal("non retry errors should not hit retry limit")
	}
}

func TestSCNMedia009DownloadClientConfigAndOptionalStrings(t *testing.T) {
	username := "user"
	password := "secret"
	apiKey := "key"
	category := "movies"

	config := downloadClientConfig(storage.DownloadClient{
		Name:     "Scenario Client",
		Type:     "sabnzbd",
		BaseURL:  "http://client.local",
		Username: &username,
		Password: &password,
		APIKey:   &apiKey,
		Category: &category,
	})

	if config.Name != "Scenario Client" || config.Type != "sabnzbd" || config.BaseURL == "" {
		t.Fatalf("download client config = %#v", config)
	}
	if config.Username == nil || *config.Username != "user" || config.Category == nil {
		t.Fatalf("download client optional config = %#v", config)
	}
	if optionalString("  download-1  ") == nil || *optionalString("download-1") != "download-1" {
		t.Fatal("non-empty optional string should be preserved")
	}
	if optionalString("   ") != nil {
		t.Fatal("blank optional string should become nil")
	}
}

func TestSCNSystem008JobArgumentKindsAreStable(t *testing.T) {
	if (ReleaseSearchArgs{}).Kind() != "media.release_search" {
		t.Fatalf("release search kind = %q", (ReleaseSearchArgs{}).Kind())
	}
	if (AutoSearchDownloadArgs{}).Kind() != "media.auto_search_download" {
		t.Fatalf("auto search kind = %q", (AutoSearchDownloadArgs{}).Kind())
	}
	if (RSSSyncArgs{}).Kind() != "media.rss_sync" {
		t.Fatalf("rss sync kind = %q", (RSSSyncArgs{}).Kind())
	}
	if (GrabReleaseArgs{}).Kind() != "media.grab_release" {
		t.Fatalf("grab release kind = %q", (GrabReleaseArgs{}).Kind())
	}
	if (DownloadActivitySyncArgs{}).Kind() != "download.activity_sync" {
		t.Fatalf("activity sync kind = %q", (DownloadActivitySyncArgs{}).Kind())
	}
	if (ReleaseBlocklistCleanupArgs{}).Kind() != "release.blocklist_cleanup" {
		t.Fatalf("blocklist cleanup kind = %q", (ReleaseBlocklistCleanupArgs{}).Kind())
	}
	if (SubtitleSearchArgs{}).Kind() != "media.subtitle_search" {
		t.Fatalf("subtitle search kind = %q", (SubtitleSearchArgs{}).Kind())
	}
	if (SubtitleRetryArgs{}).Kind() != "media.subtitle_retry" {
		t.Fatalf("subtitle retry kind = %q", (SubtitleRetryArgs{}).Kind())
	}
}

func int32Ptr(value int32) *int32 {
	return &value
}

func jobsTestStore(t *testing.T) (context.Context, *storage.SettingsStore) {
	t.Helper()
	ctx := context.Background()
	databaseURL := testdb.Create(t)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	return ctx, storage.NewSettingsStore(pool)
}
