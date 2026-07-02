package jobs

import (
	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func shouldFallbackEpisodeToSeason(
	item storage.MediaItem,
	criteria decisions.ReleaseSearchCriteria,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) bool {
	if criteria.Kind != "episode" || criteria.SeasonNumber == nil {
		return false
	}
	matches := newReleaseMatchCache(item, profile, formats, languages)
	for _, release := range releases {
		if matches.match(release).Severity != "error" {
			return false
		}
	}
	return true
}

func seasonFallbackCriteria(criteria decisions.ReleaseSearchCriteria) decisions.ReleaseSearchCriteria {
	return decisions.ReleaseSearchCriteria{
		Kind:         "season",
		Title:        criteria.Title,
		Year:         criteria.Year,
		SeasonNumber: criteria.SeasonNumber,
	}
}
