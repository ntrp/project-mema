package delivery

import (
	"strconv"
	"strings"
)

func OptionalDuration(value string) *float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || parsed <= 0 {
		return nil
	}
	return &parsed
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "unknown") || strings.EqualFold(value, "n/a") {
		return nil
	}
	return &value
}

func optionalIndex(value int32) *int32 {
	if value < 0 {
		return nil
	}
	return &value
}

func optionalInt(value int32) *int32 {
	if value <= 0 {
		return nil
	}
	return &value
}

func normalFrameRate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "0/0" {
		return ""
	}
	return value
}

func languageTag(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	if value := strings.TrimSpace(tags["language"]); value != "" {
		return value
	}
	if value := strings.TrimSpace(tags["LANGUAGE"]); value != "" {
		return value
	}
	return ""
}
