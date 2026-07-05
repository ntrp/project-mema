package jobs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

type releaseSearchBranch struct {
	criteria decisions.ReleaseSearchCriteria
	queries  []string
}

var seasonNumberPattern = regexp.MustCompile(`\d+`)

func releaseSearchBranches(
	item storage.MediaItem,
	criteria decisions.ReleaseSearchCriteria,
	query string,
) []releaseSearchBranch {
	if item.Type != "serie" || criteria.Kind != "serie" {
		return []releaseSearchBranch{{criteria: criteria, queries: decisions.SearchQueriesForCriteria(criteria, query)}}
	}
	seasons := monitoredSeasons(item.Seasons)
	if len(seasons) > 0 {
		branches := make([]releaseSearchBranch, 0, len(seasons))
		for _, season := range seasons {
			seasonNumber := season.number
			seasonCriteria := decisions.ReleaseSearchCriteria{
				Kind:         "season",
				Title:        item.Title,
				SeasonID:     season.id,
				SeasonNumber: &seasonNumber,
			}
			branches = append(branches, releaseSearchBranch{
				criteria: seasonCriteria,
				queries:  decisions.SearchQueriesForCriteria(seasonCriteria, ""),
			})
		}
		return branches
	}
	episodes := monitoredAiredEpisodes(item.Seasons)
	if len(episodes) == 0 {
		return []releaseSearchBranch{{criteria: criteria, queries: decisions.SearchQueriesForCriteria(criteria, query)}}
	}
	branches := make([]releaseSearchBranch, 0, len(episodes))
	for _, episode := range episodes {
		seasonNumber := episode.season
		episodeNumber := episode.episode
		episodeCriteria := decisions.ReleaseSearchCriteria{
			Kind:          "episode",
			Title:         item.Title,
			SeasonID:      episode.seasonID,
			EpisodeID:     episode.episodeID,
			SeasonNumber:  &seasonNumber,
			EpisodeNumber: &episodeNumber,
		}
		queries := decisions.SearchQueriesForCriteria(episodeCriteria, "")
		if episode.title != "" {
			queries = append(queries, fmt.Sprintf("%s %s", item.Title, episode.title))
		}
		branches = append(branches, releaseSearchBranch{criteria: episodeCriteria, queries: queries})
	}
	return branches
}

type monitoredSeason struct {
	id     *uuid.UUID
	number int32
}

func monitoredSeasons(seasons []storage.MediaSeason) []monitoredSeason {
	values := []monitoredSeason{}
	for _, season := range seasons {
		number, ok := seasonNumber(season)
		if !ok || !season.Monitored {
			continue
		}
		values = append(values, monitoredSeason{id: season.ID, number: number})
	}
	return values
}

type monitoredEpisode struct {
	seasonID  *uuid.UUID
	episodeID *uuid.UUID
	season    int32
	episode   int32
	title     string
}

func monitoredAiredEpisodes(seasons []storage.MediaSeason) []monitoredEpisode {
	values := []monitoredEpisode{}
	for _, season := range seasons {
		number, ok := seasonNumber(season)
		if !ok {
			continue
		}
		for _, episode := range season.Episodes {
			if !episode.Monitored || !episodeAired(episode) {
				continue
			}
			values = append(values, monitoredEpisode{
				seasonID:  season.ID,
				episodeID: episode.ID,
				season:    number,
				episode:   episode.EpisodeNumber,
				title:     episode.Name,
			})
		}
	}
	return values
}

func seasonNumber(season storage.MediaSeason) (int32, bool) {
	if season.SeasonNumber != 0 {
		return season.SeasonNumber, true
	}
	if strings.EqualFold(strings.TrimSpace(season.Name), "specials") {
		return 0, true
	}
	match := seasonNumberPattern.FindString(season.Name)
	if match == "" {
		return 0, false
	}
	value, err := strconv.ParseInt(match, 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(value), true
}

func episodeAired(episode storage.MediaEpisode) bool {
	if episode.AirDate == nil || *episode.AirDate == "" {
		return true
	}
	value, err := time.Parse("2006-01-02", *episode.AirDate)
	return err != nil || !value.After(time.Now())
}
