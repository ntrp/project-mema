package jobs

import (
	"context"
	"testing"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func TestAutoSearchDownloadSkipsManualMedia(t *testing.T) {
	err := autoSearchDownload(context.Background(), nil, nil, nil, decisions.NewEngine(), nil, storage.MediaItem{
		Title:  "Manual Only",
		Manual: true,
	})
	if err != nil {
		t.Fatalf("autoSearchDownload returned error: %v", err)
	}
}
