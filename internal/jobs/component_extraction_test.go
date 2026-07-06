package jobs

import (
	"reflect"
	"testing"
)

func TestMkvExtractTrackArgsUseControlledTrackSyntax(t *testing.T) {
	args, err := MkvExtractTrackArgs("/library/movie/source.mkv", 2, "/library/movie/audio.mka")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"tracks", "/library/movie/source.mkv", "2:/library/movie/audio.mka"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestMkvExtractTrackArgsRejectUnsafeValues(t *testing.T) {
	for _, tc := range []struct {
		name       string
		sourcePath string
		streamID   int32
		outputPath string
	}{
		{name: "relative source", sourcePath: "movie.mkv", streamID: 1, outputPath: "/out.mka"},
		{name: "option source", sourcePath: "/library/-movie.mkv", streamID: 1, outputPath: "/out.mka"},
		{name: "negative stream", sourcePath: "/library/movie.mkv", streamID: -1, outputPath: "/out.mka"},
		{name: "relative output", sourcePath: "/library/movie.mkv", streamID: 1, outputPath: "out.mka"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := MkvExtractTrackArgs(tc.sourcePath, tc.streamID, tc.outputPath); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
