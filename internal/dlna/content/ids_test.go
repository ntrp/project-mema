package content

import "testing"

func TestRootContainerIDsUseStableUPnPStyleValues(t *testing.T) {
	id := EncodeID(RootContainerRef("movies"))
	if id != "0$1" {
		t.Fatalf("root container id = %q", id)
	}
	ref, err := DecodeID(id)
	if err != nil {
		t.Fatalf("DecodeID returned error: %v", err)
	}
	if ref.Kind != "root" || ref.Key != "movies" {
		t.Fatalf("ref = %#v", ref)
	}
}
