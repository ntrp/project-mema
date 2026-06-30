package httpapi

import (
	"strconv"
	"strings"
	"time"

	"media-manager/internal/storage"
)

func applySeriesMonitoring(input storage.MediaItemInput) storage.MediaItemInput {
	if input.Type != "series" {
		return input
	}
	input.MonitorMode = seriesMonitorMode(input.MonitorMode)
	today := time.Now().Truncate(24 * time.Hour)
	for seasonIndex := range input.Seasons {
		seasonNumber := seasonNumberFromName(input.Seasons[seasonIndex].Name)
		special := seasonNumber == 0 || strings.Contains(strings.ToLower(input.Seasons[seasonIndex].Name), "special")
		seasonMonitored := seasonMatchesMonitor(input.MonitorMode, special, input.Seasons[seasonIndex].AirDate, today)
		for episodeIndex := range input.Seasons[seasonIndex].Episodes {
			episode := &input.Seasons[seasonIndex].Episodes[episodeIndex]
			episode.Monitored = episodeMatchesMonitor(input.MonitorMode, special, episode.AirDate, today)
			seasonMonitored = seasonMonitored || episode.Monitored
		}
		input.Seasons[seasonIndex].Monitored = seasonMonitored
	}
	return input
}

func seriesMonitorMode(value string) string {
	switch value {
	case "none", "all_episodes", "future_episodes", "missing_episodes", "existing_episodes", "no_specials":
		return value
	default:
		return "all_episodes"
	}
}

func seasonMatchesMonitor(mode string, special bool, airDate *string, today time.Time) bool {
	switch mode {
	case "all_episodes", "missing_episodes":
		return true
	case "no_specials":
		return !special
	case "future_episodes":
		return dateAfter(airDate, today)
	default:
		return false
	}
}

func episodeMatchesMonitor(mode string, special bool, airDate *string, today time.Time) bool {
	switch mode {
	case "all_episodes", "missing_episodes":
		return true
	case "no_specials":
		return !special
	case "future_episodes":
		return dateAfter(airDate, today)
	default:
		return false
	}
}

func dateAfter(value *string, today time.Time) bool {
	if value == nil {
		return false
	}
	airDate, err := time.Parse("2006-01-02", *value)
	return err == nil && airDate.After(today)
}

func seasonNumberFromName(value string) int {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return -1
	}
	number, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1
	}
	return number
}
