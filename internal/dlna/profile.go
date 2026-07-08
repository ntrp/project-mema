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
	Name            string
	MatchTokens     []string
	PreferHLS       bool
	AvoidHLS        bool
	DisableEventing bool
	SubtitleFormats []string
	ResponseHeaders map[string]string
}

type RendererRequest struct {
	UserAgent    string
	FriendlyName string
	Headers      http.Header
	ClientIP     string
}

func DefaultRendererProfiles() []RendererProfile {
	return []RendererProfile{
		{ID: "generic", Name: "Generic DLNA", SubtitleFormats: []string{"srt", "vtt"}},
		{ID: "vlc", Name: "VLC", MatchTokens: []string{"vlc"}, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "kodi", Name: "Kodi", MatchTokens: []string{"kodi", "xbmc"}, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "samsung", Name: "Samsung TV", MatchTokens: []string{"samsung", "tizen"}, SubtitleFormats: []string{"srt"}, ResponseHeaders: streamingHeaders()},
		{ID: "lg", Name: "LG TV", MatchTokens: []string{"lg", "webos"}, AvoidHLS: true, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "sony", Name: "Sony TV", MatchTokens: []string{"sony", "bravia"}, SubtitleFormats: []string{"srt"}, ResponseHeaders: streamingHeaders()},
		{ID: "bubbleupnp", Name: "BubbleUPnP", MatchTokens: []string{"bubbleupnp"}, SubtitleFormats: []string{"srt", "vtt"}, ResponseHeaders: streamingHeaders()},
		{ID: "chromecast", Name: "Chromecast", MatchTokens: []string{"chromecast", "google cast"}, PreferHLS: true, DisableEventing: true, SubtitleFormats: []string{"vtt"}, ResponseHeaders: streamingHeaders()},
	}
}

func MatchRendererProfile(request RendererRequest, overrides map[string]string) RendererProfile {
	profiles := DefaultRendererProfiles()
	if id := overrides[strings.TrimSpace(request.ClientIP)]; id != "" {
		if profile, ok := findProfile(profiles, id); ok {
			return profile
		}
	}
	haystack := strings.ToLower(request.UserAgent + " " + request.FriendlyName + " " + headersText(request.Headers))
	for _, profile := range profiles[1:] {
		for _, token := range profile.MatchTokens {
			if strings.Contains(haystack, strings.ToLower(token)) {
				return profile
			}
		}
	}
	return profiles[0]
}

func SourceProtocolInfosForProfile(profile RendererProfile) []string {
	values := SourceProtocolInfos()
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

func (m *Manager) SetRendererProfileOverride(clientIP string, profileID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.profileOverrides == nil {
		m.profileOverrides = map[string]string{}
	}
	m.profileOverrides[strings.TrimSpace(clientIP)] = strings.TrimSpace(profileID)
}

func (m *Manager) RendererProfile(request RendererRequest) RendererProfile {
	m.mu.Lock()
	overrides := map[string]string{}
	for key, value := range m.profileOverrides {
		overrides[key] = value
	}
	rememberedProfileID := ""
	if existing, ok := m.recentClients[strings.TrimSpace(request.ClientIP)]; ok {
		rememberedProfileID = existing.ProfileID
	}
	m.mu.Unlock()
	profile := MatchRendererProfile(request, overrides)
	if profile.ID != "generic" || rememberedProfileID == "" || rememberedProfileID == "generic" {
		return profile
	}
	if remembered, ok := findProfile(DefaultRendererProfiles(), rememberedProfileID); ok {
		return remembered
	}
	return profile
}

func (m *Manager) RendererProfileFromRequest(r *http.Request) RendererProfile {
	return m.RendererProfile(RendererRequestFromHTTP(r))
}

func (m *Manager) rendererProfileFromContext(ctx context.Context) RendererProfile {
	if r, ok := soap.RequestFromContext(ctx); ok {
		return m.RendererProfileFromRequest(r)
	}
	return MatchRendererProfile(RendererRequest{}, nil)
}

func RendererRequestFromHTTP(r *http.Request) RendererRequest {
	return RendererRequest{
		UserAgent: r.UserAgent(),
		Headers:   r.Header,
		ClientIP:  clientIP(r),
	}
}

func applyRendererHeaders(w http.ResponseWriter, profile RendererProfile) {
	for key, value := range profile.ResponseHeaders {
		w.Header().Set(key, value)
	}
}

func findProfile(profiles []RendererProfile, id string) (RendererProfile, bool) {
	for _, profile := range profiles {
		if profile.ID == id {
			return profile, true
		}
	}
	return RendererProfile{}, false
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

func streamingHeaders() map[string]string {
	return map[string]string{
		"TransferMode.DLNA.ORG":    "Streaming",
		"ContentFeatures.DLNA.ORG": "DLNA.ORG_OP=01;DLNA.ORG_CI=0",
	}
}
