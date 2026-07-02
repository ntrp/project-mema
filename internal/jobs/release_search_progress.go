package jobs

import "fmt"

type ReleaseSearchProgressEvent struct {
	Kind        string `json:"kind"`
	Message     string `json:"message"`
	IndexerName string `json:"indexerName,omitempty"`
	Query       string `json:"query,omitempty"`
	ResultCount *int   `json:"resultCount,omitempty"`
	CacheHit    *bool  `json:"cacheHit,omitempty"`
	DurationMs  *int64 `json:"durationMs,omitempty"`
}

type ReleaseSearchProgress func(event ReleaseSearchProgressEvent)

func publishReleaseSearchProgress(progress ReleaseSearchProgress, format string, args ...any) {
	if progress == nil {
		return
	}
	progress(ReleaseSearchProgressEvent{Kind: "message", Message: fmt.Sprintf(format, args...)})
}

func publishIndexerSearchStarted(progress ReleaseSearchProgress, indexerName string, query string) {
	if progress == nil {
		return
	}
	progress(ReleaseSearchProgressEvent{
		Kind:        "indexer_start",
		Message:     fmt.Sprintf("Searching %s for %q", indexerName, query),
		IndexerName: indexerName,
		Query:       query,
	})
}

func publishIndexerSearchFinished(
	progress ReleaseSearchProgress,
	indexerName string,
	query string,
	resultCount int,
	cacheHit bool,
	durationMs int64,
) {
	if progress == nil {
		return
	}
	progress(ReleaseSearchProgressEvent{
		Kind:        "indexer_finish",
		Message:     fmt.Sprintf("%s returned %d release(s)", indexerName, resultCount),
		IndexerName: indexerName,
		Query:       query,
		ResultCount: &resultCount,
		CacheHit:    &cacheHit,
		DurationMs:  &durationMs,
	})
}
