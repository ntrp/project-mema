package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings014StorageIndexerLifecycleAndHealth(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-014")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	categories := []int32{2000, 2040}
	indexer, err := store.CreateIndexer(ctx, IndexerInput{
		Name:       "Indexer " + suffix,
		Protocol:   "torrent",
		BaseURL:    "http://indexer.test/" + suffix,
		APIKey:     stringPtr("key"),
		Categories: categories,
		Enabled:    true,
		Priority:   30,
	})
	if err != nil {
		t.Fatalf("create indexer: %v", err)
	}

	statusCode := int32(429)
	failed, err := store.RecordIndexerFailure(ctx, indexer.ID, &statusCode, "rate limited", false, nil)
	if err != nil {
		t.Fatalf("record failure: %v", err)
	}
	if failed.HealthStatus != "temporary_disabled" || failed.FailureCount != 1 || failed.NextCheckAt == nil {
		t.Fatalf("failed indexer = %#v", failed)
	}

	healthy, err := store.RecordIndexerSuccess(ctx, indexer.ID)
	if err != nil {
		t.Fatalf("record success: %v", err)
	}
	if healthy.HealthStatus != "healthy" || healthy.FailureCount != 0 || healthy.LastError != nil {
		t.Fatalf("healthy indexer = %#v", healthy)
	}

	updated, err := store.UpdateIndexer(ctx, indexer.ID, IndexerInput{
		Name:       "Updated " + suffix,
		Protocol:   "torrent",
		BaseURL:    "http://indexer.test/updated/" + suffix,
		Categories: []int32{5000},
		Enabled:    false,
		Priority:   40,
	})
	if err != nil {
		t.Fatalf("update indexer: %v", err)
	}
	if updated.Enabled || updated.Priority != 40 || updated.HealthStatus != "healthy" {
		t.Fatalf("updated indexer = %#v", updated)
	}

	if err := store.DeleteIndexer(ctx, indexer.ID); err != nil {
		t.Fatalf("delete indexer: %v", err)
	}
	if _, err := store.GetIndexer(ctx, indexer.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted indexer to be missing, got %v", err)
	}
}

func TestScenarioSCNSettings014StorageDownloadClientLifecycle(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-014")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	client, err := store.CreateDownloadClient(ctx, DownloadClientInput{
		Name:     "SAB " + suffix,
		Type:     "sabnzbd",
		Protocol: "usenet",
		BaseURL:  "http://sab.test/" + suffix,
		APIKey:   stringPtr("api-key"),
		Category: stringPtr("movies"),
		Enabled:  true,
		Priority: 5,
	})
	if err != nil {
		t.Fatalf("create download client: %v", err)
	}

	updated, err := store.UpdateDownloadClient(ctx, client.ID, DownloadClientInput{
		Name:     "Transmission " + suffix,
		Type:     "transmission",
		Protocol: "torrent",
		BaseURL:  "http://transmission.test/" + suffix,
		Username: stringPtr("user"),
		Password: stringPtr("pass"),
		Enabled:  false,
		Priority: 10,
	})
	if err != nil {
		t.Fatalf("update download client: %v", err)
	}
	if updated.Enabled || updated.Type != "transmission" || updated.Protocol != "torrent" || updated.Username == nil {
		t.Fatalf("updated client = %#v", updated)
	}

	enabled, err := store.ListEnabledDownloadClients(ctx)
	if err != nil {
		t.Fatalf("list enabled clients: %v", err)
	}
	for _, item := range enabled {
		if item.ID == updated.ID {
			t.Fatalf("disabled client should not be listed as enabled: %#v", enabled)
		}
	}

	if err := store.DeleteDownloadClient(ctx, client.ID); err != nil {
		t.Fatalf("delete download client: %v", err)
	}
	if _, err := store.GetDownloadClient(ctx, client.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted client to be missing, got %v", err)
	}
}

func TestScenarioSCNSettings014StorageMetadataProviderLifecycle(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-014")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	provider, err := store.CreateMetadataProvider(ctx, MetadataProviderInput{
		Name:        "Provider " + suffix,
		Type:        "tmdb",
		BaseURL:     "http://metadata.test/" + suffix,
		APIKey:      stringPtr("key"),
		PIN:         stringPtr("1234"),
		AccessToken: stringPtr("access"),
		Enabled:     true,
		Priority:    12,
	})
	if err != nil {
		t.Fatalf("create metadata provider: %v", err)
	}

	expiresAt := time.Now().Add(time.Hour)
	if err := store.UpdateMetadataProviderSessionToken(ctx, provider.ID, "session", expiresAt); err != nil {
		t.Fatalf("update session token: %v", err)
	}
	found, err := store.GetMetadataProvider(ctx, provider.ID)
	if err != nil {
		t.Fatalf("get metadata provider: %v", err)
	}
	if found.SessionToken == nil || *found.SessionToken != "session" {
		t.Fatalf("provider session token = %#v", found.SessionToken)
	}

	updated, err := store.UpdateMetadataProvider(ctx, provider.ID, MetadataProviderInput{
		Name:     "Provider Updated " + suffix,
		Type:     "tvdb",
		BaseURL:  "http://metadata.test/updated/" + suffix,
		Enabled:  false,
		Priority: 20,
	})
	if err != nil {
		t.Fatalf("update metadata provider: %v", err)
	}
	if updated.Enabled || updated.SessionToken != nil || updated.Type != "tvdb" {
		t.Fatalf("updated provider = %#v", updated)
	}

	if err := store.DeleteMetadataProvider(ctx, provider.ID); err != nil {
		t.Fatalf("delete metadata provider: %v", err)
	}
	if _, err := store.GetMetadataProvider(ctx, provider.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted provider to be missing, got %v", err)
	}
}

