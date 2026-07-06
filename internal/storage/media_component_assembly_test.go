package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestMediaComponentAssemblyRunRequiresAllowedArtifactsAndTracksLifecycle(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	base := retainCompatibilitySource(t, ctx, store, item, "baseVideo", "Base.mkv", 7_200_000, "")
	source := retainCompatibilitySource(t, ctx, store, item, "audio", "Audio.mkv", 7_200_000, "")
	if _, err := store.EvaluateMediaComponentCompatibility(ctx, item.ID, source.ID, base.ID); err != nil {
		t.Fatal(err)
	}
	artifact, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, MediaComponentArtifactInput{
		StreamID:   1,
		StreamType: "audio",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifact.OutputPath, []byte("audio"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err = store.CompleteMediaComponentArtifact(ctx, artifact.ID, "ok")
	if err != nil {
		t.Fatal(err)
	}

	run, err := store.CreateMediaComponentAssemblyRun(ctx, item.ID, MediaComponentAssemblyRunInput{
		BaseSourceID: base.ID,
		ArtifactIDs:  []uuid.UUID{artifact.ID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != "queued" || len(run.Inputs) != 2 || filepath.Base(run.OutputPath) != "assembled.mkv" {
		t.Fatalf("run = %#v", run)
	}
	if run.Inputs[0].StreamType != "video" || run.Inputs[1].Provenance["artifactId"] != artifact.ID.String() {
		t.Fatalf("inputs = %#v", run.Inputs)
	}
	provenance, err := store.ListMediaComponentProvenance(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !componentProvenanceHas(provenance, "container") || !componentProvenanceHas(provenance, "audio") {
		t.Fatalf("provenance = %#v", provenance)
	}
	run, err = store.StartMediaComponentAssemblyRun(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != "running" {
		t.Fatalf("started run = %#v", run)
	}
	if err := os.WriteFile(run.OutputPath, []byte("muxed"), 0o644); err != nil {
		t.Fatal(err)
	}
	run, err = store.CompleteMediaComponentAssemblyRun(ctx, run.ID, "mux ok")
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != "succeeded" || run.SizeBytes == nil || *run.SizeBytes != 5 {
		t.Fatalf("completed run = %#v", run)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.AssemblyRuns) != 1 || len(loaded.AssemblyRuns[0].Inputs) != 2 {
		t.Fatalf("assembly runs not hydrated: %#v", loaded.AssemblyRuns)
	}
	if len(loaded.ComponentProvenance) < 2 {
		t.Fatalf("component provenance not hydrated: %#v", loaded.ComponentProvenance)
	}
}

func TestMediaComponentAssemblyUsesProfileFinalContainer(t *testing.T) {
	ctx, store := testDBStore(t)
	profile, err := store.CreateMediaProfile(ctx, MediaProfileInput{
		Name:            "MP4 Assembly",
		FinalContainer:  "mp4",
		UpgradesAllowed: true,
		QualityIDs:      []string{"webdl-1080p"},
		AudioTargets:    []MediaProfileAudioTarget{{LanguageID: "english"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:             "movie",
		Title:            "MP4 Assembly " + uuid.NewString(),
		Monitored:        true,
		QualityProfileID: &profile.ID,
		LibraryFolderID:  &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	base := retainCompatibilitySource(t, ctx, store, item, "baseVideo", "Base.mkv", 7_200_000, "")
	source := retainCompatibilitySource(t, ctx, store, item, "audio", "Audio.mkv", 7_200_000, "")
	if _, err := store.EvaluateMediaComponentCompatibility(ctx, item.ID, source.ID, base.ID); err != nil {
		t.Fatal(err)
	}
	artifact, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, MediaComponentArtifactInput{
		StreamID:   1,
		StreamType: "audio",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifact.OutputPath, []byte("audio"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err = store.CompleteMediaComponentArtifact(ctx, artifact.ID, "ok")
	if err != nil {
		t.Fatal(err)
	}

	run, err := store.CreateMediaComponentAssemblyRun(ctx, item.ID, MediaComponentAssemblyRunInput{
		BaseSourceID: base.ID,
		ArtifactIDs:  []uuid.UUID{artifact.ID},
	})

	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(run.OutputPath) != ".mp4" {
		t.Fatalf("output path = %q", run.OutputPath)
	}
}

func TestMediaComponentAssemblyRejectsBlockedCompatibility(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)
	base := retainCompatibilitySource(t, ctx, store, item, "baseVideo", "Base.mkv", 7_200_000, "")
	source := retainCompatibilitySource(t, ctx, store, item, "audio", "Audio.mkv", 7_230_000, "")
	if _, err := store.EvaluateMediaComponentCompatibility(ctx, item.ID, source.ID, base.ID); err != nil {
		t.Fatal(err)
	}
	artifact, err := store.CreateMediaComponentArtifact(ctx, item.ID, source.ID, MediaComponentArtifactInput{
		StreamID:   1,
		StreamType: "audio",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifact.OutputPath, []byte("audio"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err = store.CompleteMediaComponentArtifact(ctx, artifact.ID, "ok")
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.CreateMediaComponentAssemblyRun(ctx, item.ID, MediaComponentAssemblyRunInput{
		BaseSourceID: base.ID,
		ArtifactIDs:  []uuid.UUID{artifact.ID},
	})
	if err != ErrInvalidInput {
		t.Fatalf("error = %v, want invalid", err)
	}
}

func componentProvenanceHas(values []MediaComponentProvenance, componentType string) bool {
	for _, value := range values {
		if value.ComponentType == componentType && len(value.TransformationChain) > 0 {
			return true
		}
	}
	return false
}
