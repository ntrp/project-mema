package httpapi

import (
	"net/http"
	"strings"
	"testing"
)

func TestScenarioSCNSettings026AdminTracesDLNADecisions(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-025")

	var match DLNAProfileMatchTraceResponse
	client.doJSON(t, http.MethodPost, "/settings/dlna/profile-match-trace", DLNAProfileMatchTraceRequest{
		UserAgent: traceString("LG webOS TV"),
		Headers:   &map[string]string{"X-Device": "LG"},
	}, http.StatusOK, &match)
	if !strings.Contains(strings.ToLower(match.ProfileId+" "+match.ProfileName), "lg") {
		t.Fatalf("match = %#v", match)
	}
	if match.MatchSource == "" || len(match.RuleTrace) == 0 || len(match.Candidates) == 0 {
		t.Fatalf("expected match trace details, got %#v", match)
	}
	selectedCandidate := false
	for _, candidate := range match.Candidates {
		if candidate.Selected && candidate.ProfileId == match.ProfileId {
			selectedCandidate = true
			if !candidate.Qualified || candidate.Score < candidate.MinimumScore {
				t.Fatalf("selected candidate is not qualified: %#v", candidate)
			}
		}
	}
	if match.MatchSource == "match" && !selectedCandidate {
		t.Fatalf("selected automatic candidate missing: %#v", match.Candidates)
	}

	var delivery DLNADeliveryTraceResponse
	client.doJSON(t, http.MethodPost, "/settings/dlna/delivery-trace", DLNADeliveryTraceRequest{
		ProfileId:  traceString("lg-webos"),
		MediaPath:  traceString("/private/media/Madagascar.mkv"),
		ObjectId:   traceString("movie:madagascar"),
		ResourceId: traceString("resource:madagascar"),
		Probe: &DLNAMediaProbeRequest{
			Container:  traceString("mov,mp4,m4a,3gp,3g2,mj2"),
			VideoCodec: traceString("h264"),
			AudioCodec: traceString("dts"),
			Height:     traceInt32(1080),
		},
	}, http.StatusOK, &delivery)
	if delivery.MediaFileName != "Madagascar.mkv" || strings.Contains(delivery.MediaFileName, "/private") {
		t.Fatalf("delivery leaked path: %#v", delivery)
	}
	if delivery.Mode != DLNADeliveryModeTranscode || delivery.AudioCodec != "aac" {
		t.Fatalf("delivery = %#v", delivery)
	}
}

func traceString(value string) *string {
	return &value
}

func traceInt32(value int32) *int32 {
	return &value
}
