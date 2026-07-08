package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	storagegen "media-manager/internal/storage/generated"
)

type QualitySizeDefinition struct {
	ID                              string
	Name                            string
	SortOrder                       int32
	DefaultMinimumSizeMBPerMinute   float64
	DefaultPreferredSizeMBPerMinute *float64
	DefaultMaximumSizeMBPerMinute   *float64
}

type QualitySizeSetting struct {
	QualitySizeDefinition
	MinimumSizeMBPerMinute   float64
	PreferredSizeMBPerMinute *float64
	MaximumSizeMBPerMinute   *float64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type QualitySizeSettingInput struct {
	QualityID                string
	MinimumSizeMBPerMinute   float64
	PreferredSizeMBPerMinute *float64
	MaximumSizeMBPerMinute   *float64
}

func (s *SettingsStore) ListQualitySizeSettings(ctx context.Context) ([]QualitySizeSetting, error) {
	if err := s.ensureQualitySizeSettings(ctx); err != nil {
		return nil, err
	}

	rows, err := storagegen.New(s.pool).ListQualitySizeSettings(ctx)
	if err != nil {
		return nil, err
	}

	settingsByID := map[string]QualitySizeSetting{}
	for _, row := range rows {
		setting := qualitySizeSettingFromRow(row)
		settingsByID[setting.ID] = setting
	}

	definitions := QualitySizeDefinitions()
	settings := make([]QualitySizeSetting, 0, len(definitions))
	for _, definition := range definitions {
		setting := settingsByID[definition.ID]
		setting.QualitySizeDefinition = definition
		settings = append(settings, setting)
	}
	return settings, nil
}

func (s *SettingsStore) SaveQualitySizeSettings(
	ctx context.Context,
	inputs []QualitySizeSettingInput,
) ([]QualitySizeSetting, error) {
	definitions := QualitySizeDefinitionMap()
	for _, input := range inputs {
		if _, ok := definitions[input.QualityID]; !ok {
			return nil, ErrInvalidInput
		}
		if input.MinimumSizeMBPerMinute < 0 {
			return nil, ErrInvalidInput
		}
		if input.PreferredSizeMBPerMinute != nil && *input.PreferredSizeMBPerMinute < input.MinimumSizeMBPerMinute {
			return nil, ErrInvalidInput
		}
		if input.MaximumSizeMBPerMinute != nil && *input.MaximumSizeMBPerMinute < input.MinimumSizeMBPerMinute {
			return nil, ErrInvalidInput
		}
		if input.PreferredSizeMBPerMinute != nil && input.MaximumSizeMBPerMinute != nil &&
			*input.PreferredSizeMBPerMinute > *input.MaximumSizeMBPerMinute {
			return nil, ErrInvalidInput
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	queries := storagegen.New(s.pool).WithTx(tx)
	for _, input := range inputs {
		if err := queries.UpsertQualitySizeSetting(ctx, storagegen.UpsertQualitySizeSettingParams{
			QualityID:                input.QualityID,
			MinimumSizeMbPerMinute:   input.MinimumSizeMBPerMinute,
			PreferredSizeMbPerMinute: input.PreferredSizeMBPerMinute,
			MaximumSizeMbPerMinute:   input.MaximumSizeMBPerMinute,
		}); err != nil {
			return nil, normalizeQualitySizeWriteError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, normalizeQualitySizeWriteError(err)
	}
	return s.ListQualitySizeSettings(ctx)
}

func QualitySizeDefinitions() []QualitySizeDefinition {
	return []QualitySizeDefinition{
		qualitySizeDefinition("unknown", "Unknown", 1, 0, 45, 180),
		qualitySizeDefinition("workprint", "WORKPRINT", 2, 0, 18, 60),
		qualitySizeDefinition("cam", "CAM", 3, 0, 8, 25),
		qualitySizeDefinition("telesync", "TELESYNC", 4, 0, 10, 35),
		qualitySizeDefinition("telecine", "TELECINE", 5, 4, 14, 45),
		qualitySizeDefinition("regional", "REGIONAL", 6, 6, 18, 55),
		qualitySizeDefinition("dvdscr", "DVDSCR", 7, 5, 16, 50),
		qualitySizeDefinition("sdtv", "SDTV", 8, 4, 10, 22),
		qualitySizeDefinition("dvd", "DVD", 9, 6, 14, 30),
		qualitySizeDefinition("dvd-r", "DVD-R", 10, 10, 22, 45),
		qualitySizeDefinition("webdl-480p", "WEBDL-480p", 11, 6, 14, 28),
		qualitySizeDefinition("webrip-480p", "WEBRip-480p", 11, 5, 12, 26),
		qualitySizeDefinition("bluray-480p", "Bluray-480p", 12, 8, 18, 36),
		qualitySizeDefinition("bluray-576p", "Bluray-576p", 13, 10, 22, 44),
		qualitySizeDefinition("hdtv-720p", "HDTV-720p", 14, 10, 22, 42),
		qualitySizeDefinition("webdl-720p", "WEBDL-720p", 15, 13, 28, 54),
		qualitySizeDefinition("webrip-720p", "WEBRip-720p", 15, 12, 26, 50),
		qualitySizeDefinition("bluray-720p", "Bluray-720p", 16, 18, 38, 72),
		qualitySizeDefinition("hdtv-1080p", "HDTV-1080p", 17, 18, 38, 72),
		qualitySizeDefinition("webdl-1080p", "WEBDL-1080p", 18, 24, 52, 96),
		qualitySizeDefinition("webrip-1080p", "WEBRip-1080p", 18, 22, 48, 90),
		qualitySizeDefinition("bluray-1080p", "Bluray-1080p", 19, 36, 75, 135),
		qualitySizeDefinition("remux-1080p", "Remux-1080p", 20, 85, 150, 260),
		qualitySizeDefinition("hdtv-2160p", "HDTV-2160p", 21, 60, 130, 260),
		qualitySizeDefinition("webdl-2160p", "WEBDL-2160p", 22, 75, 160, 320),
		qualitySizeDefinition("webrip-2160p", "WEBRip-2160p", 22, 70, 150, 300),
		qualitySizeDefinition("bluray-2160p", "Bluray-2160p", 23, 110, 230, 420),
		qualitySizeDefinition("remux-2160p", "Remux-2160p", 24, 170, 330, 620),
		qualitySizeDefinition("br-disk", "BR-DISK", 25, 200, 400, 760),
		qualitySizeDefinition("raw-hd", "Raw-HD", 26, 220, 450, 850),
	}
}

func qualitySizeDefinition(id string, name string, sortOrder int32, minimum float64, preferred float64, maximum float64) QualitySizeDefinition {
	return QualitySizeDefinition{
		ID:                              id,
		Name:                            name,
		SortOrder:                       sortOrder,
		DefaultMinimumSizeMBPerMinute:   minimum,
		DefaultPreferredSizeMBPerMinute: &preferred,
		DefaultMaximumSizeMBPerMinute:   &maximum,
	}
}

func QualitySizeDefinitionMap() map[string]QualitySizeDefinition {
	definitions := QualitySizeDefinitions()
	byID := make(map[string]QualitySizeDefinition, len(definitions))
	for _, definition := range definitions {
		byID[definition.ID] = definition
	}
	return byID
}

func (s *SettingsStore) ensureQualitySizeSettings(ctx context.Context) error {
	queries := storagegen.New(s.pool)
	for _, definition := range QualitySizeDefinitions() {
		if err := queries.EnsureQualitySizeSetting(ctx, storagegen.EnsureQualitySizeSettingParams{
			QualityID:                definition.ID,
			MinimumSizeMbPerMinute:   definition.DefaultMinimumSizeMBPerMinute,
			PreferredSizeMbPerMinute: definition.DefaultPreferredSizeMBPerMinute,
			MaximumSizeMbPerMinute:   definition.DefaultMaximumSizeMBPerMinute,
		}); err != nil {
			return err
		}
	}
	return nil
}

func qualitySizeSettingFromRow(row storagegen.AppQualitySizeSetting) QualitySizeSetting {
	return QualitySizeSetting{
		QualitySizeDefinition:    QualitySizeDefinition{ID: row.QualityID},
		MinimumSizeMBPerMinute:   row.MinimumSizeMbPerMinute,
		PreferredSizeMBPerMinute: row.PreferredSizeMbPerMinute,
		MaximumSizeMBPerMinute:   row.MaximumSizeMbPerMinute,
		CreatedAt:                row.CreatedAt,
		UpdatedAt:                row.UpdatedAt,
	}
}

func normalizeQualitySizeWriteError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23514" {
		return ErrInvalidInput
	}
	return err
}
