package metadata

import (
	"strconv"
	"strings"
)

func yearFromDate(value string) *int32 {
	if len(value) < 4 {
		return nil
	}
	return yearFromString(value[:4])
}

func yearFromString(value string) *int32 {
	if len(value) < 4 {
		return nil
	}
	year, err := strconv.ParseInt(value[:4], 10, 32)
	if err != nil {
		return nil
	}
	result := int32(year)
	return &result
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalFloat64(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstString(values []string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
