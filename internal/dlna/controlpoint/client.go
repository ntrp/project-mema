package controlpoint

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	http HTTPDoer
}

func NewClient(httpClient HTTPDoer) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return Client{http: httpClient}
}

func (c Client) SetAVTransportURI(ctx context.Context, renderer Renderer, uri string, metadata string) error {
	return c.call(ctx, renderer.AVTransportURL, AVTransportService, "SetAVTransportURI", map[string]string{
		"InstanceID":         "0",
		"CurrentURI":         uri,
		"CurrentURIMetaData": metadata,
	})
}

func (c Client) Play(ctx context.Context, renderer Renderer) error {
	return c.transport(ctx, renderer, "Play", map[string]string{"Speed": "1"})
}

func (c Client) Pause(ctx context.Context, renderer Renderer) error {
	return c.transport(ctx, renderer, "Pause", nil)
}

func (c Client) Stop(ctx context.Context, renderer Renderer) error {
	return c.transport(ctx, renderer, "Stop", nil)
}

func (c Client) Seek(ctx context.Context, renderer Renderer, target string) error {
	return c.transport(ctx, renderer, "Seek", map[string]string{"Unit": "REL_TIME", "Target": target})
}

func (c Client) Next(ctx context.Context, renderer Renderer) error {
	return c.transport(ctx, renderer, "Next", nil)
}

func (c Client) Previous(ctx context.Context, renderer Renderer) error {
	return c.transport(ctx, renderer, "Previous", nil)
}

func (c Client) SetVolume(ctx context.Context, renderer Renderer, volume int) error {
	return c.rendering(ctx, renderer, "SetVolume", map[string]string{"Channel": "Master", "DesiredVolume": strconv.Itoa(volume)})
}

func (c Client) SetMute(ctx context.Context, renderer Renderer, mute bool) error {
	value := "0"
	if mute {
		value = "1"
	}
	return c.rendering(ctx, renderer, "SetMute", map[string]string{"Channel": "Master", "DesiredMute": value})
}

func (c Client) transport(ctx context.Context, renderer Renderer, action string, args map[string]string) error {
	values := map[string]string{"InstanceID": "0"}
	for key, value := range args {
		values[key] = value
	}
	return c.call(ctx, renderer.AVTransportURL, AVTransportService, action, values)
}

func (c Client) rendering(ctx context.Context, renderer Renderer, action string, args map[string]string) error {
	values := map[string]string{"InstanceID": "0"}
	for key, value := range args {
		values[key] = value
	}
	return c.call(ctx, renderer.RenderingControlURL, RenderingService, action, values)
}

func (c Client) call(ctx context.Context, endpoint string, service string, action string, args map[string]string) error {
	payload, err := envelope(service, action, args)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	request.Header.Set("SOAPACTION", `"`+service+"#"+action+`"`)
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return fmt.Errorf("renderer action %s failed: status %d: %s", action, response.StatusCode, string(body))
	}
	return nil
}

func envelope(service string, action string, args map[string]string) ([]byte, error) {
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
