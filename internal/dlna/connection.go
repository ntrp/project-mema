package dlna

import (
	"strings"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna/content"
)

func SourceProtocolInfo() string {
	return strings.Join(SourceProtocolInfos(), ",")
}

func SourceProtocolInfos() []string {
	mp4 := "mov,mp4,m4a,3gp,3g2,mj2"
	mkv := "matroska,webm"
	mpegts := "mpegts"
	return []string{
		content.ProtocolInfo("direct.mp4", delivery.Container{FormatName: &mp4}, delivery.Decision{
			DeliveryProtocol: delivery.ProtocolFile,
			Mode:             delivery.ModeDirect,
		}),
		content.ProtocolInfo("direct.mkv", delivery.Container{FormatName: &mkv}, delivery.Decision{
			DeliveryProtocol: delivery.ProtocolFile,
			Mode:             delivery.ModeDirect,
		}),
		content.ProtocolInfo("remux.ts", delivery.Container{FormatName: &mpegts}, delivery.Decision{
			DeliveryProtocol: delivery.ProtocolFile,
			Mode:             delivery.ModeRemux,
		}),
		content.ProtocolInfo("stream.m3u8", delivery.Container{}, delivery.Decision{
			DeliveryProtocol: delivery.ProtocolHLS,
			Mode:             delivery.ModeTranscode,
		}),
	}
}

func CurrentConnectionInfo(connectionID string) (map[string]string, error) {
	if strings.TrimSpace(connectionID) != "0" {
		return nil, errNoSuchConnection()
	}
	return map[string]string{
		"RcsID":                 "-1",
		"AVTransportID":         "-1",
		"ProtocolInfo":          SourceProtocolInfos()[0],
		"PeerConnectionManager": "",
		"PeerConnectionID":      "-1",
		"Direction":             "Output",
		"Status":                "OK",
	}, nil
}

func errNoSuchConnection() error {
	return connectionError{code: 706, description: "No Such Connection"}
}

type connectionError struct {
	code        int
	description string
}

func (e connectionError) Error() string {
	return e.description
}
