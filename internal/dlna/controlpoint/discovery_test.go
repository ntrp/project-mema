package controlpoint

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegistryExpiresAndScopesRenderers(t *testing.T) {
	r := NewRegistry()
	base := time.Unix(100, 0)
	r.now = func() time.Time { return base }
	device := Renderer{UDN: "uuid:one", FriendlyName: "Living room"}
	r.Upsert(device, net.ParseIP("192.168.1.42"), time.Minute)
	if !r.InScope(device.UDN, net.ParseIP("192.168.1.99"), base.Add(time.Second)) {
		t.Fatal("renderer should be on LAN")
	}
	if r.InScope(device.UDN, net.ParseIP("192.168.2.1"), base.Add(time.Second)) {
		t.Fatal("renderer leaked across LAN")
	}
	if got := r.List(base.Add(2 * time.Minute)); len(got) != 0 {
		t.Fatalf("expired renderers = %#v", got)
	}
	if _, ok := r.Get(device.UDN, base.Add(2*time.Minute)); ok {
		t.Fatal("expired renderer returned")
	}
}

func TestRegistryUpdatesAdvertisementExpiry(t *testing.T) {
	r := NewRegistry()
	r.now = func() time.Time { return time.Unix(0, 0) }
	r.Upsert(Renderer{UDN: "uuid:one"}, net.ParseIP("10.0.0.2"), 0)
	if _, ok := r.Get("uuid:one", time.Unix(29, 0)); !ok {
		t.Fatal("default lease should be retained")
	}
	r.Upsert(Renderer{UDN: "uuid:one", FriendlyName: "new"}, net.ParseIP("10.0.0.2"), time.Hour)
	got, ok := r.Get("uuid:one", time.Unix(60, 0))
	if !ok || got.FriendlyName != "new" {
		t.Fatalf("updated renderer = %#v, %v", got, ok)
	}
}

func TestParseAdvertisementHandlesNotifyLeaseAndByebye(t *testing.T) {
	alive := []byte("NOTIFY * HTTP/1.1\r\nNT: " + MediaRendererDevice + "\r\nNTS: ssdp:alive\r\nUSN: uuid:tv::" + MediaRendererDevice + "\r\nLOCATION: http://192.168.1.4/root.xml\r\nCACHE-CONTROL: max-age=90\r\n\r\n")
	advertisement, ok := parseAdvertisement(alive)
	if !ok || advertisement.Byebye || advertisement.UDN != "uuid:tv" || advertisement.MaxAge != 90*time.Second {
		t.Fatalf("alive advertisement = %#v, %v", advertisement, ok)
	}
	byebye := []byte("NOTIFY * HTTP/1.1\r\nNT: " + MediaRendererDevice + "\r\nNTS: ssdp:byebye\r\nUSN: uuid:tv::" + MediaRendererDevice + "\r\n\r\n")
	advertisement, ok = parseAdvertisement(byebye)
	if !ok || !advertisement.Byebye {
		t.Fatalf("byebye advertisement = %#v, %v", advertisement, ok)
	}
	responseByebye := []byte("HTTP/1.1 200 OK\r\nST: " + MediaRendererDevice + "\r\nNTS: ssdp:byebye\r\nUSN: uuid:tv::" + MediaRendererDevice + "\r\n\r\n")
	if _, ok := parseAdvertisement(responseByebye); ok {
		t.Fatal("M-SEARCH response was trusted as a byebye notification")
	}
	wrongTarget := []byte("NOTIFY * HTTP/1.1\r\nNT: upnp:rootdevice\r\nNTS: ssdp:byebye\r\nUSN: uuid:tv::" + MediaRendererDevice + "\r\n\r\n")
	if _, ok := parseAdvertisement(wrongTarget); ok {
		t.Fatal("renderer UDN bypassed the SSDP target check")
	}
}

