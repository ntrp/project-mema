package dlna

import (
	"net"
	"net/http"
)

const maxActiveStreams = 8

func (m *Manager) allowRequest(r *http.Request) bool {
	ip := net.ParseIP(clientIP(r))
	if ip == nil {
		return false
	}
	m.mu.Lock()
	nets := append([]*net.IPNet{}, m.allowedNets...)
	m.mu.Unlock()
	for _, network := range nets {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func parseAllowedCIDRs(values []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		_, network, err := net.ParseCIDR(value)
		if err == nil {
			nets = append(nets, network)
		}
	}
	return nets
}
