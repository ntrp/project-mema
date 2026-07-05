package storage

import (
	"context"
	"strings"
)

func (s *SettingsStore) ListEligibleIndexers(ctx context.Context, item MediaItem) ([]Indexer, error) {
	indexers, err := s.ListEnabledIndexers(ctx)
	if err != nil {
		return nil, err
	}
	return EligibleIndexers(indexers, item), nil
}

func EligibleIndexers(indexers []Indexer, item MediaItem) []Indexer {
	eligible := []Indexer{}
	for _, indexer := range indexers {
		if !indexerMatchesMediaType(indexer, item.Type) {
			continue
		}
		if !indexerMatchesTags(indexer, item.Tags) {
			continue
		}
		eligible = append(eligible, indexer)
	}
	return eligible
}

func indexerMatchesMediaType(indexer Indexer, mediaType string) bool {
	mediaType = strings.TrimSpace(mediaType)
	for _, scope := range normalizeIndexerMediaTypeScopes(indexer.MediaTypeScopes) {
		if scope == mediaType {
			return true
		}
	}
	return false
}

func indexerMatchesTags(indexer Indexer, tags []string) bool {
	mediaTags := normalizedTagSet(tags)
	if len(mediaTags) == 0 {
		return true
	}
	indexerTags := normalizedTagSet(indexer.TagScopes)
	if len(indexerTags) == 0 {
		return false
	}
	for tag := range mediaTags {
		if indexerTags[tag] {
			return true
		}
	}
	return false
}

func normalizedTagSet(values []string) map[string]bool {
	set := map[string]bool{}
	for _, value := range normalizeTagNames(values) {
		set[strings.ToLower(value)] = true
	}
	return set
}
