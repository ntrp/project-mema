package httpapi

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"media-manager/internal/storage"
)

func TestScenarioSCNSettings021AdminInspectsAndClearsMetadataCache(t *testing.T) {
	store := testSettingsStore(t)
	client := newAcceptanceClientWithStore(t, "SCN-SETTINGS-021", store)
	ctx := context.Background()
	provider, query, year := seedMetadataCacheState(t, ctx, store)

	var initial MetadataCacheResponse
	client.doJSON(t, http.MethodGet, "/settings/metadata-cache?cacheLimit=5&historyLimit=5", nil, http.StatusOK, &initial)
	if initial.Stats.TotalEntries == 0 || initial.HistoryTotalEntries == 0 {
		t.Fatalf("initial metadata cache response = %#v", initial)
	}

	var deleted MetadataCacheClearResponse
	path := "/settings/metadata-cache/entry?providerId=" + provider.ID.String() +
		"&mediaType=movie&query=" + url.QueryEscape(query) + "&year=2026"
	client.doJSON(t, http.MethodDelete, path, nil, http.StatusOK, &deleted)
	if deleted.DeletedCount != 1 {
		t.Fatalf("deleted metadata cache entry = %#v", deleted)
	}

	if err := store.SetMetadataSearchCache(ctx, provider.ID, "movie", query+" sequel", &year, []map[string]string{{"title": "Again"}}, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var patternCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodPost, "/settings/metadata-cache/reset", MetadataCacheClearRequest{
		Pattern: "Scenario",
	}, http.StatusOK, &patternCleared)
	if patternCleared.DeletedCount == 0 {
		t.Fatalf("pattern clear result = %#v", patternCleared)
	}

	if err := store.SetMetadataSearchCache(ctx, provider.ID, "movie", "Full Reset Scenario", &year, []map[string]string{{"title": "Reset"}}, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var cacheCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodDelete, "/settings/metadata-cache", nil, http.StatusOK, &cacheCleared)
	if cacheCleared.DeletedCount == 0 {
		t.Fatalf("full cache clear result = %#v", cacheCleared)
	}

	var historyCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodDelete, "/settings/metadata-cache/history", nil, http.StatusOK, &historyCleared)
	if historyCleared.DeletedCount == 0 {
		t.Fatalf("history clear result = %#v", historyCleared)
	}
}

func seedMetadataCacheState(t *testing.T, ctx context.Context, store *storage.SettingsStore) (storage.MetadataProvider, string, int32) {
	t.Helper()
	if err := store.EnsureDefaultMetadataProviders(ctx); err != nil {
		t.Fatal(err)
	}
	providers, err := store.ListMetadataProviders(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(providers) == 0 {
		t.Fatal("expected default metadata provider")
	}
	provider := providers[0]
	query := "Scenario Movie"
	year := int32(2026)
	if err := store.SetMetadataSearchCache(ctx, provider.ID, "movie", query, &year, []map[string]string{{"title": "Result"}}, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordMetadataSearchHistory(ctx, storage.MetadataSearchHistoryInput{
		ProviderID:   provider.ID,
		ProviderName: provider.Name,
		ProviderType: provider.Type,
		MediaType:    "movie",
		Query:        query,
		Year:         &year,
		CacheHit:     false,
		Success:      true,
		ItemCount:    1,
		Response:     []map[string]string{{"title": "Result"}},
	}); err != nil {
		t.Fatal(err)
	}
	return provider, query, year
}
