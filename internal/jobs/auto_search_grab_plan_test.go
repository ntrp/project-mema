package jobs

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

func TestGrabAutomaticTargetsQueuesEveryEpisodeTarget(t *testing.T) {
	ctx, store := jobsTestStore(t)
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{Type: "serie", Title: "The Show"})
	if err != nil {
		t.Fatal(err)
	}
	doer := &successfulSABDoer{}
	targets := []automaticReleaseTarget{
		{label: "S01E01", candidates: []storage.ReleaseCandidateInput{{
			Title: "The.Show.S01E01.1080p.WEBDL", IndexerName: "Indexer", IndexerProtocol: "usenet", DownloadURL: "https://indexer.test/1",
		}}},
		{label: "S01E02", candidates: []storage.ReleaseCandidateInput{{
			Title: "The.Show.S01E02.1080p.WEBDL", IndexerName: "Indexer", IndexerProtocol: "usenet", DownloadURL: "https://indexer.test/2",
		}}},
	}
	clients := []storage.DownloadClient{{Name: "SAB", Type: "sabnzbd", BaseURL: "http://sab.test", Enabled: true}}

	err = grabAutomaticTargets(ctx, store, downloadclients.NewService(doer), decisions.NewEngine(), nil, item, nil, nil, nil, clients, targets)

	if err != nil || doer.calls != 2 {
		t.Fatalf("err=%v calls=%d", err, doer.calls)
	}
}

type successfulSABDoer struct {
	calls int
}

func (d *successfulSABDoer) Do(_ *http.Request) (*http.Response, error) {
	d.calls++
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"status":true,"nzo_ids":["download-id"]}`)),
	}, nil
}
