package dlna

import (
	"context"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"media-manager/internal/delivery"
)

func (m *Manager) TraceRendererProfile(ctx context.Context, request RendererRequest) RendererProfileTrace {
	profiles, deviceOverrides := m.rendererProfileCache(ctx)
	overrides, rememberedProfileID := m.profileMatchState(request.ClientIP)
	match := matchRendererProfiles(profiles, deviceOverrides, request, overrides, rememberedProfileID)
	return RendererProfileTrace{
		Match:          match,
		HeadersSummary: SafeHeadersSummary(request.Headers),
		Rules:          profileRuleTraces(profiles, request),
	}
}

func (m *Manager) TraceDeliveryDecision(ctx context.Context, input DeliveryTraceInput) DeliveryDecisionTrace {
	profile := m.deliveryTraceProfile(ctx, input)
	probe := probeWithPathContainer(input.Probe, input.MediaPath)
	if probeEmpty(probe) && strings.TrimSpace(input.MediaPath) != "" {
		probe = probeWithPathContainer(delivery.Probe(input.MediaPath), input.MediaPath)
	}
	capability := EvaluateRendererCapability(profile, probe)
	return DeliveryDecisionTrace{
		ProfileID:     profile.ID,
		ProfileName:   profile.Name,
		MediaFileName: safeMediaFileName(input.MediaPath, input.ObjectID),
		ObjectID:      strings.TrimSpace(input.ObjectID),
		ResourceID:    strings.TrimSpace(input.ResourceID),
		StreamMode:    strings.TrimSpace(input.StreamMode),
		Decision:      capability.Decision,
		ReasonCodes:   append([]string{}, capability.ReasonCodes...),
		Trace:         append([]RendererCapabilityTrace{}, capability.Trace...),
	}
}

func (m *Manager) deliveryTraceProfile(ctx context.Context, input DeliveryTraceInput) RendererProfile {
	profiles, _ := m.rendererProfileCache(ctx)
	if strings.TrimSpace(input.ProfileID) != "" {
		if profile, ok := findProfile(profiles, input.ProfileID); ok {
			return profile
		}
	}
	return m.ExplainRendererProfile(ctx, input.Request).Profile
}

func profileRuleTraces(profiles []RendererProfile, request RendererRequest) []RendererProfileRuleTrace {
	traces := []RendererProfileRuleTrace{}
	for _, profile := range profiles {
		if profile.ID == "generic" {
			continue
		}
		traces = append(traces, profileRuleTrace(profile, request)...)
	}
	sort.SliceStable(traces, func(i, j int) bool {
		if traces[i].ProfileID != traces[j].ProfileID {
			return traces[i].ProfileID < traces[j].ProfileID
		}
		return traces[i].Rule < traces[j].Rule
	})
	return traces
}

func profileRuleTrace(profile RendererProfile, request RendererRequest) []RendererProfileRuleTrace {
	if len(profile.rules) > 0 {
		return structuredRuleTraces(profile, request)
	}
	traces := make([]RendererProfileRuleTrace, 0, len(profile.MatchTokens))
	for _, token := range profile.MatchTokens {
		value := safeTraceValue(request.UserAgent + " " + request.FriendlyName + " " + headersText(request.Headers))
		matched := strings.Contains(strings.ToLower(value), strings.ToLower(token))
		traces = append(traces, RendererProfileRuleTrace{
			ProfileID: profile.ID, ProfileName: profile.Name, Field: "any",
			Value: value, Rule: "token:" + token, Score: boolScore(matched, 1), Result: passFail(matched),
		})
	}
	return traces
}

func structuredRuleTraces(profile RendererProfile, request RendererRequest) []RendererProfileRuleTrace {
	traces := make([]RendererProfileRuleTrace, 0, len(profile.rules))
	for _, rule := range profile.rules {
		value := safeTraceValue(matchFieldText(rule.Field, request))
		matched := ruleMatches(rule, request)
		traces = append(traces, RendererProfileRuleTrace{
			ProfileID: profile.ID, ProfileName: profile.Name, Field: rule.Field,
			Value: value, Rule: rule.Contains, Score: boolScore(matched, rule.Score), Result: passFail(matched),
		})
	}
	return traces
}

func SafeHeadersSummary(headers http.Header) []string {
	values := make([]string, 0, len(headers))
	for key, items := range headers {
		name := strings.TrimSpace(key)
		if unsafeHeaderName(name) {
			continue
		}
		joined := safeTraceValue(strings.Join(items, ", "))
		if joined == "" {
			continue
		}
		values = append(values, name+": "+joined)
	}
	sort.Strings(values)
	if len(values) > 12 {
		return values[:12]
	}
	return values
}

func RendererMatchReason(explanation RendererProfileExplanation) string {
	parts := []string{explanation.MatchSource}
	for _, value := range []string{explanation.WinningRule, explanation.FallbackPath} {
		if strings.TrimSpace(value) != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, ":")
}

func probeEmpty(probe delivery.ProbeResult) bool {
	return probe.Container.FormatName == nil && probe.Container.Format == nil && len(probe.Tracks) == 0
}

func safeMediaFileName(mediaPath string, fallback string) string {
	name := filepath.Base(strings.TrimSpace(mediaPath))
	if name == "." || name == "/" || name == "" {
		return strings.TrimSpace(fallback)
	}
	return name
}

func safeTraceValue(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\n", " "))
	if len(value) > 160 {
		return value[:160]
	}
	return value
}

func unsafeHeaderName(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "authorization", "cookie", "set-cookie", "x-api-key":
		return true
	default:
		return false
	}
}

func boolScore(ok bool, score int) int {
	if ok {
		return score
	}
	return 0
}
