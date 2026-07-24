package controlpoint

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	http     HTTPDoer
	networks []*net.IPNet
}

type FaultError struct {
	Action      string
	StatusCode  int
	UPnPCode    string
	Description string
}

func (e *FaultError) Error() string {
	detail := strings.TrimSpace(e.Description)
	if detail == "" {
		detail = fmt.Sprintf("HTTP status %d", e.StatusCode)
	}
	return fmt.Sprintf("renderer rejected %s: %s", e.Action, detail)
}

func NewClient(httpClient HTTPDoer) Client { return NewSecureClient(httpClient, nil) }

func NewSecureClient(httpClient HTTPDoer, networks []*net.IPNet) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return Client{http: httpClient, networks: append([]*net.IPNet{}, networks...)}
}

func (c Client) SetAVTransportURI(ctx context.Context, renderer Renderer, uri, metadata string) error {
	return c.call(ctx, renderer.AVTransportURL, AVTransportService, "SetAVTransportURI", map[string]string{
		"InstanceID": "0", "CurrentURI": uri, "CurrentURIMetaData": metadata,
	})
}
func (c Client) Play(ctx context.Context, r Renderer) error {
	return c.transport(ctx, r, "Play", map[string]string{"Speed": "1"})
}
func (c Client) Pause(ctx context.Context, r Renderer) error {
	return c.transport(ctx, r, "Pause", nil)
}
func (c Client) Stop(ctx context.Context, r Renderer) error { return c.transport(ctx, r, "Stop", nil) }
func (c Client) Next(ctx context.Context, r Renderer) error { return c.transport(ctx, r, "Next", nil) }
func (c Client) Previous(ctx context.Context, r Renderer) error {
	return c.transport(ctx, r, "Previous", nil)
}
func (c Client) Seek(ctx context.Context, r Renderer, target string) error {
	return c.transport(ctx, r, "Seek", map[string]string{"Unit": "REL_TIME", "Target": target})
}
func (c Client) SetVolume(ctx context.Context, r Renderer, volume int) error {
	return c.rendering(ctx, r, "SetVolume", map[string]string{"Channel": "Master", "DesiredVolume": strconv.Itoa(volume)})
}
func (c Client) SetMute(ctx context.Context, r Renderer, mute bool) error {
	value := "0"
	if mute {
		value = "1"
	}
	return c.rendering(ctx, r, "SetMute", map[string]string{"Channel": "Master", "DesiredMute": value})
}
func (c Client) transport(ctx context.Context, r Renderer, action string, args map[string]string) error {
	return c.withInstance(ctx, r.AVTransportURL, AVTransportService, action, args)
}
func (c Client) rendering(ctx context.Context, r Renderer, action string, args map[string]string) error {
	return c.withInstance(ctx, r.RenderingControlURL, RenderingService, action, args)
}
func (c Client) withInstance(ctx context.Context, endpoint, service, action string, args map[string]string) error {
	values := map[string]string{"InstanceID": "0"}
	for key, value := range args {
		values[key] = value
	}
	return c.call(ctx, endpoint, service, action, values)
}
func (c Client) call(ctx context.Context, endpoint, service, action string, args map[string]string) error {
	response, err := c.request(ctx, endpoint, service, action, args)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 400 {
		return nil
	}
	return parseFault(action, response)
}
func (c Client) request(ctx context.Context, endpoint, service, action string, args map[string]string) (*http.Response, error) {
	if len(c.networks) > 0 {
		if err := ValidateEndpoint(endpoint, c.networks); err != nil {
			return nil, err
		}
	}
	payload, err := envelope(service, action, args)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	request.Header.Set("SOAPACTION", `"`+service+"#"+action+`"`)
	return c.http.Do(request)
}
func parseFault(action string, response *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(response.Body, 16<<10))
	var fault struct {
		Code        string `xml:"Body>Fault>detail>UPnPError>errorCode"`
		Description string `xml:"Body>Fault>detail>UPnPError>errorDescription"`
	}
	_ = xml.Unmarshal(body, &fault)
	return &FaultError{Action: action, StatusCode: response.StatusCode, UPnPCode: fault.Code, Description: fault.Description}
}

func envelope(service, action string, args map[string]string) ([]byte, error) {
	var body bytes.Buffer
	body.WriteString(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>`)
	body.WriteString(`<u:` + action + ` xmlns:u="` + service + `">`)
	for key, value := range args {
		body.WriteString(`<` + key + `>`)
		if err := xml.EscapeText(&body, []byte(value)); err != nil {
			return nil, err
		}
		body.WriteString(`</` + key + `>`)
	}
	body.WriteString(`</u:` + action + `></s:Body></s:Envelope>`)
	return body.Bytes(), nil
}
