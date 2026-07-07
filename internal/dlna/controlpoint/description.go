package controlpoint

import (
	"encoding/xml"
	"errors"
	"strings"
)

var ErrNotMediaRenderer = errors.New("device is not a MediaRenderer")

type rootDescription struct {
	Device deviceDescription `xml:"device"`
}

type deviceDescription struct {
	DeviceType   string               `xml:"deviceType"`
	FriendlyName string               `xml:"friendlyName"`
	UDN          string               `xml:"UDN"`
	Services     []serviceDescription `xml:"serviceList>service"`
}

type serviceDescription struct {
	ServiceType string `xml:"serviceType"`
	ControlURL  string `xml:"controlURL"`
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
		UDN:          strings.TrimSpace(root.Device.UDN),
		FriendlyName: strings.TrimSpace(root.Device.FriendlyName),
	}
	for _, service := range root.Device.Services {
		switch strings.TrimSpace(service.ServiceType) {
		case AVTransportService:
			renderer.AVTransportURL = resolveURL(baseURL, strings.TrimSpace(service.ControlURL))
		case RenderingService:
			renderer.RenderingControlURL = resolveURL(baseURL, strings.TrimSpace(service.ControlURL))
		}
	}
	return renderer, nil
}
