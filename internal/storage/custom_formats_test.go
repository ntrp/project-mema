package storage

import "testing"

func TestDefaultCustomFormatsLoad(t *testing.T) {
	formats, err := defaultCustomFormats()
	if err != nil {
		t.Fatalf("load default custom formats: %v", err)
	}
	if len(formats) < 300 {
		t.Fatalf("expected TRaSH custom format seed set, got %d", len(formats))
	}

	duplicateNames := map[string]int{}
	for _, format := range formats {
		if format.ID.String() == "00000000-0000-0000-0000-000000000000" {
			t.Fatalf("format %q has empty id", format.Name)
		}
		if format.Name == "" {
			t.Fatal("format has empty name")
		}
		if len(format.Name) >= 7 && (format.Name[:7] == "Radarr " || format.Name[:7] == "Sonarr ") {
			t.Fatalf("format %q should not expose source app prefix", format.Name)
		}
		duplicateNames[format.Name]++
		if len(format.IncludeSpecs) == 0 && len(format.ExcludeSpecs) == 0 {
			t.Fatalf("format %q has no specs", format.Name)
		}
		for _, spec := range append(format.IncludeSpecs, format.ExcludeSpecs...) {
			if spec.ID == "" || spec.Name == "" || spec.Type == "" || spec.Value == "" {
				t.Fatalf("format %q has invalid spec %#v", format.Name, spec)
			}
		}
	}
	if duplicateNames["WEB Tier 01"] < 2 {
		t.Fatal("expected duplicate visible names for Radarr/Sonarr variants")
	}
}
