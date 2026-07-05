package storage

import (
	"context"
	"encoding/json"
	"errors"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) ListMediaSeriesSeasons(ctx context.Context, mediaItemID uuid.UUID) ([]MediaSeriesSeason, error) {
	return listMediaSeriesSeasons(ctx, s.pool, mediaItemID)
}

func listMediaSeriesSeasons(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
) ([]MediaSeriesSeason, error) {
	queries := storagegen.New(q)
	seasonRows, err := queries.ListMediaSeasonRows(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	episodeRows, err := queries.ListMediaEpisodeRows(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	episodesBySeason := make(map[uuid.UUID][]MediaSeriesEpisode, len(seasonRows))
	for _, row := range episodeRows {
		episode := mediaSeriesEpisodeFromRow(row)
		episodesBySeason[episode.SeasonID] = append(episodesBySeason[episode.SeasonID], episode)
	}
	seasons := make([]MediaSeriesSeason, 0, len(seasonRows))
	for _, row := range seasonRows {
		season := mediaSeriesSeasonFromRow(row)
		season.Episodes = episodesBySeason[season.ID]
		seasons = append(seasons, season)
	}
	return seasons, nil
}

func (s *SettingsStore) UpsertMediaSeriesSeasons(
	ctx context.Context,
	mediaItemID uuid.UUID,
	input []MediaSeriesSeasonInput,
) ([]MediaSeriesSeason, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err := upsertMediaSeriesSeasons(ctx, tx, mediaItemID, input, false); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.ListMediaSeriesSeasons(ctx, mediaItemID)
}

func upsertMediaSeriesSeasons(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	input []MediaSeriesSeasonInput,
	preserveMonitor bool,
) error {
	if preserveMonitor {
		var err error
		input, err = preserveMediaSeriesMonitorState(ctx, q, mediaItemID, input)
		if err != nil {
			return err
		}
	}
	queries := storagegen.New(q)
	for _, season := range input {
		params, err := mediaSeriesSeasonParams(mediaItemID, season)
		if err != nil {
			return err
		}
		row, err := queries.UpsertMediaSeasonRow(ctx, params)
		if err != nil {
			return err
		}
		for _, episode := range season.Episodes {
			params, err := mediaSeriesEpisodeParams(mediaItemID, row.ID, season.SeasonNumber, episode)
			if err != nil {
				return err
			}
			if _, err := queries.UpsertMediaEpisodeRow(ctx, params); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SettingsStore) SetMediaSeriesSeasonMonitored(
	ctx context.Context,
	seasonID uuid.UUID,
	monitored bool,
) (MediaSeriesSeason, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaSeriesSeason{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(tx)
	row, err := queries.UpdateMediaSeasonMonitoredRow(ctx, storagegen.UpdateMediaSeasonMonitoredRowParams{
		ID:        seasonID,
		Monitored: monitored,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaSeriesSeason{}, ErrNotFound
	}
	if err != nil {
		return MediaSeriesSeason{}, err
	}
	if err := queries.UpdateMediaSeasonEpisodesMonitored(ctx, storagegen.UpdateMediaSeasonEpisodesMonitoredParams{
		SeasonID:  seasonID,
		Monitored: monitored,
	}); err != nil {
		return MediaSeriesSeason{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaSeriesSeason{}, err
	}
	return mediaSeriesSeasonFromRow(row), nil
}

func (s *SettingsStore) SetMediaSeriesEpisodeMonitored(
	ctx context.Context,
	episodeID uuid.UUID,
	monitored bool,
) (MediaSeriesEpisode, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaSeriesEpisode{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(tx)
	row, err := queries.UpdateMediaEpisodeMonitoredRow(ctx, storagegen.UpdateMediaEpisodeMonitoredRowParams{
		ID:        episodeID,
		Monitored: monitored,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaSeriesEpisode{}, ErrNotFound
	}
	if err != nil {
		return MediaSeriesEpisode{}, err
	}
	if _, err := queries.SyncMediaSeasonMonitoredFromEpisodes(ctx, row.SeasonID); err != nil {
		return MediaSeriesEpisode{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaSeriesEpisode{}, err
	}
	return mediaSeriesEpisodeFromRow(row), nil
}

func mediaSeriesSourcePayload(source map[string]any) ([]byte, error) {
	if source == nil {
		source = map[string]any{}
	}
	return json.Marshal(source)
}
