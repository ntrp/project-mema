package storage

import "strings"

func ReleaseCandidateSourcesForInput(release ReleaseCandidateInput) []ReleaseCandidateSource {
	if len(release.Sources) > 0 {
		return append([]ReleaseCandidateSource(nil), release.Sources...)
	}
	return []ReleaseCandidateSource{{
		IndexerID:       release.IndexerID,
		IndexerName:     release.IndexerName,
		IndexerProtocol: release.IndexerProtocol,
		Title:           release.Title,
		DownloadURL:     release.DownloadURL,
		InfoURL:         release.InfoURL,
		GUID:            release.GUID,
	}}
}

func ReleaseCandidateSourcesForStored(release ReleaseCandidate) []ReleaseCandidateSource {
	if len(release.Sources) > 0 {
		return append([]ReleaseCandidateSource(nil), release.Sources...)
	}
	return []ReleaseCandidateSource{{
		IndexerID:       release.IndexerID,
		IndexerName:     release.IndexerName,
		IndexerProtocol: release.IndexerProtocol,
		Title:           release.Title,
		DownloadURL:     release.DownloadURL,
		InfoURL:         release.InfoURL,
		GUID:            release.GUID,
	}}
}

func appendReleaseCandidateSources(
	left []ReleaseCandidateSource,
	right []ReleaseCandidateSource,
) []ReleaseCandidateSource {
	merged := append([]ReleaseCandidateSource(nil), left...)
	seen := map[string]struct{}{}
	for _, source := range merged {
		seen[releaseCandidateSourceKey(source)] = struct{}{}
	}
	for _, source := range right {
		key := releaseCandidateSourceKey(source)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, source)
	}
	return merged
}

func releaseCandidateSourceKey(source ReleaseCandidateSource) string {
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
