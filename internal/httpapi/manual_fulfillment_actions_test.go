package httpapi

import (
	"testing"

	"media-manager/internal/targets"
)

func TestManualFulfillmentActionResponseMapsCatalog(t *testing.T) {
	action := targets.ManualAction{
		ID:            "file_rescan",
		Operation:     targets.OperationFileRescan,
		Label:         "Rescan files",
		Description:   "Persist current facts.",
		Manual:        true,
		Automatic:     true,
		Available:     true,
		BlockedReason: "",
		Method:        "POST",
		Path:          "/media/items/{id}/files/rescan",
		WorkerPath:    "MediaFileRescan",
		StateEffect:   "Recalculate target satisfaction.",
	}

	got := manualFulfillmentActionResponse(action)

	if got.Id != action.ID || got.Operation != TargetOperationTypeFileRescan {
		t.Fatalf("unexpected action response: %#v", got)
	}
	if !got.Manual || !got.Automatic || !got.Available {
		t.Fatalf("manual action flags were not preserved: %#v", got)
	}
	if got.Path != action.Path || got.WorkerPath != action.WorkerPath {
		t.Fatalf("manual action path was not preserved: %#v", got)
	}
}
