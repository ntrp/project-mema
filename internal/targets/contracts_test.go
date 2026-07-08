package targets

import "testing"

func TestContractStateListExcludesNonTargetStates(t *testing.T) {
	states := []State{
		StateMissing,
		StatePartial,
		StatePending,
		StateSatisfied,
		StateUpgradeable,
		StateBlocked,
		StateFailed,
	}

	for _, state := range states {
		if state == "available" || state == "unmanaged" {
			t.Fatalf("state %q must not be a target state", state)
		}
	}
}

func TestContractSeparatesTargetAndCandidateStates(t *testing.T) {
	target := Target{Type: TypeAudio, State: StateMissing}
	candidate := Candidate{Type: CandidateAudioTrack, VisualState: VisualUnwanted}

	if State(candidate.VisualState) == target.State {
		t.Fatalf("candidate visual state %q should not equal target state %q", candidate.VisualState, target.State)
	}
}
