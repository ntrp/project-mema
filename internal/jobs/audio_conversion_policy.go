package jobs

import (
	"fmt"
	"strconv"
	"strings"

	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

type AudioConversionInput struct {
	Policy               string
	SourceCodec          string
	SourceChannels       string
	SourceBitrateKbps    int32
	TargetCodec          *string
	TargetChannels       []string
	MinimumBitrateKbps   *int32
	PreferredBitrateKbps *int32
}

type AudioConversionDecision struct {
	Needed            bool
	Allowed           bool
	Status            string
	Reason            string
	Policy            string
	SourceCodec       string
	TargetCodec       string
	TargetChannels    string
	TargetBitrateKbps int32
}

func DecideAudioConversion(input AudioConversionInput) AudioConversionDecision {
	decision := audioConversionNeed(input)
	if !decision.Needed {
		decision.Status = "satisfied"
		decision.Reason = "Audio already meets target details."
		return decision
	}
	switch input.Policy {
	case "manual":
		decision.Allowed = true
		decision.Status = "allowed"
		decision.Reason = "Manual audio conversion requested."
		return decision
	case "losslessToLossy":
		if audioCodecIsLossless(input.SourceCodec) {
			decision.Allowed = true
			decision.Status = "allowed"
			decision.Reason = "Profile allows conversion from lossless audio."
			return decision
		}
		decision.Status = "blocked"
		decision.Reason = "Profile allows conversion only from lossless audio."
	case "lossyToLossy":
		decision.Allowed = true
		decision.Status = "allowed"
		decision.Reason = "Profile allows conversion from lossy audio."
	default:
		decision.Status = "blocked"
		decision.Reason = "Audio conversion disabled by profile."
	}
	return decision
}

func audioConversionHasExecutableWork(decision AudioConversionDecision) bool {
	if !decision.Needed || !decision.Allowed || decision.TargetCodec == "" {
		return false
	}
	if decision.SourceCodec == "" || decision.SourceCodec != decision.TargetCodec {
		return true
	}
	if ffmpegChannelCount(decision.TargetChannels) != "" {
		return true
	}
	return decision.TargetBitrateKbps > 0
}

func audioConversionNeed(input AudioConversionInput) AudioConversionDecision {
	decision := AudioConversionDecision{
		Policy:      input.Policy,
		SourceCodec: normalizeJobAudioCodec(input.SourceCodec),
	}
	if input.TargetCodec != nil {
		decision.TargetCodec = normalizeJobAudioCodec(*input.TargetCodec)
		if decision.SourceCodec != "" && decision.SourceCodec != decision.TargetCodec {
			decision.Needed = true
		}
	}
	decision.TargetChannels = firstMissingChannelTarget(input.SourceChannels, input.TargetChannels)
	if decision.TargetChannels != "" {
		decision.Needed = true
	}
	decision.TargetBitrateKbps = desiredAudioBitrate(input)
	if decision.TargetBitrateKbps > 0 && input.SourceBitrateKbps > 0 &&
		input.SourceBitrateKbps < decision.TargetBitrateKbps {
		decision.Needed = true
	}
	return decision
}

func AudioConversionProvenance(
	artifact storage.MediaComponentArtifact,
	decision AudioConversionDecision,
) map[string]any {
	return map[string]any{
		"kind":              "audioConversion",
		"sourceId":          artifact.SourceID.String(),
		"artifactId":        artifact.ID.String(),
		"streamId":          artifact.StreamID,
		"streamType":        artifact.StreamType,
		"policy":            decision.Policy,
		"decisionStatus":    decision.Status,
		"reason":            decision.Reason,
		"sourceCodec":       decision.SourceCodec,
		"targetCodec":       decision.TargetCodec,
		"targetChannels":    decision.TargetChannels,
		"targetBitrateKbps": decision.TargetBitrateKbps,
	}
}

func FfmpegAudioConversionArgs(
	inputPath string,
	outputPath string,
	decision AudioConversionDecision,
) ([]string, error) {
	if !decision.Allowed {
		return nil, fmt.Errorf("audio conversion is not allowed")
	}
	if decision.TargetCodec == "" {
		return nil, fmt.Errorf("target audio codec is required")
	}
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	args := []string{"-y", "-i", inputPath, "-map", "0:a:0", "-c:a", ffmpegAudioCodec(decision.TargetCodec)}
	if channels := ffmpegChannelCount(decision.TargetChannels); channels != "" {
		args = append(args, "-ac", channels)
	}
	if decision.TargetBitrateKbps > 0 {
		args = append(args, "-b:a", strconv.Itoa(int(decision.TargetBitrateKbps))+"k")
	}
	return append(args, outputPath), nil
}

func firstMissingChannelTarget(source string, targets []string) string {
	source = storage.NormalizeAudioChannelDefinition(source)
	if source == "" || len(targets) == 0 {
		return ""
	}
	for _, target := range targets {
		normalized := storage.NormalizeAudioChannelDefinition(target)
		if source == normalized {
			return ""
		}
	}
	for _, target := range targets {
		if normalized := storage.NormalizeAudioChannelDefinition(target); normalized != "" {
			return normalized
		}
	}
	return ""
}

func desiredAudioBitrate(input AudioConversionInput) int32 {
	if input.PreferredBitrateKbps != nil {
		return *input.PreferredBitrateKbps
	}
	if input.MinimumBitrateKbps != nil {
		return *input.MinimumBitrateKbps
	}
	return 0
}

func audioCodecIsLossless(value string) bool {
	switch normalizeJobAudioCodec(value) {
	case "flac", "truehd", "pcm":
		return true
	default:
		return false
	}
}

func normalizeJobAudioCodec(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.EqualFold(trimmed, "DD+") {
		return "eac3"
	}
	switch normalizedJobToken(trimmed) {
	case "ddp", "ddplus", "eac3":
		return "eac3"
	case "dd", "ac3", "dolbydigital":
		return "ac3"
	case "truehd", "truehdatmos":
		return "truehd"
	default:
		return strings.ToLower(trimmed)
	}
}

func ffmpegAudioCodec(value string) string {
	switch normalizeJobAudioCodec(value) {
	case "opus":
		return "libopus"
	default:
		return normalizeJobAudioCodec(value)
	}
}

func ffmpegChannelCount(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1.0":
		return "1"
	case "2.0":
		return "2"
	case "5.1":
		return "6"
	case "7.1":
		return "8"
	default:
		return ""
	}
}

func normalizedJobToken(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
