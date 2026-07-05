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
