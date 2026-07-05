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

	if len(response.Rows) != 1 || response.Rows[0].Status != Safe {
		t.Fatalf("response = %#v", response)
	}
	if response.Rows[0].Messages[0] != "Ready." {
		t.Fatalf("messages = %#v", response.Rows[0].Messages)
	}
}
