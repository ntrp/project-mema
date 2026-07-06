package decisions

import (
	"strings"
	"testing"

	"media-manager/internal/storage"
)

func TestVideoTargetScoringAddsContributors(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-2160p"},
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:              []string{"hevc"},
			CodecRequired:       true,
			CodecScore:          100,
			HDRFormats:          []string{"hdr10"},
			HDRRequired:         true,
			HDRScore:            75,
			PixelFormats:        []string{"yuv420p10le"},
			PixelFormatRequired: true,
			PixelFormatScore:    25,
		},
	}

	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.2160p.WEBDL.HDR10.10bit.H265"},
		&profile,
		nil,
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
	if match.TargetScore != 200 || match.Score != 200 {
		t.Fatalf("scores = target %d total %d, want 200", match.TargetScore, match.Score)
	}
	if len(match.TargetContributors) != 3 {
		t.Fatalf("contributors = %#v", match.TargetContributors)
	}
}

func TestVideoTargetScoringRejectsRequiredKnownMismatch(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:        []string{"hevc"},
			CodecRequired: true,
		},
	}

	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.1080p.WEBDL.H264"},
		&profile,
		nil,
	)

	if match.Severity != "error" {
		t.Fatalf("expected target rejection, got %q: %v", match.Severity, match.Details)
	}
	if len(match.Details) == 0 || !strings.Contains(match.Details[0], "Video codec h264") {
		t.Fatalf("details = %#v", match.Details)
	}
}

func TestVideoTargetScoringDoesNotRejectMissingMetadata(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:        []string{"hevc"},
			CodecRequired: true,
		},
	}

	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.1080p.WEBDL"},
		&profile,
		nil,
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
}
