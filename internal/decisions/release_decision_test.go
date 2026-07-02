package decisions

import (
	"testing"

	"github.com/google/uuid"

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

func TestChooseReleaseQualityBeatsCustomFormatScore(t *testing.T) {
	engine := NewEngine()
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000201")
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p", "remux-2160p"},
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: 10000},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Preferred group",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "preferred", Name: "Preferred", Type: "releaseTitle", Value: "Preferred", Required: true,
		}},
	}}
	item := storage.MediaItem{Type: "movie", Title: "Movie"}
	decision, ok := engine.ChooseReleaseWithProfile(item, &profile, formats, []storage.ReleaseCandidateInput{
		{Title: "Movie 2026 WEB-DL 1080p Preferred", SizeBytes: 10},
		{Title: "Movie 2026 Remux 2160p", SizeBytes: 20},
	})
	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.Title != "Movie 2026 Remux 2160p" {
		t.Fatalf("expected quality to win, got %q", decision.Release.Title)
	}
}

func TestChooseReleaseCustomFormatScoresWithinSameQuality(t *testing.T) {
	engine := NewEngine()
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000202")
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: 100},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Preferred group",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "preferred", Name: "Preferred", Type: "releaseTitle", Value: "Preferred", Required: true,
		}},
	}}
	item := storage.MediaItem{Type: "movie", Title: "Movie"}
	decision, ok := engine.ChooseReleaseWithProfile(item, &profile, formats, []storage.ReleaseCandidateInput{
		{Title: "Movie 2026 WEB-DL 1080p", SizeBytes: 10},
		{Title: "Movie 2026 WEB-DL 1080p Preferred", SizeBytes: 20},
	})
	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.Title != "Movie 2026 WEB-DL 1080p Preferred" {
		t.Fatalf("expected custom format to win, got %q", decision.Release.Title)
	}
}
