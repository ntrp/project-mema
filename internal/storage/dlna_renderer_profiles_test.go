package storage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/testdb"
)

func TestDLNARendererProfilesSeedBaseline(t *testing.T) {
	ctx, store := testDBStore(t)

	profiles, err := store.ListDLNARendererProfiles(ctx)
	if err != nil {
		t.Fatalf("list renderer profiles: %v", err)
	}
	wantIDs := []string{
		"generic", "vlc", "bubbleupnp", "chromecast", "kodi",
		"samsung-tv", "sony-tv", "lg-webos", "lg-tv-2023", "lg-tv-2025",
	}
	for _, id := range wantIDs {
		if !hasDLNARendererProfile(profiles, id) {
			t.Fatalf("seeded renderer profile %q missing in %#v", id, profileIDs(profiles))
		}
	}
	if len(profiles) < len(wantIDs) {
		t.Fatalf("profile count = %d, want at least %d", len(profiles), len(wantIDs))
	}
}

func TestDLNARendererProfileEditSurvivesSeedRestart(t *testing.T) {
	ctx := context.Background()
	store, databaseURL, closeStore := newRendererProfileTestStore(t, ctx)

	edited := rendererProfileInput(requireRendererProfile(t, ctx, store, "vlc"))
	edited.Name = "VLC Custom"
	edited.Notes = "user tuned"
	edited.Priority = 12
	edited.MatchRules = rawJSON(`{"mode":"weighted","minScore":9,"tokens":[]}`)
	if _, err := store.UpdateDLNARendererProfile(ctx, "vlc", edited); err != nil {
		t.Fatalf("update renderer profile: %v", err)
	}
	closeStore()

	if err := EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatalf("ensure schema after edit: %v", err)
	}
	store, _, closeStore = newRendererProfileStoreFromURL(t, ctx, databaseURL)
	defer closeStore()
	reloaded := requireRendererProfile(t, ctx, store, "vlc")
	if reloaded.Name != "VLC Custom" || reloaded.Notes != "user tuned" || reloaded.Priority != 12 {
		t.Fatalf("profile overwritten after seed: %#v", reloaded)
	}
	if !reloaded.Customized || reloaded.Source != "user" {
		t.Fatalf("custom flags = customized %t source %q", reloaded.Customized, reloaded.Source)
	}
}

func TestDLNARendererProfileResetRestoresDefault(t *testing.T) {
	ctx, store := testDBStore(t)

	input := rendererProfileInput(requireRendererProfile(t, ctx, store, "chromecast"))
	input.Name = "Custom Cast"
	input.DeliverySettings = rawJSON(`{"preferHls":false,"avoidHls":true}`)
	if _, err := store.UpdateDLNARendererProfile(ctx, "chromecast", input); err != nil {
		t.Fatalf("update renderer profile: %v", err)
	}
	reset, err := store.ResetDLNARendererProfile(ctx, "chromecast")
	if err != nil {
		t.Fatalf("reset renderer profile: %v", err)
	}
	if reset.Name != "Chromecast" || reset.Customized || reset.Source != "mema_seed" {
		t.Fatalf("reset profile = %#v", reset)
	}
	if !json.Valid(reset.DeliverySettings) {
		t.Fatalf("reset delivery settings invalid: %s", reset.DeliverySettings)
	}
}

func TestDLNARendererDeviceOverridesAssignByIPOrUUID(t *testing.T) {
	ctx, store := testDBStore(t)
	ip := "192.0.2.44"
	rendererUUID := "uuid:living-room-tv"

	first, err := store.UpsertDLNARendererDeviceOverride(ctx, DLNARendererDeviceOverrideInput{
		IPAddress:               &ip,
		ProfileID:               "vlc",
		DisplayName:             "Desk VLC",
		Allowed:                 true,
		DeliveryPolicyOverrides: rawJSON(`{"preferHls":false}`),
	})
	if err != nil {
		t.Fatalf("upsert ip override: %v", err)
	}
	if _, err := store.UpsertDLNARendererDeviceOverride(ctx, DLNARendererDeviceOverrideInput{
		RendererUUID:            &rendererUUID,
		ProfileID:               "lg-tv-2025",
		DisplayName:             "Living Room",
		Allowed:                 true,
		DeliveryPolicyOverrides: rawJSON(`{"avoidHls":true}`),
	}); err != nil {
		t.Fatalf("upsert uuid override: %v", err)
	}
	overrides, err := store.ListDLNARendererDeviceOverrides(ctx)
	if err != nil {
		t.Fatalf("list overrides: %v", err)
	}
	if !hasOverride(overrides, "vlc", &ip, nil) {
		t.Fatalf("ip override missing: %#v", overrides)
	}
	if !hasOverride(overrides, "lg-tv-2025", nil, &rendererUUID) {
		t.Fatalf("uuid override missing: %#v", overrides)
	}
	if err := store.DeleteDLNARendererDeviceOverride(ctx, first.ID); err != nil {
		t.Fatalf("delete override: %v", err)
	}
}

func newRendererProfileTestStore(
	t *testing.T,
	ctx context.Context,
) (*SettingsStore, string, func()) {
	t.Helper()
	databaseURL := testdb.Create(t)
	if err := EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	store, _, closeStore := newRendererProfileStoreFromURL(t, ctx, databaseURL)
	return store, databaseURL, closeStore
}

func newRendererProfileStoreFromURL(
	t *testing.T,
	ctx context.Context,
	databaseURL string,
) (*SettingsStore, string, func()) {
	t.Helper()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	return NewSettingsStore(pool), databaseURL, pool.Close
}

func requireRendererProfile(
	t *testing.T,
	ctx context.Context,
	store *SettingsStore,
	id string,
) DLNARendererProfile {
	t.Helper()
	profile, err := store.GetDLNARendererProfile(ctx, id)
	if err != nil {
		t.Fatalf("get renderer profile %s: %v", id, err)
	}
	return profile
}

func rendererProfileInput(profile DLNARendererProfile) DLNARendererProfileInput {
	return DLNARendererProfileInput{
		Name: profile.Name, Vendor: profile.Vendor, DeviceClass: profile.DeviceClass,
		Enabled: profile.Enabled, Priority: profile.Priority, IconKey: profile.IconKey,
		Notes: profile.Notes, MatchRules: profile.MatchRules,
		CapabilityRules: profile.CapabilityRules, DeliverySettings: profile.DeliverySettings,
		DLNAFlags: profile.DLNAFlags, SubtitleRules: profile.SubtitleRules,
		ArtworkRules: profile.ArtworkRules, MetadataRules: profile.MetadataRules, Quirks: profile.Quirks,
	}
}

func rawJSON(value string) json.RawMessage {
	return json.RawMessage(value)
}

func hasDLNARendererProfile(profiles []DLNARendererProfile, id string) bool {
	for _, profile := range profiles {
		if profile.ID == id {
			return true
		}
	}
	return false
}

func profileIDs(profiles []DLNARendererProfile) []string {
	ids := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		ids = append(ids, profile.ID)
	}
	return ids
}

func hasOverride(
	overrides []DLNARendererDeviceOverride,
	profileID string,
	ip *string,
	rendererUUID *string,
) bool {
	for _, override := range overrides {
		if override.ProfileID != profileID {
			continue
		}
		if ip != nil && override.IPAddress != nil && *override.IPAddress == *ip {
			return true
		}
		if rendererUUID != nil && override.RendererUUID != nil && *override.RendererUUID == *rendererUUID {
			return true
		}
	}
	return false
}
