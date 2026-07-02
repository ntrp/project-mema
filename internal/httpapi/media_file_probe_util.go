package httpapi

import "strings"

func optionalProbeString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" || value == "0/0" || strings.EqualFold(value, "unknown") {
		return nil
	}
	return &value
}

func optionalProbeInt(value int32) *int32 {
	if value <= 0 {
		return nil
	}
	return &value
}

func optionalProbeIndex(value int32) *int32 {
	if value < 0 {
		return nil
	}
	return &value
}

func languageTag(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	return firstString(tags["language"], tags["LANGUAGE"])
}

func normalFrameRate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "0/0" {
		return ""
	}
	return value
}

func firstString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
