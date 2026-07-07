package controlpoint

import (
	"strings"
	"testing"

	"media-manager/internal/dlna"
)

func TestMediaServerDescriptionDoesNotAdvertiseRendererServices(t *testing.T) {
	payload, err := dlna.RootDeviceXML("Mema", "uuid:test", "http://127.0.0.1:18080")
	if err != nil {
		t.Fatal(err)
	}
	body := string(payload)
	for _, forbidden := range []string{MediaRendererDevice, AVTransportService, RenderingService} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("server description advertises renderer service %q:\n%s", forbidden, body)
		}
	}
}
