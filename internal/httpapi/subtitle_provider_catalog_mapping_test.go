package httpapi

import (
	"testing"

	"media-manager/internal/subtitles/catalog"
)

func TestSubtitleProviderCatalogMapping(t *testing.T) {
	entry, ok := catalog.Lookup("opensubtitlescom")
	if !ok {
		t.Fatal("opensubtitlescom catalog entry missing")
	}
	mapped := subtitleProviderCatalogEntry(entry)
	if mapped.Key != "opensubtitlescom" || mapped.RuntimeStatus != Supported {
		t.Fatalf("unexpected mapped catalog entry: %#v", mapped)
	}
	if mapped.ProvenanceCommit == nil || *mapped.ProvenanceCommit == "" {
		t.Fatalf("expected provenance commit: %#v", mapped.ProvenanceCommit)
	}
	if len(mapped.Fields) == 0 {
		t.Fatalf("expected field definitions")
	}
	if mapped.OutboundPolicy.AllowedBaseHosts == nil || len(*mapped.OutboundPolicy.AllowedBaseHosts) == 0 {
		t.Fatalf("expected outbound policy hosts")
	}
}
