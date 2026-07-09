package jobs

import (
	"reflect"
	"testing"
)

func TestFfmpegProgressArgsWritesProgressBeforeOutput(t *testing.T) {
	got := ffmpegProgressArgs([]string{"-y", "-i", "/in.mkv", "/out.mkv"})
	want := []string{"-y", "-i", "/in.mkv", "-nostats", "-progress", "pipe:1", "/out.mkv"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
}

func TestAudioTranscodeProgressParsesPercent(t *testing.T) {
	progress := audioTranscodeProgress{durationMs: 10_000}
	got, ok := progress.percent("out_time_us=2500000")

	if !ok || got != 25 {
		t.Fatalf("percent = %d ok=%t", got, ok)
	}
}

func TestAudioTranscodeProgressParsesClockPercent(t *testing.T) {
	progress := audioTranscodeProgress{durationMs: 10_000}
	got, ok := progress.percent("out_time=00:00:02.500000")

	if !ok || got != 25 {
		t.Fatalf("percent = %d ok=%t", got, ok)
	}
}

func TestAudioTranscodeProgressClampsBeforeDone(t *testing.T) {
	progress := audioTranscodeProgress{durationMs: 10_000}
	got, ok := progress.percent("out_time_us=11000000")

	if !ok || got != 99 {
		t.Fatalf("percent = %d ok=%t", got, ok)
	}
}
