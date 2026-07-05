package tools

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunOutputUsesArgvAndBoundedOutput(t *testing.T) {
	dir := t.TempDir()
	writeRunnerTool(t, dir, "echo-args", "printf '%s' \"$1\"\n")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	output, err := RunOutput(context.Background(), CommandSpec{
		Name:           "echo-args",
		Args:           []string{"file name; still one arg"},
		MaxOutputBytes: 64,
		MaxStderrBytes: 64,
	})

	if err != nil {
		t.Fatalf("RunOutput error = %v", err)
	}
	if string(output) != "file name; still one arg" {
		t.Fatalf("output = %q", output)
	}
}

func TestRunOutputRejectsUnsafeToolNames(t *testing.T) {
	for _, name := range []string{"", "-ffmpeg", "/usr/bin/ffmpeg", `dir\ffmpeg`} {
		_, err := RunOutput(context.Background(), CommandSpec{Name: name})
		if err == nil {
			t.Fatalf("expected %q to be rejected", name)
		}
	}
}

func TestRunOutputCapsOutput(t *testing.T) {
	dir := t.TempDir()
	writeRunnerTool(t, dir, "too-loud", "printf '1234567890'\n")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	_, err := RunOutput(context.Background(), CommandSpec{
		Name:           "too-loud",
		MaxOutputBytes: 4,
		MaxStderrBytes: 64,
	})

	if !errors.Is(err, ErrOutputLimit) {
		t.Fatalf("error = %v, want ErrOutputLimit", err)
	}
}

func TestRunOutputReturnsBoundedStderr(t *testing.T) {
	dir := t.TempDir()
	writeRunnerTool(t, dir, "fail-tool", "printf 'controlled failure' >&2\nexit 2\n")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	_, err := RunOutput(context.Background(), CommandSpec{
		Name:           "fail-tool",
		MaxOutputBytes: 64,
		MaxStderrBytes: 64,
	})

	if err == nil || !strings.Contains(err.Error(), "controlled failure") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunOutputHonorsTimeout(t *testing.T) {
	dir := t.TempDir()
	writeRunnerTool(t, dir, "slow-tool", "sleep 1\n")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	_, err := RunOutput(context.Background(), CommandSpec{
		Name:           "slow-tool",
		Timeout:        10 * time.Millisecond,
		MaxOutputBytes: 64,
		MaxStderrBytes: 64,
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want deadline exceeded", err)
	}
}

func TestSafePathArgRejectsArgumentLikePaths(t *testing.T) {
	if err := SafePathArg(filepath.Join(t.TempDir(), "movie.mkv")); err != nil {
		t.Fatalf("SafePathArg normal path = %v", err)
	}
	for _, path := range []string{"movie.mkv", filepath.Join(t.TempDir(), "-movie.mkv"), ""} {
		if err := SafePathArg(path); err == nil {
			t.Fatalf("expected %q to be rejected", path)
		}
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
