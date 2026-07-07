package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestLibraryImportResetClearsRowAndKeepsFile(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-016")
	libraryPath := filepath.Join(t.TempDir(), "library")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(libraryPath, "Reset.Movie.2026.mkv")
	if err := os.WriteFile(filePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 {
		t.Fatal("expected seeded media profile")
	}

	var created LibraryFolderCreateResponse
	client.doJSON(t, http.MethodPost, "/settings/library/folders", LibraryFolderRequest{
		Path: libraryPath,
		Kind: LibraryFolderKindMovie,
	}, http.StatusCreated, &created)
	if len(created.Scan.Items) == 0 {
		t.Fatalf("created scan has no items: %#v", created.Scan)
	}
	itemID := created.Scan.Items[0].Id

	var imported LibraryScanImportResponse
	client.doJSON(t, http.MethodPost, "/settings/library/scans/"+created.Scan.Id.String()+"/import", LibraryScanImportRequest{
		Items: []LibraryScanImportRowRequest{{
			ItemId: itemID,
			Match: LibraryScanItemMatchRequest{
				MediaKind:           LibraryMediaKindMovie,
				Title:               "Reset Movie",
				Year:                int32Ptr(2026),
				Monitored:           true,
				QualityProfileId:    profiles.Profiles[0].Id,
				MonitorMode:         OnlyMedia,
				MinimumAvailability: Released,
			},
		}},
	}, http.StatusOK, &imported)
	if imported.ImportedCount != 1 || imported.Scan.Items[0].Imported == false {
		t.Fatalf("import response = %#v", imported)
	}

	var reset LibraryScanItemResetResponse
	client.doJSON(t, http.MethodPost, "/settings/library/scans/"+created.Scan.Id.String()+"/items/"+itemID.String()+"/reset", nil, http.StatusOK, &reset)
	if reset.Item.Imported || reset.Item.Status != LibraryScanItemStatusPending || reset.Item.MediaItemId != nil {
		t.Fatalf("reset item = %#v", reset.Item)
	}
	if reset.RemovedMediaItemId == nil {
		t.Fatalf("reset did not remove imported media item: %#v", reset)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("reset touched media file: %v", err)
	}
}
