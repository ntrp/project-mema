package storage

import (
	"context"
	"testing"
)

func TestDefaultMediaProfilesCarryTargetDetails(t *testing.T) {
	ctx, store := testDBStore(t)

	hd := requireMediaProfile(t, ctx, store, "hd-1080p")
	expectStrings(t, hd.VideoTarget.Codecs, []string{"h264", "hevc"})
	expectStrings(t, hd.VideoTarget.PixelFormats, []string{"yuv420p", "yuv420p10le"})
	if hd.AudioLossyTranscodePolicy != "losslessToLossy" {
		t.Fatalf("hd conversion policy = %q", hd.AudioLossyTranscodePolicy)
	}
	requireAudioTarget(t, hd, "english", "aac", []string{"2.0", "5.1"}, 192, 640)

	uhd := requireMediaProfile(t, ctx, store, "uhd-4k")
	expectStrings(t, uhd.VideoTarget.Codecs, []string{"hevc", "av1"})
	expectStrings(t, uhd.VideoTarget.HDRFormats, []string{"hdr10", "hdr10plus", "dolby-vision"})
	requireAudioTarget(t, uhd, "english", "eac3", []string{"5.1", "7.1"}, 640, 1536)

	anime := requireMediaProfile(t, ctx, store, "anime-1080p")
	requireAudioTarget(t, anime, "japanese", "aac", []string{"2.0"}, 160, 256)
	requireAudioTarget(t, anime, "english", "aac", []string{"2.0"}, 160, 256)
	requireSubtitleTarget(t, anime, "english", []string{"ass", "subrip"})
}

func requireMediaProfile(
	t *testing.T,
	ctx context.Context,
	store *SettingsStore,
	id string,
) MediaProfile {
	t.Helper()
	profile, err := store.GetMediaProfile(ctx, id)
	if err != nil {
		t.Fatalf("get profile %s: %v", id, err)
	}
	return profile
}

func requireAudioTarget(
	t *testing.T,
	profile MediaProfile,
	language string,
	codec string,
	channels []string,
	minimum int32,
	preferred int32,
) {
	t.Helper()
	for _, target := range profile.AudioTargets {
		if target.LanguageID != language {
			continue
		}
		if target.TargetCodec == nil || *target.TargetCodec != codec {
			t.Fatalf("%s audio codec = %#v", language, target.TargetCodec)
		}
		expectStrings(t, target.TargetChannels, channels)
		if target.MinimumBitrateKbps == nil || *target.MinimumBitrateKbps != minimum {
			t.Fatalf("%s minimum bitrate = %#v", language, target.MinimumBitrateKbps)
		}
		if target.PreferredBitrateKbps == nil || *target.PreferredBitrateKbps != preferred {
			t.Fatalf("%s preferred bitrate = %#v", language, target.PreferredBitrateKbps)
		}
		return
	}
	t.Fatalf("audio target %s missing from %s", language, profile.ID)
}

func requireSubtitleTarget(t *testing.T, profile MediaProfile, language string, formats []string) {
	t.Helper()
	for _, target := range profile.SubtitleTargets {
		if target.LanguageID != language {
			continue
		}
		expectStrings(t, target.Formats, formats)
		return
	}
	t.Fatalf("subtitle target %s missing from %s", language, profile.ID)
}
