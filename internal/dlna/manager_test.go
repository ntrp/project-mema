package dlna

import (
	"context"
	"testing"

	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"
)

func TestManagerKeepsDLNADisabledByDefault(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")

	err := manager.ApplySettings(context.Background(), storage.DLNASettings{})
	if err != nil {
		t.Fatalf("ApplySettings returned error: %v", err)
	}
	if status := manager.Status(); status.Running {
		t.Fatalf("status = %#v, want stopped", status)
	}
}

func TestManagerStartsWithConfiguredInterfaces(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.startSSDP = func(ctx context.Context, config ssdp.Config) (ssdpRuntime, error) {
		if config.HTTPPort != "18080" {
			t.Fatalf("HTTPPort = %q", config.HTTPPort)
		}
		return &fakeSSDPRuntime{ifaces: []ssdp.Interface{{
			Name:     "lo0",
			Location: "http://127.0.0.1:18080/dlna/rootDesc.xml",
		}}}, nil
	}

	err := manager.ApplySettings(context.Background(), storage.DLNASettings{
		Enabled:    true,
		Interfaces: []string{"lo0"},
	})
	if err != nil {
		t.Fatalf("ApplySettings returned error: %v", err)
	}
	status := manager.Status()
	if !status.Running || len(status.BoundInterfaces) != 1 || status.BoundInterfaces[0] != "lo0" {
		t.Fatalf("status = %#v, want running on lo0", status)
	}
	if len(status.AdvertisedURLs) != 1 || status.AdvertisedURLs[0] != "http://127.0.0.1:18080/dlna/rootDesc.xml" {
		t.Fatalf("advertised urls = %#v", status.AdvertisedURLs)
	}
}

func TestManagerStopsSSDPOnDisable(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	runtime := &fakeSSDPRuntime{ifaces: []ssdp.Interface{{Name: "en0"}}}
	manager.startSSDP = func(ctx context.Context, config ssdp.Config) (ssdpRuntime, error) {
		return runtime, nil
	}

	if err := manager.ApplySettings(context.Background(), storage.DLNASettings{Enabled: true}); err != nil {
		t.Fatalf("ApplySettings enable returned error: %v", err)
	}
	if err := manager.ApplySettings(context.Background(), storage.DLNASettings{}); err != nil {
		t.Fatalf("ApplySettings disable returned error: %v", err)
	}
	if !runtime.stopped {
		t.Fatal("expected SSDP runtime to stop")
	}
}

type fakeSSDPRuntime struct {
	ifaces  []ssdp.Interface
	stopped bool
}

func (f *fakeSSDPRuntime) Interfaces() []ssdp.Interface {
	return f.ifaces
}

func (f *fakeSSDPRuntime) Stop(ctx context.Context) error {
	f.stopped = true
	return nil
}
