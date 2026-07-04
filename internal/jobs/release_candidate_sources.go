package jobs

import (
	"strings"

	"media-manager/internal/storage"
)

func appendUniqueReleaseSources(
	left []storage.ReleaseCandidateSource,
	right []storage.ReleaseCandidateSource,
) []storage.ReleaseCandidateSource {
	merged := append([]storage.ReleaseCandidateSource(nil), left...)
	seen := map[string]struct{}{}
	for _, source := range merged {
		seen[releaseSourceKey(source)] = struct{}{}
	}
	for _, source := range right {
		key := releaseSourceKey(source)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, source)
	}
	return merged
}

func releaseSourceKey(source storage.ReleaseCandidateSource) string {
	indexerID := ""
	if source.IndexerID != nil {
		indexerID = source.IndexerID.String()
	}
	values := []string{
		indexerID,
		strings.ToLower(strings.TrimSpace(source.IndexerName)),
		strings.ToLower(strings.TrimSpace(source.IndexerProtocol)),
		strings.ToLower(strings.TrimSpace(source.Title)),
		strings.ToLower(strings.TrimSpace(source.DownloadURL)),
	}
	if source.InfoURL != nil {
		values = append(values, strings.ToLower(strings.TrimSpace(*source.InfoURL)))
	}
	if source.GUID != nil {
		values = append(values, strings.ToLower(strings.TrimSpace(*source.GUID)))
	}
	return strings.Join(values, "\x00")
}
