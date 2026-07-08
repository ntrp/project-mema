package dlna

import (
	"strings"
	"testing"

	"media-manager/internal/delivery"
)

func TestRendererCapabilityEvaluatorDirectPlaysCompatibleMedia(t *testing.T) {
	profile := capabilityTestProfile()
	probe := capabilityProbe("mov,mp4,m4a,3gp,3g2,mj2", "h264", "aac", 1080)

	result := EvaluateRendererCapability(profile, probe)

	if result.Decision.Mode != delivery.ModeDirect || len(result.ReasonCodes) != 0 {
		t.Fatalf("result = %#v", result)
	}
	requireTrace(t, result.Trace, "container", "pass")
	requireTrace(t, result.Trace, "videoCodec", "pass")
	requireTrace(t, result.Trace, "audioCodec", "pass")
}

func TestRendererCapabilityEvaluatorChoosesAudioTranscode(t *testing.T) {
	profile := capabilityTestProfile()
	probe := capabilityProbe("mov,mp4,m4a,3gp,3g2,mj2", "h264", "dts", 1080)

	result := EvaluateRendererCapability(profile, probe)

	if result.Decision.Mode != delivery.ModeTranscode || result.Decision.Plan.VideoCodec != "copy" {
		t.Fatalf("result = %#v", result)
	}
	if result.Decision.Plan.AudioCodec != "aac" || !hasReason(result.ReasonCodes, "audio_codec_not_supported") {
		t.Fatalf("result = %#v", result)
	}
}

func TestRendererCapabilityEvaluatorChoosesFullTranscode(t *testing.T) {
	profile := capabilityTestProfile()
	probe := capabilityProbe("mov,mp4,m4a,3gp,3g2,mj2", "hevc", "aac", 1080)

	result := EvaluateRendererCapability(profile, probe)

	if result.Decision.Mode != delivery.ModeTranscode || result.Decision.DeliveryProtocol != delivery.ProtocolHLS {
		t.Fatalf("result = %#v", result)
	}
	if result.Decision.Plan.VideoCodec != "libx264" || !hasReason(result.ReasonCodes, "video_codec_not_supported") {
		t.Fatalf("result = %#v", result)
	}
}

func TestRendererCapabilityEvaluatorChoosesRemux(t *testing.T) {
	profile := capabilityTestProfile()
	probe := capabilityProbe("matroska,webm", "h264", "aac", 1080)

	result := EvaluateRendererCapability(profile, probe)

	if result.Decision.Mode != delivery.ModeRemux || result.Decision.DeliveryProtocol != delivery.ProtocolFile {
		t.Fatalf("result = %#v", result)
	}
	if !hasReason(result.ReasonCodes, "container_not_supported") {
		t.Fatalf("reasons = %#v", result.ReasonCodes)
	}
}

func TestConnectionProtocolInfoUsesProfileCapabilities(t *testing.T) {
	profile := capabilityTestProfile()
	profile.AvoidHLS = true
	source := SourceProtocolInfoForProfile(profile)

	if !strings.Contains(source, "video/mp4") {
		t.Fatalf("source = %s", source)
	}
	if strings.Contains(source, "application/vnd.apple.mpegurl") {
		t.Fatalf("source should omit HLS: %s", source)
	}
}

func capabilityTestProfile() RendererProfile {
	return RendererProfile{
		ID: "test",
		Capabilities: RendererCapabilities{
			Containers:    []string{"mp4"},
			VideoCodecs:   []string{"h264"},
			AudioCodecs:   []string{"aac", "ac3"},
			MaxResolution: "1080p",
		},
		DeliveryRules: RendererDeliveryRules{DirectPlay: true, Transcode: true},
	}
}

func capabilityProbe(container string, videoCodec string, audioCodec string, height int32) delivery.ProbeResult {
	return delivery.ProbeResult{
		Container: delivery.Container{FormatName: &container},
		Tracks: []delivery.Track{
			{Type: delivery.TrackVideo, Codec: &videoCodec, Height: &height},
			{Type: delivery.TrackAudio, Codec: &audioCodec},
		},
	}
}

func requireTrace(t *testing.T, traces []RendererCapabilityTrace, field string, result string) {
	t.Helper()
	for _, trace := range traces {
		if trace.Field == field && trace.Result == result {
			return
		}
	}
	t.Fatalf("trace %s=%s missing in %#v", field, result, traces)
}

func hasReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if reason == want {
			return true
		}
	}
	return false
}
