package storage

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) UpsertMediaFileFact(ctx context.Context, input MediaFileFactInput) (MediaFileFact, error) {
	return upsertMediaFileFact(ctx, s.pool, input)
}

func (s *SettingsStore) ListMediaFileFacts(ctx context.Context, mediaItemID uuid.UUID) ([]MediaFileFact, error) {
	return listMediaFileFacts(ctx, s.pool, mediaItemID)
}

func upsertMediaFileFact(ctx context.Context, q storagegen.DBTX, input MediaFileFactInput) (MediaFileFact, error) {
	row, err := storagegen.New(q).UpsertMediaFileFact(ctx, mediaFileFactParams(input))
	if err != nil {
		return MediaFileFact{}, err
	}
	if err := storagegen.New(q).DeleteMediaFileTracksForFact(ctx, row.ID); err != nil {
		return MediaFileFact{}, err
	}
	fact := mediaFileFactFromRow(row)
	for _, track := range input.Tracks {
		trackRow, err := storagegen.New(q).InsertMediaFileTrack(ctx, mediaFileTrackParams(fact, track))
		if err != nil {
			return MediaFileFact{}, err
		}
		fact.Tracks = append(fact.Tracks, mediaFileTrackFromRow(trackRow))
	}
	return fact, nil
}

func listMediaFileFacts(ctx context.Context, q storagegen.DBTX, mediaItemID uuid.UUID) ([]MediaFileFact, error) {
	rows, err := storagegen.New(q).ListMediaFileFactsForItem(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	tracks, err := storagegen.New(q).ListMediaFileTracksForItem(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	items := make([]MediaFileFact, 0, len(rows))
	byID := map[uuid.UUID]int{}
	for _, row := range rows {
		byID[row.ID] = len(items)
		items = append(items, mediaFileFactFromRow(row))
	}
	for _, track := range tracks {
		index, ok := byID[track.MediaFileFactID]
		if ok {
			items[index].Tracks = append(items[index].Tracks, mediaFileTrackFromRow(track))
		}
	}
	return items, nil
}

func recordMediaFileFactFromPath(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	seasonID *uuid.UUID,
	episodeID *uuid.UUID,
	filePath string,
	sourceKind string,
) error {
	var size *int64
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		value := info.Size()
		size = &value
	}
	_, err := upsertMediaFileFact(ctx, q, MediaFileFactInput{
		MediaItemID: mediaItemID,
		SeasonID:    seasonID,
		EpisodeID:   episodeID,
		FilePath:    absoluteCleanPathOrClean(filePath),
		SizeBytes:   size,
		SourceKind:  sourceKind,
		ProbedAt:    time.Now().UTC(),
	})
	return err
}

func mediaFileFactParams(input MediaFileFactInput) storagegen.UpsertMediaFileFactParams {
	probedAt := input.ProbedAt
	if probedAt.IsZero() {
		probedAt = time.Now().UTC()
	}
	return storagegen.UpsertMediaFileFactParams{
		ID:                  uuid.New(),
		MediaItemID:         input.MediaItemID,
		SeasonID:            input.SeasonID,
		EpisodeID:           input.EpisodeID,
		FilePath:            absoluteCleanPathOrClean(input.FilePath),
		QualityID:           textValue(input.QualityID),
		ContainerFormat:     textValue(input.ContainerFormat),
		ContainerFormatName: textValue(input.ContainerFormatName),
		ContainerBitrate:    int8Value(input.ContainerBitrate),
		DurationMs:          int8Value(input.DurationMs),
		SizeBytes:           int8Value(input.SizeBytes),
		SourceKind:          normalizedMediaFileFactSource(input.SourceKind),
		ProbedAt:            probedAt,
	}
}

func mediaFileTrackParams(fact MediaFileFact, input MediaFileTrackFactInput) storagegen.InsertMediaFileTrackParams {
	return storagegen.InsertMediaFileTrackParams{
		ID:              uuid.New(),
		MediaFileFactID: fact.ID,
		MediaItemID:     fact.MediaItemID,
		FilePath:        fact.FilePath,
		StreamIndex:     input.StreamIndex,
		TrackType:       strings.TrimSpace(input.TrackType),
		LanguageID:      textValue(input.LanguageID),
		Codec:           textValue(input.Codec),
		Channels:        textValue(input.Channels),
		DurationMs:      int8Value(input.DurationMs),
		BitrateKbps:     int4Value(input.BitrateKbps),
		Width:           int4Value(input.Width),
		Height:          int4Value(input.Height),
		HdrFormat:       textValue(input.HDRFormat),
		PixelFormat:     textValue(input.PixelFormat),
		BitDepth:        int4Value(input.BitDepth),
		Format:          textValue(input.Format),
		Title:           textValue(input.Title),
		Disposition:     jsonObject(input.Disposition),
	}
}

func mediaFileFactFromRow(row storagegen.AppMediaFileFact) MediaFileFact {
	return MediaFileFact{
		ID:                  row.ID,
		MediaItemID:         row.MediaItemID,
		SeasonID:            row.SeasonID,
		EpisodeID:           row.EpisodeID,
		FilePath:            row.FilePath,
		QualityID:           textPtr(row.QualityID),
		ContainerFormat:     textPtr(row.ContainerFormat),
		ContainerFormatName: textPtr(row.ContainerFormatName),
		ContainerBitrate:    int8Ptr(row.ContainerBitrate),
		DurationMs:          int8Ptr(row.DurationMs),
		SizeBytes:           int8Ptr(row.SizeBytes),
		SourceKind:          row.SourceKind,
		ProbedAt:            row.ProbedAt,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func mediaFileTrackFromRow(row storagegen.AppMediaFileTrack) MediaFileTrackFact {
	return MediaFileTrackFact{
		ID:              row.ID,
		MediaFileFactID: row.MediaFileFactID,
		MediaItemID:     row.MediaItemID,
		FilePath:        row.FilePath,
		StreamIndex:     row.StreamIndex,
		TrackType:       row.TrackType,
		LanguageID:      textPtr(row.LanguageID),
		Codec:           textPtr(row.Codec),
		Channels:        textPtr(row.Channels),
		DurationMs:      int8Ptr(row.DurationMs),
		BitrateKbps:     int4Ptr(row.BitrateKbps),
		Width:           int4Ptr(row.Width),
		Height:          int4Ptr(row.Height),
		HDRFormat:       textPtr(row.HdrFormat),
		PixelFormat:     textPtr(row.PixelFormat),
		BitDepth:        int4Ptr(row.BitDepth),
		Format:          textPtr(row.Format),
		Title:           textPtr(row.Title),
		Disposition:     jsonMap(row.Disposition),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func normalizedMediaFileFactSource(value string) string {
	switch strings.TrimSpace(value) {
	case "import", "probe", "manual":
		return strings.TrimSpace(value)
	default:
		return "rescan"
	}
}
