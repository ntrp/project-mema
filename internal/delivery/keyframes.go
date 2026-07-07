package delivery

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	mediatools "media-manager/internal/tools"
)

type ffprobeKeyframeOutput struct {
	Frames []ffprobeKeyframe `json:"frames"`
}

type ffprobeKeyframe struct {
	BestEffortTimestampTime string `json:"best_effort_timestamp_time"`
	PktPtsTime              string `json:"pkt_pts_time"`
}

func VideoKeyframes(target string) []float64 {
	if _, err := mediatools.LookPath("ffprobe"); err != nil {
		return nil
	}
	if err := mediatools.SafePathArg(target); err != nil {
		return nil
	}
	output, err := mediatools.RunOutput(context.Background(), mediatools.CommandSpec{
		Name: "ffprobe",
		Args: []string{
			"-v", "error",
			"-select_streams", "v:0",
			"-skip_frame", "nokey",
			"-show_frames",
			"-show_entries", "frame=best_effort_timestamp_time,pkt_pts_time",
			"-of", "json",
			target,
		},
		Timeout:        5 * time.Second,
		MaxOutputBytes: 8 * 1024 * 1024,
		MaxStderrBytes: 64 * 1024,
	})
	if err != nil {
		return nil
	}
	var payload ffprobeKeyframeOutput
	if err := json.Unmarshal(output, &payload); err != nil {
		return nil
	}
	return NormalizeKeyframes(payload.Frames)
}

func NormalizeKeyframes(frames []ffprobeKeyframe) []float64 {
	values := make([]float64, 0, len(frames))
	for _, frame := range frames {
		value, ok := keyframeTime(frame)
		if ok {
			values = append(values, value)
		}
	}
	sort.Float64s(values)
	return compactKeyframes(values)
}

func keyframeTime(frame ffprobeKeyframe) (float64, bool) {
	for _, value := range []string{frame.BestEffortTimestampTime, frame.PktPtsTime} {
		parsed, err := strconv.ParseFloat(value, 64)
		if err == nil && parsed >= 0 {
			return parsed, true
		}
	}
	return 0, false
}

func compactKeyframes(values []float64) []float64 {
	keyframes := []float64{}
	for _, value := range values {
		if len(keyframes) == 0 || value-keyframes[len(keyframes)-1] > 0.05 {
			keyframes = append(keyframes, value)
		}
	}
	return keyframes
}
