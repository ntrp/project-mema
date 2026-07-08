package satisfaction

import (
	"testing"

	"media-manager/internal/targets"
)

func TestMediaFileRollupMissingAndDownloading(t *testing.T) {
	if got := RollupMediaFileState(false, nil, false); got != MediaFileMissing {
		t.Fatalf("state = %s", got)
	}
	if got := RollupMediaFileState(false, nil, true); got != MediaFileDownloading {
		t.Fatalf("state = %s", got)
	}
}

func TestMediaFileRollupPartialBeatsDownloaded(t *testing.T) {
	got := RollupMediaFileState(true, []targets.State{
		targets.StateSatisfied,
		targets.StatePartial,
	}, false)
	if got != MediaFilePartial {
		t.Fatalf("state = %s", got)
	}
}

func TestMediaFileRollupUpgradeable(t *testing.T) {
	got := RollupMediaFileState(true, []targets.State{
		targets.StateSatisfied,
		targets.StateUpgradeable,
	}, false)
	if got != MediaFileUpgradeable {
		t.Fatalf("state = %s", got)
	}
}

func TestMediaItemRollupMovieFromSingleRequiredFile(t *testing.T) {
	if got := RollupMediaItemState([]MediaFileState{MediaFileDownloaded}, false); got != MediaItemDownloaded {
		t.Fatalf("state = %s", got)
	}
	if got := RollupMediaItemState([]MediaFileState{MediaFilePartial}, false); got != MediaItemPartial {
		t.Fatalf("state = %s", got)
	}
}

func TestMediaItemRollupSeriesFromExpectedFiles(t *testing.T) {
	got := RollupMediaItemState([]MediaFileState{
		MediaFileDownloaded,
		MediaFileMissing,
	}, false)
	if got != MediaItemMissing {
		t.Fatalf("state = %s", got)
	}
}

func TestMediaItemRollupDownloadingDoesNotHidePartial(t *testing.T) {
	got := RollupMediaItemState([]MediaFileState{
		MediaFilePartial,
		MediaFileDownloading,
	}, true)
	if got != MediaItemPartial {
		t.Fatalf("state = %s", got)
	}
}
