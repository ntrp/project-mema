package jobs

import (
	"testing"

	"media-manager/internal/storage"
)

func TestMarkActiveEpisodeDownloadsKeepsSiblingRetryAvailable(t *testing.T) {
	item := automaticPlanSeries(2)
	updated, covered := markActiveEpisodeDownloads(item, []storage.DownloadActivity{{
		MediaItemID: item.ID, ReleaseTitle: "The.Show.S01E01.1080p.WEBDL", Status: "queued",
	}})

	if covered || len(updated.FileFacts) != 1 || updated.FileFacts[0].EpisodeID == nil ||
		*updated.FileFacts[0].EpisodeID != *item.Seasons[0].Episodes[0].ID {
		t.Fatalf("covered=%v fileFacts=%#v", covered, updated.FileFacts)
	}
}

func TestMarkActiveEpisodeDownloadsStopsOverlappingPackSearch(t *testing.T) {
	item := automaticPlanSeries(2)
	_, covered := markActiveEpisodeDownloads(item, []storage.DownloadActivity{{
		MediaItemID: item.ID, ReleaseTitle: "The.Show.S01.1080p.WEBDL", Status: "downloading",
	}})
	if !covered {
		t.Fatal("active season pack should cover the pending automatic search")
	}
}

func TestMarkActiveEpisodeDownloadsStopsOverlappingMultiEpisodeSearch(t *testing.T) {
	item := automaticPlanSeries(2)
	_, covered := markActiveEpisodeDownloads(item, []storage.DownloadActivity{{
		MediaItemID: item.ID, ReleaseTitle: "The.Show.S01E01-E02.1080p.WEBDL", Status: "grabbed",
	}})
	if !covered {
		t.Fatal("active multi-episode release should cover the pending automatic search")
	}
}

func TestMarkActiveEpisodeDownloadsIgnoresCompletedHistory(t *testing.T) {
	item := automaticPlanSeries(2)
	updated, covered := markActiveEpisodeDownloads(item, []storage.DownloadActivity{{
		MediaItemID: item.ID, ReleaseTitle: "The.Show.S01.1080p.WEBDL", Status: "completed",
	}})
	if covered || len(updated.FileFacts) != 0 {
		t.Fatalf("covered=%v fileFacts=%#v", covered, updated.FileFacts)
	}
}
