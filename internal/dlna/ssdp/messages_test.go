package ssdp

import (
	"strings"
	"testing"
)

func TestSearchResponseIncludesRequiredHeaders(t *testing.T) {
	config := Config{UUID: "device-1", AnnounceSeconds: 1800}
	iface := Interface{Name: "en0", Location: "http://192.168.1.2:18080/dlna/rootDesc.xml"}

	packet, ok := SearchResponse(config, iface, "upnp:rootdevice")
	if !ok {
		t.Fatal("expected rootdevice search response")
	}
	text := PacketText(packet, "HTTP/1.1 200 OK")
	for _, header := range []string{"CACHE-CONTROL:", "EXT:", "LOCATION:", "SERVER:", "ST:", "USN:"} {
		if !strings.Contains(text, header) {
			t.Fatalf("response missing %s:\n%s", header, text)
		}
	}
	if !strings.Contains(text, "USN: uuid:device-1::upnp:rootdevice") {
		t.Fatalf("response USN mismatch:\n%s", text)
	}
}

func TestAliveAndByebyeAdvertiseAllTargets(t *testing.T) {
	config := Config{UUID: "device-1", AnnounceSeconds: 1800}
	iface := Interface{Name: "en0", Location: "http://192.168.1.2:18080/dlna/rootDesc.xml"}

	alive := AlivePackets(config, iface)
	byebye := ByebyePackets(config, iface)
	want := len(ServiceTargets) + 1
	if len(alive) != want || len(byebye) != want {
		t.Fatalf("alive=%d byebye=%d want %d", len(alive), len(byebye), want)
	}
	for _, packet := range alive {
		if packet.Headers["NTS"] != "ssdp:alive" || packet.Headers["LOCATION"] == "" {
			t.Fatalf("alive packet = %#v", packet)
		}
	}
	for _, packet := range byebye {
		if packet.Headers["NTS"] != "ssdp:byebye" {
			t.Fatalf("byebye packet = %#v", packet)
		}
	}
}

func TestSupportsExpectedSearchTargets(t *testing.T) {
	for _, target := range []string{"ssdp:all", "uuid:device-1", "upnp:rootdevice", MediaServer, ContentDir, Connection} {
		if !SupportsTarget("device-1", target) {
			t.Fatalf("expected target %q to be supported", target)
		}
	}
	if SupportsTarget("device-1", "urn:schemas-upnp-org:service:Unknown:1") {
		t.Fatal("unexpected unsupported target match")
	}
}
