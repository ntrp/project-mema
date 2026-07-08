package storage

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestDLNARendererSeedCatalogCoverage(t *testing.T) {
	ctx, store := testDBStore(t)
	profiles, err := store.ListDLNARendererProfiles(ctx)
	if err != nil {
		t.Fatalf("list renderer profiles: %v", err)
	}
	wantIDs := []string{
		"amazon-fire-tv-vimu", "android-generic", "chromecast-ultra",
		"apple-ios", "vlc-appletv", "windows-media-player", "xbox-one",
		"panasonic-tv", "philips-tv", "roku-tv", "roku-ultra",
		"samsung-tv-modern", "samsung-neo-qled", "sony-bravia-modern",
		"sony-playstation", "lg-bluray", "lg-oled-2022", "yamaha-av-receiver",
	}
	for _, id := range wantIDs {
		if !hasDLNARendererProfile(profiles, id) {
			t.Fatalf("seeded renderer profile %q missing in %#v", id, profileIDs(profiles))
		}
	}
	if len(profiles) < 80 {
		t.Fatalf("profile count = %d, want full family catalog", len(profiles))
	}
}

func TestDLNARendererSeedCatalogRulesAreCleanRoomData(t *testing.T) {
	seedSQL, err := os.ReadFile("seeds/defaults.sql")
	if err != nil {
		t.Fatalf("read defaults seed: %v", err)
	}
	for _, forbidden := range []string{
		"Supported =", "UserAgentSearch", "TranscodeVideo", "RendererName",
		"MimeTypesChanges", "DLNALocalizationRequired",
	} {
		if strings.Contains(string(seedSQL), forbidden) {
			t.Fatalf("seed contains UMS-style assignment token %q", forbidden)
		}
	}

	ctx, store := testDBStore(t)
	profiles, err := store.ListDLNARendererProfiles(ctx)
	if err != nil {
		t.Fatalf("list renderer profiles: %v", err)
	}
	for _, profile := range profiles {
		requireProfileJSON(t, profile)
		if profile.ID != "generic" && tokenCount(t, profile.MatchRules) == 0 {
			t.Fatalf("profile %s has no match tokens", profile.ID)
		}
	}
}

func requireProfileJSON(t *testing.T, profile DLNARendererProfile) {
	t.Helper()
	payloads := [][]byte{
		profile.MatchRules, profile.CapabilityRules, profile.DeliverySettings,
		profile.DLNAFlags, profile.SubtitleRules, profile.ArtworkRules,
		profile.MetadataRules, profile.Quirks,
	}
	for _, payload := range payloads {
		if !json.Valid(payload) {
			t.Fatalf("profile %s has invalid json: %s", profile.ID, payload)
		}
	}
}

func tokenCount(t *testing.T, payload []byte) int {
	t.Helper()
	var parsed struct {
		Tokens []struct{} `json:"tokens"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		t.Fatalf("parse match rules: %v", err)
	}
	return len(parsed.Tokens)
}
