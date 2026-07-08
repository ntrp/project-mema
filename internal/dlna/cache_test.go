package dlna

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPruneCacheDirRemovesOldAndOversizedEntries(t *testing.T) {
	dir := t.TempDir()
	now := time.Unix(1000, 0)
	old := writeCacheFile(t, dir, "old.bin", 6, now.Add(-48*time.Hour))
	newer := writeCacheFile(t, dir, "newer.bin", 6, now.Add(-2*time.Hour))
	newest := writeCacheFile(t, dir, "newest.bin", 6, now.Add(-time.Hour))

	result, err := pruneCacheDir(dir, now, 24*time.Hour, 8)
	if err != nil {
		t.Fatal(err)
	}
	if result.DeletedFiles != 2 || result.DeletedBytes != 12 {
		t.Fatalf("result = %#v", result)
	}
	if fileExists(old) || fileExists(newer) || !fileExists(newest) {
		t.Fatalf("cache files old=%t newer=%t newest=%t", fileExists(old), fileExists(newer), fileExists(newest))
	}
}

func TestManagerPruneCachesCoversThumbnailAndOutputDirs(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.thumbDir = t.TempDir()
	manager.remuxDir = t.TempDir()
	now := time.Unix(2000, 0)
	writeCacheFile(t, manager.thumbDir, "thumb.jpg", 3, now.Add(-30*24*time.Hour))
	writeCacheFile(t, manager.remuxDir, "stream.ts", 5, now.Add(-10*24*time.Hour))

	result, err := manager.PruneCaches(now)
	if err != nil {
		t.Fatal(err)
	}
	if result.DeletedFiles != 2 || result.DeletedBytes != 8 {
		t.Fatalf("result = %#v", result)
	}
}

func writeCacheFile(t *testing.T, dir string, name string, size int, modTime time.Time) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatal(err)
	}
	return path
}
