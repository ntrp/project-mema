package decisions

import (
	"testing"

	"media-manager/internal/storage"
)

func TestChooseReleasePrefersHighestDetectedQuality(t *testing.T) {
	engine := NewEngine()
	seeders := int32(100)
	item := storage.MediaItem{Type: "movie", Title: "Movie"}
	decision, ok := engine.ChooseRelease(item, []storage.ReleaseCandidateInput{
		{Title: "Movie 2026 WEB-DL 1080p", SizeBytes: 20, Seeders: &seeders},
		{Title: "Movie 2026 Remux 2160p", SizeBytes: 10},
	})
	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.Title != "Movie 2026 Remux 2160p" {
		t.Fatalf("expected 2160p remux, got %q", decision.Release.Title)
	}
}

func TestChooseReleaseUsesSeedersAndSizeAsTiebreakers(t *testing.T) {
	engine := NewEngine()
	lowSeeders := int32(3)
	highSeeders := int32(8)
	item := storage.MediaItem{Type: "movie", Title: "Movie"}
	decision, ok := engine.ChooseRelease(item, []storage.ReleaseCandidateInput{
		{Title: "Movie 2026 WEB-DL 1080p", SizeBytes: 200, Seeders: &lowSeeders},
		{Title: "Movie 2026 WEB-DL 1080p Proper", SizeBytes: 100, Seeders: &highSeeders},
	})
	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.Title != "Movie 2026 WEB-DL 1080p Proper" {
		t.Fatalf("expected higher seeder release, got %q", decision.Release.Title)
	}
}

func TestChooseReleaseRejectsMismatchedResource(t *testing.T) {
	engine := NewEngine()
	item := storage.MediaItem{Type: "movie", Title: "Expected Movie"}
	_, ok := engine.ChooseRelease(item, []storage.ReleaseCandidateInput{
		{Title: "Other Movie 2026 Remux 2160p", SizeBytes: 10},
	})
	if ok {
		t.Fatal("expected mismatched release to be rejected")
	}
}
