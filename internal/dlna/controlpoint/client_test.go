package controlpoint

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseRendererDescriptionFindsControlURLs(t *testing.T) {
	payload := []byte(`<root><device><deviceType>urn:schemas-upnp-org:device:MediaRenderer:1</deviceType><friendlyName>Scenario Renderer</friendlyName><UDN>uuid:renderer</UDN><serviceList><service><serviceType>urn:schemas-upnp-org:service:AVTransport:1</serviceType><controlURL>/upnp/av</controlURL></service><service><serviceType>urn:schemas-upnp-org:service:RenderingControl:1</serviceType><controlURL>/upnp/render</controlURL></service></serviceList></device></root>`)

	renderer, err := ParseRendererDescription("http://127.0.0.1:1400/root.xml", payload)

	if err != nil {
		t.Fatal(err)
	}
	if renderer.FriendlyName != "Scenario Renderer" || renderer.AVTransportURL != "http://127.0.0.1:1400/upnp/av" ||
		renderer.RenderingControlURL != "http://127.0.0.1:1400/upnp/render" {
		t.Fatalf("renderer = %#v", renderer)
	}
}

func TestClientSendsAVTransportActionsToFakeRenderer(t *testing.T) {
	var actions, payloads []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actions = append(actions, r.Header.Get("SOAPACTION"))
		payload := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(payload)
		payloads = append(payloads, string(payload))
		if strings.Contains(r.Header.Get("SOAPACTION"), "SetAVTransportURI") &&
			!strings.Contains(string(payload), "http://mema.local/dlna/resource/item") {
			t.Fatalf("payload = %s", string(payload))
		}
		_, _ = w.Write([]byte(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body/></s:Envelope>`))
	}))
	defer server.Close()
	renderer := Renderer{AVTransportURL: server.URL, RenderingControlURL: server.URL}
	client := NewClient(server.Client())

	for _, run := range []func() error{
		func() error {
			return client.SetAVTransportURI(context.Background(), renderer, "http://mema.local/dlna/resource/item", "")
		},
		func() error { return client.Play(context.Background(), renderer) },
		func() error { return client.Pause(context.Background(), renderer) },
		func() error { return client.Seek(context.Background(), renderer, "00:01:00") },
		func() error { return client.Stop(context.Background(), renderer) },
		func() error { return client.Next(context.Background(), renderer) },
		func() error { return client.Previous(context.Background(), renderer) },
	} {
		if err := run(); err != nil {
			t.Fatal(err)
		}
	}
	if len(actions) != 7 || !strings.Contains(actions[0], "SetAVTransportURI") || !strings.Contains(actions[1], "Play") {
		t.Fatalf("actions = %#v", actions)
	}
	if !strings.HasPrefix(payloads[0], `<?xml version="1.0"`) || !strings.Contains(payloads[0], `s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"`) {
		t.Fatalf("SOAP compatibility attributes missing: %s", payloads[0])
	}
	instance, uri, metadata := strings.Index(payloads[0], "<InstanceID>"), strings.Index(payloads[0], "<CurrentURI>"), strings.Index(payloads[0], "<CurrentURIMetaData>")
	if instance < 0 || instance > uri || uri > metadata {
		t.Fatalf("SetAVTransportURI arguments are out of order: %s", payloads[0])
	}
}

func TestClientReportsNamespacedUPnPFault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><s:Fault><detail><UPnPError xmlns="urn:schemas-upnp-org:control-1-0"><errorCode>714</errorCode><errorDescription>Illegal MIME-type</errorDescription></UPnPError></detail></s:Fault></s:Body></s:Envelope>`))
	}))
	defer server.Close()

	err := NewClient(server.Client()).Stop(context.Background(), Renderer{AVTransportURL: server.URL})
	if err == nil || !strings.Contains(err.Error(), "Illegal MIME-type (UPnP 714)") {
		t.Fatalf("error = %v", err)
	}
}

func TestClientSendsRenderingControlActionsToFakeRenderer(t *testing.T) {
	var body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(payload)
		body += string(payload)
		if !strings.Contains(r.Header.Get("SOAPACTION"), RenderingService) {
			t.Fatalf("SOAPACTION = %s", r.Header.Get("SOAPACTION"))
		}
	}))
	defer server.Close()
	client := NewClient(server.Client())
	renderer := Renderer{RenderingControlURL: server.URL}

	if err := client.SetVolume(context.Background(), renderer, 35); err != nil {
		t.Fatal(err)
	}
	if err := client.SetMute(context.Background(), renderer, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "<DesiredVolume>35</DesiredVolume>") || !strings.Contains(body, "<DesiredMute>1</DesiredMute>") {
		t.Fatalf("body = %s", body)
	}
}
