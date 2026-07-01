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
	_, err := s.pool.Exec(ctx, `
		insert into app.file_naming_settings (
			id,
			movie_file_format,
			movie_folder_format,
			series_episode_format,
			daily_episode_format,
			anime_episode_format,
			series_folder_format,
			season_folder_format,
			specials_folder_format
		)
		values (1, $1, $2, $3, $4, $5, $6, $7, $8)
		on conflict do nothing
	`,
		defaults.MovieFileFormat,
		defaults.MovieFolderFormat,
		defaults.SeriesEpisodeFormat,
		defaults.DailyEpisodeFormat,
		defaults.AnimeEpisodeFormat,
		defaults.SeriesFolderFormat,
		defaults.SeasonFolderFormat,
		defaults.SpecialsFolderFormat,
	)
	return err
}

func (s *SettingsStore) GetFileNamingSettings(ctx context.Context) (FileNamingSettings, error) {
	if err := s.EnsureDefaultFileNamingSettings(ctx); err != nil {
		return FileNamingSettings{}, err
	}
	settings, err := scanFileNamingSettings(s.pool.QueryRow(ctx, `
		select movie_file_format,
			movie_folder_format,
			series_episode_format,
			daily_episode_format,
			anime_episode_format,
			series_folder_format,
			season_folder_format,
			specials_folder_format,
			created_at,
			updated_at
		from app.file_naming_settings
		where id = 1
	`))
	if errors.Is(err, pgx.ErrNoRows) {
		return FileNamingSettings{}, ErrNotFound
	}
	return settings, err
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
	return scanFileNamingSettings(s.pool.QueryRow(ctx, `
		update app.file_naming_settings
		set movie_file_format = $1,
			movie_folder_format = $2,
			series_episode_format = $3,
			daily_episode_format = $4,
			anime_episode_format = $5,
			series_folder_format = $6,
			season_folder_format = $7,
			specials_folder_format = $8,
			updated_at = now()
		where id = 1
		returning movie_file_format,
			movie_folder_format,
			series_episode_format,
			daily_episode_format,
			anime_episode_format,
			series_folder_format,
			season_folder_format,
			specials_folder_format,
			created_at,
			updated_at
	`,
		normalized.MovieFileFormat,
		normalized.MovieFolderFormat,
		normalized.SeriesEpisodeFormat,
		normalized.DailyEpisodeFormat,
		normalized.AnimeEpisodeFormat,
		normalized.SeriesFolderFormat,
		normalized.SeasonFolderFormat,
		normalized.SpecialsFolderFormat,
	))
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

func scanFileNamingSettings(row pgx.Row) (FileNamingSettings, error) {
	var settings FileNamingSettings
	err := row.Scan(
		&settings.MovieFileFormat,
		&settings.MovieFolderFormat,
		&settings.SeriesEpisodeFormat,
		&settings.DailyEpisodeFormat,
		&settings.AnimeEpisodeFormat,
		&settings.SeriesFolderFormat,
		&settings.SeasonFolderFormat,
		&settings.SpecialsFolderFormat,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	return settings, err
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
	var root string
	if err := q.QueryRow(ctx, `
		select path
		from app.library_folders
		where id = $1
	`, input.LibraryFolderID).Scan(&root); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	settings, err := getFileNamingSettings(ctx, q)
	if err != nil {
		return nil, err
	}
	path := mediaMainFolderPath(root, settings, input)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}
	return &path, nil
}

func getFileNamingSettings(ctx context.Context, q mediaItemQuerier) (FileNamingSettings, error) {
	settings, err := scanFileNamingSettings(q.QueryRow(ctx, `
		select movie_file_format,
			movie_folder_format,
			series_episode_format,
			daily_episode_format,
			anime_episode_format,
			series_folder_format,
			season_folder_format,
			specials_folder_format,
			created_at,
			updated_at
		from app.file_naming_settings
		where id = 1
	`))
	if errors.Is(err, pgx.ErrNoRows) {
		defaults := DefaultFileNamingSettings()
		return defaults, nil
	}
	return settings, err
}

func formatInt32(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}
