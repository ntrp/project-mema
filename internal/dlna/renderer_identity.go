package dlna

import "strings"

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
	if rendererUUID == "" || clientIP == "" {
		return false
	}
	for _, override := range overrides {
		if normalizeRendererUUID(override.RendererUUID) == rendererUUID && strings.TrimSpace(override.IPAddress) == clientIP {
			return true
		}
	}
	return false
}
