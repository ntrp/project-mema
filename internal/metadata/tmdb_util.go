package metadata

import (
	"strconv"
	"strings"
)

func tmdbCountries(items []tmdbCountry) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name := strings.TrimSpace(item.Name); name != "" {
			names = append(names, name)
		}
	}
	return names
}

func appendUniqueLimit(values []string, value string, limit int) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	if len(values) >= limit {
		return values
	}
	return append(values, value)
}

func formatMoney(value int64) string {
	digits := strconv.FormatInt(value, 10)
	parts := []string{}
	for len(digits) > 3 {
		parts = append([]string{digits[len(digits)-3:]}, parts...)
		digits = digits[:len(digits)-3]
	}
	parts = append([]string{digits}, parts...)
	return "$" + strings.Join(parts, ",") + ".00"
}

func releaseDateOnly(value string) string {
	if len(value) < 10 {
		return ""
	}
	return value[:10]
}

func earliestDate(current string, candidate string) string {
	if current == "" || candidate < current {
		return candidate
	}
	return current
}

func tmdbCollectionID(collection *tmdbCollection) *string {
	if collection == nil || collection.ID == 0 {
		return nil
	}
	return optionalString(strconv.FormatInt(collection.ID, 10))
}

func tmdbCollectionName(collection *tmdbCollection) *string {
	if collection == nil {
		return nil
	}
	return optionalString(collection.Name)
}

func tmdbNames(items []tmdbName) []string {
	names := []string{}
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
