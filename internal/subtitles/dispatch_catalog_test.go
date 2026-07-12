package subtitles

import (
	"testing"

	"media-manager/internal/subtitles/catalog"
)

func TestCatalogProvidersAreRuntimeDispatchable(t *testing.T) {
	entries := catalog.MustAll()
	if len(entries) != 59 {
		t.Fatalf("catalog entries = %d, want 59", len(entries))
	}
	for _, entry := range entries {
		if !RuntimeSupported(entry.Key) {
			t.Fatalf("catalog provider %q is not runtime-dispatchable: %v", entry.Key, UnsupportedRuntimeError(entry.Key))
		}
	}
}
