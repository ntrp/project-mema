package storage

import "testing"

func TestQualitySizeSettingsUseGeneratedQueriesWithoutResettingUpdates(t *testing.T) {
	ctx, store := testDBStore(t)

	preferred := 2.5
	maximum := 4.75
	updated, err := store.SaveQualitySizeSettings(ctx, []QualitySizeSettingInput{{
		QualityID:                "unknown",
		MinimumSizeMBPerMinute:   1.25,
		PreferredSizeMBPerMinute: &preferred,
		MaximumSizeMBPerMinute:   &maximum,
	}})
	if err != nil {
		t.Fatal(err)
	}
	requireQualitySizeMinimum(t, updated, "unknown", 1.25)

	listed, err := store.ListQualitySizeSettings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	requireQualitySizeMinimum(t, listed, "unknown", 1.25)
}

func TestQualitySizeDefinitionsUseSeededDefaults(t *testing.T) {
	definitions := QualitySizeDefinitionMap()
	requireQualitySizeDefault(t, definitions, "unknown", 0, 45, 180)
	requireQualitySizeDefault(t, definitions, "sdtv", 4, 10, 22)
	requireQualitySizeDefault(t, definitions, "webdl-1080p", 24, 52, 96)
	requireQualitySizeDefault(t, definitions, "bluray-1080p", 36, 75, 135)
	requireQualitySizeDefault(t, definitions, "remux-2160p", 170, 330, 620)
	requireQualitySizeDefault(t, definitions, "br-disk", 200, 400, 760)
	requireQualitySizeSourceOrder(t, definitions, "720p")
	requireQualitySizeSourceOrder(t, definitions, "1080p")
	requireQualitySizeSourceOrder(t, definitions, "2160p")
	requireQualitySizeAbove(t, definitions, "br-disk", "remux-2160p")
	requireQualitySizeAbove(t, definitions, "raw-hd", "br-disk")
}

func TestQualityResolutionForIDUsesQualityResolution(t *testing.T) {
	bounds, ok := QualityResolutionForID("webdl-1080p")

	if !ok {
		t.Fatalf("expected resolution bounds")
	}
	if bounds.MinWidth != 1920 || bounds.MinHeight != 1080 {
		t.Fatalf("bounds = %#v", bounds)
	}
}

func TestQualityIDFromPathUsesBestMatchingQuality(t *testing.T) {
	got := QualityIDFromPath("/media/Movie.2020.1080p.WEBDL.mkv")

	if got != "webdl-1080p" {
		t.Fatalf("quality id = %q", got)
	}
}

func requireQualitySizeMinimum(
	t *testing.T,
	settings []QualitySizeSetting,
	qualityID string,
	minimum float64,
) {
	t.Helper()
	for _, setting := range settings {
		if setting.ID == qualityID {
			if setting.MinimumSizeMBPerMinute != minimum {
				t.Fatalf("%s minimum = %v, want %v", qualityID, setting.MinimumSizeMBPerMinute, minimum)
			}
			return
		}
	}
	t.Fatalf("%s missing from quality size settings", qualityID)
}

func requireQualitySizeDefault(
	t *testing.T,
	definitions map[string]QualitySizeDefinition,
	qualityID string,
	minimum float64,
	preferred float64,
	maximum float64,
) {
	t.Helper()
	definition, ok := definitions[qualityID]
	if !ok {
		t.Fatalf("%s missing from quality size definitions", qualityID)
	}
	if definition.DefaultMinimumSizeMBPerMinute != minimum {
		t.Fatalf("%s minimum = %v, want %v", qualityID, definition.DefaultMinimumSizeMBPerMinute, minimum)
	}
	if definition.DefaultPreferredSizeMBPerMinute == nil || *definition.DefaultPreferredSizeMBPerMinute != preferred {
		t.Fatalf("%s preferred = %v, want %v", qualityID, definition.DefaultPreferredSizeMBPerMinute, preferred)
	}
	if definition.DefaultMaximumSizeMBPerMinute == nil || *definition.DefaultMaximumSizeMBPerMinute != maximum {
		t.Fatalf("%s maximum = %v, want %v", qualityID, definition.DefaultMaximumSizeMBPerMinute, maximum)
	}
}

func requireQualitySizeSourceOrder(
	t *testing.T,
	definitions map[string]QualitySizeDefinition,
	resolution string,
) {
	t.Helper()
	requireQualitySizeAbove(t, definitions, "webrip-"+resolution, "hdtv-"+resolution)
	requireQualitySizeAbove(t, definitions, "webdl-"+resolution, "webrip-"+resolution)
	requireQualitySizeAbove(t, definitions, "bluray-"+resolution, "webdl-"+resolution)
}

func requireQualitySizeAbove(
	t *testing.T,
	definitions map[string]QualitySizeDefinition,
	higherID string,
	lowerID string,
) {
	t.Helper()
	higher := definitions[higherID]
	lower := definitions[lowerID]
	if higher.DefaultPreferredSizeMBPerMinute == nil || lower.DefaultPreferredSizeMBPerMinute == nil ||
		*higher.DefaultPreferredSizeMBPerMinute <= *lower.DefaultPreferredSizeMBPerMinute {
		t.Fatalf("%s preferred should be above %s", higherID, lowerID)
	}
	if higher.DefaultMaximumSizeMBPerMinute == nil || lower.DefaultMaximumSizeMBPerMinute == nil ||
		*higher.DefaultMaximumSizeMBPerMinute <= *lower.DefaultMaximumSizeMBPerMinute {
		t.Fatalf("%s maximum should be above %s", higherID, lowerID)
	}
}
