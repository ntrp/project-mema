package storage

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

var seasonNumberPattern = regexp.MustCompile(`(?i)(?:season|series)\s*(\d+)$`)

func materializeMediaSeriesSnapshot(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	input MediaItemInput,
) error {
	if input.Type != "serie" || len(input.Seasons) == 0 {
		return nil
	}
	return upsertMediaSeriesSeasons(ctx, q, mediaItemID, mediaSeriesInputsFromSnapshot(input), true)
}

func mediaSeriesInputsFromSnapshot(input MediaItemInput) []MediaSeriesSeasonInput {
	items := make([]MediaSeriesSeasonInput, 0, len(input.Seasons))
	for index, season := range input.Seasons {
		seasonNumber := mediaSeasonSnapshotNumber(index, season)
		items = append(items, MediaSeriesSeasonInput{
			ExternalProvider: input.ExternalProvider,
			ExternalID:       input.ExternalID,
			SeasonNumber:     seasonNumber,
			Name:             season.Name,
			Overview:         nil,
			AirDate:          season.AirDate,
			PosterPath:       season.PosterPath,
			EpisodeCount:     season.EpisodeCount,
			Monitored:        season.Monitored,
			Source:           mediaSeriesSource(input),
			Episodes:         mediaSeriesEpisodeInputs(season, input),
		})
	}
	return items
}

func mediaSeriesEpisodeInputs(season MediaSeason, input MediaItemInput) []MediaSeriesEpisodeInput {
	items := make([]MediaSeriesEpisodeInput, 0, len(season.Episodes))
	for _, episode := range season.Episodes {
		items = append(items, MediaSeriesEpisodeInput{
			ExternalProvider: input.ExternalProvider,
			ExternalID:       input.ExternalID,
			EpisodeNumber:    episode.EpisodeNumber,
			Name:             episode.Name,
			Overview:         episode.Overview,
			AirDate:          episode.AirDate,
			StillPath:        episode.StillPath,
			Monitored:        episode.Monitored,
			Source:           mediaSeriesSource(input),
		})
	}
	return items
}

func mediaSeasonSnapshotNumber(index int, season MediaSeason) int32 {
	if season.SeasonNumber != 0 {
		return season.SeasonNumber
	}
	name := strings.TrimSpace(strings.ToLower(season.Name))
	if name == "specials" {
		return 0
	}
	matches := seasonNumberPattern.FindStringSubmatch(strings.TrimSpace(season.Name))
	if len(matches) == 2 {
		if value, err := strconv.Atoi(matches[1]); err == nil {
			return int32(value)
		}
	}
	return int32(index + 1)
}

func mediaSeriesSource(input MediaItemInput) map[string]any {
	source := map[string]any{}
	if input.ExternalProvider != nil {
		source["provider"] = *input.ExternalProvider
	}
	if input.ExternalID != nil {
		source["externalId"] = *input.ExternalID
	}
	return source
}

func preserveMediaSeriesMonitorState(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	input []MediaSeriesSeasonInput,
) ([]MediaSeriesSeasonInput, error) {
	existing, err := listMediaSeriesSeasons(ctx, q, mediaItemID)
	if err != nil {
		return nil, err
	}
	bySeason := make(map[int32]MediaSeriesSeason, len(existing))
	byEpisode := map[[2]int32]MediaSeriesEpisode{}
	for _, season := range existing {
		bySeason[season.SeasonNumber] = season
		for _, episode := range season.Episodes {
			byEpisode[[2]int32{season.SeasonNumber, episode.EpisodeNumber}] = episode
		}
	}
	for seasonIndex := range input {
		season, ok := bySeason[input[seasonIndex].SeasonNumber]
		if ok {
			input[seasonIndex].Monitored = season.Monitored
		}
		for episodeIndex := range input[seasonIndex].Episodes {
			key := [2]int32{input[seasonIndex].SeasonNumber, input[seasonIndex].Episodes[episodeIndex].EpisodeNumber}
			if episode, ok := byEpisode[key]; ok {
				input[seasonIndex].Episodes[episodeIndex].Monitored = episode.Monitored
			}
		}
	}
	return input, nil
}

