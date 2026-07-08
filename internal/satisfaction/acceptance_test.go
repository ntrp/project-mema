package satisfaction

import (
	"testing"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestTargetStateAcceptanceMatrixFeedsRollupAndWantedRows(t *testing.T) {
	wantedStates := map[targets.State]bool{
		targets.StateMissing: true,
		targets.StatePartial: true,
		targets.StatePending: true,
		targets.StateBlocked: true,
		targets.StateFailed:  true,
	}
	tests := []struct {
		state      targets.State
		fileState  MediaFileState
		wantedRows int
	}{
		{targets.StateMissing, MediaFilePartial, 1},
		{targets.StatePartial, MediaFilePartial, 1},
		{targets.StatePending, MediaFilePartial, 1},
		{targets.StateSatisfied, MediaFileDownloaded, 0},
		{targets.StateUpgradeable, MediaFileUpgradeable, 0},
		{targets.StateBlocked, MediaFilePartial, 1},
		{targets.StateFailed, MediaFilePartial, 1},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := RollupMediaFileState(true, []targets.State{tt.state}, false); got != tt.fileState {
				t.Fatalf("file state = %s", got)
			}
			rows := BuildWantedRows(WantedRowsInput{
				Item:           wantedTestItem(),
				HasUsableMedia: true,
				Targets: []WantedTargetInput{{
					Target: targets.Target{ID: "target:" + string(tt.state), Type: targets.TypeVideo, State: tt.state},
				}},
			})
			if len(rows) != tt.wantedRows {
				t.Fatalf("rows = %#v", rows)
			}
			if wantedStates[tt.state] && rows[0].TargetState != tt.state {
				t.Fatalf("target state = %s", rows[0].TargetState)
			}
		})
	}
}

func TestSubtitleModeTransitionsRecalculateFromSamePersistedFacts(t *testing.T) {
	item := mediaItem()
	item.ExternalSubtitles = []storage.MediaItemSubtitle{externalSubtitle("english", "srt")}
	fact := subtitleFact(subtitleTrack(0, "english"))

	tests := []struct {
		mode         string
		wantExternal targets.State
		wantEmbedded targets.State
	}{
		{"embedded", targets.StatePending, targets.StateSatisfied},
		{"mixed", targets.StateSatisfied, targets.StateSatisfied},
		{"external", targets.StateSatisfied, targets.StatePending},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			externalOnly := EvaluateSubtitleTargets(item, subtitleProfile(tt.mode), subtitleFact())
			if got := externalOnly.Results[0].Target.State; got != tt.wantExternal {
				t.Fatalf("external state = %s", got)
			}
			embeddedOnly := EvaluateSubtitleTargets(mediaItem(), subtitleProfile(tt.mode), fact)
			if got := embeddedOnly.Results[0].Target.State; got != tt.wantEmbedded {
				t.Fatalf("embedded state = %s", got)
			}
		})
	}
}
