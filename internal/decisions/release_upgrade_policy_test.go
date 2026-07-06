package decisions

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestReleaseUpgradePolicyRejectsDisabledUpgrades(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs:      []string{"webdl-1080p", "remux-2160p"},
		UpgradesAllowed: false,
	}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Movie", FilePaths: []string{"/media/Movie.2026.WEB-DL.1080p.mkv"}},
		storage.ReleaseCandidate{Title: "Movie 2026 Remux 2160p"},
		&profile,
		nil,
	)
	requireBlockedRelease(t, match)
	requireReleaseDecisionDetail(t, match, "Upgrades are disabled")
}

func TestReleaseUpgradePolicyStopsAtQualityTarget(t *testing.T) {
	target := "webdl-1080p"
	profile := storage.MediaProfile{
		QualityIDs:                        []string{"webdl-720p", "webdl-1080p", "remux-2160p"},
		UpgradesAllowed:                   true,
		UpgradeUntilQualityID:             &target,
		MinimumCustomFormatScoreIncrement: 1,
	}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Movie", FilePaths: []string{"/media/Movie.2026.WEB-DL.1080p.mkv"}},
		storage.ReleaseCandidate{Title: "Movie 2026 Remux 2160p"},
		&profile,
		nil,
	)
	requireBlockedRelease(t, match)
	requireReleaseDecisionDetail(t, match, "quality upgrade target")
}

func TestReleaseUpgradePolicyRequiresCustomFormatIncrement(t *testing.T) {
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000301")
	profile := storage.MediaProfile{
		QualityIDs:                        []string{"webdl-1080p"},
		UpgradesAllowed:                   true,
		MinimumCustomFormatScoreIncrement: 50,
		CustomFormatScores:                []storage.MediaProfileCustomFormatScore{{CustomFormatID: formatID, Score: 25}},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Preferred",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "preferred", Name: "Preferred", Type: "releaseTitle", Value: "Preferred", Required: true,
		}},
	}}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Movie", FilePaths: []string{"/media/Movie.2026.WEB-DL.1080p.mkv"}},
		storage.ReleaseCandidate{Title: "Movie 2026 WEB-DL 1080p Preferred"},
		&profile,
		formats,
	)
	requireBlockedRelease(t, match)
	requireReleaseDecisionDetail(t, match, "below the profile minimum 50")
}

func TestReleaseUpgradePolicyStopsAtCustomFormatTarget(t *testing.T) {
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000302")
	profile := storage.MediaProfile{
		QualityIDs:                    []string{"webdl-1080p"},
		UpgradesAllowed:               true,
		UpgradeUntilCustomFormatScore: 25,
		CustomFormatScores:            []storage.MediaProfileCustomFormatScore{{CustomFormatID: formatID, Score: 25}},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Preferred",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "preferred", Name: "Preferred", Type: "releaseTitle", Value: "Preferred", Required: true,
		}},
	}}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Movie", FilePaths: []string{"/media/Movie.2026.WEB-DL.1080p.Preferred.mkv"}},
		storage.ReleaseCandidate{Title: "Movie 2026 WEB-DL 1080p Preferred"},
		&profile,
		formats,
	)
	requireBlockedRelease(t, match)
	requireReleaseDecisionDetail(t, match, "custom format upgrade target")
}

func TestReleaseUpgradePolicyAllowsQualityUpgrade(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs:                        []string{"webdl-1080p", "remux-2160p"},
		UpgradesAllowed:                   true,
		MinimumCustomFormatScoreIncrement: 1,
	}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Movie", FilePaths: []string{"/media/Movie.2026.WEB-DL.1080p.mkv"}},
		storage.ReleaseCandidate{Title: "Movie 2026 Remux 2160p"},
		&profile,
		nil,
	)
	if match.Severity != "info" {
		t.Fatalf("expected upgradeable release, got %q: %v", match.Severity, match.Details)
	}
	requireReleaseDecisionDetail(t, match, "higher than the current file")
}

func requireBlockedRelease(t *testing.T, match ReleaseMatch) {
	t.Helper()
	if match.Severity != "error" {
		t.Fatalf("expected blocked release, got %q: %v", match.Severity, match.Details)
	}
}

func requireReleaseDecisionDetail(t *testing.T, match ReleaseMatch, value string) {
	t.Helper()
	for _, detail := range match.Details {
		if strings.Contains(detail, value) {
			return
		}
	}
	t.Fatalf("detail containing %q missing from %v", value, match.Details)
}