func TestDiscoveryKeepsCapabilityProvenanceWhenOneSCPDFails(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/root.xml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `<root><device><deviceType>%s</deviceType><friendlyName>Partial</friendlyName><UDN>uuid:partial</UDN><serviceList><service><serviceType>%s</serviceType><controlURL>/av</controlURL><SCPDURL>/av.xml</SCPDURL></service><service><serviceType>%s</serviceType><controlURL>/render</controlURL><SCPDURL>/render.xml</SCPDURL></service></serviceList></device></root>`, MediaRendererDevice, AVTransportService, RenderingService)
	})
	mux.HandleFunc("/av.xml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<scpd><actionList><action><name>Play</name></action></actionList></scpd>`))
	})
	mux.HandleFunc("/render.xml", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) })
	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	service := NewService(nil, []*net.IPNet{network}, server.Client())
	service.process(context.Background(), rendererPacket("uuid:partial", server.URL+"/root.xml"), net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"))

	var renderer Renderer
	for range 50 {
		if value, ok := service.Registry().Get("uuid:partial", time.Time{}); ok {
			renderer = value
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !renderer.AVTransportActionsKnown || renderer.RenderingActionsKnown || SupportsTransport(renderer, "Pause") || !SupportsRendering(renderer, "SetMute") {
		t.Fatalf("partial capabilities = %#v", renderer)
	}
}

func TestDiscoveryRejectsSenderLocationMismatchAndForgedByebye(t *testing.T) {
	_, network, _ := net.ParseCIDR("192.168.1.0/24")
	service := NewService(nil, []*net.IPNet{network}, nil)
	owner, attacker := net.ParseIP("192.168.1.20"), net.ParseIP("192.168.1.21")
	service.Registry().Upsert(Renderer{UDN: "uuid:tv"}, owner, time.Minute)
	byebye := []byte("NOTIFY * HTTP/1.1\r\nNT: " + MediaRendererDevice + "\r\nNTS: ssdp:byebye\r\nUSN: uuid:tv::" + MediaRendererDevice + "\r\n\r\n")
	service.process(context.Background(), byebye, attacker, owner)
	if _, ok := service.Registry().Get("uuid:tv", time.Time{}); !ok {
		t.Fatal("same-CIDR attacker removed another sender's renderer")
	}
	service.process(context.Background(), rendererPacket("uuid:attacker", "http://192.168.1.20/root.xml"), attacker, owner)
	time.Sleep(10 * time.Millisecond)
	if _, ok := service.Registry().Get("uuid:attacker", time.Time{}); ok {
		t.Fatal("sender advertised a LOCATION owned by another host")
	}
	service.process(context.Background(), byebye, owner, owner)
	if _, ok := service.Registry().Get("uuid:tv", time.Time{}); ok {
		t.Fatal("owner byebye did not remove renderer")
	}
}

func TestServiceStartsAndStopsUDPListeners(t *testing.T) {
	iface, ok := loopbackInterface(t)
	if !ok {
		t.Skip("no IPv4 loopback interface")
	}
	_, network, _ := net.ParseCIDR(iface.Address.String() + "/8")
	service := NewService([]Interface{iface}, []*net.IPNet{network}, nil)
	service.Start(context.Background())
	service.Start(context.Background())
	time.Sleep(20 * time.Millisecond)
	service.Stop()
	service.Stop()
	if service.cancel != nil || !sourceOnInterface(iface.Address, iface) {
		t.Fatal("UDP discovery lifecycle or interface scope failed")
	}
}

func rendererPacket(udn, location string) []byte {
	return []byte("HTTP/1.1 200 OK\r\nST: " + MediaRendererDevice + "\r\nUSN: " + udn + "::" + MediaRendererDevice + "\r\nLOCATION: " + location + "\r\n\r\n")
}

func loopbackInterface(t *testing.T) (Interface, bool) {
	t.Helper()
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for _, hardware := range interfaces {
		addresses, _ := hardware.Addrs()
		for _, address := range addresses {
			if network, ok := address.(*net.IPNet); ok && network.IP.To4() != nil && network.IP.IsLoopback() {
				return Interface{Name: hardware.Name, Hardware: hardware, Address: network.IP}, true
			}
		}
	}
	return Interface{}, false
}
