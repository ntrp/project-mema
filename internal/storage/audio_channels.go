package storage

import "strings"

func NormalizeAudioChannelDefinitions(values []string) []string {
	seen := map[string]struct{}{}
	channels := []string{}
	for _, value := range values {
		channel := NormalizeAudioChannelDefinition(value)
		if channel == "" {
			continue
		}
		if _, ok := seen[channel]; ok {
			continue
		}
		seen[channel] = struct{}{}
		channels = append(channels, channel)
	}
	return channels
}

func NormalizeAudioChannelDefinition(value string) string {
	switch normalizedChannelToken(value) {
	case "mono", "10", "10mono", "1", "1ch", "1channel", "1channels":
		return "1.0"
	case "stereo", "20", "20stereo", "2", "2ch", "2channel", "2channels":
		return "2.0"
	case "30", "30ch", "30sound", "3", "3ch", "3channel", "3channels":
		return "3.0"
	case "40", "40ch", "40sound", "4", "4ch", "4channel", "4channels":
		return "4.0"
	case "50", "50ch", "5", "5ch", "5channel", "5channels":
		return "5.0"
	case "51", "51ch", "51surround":
		return "5.1"
	case "61", "61ch", "61surround":
		return "6.1"
	case "71", "71ch", "71surround":
		return "7.1"
	default:
		return ""
	}
}

func normalizedChannelToken(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
