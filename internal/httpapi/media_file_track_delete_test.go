package httpapi

import (
	"reflect"
	"testing"
)

func TestMediaFileTrackDeleteArgsRemovesAudioStream(t *testing.T) {
	trackIndex := int32(2)
	command, err := mediaFileTrackDeleteArgs("/media/movie.mkv", "/media/movie.tmp.mkv", MediaFileTrackDeleteRequest{
		Path:       "/media/movie.mkv",
		TargetType: MediaFileTrackDeleteRequestTargetTypeAudio,
		TrackIndex: &trackIndex,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"-hide_banner", "-loglevel", "error", "-y", "-i", "/media/movie.mkv",
		"-map", "0", "-map", "-0:2", "-c", "copy", "/media/movie.tmp.mkv",
	}
	if !reflect.DeepEqual(command.args, want) {
		t.Fatalf("args = %#v", command.args)
	}
}

func TestMediaFileTrackDeleteArgsRemovesAllChapters(t *testing.T) {
	command, err := mediaFileTrackDeleteArgs("/media/movie.mkv", "/media/movie.tmp.mkv", MediaFileTrackDeleteRequest{
		Path:       "/media/movie.mkv",
		TargetType: MediaFileTrackDeleteRequestTargetTypeChapters,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"-hide_banner", "-loglevel", "error", "-y", "-i", "/media/movie.mkv",
		"-map", "0", "-map_chapters", "-1", "-c", "copy", "/media/movie.tmp.mkv",
	}
	if !reflect.DeepEqual(command.args, want) {
		t.Fatalf("args = %#v", command.args)
	}
}

func TestFFMetadataChaptersEscapesTitles(t *testing.T) {
	start := "0.000000"
	end := "60.500000"
	title := "Opening=One; #1"
	metadata, err := ffmetadataChapters([]MediaFileChapter{{
		Index:     0,
		StartTime: &start,
		EndTime:   &end,
		Title:     &title,
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := ";FFMETADATA1\n[CHAPTER]\nTIMEBASE=1/1000\nSTART=0\nEND=60500\ntitle=Opening\\=One\\; \\#1\n"
	if metadata != want {
		t.Fatalf("metadata = %q", metadata)
	}
}
