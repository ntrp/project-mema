package httpapi

import (
	"net/http"
	"testing"
)

func TestScenarioSCNSettings019AdminUpdatesFileNamingTemplates(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-019")

	var current FileNamingSettings
	client.doJSON(t, http.MethodGet, "/settings/file-naming", nil, http.StatusOK, &current)
	if current.MovieFileFormat == "" || current.SeriesEpisodeFormat == "" {
		t.Fatalf("file naming defaults = %#v", current)
	}

	var updated FileNamingSettings
	client.doJSON(t, http.MethodPut, "/settings/file-naming", FileNamingSettingsRequest{
		MovieFileFormat:      "  {movie_title} ({release_year}) [{quality_full}]  ",
		MovieFolderFormat:    "{movie_title} ({release_year})",
		SeriesEpisodeFormat:  "{series_title} - S{season:00}E{episode:00}",
		DailyEpisodeFormat:   "{series_title} - {air_date}",
		AnimeEpisodeFormat:   "{series_title} - {absolute_episode:000}",
		SeriesFolderFormat:   "{series_title}",
		SeasonFolderFormat:   "Season {season:00}",
		SpecialsFolderFormat: "Specials",
	}, http.StatusOK, &updated)
	if updated.MovieFileFormat != "{movie_title} ({release_year}) [{quality_full}]" {
		t.Fatalf("updated file naming settings = %#v", updated)
	}
}
