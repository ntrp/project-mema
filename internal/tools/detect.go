package tools

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

type Tool struct {
	Name     string
	Required bool
}

type Status struct {
	Name      string
	Required  bool
	Available bool
	Version   string
	Path      string
	Error     string
}

var DefaultTools = []Tool{
	{Name: "ffmpeg", Required: true},
	{Name: "ffprobe", Required: true},
	{Name: "mkvmerge", Required: true},
	{Name: "mkvextract", Required: true},
	{Name: "mediainfo", Required: false},
}

func Detect(ctx context.Context, wanted []Tool) []Status {
	statuses := make([]Status, 0, len(wanted))
	for _, tool := range wanted {
		status := Status{Name: tool.Name, Required: tool.Required}

		path, err := exec.LookPath(tool.Name)
		if err != nil {
			status.Error = err.Error()
			statuses = append(statuses, status)
			continue
		}

		status.Available = true
		status.Path = path
		status.Version = versionLine(ctx, tool.Name)
		statuses = append(statuses, status)
	}
	return statuses
}

func versionLine(ctx context.Context, name string) string {
	args := []string{"-version"}
	if name == "mkvmerge" || name == "mkvextract" {
		args = []string{"--version"}
	}
	if name == "mediainfo" {
		args = []string{"--Version"}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		return ""
	}

	line, _, _ := strings.Cut(string(output), "\n")
	return strings.TrimSpace(line)
}
