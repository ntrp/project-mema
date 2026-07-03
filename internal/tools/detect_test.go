package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScenarioSCNSystem005DetectReportsMissingTool(t *testing.T) {
	statuses := Detect(context.Background(), []Tool{{Name: "definitely-not-a-real-media-tool", Required: true}})

	if len(statuses) != 1 {
		t.Fatalf("statuses = %#v", statuses)
	}
	status := statuses[0]
	if status.Available || !status.Required || status.Path != "" || status.Error == "" {
		t.Fatalf("status = %#v", status)
	}
}

func TestScenarioSCNSystem005VersionLineReturnsEmptyForMissingTool(t *testing.T) {
	if got := versionLine(context.Background(), "definitely-not-a-real-media-tool"); got != "" {
		t.Fatalf("version = %q", got)
	}
}

func TestScenarioSCNSystem005DetectReportsAvailableToolVersions(t *testing.T) {
	dir := t.TempDir()
	createFakeTool(t, dir, "scenario-tool", "-version", "scenario-tool version 1.2.3")
	createFakeTool(t, dir, "mkvmerge", "--version", "mkvmerge v99")
	createFakeTool(t, dir, "mediainfo", "--Version", "mediainfo v24")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	statuses := Detect(context.Background(), []Tool{
		{Name: "scenario-tool"},
		{Name: "mkvmerge", Required: true},
		{Name: "mediainfo"},
	})

	if len(statuses) != 3 {
		t.Fatalf("statuses = %#v", statuses)
	}
	for _, status := range statuses {
		if !status.Available || status.Path == "" || status.Error != "" {
			t.Fatalf("status = %#v", status)
		}
	}
	if statuses[0].Version != "scenario-tool version 1.2.3" {
		t.Fatalf("scenario-tool version = %q", statuses[0].Version)
	}
	if statuses[1].Version != "mkvmerge v99" || !statuses[1].Required {
		t.Fatalf("mkvmerge status = %#v", statuses[1])
	}
	if statuses[2].Version != "mediainfo v24" {
		t.Fatalf("mediainfo status = %#v", statuses[2])
	}
}

func TestScenarioSCNSystem005DefaultToolListIncludesRequiredMediaTools(t *testing.T) {
	names := make([]string, 0, len(DefaultTools))
	for _, tool := range DefaultTools {
		if tool.Required {
			names = append(names, tool.Name)
		}
	}
	got := strings.Join(names, ",")
	for _, want := range []string{"ffmpeg", "ffprobe", "mkvmerge", "mkvextract"} {
		if !strings.Contains(got, want) {
			t.Fatalf("required tool %q missing from %q", want, got)
		}
	}
}

func createFakeTool(t *testing.T, dir string, name string, expectedArg string, output string) {
	t.Helper()
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\n" +
		"if [ \"$1\" != \"" + expectedArg + "\" ]; then\n" +
		"  echo unexpected argument \"$1\"\n" +
		"  exit 2\n" +
		"fi\n" +
		"echo \"" + output + "\"\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake tool: %v", err)
	}
}
