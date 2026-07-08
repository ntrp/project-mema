package httpapi

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestTargetSatisfactionSummaryResponseMapsTargetsAndCandidates(t *testing.T) {
	itemID := uuid.New()
	fileID := uuid.New()
	operation := &targets.Operation{
		Type:      targets.OperationAudioTranscode,
		Manual:    true,
		Automatic: true,
		Reason:    "Transcode audio.",
	}

	got := targetSatisfactionSummaryResponse(
		[]targets.Target{{
			ID:                "audio:" + fileID.String() + ":english",
			Type:              targets.TypeAudio,
			State:             targets.StatePending,
			MediaItemID:       itemID.String(),
			MediaFileID:       fileID.String(),
			LanguageID:        "english",
			RequiredOperation: operation,
			Reasons:           []string{"Transcode audio."},
		}},
		[]targets.Candidate{{
			ID:          fileID.String() + ":stream:1",
			Type:        targets.CandidateAudioTrack,
			VisualState: targets.VisualPendingOperation,
			TargetIDs:   []string{"audio:" + fileID.String() + ":english"},
			LanguageID:  "english",
			Operation:   operation,
		}},
	)

	if len(got.Targets) != 1 || got.Targets[0].MediaItemId != itemID {
		t.Fatalf("target summary = %#v", got)
	}
	if got.Targets[0].RequiredOperation == nil || got.Targets[0].RequiredOperation.Type != TargetOperationTypeAudioTranscode {
		t.Fatalf("operation metadata missing: %#v", got.Targets[0])
	}
	if len(got.Candidates) != 1 || got.Candidates[0].VisualState != TargetCandidateVisualStatePendingOperation {
		t.Fatalf("candidate summary = %#v", got)
	}
}

func TestRollupAndWantedResponsesCoverContractVariants(t *testing.T) {
	itemID := uuid.New()
	state := mediaRollupSummaryResponse(MediaRollupStatePartial, []targets.State{
		targets.StateMissing,
		targets.StatePartial,
		targets.StateSatisfied,
	}, []string{"missing audio"})
	if state.TargetCounts.Missing != 1 || state.TargetCounts.Partial != 1 || state.TargetCounts.Satisfied != 1 {
		t.Fatalf("target counts = %#v", state.TargetCounts)
	}

	rows := []satisfaction.WantedRow{
		{ID: "media:" + itemID.String(), Kind: satisfaction.WantedRowMedia, MediaItemID: itemID.String(), MediaTitle: "Movie", MediaType: "movie"},
		{ID: "target:audio", Kind: satisfaction.WantedRowTarget, MediaItemID: itemID.String(), MediaTitle: "Movie", MediaType: "movie", TargetType: targets.TypeAudio, TargetState: targets.StateMissing, LanguageID: "english"},
		{ID: "custom-format:movie:file", Kind: satisfaction.WantedRowCustomFormatUpgrade, MediaItemID: itemID.String(), MediaTitle: "Movie", MediaType: "movie", CurrentScore: contractInt32Ptr(10), TargetScore: contractInt32Ptr(25)},
	}

	for _, row := range rows {
		if got := wantedRowResponse(row); got.Kind != WantedRowKind(row.Kind) {
			t.Fatalf("wanted row kind mismatch: %#v", got)
		}
	}
}

func TestWantedRowsResponseBuildsMediaRows(t *testing.T) {
	itemID := uuid.New()
	got := wantedRowsResponse([]storage.MediaItem{{
		ID:    itemID,
		Type:  "movie",
		Title: "Missing Movie",
	}})

	if len(got.Rows) != 1 || got.Rows[0].Kind != WantedRowKindMedia || got.Rows[0].MediaItemId != itemID {
		t.Fatalf("wanted response = %#v", got)
	}
}

func contractInt32Ptr(value int32) *int32 {
	return &value
}
