package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/jackc/pgx/v5"
)

type FileNamingSettings struct {
	MovieFileFormat      string
	MovieFolderFormat    string
	SeriesEpisodeFormat  string
	DailyEpisodeFormat   string
	AnimeEpisodeFormat   string
	SeriesFolderFormat   string
	SeasonFolderFormat   string
	SpecialsFolderFormat string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type FileNamingSettingsInput struct {
	MovieFileFormat      string
	MovieFolderFormat    string
	SeriesEpisodeFormat  string
	DailyEpisodeFormat   string
	AnimeEpisodeFormat   string
	SeriesFolderFormat   string
	SeasonFolderFormat   string
	SpecialsFolderFormat string
}

var fileNamingTokenPattern = regexp.MustCompile(`\{([^{}]+)\}`)

func DefaultFileNamingSettings() FileNamingSettings {
	return FileNamingSettings{
		MovieFileFormat:      "{movie_title} ({release_year}) {quality_full}",
		MovieFolderFormat:    "{movie_title} ({release_year})",
		SeriesEpisodeFormat:  "{series_title} - S{season:00}E{episode:00} - {episode_title} {quality_full}",
		DailyEpisodeFormat:   "{series_title} - {air_date} - {episode_title} {quality_full}",
		AnimeEpisodeFormat:   "{series_title} - S{season:00}E{episode:00} - {episode_title} {quality_full}",
		SeriesFolderFormat:   "{series_title}",
		SeasonFolderFormat:   "Season {season}",
		SpecialsFolderFormat: "Specials",
	}
}

func (s *SettingsStore) EnsureDefaultFileNamingSettings(ctx context.Context) error {
	defaults := DefaultFileNamingSettings()
	return storagegen.New(s.pool).EnsureDefaultFileNamingSettings(ctx, fileNamingDefaultsParams(defaults))
}

func (s *SettingsStore) GetFileNamingSettings(ctx context.Context) (FileNamingSettings, error) {
	if err := s.EnsureDefaultFileNamingSettings(ctx); err != nil {
		return FileNamingSettings{}, err
	}
	row, err := storagegen.New(s.pool).GetFileNamingSettings(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return FileNamingSettings{}, ErrNotFound
	}
	return fileNamingSettingsFromGetRow(row), err
}

func (s *SettingsStore) SaveFileNamingSettings(
	ctx context.Context,
	input FileNamingSettingsInput,
) (FileNamingSettings, error) {
	normalized, err := normalizeFileNamingSettings(input)
	if err != nil {
		return FileNamingSettings{}, err
	}
	if err := s.EnsureDefaultFileNamingSettings(ctx); err != nil {
		return FileNamingSettings{}, err
	}
	row, err := storagegen.New(s.pool).UpdateFileNamingSettings(ctx, fileNamingInputParams(normalized))
	return fileNamingSettingsFromUpdateRow(row), err
}

func normalizeFileNamingSettings(input FileNamingSettingsInput) (FileNamingSettingsInput, error) {
	normalized := FileNamingSettingsInput{
		MovieFileFormat:      normalizeTemplate(input.MovieFileFormat),
		MovieFolderFormat:    normalizeTemplate(input.MovieFolderFormat),
		SeriesEpisodeFormat:  normalizeTemplate(input.SeriesEpisodeFormat),
		DailyEpisodeFormat:   normalizeTemplate(input.DailyEpisodeFormat),
		AnimeEpisodeFormat:   normalizeTemplate(input.AnimeEpisodeFormat),
		SeriesFolderFormat:   normalizeTemplate(input.SeriesFolderFormat),
		SeasonFolderFormat:   normalizeTemplate(input.SeasonFolderFormat),
		SpecialsFolderFormat: normalizeTemplate(input.SpecialsFolderFormat),
	}
	if normalized.MovieFileFormat == "" ||
		normalized.MovieFolderFormat == "" ||
		normalized.SeriesEpisodeFormat == "" ||
		normalized.DailyEpisodeFormat == "" ||
		normalized.AnimeEpisodeFormat == "" ||
		normalized.SeriesFolderFormat == "" ||
		normalized.SeasonFolderFormat == "" ||
		normalized.SpecialsFolderFormat == "" {
		return FileNamingSettingsInput{}, ErrInvalidInput
	}
	return normalized, nil
}

func normalizeTemplate(value string) string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	return normalizeTemplateTokens(normalized)
}

func mediaMainFolderPath(root string, settings FileNamingSettings, input MediaItemInput) string {
	template := settings.SeriesFolderFormat
	if input.Type == "movie" {
		template = settings.MovieFolderFormat
	}
	rendered := renderMediaTemplate(template, input)
	return filepath.Join(root, sanitizePathSegment(rendered))
}

