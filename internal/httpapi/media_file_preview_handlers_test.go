package httpapi

import "testing"

func TestSCNMedia012PreviewInfoReportsDirectModeForCompatibleMp4(t *testing.T) {
	track := int32(1)
	info := mediaPreviewInfoFromTracks("/media/movie.mp4", []MediaFileTrack{
		{
			Type:        Video,
			Codec:       previewString("h264"),
			PixelFormat: previewString("yuv420p"),
			BitRate:     previewString("4000000"),
		},
		{Type: Audio, Index: &track, Codec: previewString("aac"), BitRate: previewString("640000")},
	}, &track, Browser)

	if info.StreamingMode != Direct || info.DeliveryProtocol != File {
		t.Fatalf("preview info = %#v, want direct file playback", info)
	}
}

func TestSCNMedia012PreviewInfoReportsRemuxModeAndSelectedBitrate(t *testing.T) {
	track := int32(2)
	info := mediaPreviewInfoFromTracks("/media/movie.mkv", []MediaFileTrack{
		{
			Type:        Video,
			Codec:       previewString("h264"),
			PixelFormat: previewString("yuv420p"),
			BitRate:     previewString("4000000"),
		},
		{Type: Audio, Index: &track, Codec: previewString("aac"), BitRate: previewString("640000")},
	}, &track, Browser)

	if info.StreamingMode != Remux || info.DeliveryProtocol != Hls {
		t.Fatalf("preview info = %#v, want HLS remux", info)
	}
	if info.LiveBitRate == nil || *info.LiveBitRate != "4640000" {
		t.Fatalf("live bit rate = %#v, want 4640000", info.LiveBitRate)
	}
	if info.VideoTrack == nil || info.AudioTrack == nil {
		t.Fatalf("expected selected video and audio tracks, got %#v", info)
	}
}

func previewString(value string) *string {
	return &value
}
