package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestMediaRenamePreviewResponseMapsRows(t *testing.T) {
	response := mediaRenamePreviewResponse(storage.MediaRenamePreview{
		Rows: []storage.MediaRenamePreviewRow{{
			CurrentPath:  "/library/old.mkv",
			ProposedPath: "/library/new.mkv",
			Status:       "safe",
			Messages:     []string{"Ready."},
		}},
	})

	if len(response.Rows) != 1 || response.Rows[0].Status != MediaRenamePreviewRowStatusSafe {
		t.Fatalf("response = %#v", response)
	}
	if response.Rows[0].Messages[0] != "Ready." {
		t.Fatalf("messages = %#v", response.Rows[0].Messages)
	}
}

func TestMediaRenameApplyResponseMapsCountsAndRows(t *testing.T) {
	response := mediaRenameApplyResponse(storage.MediaRenameApplyResult{
		AppliedCount: 1,
		SkippedCount: 2,
		FailedCount:  1,
		Rows: []storage.MediaRenamePreviewRow{{
			CurrentPath:  "/library/old.mkv",
			ProposedPath: "/library/new.mkv",
			Status:       "applied",
			Messages:     []string{"File renamed."},
		}},
	})

	if response.AppliedCount != 1 || response.SkippedCount != 2 || response.FailedCount != 1 {
		t.Fatalf("response = %#v", response)
	}
	if response.Rows[0].Status != MediaRenamePreviewRowStatusApplied {
		t.Fatalf("row = %#v", response.Rows[0])
	}
}
