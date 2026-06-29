package decisions

import (
	"testing"

	"media-manager/internal/storage"
)

func TestChooseReleasePrefersHighestDetectedQuality(t *testing.T) {
	engine := NewEngine()
	seeders := int32(100)
	decision, ok := engine.ChooseRelease([]storage.ReleaseCandidateInput{
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
	decision, ok := engine.ChooseRelease([]storage.ReleaseCandidateInput{
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
