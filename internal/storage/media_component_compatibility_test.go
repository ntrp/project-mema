package storage

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"
)

func TestMediaComponentCompatibilityStates(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	base := retainCompatibilitySource(t, ctx, store, item, "baseVideo", "Base.mkv", 7_200_000, "Theatrical")

	for _, tc := range []struct {
		name       string
		durationMs int
		title      string
		confidence string
		automation string
		review     string
	}{
		{name: "exact", durationMs: 7_200_500, title: "Theatrical", confidence: "exact", automation: "allowed", review: "notRequired"},
		{name: "likely", durationMs: 7_202_500, title: "Theatrical", confidence: "likely", automation: "allowed", review: "notRequired"},
		{name: "runtime mismatch", durationMs: 7_230_000, title: "Theatrical", confidence: "incompatible", automation: "blocked", review: "pending"},
		{name: "cut mismatch", durationMs: 7_200_000, title: "Extended", confidence: "incompatible", automation: "blocked", review: "pending"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			component := retainCompatibilitySource(t, ctx, store, item, "audio", tc.name+".mkv", tc.durationMs, tc.title)
			decision, err := store.EvaluateMediaComponentCompatibility(ctx, item.ID, component.ID, base.ID)
			if err != nil {
				t.Fatal(err)
			}
			if decision.ConfidenceState != tc.confidence ||
				decision.AutomationState != tc.automation ||
				decision.ReviewState != tc.review {
				t.Fatalf("decision = %#v", decision)
			}
		})
	}
}

func TestMediaComponentCompatibilityManualReview(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	base := retainCompatibilitySource(t, ctx, store, item, "baseVideo", "Base.mkv", 7_200_000, "")
	component := retainCompatibilitySource(t, ctx, store, item, "audio", "Candidate.mkv", 7_206_000, "")

	decision, err := store.EvaluateMediaComponentCompatibility(ctx, item.ID, component.ID, base.ID)
	if err != nil {
		t.Fatal(err)
	}
	if decision.ConfidenceState != "uncertain" || decision.AutomationState != "blocked" {
		t.Fatalf("decision = %#v", decision)
	}
	reviewed, err := store.ReviewMediaComponentCompatibility(ctx, item.ID, component.ID, decision.ID, MediaComponentCompatibilityReviewInput{
		ReviewState: "approved",
		Reason:      stringPtr("manual sync check passed"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if reviewed.ReviewState != "approved" || reviewed.AutomationState != "allowed" ||
		reviewed.ReviewReason == nil || *reviewed.ReviewReason != "manual sync check passed" {
		t.Fatalf("reviewed = %#v", reviewed)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.ComponentSources) < 2 {
		t.Fatalf("sources = %#v", loaded.ComponentSources)
	}
	for _, source := range loaded.ComponentSources {
		if source.ID == component.ID && len(source.Compatibility) != 1 {
			t.Fatalf("compatibility not hydrated: %#v", source)
		}
	}
}

func retainCompatibilitySource(
	t *testing.T,
	ctx context.Context,
	store *SettingsStore,
	item MediaItem,
	role string,
	name string,
	durationMs int,
	title string,
) MediaComponentSource {
	t.Helper()
	path := filepath.Join(*item.MediaFolderPath, name)
	writeTestFile(t, path)
	source, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:     role,
		SourceFilePath: path,
		ReleaseTitle:   stringPtr(title),
		StreamInventory: `{"format":{"duration":"` +
			formatDurationSeconds(durationMs) + `"},"streams":[{"id":0,"type":"video"}]}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return source
}

func formatDurationSeconds(durationMs int) string {
	return strconv.FormatFloat(float64(durationMs)/1000, 'f', 3, 64)
}
