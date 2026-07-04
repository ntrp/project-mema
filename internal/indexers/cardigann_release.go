package indexers

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

func applyCardigannField(release *Release, name string, value string, searchURL string) {
	field := canonicalCardigannField(name)
	if field != name || strings.HasPrefix(name, "_") {
		return
	}
	switch field {
	case "title":
		release.Title = value
	case "details", "comments":
		release.InfoURL = resolveCardigannLink(searchURL, value)
	case "download":
		release.DownloadURL = resolveCardigannLink(searchURL, value)
	case "guid":
		release.GUID = value
	case "size":
		release.SizeBytes = parseSizeBytes(value)
	case "seeders":
		release.Seeders = int32Value(value)
	case "leechers", "peers":
		release.Peers = int32Value(value)
	case "date", "publishdate":
		release.PublishedAt = cardigannPublishedAt(value)
	case "infohash":
		if release.DownloadURL == "" {
			release.DownloadURL = "magnet:?xt=urn:btih:" + strings.TrimSpace(value)
		}
	}
}

func finalizeCardigannRelease(release Release) (Release, bool) {
	if release.Title == "" {
		return release, false
	}
	if release.DownloadURL == "" {
		return release, false
	}
	if release.InfoURL == "" {
		release.InfoURL = release.DownloadURL
	}
	if release.GUID == "" {
		release.GUID = firstNonEmpty(release.InfoURL, release.DownloadURL, release.Title)
	}
	return release, true
}

func canonicalCardigannField(name string) string {
	if before, _, ok := strings.Cut(name, "_"); ok {
		return before
	}
	return name
}

func resolveCardigannLink(base string, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "magnet:") {
		return value
	}
	parsed, err := url.Parse(value)
	if err == nil && parsed.IsAbs() {
		return parsed.String()
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return value
	}
	relative, err := url.Parse(value)
	if err != nil {
		return value
	}
	return baseURL.ResolveReference(relative).String()
}

func int32Value(value string) *int32 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 32)
	if err != nil {
		return nil
	}
	result := int32(parsed)
	return &result
}

func cardigannPublishedAt(value string) *time.Time {
	value = strings.TrimSpace(value)
	if strings.EqualFold(value, "now") {
		now := time.Now()
		return &now
	}
	parsed, ok := parseCardigannDate(value, "")
	if !ok {
		return nil
	}
	return &parsed
}
