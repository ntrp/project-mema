package dlna

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"media-manager/internal/storage"
)

func TestRepositoryRendererProfilesKeepCurrentMatches(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")

	tests := []struct {
		name    string
		request RendererRequest
		want    string
	}{
		{name: "vlc", request: RendererRequest{UserAgent: "VLC/3.0"}, want: "vlc"},
		{name: "samsung", request: RendererRequest{UserAgent: "Samsung Tizen DLNADOC"}, want: "samsung"},
		{name: "lg", request: RendererRequest{UserAgent: "LG webOS TV"}, want: "lg"},
		{name: "sony", request: RendererRequest{UserAgent: "Sony BRAVIA"}, want: "sony"},
		{name: "chromecast", request: RendererRequest{Headers: http.Header{"X-Device": []string{"Google Cast"}}}, want: "chromecast"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := manager.ExplainRendererProfile(ctx, test.request).Profile.ID; got != test.want {
				t.Fatalf("profile = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRepositoryRendererProfileOverridesWin(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	ip := "192.0.2.55"
	rendererUUID := "uuid:living-room"
	upsertOverride(t, ctx, store, storage.DLNARendererDeviceOverrideInput{
		IPAddress:               &ip,
		ProfileID:               "chromecast",
		Allowed:                 true,
		DeliveryPolicyOverrides: rawObject(),
	})
	upsertOverride(t, ctx, store, storage.DLNARendererDeviceOverrideInput{
		RendererUUID:            &rendererUUID,
		ProfileID:               "lg-tv-2025",
		Allowed:                 true,
		DeliveryPolicyOverrides: rawObject(),
	})
	if err := manager.RefreshRendererProfiles(ctx); err != nil {
		t.Fatalf("refresh profiles: %v", err)
	}

	ipMatch := manager.ExplainRendererProfile(ctx, RendererRequest{ClientIP: ip, UserAgent: "VLC/3.0"})
	if ipMatch.Profile.ID != "chromecast" || ipMatch.Explanation.MatchSource != "manual_ip" {
		t.Fatalf("ip match = %#v", ipMatch)
	}
	uuidMatch := manager.ExplainRendererProfile(ctx, RendererRequest{RendererUUID: rendererUUID, UserAgent: "VLC/3.0"})
	if uuidMatch.Profile.ID != "lg-tv-2025" || uuidMatch.Explanation.MatchSource != "manual_uuid" {
		t.Fatalf("uuid match = %#v", uuidMatch)
	}
}

func TestRepositoryRendererProfilePriorityAndExplanation(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	request := RendererRequest{
		UserAgent: "LG webOS TV",
		Headers:   http.Header{"X-Model-Year": []string{"LG 2025"}},
	}

	match := manager.ExplainRendererProfile(ctx, request)

	if match.Profile.ID != "lg-tv-2025" {
		t.Fatalf("profile = %q, explanation=%#v", match.Profile.ID, match.Explanation)
	}
	if match.Explanation.MatchSource != "match" || match.Explanation.WinningRule != "headers:2025" {
		t.Fatalf("explanation = %#v", match.Explanation)
	}
	if !containsString(match.Explanation.CandidateProfileIDs, "lg") {
		t.Fatalf("candidates = %#v", match.Explanation.CandidateProfileIDs)
	}
}

func TestRepositoryRendererProfileStickyIPRegression(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	request := httptest.NewRequest(http.MethodPost, "/dlna/control/content-directory", nil)
	request.RemoteAddr = "192.0.2.77:1234"
	request.Header.Set("User-Agent", "LG webOS TV")
	manager.recordClient(request, "Browse", nil, nil)

	match := manager.ExplainRendererProfile(ctx, RendererRequest{
		ClientIP:  "192.0.2.77",
		UserAgent: "Lavf/60.0",
	})

	if match.Profile.ID != "lg" || match.Explanation.MatchSource != "sticky_ip" {
		t.Fatalf("sticky match = %#v", match)
	}
}

func upsertOverride(
	t *testing.T,
	ctx context.Context,
	store *storage.SettingsStore,
	input storage.DLNARendererDeviceOverrideInput,
) {
	t.Helper()
	if _, err := store.UpsertDLNARendererDeviceOverride(ctx, input); err != nil {
		t.Fatalf("upsert override: %v", err)
	}
}

func rawObject() json.RawMessage {
	return json.RawMessage(`{}`)
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
