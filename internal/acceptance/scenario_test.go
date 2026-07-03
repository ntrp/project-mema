package acceptance

import "testing"

func TestRequireScenarioFindsTaggedScenario(t *testing.T) {
	scenario, err := RequireScenario("features/behavior", "SCN-AUTH-001")
	if err != nil {
		t.Fatal(err)
	}
	if scenario.Feature != "Authentication" {
		t.Fatalf("feature = %q", scenario.Feature)
	}
	if !scenario.HasTag("api") || !scenario.HasTag("@e2e") {
		t.Fatalf("tags = %#v", scenario.Tags)
	}
	if len(scenario.Steps) == 0 {
		t.Fatal("expected parsed steps")
	}
}
