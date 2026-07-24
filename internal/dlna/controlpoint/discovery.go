package controlpoint

import (
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// Registry retains renderers only while they are visible on a local interface.
type Registry struct {
	mu      sync.RWMutex
	devices map[string]entry
	now     func() time.Time
}

type entry struct {
	renderer Renderer
	source   net.IP
	scope    string
	expires  time.Time
}

func NewRegistry() *Registry { return &Registry{devices: make(map[string]entry), now: time.Now} }

func (r *Registry) Upsert(renderer Renderer, source net.IP, maxAge time.Duration) {
	if r == nil || strings.TrimSpace(renderer.UDN) == "" || source == nil {
		return
	}
	if maxAge <= 0 {
		maxAge = 30 * time.Second
	}
	if maxAge > 30*time.Minute {
		maxAge = 30 * time.Minute
	}
	now := r.now()
	renderer.SourceIP = append(net.IP(nil), source...)
	renderer.LastSeen = now
	renderer.ExpiresAt = now.Add(maxAge)
	r.mu.Lock()
	r.devices[renderer.UDN] = entry{
		renderer: renderer, source: append(net.IP(nil), source...),
		scope: scopeFor(source), expires: renderer.ExpiresAt,
	}
	r.mu.Unlock()
}
func (r *Registry) RemoveFrom(udn string, source net.IP) bool {
	if r == nil || source == nil {
		return false
	}
	udn = strings.TrimSpace(udn)
	r.mu.Lock()
	defer r.mu.Unlock()
	device, ok := r.devices[udn]
	if !ok || !device.source.Equal(source) {
		return false
	}
	delete(r.devices, udn)
	return true
}
func (r *Registry) List(now time.Time) []Renderer {
	if r == nil {
		return []Renderer{}
	}
	if now.IsZero() {
		now = r.now()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Renderer, 0, len(r.devices))
	for udn, device := range r.devices {
		if !device.expires.After(now) {
			delete(r.devices, udn)
			continue
		}
		out = append(out, device.renderer)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FriendlyName < out[j].FriendlyName })
	return out
}
func (r *Registry) Get(udn string, now time.Time) (Renderer, bool) {
	if r == nil {
		return Renderer{}, false
	}
	if now.IsZero() {
		now = r.now()
	}
	r.mu.RLock()
	device, ok := r.devices[udn]
	r.mu.RUnlock()
	return device.renderer, ok && device.expires.After(now)
}
func scopeFor(ip net.IP) string {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4.Mask(net.CIDRMask(24, 32)).String() + "/24"
	}
	return ip.String()
}
func (r *Registry) InScope(udn string, source net.IP, now time.Time) bool {
	if r == nil || source == nil {
		return false
	}
	r.mu.RLock()
	device, ok := r.devices[udn]
	r.mu.RUnlock()
	return ok && device.scope == scopeFor(source) && device.expires.After(now)
}