func renderMediaTemplate(template string, input MediaItemInput) string {
	title := input.Title
	year := ""
	if input.Year != nil {
		year = formatInt32(*input.Year)
	}
	values := map[string]string{
		"movie_title":  title,
		"quality_full": strings.TrimSpace(input.QualityFull),
		"release_year": year,
		"series_title": title,
		"year":         year,
	}
	rendered := fileNamingTokenPattern.ReplaceAllStringFunc(template, func(token string) string {
		key := normalizeTemplateTokenName(strings.Trim(token, "{}"))
		if value, ok := values[key]; ok {
			return value
		}
		return token
	})
	return strings.Join(strings.Fields(rendered), " ")
}

func sanitizePathSegment(value string) string {
	replacer := strings.NewReplacer(
		"/", " ",
		"\\", " ",
		":", " -",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	value = strings.TrimSpace(replacer.Replace(value))
	value = strings.Trim(value, ".")
	if value == "" {
		return "Untitled"
	}
	return value
}

func ensureMediaMainFolder(ctx context.Context, q mediaItemQuerier, input MediaItemInput) (*string, error) {
	if input.LibraryFolderID == nil {
		return nil, nil
	}
	folder, err := storagegen.New(q).GetLibraryFolder(ctx, *input.LibraryFolderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	settings, err := getFileNamingSettings(ctx, q)
	if err != nil {
		return nil, err
	}
	path := mediaMainFolderPath(folder.Path, settings, input)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}
	return &path, nil
}

func getFileNamingSettings(ctx context.Context, q mediaItemQuerier) (FileNamingSettings, error) {
	row, err := storagegen.New(q).GetFileNamingSettings(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		defaults := DefaultFileNamingSettings()
		return defaults, nil
	}
	return fileNamingSettingsFromGetRow(row), err
}

func formatInt32(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func fileNamingDefaultsParams(settings FileNamingSettings) storagegen.EnsureDefaultFileNamingSettingsParams {
	return storagegen.EnsureDefaultFileNamingSettingsParams{
		MovieFileFormat:      settings.MovieFileFormat,
		MovieFolderFormat:    settings.MovieFolderFormat,
		SeriesEpisodeFormat:  settings.SeriesEpisodeFormat,
		DailyEpisodeFormat:   settings.DailyEpisodeFormat,
		AnimeEpisodeFormat:   settings.AnimeEpisodeFormat,
		SeriesFolderFormat:   settings.SeriesFolderFormat,
		SeasonFolderFormat:   settings.SeasonFolderFormat,
		SpecialsFolderFormat: settings.SpecialsFolderFormat,
	}
}

func fileNamingInputParams(input FileNamingSettingsInput) storagegen.UpdateFileNamingSettingsParams {
	return storagegen.UpdateFileNamingSettingsParams{
		MovieFileFormat:      input.MovieFileFormat,
		MovieFolderFormat:    input.MovieFolderFormat,
		SeriesEpisodeFormat:  input.SeriesEpisodeFormat,
		DailyEpisodeFormat:   input.DailyEpisodeFormat,
		AnimeEpisodeFormat:   input.AnimeEpisodeFormat,
		SeriesFolderFormat:   input.SeriesFolderFormat,
		SeasonFolderFormat:   input.SeasonFolderFormat,
		SpecialsFolderFormat: input.SpecialsFolderFormat,
	}
}

func fileNamingSettingsFromGetRow(row storagegen.GetFileNamingSettingsRow) FileNamingSettings {
	return FileNamingSettings{
		MovieFileFormat:      row.MovieFileFormat,
		MovieFolderFormat:    row.MovieFolderFormat,
		SeriesEpisodeFormat:  row.SeriesEpisodeFormat,
		DailyEpisodeFormat:   row.DailyEpisodeFormat,
		AnimeEpisodeFormat:   row.AnimeEpisodeFormat,
		SeriesFolderFormat:   row.SeriesFolderFormat,
		SeasonFolderFormat:   row.SeasonFolderFormat,
		SpecialsFolderFormat: row.SpecialsFolderFormat,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
}

func fileNamingSettingsFromUpdateRow(row storagegen.UpdateFileNamingSettingsRow) FileNamingSettings {
	return FileNamingSettings{
		MovieFileFormat:      row.MovieFileFormat,
		MovieFolderFormat:    row.MovieFolderFormat,
		SeriesEpisodeFormat:  row.SeriesEpisodeFormat,
		DailyEpisodeFormat:   row.DailyEpisodeFormat,
		AnimeEpisodeFormat:   row.AnimeEpisodeFormat,
		SeriesFolderFormat:   row.SeriesFolderFormat,
		SeasonFolderFormat:   row.SeasonFolderFormat,
		SpecialsFolderFormat: row.SpecialsFolderFormat,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
}
