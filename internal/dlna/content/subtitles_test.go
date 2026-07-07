package content

import (
	"strings"
	"testing"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestSubtitlesForFileMapsExternalSRT(t *testing.T) {
	lang := "eng"
	format := "srt"
	item := storage.MediaItem{
		ID: uuid.New(),
		Sidecars: []storage.MediaItemSidecar{{
			MediaFilePath: "/media/movie.mkv",
			FilePath:      "/media/movie.eng.srt",
			SidecarType:   storage.MediaSidecarSubtitle,
			LanguageID:    &lang,
			Format:        &format,
		}},
	}

	subtitles := SubtitlesForFile(item, "/media/movie.mkv")

	if len(subtitles) != 1 || subtitles[0].Plan != SubtitleDirect || subtitles[0].Language != "eng" {
		t.Fatalf("subtitles = %#v", subtitles)
	}
}

func TestSubtitlePlannerConvertsOrOmitsUnsupportedFormats(t *testing.T) {
	if PlanSubtitle("ass") != SubtitleConvert {
		t.Fatalf("ass plan = %s", PlanSubtitle("ass"))
	}
	if PlanSubtitle("idx") != SubtitleOmit {
		t.Fatalf("idx plan = %s", PlanSubtitle("idx"))
	}
}

func TestDIDLIncludesSubtitleResource(t *testing.T) {
	object := Object{
		ID:       "file-1",
		ParentID: "item-1",
		Title:    "Movie.mkv",
		Class:    "object.item.videoItem",
		Kind:     ObjectItem,
		Subtitles: []Subtitle{{
			URL:    "http://127.0.0.1:18080/dlna/subtitle/file-1/0",
			Format: "srt",
		}},
	}

	payload, err := RenderDIDL([]Object{object}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := string(payload)
	if !strings.Contains(got, `protocolInfo="http-get:*:application/x-subrip:*"`) ||
		!strings.Contains(got, `/dlna/subtitle/file-1/0`) {
		t.Fatalf("DIDL missing subtitle resource:\n%s", got)
	}
}
