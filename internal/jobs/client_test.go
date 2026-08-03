package jobs

import (
	"testing"

	"github.com/riverqueue/river/rivertype"
)

func TestAutoSearchCanRunAgainAfterCompletion(t *testing.T) {
	opts := autoSearchDownloadInsertOpts()
	if !opts.UniqueOpts.ByArgs {
		t.Fatal("concurrent searches for the same media item must remain unique")
	}
	states := opts.UniqueOpts.ByState
	if containsJobState(states, rivertype.JobStateCompleted) {
		t.Fatal("completed jobs must not block a new automatic search")
	}
	for _, required := range []rivertype.JobState{
		rivertype.JobStateAvailable,
		rivertype.JobStatePending,
		rivertype.JobStateRunning,
		rivertype.JobStateScheduled,
		rivertype.JobStateRetryable,
	} {
		if !containsJobState(states, required) {
			t.Fatalf("active state %q must remain unique", required)
		}
	}
}

func containsJobState(states []rivertype.JobState, target rivertype.JobState) bool {
	for _, state := range states {
		if state == target {
			return true
		}
	}
	return false
}
