package clusterapi

import (
	"bytes"
	"encoding/json"
	"path"
	"strings"

	"media-manager/internal/subtitles/providercore"
)

func parseCandidates(provider, fallbackLang string, data []byte) ([]providercore.Candidate, error) {
	var v any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	objects := collectObjects(v)
	out := make([]providercore.Candidate, 0, len(objects))
	for _, obj := range objects {
		urlv := firstString(obj, "download_url", "download", "url", "link", "file", "subtitle_url")
		if urlv == "" {
			continue
		}
		lang := firstString(obj, "language", "lang", "language_id", "locale")
		if lang == "" {
			lang = fallbackLang
		}
		id := firstInt(obj, "id", "file_id", "subtitle_id")
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, FileID: id, Format: formatFromURL(urlv), ReleaseName: firstString(obj, "release", "release_name", "filename", "name", "title"), DownloadCount: int(firstInt(obj, "downloads", "download_count")), SourceURL: urlv})
	}
	return out, nil
}

func collectObjects(v any) []map[string]any {
	var out []map[string]any
	switch x := v.(type) {
	case map[string]any:
		out = append(out, x)
		for _, child := range x {
			out = append(out, collectObjects(child)...)
		}
	case []any:
		for _, child := range x {
			out = append(out, collectObjects(child)...)
		}
	}
	return out
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch x := v.(type) {
			case string:
				return strings.TrimSpace(x)
			case json.Number:
				return x.String()
			}
		}
	}
	return ""
}

func firstInt(m map[string]any, keys ...string) int64 {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch x := v.(type) {
			case json.Number:
				n, _ := x.Int64()
				return n
			case float64:
				return int64(x)
			}
		}
	}
	return 0
}

func formatFromURL(raw string) string {
	ext := strings.TrimPrefix(path.Ext(raw), ".")
	if ext == "" || len(ext) > 4 {
		return "srt"
	}
	return ext
}
