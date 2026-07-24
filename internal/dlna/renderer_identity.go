package dlna

import (
	"net"
	"strings"
	"time"
)

func (m *Manager) verifiedRendererRequest(request RendererRequest, overrides []rendererDeviceOverride) RendererRequest {
	claimedUUID := normalizeRendererUUID(request.RendererUUID)
	request.RendererUUID = ""
	request.rendererUUIDVerified = false
	if claimedUUID != "" {
		if m.rendererUUIDBoundToIP(claimedUUID, request.ClientIP, overrides) {
			request.RendererUUID = claimedUUID
			request.rendererUUIDVerified = true
		}
		return request
	}

	m.mu.Lock()
	rememberedUUID := normalizeRendererUUID(m.recentClients[strings.TrimSpace(request.ClientIP)].RendererUUID)
	m.mu.Unlock()
	if rememberedUUID != "" && m.rendererUUIDBoundToIP(rememberedUUID, request.ClientIP, overrides) {
		request.RendererUUID = rememberedUUID
		request.rendererUUIDVerified = true
	}
	return request
}

func (m *Manager) rendererUUIDBoundToIP(rendererUUID, clientIP string, overrides []rendererDeviceOverride) bool {
	clientIP = strings.TrimSpace(clientIP)
	parsedIP := net.ParseIP(clientIP)
	if rendererUUID == "" || parsedIP == nil {
		return false
	}
	for _, override := range overrides {
		if normalizeRendererUUID(override.RendererUUID) == rendererUUID && strings.TrimSpace(override.IPAddress) == clientIP {
			return true
		}
	}
	service := m.controlPointSnapshot()
	if service == nil {
		return false
	}
	for _, renderer := range service.Registry().List(time.Time{}) {
		if normalizeRendererUUID(renderer.UDN) == rendererUUID && renderer.SourceIP.Equal(parsedIP) {
			return true
		}
	}
	return false
}
