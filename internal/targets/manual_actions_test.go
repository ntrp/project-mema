package targets

import "testing"

func TestManualActionsCoverEveryAutomaticOperation(t *testing.T) {
	covered := map[OperationType]bool{}
	for _, action := range ManualActions() {
		if action.Automatic && !action.Manual {
			t.Fatalf("action %q is automatic-only", action.ID)
		}
		if action.Available {
			if action.Method == "" || action.Path == "" || action.WorkerPath == "" {
				t.Fatalf("available action %q must expose method, path, and worker path", action.ID)
			}
			covered[action.Operation] = true
		}
		if !action.Available && action.BlockedReason == "" {
			t.Fatalf("blocked action %q must explain why it is unavailable", action.ID)
		}
	}
	for _, operation := range OperationTypes() {
		if !covered[operation] {
			t.Fatalf("operation %q has no available manual action", operation)
		}
	}
}

func TestManualActionsIncludeReleaseAndSubtitleChoices(t *testing.T) {
	required := map[string]OperationType{
		"release_search":      OperationReleaseSearch,
		"release_grab":        OperationReleaseSearch,
		"release_import":      OperationReleaseSearch,
		"subtitle_download":   OperationSubtitleDownload,
		"subtitle_grab":       OperationSubtitleDownload,
		"subtitle_extraction": OperationSubtitleExtraction,
	}
	for _, action := range ManualActions() {
		if operation, ok := required[action.ID]; ok && action.Operation == operation {
			delete(required, action.ID)
		}
	}
	for id := range required {
		t.Fatalf("missing manual action %q", id)
	}
}
