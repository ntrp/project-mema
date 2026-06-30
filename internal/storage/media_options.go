package storage

func normalizeMediaItemOptions(input MediaItemInput) MediaItemInput {
	input.MonitorMode = normalizeMonitorMode(input.MonitorMode)
	input.MinimumAvailability = normalizeMinimumAvailability(input.MinimumAvailability)
	input.Monitored = input.MonitorMode != "none"
	return input
}

func normalizeMediaRequestOptions(input MediaRequestInput) MediaRequestInput {
	input.MonitorMode = normalizeMonitorMode(input.MonitorMode)
	input.MinimumAvailability = normalizeMinimumAvailability(input.MinimumAvailability)
	return input
}

func normalizeMonitorMode(value string) string {
	switch value {
	case "none", "collection":
		return value
	default:
		return "only_media"
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
