package delivery

import "testing"

func TestSCNMedia001ProbeStringNormalization(t *testing.T) {
	for _, value := range []string{"", " ", "unknown", "UNKNOWN", "n/a"} {
		if optionalString(value) != nil {
			t.Fatalf("expected %q to normalize to nil", value)
		}
	}

	got := optionalString(" eng ")
	if got == nil || *got != "eng" {
		t.Fatalf("optionalString = %#v, want eng", got)
	}
}

func TestSCNMedia001ProbeNumberNormalization(t *testing.T) {
	if optionalInt(0) != nil || optionalInt(-1) != nil {
		t.Fatal("expected non-positive probe integers to normalize to nil")
	}
	if optionalIndex(-1) != nil {
		t.Fatal("expected negative probe index to normalize to nil")
	}
	if got := optionalInt(24); got == nil || *got != 24 {
		t.Fatalf("optionalInt = %#v, want 24", got)
	}
	if got := optionalIndex(0); got == nil || *got != 0 {
		t.Fatalf("optionalIndex = %#v, want 0", got)
	}
}

func TestSCNMedia001ProbeLanguageAndFrameRateNormalization(t *testing.T) {
	if got := languageTag(map[string]string{"LANGUAGE": "deu"}); got != "deu" {
		t.Fatalf("languageTag uppercase = %q, want deu", got)
	}
	if got := languageTag(map[string]string{"language": "eng", "LANGUAGE": "deu"}); got != "eng" {
		t.Fatalf("languageTag lowercase precedence = %q, want eng", got)
	}
	if got := normalFrameRate(" 24000/1001 "); got != "24000/1001" {
		t.Fatalf("normalFrameRate = %q, want 24000/1001", got)
	}
	if normalFrameRate("0/0") != "" {
		t.Fatal("expected zero frame rate to normalize to empty string")
	}
}

func TestSCNMedia001ProbeTracksAndChapters(t *testing.T) {
	tracks := tracks([]ffprobeStream{
		{
			Index:         0,
			CodecName:     "h264",
			CodecType:     "video",
			Profile:       "High",
			Width:         1920,
			Height:        1080,
			PixelFormat:   "yuv420p",
			FrameRate:     "24000/1001",
			BitRate:       "8000000",
			Tags:          map[string]string{"language": "eng", "title": "Main video"},
			ChannelLayout: "ignored",
		},
		{
			Index:         1,
			CodecName:     "aac",
			CodecType:     "audio",
			Channels:      6,
			ChannelLayout: "5.1",
			Tags:          map[string]string{"LANGUAGE": "deu", "BPS-eng": "640000"},
		},
		{Index: 2, CodecType: "data"},
	})

	if len(tracks) != 2 {
		t.Fatalf("tracks = %#v, want 2 supported tracks", tracks)
	}
	if tracks[0].Type != TrackVideo || tracks[0].Codec == nil || *tracks[0].Codec != "h264" {
		t.Fatalf("video track = %#v", tracks[0])
	}
	if tracks[0].Language == nil || *tracks[0].Language != "eng" ||
		tracks[0].Height == nil || *tracks[0].Height != 1080 {
		t.Fatalf("video track details = %#v", tracks[0])
	}
	if tracks[1].Type != TrackAudio || tracks[1].Language == nil || *tracks[1].Language != "deu" {
		t.Fatalf("audio track = %#v", tracks[1])
	}
	if tracks[1].BitRate == nil || *tracks[1].BitRate != "640000" {
		t.Fatalf("audio bitrate = %#v", tracks[1].BitRate)
	}

	chapters := chapters([]ffprobeChapter{
		{ID: 0, StartTime: "0.0", EndTime: "60.0", Tags: map[string]string{"title": "Intro"}},
		{ID: 3, StartTime: "60.0", EndTime: "120.0", Tags: map[string]string{}},
	})
	if len(chapters) != 2 || chapters[0].Index != 0 || chapters[1].Index != 3 {
		t.Fatalf("chapters = %#v", chapters)
	}
	if chapters[0].Title == nil || *chapters[0].Title != "Intro" {
		t.Fatalf("chapter title = %#v", chapters[0])
	}
}

func TestSCNMedia001ProbeTrackBitRateFromByteTags(t *testing.T) {
	track, ok := trackFromStream(ffprobeStream{
		Index:     1,
		CodecName: "flac",
		CodecType: "audio",
		Tags: map[string]string{
			"NUMBER_OF_BYTES-eng": "9600000",
			"DURATION-eng":        "00:02:00.000000000",
		},
	})

	if !ok {
		t.Fatal("expected audio track")
	}
	if track.BitRate == nil || *track.BitRate != "640000" {
		t.Fatalf("derived audio bitrate = %#v, want 640000", track.BitRate)
	}
}

func TestSCNMedia001ProbeTrackTypeMapping(t *testing.T) {
	for input, want := range map[string]TrackType{
		"VIDEO":    TrackVideo,
		"audio":    TrackAudio,
		"Subtitle": TrackSubtitle,
	} {
		got, ok := trackTypeFromCodec(input)
		if !ok || got != want {
			t.Fatalf("trackTypeFromCodec(%q) = %q, %v; want %q, true", input, got, ok, want)
		}
	}
	if _, ok := trackTypeFromCodec("data"); ok {
		t.Fatal("expected unsupported track type to be rejected")
	}
}
