package controlpoint

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TransportState struct {
	State  string
	Status string
	Speed  string
}
type PositionInfo struct {
	Track         int
	RelTime       string
	TrackDuration string
}
type VolumeState struct {
	Volume int
	Mute   bool
}
type RendererState struct {
	Transport          TransportState
	Position           PositionInfo
	Volume             VolumeState
	TransportAvailable bool
	PositionAvailable  bool
	VolumeAvailable    bool
	MuteAvailable      bool
}

func (c Client) GetTransportInfo(ctx context.Context, renderer Renderer) (TransportState, error) {
	values, err := c.read(ctx, renderer, renderer.AVTransportURL, AVTransportService, "GetTransportInfo", map[string]string{"InstanceID": "0"})
	return TransportState{State: values["CurrentTransportState"], Status: values["CurrentTransportStatus"], Speed: values["CurrentSpeed"]}, err
}
func (c Client) GetPositionInfo(ctx context.Context, renderer Renderer) (PositionInfo, error) {
	values, err := c.read(ctx, renderer, renderer.AVTransportURL, AVTransportService, "GetPositionInfo", map[string]string{"InstanceID": "0"})
	track, _ := strconv.Atoi(values["Track"])
	return PositionInfo{Track: track, RelTime: values["RelTime"], TrackDuration: values["TrackDuration"]}, err
}
func (c Client) GetVolume(ctx context.Context, renderer Renderer) (int, error) {
	values, err := c.read(ctx, renderer, renderer.RenderingControlURL, RenderingService, "GetVolume", map[string]string{"InstanceID": "0", "Channel": "Master"})
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(values["CurrentVolume"])
}
func (c Client) GetMute(ctx context.Context, renderer Renderer) (bool, error) {
	values, err := c.read(ctx, renderer, renderer.RenderingControlURL, RenderingService, "GetMute", map[string]string{"InstanceID": "0", "Channel": "Master"})
	if err != nil {
		return false, err
	}
	value := strings.ToLower(values["CurrentMute"])
	return value == "1" || value == "true", nil
}
func (c Client) read(ctx context.Context, renderer Renderer, endpoint, service, action string, args map[string]string) (map[string]string, error) {
	response, err := c.request(ctx, renderer, endpoint, service, action, args)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil, parseFault(action, response)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 64<<10))
	if err != nil {
		return nil, err
	}
	var root xmlNode
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	values := map[string]string{}
	root.flatten(values)
	if len(values) == 0 {
		return nil, fmt.Errorf("renderer returned an empty %s response", action)
	}
	return values, nil
}

type xmlNode struct {
	XMLName  xml.Name
	Text     string    `xml:",chardata"`
	Children []xmlNode `xml:",any"`
}

func (n xmlNode) flatten(values map[string]string) {
	if value := strings.TrimSpace(n.Text); value != "" {
		values[n.XMLName.Local] = value
	}
	for _, child := range n.Children {
		child.flatten(values)
	}
}
