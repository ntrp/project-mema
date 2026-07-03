package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings015StorageIndexerSearchCacheAndHistory(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-015")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	indexer, err := store.CreateIndexer(ctx, IndexerInput{
		Name:       "Cache Indexer " + suffix,
		Protocol:   "torrent",
		BaseURL:    "http://indexer.cache/" + suffix,
		Categories: []int32{},
		Enabled:    true,
		Priority:   1,
	})
	if err != nil {
		t.Fatalf("create indexer: %v", err)
	}

	query := "cache-query-" + suffix
	entry, err := store.SetIndexerSearchCache(ctx, indexer.ID, "movie", query, []map[string]any{{"title": "Result"}}, 1, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("set cache: %v", err)
	}
	if entry.Query != query || entry.ResultCount != 1 || entry.Expired {
		t.Fatalf("cache entry = %#v", entry)
	}

	var cached []map[string]any
	hit, err := store.GetIndexerSearchCache(ctx, indexer.ID, "movie", query, &cached)
	if err != nil {
		t.Fatalf("get cache: %v", err)
	}
	if !hit || cached[0]["title"] != "Result" {
		t.Fatalf("cache hit=%v cached=%#v", hit, cached)
	}

	errorText := "timeout"
	if _, err := store.RecordIndexerSearchHistory(ctx, IndexerSearchHistoryInput{
		IndexerID:       indexer.ID,
		IndexerName:     indexer.Name,
		IndexerProtocol: indexer.Protocol,
		MediaType:       "movie",
		Query:           query,
		CacheHit:        true,
		Success:         false,
		ResultCount:     1,
		Error:           &errorText,
		Response:        map[string]any{"error": errorText},
	}); err != nil {
		t.Fatalf("record history: %v", err)
	}

	stats, err := store.IndexerSearchHistoryStats(ctx)
	if err != nil {
		t.Fatalf("history stats: %v", err)
	}
	if stats.TotalEntries == 0 || stats.CacheHits == 0 || stats.Failures == 0 {
		t.Fatalf("history stats = %#v", stats)
	}

	deleted, err := store.DeleteIndexerSearchCacheEntry(ctx, indexer.ID, "movie", query)
	if err != nil {
		t.Fatalf("delete cache entry: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted cache entries = %d", deleted)
	}

	if _, err := store.SetIndexerSearchCache(ctx, indexer.ID, "movie", query+"-clear", []map[string]any{{"title": "Result"}}, 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("set cache before clear: %v", err)
	}
	cleared, err := store.ClearIndexerSearchCache(ctx)
	if err != nil {
		t.Fatalf("clear cache: %v", err)
	}
	if cleared < 1 {
		t.Fatalf("cleared indexer cache entries = %d, want at least 1", cleared)
	}
}

func TestScenarioSCNSettings015StorageMetadataCacheAndHistory(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-015")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	provider, err := store.CreateMetadataProvider(ctx, MetadataProviderInput{
		Name:     "Cache Provider " + suffix,
		Type:     "tmdb",
		BaseURL:  "http://metadata.cache/" + suffix,
		Enabled:  true,
		Priority: 1,
	})
	if err != nil {
		t.Fatalf("create provider: %v", err)
	}

	query := "metadata-query-" + suffix
	year := int32(2026)
	if err := store.SetMetadataSearchCache(ctx, provider.ID, "movie", query, &year, []map[string]any{{"title": "Movie"}}, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("set metadata cache: %v", err)
	}
	var cached []map[string]any
	hit, err := store.GetMetadataSearchCache(ctx, provider.ID, "movie", query, &year, &cached)
	if err != nil {
		t.Fatalf("get metadata cache: %v", err)
	}
	if !hit || cached[0]["title"] != "Movie" {
		t.Fatalf("metadata cache hit=%v cached=%#v", hit, cached)
	}

	if _, err := store.RecordMetadataSearchHistory(ctx, MetadataSearchHistoryInput{
		ProviderID:   provider.ID,
		ProviderName: provider.Name,
		ProviderType: provider.Type,
		MediaType:    "movie",
		Query:        "details:" + query,
		Year:         &year,
		CacheHit:     false,
		Success:      true,
		ItemCount:    1,
		Response:     []map[string]any{{"title": "Movie"}},
	}); err != nil {
		t.Fatalf("record metadata history: %v", err)
	}

	entries, err := store.ListMetadataSearchHistoryEntries(ctx, 1)
	if err != nil {
		t.Fatalf("list metadata history: %v", err)
	}
	if len(entries) == 0 || entries[0].CacheKind != "details" {
		t.Fatalf("metadata history entries = %#v", entries)
	}

	deleted, err := store.DeleteMetadataCacheEntry(ctx, provider.ID, "movie", query, year)
	if err != nil {
		t.Fatalf("delete metadata cache entry: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted metadata cache entries = %d", deleted)
	}

	if err := store.SetMetadataSearchCache(ctx, provider.ID, "movie", query+"-clear", &year, []map[string]any{{"title": "Movie"}}, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("set metadata cache before clear: %v", err)
	}
	cleared, err := store.ClearMetadataCache(ctx)
	if err != nil {
		t.Fatalf("clear metadata cache: %v", err)
	}
	if cleared < 1 {
		t.Fatalf("cleared metadata cache entries = %d, want at least 1", cleared)
	}
}
