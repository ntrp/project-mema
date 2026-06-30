package decisions

import "strings"

func containsAnyNormalized(value string, candidates ...string) bool {
	normalizedValue := normalizedToken(value)
	for _, candidate := range candidates {
		if token := normalizedToken(candidate); token != "" && strings.Contains(normalizedValue, token) {
			return true
		}
	}
	return false
}

func containsString(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func tokenSet(value string) map[string]struct{} {
	fields := strings.Fields(strings.ToLower(releaseSeparator.ReplaceAllString(value, " ")))
	tokens := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		token := normalizedToken(field)
		if token != "" {
			tokens[token] = struct{}{}
		}
	}
	return tokens
}
