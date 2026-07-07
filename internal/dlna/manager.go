package dlna

import (
	"context"
	"fmt"
	"net"
	"sync"

	"media-manager/internal/storage"
)

type Status struct {
	Running         bool
	BoundInterfaces []string
	AdvertisedURLs  []string
	LastError       *string
}

type Manager struct {
	store   *storage.SettingsStore
	baseURL string
	mu      sync.Mutex
	status  Status
}

func NewManager(store *storage.SettingsStore, baseURL string) *Manager {
	return &Manager{store: store, baseURL: baseURL}
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
	if !settings.Enabled {
		m.status = Status{}
		return nil
	}
	status, err := m.statusForSettings(settings)
	if err != nil {
		message := err.Error()
		m.status = Status{LastError: &message}
		return err
	}
	m.status = status
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
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

func (m *Manager) statusForSettings(settings storage.DLNASettings) (Status, error) {
	names := settings.Interfaces
	if len(names) == 0 {
		discovered, err := activeInterfaceNames()
		if err != nil {
			return Status{}, err
		}
		names = discovered
	}
	urls := []string{}
	if m.baseURL != "" {
		urls = append(urls, m.baseURL+"/dlna/rootDesc.xml")
	}
	return Status{
		Running:         true,
		BoundInterfaces: append([]string{}, names...),
		AdvertisedURLs:  urls,
	}, nil
}

func activeInterfaceNames() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("network interfaces unavailable: %w", err)
	}
	names := []string{}
	for _, item := range interfaces {
		if item.Flags&net.FlagUp == 0 || item.Flags&net.FlagLoopback != 0 {
			continue
		}
		names = append(names, item.Name)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no active non-loopback interfaces found")
	}
	return names, nil
}

func (m *Manager) setError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	message := err.Error()
	m.status = Status{LastError: &message}
}
