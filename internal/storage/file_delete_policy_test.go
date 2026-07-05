package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileDeleteSettingsValidateModesAndRecycleFolder(t *testing.T) {
	ctx, store := testDBStore(t)
	for _, mode := range []string{FileDeleteModePermanent, FileDeleteModeRecycle, FileDeleteModeKeep} {
		if _, err := store.UpdateFileDeleteSettings(ctx, FileDeleteSettingsInput{
			Mode:          mode,
			RecycleFolder: ".trash",
		}); err != nil {
			t.Fatalf("mode %q rejected: %v", mode, err)
		}
	}
	for _, input := range []FileDeleteSettingsInput{
		{Mode: "unknown", RecycleFolder: ".trash"},
		{Mode: FileDeleteModeRecycle, RecycleFolder: "../trash"},
		{Mode: FileDeleteModeRecycle, RecycleFolder: "/tmp/trash"},
		{Mode: FileDeleteModeRecycle, RecycleFolder: "trash"},
	} {
		if _, err := store.UpdateFileDeleteSettings(ctx, input); err == nil {
			t.Fatalf("expected %#v to be rejected", input)
		}
	}
}

func TestDeleteMediaItemFilePermanentPolicyRemovesFile(t *testing.T) {
	ctx, store := testDBStore(t)
	item, path := deletePolicyItem(t, ctx, store)

	if _, err := store.UpdateFileDeleteSettings(ctx, FileDeleteSettingsInput{Mode: FileDeleteModePermanent}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DeleteMediaItemFile(ctx, item.ID, path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("file stat err = %v, want missing", err)
	}
	history, err := store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !historyHasDeleteMode(history, FileDeleteModePermanent, "succeeded") {
		t.Fatalf("history = %#v", history)
	}
}

func TestDeleteMediaItemFileRecyclePolicyMovesFile(t *testing.T) {
	ctx, store := testDBStore(t)
	item, path := deletePolicyItem(t, ctx, store)

	if _, err := store.UpdateFileDeleteSettings(ctx, FileDeleteSettingsInput{
		Mode:          FileDeleteModeRecycle,
		RecycleFolder: ".trash",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DeleteMediaItemFile(ctx, item.ID, path); err != nil {
		t.Fatal(err)
	}
	recycled := filepath.Join(*item.LibraryFolderPath, ".trash", "Policy Movie (2026)", filepath.Base(path))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("source stat err = %v, want missing", err)
	}
	if _, err := os.Stat(recycled); err != nil {
		t.Fatalf("recycled file missing: %v", err)
	}
	history, err := store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !historyHasDeleteMode(history, FileDeleteModeRecycle, "succeeded") {
		t.Fatalf("history = %#v", history)
	}
}

func TestDeleteMediaItemFileKeepPolicySkipsRemoval(t *testing.T) {
	ctx, store := testDBStore(t)
	item, path := deletePolicyItem(t, ctx, store)

	if _, err := store.UpdateFileDeleteSettings(ctx, FileDeleteSettingsInput{Mode: FileDeleteModeKeep}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DeleteMediaItemFile(ctx, item.ID, path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should remain: %v", err)
	}
	history, err := store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !historyHasDeleteMode(history, FileDeleteModeKeep, "skipped") {
		t.Fatalf("history = %#v", history)
	}
}

func historyHasDeleteMode(history []MediaFileHistoryEntry, mode string, status string) bool {
	for _, entry := range history {
		if entry.Operation != "deleted" || entry.Status != status {
			continue
		}
		if entry.Details["deleteMode"] == mode {
			return true
		}
	}
	return false
}
