package playback

import "strings"

func BuildVideoStreamInfo(source MediaSource, options MediaOptions) StreamInfo {
	options = normalizedOptions(options)
	audio := selectAudioStream(source.AudioStreams, options.AudioStreamIndex)
	directReasons := directPlayFailureReasons(source, options, audio)
	if len(directReasons) == 0 && options.EnableDirectPlay && source.SupportsDirectPlay {
		return StreamInfo{
			PlayMethod:       PlayMethodDirectPlay,
			Container:        normalizedContainer(source.Container),
			Protocol:         ProtocolHTTP,
			OutputVideoCodec: "copy",
			OutputAudioCodec: "copy",
			AudioStreamIndex: audioIndex(audio),
		}
	}
	if info, ok := directStreamInfo(source, options, audio, directReasons); ok {
		return info
	}
	return transcodeInfo(source, options, audio, directReasons)
}

func normalizedOptions(options MediaOptions) MediaOptions {
	if options.Profile.Name == "" {
		options.Profile = BrowserVideoProfile()
	}
	if !options.EnableDirectPlay && !options.EnableDirectStream && !options.EnableTranscoding {
		options.EnableDirectPlay = true
		options.EnableDirectStream = true
		options.EnableTranscoding = true
	}
	if !options.AllowVideoStreamCopy && !options.AllowAudioStreamCopy {
		options.AllowVideoStreamCopy = true
		options.AllowAudioStreamCopy = true
	}
	return options
}

func directStreamInfo(source MediaSource, options MediaOptions, audio *MediaStream, directReasons []TranscodeReason) (StreamInfo, bool) {
	if !options.EnableDirectStream || !source.SupportsDirectStream {
		return StreamInfo{}, false
	}
	profile, ok := firstTranscodingProfile(options.Profile, ProfileVideo)
	if !ok {
		return StreamInfo{}, false
	}
	if !videoCopySupported(source.Video, profile, options.Profile) {
		return StreamInfo{}, false
	}
	if !audioCopySupported(audio, profile) {
		return StreamInfo{}, false
	}
	return StreamInfo{
		PlayMethod:       PlayMethodDirectStream,
		Container:        profile.Container,
		Protocol:         profile.Protocol,
		OutputVideoCodec: "copy",
		OutputAudioCodec: "copy",
		AudioStreamIndex: audioIndex(audio),
		SegmentLength:    profile.SegmentLength,
		TranscodeReasons: compactReasons(directReasons),
	}, true
}

func transcodeInfo(source MediaSource, options MediaOptions, audio *MediaStream, directReasons []TranscodeReason) StreamInfo {
	profile, ok := firstTranscodingProfile(options.Profile, ProfileVideo)
	if !ok || !options.EnableTranscoding || !source.SupportsTranscoding {
		return StreamInfo{TranscodeReasons: compactReasons(directReasons)}
	}
	videoCodec := firstCodec(profile.VideoCodecs)
	if options.AllowVideoStreamCopy && videoCopySupported(source.Video, profile, options.Profile) {
		videoCodec = "copy"
	}
	audioCodec := firstCodec(profile.AudioCodecs)
	if audio == nil || (options.AllowAudioStreamCopy && audioCopySupported(audio, profile)) {
		audioCodec = "copy"
	}
	return StreamInfo{
		PlayMethod:       PlayMethodTranscode,
		Container:        profile.Container,
		Protocol:         profile.Protocol,
		OutputVideoCodec: videoCodec,
		OutputAudioCodec: audioCodec,
		AudioStreamIndex: audioIndex(audio),
		SegmentLength:    profile.SegmentLength,
		TranscodeReasons: compactReasons(directReasons),
	}
}

func directPlayFailureReasons(source MediaSource, options MediaOptions, audio *MediaStream) []TranscodeReason {
	reasons := []TranscodeReason{}
	if bitrateExceeded(source, options) {
		reasons = append(reasons, ReasonContainerBitrateExceeded)
	}
	if _, ok := matchingDirectPlayProfile(source, options.Profile, audio); !ok {
		reasons = append(reasons, directProfileReasons(source, options.Profile, audio)...)
	}
	if audio != nil && !selectedAudioIsFirst(source.AudioStreams, audio.Index) {
		reasons = append(reasons, ReasonSecondaryAudioNotSupported)
	}
	if profileSupportsCodecConditions(source.Video, options.Profile) != "" {
		reasons = append(reasons, ReasonVideoPixelFormatUnsupported)
	}
	return compactReasons(reasons)
}

