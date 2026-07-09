package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubtitleEmbedCommandSetsAppendedSubtitleLanguage(t *testing.T) {
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args.txt")
	writeRunnerTool(t, dir, "ffmpeg", fmt.Sprintf("printf '%%s\\n' \"$@\" > %q\n", argsFile))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	err := runSubtitleEmbedCommand(t.Context(), "/library/movie.mkv", "/library/movie.eng.srt", "/library/out.mkv", "english", 2)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	args := strings.Split(strings.TrimSpace(string(content)), "\n")
	if !hasArgPair(args, "-metadata:s:s:2", "language=eng") {
		t.Fatalf("expected subtitle language metadata args, got %#v", args)
	}
}

func writeRunnerTool(t *testing.T, dir string, name string, body string) {
	t.Helper()
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\n" + body
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake tool: %v", err)
	}
}
