package playback

func BrowserVideoProfile() DeviceProfile {
	return DeviceProfile{
		Name:                "Browser HLS",
		MaxStreamingBitrate: 120_000_000,
		DirectPlayProfiles: []DirectPlayProfile{
			{
				Type:        ProfileVideo,
				Containers:  []string{"mp4", "m4v"},
				VideoCodecs: []string{"h264", "avc1"},
				AudioCodecs: []string{"aac"},
			},
		},
		TranscodingProfiles: []TranscodingProfile{
			{
				Type:          ProfileVideo,
				Container:     "ts",
				Protocol:      ProtocolHLS,
				VideoCodecs:   []string{"h264", "avc1"},
				AudioCodecs:   []string{"aac"},
				SegmentLength: 6,
			},
		},
		CodecProfiles: []CodecProfile{
			{
				Type:   CodecVideo,
				Codecs: []string{"h264", "avc1"},
				Conditions: []ProfileCondition{
					{
						Property: ConditionPixelFormat,
						Allowed:  []string{"", "yuv420p", "yuvj420p"},
					},
				},
			},
		},
	}
}
