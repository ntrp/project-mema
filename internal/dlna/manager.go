package dlna

import (
	"context"
	"net"
	"net/url"
	"strings"
	"sync"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"
)

type Status struct {
	Running         bool
	BoundInterfaces []string
	AdvertisedURLs  []string
	LastError       *string
}

type Manager struct {
	store     *storage.SettingsStore
	source    content.LibrarySource
	baseURL   string
	httpPort  string
	thumbDir  string
	events    *EventManager
	profileOverrides map[string]string
	startSSDP func(context.Context, ssdp.Config) (ssdpRuntime, error)
	ssdp      ssdpRuntime
	mu        sync.Mutex
	status    Status
}

func NewManager(store *storage.SettingsStore, baseURL string) *Manager {
	return &Manager{
		store:     store,
		baseURL:   strings.TrimRight(baseURL, "/"),
		httpPort:  portFromBaseURL(baseURL),
		thumbDir:  ".data/dlna-thumbnails",
		events:    NewEventManager(),
		startSSDP: startSSDPRuntime,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	settings, err := m.store.GetDLNASettings(ctx)
	if err != nil {
		m.setError(err)
		return err
	}
	return m.ApplySettings(ctx, settings)
}

func (m *Manager) ApplySettings(ctx context.Context, settings storage.DLNASettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopSSDP(ctx)
	if !settings.Enabled {
		m.status = Status{}
		return nil
	}
	runtime, err := m.startSSDP(ctx, ssdp.Config{
		FriendlyName:    settings.FriendlyName,
		HTTPPort:        m.httpPort,
		Interfaces:      settings.Interfaces,
		AnnounceSeconds: settings.AnnounceIntervalSeconds,
		UUID:            deviceUUID(settings),
	})
	if err != nil {
		message := err.Error()
		m.status = Status{LastError: &message}
		return err
	}
	m.ssdp = runtime
	m.status = m.statusForSettings(settings, runtime.Interfaces())
	return nil
}

func deviceUUID(settings storage.DLNASettings) string {
	if strings.TrimSpace(settings.DeviceUUID) != "" {
		return settings.DeviceUUID
	}
	return "00000000-0000-4000-8000-000000000001"
}

func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopSSDP(ctx)
	m.status = Status{}
	return nil
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	status := m.status
	status.BoundInterfaces = append([]string{}, status.BoundInterfaces...)
	status.AdvertisedURLs = append([]string{}, status.AdvertisedURLs...)
	return status
}

func (m *Manager) statusForSettings(settings storage.DLNASettings, ifaces []ssdp.Interface) Status {
	names := settings.Interfaces
	if len(names) == 0 {
		names = make([]string, 0, len(ifaces))
		for _, item := range ifaces {
			names = append(names, item.Name)
		}
	}
	urls := make([]string, 0, len(ifaces))
	for _, item := range ifaces {
		urls = append(urls, item.Location)
	}
	if len(urls) == 0 && m.baseURL != "" {
		urls = append(urls, m.baseURL+"/dlna/rootDesc.xml")
	}
	return Status{Running: true, BoundInterfaces: append([]string{}, names...), AdvertisedURLs: urls}
}

func (m *Manager) stopSSDP(ctx context.Context) {
	if m.ssdp == nil {
		return
	}
	_ = m.ssdp.Stop(ctx)
	m.ssdp = nil
}

func (m *Manager) setError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	message := err.Error()
	m.status = Status{LastError: &message}
}

type ssdpRuntime interface {
	Stop(context.Context) error
	Interfaces() []ssdp.Interface
}

func startSSDPRuntime(ctx context.Context, config ssdp.Config) (ssdpRuntime, error) {
	return ssdp.Start(ctx, config)
}

func portFromBaseURL(baseURL string) string {
	parsed, err := url.Parse(baseURL)
	if err == nil && parsed.Port() != "" {
		return parsed.Port()
	}
	_, port, err := net.SplitHostPort(strings.TrimPrefix(baseURL, "http://"))
	if err == nil && port != "" {
		return port
	}
	if strings.HasPrefix(baseURL, "https://") {
		return "443"
	}
	return "80"
}
