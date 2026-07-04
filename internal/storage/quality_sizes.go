package storage

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

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
		setting, err := qualitySizeSetting(row)
		if err != nil {
			return nil, err
		}
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
			MinimumSizeMbPerMinute:   qualitySizeNumeric(input.MinimumSizeMBPerMinute),
			PreferredSizeMbPerMinute: optionalQualitySizeNumeric(input.PreferredSizeMBPerMinute),
			MaximumSizeMbPerMinute:   optionalQualitySizeNumeric(input.MaximumSizeMBPerMinute),
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
	preferred95 := float64(95)
	maximum100 := float64(100)
	return []QualitySizeDefinition{
		{ID: "unknown", Name: "Unknown", SortOrder: 1, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "workprint", Name: "WORKPRINT", SortOrder: 2, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "cam", Name: "CAM", SortOrder: 3, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "telesync", Name: "TELESYNC", SortOrder: 4, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "telecine", Name: "TELECINE", SortOrder: 5, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "regional", Name: "REGIONAL", SortOrder: 6, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "dvdscr", Name: "DVDSCR", SortOrder: 7, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "sdtv", Name: "SDTV", SortOrder: 8, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "dvd", Name: "DVD", SortOrder: 9, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "dvd-r", Name: "DVD-R", SortOrder: 10, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webdl-480p", Name: "WEBDL-480p", SortOrder: 11, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webrip-480p", Name: "WEBRip-480p", SortOrder: 11, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "bluray-480p", Name: "Bluray-480p", SortOrder: 12, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "bluray-576p", Name: "Bluray-576p", SortOrder: 13, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "hdtv-720p", Name: "HDTV-720p", SortOrder: 14, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webdl-720p", Name: "WEBDL-720p", SortOrder: 15, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webrip-720p", Name: "WEBRip-720p", SortOrder: 15, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "bluray-720p", Name: "Bluray-720p", SortOrder: 16, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "hdtv-1080p", Name: "HDTV-1080p", SortOrder: 17, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webdl-1080p", Name: "WEBDL-1080p", SortOrder: 18, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "webrip-1080p", Name: "WEBRip-1080p", SortOrder: 18, DefaultMaximumSizeMBPerMinute: &maximum100, DefaultPreferredSizeMBPerMinute: &preferred95},
		{ID: "bluray-1080p", Name: "Bluray-1080p", SortOrder: 19},
		{ID: "remux-1080p", Name: "Remux-1080p", SortOrder: 20},
		{ID: "hdtv-2160p", Name: "HDTV-2160p", SortOrder: 21},
		{ID: "webdl-2160p", Name: "WEBDL-2160p", SortOrder: 22},
		{ID: "webrip-2160p", Name: "WEBRip-2160p", SortOrder: 22},
		{ID: "bluray-2160p", Name: "Bluray-2160p", SortOrder: 23},
		{ID: "remux-2160p", Name: "Remux-2160p", SortOrder: 24},
		{ID: "br-disk", Name: "BR-DISK", SortOrder: 25},
		{ID: "raw-hd", Name: "Raw-HD", SortOrder: 26},
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
			MinimumSizeMbPerMinute:   qualitySizeNumeric(definition.DefaultMinimumSizeMBPerMinute),
			PreferredSizeMbPerMinute: optionalQualitySizeNumeric(definition.DefaultPreferredSizeMBPerMinute),
			MaximumSizeMbPerMinute:   optionalQualitySizeNumeric(definition.DefaultMaximumSizeMBPerMinute),
		}); err != nil {
			return err
		}
	}
	return nil
}

func qualitySizeSetting(row storagegen.AppQualitySizeSetting) (QualitySizeSetting, error) {
	minimum, err := qualitySizeFloat(row.MinimumSizeMbPerMinute)
	if err != nil {
		return QualitySizeSetting{}, err
	}
	preferred, err := optionalQualitySizeFloat(row.PreferredSizeMbPerMinute)
	if err != nil {
		return QualitySizeSetting{}, err
	}
	maximum, err := optionalQualitySizeFloat(row.MaximumSizeMbPerMinute)
	if err != nil {
		return QualitySizeSetting{}, err
	}
	return QualitySizeSetting{
		QualitySizeDefinition:    QualitySizeDefinition{ID: row.QualityID},
		MinimumSizeMBPerMinute:   minimum,
		PreferredSizeMBPerMinute: preferred,
		MaximumSizeMBPerMinute:   maximum,
		CreatedAt:                row.CreatedAt.Time,
		UpdatedAt:                row.UpdatedAt.Time,
	}, nil
}

func qualitySizeNumeric(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 100))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -2, Valid: true}
}

func optionalQualitySizeNumeric(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{}
	}
	return qualitySizeNumeric(*value)
}

func qualitySizeFloat(value pgtype.Numeric) (float64, error) {
	converted, err := value.Float64Value()
	if err != nil {
		return 0, err
	}
	if !converted.Valid {
		return 0, errors.New("quality size numeric is null")
	}
	return converted.Float64, nil
}

func optionalQualitySizeFloat(value pgtype.Numeric) (*float64, error) {
	converted, err := value.Float64Value()
	if err != nil {
		return nil, err
	}
	if !converted.Valid {
		return nil, err
	}
	return &converted.Float64, nil
}

func normalizeQualitySizeWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23514" {
		return ErrInvalidInput
	}
	return err
}
