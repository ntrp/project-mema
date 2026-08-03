package metadata

import (
	"path/filepath"
	"testing"
)

func TestAniDBArtworkFilePathAcceptsOnlyNumericJPEGNames(t *testing.T) {
	root := t.TempDir()
	path, ok := AniDBArtworkFilePath(root, "221067.jpg")
	if !ok || path != filepath.Join(root, ".mema", "anidb", "artwork", "221067.jpg") {
		t.Fatalf("path=%q ok=%v", path, ok)
	}
	for _, name := range []string{"../secret.jpg", "actor.jpg", "221067.png", ""} {
		if path, ok := AniDBArtworkFilePath(root, name); ok || path != "" {
			t.Fatalf("unsafe name %q produced %q", name, path)
		}
	}
}