func TestSubtitleProviderLifecycle(t *testing.T) {
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	provider, err := store.CreateSubtitleProvider(ctx, SubtitleProviderInput{
		Name:     "OpenSubtitles " + suffix,
		Type:     "opensubtitles",
		BaseURL:  "https://api.opensubtitles.com",
		Username: stringPtr("user"),
		Password: stringPtr("secret"),
		APIKey:   stringPtr("key"),
		Enabled:  true,
		Priority: 10,
	})
	if err != nil {
		t.Fatalf("create subtitle provider: %v", err)
	}
	providers, err := store.ListSubtitleProviders(ctx)
	if err != nil {
		t.Fatalf("list subtitle providers: %v", err)
	}
	if len(providers) != 1 || providers[0].APIKey == nil || providers[0].Password == nil {
		t.Fatalf("providers = %#v", providers)
	}
	updated, err := store.UpdateSubtitleProvider(ctx, provider.ID, SubtitleProviderInput{
		Name:     "Updated " + suffix,
		Type:     "opensubtitles",
		BaseURL:  "https://api.opensubtitles.com",
		Enabled:  false,
		Priority: 20,
	})
	if err != nil {
		t.Fatalf("update subtitle provider: %v", err)
	}
	if updated.Enabled || updated.APIKey == nil || *updated.APIKey != "key" || updated.Password == nil || *updated.Password != "secret" {
		t.Fatalf("updated provider should preserve omitted secrets, got %#v", updated)
	}
	cleared, err := store.UpdateSubtitleProvider(ctx, provider.ID, SubtitleProviderInput{
		Name:              "Updated " + suffix,
		Type:              "opensubtitles",
		BaseURL:           "https://api.opensubtitles.com",
		Enabled:           false,
		Priority:          20,
		ClearSecretFields: []string{"apiKey", "password"},
	})
	if err != nil {
		t.Fatalf("clear subtitle provider secrets: %v", err)
	}
	if cleared.APIKey != nil || cleared.Password != nil || len(cleared.SecretFieldsSet) != 0 {
		t.Fatalf("cleared provider secrets = %#v", cleared)
	}
	if err := store.DeleteSubtitleProvider(ctx, provider.ID); err != nil {
		t.Fatalf("delete subtitle provider: %v", err)
	}
	if _, err := store.GetSubtitleProvider(ctx, provider.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted provider to be missing, got %v", err)
	}
}

func TestMockSubtitleProviderRowsLifecycle(t *testing.T) {
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	provider, err := store.CreateSubtitleProvider(ctx, SubtitleProviderInput{
		Name:    "Mock Subtitles " + suffix,
		Type:    "mock",
		BaseURL: "mock://subtitles",
		Enabled: true,
		MockSubtitles: []MockSubtitleProviderRowInput{{
			Title:      "Scenario Movie",
			LanguageID: "english",
			Format:     "vtt",
		}},
	})
	if err != nil {
		t.Fatalf("create subtitle provider: %v", err)
	}
	if len(provider.MockSubtitles) != 1 || provider.MockSubtitles[0].Format != "vtt" {
		t.Fatalf("created provider = %#v", provider)
	}
	found, err := store.GetSubtitleProvider(ctx, provider.ID)
	if err != nil {
		t.Fatalf("get subtitle provider: %v", err)
	}
	if len(found.MockSubtitles) != 1 || found.MockSubtitles[0].Title != "Scenario Movie" {
		t.Fatalf("found provider = %#v", found)
	}
	updated, err := store.UpdateSubtitleProvider(ctx, provider.ID, SubtitleProviderInput{
		Name:    "Mock Subtitles " + suffix,
		Type:    "mock",
		BaseURL: "mock://subtitles",
		Enabled: true,
		MockSubtitles: []MockSubtitleProviderRowInput{{
			Title:      "Scenario Movie",
			LanguageID: "german",
			Format:     "srt",
		}},
	})
	if err != nil {
		t.Fatalf("update subtitle provider: %v", err)
	}
	if len(updated.MockSubtitles) != 1 || updated.MockSubtitles[0].LanguageID != "german" {
		t.Fatalf("updated provider = %#v", updated)
	}
}
