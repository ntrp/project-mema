package content

import (
	"strings"
	"testing"

	"media-manager/internal/delivery"
)

func TestRenderDIDLSerializesContainerNamespacesAndEscapesXML(t *testing.T) {
	payload, err := RenderDIDL([]Object{{
		ID:         "container-1",
		ParentID:   RootID,
		Title:      `Rock & Roll <Movies>`,
		Class:      "object.container.storageFolder",
		Kind:       ObjectContainer,
		ChildCount: 2,
	}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := string(payload)
	for _, want := range []string{
		`xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"`,
		`xmlns:dc="http://purl.org/dc/elements/1.1/"`,
		`xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/"`,
		`xmlns:dlna="urn:schemas-dlna-org:metadata-1-0/"`,
		`<dc:title>Rock &amp; Roll &lt;Movies&gt;</dc:title>`,
		`childCount="2"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("DIDL missing %q:\n%s", want, got)
		}
	}
}

func TestRenderDIDLSerializesItemMetadataAndResources(t *testing.T) {
	date := "2026-07-07"
	album := "Scenario Collection"
	artwork := "http://127.0.0.1:18080/dlna/artwork/poster.jpg?size=large"
	size := int64(4096)
	duration := "0:01:05.250"
	bitrate := int64(800000)
	channels := int32(2)
	resolution := "1920x1080"
	object := Object{
		ID:       "item-1",
		ParentID: "container-1",
		Title:    "Scenario & Movie",
		Class:    "object.item.videoItem.movie",
		Kind:     ObjectItem,
		Date:     &date,
		Genres:   []string{"Drama"},
		Artists:  []string{"Actor One"},
		Album:    &album,
		Artwork:  &artwork,
	}
	resource := Resource{
		URL:           "http://127.0.0.1:18080/dlna/resource/item-1?mode=direct&x=1",
		ProtocolInfo:  "http-get:*:video/mp4:DLNA.ORG_OP=01;DLNA.ORG_CI=0",
		SizeBytes:     &size,
		Duration:      &duration,
		BitRate:       &bitrate,
		Resolution:    &resolution,
		AudioChannels: &channels,
	}

	payload, err := RenderDIDL([]Object{object}, map[string][]Resource{object.ID: []Resource{resource}})
	if err != nil {
		t.Fatal(err)
	}
	got := string(payload)
	for _, want := range []string{
		`<dc:title>Scenario &amp; Movie</dc:title>`,
		`<dc:date>2026-07-07</dc:date>`,
		`<upnp:genre>Drama</upnp:genre>`,
		`<upnp:artist>Actor One</upnp:artist>`,
		`<upnp:album>Scenario Collection</upnp:album>`,
		`<upnp:albumArtURI>http://127.0.0.1:18080/dlna/artwork/poster.jpg?size=large</upnp:albumArtURI>`,
		`protocolInfo="http-get:*:video/mp4:DLNA.ORG_OP=01;DLNA.ORG_CI=0"`,
		`size="4096"`,
		`duration="0:01:05.250"`,
		`bitrate="800000"`,
		`resolution="1920x1080"`,
		`nrAudioChannels="2"`,
		`http://127.0.0.1:18080/dlna/resource/item-1?mode=direct&amp;x=1`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("DIDL missing %q:\n%s", want, got)
		}
	}
}

func stringPtr(value string) *string {
	return &value
}

func TestResourceFromDeliveryMapsProbeAttributes(t *testing.T) {
	width := int32(1920)
	height := int32(1080)
	channels := int32(6)
	videoBitrate := "700000"
	audioBitrate := "300000"
	duration := 65.25
	size := int64(12345)

	resource := ResourceFromDelivery(ResourceInput{
		URL:       "http://mema.local/video.mp4",
		SizeBytes: &size,
		Probe: delivery.ProbeResult{
			Container:       delivery.Container{FormatName: stringPtr("mov,mp4,m4a,3gp,3g2,mj2")},
			DurationSeconds: &duration,
			Tracks: []delivery.Track{
				{Type: delivery.TrackVideo, Width: &width, Height: &height, BitRate: &videoBitrate},
				{Type: delivery.TrackAudio, Channels: &channels, BitRate: &audioBitrate},
			},
		},
		Decision: delivery.Decision{DeliveryProtocol: delivery.ProtocolFile, Mode: delivery.ModeDirect},
	})

	if resource.ProtocolInfo != "http-get:*:video/mp4:DLNA.ORG_OP=01;DLNA.ORG_CI=0" {
		t.Fatalf("protocolInfo = %q", resource.ProtocolInfo)
	}
	if resource.SizeBytes == nil || *resource.SizeBytes != size {
		t.Fatalf("size = %#v", resource.SizeBytes)
	}
	if resource.Duration == nil || *resource.Duration != "0:01:05.250" {
		t.Fatalf("duration = %#v", resource.Duration)
	}
	if resource.BitRate == nil || *resource.BitRate != 1000000 {
		t.Fatalf("bitrate = %#v", resource.BitRate)
	}
	if resource.Resolution == nil || *resource.Resolution != "1920x1080" {
		t.Fatalf("resolution = %#v", resource.Resolution)
	}
	if resource.AudioChannels == nil || *resource.AudioChannels != 6 {
		t.Fatalf("channels = %#v", resource.AudioChannels)
	}
}

func TestResourceFromDeliveryMapsHLSAlternate(t *testing.T) {
	size := int64(12345)

	resource := ResourceFromDelivery(ResourceInput{
		URL:       "http://mema.local/video.m3u8",
		SizeBytes: &size,
		Decision:  delivery.Decision{DeliveryProtocol: delivery.ProtocolHLS, Mode: delivery.ModeTranscode},
	})

	if resource.ProtocolInfo != "http-get:*:application/vnd.apple.mpegurl:DLNA.ORG_OP=01;DLNA.ORG_CI=1" {
		t.Fatalf("protocolInfo = %q", resource.ProtocolInfo)
	}
	if resource.SizeBytes != nil {
		t.Fatalf("HLS size = %#v", resource.SizeBytes)
	}
}
