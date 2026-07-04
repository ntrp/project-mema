package cardigann

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseCardigannDate(value string, layout string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	if layout != "" {
		if parsed, ok := parseNamedCardigannDate(value, layout); ok {
			return parsed, true
		}
	}
	if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(unix, 0).UTC(), true
	}
	for _, candidate := range []string{time.RFC3339, time.RFC1123Z, time.RFC1123, "2006-01-02", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(candidate, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func parseNamedCardigannDate(value string, layout string) (time.Time, bool) {
	now := time.Now()
	layout = strings.TrimSpace(layout)
	switch layout {
	case "unix":
		seconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, false
		}
		return time.Unix(seconds, 0).UTC(), true
	case "htt MMM. d":
		cleaned := normalizeCardigannMonth(value)
		parsed, err := time.Parse("3pm Jan 2 2006", cleaned+" "+strconv.Itoa(now.Year()))
		return parsed, err == nil
	case "MMM. d yy":
		cleaned := normalizeCardigannMonth(value)
		parsed, err := time.Parse("Jan 2 06", cleaned)
		return parsed, err == nil
	default:
		goLayout := strings.NewReplacer("yyyy", "2006", "yy", "06", "MMM", "Jan", "MM", "01", "dd", "02").Replace(layout)
		parsed, err := time.Parse(goLayout, value)
		return parsed, err == nil
	}
}

func normalizeCardigannMonth(value string) string {
	replacer := strings.NewReplacer(".", "", "st", "", "nd", "", "rd", "", "th", "")
	return strings.Join(strings.Fields(replacer.Replace(value)), " ")
}

func parseFuzzyTime(value string, now time.Time) (time.Time, bool) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return time.Time{}, false
	}
	if parsed, err := time.Parse("3:04pm", value); err == nil {
		return time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location()), true
	}
	re := regexp.MustCompile(`(\d+)\s*(minute|hour|day|week|month|year)s?\s+ago`)
	match := re.FindStringSubmatch(value)
	if len(match) != 3 {
		return time.Time{}, false
	}
	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return time.Time{}, false
	}
	switch match[2] {
	case "minute":
		return now.Add(-time.Duration(amount) * time.Minute), true
	case "hour":
		return now.Add(-time.Duration(amount) * time.Hour), true
	case "day":
		return now.AddDate(0, 0, -amount), true
	case "week":
		return now.AddDate(0, 0, -amount*7), true
	case "month":
		return now.AddDate(0, -amount, 0), true
	case "year":
		return now.AddDate(-amount, 0, 0), true
	default:
		return time.Time{}, false
	}
}

func parseSizeBytes(value string) int64 {
	cleaned := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(value, ",", "")))
	if cleaned == "" {
		return 0
	}
	if parsed, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
		return parsed
	}
	re := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*([kmgtp]?i?b|[kmgtp])`)
	match := re.FindStringSubmatch(cleaned)
	if len(match) != 3 {
		return 0
	}
	number, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0
	}
	multiplier := float64(1)
	switch strings.TrimSuffix(match[2], "b") {
	case "k", "ki":
		multiplier = 1024
	case "m", "mi":
		multiplier = 1024 * 1024
	case "g", "gi":
		multiplier = 1024 * 1024 * 1024
	case "t", "ti":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "p", "pi":
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024
	}
	return int64(number * multiplier)
}
