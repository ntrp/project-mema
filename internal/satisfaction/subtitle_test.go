package satisfaction

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestSubtitleTargetsReturnEmptyWhenNoTargetConfigured(t *testing.T) {
	if got := EvaluateSubtitleTargets(mediaItem(), &storage.MediaProfile{}, subtitleFact()); len(got.Results) != 0 {
		t.Fatalf("evaluation = %#v", got)
	}
}

func TestExternalSubtitleSatisfiesMixedAndExternalModes(t *testing.T) {
	item := mediaItem()
	item.ExternalSubtitles = []storage.MediaItemSubtitle{externalSubtitle("english", "srt")}
	for _, mode := range []string{"mixed", "external"} {
		profile := subtitleProfile(mode)
		evaluation := EvaluateSubtitleTargets(item, profile, subtitleFact())
		if evaluation.Results[0].Target.State != targets.StateSatisfied {
			t.Fatalf("mode %s evaluation = %#v", mode, evaluation)
		}
	}
}

func TestExternalSubtitlePendingInEmbeddedMode(t *testing.T) {
	item := mediaItem()
	item.ExternalSubtitles = []storage.MediaItemSubtitle{externalSubtitle("english", "srt")}
	evaluation := EvaluateSubtitleTargets(item, subtitleProfile("embedded"), subtitleFact())

	result := evaluation.Results[0]
	if result.Target.State != targets.StatePending || result.Target.RequiredOperation == nil {
		t.Fatalf("result = %#v", result)
	}
	if result.Target.RequiredOperation.Type != targets.OperationSubtitleEmbed {
		t.Fatalf("operation = %#v", result.Target.RequiredOperation)
	}
}

func TestEmbeddedSubtitlePendingInExternalMode(t *testing.T) {
	evaluation := EvaluateSubtitleTargets(mediaItem(), subtitleProfile("external"), subtitleFact(
		subtitleTrack(0, "english"),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StatePending || result.Target.RequiredOperation.Type != targets.OperationSubtitleExtraction {
		t.Fatalf("result = %#v", result)
	}
}

func TestSubtitleFormatMismatchNamesConversion(t *testing.T) {
	profile := subtitleProfile("mixed")
	profile.SubtitleTargets[0].Formats = []string{"srt"}
	evaluation := EvaluateSubtitleTargets(mediaItem(), profile, subtitleFact(
		subtitleTrackWithFormat(0, "english", "ass"),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StatePending || result.Target.RequiredOperation.Type != targets.OperationSubtitleConversion {
		t.Fatalf("result = %#v", result)
	}
}

func TestSubtitleFormatAliasesSatisfyTarget(t *testing.T) {
	profile := subtitleProfile("mixed")
	profile.SubtitleTargets[0].Formats = []string{"srt"}
	evaluation := EvaluateSubtitleTargets(mediaItem(), profile, subtitleFact(
		subtitleTrackWithFormat(0, "english", "subrip"),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StateSatisfied {
		t.Fatalf("result = %#v", result)
	}
}

func TestSubtitleFormatMismatchBlocksNonTextConversion(t *testing.T) {
	profile := subtitleProfile("mixed")
	profile.SubtitleTargets[0].Formats = []string{"srt"}
	evaluation := EvaluateSubtitleTargets(mediaItem(), profile, subtitleFact(
		subtitleTrackWithFormat(0, "english", "pgs"),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StateBlocked || result.Target.RequiredOperation != nil {
		t.Fatalf("result = %#v", result)
	}
	if result.Target.Reasons[0] != "subtitle format requires non-text conversion support" {
		t.Fatalf("reasons = %#v", result.Target.Reasons)
	}
}

func TestSubtitleFormatMismatchBlocksBitmapTargetConversion(t *testing.T) {
	profile := subtitleProfile("mixed")
	profile.SubtitleTargets[0].Formats = []string{"pgs"}
	evaluation := EvaluateSubtitleTargets(mediaItem(), profile, subtitleFact(
		subtitleTrackWithFormat(0, "english", "srt"),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StateBlocked || result.Target.RequiredOperation != nil {
		t.Fatalf("result = %#v", result)
	}
	if result.Target.Reasons[0] != "subtitle target format requires non-text conversion support" {
		t.Fatalf("reasons = %#v", result.Target.Reasons)
	}
}

func TestSubtitleUnwantedCandidatesDoNotChangeTargetState(t *testing.T) {
	profile := subtitleProfile("mixed")
	profile.RemoveUnwantedSubtitles = true
	evaluation := EvaluateSubtitleTargets(mediaItem(), profile, subtitleFact(
		subtitleTrack(0, "english"),
		subtitleTrack(1, "japanese"),
	))

	if evaluation.Results[0].Target.State != targets.StateSatisfied {
		t.Fatalf("result = %#v", evaluation.Results[0])
	}
	found := false
	for _, candidate := range evaluation.Candidates {
		if candidate.VisualState == targets.VisualUnwanted {
			found = true
		}
	}
	if !found {
		t.Fatalf("candidates = %#v", evaluation.Candidates)
	}
}

func subtitleProfile(mode string) *storage.MediaProfile {
	return &storage.MediaProfile{
		SubtitleMode: mode,
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: "english",
		}},
	}
}

func subtitleFact(tracks ...storage.MediaFileTrackFact) storage.MediaFileFact {
	factID := uuid.New()
	for index := range tracks {
		tracks[index].MediaFileFactID = factID
	}
	return storage.MediaFileFact{ID: factID, FilePath: "/media/movie.mkv", Tracks: tracks}
}

func subtitleTrack(index int32, language string) storage.MediaFileTrackFact {
	return subtitleTrackWithFormat(index, language, "srt")
}

func subtitleTrackWithFormat(index int32, language string, format string) storage.MediaFileTrackFact {
	return storage.MediaFileTrackFact{
		ID:          uuid.New(),
		StreamIndex: index,
		TrackType:   "subtitle",
		LanguageID:  &language,
		Format:      &format,
	}
}

func externalSubtitle(language string, format string) storage.MediaItemSubtitle {
	return storage.MediaItemSubtitle{
		ID:           uuid.New(),
		LanguageID:   language,
		Format:       format,
		FilePath:     "/media/" + language + "." + format,
		Selected:     true,
		DownloadedAt: time.Now().UTC(),
	}
}
