package indexers

import (
	"context"
	"testing"

	"media-manager/internal/acceptance"
	"media-manager/internal/testmocks"
)

func TestScenarioSCNIntegrations001TorznabMockCapabilities(t *testing.T) {
	scenario, err := acceptance.RequireScenario("features/behavior", "SCN-INTEGRATIONS-001")
	if err != nil {
		t.Fatal(err)
	}
	if !scenario.HasTag("integration") {
		t.Fatal("scenario must be tagged @integration")
	}

	mock := testmocks.NewProviderServer()
	t.Cleanup(mock.Close)

	result := NewService(mock.Client()).Test(context.Background(), Config{
		Type:    "torznab",
		BaseURL: mock.URL + "/torznab/api",
	})

	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
	if got := result.Details["categoryCount"]; got != 4 {
		t.Fatalf("categoryCount = %v", got)
	}

	releases, err := NewService(mock.Client()).Search(context.Background(), Config{
		ID:         "scenario-indexer",
		Name:       "Scenario Torznab",
		Type:       "torznab",
		BaseURL:    mock.URL + "/torznab/api",
		Categories: []int32{2000, 2040},
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d, releases = %#v", len(releases), releases)
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026.1080p.WEB-DL" || release.DownloadURL == "" {
		t.Fatalf("release = %#v", release)
	}
	if release.Seeders == nil || *release.Seeders != 42 || release.Peers == nil || *release.Peers != 7 {
		t.Fatalf("release peers/seeders = %#v", release)
	}
}
