package downloadclients

import (
	"math"
	"strconv"
	"strings"
)

func statusLookupResultWithProgress(status string, progressPercent *int, message string, pairs ...interface{}) StatusResult {
	result := statusLookupResult(status, message, pairs...)
	result.ProgressPercent = progressPercent
	return result
}

func progressFromFraction(value float64) *int {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	percent := int(math.Round(value * 100))
	return &percent
}

func progressFromPercentString(value string) *int {
	trimmed := strings.TrimSpace(strings.TrimSuffix(value, "%"))
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil
	}
	if parsed < 0 {
		parsed = 0
	}
	if parsed > 100 {
		parsed = 100
	}
	percent := int(math.Round(parsed))
	return &percent
}

func completedProgress() *int {
	percent := 100
	return &percent
}
