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
	served := []string{
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
	compatibility := []string{
		"http-get:*:*:*",
		"http-get:*:image/*:*",
		"http-get:*:audio:*",
		"http-get:*:audio/*:*",
		"http-get:*:video/*:*",
		"http-get:*:image/jpeg:DLNA.ORG_PN=JPEG_TN",
		"http-get:*:image/jpeg:DLNA.ORG_PN=JPEG_SM",
		"http-get:*:image/jpeg:DLNA.ORG_PN=JPEG_MED",
		"http-get:*:image/jpeg:DLNA.ORG_PN=JPEG_LRG",
		"http-get:*:image/png:DLNA.ORG_PN=PNG_TN",
		"http-get:*:image/png:DLNA.ORG_PN=PNG_LRG",
		"http-get:*:image/gif:DLNA.ORG_PN=GIF_LRG",
		"http-get:*:audio/mpeg:DLNA.ORG_PN=MP3",
		"http-get:*:audio/L16:DLNA.ORG_PN=LPCM",
		"http-get:*:video/mpeg:DLNA.ORG_PN=MPEG_PS_PAL",
		"http-get:*:video/mpeg:DLNA.ORG_PN=MPEG_PS_NTSC",
		"http-get:*:video/mpeg:DLNA.ORG_PN=MPEG_TS_SD_EU_ISO",
		"http-get:*:video/vnd.dlna.mpeg-tts:DLNA.ORG_PN=MPEG_TS_SD_EU",
		"http-get:*:video/vnd.dlna.mpeg-tts:DLNA.ORG_PN=MPEG_TS_SD_EU_T",
		"http-get:*:video/x-matroska:*",
		"http-get:*:video/mp4:*",
		"http-get:*:video/mpeg:*",
		"http-get:*:video/vnd.dlna.mpeg-tts:*",
		"http-get:*:audio/flac:*",
		"http-get:*:audio/x-flac:*",
		"http-get:*:audio/mp3:*",
		"http-get:*:image/png:*",
		"http-get:*:image/jpeg:*",
	}
	return append(served, compatibility...)
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
