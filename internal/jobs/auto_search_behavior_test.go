package jobs

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
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
	if (GrabReleaseArgs{}).Kind() != "media.grab_release" {
		t.Fatalf("grab release kind = %q", (GrabReleaseArgs{}).Kind())
	}
	if (MissingMediaRetryArgs{}).Kind() != "media.missing_media_retry" {
		t.Fatalf("missing retry kind = %q", (MissingMediaRetryArgs{}).Kind())
	}
	if (DownloadActivitySyncArgs{}).Kind() != "download.activity_sync" {
		t.Fatalf("activity sync kind = %q", (DownloadActivitySyncArgs{}).Kind())
	}
}

func int32Ptr(value int32) *int32 {
	return &value
}
