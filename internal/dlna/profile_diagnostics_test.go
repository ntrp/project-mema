package dlna

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/delivery"
)

func TestRendererProfileTraceShowsRuleReasonsAndSafeHeaders(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, profiles: []RendererProfile{
		{
			ID: "lg-test", Name: "LG Test", Priority: 10, MatchMinScore: 3,
			rules: []rendererMatchRule{{Field: "userAgent", Contains: "webos", Score: 5}},
		},
		{ID: "sony-test", Name: "Sony Test", rules: []rendererMatchRule{{Field: "userAgent", Contains: "bravia", Score: 5}}},
		{ID: "generic", Name: "Generic DLNA"},
	}}

	trace := manager.TraceRendererProfile(context.Background(), RendererRequest{
		UserAgent: "LG webOS TV",
		Headers: http.Header{
			"X-Device":      []string{"LG"},
			"Authorization": []string{"secret"},
		},
	})

	if trace.Match.Profile.ID != "lg-test" || trace.Match.Explanation.WinningRule != "userAgent:webos" {
		t.Fatalf("trace match = %#v", trace.Match)
	}
	requireProfileRuleTrace(t, trace.Rules, "lg-test", "pass")
	requireProfileRuleTrace(t, trace.Rules, "sony-test", "fail")
	for _, header := range trace.HeadersSummary {
		if strings.Contains(strings.ToLower(header), "authorization") || strings.Contains(header, "secret") {
			t.Fatalf("unsafe header summary = %#v", trace.HeadersSummary)
		}
	}
}

func TestDeliveryDecisionTraceSanitizesMediaPathAndExplainsAudioTranscode(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, profiles: []RendererProfile{capabilityTestProfile()}}

	trace := manager.TraceDeliveryDecision(context.Background(), DeliveryTraceInput{
		ProfileID:  "test",
		MediaPath:  "/private/media/Madagascar.mkv",
		ObjectID:   "movie:1",
		ResourceID: "resource:1",
		StreamMode: "direct",
		Probe:      capabilityProbe("mov,mp4,m4a,3gp,3g2,mj2", "h264", "dts", 1080),
	})

	if trace.MediaFileName != "Madagascar.mkv" || strings.Contains(trace.MediaFileName, "/private") {
		t.Fatalf("media file name = %q", trace.MediaFileName)
	}
	if trace.Decision.Mode != delivery.ModeTranscode || trace.Decision.Plan.AudioCodec != "aac" {
		t.Fatalf("decision = %#v", trace.Decision)
	}
	requireTrace(t, trace.Trace, "audioCodec", "fail")
}

func requireProfileRuleTrace(t *testing.T, traces []RendererProfileRuleTrace, profileID string, result string) {
	t.Helper()
	for _, trace := range traces {
		if trace.ProfileID == profileID && trace.Result == result {
			return
		}
	}
	t.Fatalf("profile trace %s=%s missing in %#v", profileID, result, traces)
}
