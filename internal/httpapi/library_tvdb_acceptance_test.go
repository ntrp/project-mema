package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/testmocks"
)

func TestLibraryImportStoresTVDBMetadata(t *testing.T) {
	provider := testmocks.NewProviderServer()
	t.Cleanup(provider.Close)
	client := newAcceptanceClientWithProviders(t, "SCN-SETTINGS-016", provider)
	createScenarioTVDBProvider(t, client, provider.URL+"/tvdb/v4")

	libraryPath := filepath.Join(t.TempDir(), "library")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "Example.TVDB.Movie.2026.mkv"), []byte("video"), 0o644); err != nil {
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

	var imported LibraryScanImportResponse
	client.doJSON(t, http.MethodPost, "/settings/library/scans/"+created.Scan.Id.String()+"/import", LibraryScanImportRequest{
		Items: []LibraryScanImportRowRequest{{
			ItemId: created.Scan.Items[0].Id,
			Match: LibraryScanItemMatchRequest{
				MediaKind:           LibraryMediaKindMovie,
				Title:               "Example TVDB Movie",
				Year:                int32Ptr(2026),
				Monitored:           true,
				QualityProfileId:    profiles.Profiles[0].Id,
				MonitorMode:         OnlyMedia,
				MinimumAvailability: Released,
				ExternalProvider:    stringPtr("tvdb"),
				ExternalId:          stringPtr("900"),
			},
		}},
	}, http.StatusOK, &imported)
	if imported.ImportedCount != 1 || len(imported.MediaItems) != 1 {
		t.Fatalf("import response = %#v", imported)
	}
	item := imported.MediaItems[0]
	if item.Overview == nil || *item.Overview != "A realistic TVDB metadata detail response." {
		t.Fatalf("imported TVDB media missing enriched overview: %#v", item)
	}
	if item.RuntimeMinutes == nil || *item.RuntimeMinutes != 101 {
		t.Fatalf("imported TVDB media missing runtime: %#v", item)
	}
	if item.BackdropPath == nil || *item.BackdropPath != "/tvdb-backdrop.jpg" {
		t.Fatalf("imported TVDB media missing backdrop: %#v", item)
	}
	if item.VoteAverage != nil {
		t.Fatalf("TVDB popularity score stored as rating: %#v", item)
	}
	if mediaItemFact(item, "Revenue") != "$533,300,000.00" || mediaItemFact(item, "Budget") != "$120,000,000.00" {
		t.Fatalf("imported TVDB media missing formatted money facts: %#v", item.Facts)
	}
	if mediaItemFact(item, "Production Countries") != "🇺🇸 United States" {
		t.Fatalf("imported TVDB media missing flagged country fact: %#v", item.Facts)
	}
}

func createScenarioTVDBProvider(t *testing.T, client acceptanceClient, baseURL string) {
	t.Helper()
	apiKey := "scenario-tvdb-key"
	var created MetadataProvider
	client.doJSON(t, http.MethodPost, "/settings/metadata-providers", MetadataProviderRequest{
		Name:     "AAA Scenario TVDB",
		Type:     Tvdb,
		BaseUrl:  baseURL,
		ApiKey:   &apiKey,
		Enabled:  true,
		Priority: 0,
	}, http.StatusCreated, &created)
}

func mediaItemFact(item MediaItem, label string) string {
	if item.Facts == nil {
		return ""
	}
	for _, fact := range *item.Facts {
		if fact.Label == label {
			return fact.Value
		}
	}
	return ""
}
