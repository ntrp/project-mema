package storage

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
