package dlna

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestRootDeviceXMLIncludesStableUDNAndServices(t *testing.T) {
	payload, err := RootDeviceXML("Mema Test", "device-1", "http://127.0.0.1:18080")
	if err != nil {
		t.Fatalf("RootDeviceXML returned error: %v", err)
	}
	var doc DeviceDocument
	if err := xml.Unmarshal(payload, &doc); err != nil {
		t.Fatalf("root XML did not parse: %v\n%s", err, payload)
	}
	if doc.Device.UDN != "uuid:device-1" || doc.Device.FriendlyName != "Mema Test" {
		t.Fatalf("device = %#v", doc.Device)
	}
	if len(doc.Device.Services) != 3 {
		t.Fatalf("services = %#v", doc.Device.Services)
	}
	if !strings.Contains(string(payload), "/dlna/contentDirectory.xml") {
		t.Fatalf("root XML missing ContentDirectory SCPD URL:\n%s", payload)
	}
}

func TestSCPDDocumentsParse(t *testing.T) {
	for name, build := range map[string]func() ([]byte, error){
		"content":    ContentDirectorySCPDXML,
		"connection": ConnectionManagerSCPDXML,
		"registrar":  MediaReceiverRegistrarSCPDXML,
	} {
		payload, err := build()
		if err != nil {
			t.Fatalf("%s SCPD returned error: %v", name, err)
		}
		var doc SCPDDocument
		if err := xml.Unmarshal(payload, &doc); err != nil {
			t.Fatalf("%s SCPD did not parse: %v\n%s", name, err, payload)
		}
		if len(doc.Actions) == 0 || len(doc.State) == 0 {
			t.Fatalf("%s SCPD incomplete: %#v", name, doc)
		}
	}
}
