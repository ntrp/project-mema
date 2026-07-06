package jobs

import (
	"reflect"
	"testing"

	"media-manager/internal/storage"
)

func TestMkvMergeArgsUseControlledOutputAndInputs(t *testing.T) {
	args, err := MkvMergeArgs("/library/out.mkv", []string{"/library/base.mkv", "/library/audio.mka"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-o", "/library/out.mkv", "/library/base.mkv", "/library/audio.mka"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestMkvMergeArgsRejectUnsafeValues(t *testing.T) {
	for _, tc := range []struct {
		name   string
		output string
		inputs []string
	}{
		{name: "few inputs", output: "/out.mkv", inputs: []string{"/base.mkv"}},
		{name: "relative output", output: "out.mkv", inputs: []string{"/base.mkv", "/audio.mka"}},
		{name: "option input", output: "/out.mkv", inputs: []string{"/base.mkv", "/tmp/-audio.mka"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := MkvMergeArgs(tc.output, tc.inputs); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestComponentMuxCommandUsesFfmpegForMP4(t *testing.T) {
	command, args, err := ComponentMuxCommand("/library/out.mp4", []storage.MediaComponentAssemblyInput{
		{StreamType: "video", InputPath: "/library/base.mkv"},
		{StreamType: "audio", InputPath: "/library/audio.mka"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if command != "ffmpeg" {
		t.Fatalf("command = %q", command)
	}
	want := []string{
		"-y", "-i", "/library/base.mkv", "-i", "/library/audio.mka",
		"-map", "0:v:0?", "-map", "1:a:0?", "-c", "copy", "/library/out.mp4",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestComponentMuxCommandRejectsMP4SubtitleInput(t *testing.T) {
	_, _, err := ComponentMuxCommand("/library/out.mp4", []storage.MediaComponentAssemblyInput{
		{StreamType: "video", InputPath: "/library/base.mkv"},
		{StreamType: "subtitle", InputPath: "/library/subs.srt"},
	})
	if err == nil {
		t.Fatal("expected mp4 subtitle rejection")
	}
}
