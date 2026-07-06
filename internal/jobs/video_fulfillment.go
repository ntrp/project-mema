package jobs

import (
	"encoding/json"
	"fmt"
	"strings"

	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

type VideoFulfillmentPlan struct {
	Status      string
	Reason      string
	SourcePath  string
	TargetCodec string
	TargetHDR   string
	TargetPixel string
	Provenance  map[string]any
}

type VideoStreamInfo struct {
	Codec       string
	HDRFormat   string
	PixelFormat string
}

func PlanVideoFulfillment(
	item storage.MediaItem,
	profile *storage.MediaProfile,
) VideoFulfillmentPlan {
	if profile == nil {
		return VideoFulfillmentPlan{Status: "skipped", Reason: "No media profile is assigned."}
	}
	source := firstBaseVideoSource(item)
	if source.RetainedPath == "" {
		return VideoFulfillmentPlan{Status: "blocked", Reason: "No retained base video source is available."}
	}
	stream := videoStreamFromInventory(source.StreamInventory)
	target := profile.VideoTarget
	plan := VideoFulfillmentPlan{
		Status:      "satisfied",
		Reason:      "Video already satisfies target.",
		SourcePath:  source.RetainedPath,
		TargetCodec: firstString(target.Codecs),
		TargetHDR:   firstString(target.HDRFormats),
		TargetPixel: firstString(target.PixelFormats),
		Provenance:  videoFulfillmentProvenance(source, "satisfied"),
	}
	missing := videoTargetMisses(stream, target)
	if len(missing) == 0 {
		return plan
	}
	plan.Status = "transcodeRequired"
	plan.Reason = "Video target mismatch: " + strings.Join(missing, ", ")
	plan.Provenance = videoFulfillmentProvenance(source, plan.Status)
	return plan
}

func FfprobeVideoArgs(path string) ([]string, error) {
	if err := mediatools.SafePathArg(path); err != nil {
		return nil, err
	}
	return []string{"-v", "error", "-select_streams", "v:0", "-show_streams", "-of", "json", path}, nil
}

func FfmpegVideoTranscodeArgs(inputPath string, outputPath string, plan VideoFulfillmentPlan) ([]string, error) {
	if plan.Status != "transcodeRequired" {
		return nil, fmt.Errorf("video transcode is not required")
	}
	if err := mediatools.SafePathArg(inputPath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	args := []string{"-y", "-i", inputPath, "-map", "0", "-c", "copy"}
	if plan.TargetCodec != "" {
		args = append(args, "-c:v", ffmpegVideoCodec(plan.TargetCodec))
	}
	if plan.TargetPixel != "" {
		args = append(args, "-pix_fmt", plan.TargetPixel)
	}
	return append(args, outputPath), nil
}

func firstBaseVideoSource(item storage.MediaItem) storage.MediaComponentSource {
	for _, source := range item.ComponentSources {
		if source.SourceRole == "baseVideo" && source.RetentionState == "retained" {
			return source
		}
	}
	return storage.MediaComponentSource{}
}

func videoTargetMisses(stream VideoStreamInfo, target storage.MediaProfileVideoTarget) []string {
	misses := []string{}
	codec := normalizeVideoCodecName(stream.Codec)
	if len(target.Codecs) > 0 && codec != "" && !stringListHas(target.Codecs, codec) {
		misses = append(misses, "codec")
	}
	if len(target.HDRFormats) > 0 && stream.HDRFormat != "" && !stringListHas(target.HDRFormats, stream.HDRFormat) {
		misses = append(misses, "hdr")
	}
	if len(target.PixelFormats) > 0 && stream.PixelFormat != "" && !stringListHas(target.PixelFormats, stream.PixelFormat) {
		misses = append(misses, "pixel")
	}
	return misses
}

func videoStreamFromInventory(payload string) VideoStreamInfo {
	for _, stream := range videoInventoryStreams(payload) {
		if stream.Type == "video" {
			return VideoStreamInfo{Codec: stream.Codec, HDRFormat: stream.HDRFormat, PixelFormat: stream.PixelFormat}
		}
	}
	return VideoStreamInfo{}
}

type videoInventoryStream struct {
	Type        string `json:"type"`
	Codec       string `json:"codec"`
	HDRFormat   string `json:"hdrFormat"`
	PixelFormat string `json:"pixelFormat"`
}

func videoInventoryStreams(payload string) []videoInventoryStream {
	var list []videoInventoryStream
	if err := json.Unmarshal([]byte(payload), &list); err == nil {
		return list
	}
	var wrapped struct {
		Streams []videoInventoryStream `json:"streams"`
	}
	if err := json.Unmarshal([]byte(payload), &wrapped); err != nil {
		return nil
	}
	return wrapped.Streams
}

func videoFulfillmentProvenance(source storage.MediaComponentSource, status string) map[string]any {
	return map[string]any{
		"kind":       "videoFulfillment",
		"sourceId":   source.ID.String(),
		"sourcePath": source.RetainedPath,
		"status":     status,
	}
}

func normalizeVideoCodecName(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "x264", "h264", "avc":
		return "h264"
	case "x265", "h265", "hevc":
		return "hevc"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func ffmpegVideoCodec(value string) string {
	switch normalizeVideoCodecName(value) {
	case "h264":
		return "libx264"
	case "hevc":
		return "libx265"
	default:
		return normalizeVideoCodecName(value)
	}
}