func matchingDirectPlayProfile(source MediaSource, profile DeviceProfile, audio *MediaStream) (DirectPlayProfile, bool) {
	for _, playProfile := range profile.DirectPlayProfiles {
		if playProfile.Type != ProfileVideo {
			continue
		}
		if !playProfile.SupportsContainer(normalizedContainer(source.Container)) {
			continue
		}
		if source.Video != nil && !playProfile.SupportsVideoCodec(source.Video.Codec) {
			continue
		}
		if audio != nil && !playProfile.SupportsAudioCodec(audio.Codec) {
			continue
		}
		return playProfile, true
	}
	return DirectPlayProfile{}, false
}

func directProfileReasons(source MediaSource, profile DeviceProfile, audio *MediaStream) []TranscodeReason {
	container := normalizedContainer(source.Container)
	containerSupported := false
	videoSupported := false
	audioSupported := audio == nil
	for _, playProfile := range profile.DirectPlayProfiles {
		if playProfile.Type != ProfileVideo {
			continue
		}
		if playProfile.SupportsContainer(container) {
			containerSupported = true
		}
		if source.Video != nil && playProfile.SupportsVideoCodec(source.Video.Codec) {
			videoSupported = true
		}
		if audio == nil || playProfile.SupportsAudioCodec(audio.Codec) {
			audioSupported = true
		}
	}
	reasons := []TranscodeReason{}
	if !containerSupported {
		reasons = append(reasons, ReasonContainerNotSupported)
	}
	if source.Video != nil && !videoSupported {
		reasons = append(reasons, ReasonVideoCodecNotSupported)
	}
	if !audioSupported {
		reasons = append(reasons, ReasonAudioCodecNotSupported)
	}
	return reasons
}

func firstTranscodingProfile(profile DeviceProfile, profileType ProfileType) (TranscodingProfile, bool) {
	for _, transcodingProfile := range profile.TranscodingProfiles {
		if transcodingProfile.Type == profileType {
			return transcodingProfile, true
		}
	}
	return TranscodingProfile{}, false
}

func videoCopySupported(video *MediaStream, profile TranscodingProfile, device DeviceProfile) bool {
	return video != nil && profile.SupportsVideoCodec(video.Codec) && profileSupportsCodecConditions(video, device) == ""
}

func audioCopySupported(audio *MediaStream, profile TranscodingProfile) bool {
	return audio == nil || profile.SupportsAudioCodec(audio.Codec)
}

func profileSupportsCodecConditions(stream *MediaStream, profile DeviceProfile) TranscodeReason {
	if stream == nil {
		return ""
	}
	for _, codecProfile := range profile.CodecProfiles {
		if codecProfile.Type != streamCodecType(stream.Type) || !containsToken(codecProfile.Codecs, stream.Codec) {
			continue
		}
		for _, condition := range codecProfile.Conditions {
			if !conditionPasses(stream, condition) {
				if condition.Property == ConditionPixelFormat {
					return ReasonVideoPixelFormatUnsupported
				}
			}
		}
	}
	return ""
}

func conditionPasses(stream *MediaStream, condition ProfileCondition) bool {
	switch condition.Property {
	case ConditionPixelFormat:
		return containsToken(condition.Allowed, stream.PixelFormat)
	case ConditionAudioChannels:
		return condition.Maximum <= 0 || stream.Channels <= condition.Maximum
	default:
		return true
	}
}

func selectAudioStream(streams []MediaStream, index *int) *MediaStream {
	for i := range streams {
		if index == nil || streams[i].Index == *index {
			return &streams[i]
		}
	}
	return nil
}

func selectedAudioIsFirst(streams []MediaStream, index int) bool {
	return len(streams) > 0 && streams[0].Index == index
}

func bitrateExceeded(source MediaSource, options MediaOptions) bool {
	limit := options.MaxStreamingBitrate
	if limit <= 0 {
		limit = options.Profile.MaxStreamingBitrate
	}
	return limit > 0 && source.BitRate > 0 && source.BitRate > limit
}

func compactReasons(reasons []TranscodeReason) []TranscodeReason {
	seen := map[TranscodeReason]bool{}
	compacted := []TranscodeReason{}
	for _, reason := range reasons {
		if reason == "" || seen[reason] {
			continue
		}
		seen[reason] = true
		compacted = append(compacted, reason)
	}
	return compacted
}

func audioIndex(audio *MediaStream) *int {
	if audio == nil {
		return nil
	}
	index := audio.Index
	return &index
}

func firstCodec(codecs []string) string {
	if len(codecs) == 0 {
		return ""
	}
	return normalizeToken(codecs[0])
}

func normalizedContainer(container string) string {
	parts := strings.Split(container, ",")
	return normalizeToken(parts[0])
}

func streamCodecType(streamType StreamType) CodecType {
	if streamType == StreamAudio {
		return CodecAudio
	}
	return CodecVideo
}
