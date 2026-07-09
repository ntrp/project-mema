package jobs

import "media-manager/internal/storage"

func audioConversionInputForTrack(
	policy string,
	target storage.MediaProfileAudioTarget,
	track storage.MediaFileTrackFact,
) AudioConversionInput {
	targetCodec := target.TargetCodec
	if targetCodec == nil && track.Codec != nil {
		codec := *track.Codec
		targetCodec = &codec
	}
	return AudioConversionInput{
		Policy:               policy,
		SourceCodec:          stringPtrValue(track.Codec),
		SourceChannels:       stringPtrValue(track.Channels),
		SourceBitrateKbps:    int32PtrValue(track.BitrateKbps),
		TargetCodec:          targetCodec,
		TargetChannels:       target.TargetChannels,
		MinimumBitrateKbps:   target.MinimumBitrateKbps,
		PreferredBitrateKbps: target.PreferredBitrateKbps,
	}
}
