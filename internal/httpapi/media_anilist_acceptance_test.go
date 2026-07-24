package httpapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"media-manager/internal/storage"
	"media-manager/internal/testmocks"
)

func TestAniListAPIUsesFixtureAndReportsCacheMissAndHit(t *testing.T) {
	providerServer := testmocks.NewProviderServer()
	t.Cleanup(providerServer.Close)
	store := testSettingsStore(t)
	client := newAcceptanceClientWithStore(t, "SCN-MEDIA-008", store)

	var provider MetadataProvider
	client.doJSON(t, http.MethodPost, "/settings/metadata-providers", MetadataProviderRequest{
		Name: "AAA Fixture AniList", Type: Anilist, BaseUrl: providerServer.URL + "/anilist/graphql", Enabled: true, Priority: 0,
	}, http.StatusCreated, &provider)

	var selected MediaSearchResult
	for range 2 {
		var search MediaSearchResponse
		client.doJSON(t, http.MethodPost, "/media/search", MediaSearchRequest{Query: "Frieren", Type: MediaTypeSerie}, http.StatusOK, &search)
		if len(search.Results) != 1 || search.Results[0].ExternalProvider == nil || *search.Results[0].ExternalProvider != "anilist" || search.Results[0].Aliases == nil || len(*search.Results[0].Aliases) != 3 {
			t.Fatalf("AniList search = %#v", search.Results)
		}
		selected = search.Results[0]
	}

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 {
		t.Fatal("expected seeded media profile")
	}
	libraryPath := filepath.Join(t.TempDir(), "series")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	var folder LibraryFolderCreateResponse
	client.doJSON(t, http.MethodPost, "/settings/library/folders", LibraryFolderRequest{Path: libraryPath, Kind: LibraryFolderKindSeries}, http.StatusCreated, &folder)
	var created MediaItem
	client.doJSON(t, http.MethodPost, "/media/items", MediaItemCreateRequest{
		Title: selected.Title, Type: selected.Type, Year: selected.Year, Monitored: true, MonitorMode: AllEpisodes,
		MinimumAvailability: Released, ExternalProvider: selected.ExternalProvider, ExternalId: selected.ExternalId,
		Overview: selected.Overview, PosterPath: selected.PosterPath, QualityProfileId: &profiles.Profiles[0].Id,
		LibraryFolderId: &folder.Folder.Id,
	}, http.StatusCreated, &created)
	if created.ExternalProvider == nil || *created.ExternalProvider != "anilist" || created.ExternalId == nil || *created.ExternalId != "154587" || created.ContentKind == nil || *created.ContentKind != MediaContentKindAnime || created.Aliases == nil || len(*created.Aliases) < 5 {
		t.Fatalf("created selected AniList result = %#v", created)
	}

	var details MediaMetadataDetails
	client.doJSON(t, http.MethodGet, "/media/metadata/anilist/serie/154587", nil, http.StatusOK, &details)
	if details.Title != "Frieren: Beyond Journey's End" || details.ExternalId != "154587" || details.Type != MediaTypeSerie {
		t.Fatalf("AniList details = %#v", details)
	}

	var diagnostics MetadataCacheResponse
	client.doJSON(t, http.MethodGet, "/settings/metadata-cache?cacheLimit=20&historyLimit=20", nil, http.StatusOK, &diagnostics)
	miss, hit := false, false
	for _, entry := range diagnostics.HistoryEntries {
		if entry.ProviderType != Anilist || entry.MediaType != MediaTypeSerie || !strings.Contains(entry.Query, "frieren") {
			continue
		}
		if entry.CacheHit {
			hit = true
		} else {
			miss = true
		}
	}
	if !miss || !hit {
		t.Fatalf("AniList cache diagnostics missing miss/hit: %#v", diagnostics.HistoryEntries)
	}

	primarySource := map[string]any{"source": "tmdb_import", "raw": map[string]any{"id": float64(209867), "title": "TMDB Primary Title"}}
	item, err := store.CreateMediaItem(t.Context(), storage.MediaItemInput{
		Type: "serie", ContentKind: "anime", Title: "TMDB Primary Title", Monitored: true,
		ExternalProvider: stringPointer("tmdb"), ExternalID: stringPointer("209867"),
		ProviderMappings: []storage.MediaProviderMappingInput{
			{EntityType: "media_item", ProviderName: "tmdb", ProviderEntityType: "serie", ExternalID: "209867", Canonical: true, Source: primarySource},
			{EntityType: "media_item", ProviderName: "anilist", ProviderEntityType: "serie", ExternalID: "154587"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.doJSON(t, http.MethodPost, "/media/items/"+item.ID.String()+"/metadata/refresh", nil, http.StatusOK, nil)
	refreshed, err := store.GetMediaItem(t.Context(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.ExternalProvider == nil || *refreshed.ExternalProvider != "tmdb" || len(refreshed.Aliases) < 5 {
		t.Fatalf("refreshed item did not preserve TMDB and persist aliases: %#v", refreshed)
	}
	anilistMappingFound, primaryProvenancePreserved := false, false
	for _, mapping := range refreshed.ProviderMappings {
		if mapping.ProviderName == "anilist" && mapping.ExternalID == "154587" && mapping.Source["raw"] != nil {
			anilistMappingFound = true
		}
		if mapping.ProviderName == "tmdb" && mapping.ExternalID == "209867" && reflect.DeepEqual(mapping.Source, primarySource) {
			primaryProvenancePreserved = true
		}
	}
	if !anilistMappingFound || !primaryProvenancePreserved {
		t.Fatalf("refreshed mapping provenance missing: %#v", refreshed.ProviderMappings)
	}

	failureServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"fixture failure"}],"data":null}`))
	}))
	t.Cleanup(failureServer.Close)
	client.doJSON(t, http.MethodPut, "/settings/metadata-providers/"+provider.Id.String(), MetadataProviderRequest{
		Name: provider.Name, Type: Anilist, BaseUrl: failureServer.URL, Enabled: true, Priority: 0,
	}, http.StatusOK, &provider)
	client.doJSON(t, http.MethodPost, "/media/items/"+item.ID.String()+"/metadata/refresh", nil, http.StatusBadGateway, nil)
	afterFailure, err := store.GetMediaItem(t.Context(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterFailure.Title != refreshed.Title || len(afterFailure.Aliases) != len(refreshed.Aliases) || len(afterFailure.ProviderMappings) != len(refreshed.ProviderMappings) {
		t.Fatalf("failed refresh partially wrote metadata: before=%#v after=%#v", refreshed, afterFailure)
	}
}
