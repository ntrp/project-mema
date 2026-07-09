package jobs

import (
	"strings"

	"media-manager/internal/storage"
)

type VideoConversionInput struct {
	SourceCodec string
	SourcePixel string
	TargetCodec string
	TargetPixel string
}

type VideoConversionDecision struct {
	Needed      bool
	Allowed     bool
	Status      string
	Reason      string
	SourceCodec string
	SourcePixel string
	TargetCodec string
	TargetPixel string
}

func DecideVideoConversion(input VideoConversionInput) VideoConversionDecision {
	decision := VideoConversionDecision{
		SourceCodec: normalizeVideoCodecName(input.SourceCodec),
		SourcePixel: strings.TrimSpace(input.SourcePixel),
		TargetCodec: normalizeVideoCodecName(input.TargetCodec),
		TargetPixel: strings.TrimSpace(input.TargetPixel),
	}
	if decision.TargetCodec != "" && decision.SourceCodec != "" && decision.SourceCodec != decision.TargetCodec {
		decision.Needed = true
	}
	if decision.TargetPixel != "" && decision.SourcePixel != "" && decision.SourcePixel != decision.TargetPixel {
		decision.Needed = true
	}
	if !decision.Needed {
		decision.Status = "satisfied"
		decision.Reason = "Video already meets supported target details."
		return decision
	}
	if decision.TargetCodec == "" && decision.TargetPixel == "" {
		decision.Status = "blocked"
		decision.Reason = "Video target mismatch has no supported transcode target."
		return decision
	}
	decision.Allowed = true
	decision.Status = "allowed"
	decision.Reason = "Video codec or pixel format conversion is supported."
	return decision
}

func videoConversionInputForTrack(
	target storage.MediaProfileVideoTarget,
	track storage.MediaFileTrackFact,
) VideoConversionInput {
	return VideoConversionInput{
		SourceCodec: stringPtrValue(track.Codec),
		SourcePixel: stringPtrValue(track.PixelFormat),
		TargetCodec: firstString(target.Codecs),
		TargetPixel: firstString(target.PixelFormats),
	}
}
