package dlna

import (
	"net/http"
	"net/url"
	"strings"

	"media-manager/internal/delivery"
)

type dlnaOutputTarget struct {
	Container    string
	Extension    string
	ContentType  string
	StreamMode   string
	FormatName   string
	DLNAFeatures string
}

func remuxOutputTarget(profile RendererProfile) dlnaOutputTarget {
	return outputTargetFromContainer(profile.DeliveryRules.RemuxContainer, mpegtsOutputTarget())
}

func transcodeOutputTarget(profile RendererProfile) dlnaOutputTarget {
	return outputTargetFromContainer(profile.DeliveryRules.TranscodeContainer, matroskaOutputTarget())
}

func outputTargetFromContainer(value string, fallback dlnaOutputTarget) dlnaOutputTarget {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "mpegts", "ts":
		return mpegtsOutputTarget()
	case "matroska", "mkv":
		return matroskaOutputTarget()
	case "mp4":
		return mp4OutputTarget()
	default:
		return fallback
	}
}

func mpegtsOutputTarget() dlnaOutputTarget {
	return dlnaOutputTarget{
		Container:    "mpegts",
		Extension:    ".ts",
		ContentType:  "video/mp2t",
		StreamMode:   "mpegts_remux",
		FormatName:   "mpegts",
		DLNAFeatures: "DLNA.ORG_OP=01;DLNA.ORG_CI=1",
	}
}

func matroskaOutputTarget() dlnaOutputTarget {
	return dlnaOutputTarget{
		Container:    "matroska",
		Extension:    ".mkv",
		ContentType:  "video/x-matroska",
		StreamMode:   "matroska_transcode",
		FormatName:   "matroska,webm",
		DLNAFeatures: "DLNA.ORG_OP=01;DLNA.ORG_CI=1",
	}
}

func mp4OutputTarget() dlnaOutputTarget {
	return dlnaOutputTarget{
		Container:    "mp4",
		Extension:    ".mp4",
		ContentType:  "video/mp4",
		StreamMode:   "mp4_transcode",
		FormatName:   "mov,mp4,m4a,3gp,3g2,mj2",
		DLNAFeatures: "DLNA.ORG_OP=01;DLNA.ORG_CI=1",
	}
}

func probeForDecision(probe delivery.ProbeResult, decision delivery.Decision, profile RendererProfile) delivery.ProbeResult {
	if decision.Mode == delivery.ModeRemux {
		return probeWithOutputContainer(probe, remuxOutputTarget(profile))
	}
	if decision.Mode == delivery.ModeTranscode && decision.DeliveryProtocol == delivery.ProtocolFile {
		return probeWithOutputContainer(probe, transcodeOutputTarget(profile))
	}
	return probe
}

func probeWithOutputContainer(probe delivery.ProbeResult, target dlnaOutputTarget) delivery.ProbeResult {
	probe.Container.FormatName = &target.FormatName
	return probe
}

func resourceURLForDecision(resourceURL string, decision delivery.Decision) string {
	switch {
	case decision.DeliveryProtocol == delivery.ProtocolHLS:
		return resourceURLWithMode(resourceURL, "hls")
	case decision.Mode == delivery.ModeRemux:
		return resourceURLWithMode(resourceURL, "remux")
	case decision.Mode == delivery.ModeTranscode:
		return resourceURLWithMode(resourceURL, "transcode")
	default:
		return resourceURL
	}
}

func resourceURLWithMode(resourceURL string, mode string) string {
	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return resourceURL
	}
	values := parsed.Query()
	values.Set("mode", mode)
	parsed.RawQuery = values.Encode()
	return parsed.String()
}

func setDLNAOutputHeaders(w http.ResponseWriter, target dlnaOutputTarget) {
	w.Header().Set("Content-Type", target.ContentType)
	w.Header().Set("TransferMode.DLNA.ORG", "Streaming")
	w.Header().Set("ContentFeatures.DLNA.ORG", target.DLNAFeatures)
}
