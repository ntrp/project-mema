package controlpoint

import "net/url"

const (
	MediaRendererDevice = "urn:schemas-upnp-org:device:MediaRenderer:1"
	AVTransportService  = "urn:schemas-upnp-org:service:AVTransport:1"
	RenderingService    = "urn:schemas-upnp-org:service:RenderingControl:1"
)

type Renderer struct {
	UDN                 string
	FriendlyName        string
	AVTransportURL      string
	RenderingControlURL string
}

func resolveURL(base string, value string) string {
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
