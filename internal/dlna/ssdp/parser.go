package ssdp

import (
	"strconv"
	"strings"
	"time"
)

func ParseSearch(payload []byte) (SearchRequest, bool) {
	text := strings.ReplaceAll(string(payload), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || !strings.EqualFold(strings.TrimSpace(lines[0]), "M-SEARCH * HTTP/1.1") {
		return SearchRequest{}, false
	}
	headers := map[string]string{}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		headers[strings.ToLower(strings.TrimSpace(name))] = strings.Trim(strings.TrimSpace(value), `"`)
	}
	if !strings.EqualFold(headers["man"], "ssdp:discover") {
		return SearchRequest{}, false
	}
	target := strings.TrimSpace(headers["st"])
	if target == "" {
		return SearchRequest{}, false
	}
	return SearchRequest{Target: target, MX: parseMX(headers["mx"])}, true
}

func parseMX(value string) time.Duration {
	seconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || seconds <= 0 {
		return 0
	}
	if seconds > 5 {
		seconds = 5
	}
	return time.Duration(seconds) * time.Second
}
