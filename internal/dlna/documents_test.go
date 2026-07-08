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
	if doc.Device.ModelDescription == "" || doc.Device.ManufacturerURL == "" || doc.Device.ModelURL == "" {
		t.Fatalf("device identity descriptors missing: %#v", doc.Device)
	}
	if len(doc.Device.Services) != 3 {
		t.Fatalf("services = %#v", doc.Device.Services)
	}
	if len(doc.Device.Icons) != 4 {
		t.Fatalf("icons = %#v", doc.Device.Icons)
	}
	if !strings.Contains(string(payload), "/dlna/contentDirectory.xml") {
		t.Fatalf("root XML missing ContentDirectory SCPD URL:\n%s", payload)
	}
	for _, want := range []string{
		`xmlns:sec="http://www.sec.co.kr/"`,
		`<dlna:X_DLNADOC>DMS-1.50</dlna:X_DLNADOC>`,
		`<dlna:X_DLNADOC>M-DMS-1.50</dlna:X_DLNADOC>`,
		`<sec:ProductCap>smi,DCM10,getMediaInfo.sec,getCaptionInfo.sec</sec:ProductCap>`,
	} {
		if !strings.Contains(string(payload), want) {
			t.Fatalf("root XML missing %q:\n%s", want, payload)
		}
	}
	for _, forbidden := range []string{"/dlna/events/content-directory", "/dlna/events/connection-manager", "/dlna/events/media-receiver-registrar"} {
		if strings.Contains(string(payload), forbidden) {
			t.Fatalf("root XML advertises unsupported event URL %q:\n%s", forbidden, payload)
		}
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
