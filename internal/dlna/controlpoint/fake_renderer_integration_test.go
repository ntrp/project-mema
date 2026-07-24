package controlpoint

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestFakeRendererDiscoveryStateActionsExpiryAndFaults(t *testing.T) {
	var fault atomic.Bool
	var actions atomic.Int32
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/root.xml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `<root><device><deviceType>%s</deviceType><friendlyName>Fake TV</friendlyName><manufacturer>Mema Test</manufacturer><modelName>F1</modelName><UDN>uuid:fake</UDN><serviceList><service><serviceType>%s</serviceType><controlURL>/av</controlURL><SCPDURL>/av.xml</SCPDURL></service><service><serviceType>%s</serviceType><controlURL>/render</controlURL><SCPDURL>/render.xml</SCPDURL></service></serviceList></device></root>`, MediaRendererDevice, AVTransportService, RenderingService)
	})
	mux.HandleFunc("/av.xml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<scpd><actionList><action><name>SetAVTransportURI</name></action><action><name>Play</name></action><action><name>GetTransportInfo</name></action><action><name>GetPositionInfo</name></action></actionList></scpd>`))
	})
	mux.HandleFunc("/render.xml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<scpd><actionList><action><name>GetVolume</name></action><action><name>GetMute</name></action><action><name>SetVolume</name></action></actionList></scpd>`))
	})
	mux.HandleFunc("/av", func(w http.ResponseWriter, r *http.Request) {
		action := r.Header.Get("SOAPACTION")
		if fault.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`<Envelope><Body><Fault><detail><UPnPError><errorCode>701</errorCode><errorDescription>Transition not available</errorDescription></UPnPError></detail></Fault></Body></Envelope>`))
			return
		}
		actions.Add(1)
		if strings.Contains(action, "GetTransportInfo") {
			_, _ = w.Write([]byte(`<Envelope><Body><GetTransportInfoResponse><CurrentTransportState>PLAYING</CurrentTransportState><CurrentTransportStatus>OK</CurrentTransportStatus><CurrentSpeed>1</CurrentSpeed></GetTransportInfoResponse></Body></Envelope>`))
			return
		}
		if strings.Contains(action, "GetPositionInfo") {
			_, _ = w.Write([]byte(`<Envelope><Body><GetPositionInfoResponse><Track>1</Track><RelTime>0:00:12</RelTime><TrackDuration>0:42:00</TrackDuration></GetPositionInfoResponse></Body></Envelope>`))
		}
	})
	mux.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("SOAPACTION"), "GetVolume") {
			_, _ = w.Write([]byte(`<Envelope><Body><GetVolumeResponse><CurrentVolume>33</CurrentVolume></GetVolumeResponse></Body></Envelope>`))
			return
		}
		_, _ = w.Write([]byte(`<Envelope><Body><GetMuteResponse><CurrentMute>1</CurrentMute></GetMuteResponse></Body></Envelope>`))
	})

	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	service := NewService(nil, []*net.IPNet{network}, server.Client())
	packet := []byte("HTTP/1.1 200 OK\r\nST: " + MediaRendererDevice + "\r\nUSN: uuid:fake::" + MediaRendererDevice + "\r\nLOCATION: " + server.URL + "/root.xml\r\nCACHE-CONTROL: max-age=1\r\n\r\n")
	service.process(context.Background(), packet, net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"))
	var renderer Renderer
	for range 50 {
		if value, ok := service.registry.Get("uuid:fake", time.Time{}); ok {
			renderer = value
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if renderer.FriendlyName != "Fake TV" || !renderer.AVTransportActionsKnown || !renderer.RenderingActionsKnown ||
		!SupportsRendering(renderer, "SetVolume") || SupportsRendering(renderer, "SetMute") {
		t.Fatalf("discovered renderer = %#v", renderer)
	}
	client := service.Client()
	transport, err := client.GetTransportInfo(context.Background(), renderer)
	if err != nil || transport.State != "PLAYING" {
		t.Fatalf("transport = %#v, %v", transport, err)
	}
	position, err := client.GetPositionInfo(context.Background(), renderer)
	if err != nil || position.RelTime != "0:00:12" {
		t.Fatalf("position = %#v, %v", position, err)
	}
	volume, err := client.GetVolume(context.Background(), renderer)
	if err != nil || volume != 33 {
		t.Fatalf("volume = %d, %v", volume, err)
	}
	if err := client.Play(context.Background(), renderer); err != nil || actions.Load() < 3 {
		t.Fatalf("play failed: %v", err)
	}
	fault.Store(true)
	err = client.Play(context.Background(), renderer)
	var soapFault *FaultError
	if !errors.As(err, &soapFault) || soapFault.Description != "Transition not available" {
		t.Fatalf("fault = %#v", err)
	}
	if got := service.registry.List(renderer.ExpiresAt.Add(time.Millisecond)); len(got) != 0 {
		t.Fatalf("expired renderers = %#v", got)
	}
}

func TestDescriptionRedirectCannotChangeTrustedSenderWithinCIDR(t *testing.T) {
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.2/private", http.StatusFound)
	}))
	defer redirect.Close()
	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	service := NewService(nil, []*net.IPNet{network}, redirect.Client())
	if _, err := service.fetchURL(context.Background(), redirect.URL, net.ParseIP("127.0.0.1")); err == nil {
		t.Fatal("same-CIDR redirect changed the trusted SSDP sender")
	}
}

func TestDiscoveryRejectsLocationsOutsideConfiguredCIDR(t *testing.T) {
	_, network, _ := net.ParseCIDR("192.168.0.0/16")
	service := NewService(nil, []*net.IPNet{network}, nil)
	packet := []byte("HTTP/1.1 200 OK\r\nST: " + MediaRendererDevice + "\r\nUSN: uuid:fake::" + MediaRendererDevice + "\r\nLOCATION: http://10.0.0.2/root.xml\r\n\r\n")
	service.process(context.Background(), packet, net.ParseIP("10.0.0.2"), net.ParseIP("192.168.1.2"))
	time.Sleep(10 * time.Millisecond)
	if got := service.registry.List(time.Time{}); len(got) != 0 {
		t.Fatalf("out-of-scope renderer discovered: %#v", got)
	}
}
