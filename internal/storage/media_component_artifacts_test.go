package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMediaComponentArtifactsTrackExtractionLifecycle(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	sourcePath := filepath.Join(*item.MediaFolderPath, "Artifact.Source.mkv")
	writeTestFile(t, sourcePath)
	source, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:     "audio",
		SourceFilePath: sourcePath,
		StreamInventory: `{"streams":[{"id":1,"type":"audio","language":"eng"},` +
			`{"index":2,"codec_type":"subtitle","language":"jpn"}]}`,
	})
	if err != nil {
		t.Fatal(err)
	}

	artifact, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, MediaComponentArtifactInput{
		StreamID:   1,
		StreamType: "audio",
		Language:   stringPtr("eng"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Status != "queued" || filepath.Ext(artifact.OutputPath) != ".mka" {
		t.Fatalf("queued artifact = %#v", artifact)
	}
	artifact, err = store.AssignMediaComponentArtifactJob(ctx, artifact.ID, "42")
	if err != nil {
		t.Fatal(err)
	}
	if artifact.JobID == nil || *artifact.JobID != "42" {
		t.Fatalf("job id not assigned: %#v", artifact)
	}
	artifact, err = store.StartMediaComponentArtifact(ctx, artifact.ID)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Status != "running" {
		t.Fatalf("started artifact = %#v", artifact)
	}
	if err := os.WriteFile(artifact.OutputPath, []byte("extracted"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err = store.CompleteMediaComponentArtifact(ctx, artifact.ID, "track extracted")
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Status != "succeeded" || artifact.SizeBytes == nil || *artifact.SizeBytes != 9 {
		t.Fatalf("completed artifact = %#v", artifact)
	}

	subtitle, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, MediaComponentArtifactInput{
		StreamID:   2,
		StreamType: "subtitle",
		Language:   stringPtr("jpn"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(subtitle.OutputPath) != ".mks" {
		t.Fatalf("subtitle artifact path = %s", subtitle.OutputPath)
	}
	failed, err := store.FailMediaComponentArtifact(ctx, subtitle.ID, "stderr", "tool failed")
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != "failed" || failed.ErrorMessage == nil || *failed.ErrorMessage != "tool failed" {
		t.Fatalf("failed artifact = %#v", failed)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.ComponentSources) != 1 || len(loaded.ComponentSources[0].Artifacts) != 2 {
		t.Fatalf("hydrated sources = %#v", loaded.ComponentSources)
	}
}

func TestMediaComponentArtifactRejectsUnsafeStreamSelection(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	sourcePath := filepath.Join(*item.MediaFolderPath, "Unsafe.Source.mkv")
	writeTestFile(t, sourcePath)
	source, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:      "audio",
		SourceFilePath:  sourcePath,
		StreamInventory: `{"streams":[{"id":1,"type":"audio","language":"eng"}]}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range []MediaComponentArtifactInput{
		{StreamID: -1, StreamType: "audio"},
		{StreamID: 1, StreamType: "video"},
		{StreamID: 99, StreamType: "audio"},
		{StreamID: 1, StreamType: "audio", Language: stringPtr("deu")},
	} {
		_, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, input)
		if !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("input %#v error = %v, want invalid", input, err)
		}
	}
}
