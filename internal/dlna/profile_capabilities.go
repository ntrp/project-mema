package dlna

import (
	"strconv"
	"strings"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna/content"
)

type RendererCapabilityDecision struct {
	Decision    delivery.Decision
	ReasonCodes []string
	Trace       []RendererCapabilityTrace
}

type RendererCapabilityTrace struct {
	Field  string
	Value  string
	Rule   string
	Result string
}

func EvaluateRendererCapability(profile RendererProfile, probe delivery.ProbeResult) RendererCapabilityDecision {
	if !profile.DeliveryRules.DirectPlay && !profile.DeliveryRules.Transcode {
		return directCapabilityDecision("profile_rules_empty")
	}
	checks := capabilityChecks(profile, probe)
	reasons := failedCapabilityReasons(checks)
	if len(reasons) == 0 && profile.DeliveryRules.DirectPlay {
		return RendererCapabilityDecision{Decision: directDecision(), Trace: checks}
	}
	if !profile.DeliveryRules.Transcode {
		return RendererCapabilityDecision{
			Decision:    directDecision(),
			ReasonCodes: append(reasons, "transcode_disabled"),
			Trace:       checks,
		}
	}
	return RendererCapabilityDecision{
		Decision:    fallbackCapabilityDecision(profile, reasons),
		ReasonCodes: reasons,
		Trace:       checks,
	}
}

func capabilityChecks(profile RendererProfile, probe delivery.ProbeResult) []RendererCapabilityTrace {
	video := delivery.FirstTrackByType(probe.Tracks, delivery.TrackVideo, nil)
	audio := delivery.FirstTrackByType(probe.Tracks, delivery.TrackAudio, nil)
	return []RendererCapabilityTrace{
		containerCapabilityTrace(profile, probe),
		codecCapabilityTrace("videoCodec", optionalCodec(video), profile.Capabilities.VideoCodecs),
		codecCapabilityTrace("audioCodec", optionalCodec(audio), profile.Capabilities.AudioCodecs),
		resolutionCapabilityTrace(profile, video),
	}
}

func failedCapabilityReasons(checks []RendererCapabilityTrace) []string {
	reasons := []string{}
	for _, check := range checks {
		if check.Result == "pass" {
			continue
		}
		switch check.Field {
		case "container":
			reasons = append(reasons, "container_not_supported")
		case "videoCodec":
			reasons = append(reasons, "video_codec_not_supported")
		case "audioCodec":
			reasons = append(reasons, "audio_codec_not_supported")
		case "resolution":
			reasons = append(reasons, "resolution_not_supported")
		}
	}
	return reasons
}

func fallbackCapabilityDecision(profile RendererProfile, reasons []string) delivery.Decision {
	mode := delivery.ModeTranscode
	protocol := delivery.ProtocolFile
	if containsReason(reasons, "container_not_supported") && !containsVideoReason(reasons) {
		mode = delivery.ModeRemux
	} else if containsVideoReason(reasons) && !profile.AvoidHLS {
		protocol = delivery.ProtocolHLS
	}
	return delivery.Decision{
		DeliveryProtocol: protocol,
		Mode:             mode,
		Plan: delivery.TranscodePlan{
			VideoCodec: videoPlanCodec(reasons),
			AudioCodec: audioPlanCodec(reasons),
		},
		Reasons: append([]string{}, reasons...),
	}
}

func directCapabilityDecision(reason string) RendererCapabilityDecision {
	return RendererCapabilityDecision{
		Decision:    directDecision(),
		ReasonCodes: []string{reason},
	}
}

func containerCapabilityTrace(profile RendererProfile, probe delivery.ProbeResult) RendererCapabilityTrace {
	value := normalizedContainer(probe.Container.FormatName)
	return RendererCapabilityTrace{
		Field:  "container",
		Value:  value,
		Rule:   strings.Join(profile.Capabilities.Containers, ","),
		Result: passFail(value == "" || containsAnyToken(profile.Capabilities.Containers, value)),
	}
}

