package dlna

import (
	"net/http"
	"strings"
)

type rendererMatchRule struct {
	Field    string
	Contains string
	Score    int
}

type rendererDeviceOverride struct {
	RendererUUID string
	IPAddress    string
	ProfileID    string
}

type rendererProfileCandidate struct {
	profile     RendererProfile
	score       int
	winningRule string
}

func matchRendererProfiles(
	profiles []RendererProfile,
	deviceOverrides []rendererDeviceOverride,
	request RendererRequest,
	overrides map[string]string,
	rememberedProfileID string,
) RendererProfileMatch {
	if profile, ok := manualProfile(profiles, deviceOverrides, request); ok {
		return explainProfile(profile, "manual_uuid", "uuid override", "", 0, nil)
	}
	if profile, ok := ipOverrideProfile(profiles, deviceOverrides, request, overrides); ok {
		return explainProfile(profile, "manual_ip", "ip override", "", 0, nil)
	}
	if candidate, candidates, ok := bestAutomaticProfile(profiles, request); ok {
		return explainProfile(candidate.profile, "match", candidate.winningRule, "", candidate.score, candidates)
	}
	if rememberedProfileID != "" && rememberedProfileID != "generic" {
		if profile, ok := findProfile(profiles, rememberedProfileID); ok {
			return explainProfile(profile, "sticky_ip", "recent client", "generic", 0, nil)
		}
	}
	if profile, ok := findProfile(profiles, "generic"); ok {
		return explainProfile(profile, "default", "", "generic", 0, nil)
	}
	return explainProfile(RendererProfile{ID: "generic", Name: "Generic DLNA"}, "default", "", "empty_cache", 0, nil)
}

func manualProfile(
	profiles []RendererProfile,
	overrides []rendererDeviceOverride,
	request RendererRequest,
) (RendererProfile, bool) {
	uuid := normalizeRendererUUID(request.RendererUUID)
	if uuid == "" {
		return RendererProfile{}, false
	}
	for _, override := range overrides {
		if normalizeRendererUUID(override.RendererUUID) != uuid {
			continue
		}
		return findProfile(profiles, override.ProfileID)
	}
	return RendererProfile{}, false
}

func ipOverrideProfile(
	profiles []RendererProfile,
	deviceOverrides []rendererDeviceOverride,
	request RendererRequest,
	overrides map[string]string,
) (RendererProfile, bool) {
	ip := strings.TrimSpace(request.ClientIP)
	if ip == "" {
		return RendererProfile{}, false
	}
	for _, override := range deviceOverrides {
		if strings.TrimSpace(override.IPAddress) != ip {
			continue
		}
		return findProfile(profiles, override.ProfileID)
	}
	if id := overrides[ip]; id != "" {
		return findProfile(profiles, id)
	}
	return RendererProfile{}, false
}

func bestAutomaticProfile(
	profiles []RendererProfile,
	request RendererRequest,
) (rendererProfileCandidate, []string, bool) {
	best := rendererProfileCandidate{}
	candidateIDs := []string{}
	for _, profile := range profiles {
		if profile.ID == "generic" {
			continue
		}
		score, winningRule := profileMatchScore(profile, request)
		if score <= 0 || score < profileMinScore(profile) {
			continue
		}
		candidate := rendererProfileCandidate{profile: profile, score: score, winningRule: winningRule}
		candidateIDs = append(candidateIDs, profile.ID)
		if betterProfileCandidate(candidate, best) {
			best = candidate
		}
	}
	if best.profile.ID == "" {
		return rendererProfileCandidate{}, candidateIDs, false
	}
	return best, candidateIDs, true
}

func profileMinScore(profile RendererProfile) int {
	if profile.MatchMinScore > 0 {
		return profile.MatchMinScore
	}
	return 1
}

func profileMatchScore(profile RendererProfile, request RendererRequest) (int, string) {
	if len(profile.rules) > 0 {
		return structuredMatchScore(profile.rules, request)
	}
	haystack := strings.ToLower(request.UserAgent + " " + request.FriendlyName + " " + headersText(request.Headers))
	for _, token := range profile.MatchTokens {
		if strings.Contains(haystack, strings.ToLower(token)) {
			return 1, "token:" + token
		}
	}
	return 0, ""
}

func structuredMatchScore(rules []rendererMatchRule, request RendererRequest) (int, string) {
	score := 0
	winningScore := 0
	winningRule := ""
	for _, rule := range rules {
		if !ruleMatches(rule, request) {
			continue
		}
		score += rule.Score
		if rule.Score > winningScore {
			winningScore = rule.Score
			winningRule = rule.Field + ":" + rule.Contains
		}
	}
	return score, winningRule
}

func ruleMatches(rule rendererMatchRule, request RendererRequest) bool {
	needle := strings.ToLower(strings.TrimSpace(rule.Contains))
	if needle == "" {
		return false
	}
	return strings.Contains(strings.ToLower(matchFieldText(rule.Field, request)), needle)
}

func matchFieldText(field string, request RendererRequest) string {
	switch strings.ToLower(strings.TrimSpace(field)) {
	case "useragent", "user_agent":
		return request.UserAgent
	case "friendlyname", "friendly_name":
		return request.FriendlyName
	case "headers", "header":
		return headersText(request.Headers)
	case "uuid", "rendereruuid", "renderer_uuid":
		return request.RendererUUID
	default:
		return request.UserAgent + " " + request.FriendlyName + " " + headersText(request.Headers)
	}
}

func betterProfileCandidate(candidate rendererProfileCandidate, best rendererProfileCandidate) bool {
	if best.profile.ID == "" {
		return true
	}
	if candidate.profile.Priority != best.profile.Priority {
		return candidate.profile.Priority > best.profile.Priority
	}
	if candidate.score != best.score {
		return candidate.score > best.score
	}
	return candidate.profile.ID < best.profile.ID
}

func explainProfile(
	profile RendererProfile,
	source string,
	winningRule string,
	fallback string,
	score int,
	candidates []string,
) RendererProfileMatch {
	return RendererProfileMatch{
		Profile: profile,
		Explanation: RendererProfileExplanation{
			SelectedProfileID:   profile.ID,
			SourceProfileID:     profile.SourceID,
			MatchSource:         source,
			WinningRule:         winningRule,
			FallbackPath:        fallback,
			Score:               score,
			CandidateProfileIDs: append([]string{}, candidates...),
		},
	}
}

func headersText(headers http.Header) string {
	var builder strings.Builder
	for key, values := range headers {
		builder.WriteString(key)
		builder.WriteByte(' ')
		builder.WriteString(strings.Join(values, " "))
		builder.WriteByte(' ')
	}
	return builder.String()
}

func rendererUUIDFromHeaders(headers http.Header) string {
	for _, name := range []string{"X-Mema-DLNA-Renderer-UUID", "X-Renderer-UUID", "USN"} {
		if value := strings.TrimSpace(headers.Get(name)); value != "" {
			return normalizeRendererUUID(value)
		}
	}
	return ""
}

func normalizeRendererUUID(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if index := strings.Index(value, "::"); index >= 0 {
		value = value[:index]
	}
	return strings.TrimSpace(value)
}
