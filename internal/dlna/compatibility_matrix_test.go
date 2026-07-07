package dlna

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestCompatibilityMatrixProfiles(t *testing.T) {
	tests := []struct {
		client     string
		userAgent  string
		profileID  string
		wantFirst  string
		noEventing bool
	}{
		{client: "VLC", userAgent: "VLC/3.0.20 LibVLC", profileID: "vlc", wantFirst: "video/mp4"},
		{client: "Kodi", userAgent: "Kodi/21.0", profileID: "kodi", wantFirst: "video/mp4"},
		{client: "Samsung", userAgent: "Samsung Tizen DLNADOC/1.50", profileID: "samsung", wantFirst: "video/mp4"},
		{client: "LG", userAgent: "LG webOS TV", profileID: "lg", wantFirst: "video/mp4"},
		{client: "Sony", userAgent: "Sony BRAVIA", profileID: "sony", wantFirst: "video/mp4"},
		{client: "BubbleUPnP", userAgent: "BubbleUPnP/4.3", profileID: "bubbleupnp", wantFirst: "video/mp4"},
		{client: "iOS/tvOS", userAgent: "Chromecast", profileID: "chromecast", wantFirst: "application/vnd.apple.mpegurl", noEventing: true},
	}

	for _, test := range tests {
		t.Run(test.client, func(t *testing.T) {
			profile := MatchRendererProfile(RendererRequest{UserAgent: test.userAgent}, nil)
			if profile.ID != test.profileID {
				t.Fatalf("profile = %q, want %q", profile.ID, test.profileID)
			}
			protocols := SourceProtocolInfosForProfile(profile)
			if len(protocols) == 0 || !strings.Contains(protocols[0], test.wantFirst) {
				t.Fatalf("protocols[0] = %q, want %q first", protocols[0], test.wantFirst)
			}
			if profile.DisableEventing != test.noEventing {
				t.Fatalf("DisableEventing = %v, want %v", profile.DisableEventing, test.noEventing)
			}
		})
	}
}

func TestCompatibilityMediaFixturesCoverDeliveryModes(t *testing.T) {
	payload, err := os.ReadFile("testdata/compatibility_media.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Cases []struct {
			Name     string `json:"name"`
			Delivery string `json:"delivery"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(payload, &fixture); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, item := range fixture.Cases {
		seen[item.Delivery] = true
	}
	for _, delivery := range []string{"direct", "remux", "transcode", "subtitle"} {
		if !seen[delivery] {
			t.Fatalf("fixture missing %s case: %#v", delivery, fixture.Cases)
		}
	}
}
