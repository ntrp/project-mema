package httpapi

import "media-manager/internal/delivery"

func deliveryTracks(tracks []MediaFileTrack) []delivery.Track {
	results := make([]delivery.Track, 0, len(tracks))
	for _, track := range tracks {
		results = append(results, delivery.Track{
			Index:         track.Index,
			Type:          deliveryTrackType(track.Type),
			Codec:         track.Codec,
			Language:      track.Language,
			Title:         track.Title,
			BitRate:       track.BitRate,
			ChannelLayout: track.ChannelLayout,
			FrameRate:     track.FrameRate,
			Height:        track.Height,
			Width:         track.Width,
			PixelFormat:   track.PixelFormat,
			Profile:       track.Profile,
			Channels:      track.Channels,
		})
	}
	return results
}

func mediaFileTracksFromDelivery(tracks []delivery.Track) []MediaFileTrack {
	results := make([]MediaFileTrack, 0, len(tracks))
	for _, track := range tracks {
		results = append(results, MediaFileTrack{
			Index:         track.Index,
			Type:          mediaFileTrackTypeFromDelivery(track.Type),
			Codec:         track.Codec,
			Language:      track.Language,
			Title:         track.Title,
			BitRate:       track.BitRate,
			ChannelLayout: track.ChannelLayout,
			FrameRate:     track.FrameRate,
			Height:        track.Height,
			Width:         track.Width,
			PixelFormat:   track.PixelFormat,
			Profile:       track.Profile,
			Channels:      track.Channels,
		})
	}
	return results
}

func mediaFileChaptersFromDelivery(chapters []delivery.Chapter) []MediaFileChapter {
	results := make([]MediaFileChapter, 0, len(chapters))
	for _, chapter := range chapters {
		results = append(results, MediaFileChapter{
			Index:     chapter.Index,
			Title:     chapter.Title,
			StartTime: chapter.StartTime,
			EndTime:   chapter.EndTime,
		})
	}
	return results
}

func deliveryTrack(track *MediaFileTrack) *delivery.Track {
	if track == nil {
		return nil
	}
	result := deliveryTracks([]MediaFileTrack{*track})
	return &result[0]
}

func mediaFileTrackFromDelivery(track *delivery.Track) *MediaFileTrack {
	if track == nil {
		return nil
	}
	result := mediaFileTracksFromDelivery([]delivery.Track{*track})
	return &result[0]
}

func deliveryTrackType(value MediaFileTrackType) delivery.TrackType {
	switch value {
	case Video:
		return delivery.TrackVideo
	case Audio:
		return delivery.TrackAudio
	case Subtitle:
		return delivery.TrackSubtitle
	default:
		return ""
	}
}

func mediaFileTrackTypeFromDelivery(value delivery.TrackType) MediaFileTrackType {
	switch value {
	case delivery.TrackVideo:
		return Video
	case delivery.TrackAudio:
		return Audio
	case delivery.TrackSubtitle:
		return Subtitle
	default:
		return ""
	}
}

func deliveryClientProfile(value MediaFilePreviewClientProfile) delivery.ClientProfile {
	switch value {
	case Webkit:
		return delivery.ClientWebKit
	default:
		return delivery.ClientBrowser
	}
}

func deliveryClientProfilePointer(value *MediaFilePreviewClientProfile) *delivery.ClientProfile {
	if value == nil {
		return nil
	}
	profile := deliveryClientProfile(*value)
	return &profile
}

func mediaPreviewClientProfile(explicit *MediaFilePreviewClientProfile, userAgent string) MediaFilePreviewClientProfile {
	profile := delivery.ClientProfileForRequest(deliveryClientProfilePointer(explicit), userAgent)
	if profile == delivery.ClientWebKit {
		return Webkit
	}
	return Browser
}

func mediaPreviewMode(value delivery.Mode) MediaFilePreviewMode {
	switch value {
	case delivery.ModeDirect:
		return Direct
	case delivery.ModeRemux:
		return Remux
	default:
		return Transcode
	}
}

func mediaPreviewProtocol(value delivery.Protocol) MediaFilePreviewDeliveryProtocol {
	switch value {
	case delivery.ProtocolFile:
		return File
	default:
		return Hls
	}
}
