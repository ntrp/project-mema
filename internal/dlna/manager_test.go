package dlna

import (
	"context"
	"testing"

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
