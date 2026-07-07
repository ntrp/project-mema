package dlna

import (
	"context"
	"encoding/xml"
	"net/http"
	"strings"

	"media-manager/internal/storage"
)

func (m *Manager) Handler() http.Handler {
	mux := http.NewServeMux()
	dispatcher := m.SOAPDispatcher()
	for _, prefix := range []string{"", "/dlna"} {
		mux.HandleFunc(prefix+"/rootDesc.xml", m.rootDescription)
		mux.HandleFunc(prefix+"/contentDirectory.xml", serveXML(ContentDirectorySCPDXML))
		mux.HandleFunc(prefix+"/connectionManager.xml", serveXML(ConnectionManagerSCPDXML))
		mux.HandleFunc(prefix+"/mediaReceiverRegistrar.xml", serveXML(MediaReceiverRegistrarSCPDXML))
		mux.Handle(prefix+"/control/content-directory", dispatcher)
		mux.Handle(prefix+"/control/connection-manager", dispatcher)
		mux.Handle(prefix+"/control/media-receiver-registrar", dispatcher)
		mux.HandleFunc(prefix+"/resource/", m.resource)
		mux.HandleFunc(prefix+"/artwork/", m.artwork)
		mux.HandleFunc(prefix+"/subtitle/", m.subtitle)
		mux.HandleFunc(prefix+"/events/content-directory", m.eventHandler)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.recordHTTPClient(r)
		applyRendererHeaders(w, m.RendererProfileFromRequest(r))
		mux.ServeHTTP(w, r)
	})
}

func (m *Manager) eventHandler(w http.ResponseWriter, r *http.Request) {
	if m.RendererProfileFromRequest(r).DisableEventing {
		http.Error(w, "renderer profile disables eventing", http.StatusForbidden)
		return
	}
	m.events.Handle(w, r)
}

func (m *Manager) rootDescription(w http.ResponseWriter, r *http.Request) {
	settings := m.currentSettings(r.Context())
	payload, err := RootDeviceXML(settings.FriendlyName, settings.DeviceUUID, requestBaseURL(r))
	if err != nil {
		http.Error(w, "DLNA root description unavailable", http.StatusInternalServerError)
		return
	}
	writeXML(w, payload)
}

func (m *Manager) currentSettings(ctx context.Context) storage.DLNASettings {
	if m.store == nil {
		return storage.DLNASettings{FriendlyName: storage.DefaultDLNAFriendlyName}
	}
	settings, err := m.store.GetDLNASettings(ctx)
	if err != nil {
		return storage.DLNASettings{FriendlyName: storage.DefaultDLNAFriendlyName}
	}
	return settings
}

func serveXML(build func() ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := build()
		if err != nil {
			http.Error(w, "DLNA document unavailable", http.StatusInternalServerError)
			return
		}
		writeXML(w, payload)
	}
}

func writeXML(w http.ResponseWriter, payload []byte) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(payload)
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if value := forwardedHeader(r, "X-Forwarded-Proto"); value != "" {
		scheme = value
	}
	host := r.Host
	if value := forwardedHeader(r, "X-Forwarded-Host"); value != "" {
		host = value
	}
	return scheme + "://" + host
}

func forwardedHeader(r *http.Request, name string) string {
	value := strings.TrimSpace(r.Header.Get(name))
	if index := strings.Index(value, ","); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return value
}
