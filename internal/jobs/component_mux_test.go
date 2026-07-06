package jobs

import (
	"reflect"
	"testing"
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
