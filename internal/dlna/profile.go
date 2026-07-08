package dlna

import (
	"context"
	"net"
	"net/http"
	"strings"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna/soap"
)

type RendererProfile struct {
	ID              string
	SourceID        string
	Name            string
	MatchTokens     []string
	MatchMinScore   int
	Priority        int
	PreferHLS       bool
	AvoidHLS        bool
	DisableEventing bool
	SubtitleFormats []string
	ResponseHeaders map[string]string
	Capabilities    RendererCapabilities
	DeliveryRules   RendererDeliveryRules
	rules           []rendererMatchRule
}

type RendererCapabilities struct {
	Containers    []string
	VideoCodecs   []string
	AudioCodecs   []string
	MaxResolution string
}

type RendererDeliveryRules struct {
	DirectPlay         bool
	Transcode          bool
	SeekMode           string
	RemuxContainer     string
	TranscodeContainer string
}

type RendererProfileExplanation struct {
	SelectedProfileID   string
	SourceProfileID     string
	MatchSource         string
	WinningRule         string
	FallbackPath        string
	Score               int
	CandidateProfileIDs []string
}

type RendererProfileMatch struct {
	Profile     RendererProfile
	Explanation RendererProfileExplanation
}

type RendererProfileRuleTrace struct {
	ProfileID   string
	ProfileName string
	Field       string
	Value       string
	Rule        string
	Score       int
	Result      string
}

type RendererProfileTrace struct {
	Match          RendererProfileMatch
	HeadersSummary []string
	Rules          []RendererProfileRuleTrace
}

type DeliveryDecisionTrace struct {
	ProfileID     string
	ProfileName   string
	MediaFileName string
	ObjectID      string
	ResourceID    string
	StreamMode    string
	Decision      delivery.Decision
	ReasonCodes   []string
	Trace         []RendererCapabilityTrace
}

type DeliveryTraceInput struct {
	Request    RendererRequest
	ProfileID  string
	MediaPath  string
	ObjectID   string
	ResourceID string
	StreamMode string
	Probe      delivery.ProbeResult
}

type RendererRequest struct {
	UserAgent    string
	FriendlyName string
	Headers      http.Header
	ClientIP     string
	RendererUUID string
}

func DefaultRendererProfiles() []RendererProfile {
	return []RendererProfile{
		{ID: "generic", Name: "Generic DLNA", SubtitleFormats: []string{"srt", "vtt"}},
		{ID: "vlc", Name: "VLC", MatchTokens: []string{"vlc"}, Priority: 100, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "kodi", Name: "Kodi", MatchTokens: []string{"kodi", "xbmc"}, Priority: 100, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "samsung", SourceID: "samsung-tv", Name: "Samsung TV", MatchTokens: []string{"samsung", "tizen"}, Priority: 100, SubtitleFormats: []string{"srt"}, ResponseHeaders: streamingHeaders()},
		{ID: "lg", SourceID: "lg-webos", Name: "LG TV", MatchTokens: []string{"lg", "webos"}, Priority: 100, AvoidHLS: true, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "sony", SourceID: "sony-tv", Name: "Sony TV", MatchTokens: []string{"sony", "bravia"}, Priority: 100, SubtitleFormats: []string{"srt"}, ResponseHeaders: streamingHeaders()},
		{ID: "bubbleupnp", Name: "BubbleUPnP", MatchTokens: []string{"bubbleupnp"}, Priority: 100, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "chromecast", Name: "Chromecast", MatchTokens: []string{"chromecast", "google cast"}, Priority: 100, PreferHLS: true, DisableEventing: true, SubtitleFormats: []string{"vtt"}, ResponseHeaders: streamingHeaders()},
	}
}

func MatchRendererProfile(request RendererRequest, overrides map[string]string) RendererProfile {
	return matchRendererProfiles(DefaultRendererProfiles(), nil, request, overrides, "").Profile
}

