package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) populateMediaProfile(ctx context.Context, profile *MediaProfile) error {
	qualities, err := loadMediaProfileQualities(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	videoTarget, err := loadMediaProfileVideoTarget(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	audioTargets, err := loadMediaProfileAudioTargets(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	subtitleTargets, err := loadMediaProfileSubtitleTargets(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	scores, err := loadMediaProfileCustomFormats(ctx, s.pool, profile.ID)
	if err != nil {
		return err
	}
	profile.QualityIDs = qualities
	profile.VideoTarget = videoTarget
	profile.AudioTargets = audioTargets
	profile.SubtitleTargets = subtitleTargets
	profile.CustomFormatScores = scores
	return nil
}

func loadMediaProfileQualities(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]string, error) {
	return storagegen.New(q).ListMediaProfileQualities(ctx, profileID)
}

func loadMediaProfileVideoTarget(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) (MediaProfileVideoTarget, error) {
	row, err := storagegen.New(q).GetMediaProfileVideoTarget(ctx, profileID)
	if err != nil {
		return MediaProfileVideoTarget{}, nil
	}
	return MediaProfileVideoTarget{
		Codecs:              row.Codecs,
		CodecRequired:       row.CodecRequired,
		CodecScore:          row.CodecScore,
		HDRFormats:          row.HdrFormats,
		HDRRequired:         row.HdrRequired,
		HDRScore:            row.HdrScore,
		PixelFormats:        row.PixelFormats,
		PixelFormatRequired: row.PixelFormatRequired,
		PixelFormatScore:    row.PixelFormatScore,
	}, nil
}

func loadMediaProfileAudioTargets(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]MediaProfileAudioTarget, error) {
	rows, err := storagegen.New(q).ListMediaProfileAudioTargets(ctx, profileID)
	if err != nil {
		return nil, err
	}
	targets := make([]MediaProfileAudioTarget, 0, len(rows))
	for _, row := range rows {
		targets = append(targets, MediaProfileAudioTarget{
			LanguageID:           row.LanguageID,
			Score:                row.Score,
			Required:             row.Required,
			TargetCodec:          textPtr(row.TargetCodec),
			TargetChannels:       row.TargetChannels,
			MinimumBitrateKbps:   int4Ptr(row.MinimumBitrateKbps),
			PreferredBitrateKbps: int4Ptr(row.PreferredBitrateKbps),
		})
	}
	return targets, nil
}

func loadMediaProfileSubtitleTargets(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]MediaProfileSubtitleTarget, error) {
	rows, err := storagegen.New(q).ListMediaProfileSubtitleTargets(ctx, profileID)
	if err != nil {
		return nil, err
	}
	targets := make([]MediaProfileSubtitleTarget, 0, len(rows))
	for _, row := range rows {
		targets = append(targets, MediaProfileSubtitleTarget{
			LanguageID: row.LanguageID,
			Score:      row.Score,
			Required:   row.Required,
			Source:     row.Source,
			Formats:    row.Formats,
		})
	}
	return targets, nil
}

func loadMediaProfileCustomFormats(
	ctx context.Context,
	q storagegen.DBTX,
	profileID string,
) ([]MediaProfileCustomFormatScore, error) {
	rows, err := storagegen.New(q).ListMediaProfileCustomFormats(ctx, profileID)
	if err != nil {
		return nil, err
	}
	scores := make([]MediaProfileCustomFormatScore, 0, len(rows))
	for _, row := range rows {
		scores = append(scores, MediaProfileCustomFormatScore{
			CustomFormatID: row.CustomFormatID,
			Score:          row.Score,
		})
	}
	return scores, nil
}

func replaceMediaProfileTargets(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	input MediaProfileInput,
) error {
	queries := storagegen.New(q)
	if err := queries.UpsertMediaProfileVideoTarget(ctx, storagegen.UpsertMediaProfileVideoTargetParams{
		ProfileID:           profileID,
		Codecs:              input.VideoTarget.Codecs,
		CodecRequired:       input.VideoTarget.CodecRequired,
		CodecScore:          input.VideoTarget.CodecScore,
		HdrFormats:          input.VideoTarget.HDRFormats,
		HdrRequired:         input.VideoTarget.HDRRequired,
		HdrScore:            input.VideoTarget.HDRScore,
		PixelFormats:        input.VideoTarget.PixelFormats,
		PixelFormatRequired: input.VideoTarget.PixelFormatRequired,
		PixelFormatScore:    input.VideoTarget.PixelFormatScore,
	}); err != nil {
		return err
	}
	if err := queries.ClearMediaProfileAudioTargets(ctx, profileID); err != nil {
		return err
	}
	for index, target := range input.AudioTargets {
		if err := queries.AddMediaProfileAudioTarget(ctx, storagegen.AddMediaProfileAudioTargetParams{
			ProfileID:            profileID,
			LanguageID:           target.LanguageID,
			Score:                target.Score,
			Required:             target.Required,
			TargetCodec:          textValue(target.TargetCodec),
			TargetChannels:       target.TargetChannels,
			MinimumBitrateKbps:   int4Value(target.MinimumBitrateKbps),
			PreferredBitrateKbps: int4Value(target.PreferredBitrateKbps),
			SortOrder:            int32(index),
		}); err != nil {
			return normalizeMediaProfileWriteError(err)
		}
	}
	if err := queries.ClearMediaProfileSubtitleTargets(ctx, profileID); err != nil {
		return err
	}
	for index, target := range input.SubtitleTargets {
		if err := queries.AddMediaProfileSubtitleTarget(ctx, storagegen.AddMediaProfileSubtitleTargetParams{
			ProfileID:  profileID,
			LanguageID: target.LanguageID,
			Score:      target.Score,
			Required:   target.Required,
			Source:     target.Source,
			Formats:    target.Formats,
			SortOrder:  int32(index),
		}); err != nil {
			return normalizeMediaProfileWriteError(err)
		}
	}
	return nil
}

func replaceMediaProfileCustomFormats(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	scores []MediaProfileCustomFormatScore,
) error {
	queries := storagegen.New(q)
	if err := queries.ClearMediaProfileCustomFormats(ctx, profileID); err != nil {
		return err
	}
	for _, score := range scores {
		if err := queries.AddMediaProfileCustomFormat(ctx, storagegen.AddMediaProfileCustomFormatParams{
			ProfileID:      profileID,
			CustomFormatID: score.CustomFormatID,
			Score:          score.Score,
		}); err != nil {
			return normalizeMediaProfileWriteError(err)
		}
	}
	return nil
}

func normalizeCustomFormatScores(values []MediaProfileCustomFormatScore) []MediaProfileCustomFormatScore {
	seen := map[uuid.UUID]int{}
	scores := []MediaProfileCustomFormatScore{}
	for _, value := range values {
		if value.CustomFormatID == uuid.Nil {
			continue
		}
		if index, ok := seen[value.CustomFormatID]; ok {
			scores[index].Score = value.Score
			continue
		}
		seen[value.CustomFormatID] = len(scores)
		scores = append(scores, value)
	}
	return scores
}
