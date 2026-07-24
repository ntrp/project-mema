package controlpoint

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type advertisement struct {
	Location string
	UDN      string
	MaxAge   time.Duration
	Byebye   bool
}

func parseAdvertisement(payload []byte) (advertisement, bool) {
	text := strings.ReplaceAll(string(payload), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return advertisement{}, false
	}
	start := strings.ToUpper(strings.TrimSpace(lines[0]))
	isNotify := start == "NOTIFY * HTTP/1.1"
	if start != "HTTP/1.1 200 OK" && !isNotify {
		return advertisement{}, false
	}
	headers := http.Header{}
	for _, line := range lines[1:] {
		key, value, ok := strings.Cut(line, ":")
		if ok {
			headers.Add(strings.TrimSpace(key), strings.TrimSpace(value))
		}
	}
	target := headers.Get("ST")
	if isNotify {
		target = headers.Get("NT")
	}
	if !strings.EqualFold(target, MediaRendererDevice) {
		return advertisement{}, false
	}
	usn := headers.Get("USN")
	udn, _, _ := strings.Cut(usn, "::")
	notification := strings.ToLower(strings.TrimSpace(headers.Get("NTS")))
	if isNotify && notification != "ssdp:alive" && notification != "ssdp:update" && notification != "ssdp:byebye" {
		return advertisement{}, false
	}
	result := advertisement{
		Location: strings.TrimSpace(headers.Get("LOCATION")),
		UDN:      strings.TrimSpace(udn),
		MaxAge:   cacheMaxAge(headers.Get("CACHE-CONTROL")),
		Byebye:   isNotify && notification == "ssdp:byebye",
	}
	if result.UDN == "" || !result.Byebye && result.Location == "" {
		return advertisement{}, false
	}
	return result, true
}

func cacheMaxAge(value string) time.Duration {
	for _, part := range strings.Split(value, ",") {
		key, raw, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok || !strings.EqualFold(key, "max-age") {
			continue
		}
		seconds, err := strconv.Atoi(strings.Trim(raw, `"`))
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return 30 * time.Second
}