func codecCapabilityTrace(field string, codec string, allowed []string) RendererCapabilityTrace {
	return RendererCapabilityTrace{
		Field:  field,
		Value:  codec,
		Rule:   strings.Join(allowed, ","),
		Result: passFail(codec == "" || containsCodec(allowed, codec)),
	}
}

func resolutionCapabilityTrace(profile RendererProfile, video *delivery.Track) RendererCapabilityTrace {
	height := int32(0)
	if video != nil && video.Height != nil {
		height = *video.Height
	}
	maxHeight := maxProfileHeight(profile.Capabilities.MaxResolution)
	return RendererCapabilityTrace{
		Field:  "resolution",
		Value:  strconv.Itoa(int(height)),
		Rule:   profile.Capabilities.MaxResolution,
		Result: passFail(height == 0 || maxHeight == 0 || height <= maxHeight),
	}
}

func normalizedContainer(value *string) string {
	if value == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(*value))
}

func optionalCodec(track *delivery.Track) string {
	if track == nil || track.Codec == nil {
		return ""
	}
	return normalizeCodec(*track.Codec)
}

func normalizeCodec(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "avc1":
		return "h264"
	case "h265":
		return "hevc"
	case "e-ac-3":
		return "eac3"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func containsAnyToken(allowed []string, value string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, item := range allowed {
		if strings.Contains(value, item) {
			return true
		}
	}
	return false
}

func containsCodec(allowed []string, codec string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, item := range allowed {
		if normalizeCodec(item) == codec {
			return true
		}
	}
	return false
}

func maxProfileHeight(value string) int32 {
	value = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(value)), "p")
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return int32(parsed)
}

func passFail(ok bool) string {
	if ok {
		return "pass"
	}
	return "fail"
}

func containsReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if reason == want {
			return true
		}
	}
	return false
}

func containsVideoReason(reasons []string) bool {
	return containsReason(reasons, "video_codec_not_supported") ||
		containsReason(reasons, "resolution_not_supported")
}

func videoPlanCodec(reasons []string) string {
	if containsVideoReason(reasons) {
		return "libx264"
	}
	return "copy"
}

func audioPlanCodec(reasons []string) string {
	if containsReason(reasons, "audio_codec_not_supported") {
		return "aac"
	}
	return "copy"
}

func SourceProtocolInfosForCapabilities(profile RendererProfile) []string {
	if len(profile.Capabilities.Containers) == 0 {
		return SourceProtocolInfos()
	}
	values := []string{}
	if profile.DeliveryRules.DirectPlay {
		for _, container := range profile.Capabilities.Containers {
			values = append(values, protocolInfoForContainer(container, delivery.ModeDirect))
		}
	}
	if profile.DeliveryRules.Transcode {
		values = append(values, protocolInfoForContainer("mpegts", delivery.ModeRemux))
		if !profile.AvoidHLS {
			values = append(values, content.ProtocolInfo("stream.m3u8", delivery.Container{}, delivery.Decision{
				DeliveryProtocol: delivery.ProtocolHLS,
				Mode:             delivery.ModeTranscode,
			}))
		}
	}
	if len(values) == 0 {
		return SourceProtocolInfos()
	}
	return values
}

func protocolInfoForContainer(container string, mode delivery.Mode) string {
	format := containerFormatForProtocol(container)
	return content.ProtocolInfo("direct."+containerExtension(container), delivery.Container{FormatName: &format}, delivery.Decision{
		DeliveryProtocol: delivery.ProtocolFile,
		Mode:             mode,
	})
}

func containerFormatForProtocol(container string) string {
	switch strings.ToLower(strings.TrimSpace(container)) {
	case "mp4":
		return "mov,mp4,m4a,3gp,3g2,mj2"
	case "mkv":
		return "matroska,webm"
	case "mpegts":
		return "mpegts"
	default:
		return container
	}
}

func containerExtension(container string) string {
	switch strings.ToLower(strings.TrimSpace(container)) {
	case "mkv":
		return "mkv"
	case "mpegts":
		return "ts"
	default:
		return strings.TrimSpace(container)
	}
}
