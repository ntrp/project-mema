package jobs

import "testing"

func TestMediaRefreshScheduleIsEnabled(t *testing.T) {
	definition, ok := fixedJobDefinitionByID("media_refresh")
	if !ok {
		t.Fatal("media refresh schedule missing")
	}
	if definition.Name != "Media Refresh" || definition.Kind != (MediaRefreshArgs{}).Kind() {
		t.Fatalf("media refresh definition = %#v", definition.SystemJobScheduleDefinition)
	}
	if definition.PausedByDefault || !definition.Automatic || !definition.ManualActionAvailable {
		t.Fatalf("media refresh flags = %#v", definition.SystemJobScheduleDefinition)
	}
}

func TestMediaRefreshProgressPercent(t *testing.T) {
	cases := []struct {
		processed int
		total     int
		want      int32
	}{
		{0, 0, 100},
		{0, 4, 0},
		{1, 4, 25},
		{4, 4, 100},
	}
	for _, tc := range cases {
		got := mediaRefreshProgressPercent(tc.processed, tc.total)
		if got != tc.want {
			t.Fatalf("mediaRefreshProgressPercent(%d, %d) = %d, want %d", tc.processed, tc.total, got, tc.want)
		}
	}
}
