package dlna

import (
	"context"
	"encoding/json"
	"strings"

	"media-manager/internal/storage"
)

type rendererProfileCacheState struct {
	profiles  []RendererProfile
	overrides []rendererDeviceOverride
	loaded    bool
}

type rendererMatchRulesJSON struct {
	Mode     string                   `json:"mode"`
	MinScore int                      `json:"minScore"`
	Tokens   []rendererMatchTokenJSON `json:"tokens"`
}

type rendererMatchTokenJSON struct {
	Field    string `json:"field"`
	Contains string `json:"contains"`
	Value    string `json:"value"`
	Score    int    `json:"score"`
}

type rendererDeliveryJSON struct {
	PreferHLS        bool `json:"preferHls"`
	AvoidHLS         bool `json:"avoidHls"`
	DirectPlay       bool `json:"directPlay"`
	Transcode        bool `json:"transcode"`
	StreamingHeaders bool `json:"streamingHeaders"`
}

type rendererCapabilitiesJSON struct {
	Containers    []string `json:"containers"`
	VideoCodecs   []string `json:"videoCodecs"`
	AudioCodecs   []string `json:"audioCodecs"`
	MaxResolution string   `json:"maxResolution"`
}

type rendererSubtitlesJSON struct {
	Formats []string `json:"formats"`
}

type rendererQuirksJSON struct {
	DisableEventing bool              `json:"disableEventing"`
	ExtraHeaders    map[string]string `json:"extraHeaders"`
}

func (m *Manager) RefreshRendererProfiles(ctx context.Context) error {
	profiles, overrides, err := m.loadRendererProfiles(ctx)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.profileCache = rendererProfileCacheState{
		profiles:  append([]RendererProfile{}, profiles...),
		overrides: append([]rendererDeviceOverride{}, overrides...),
		loaded:    true,
	}
	m.mu.Unlock()
	return nil
}

func (m *Manager) rendererProfileCache(ctx context.Context) ([]RendererProfile, []rendererDeviceOverride) {
	m.mu.Lock()
	if m.profileCache.loaded {
		profiles := append([]RendererProfile{}, m.profileCache.profiles...)
		overrides := append([]rendererDeviceOverride{}, m.profileCache.overrides...)
		m.mu.Unlock()
		return profiles, overrides
	}
	m.mu.Unlock()
	if err := m.RefreshRendererProfiles(ctx); err != nil {
		m.setError(err)
		return DefaultRendererProfiles(), nil
	}
	m.mu.Lock()
	profiles := append([]RendererProfile{}, m.profileCache.profiles...)
	overrides := append([]rendererDeviceOverride{}, m.profileCache.overrides...)
	m.mu.Unlock()
	return profiles, overrides
}

func (m *Manager) profileMatchState(clientIP string) (map[string]string, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	overrides := map[string]string{}
	for key, value := range m.profileOverrides {
		overrides[key] = value
	}
	remembered := ""
	if existing, ok := m.recentClients[strings.TrimSpace(clientIP)]; ok {
		remembered = existing.ProfileID
	}
	return overrides, remembered
}

func (m *Manager) loadRendererProfiles(ctx context.Context) ([]RendererProfile, []rendererDeviceOverride, error) {
	if m.store == nil {
		return DefaultRendererProfiles(), nil, nil
	}
	rows, err := m.store.ListDLNARendererProfiles(ctx)
	if err != nil {
		return nil, nil, err
	}
	overrides, err := m.store.ListDLNARendererDeviceOverrides(ctx)
	if err != nil {
		return nil, nil, err
	}
	profiles := make([]RendererProfile, 0, len(rows))
	for _, row := range rows {
		if !row.Enabled {
			continue
		}
		profiles = append(profiles, rendererProfileFromStorage(row))
	}
	if len(profiles) == 0 {
		profiles = DefaultRendererProfiles()
	}
	return profiles, rendererOverridesFromStorage(overrides), nil
}

