package satisfaction

import (
	"testing"

	"media-manager/internal/storage"
	"media-manager/internal/targets"

	"github.com/google/uuid"
)

func TestBuildWantedRowsKeepsMissingMediaAsMediaRow(t *testing.T) {
	item := wantedTestItem()
	rows := BuildWantedRows(WantedRowsInput{Item: item})

	if len(rows) != 1 {
		t.Fatalf("rows len = %d", len(rows))
	}
	row := rows[0]
	if row.Kind != WantedRowMedia || row.ID != "media:"+item.ID.String() {
		t.Fatalf("media row = %#v", row)
	}
	if row.MediaTitle != item.Title || row.MediaType != item.Type {
		t.Fatalf("media context = %#v", row)
	}
}

func TestBuildWantedRowsForItemIncludesMissingAudioTarget(t *testing.T) {
	item := wantedTestItem()
	item.FilePaths = []string{"/media/movie.mkv"}
	item.AudioTargets = []storage.MediaProfileAudioTarget{{LanguageID: "english"}}
	item.FileFacts = []storage.MediaFileFact{{
		ID:          uuid.New(),
		MediaItemID: item.ID,
		FilePath:    "/media/movie.mkv",
		Tracks: []storage.MediaFileTrackFact{{
			ID:        uuid.New(),
			TrackType: "video",
		}},
	}}

	rows := BuildWantedRowsForItem(item)
	if len(rows) != 1 || rows[0].Kind != WantedRowTarget || rows[0].TargetType != targets.TypeAudio {
		t.Fatalf("wanted rows = %#v", rows)
	}
	if rows[0].TargetState != targets.StateMissing || rows[0].LanguageID != "english" {
		t.Fatalf("wanted target = %#v", rows[0])
	}
}

func TestBuildWantedRowsAddsExistingFileTargetProblems(t *testing.T) {
	operation := &targets.Operation{
		Type:      targets.OperationAudioTranscode,
		Manual:    true,
		Automatic: true,
		Reason:    "Transcode audio.",
	}
	season := int32(2)
	episode := int32(5)
	rows := BuildWantedRows(WantedRowsInput{
		Item:           wantedTestItem(),
		HasUsableMedia: true,
		Targets: []WantedTargetInput{{
			Target: targets.Target{
				ID:                "audio:ita",
				Type:              targets.TypeAudio,
				State:             targets.StatePending,
				LanguageID:        "ita",
				RequiredOperation: operation,
			},
			FilePath:      "/library/show/Season 02/Episode.mkv",
			SeasonNumber:  &season,
			EpisodeNumber: &episode,
		}},
	})

	if len(rows) != 1 {
		t.Fatalf("rows len = %d", len(rows))
	}
	row := rows[0]
	if row.Kind != WantedRowTarget || row.TargetType != targets.TypeAudio || row.TargetState != targets.StatePending {
		t.Fatalf("target row = %#v", row)
	}
	if row.LanguageID != "ita" || row.RequiredOperation != operation {
		t.Fatalf("target context = %#v", row)
	}
	if row.FileLabel != "Episode.mkv" || row.SeasonNumber == nil || row.EpisodeNumber == nil {
		t.Fatalf("parent context = %#v", row)
	}
}

func TestBuildWantedRowsAddsCustomFormatUpgradeSeparately(t *testing.T) {
	rows := BuildWantedRows(WantedRowsInput{
		Item:           wantedTestItem(),
		HasUsableMedia: true,
		Targets: []WantedTargetInput{{
			Target: targets.Target{
				ID:    "video:file",
				Type:  targets.TypeVideo,
				State: targets.StateUpgradeable,
			},
			FilePath: "/library/movie/Movie.mkv",
		}},
		CustomFormatUpgrades: []WantedCustomFormatUpgrade{{
			FilePath:     "/library/movie/Movie.mkv",
			CurrentScore: 25,
			TargetScore:  100,
		}},
	})

	if len(rows) != 1 {
		t.Fatalf("rows len = %d", len(rows))
	}
	row := rows[0]
	if row.Kind != WantedRowCustomFormatUpgrade {
		t.Fatalf("upgrade row = %#v", row)
	}
	if row.CurrentScore == nil || *row.CurrentScore != 25 || row.TargetScore == nil || *row.TargetScore != 100 {
		t.Fatalf("scores = %#v", row)
	}
}

func TestBuildWantedRowsDropsResolvedOrRemovedTargets(t *testing.T) {
	rows := BuildWantedRows(WantedRowsInput{
		Item:           wantedTestItem(),
		HasUsableMedia: true,
		Targets: []WantedTargetInput{{
			Target: targets.Target{ID: "video:file", State: targets.StateSatisfied},
		}},
		CustomFormatUpgrades: []WantedCustomFormatUpgrade{{
			FilePath:     "/library/movie/Movie.mkv",
			CurrentScore: 100,
			TargetScore:  100,
		}},
	})

	if len(rows) != 0 {
		t.Fatalf("rows = %#v", rows)
	}
}

func wantedTestItem() storage.MediaItem {
	return storage.MediaItem{
		ID:    uuid.MustParse("10000000-0000-4000-8000-000000000130"),
		Type:  "movie",
		Title: "Wanted Test",
	}
}
