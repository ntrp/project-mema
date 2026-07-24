package controlpoint

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	MediaRendererDevice = "urn:schemas-upnp-org:device:MediaRenderer:1"
	AVTransportService  = "urn:schemas-upnp-org:service:AVTransport:1"
	RenderingService    = "urn:schemas-upnp-org:service:RenderingControl:1"
)

var ErrNotMediaRenderer = errors.New("device is not a MediaRenderer")

type Renderer struct {
	UDN, FriendlyName, Manufacturer, ModelName           string
	AVTransportURL, RenderingControlURL                  string
	AVTransportSCPDURL, RenderingSCPDURL, DescriptionURL string
	AVTransportActions, RenderingActions                 []string
	AVTransportActionsKnown, RenderingActionsKnown       bool
	SourceIP, LocalIP                                    net.IP
	LastSeen, ExpiresAt                                  time.Time
}

type rootDescription struct {
	Device deviceDescription `xml:"device"`
}
type deviceDescription struct {
	DeviceType   string               `xml:"deviceType"`
	FriendlyName string               `xml:"friendlyName"`
	Manufacturer string               `xml:"manufacturer"`
	ModelName    string               `xml:"modelName"`
	UDN          string               `xml:"UDN"`
	Services     []serviceDescription `xml:"serviceList>service"`
}
type serviceDescription struct {
	ServiceType string `xml:"serviceType"`
	ControlURL  string `xml:"controlURL"`
	SCPDURL     string `xml:"SCPDURL"`
}

func ParseRendererDescription(baseURL string, payload []byte) (Renderer, error) {
	var root rootDescription
	if err := xml.Unmarshal(payload, &root); err != nil {
		return Renderer{}, err
	}
	if strings.TrimSpace(root.Device.DeviceType) != MediaRendererDevice {
		return Renderer{}, ErrNotMediaRenderer
	}
	renderer := Renderer{
		UDN: strings.TrimSpace(root.Device.UDN), FriendlyName: strings.TrimSpace(root.Device.FriendlyName),
		Manufacturer: strings.TrimSpace(root.Device.Manufacturer), ModelName: strings.TrimSpace(root.Device.ModelName),
	}
	for _, service := range root.Device.Services {
		controlURL, scpdURL := strings.TrimSpace(service.ControlURL), strings.TrimSpace(service.SCPDURL)
		switch strings.TrimSpace(service.ServiceType) {
		case AVTransportService:
			renderer.AVTransportURL = resolveURL(baseURL, controlURL)
			if scpdURL != "" {
				renderer.AVTransportSCPDURL = resolveURL(baseURL, scpdURL)
			}
		case RenderingService:
			renderer.RenderingControlURL = resolveURL(baseURL, controlURL)
			if scpdURL != "" {
				renderer.RenderingSCPDURL = resolveURL(baseURL, scpdURL)
			}
		}
	}
	if renderer.UDN == "" || renderer.AVTransportURL == "" {
		return Renderer{}, errors.New("renderer description is missing identity or AVTransport service")
	}
	return renderer, nil
}

type serviceDescriptionDocument struct {
	Actions []struct {
		Name string `xml:"name"`
	} `xml:"actionList>action"`
}

func ParseServiceActions(payload []byte) ([]string, error) {
	var document serviceDescriptionDocument
	if err := xml.Unmarshal(payload, &document); err != nil {
		return nil, err
	}
	actions, seen := make([]string, 0, len(document.Actions)), map[string]bool{}
	for _, action := range document.Actions {
		if name := strings.TrimSpace(action.Name); name != "" && !seen[name] {
			actions, seen[name] = append(actions, name), true
		}
	}
	sort.Strings(actions)
	return actions, nil
}

func SupportsTransport(renderer Renderer, action string) bool {
	return !renderer.AVTransportActionsKnown || containsAction(renderer.AVTransportActions, action)
}
func SupportsRendering(renderer Renderer, action string) bool {
	return !renderer.RenderingActionsKnown || containsAction(renderer.RenderingActions, action)
}
func containsAction(actions []string, action string) bool {
	for _, supported := range actions {
		if strings.EqualFold(supported, action) {
			return true
		}
	}
	return false
}

func resolveURL(base, value string) string {
	parsed, err := url.Parse(value)
	if err == nil && parsed.IsAbs() {
		return parsed.String()
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return value
	}
	relative, err := url.Parse(value)
	if err != nil {
		return value
	}
	return baseURL.ResolveReference(relative).String()
}

// ValidateEndpoint permits only literal unicast IP hosts inside configured LANs.
func ValidateEndpoint(raw string, networks []*net.IPNet) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("invalid renderer endpoint")
	}
	if u.User != nil || u.Fragment != "" {
		return fmt.Errorf("renderer endpoint contains unsupported credentials or fragment")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || ip.IsUnspecified() || ip.IsMulticast() {
		return fmt.Errorf("renderer endpoint host must be a unicast IP address")
	}
	for _, network := range networks {
		if network != nil && network.Contains(ip) {
			return nil
		}
	}
	return fmt.Errorf("renderer endpoint is outside configured LAN")
}

func validateRenderer(renderer Renderer, source net.IP, networks []*net.IPNet) error {
	for _, endpoint := range []string{renderer.DescriptionURL, renderer.AVTransportURL, renderer.RenderingControlURL, renderer.AVTransportSCPDURL, renderer.RenderingSCPDURL} {
		if endpoint == "" {
			continue
		}
		if err := validateEndpointSource(endpoint, source, networks); err != nil {
			return err
		}
	}
	return nil
}

func validateEndpointSource(endpoint string, source net.IP, networks []*net.IPNet) error {
	if err := ValidateEndpoint(endpoint, networks); err != nil {
		return err
	}
	u, _ := url.Parse(endpoint)
	if source == nil || !net.ParseIP(u.Hostname()).Equal(source) {
		return fmt.Errorf("renderer endpoint host does not match SSDP sender")
	}
	return nil
}
