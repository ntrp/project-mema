package storage

import "testing"

func TestDLNASettingsDefaultsAndNormalization(t *testing.T) {
	input := normalizeDLNASettings(DLNASettingsInput{
		FriendlyName:     "  ",
		Interfaces:       []string{" en0 ", "en0", ""},
		AllowedCIDRs:     nil,
		TranscodeEnabled: true,
	})

	if input.FriendlyName != DefaultDLNAFriendlyName {
		t.Fatalf("friendly name = %q", input.FriendlyName)
	}
	if len(input.Interfaces) != 1 || input.Interfaces[0] != "en0" {
		t.Fatalf("interfaces = %#v", input.Interfaces)
	}
	if len(input.AllowedCIDRs) != len(DefaultDLNAAllowedCIDRs) {
		t.Fatalf("allowed CIDRs = %#v", input.AllowedCIDRs)
	}
	if input.AnnounceIntervalSeconds != DefaultDLNAAnnounceIntervalSeconds {
		t.Fatalf("announce interval = %d", input.AnnounceIntervalSeconds)
	}
	if input.DefaultRendererProfile != DefaultDLNARendererProfile {
		t.Fatalf("renderer profile = %q", input.DefaultRendererProfile)
	}
}

func TestDLNASettingsRejectInvalidCIDRAndInterval(t *testing.T) {
	input := normalizeDLNASettings(DLNASettingsInput{
		FriendlyName:            "Mema",
		AllowedCIDRs:            []string{"not-a-cidr"},
		AnnounceIntervalSeconds: 1800,
		DefaultRendererProfile:  "generic",
	})
	if err := validateDLNASettings(input); err == nil {
		t.Fatal("expected invalid CIDR to fail")
	}

	input.AllowedCIDRs = []string{"192.168.0.0/16"}
	input.AnnounceIntervalSeconds = 10
	if err := validateDLNASettings(input); err == nil {
		t.Fatal("expected short announce interval to fail")
	}
}

func TestDLNASettingsRejectMissingInterface(t *testing.T) {
	input := normalizeDLNASettings(DLNASettingsInput{
		FriendlyName:            "Mema",
		Interfaces:              []string{"mema-missing-interface"},
		AllowedCIDRs:            []string{"192.168.0.0/16"},
		AnnounceIntervalSeconds: 1800,
		DefaultRendererProfile:  "generic",
	})
	if err := validateDLNASettings(input); err == nil {
		t.Fatal("expected missing interface to fail")
	}
}
