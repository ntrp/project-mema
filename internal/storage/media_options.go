package storage

func normalizeMediaItemOptions(input MediaItemInput) MediaItemInput {
	input.ContentKind = normalizeContentKind(input.ContentKind)
	input.MonitorMode = normalizeMonitorMode(input.Type, input.MonitorMode)
	input.SeriesType = normalizeSeriesType(input.Type, input.SeriesType)
	input.NumberingStrategy = normalizeNumberingStrategy(input.Type, input.ContentKind, input.NumberingStrategy)
	input.MinimumAvailability = normalizeMinimumAvailability(input.MinimumAvailability)
	input.Monitored = input.MonitorMode != "none"
	return input
}

func normalizeMediaRequestOptions(input MediaRequestInput) MediaRequestInput {
	input.MonitorMode = normalizeMonitorMode(input.Type, input.MonitorMode)
	input.SeriesType = normalizeSeriesType(input.Type, input.SeriesType)
	input.MinimumAvailability = normalizeMinimumAvailability(input.MinimumAvailability)
	return input
}

func normalizeMonitorMode(mediaType string, value string) string {
	if mediaType == "serie" {
		switch value {
		case "none", "all_episodes", "future_episodes", "missing_episodes", "existing_episodes", "no_specials":
			return value
		default:
			return "all_episodes"
		}
	}
	switch value {
	case "none", "collection":
		return value
	default:
		return "only_media"
	}
}

func normalizeSeriesType(mediaType string, value *string) *string {
	if mediaType != "serie" || value == nil {
		return nil
	}
	switch *value {
	case "standard", "daily", "absolute":
		return value
	default:
		fallback := "standard"
		return &fallback
	}
}

func normalizeContentKind(value string) string {
	if value == "anime" {
		return "anime"
	}
	return "standard"
}

func normalizeNumberingStrategy(mediaType string, contentKind string, value *string) *string {
	if mediaType != "serie" {
		return nil
	}
	if value != nil {
		switch *value {
		case "tmdb_season_episode", "tvdb_season_episode", "anidb_absolute", "manual":
			return value
		}
	}
	fallback := "tmdb_season_episode"
	if contentKind == "anime" {
		fallback = "anidb_absolute"
	}
	return &fallback
}

func normalizeMinimumAvailability(value string) string {
	switch value {
	case "announced", "in_cinema", "released":
		return value
	default:
		return "released"
	}
}