func applyMediaSeriesMonitorSnapshot(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	seasons []MediaSeason,
) error {
	existing, err := listMediaSeriesSeasons(ctx, q, mediaItemID)
	if err != nil {
		return err
	}
	bySeason := make(map[int32]MediaSeriesSeason, len(existing))
	byEpisode := map[[2]int32]MediaSeriesEpisode{}
	for _, season := range existing {
		bySeason[season.SeasonNumber] = season
		for _, episode := range season.Episodes {
			byEpisode[[2]int32{season.SeasonNumber, episode.EpisodeNumber}] = episode
		}
	}
	queries := storagegen.New(q)
	for index, season := range seasons {
		seasonNumber := mediaSeasonSnapshotNumber(index, season)
		relational, ok := bySeason[seasonNumber]
		if !ok {
			continue
		}
		if len(season.Episodes) == 0 {
			if _, err := queries.UpdateMediaSeasonMonitoredRow(ctx, storagegen.UpdateMediaSeasonMonitoredRowParams{
				ID:        relational.ID,
				Monitored: season.Monitored,
			}); err != nil {
				return err
			}
			if err := queries.UpdateMediaSeasonEpisodesMonitored(ctx, storagegen.UpdateMediaSeasonEpisodesMonitoredParams{
				SeasonID:  relational.ID,
				Monitored: season.Monitored,
			}); err != nil {
				return err
			}
			continue
		}
		for _, episode := range season.Episodes {
			key := [2]int32{seasonNumber, episode.EpisodeNumber}
			if row, ok := byEpisode[key]; ok {
				if _, err := queries.UpdateMediaEpisodeMonitoredRow(ctx, storagegen.UpdateMediaEpisodeMonitoredRowParams{
					ID:        row.ID,
					Monitored: episode.Monitored,
				}); err != nil {
					return err
				}
			}
		}
		if _, err := queries.SyncMediaSeasonMonitoredFromEpisodes(ctx, relational.ID); err != nil {
			return err
		}
	}
	return nil
}

func hydrateMediaItemSeries(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	seasons, err := listMediaSeriesSeasons(ctx, q, item.ID)
	if err != nil || len(seasons) == 0 {
		return item, err
	}
	item.Seasons = mediaSeasonsFromSeriesRows(seasons)
	return item, nil
}

func hydrateMediaItemsSeries(
	ctx context.Context,
	q storagegen.DBTX,
	items []MediaItem,
) ([]MediaItem, error) {
	for index := range items {
		item, err := hydrateMediaItemSeries(ctx, q, items[index])
		if err != nil {
			return nil, err
		}
		items[index] = item
	}
	return items, nil
}

func mediaSeasonsFromSeriesRows(seasons []MediaSeriesSeason) []MediaSeason {
	items := make([]MediaSeason, 0, len(seasons))
	for _, season := range seasons {
		items = append(items, MediaSeason{
			Name:         season.Name,
			SeasonNumber: season.SeasonNumber,
			EpisodeCount: season.EpisodeCount,
			AirDate:      season.AirDate,
			PosterPath:   season.PosterPath,
			Monitored:    season.Monitored,
			Episodes:     mediaEpisodesFromSeriesRows(season.Episodes),
		})
	}
	return items
}

func mediaEpisodesFromSeriesRows(episodes []MediaSeriesEpisode) []MediaEpisode {
	items := make([]MediaEpisode, 0, len(episodes))
	for _, episode := range episodes {
		items = append(items, MediaEpisode{
			Name:          episode.Name,
			EpisodeNumber: episode.EpisodeNumber,
			Overview:      episode.Overview,
			AirDate:       episode.AirDate,
			StillPath:     episode.StillPath,
			Monitored:     episode.Monitored,
		})
	}
	return items
}