func rendererProfileFromStorage(row storage.DLNARendererProfile) RendererProfile {
	matchRules := parseRendererMatchRules(row.MatchRules)
	capabilities := parseRendererCapabilities(row.CapabilityRules)
	delivery := parseRendererDelivery(row.DeliverySettings)
	subtitles := parseRendererSubtitles(row.SubtitleRules)
	quirks := parseRendererQuirks(row.Quirks)
	headers := map[string]string{}
	if delivery.StreamingHeaders {
		for key, value := range streamingHeaders() {
			headers[key] = value
		}
	}
	for key, value := range quirks.ExtraHeaders {
		headers[key] = value
	}
	return RendererProfile{
		ID:              runtimeRendererProfileID(row.ID),
		SourceID:        row.ID,
		Name:            row.Name,
		MatchMinScore:   matchRules.MinScore,
		Priority:        int(row.Priority),
		PreferHLS:       delivery.PreferHLS,
		AvoidHLS:        delivery.AvoidHLS,
		DisableEventing: quirks.DisableEventing,
		SubtitleFormats: append([]string{}, subtitles.Formats...),
		ResponseHeaders: headers,
		Capabilities:    capabilities,
		DeliveryRules:   RendererDeliveryRules{DirectPlay: delivery.DirectPlay, Transcode: delivery.Transcode},
		rules:           matchRules.Rules,
	}
}

func rendererOverridesFromStorage(rows []storage.DLNARendererDeviceOverride) []rendererDeviceOverride {
	overrides := make([]rendererDeviceOverride, 0, len(rows))
	for _, row := range rows {
		if !row.Allowed {
			continue
		}
		override := rendererDeviceOverride{ProfileID: row.ProfileID}
		if row.RendererUUID != nil {
			override.RendererUUID = *row.RendererUUID
		}
		if row.IPAddress != nil {
			override.IPAddress = *row.IPAddress
		}
		overrides = append(overrides, override)
	}
	return overrides
}

type parsedRendererMatchRules struct {
	MinScore int
	Rules    []rendererMatchRule
}

func parseRendererMatchRules(payload []byte) parsedRendererMatchRules {
	var raw rendererMatchRulesJSON
	if err := json.Unmarshal(payload, &raw); err != nil {
		return parsedRendererMatchRules{MinScore: 1}
	}
	rules := make([]rendererMatchRule, 0, len(raw.Tokens))
	for _, token := range raw.Tokens {
		contains := firstNonEmpty(token.Contains, token.Value)
		score := token.Score
		if score == 0 {
			score = 1
		}
		rules = append(rules, rendererMatchRule{Field: token.Field, Contains: contains, Score: score})
	}
	if raw.MinScore == 0 {
		raw.MinScore = 1
	}
	return parsedRendererMatchRules{MinScore: raw.MinScore, Rules: rules}
}

func parseRendererDelivery(payload []byte) rendererDeliveryJSON {
	var raw rendererDeliveryJSON
	_ = json.Unmarshal(payload, &raw)
	return raw
}

func parseRendererCapabilities(payload []byte) RendererCapabilities {
	var raw rendererCapabilitiesJSON
	_ = json.Unmarshal(payload, &raw)
	return RendererCapabilities{
		Containers:    normalizedLowerList(raw.Containers),
		VideoCodecs:   normalizedLowerList(raw.VideoCodecs),
		AudioCodecs:   normalizedLowerList(raw.AudioCodecs),
		MaxResolution: strings.ToLower(strings.TrimSpace(raw.MaxResolution)),
	}
}

func parseRendererSubtitles(payload []byte) rendererSubtitlesJSON {
	var raw rendererSubtitlesJSON
	_ = json.Unmarshal(payload, &raw)
	return raw
}

func parseRendererQuirks(payload []byte) rendererQuirksJSON {
	var raw rendererQuirksJSON
	_ = json.Unmarshal(payload, &raw)
	return raw
}

func runtimeRendererProfileID(id string) string {
	switch id {
	case "samsung-tv":
		return "samsung"
	case "lg-webos":
		return "lg"
	case "sony-tv":
		return "sony"
	default:
		return id
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func normalizedLowerList(values []string) []string {
	results := make([]string, 0, len(values))
	for _, value := range values {
		cleaned := strings.ToLower(strings.TrimSpace(value))
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
