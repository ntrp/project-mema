package storage

func normalizeMediaItemOptions(input MediaItemInput) MediaItemInput {
	input.MonitorMode = normalizeMonitorMode(input.Type, input.MonitorMode)
	input.SeriesType = normalizeSeriesType(input.Type, input.SeriesType)
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
	if mediaType == "series" {
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
	if mediaType != "series" || value == nil {
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

func normalizeMinimumAvailability(value string) string {
	switch value {
	case "announced", "in_cinema", "released":
		return value
	default:
		return "released"
	}
}