func (m *Manager) SetRendererProfileOverride(clientIP string, profileID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.profileOverrides == nil {
		m.profileOverrides = map[string]string{}
	}
	m.profileOverrides[strings.TrimSpace(clientIP)] = strings.TrimSpace(profileID)
}

func (m *Manager) RendererProfile(request RendererRequest) RendererProfile {
	return m.ExplainRendererProfile(context.Background(), request).Profile
}

func (m *Manager) ExplainRendererProfile(ctx context.Context, request RendererRequest) RendererProfileMatch {
	profiles, deviceOverrides := m.rendererProfileCache(ctx)
	overrides, rememberedProfileID := m.profileMatchState(request.ClientIP)
	return matchRendererProfiles(profiles, deviceOverrides, request, overrides, rememberedProfileID)
}

func (m *Manager) RendererProfileFromRequest(r *http.Request) RendererProfile {
	return m.ExplainRendererProfile(r.Context(), RendererRequestFromHTTP(r)).Profile
}

func (m *Manager) rendererProfileFromContext(ctx context.Context) RendererProfile {
	if r, ok := soap.RequestFromContext(ctx); ok {
		return m.RendererProfileFromRequest(r)
	}
	return m.ExplainRendererProfile(ctx, RendererRequest{}).Profile
}

func RendererRequestFromHTTP(r *http.Request) RendererRequest {
	return RendererRequest{
		UserAgent:    r.UserAgent(),
		FriendlyName: friendlyNameFromHeaders(r.Header),
		Headers:      r.Header,
		ClientIP:     clientIP(r),
		RendererUUID: rendererUUIDFromHeaders(r.Header),
	}
}

func SourceProtocolInfosForProfile(profile RendererProfile) []string {
	values := SourceProtocolInfosForCapabilities(profile)
	if profile.AvoidHLS {
		return protocolInfosWithoutHLS(values)
	}
	if !profile.PreferHLS {
		return values
	}
	for index, value := range values {
		if strings.Contains(value, "application/vnd.apple.mpegurl") {
			return append([]string{value}, append(values[:index], values[index+1:]...)...)
		}
	}
	return values
}

func protocolInfosWithoutHLS(values []string) []string {
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		if strings.Contains(value, "application/vnd.apple.mpegurl") {
			continue
		}
		filtered = append(filtered, value)
	}
	return filtered
}

func SourceProtocolInfoForProfile(profile RendererProfile) string {
	return strings.Join(SourceProtocolInfosForProfile(profile), ",")
}

func DeliveryClientProfile(profile RendererProfile) delivery.ClientProfile {
	if profile.PreferHLS {
		return delivery.ClientWebKit
	}
	return delivery.ClientBrowser
}

func applyRendererHeaders(w http.ResponseWriter, profile RendererProfile) {
	for key, value := range profile.ResponseHeaders {
		w.Header().Set(key, value)
	}
}

func findProfile(profiles []RendererProfile, id string) (RendererProfile, bool) {
	id = strings.TrimSpace(id)
	for _, profile := range profiles {
		if profile.ID == id || profile.SourceID == id {
			return profile, true
		}
	}
	return RendererProfile{}, false
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		if index := strings.Index(forwarded, ","); index >= 0 {
			forwarded = forwarded[:index]
		}
		return strings.TrimSpace(forwarded)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func friendlyNameFromHeaders(headers http.Header) string {
	for _, name := range []string{"X-Mema-DLNA-Friendly-Name", "X-Renderer-Name", "FriendlyName"} {
		if value := strings.TrimSpace(headers.Get(name)); value != "" {
			return value
		}
	}
	return ""
}

func streamingHeaders() map[string]string {
	return map[string]string{
		"TransferMode.DLNA.ORG":    "Streaming",
		"ContentFeatures.DLNA.ORG": "DLNA.ORG_OP=01;DLNA.ORG_CI=0",
	}
}
