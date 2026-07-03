package httpapi

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"media-manager/internal/storage"
)

func TestScenarioSCNSettings020AdminInspectsAndClearsIndexerSearchCache(t *testing.T) {
	store := testSettingsStore(t)
	client := newAcceptanceClientWithStore(t, "SCN-SETTINGS-020", store)
	ctx := context.Background()
	indexer, query := seedIndexerSearchState(t, ctx, store)

	var initial IndexerSearchResponse
	client.doJSON(t, http.MethodGet, "/settings/indexer-search?cacheLimit=5&historyLimit=5", nil, http.StatusOK, &initial)
	if initial.Stats.TotalEntries == 0 || initial.HistoryTotalEntries == 0 {
		t.Fatalf("initial indexer search response = %#v", initial)
	}

	var updated IndexerSearchResponse
	client.doJSON(t, http.MethodPut, "/settings/indexer-search", IndexerSearchSettings{
		CacheDurationMinutes: 60,
		HistoryRetentionDays: 14,
	}, http.StatusOK, &updated)
	if updated.Settings.CacheDurationMinutes != 60 || updated.Settings.HistoryRetentionDays != 14 {
		t.Fatalf("updated indexer search settings = %#v", updated.Settings)
	}

	var deleted MetadataCacheClearResponse
	path := "/settings/indexer-search/cache/entry?indexerId=" + indexer.ID.String() +
		"&mediaType=movie&query=" + url.QueryEscape(query)
	client.doJSON(t, http.MethodDelete, path, nil, http.StatusOK, &deleted)
	if deleted.DeletedCount != 1 {
		t.Fatalf("deleted cache entry = %#v", deleted)
	}

	seeded, _ := store.SetIndexerSearchCache(ctx, indexer.ID, "movie", query+" again", []map[string]string{{"title": "Again"}}, 1, time.Now().Add(time.Hour))
	if seeded.Query == "" {
		t.Fatal("expected seeded cache entry")
	}
	var patternCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodPost, "/settings/indexer-search/cache/reset", MetadataCacheClearRequest{
		Pattern: "Scenario",
	}, http.StatusOK, &patternCleared)
	if patternCleared.DeletedCount == 0 {
		t.Fatalf("pattern clear result = %#v", patternCleared)
	}

	if _, err := store.SetIndexerSearchCache(ctx, indexer.ID, "movie", "Full Reset Scenario", []map[string]string{{"title": "Reset"}}, 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var cacheCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodDelete, "/settings/indexer-search/cache", nil, http.StatusOK, &cacheCleared)
	if cacheCleared.DeletedCount == 0 {
		t.Fatalf("full cache clear result = %#v", cacheCleared)
	}

	var historyCleared MetadataCacheClearResponse
	client.doJSON(t, http.MethodDelete, "/settings/indexer-search/history", nil, http.StatusOK, &historyCleared)
	if historyCleared.DeletedCount == 0 {
		t.Fatalf("history clear result = %#v", historyCleared)
	}
}

func seedIndexerSearchState(t *testing.T, ctx context.Context, store *storage.SettingsStore) (storage.Indexer, string) {
	t.Helper()
	indexer, err := store.CreateIndexer(ctx, storage.IndexerInput{
		Name:       "Scenario Search Indexer",
		Type:       "torznab",
		BaseURL:    "http://indexer.local",
		Categories: []int32{2000},
		Enabled:    true,
		Priority:   10,
	})
	if err != nil {
		t.Fatal(err)
	}
	query := "Scenario Movie"
	if _, err := store.SetIndexerSearchCache(ctx, indexer.ID, "movie", query, []map[string]string{{"title": "Result"}}, 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordIndexerSearchHistory(ctx, storage.IndexerSearchHistoryInput{
		IndexerID:   indexer.ID,
		IndexerName: indexer.Name,
		IndexerType: indexer.Type,
		MediaType:   "movie",
		Query:       query,
		CacheHit:    false,
		Success:     true,
		ResultCount: 1,
		Response:    []map[string]string{{"title": "Result"}},
	}); err != nil {
		t.Fatal(err)
	}
	return indexer, query
}
