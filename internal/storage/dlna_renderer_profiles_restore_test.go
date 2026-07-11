package storage

import "testing"

func TestDLNARendererProfileDeleteIsBlockedForSeededProfiles(t *testing.T) {
	ctx, store := testDBStore(t)

	if err := store.DeleteDLNARendererProfile(ctx, "generic"); err != ErrInvalidInput {
		t.Fatalf("delete seeded profile error = %v, want %v", err, ErrInvalidInput)
	}

	created, err := store.CreateDLNARendererProfile(ctx, "user-profile", DLNARendererProfileInput{
		Name:             "User Profile",
		Vendor:           "Mema",
		DeviceClass:      "software",
		Enabled:          true,
		Priority:         50,
		IconKey:          "device",
		Notes:            "temporary",
		MatchRules:       rawJSON(`{}`),
		CapabilityRules:  rawJSON(`{}`),
		DeliverySettings: rawJSON(`{}`),
		DLNAFlags:        rawJSON(`{}`),
		SubtitleRules:    rawJSON(`{}`),
		ArtworkRules:     rawJSON(`{}`),
		MetadataRules:    rawJSON(`{}`),
		Quirks:           rawJSON(`{}`),
	})
	if err != nil {
		t.Fatalf("create user profile: %v", err)
	}
	if err := store.DeleteDLNARendererProfile(ctx, created.ID); err != nil {
		t.Fatalf("delete user profile: %v", err)
	}
}

func TestDLNARendererProfileRestoreAllRestoresSeededProfiles(t *testing.T) {
	ctx, store := testDBStore(t)

	edited := rendererProfileInput(requireRendererProfile(t, ctx, store, "vlc"))
	edited.Name = "VLC Custom"
	edited.Notes = "user tuned"
	if _, err := store.UpdateDLNARendererProfile(ctx, "vlc", edited); err != nil {
		t.Fatalf("update renderer profile: %v", err)
	}
	if _, err := store.CreateDLNARendererProfile(ctx, "user-profile", DLNARendererProfileInput{
		Name:             "User Profile",
		Vendor:           "Mema",
		DeviceClass:      "software",
		Enabled:          true,
		Priority:         50,
		IconKey:          "device",
		Notes:            "temporary",
		MatchRules:       rawJSON(`{}`),
		CapabilityRules:  rawJSON(`{}`),
		DeliverySettings: rawJSON(`{}`),
		DLNAFlags:        rawJSON(`{}`),
		SubtitleRules:    rawJSON(`{}`),
		ArtworkRules:     rawJSON(`{}`),
		MetadataRules:    rawJSON(`{}`),
		Quirks:           rawJSON(`{}`),
	}); err != nil {
		t.Fatalf("create user profile: %v", err)
	}

	if err := store.RestoreDLNARendererProfiles(ctx); err != nil {
		t.Fatalf("restore all renderer profiles: %v", err)
	}

	restored := requireRendererProfile(t, ctx, store, "vlc")
	if restored.Name != "VLC" || restored.Notes == "user tuned" || restored.Customized {
		t.Fatalf("restored seeded profile = %#v", restored)
	}
	if _, err := store.GetDLNARendererProfile(ctx, "user-profile"); err != nil {
		t.Fatalf("user profile should remain after restore: %v", err)
	}
}
