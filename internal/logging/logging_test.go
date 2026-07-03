package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestSCNSystem007LogManagerPublishesStructuredEntries(t *testing.T) {
	manager := NewManager()
	if err := manager.SetLevel(LevelDebug); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	logger := slog.New(manager.Handler(&output))
	entries, unsubscribe := manager.Subscribe()
	defer unsubscribe()

	logger.InfoContext(context.Background(), "scenario event",
		"component", "tests",
		"count", 2,
		"elapsed", time.Second,
		"group", slog.GroupValue(slog.String("child", "value")),
	)

	entry := <-entries
	if entry.Message != "scenario event" || entry.Level != LevelInfo {
		t.Fatalf("entry = %#v", entry)
	}
	if entry.Attributes["component"] != "tests" || entry.Attributes["count"] != int64(2) {
		t.Fatalf("attributes = %#v", entry.Attributes)
	}
	group, ok := entry.Attributes["group"].(map[string]any)
	if !ok || group["child"] != "value" {
		t.Fatalf("group = %#v", entry.Attributes["group"])
	}
	if !strings.Contains(output.String(), "scenario event") {
		t.Fatalf("handler output = %q", output.String())
	}

	buffered, closeBuffered := manager.Subscribe()
	defer closeBuffered()
	if got := <-buffered; got.ID != entry.ID {
		t.Fatalf("buffered entry = %#v, want id %s", got, entry.ID)
	}
}

func TestSCNSystem007LogFilesAreWrittenAndListed(t *testing.T) {
	manager := NewManager()
	directory := t.TempDir()
	if err := manager.ConfigureFile(FileSettings{
		Enabled:       true,
		Directory:     directory,
		RetentionDays: 7,
	}); err != nil {
		t.Fatal(err)
	}
	settings := manager.LogFileSettings()
	if !settings.Enabled || settings.Directory != directory || settings.RetentionDays != 7 {
		t.Fatalf("settings = %#v", settings)
	}

	logger := slog.New(manager.Handler(&bytes.Buffer{}))
	logger.Warn("file scenario", "kind", "unit")

	files, err := ListFiles(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].SizeBytes == 0 {
		t.Fatalf("files = %#v", files)
	}
	path, ok := FilePath(directory, files[0].Name)
	if !ok || path != files[0].Path {
		t.Fatalf("safe path = %q, %v, file = %#v", path, ok, files[0])
	}
	if _, ok := FilePath(directory, "../"+files[0].Name); ok {
		t.Fatal("path traversal should be rejected")
	}
	if _, ok := FilePath(directory, "notes.txt"); ok {
		t.Fatal("non-log file should be rejected")
	}
}
