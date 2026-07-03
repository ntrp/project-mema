package decisions

import "testing"

func TestSCNMedia002ParseReleaseFileNameExtractsReleaseMetadata(t *testing.T) {
	cases := []struct {
		name       string
		title      string
		want       ParsedRelease
		wantLangs  []string
		wantSeason *int32
		wantEp     *int32
	}{
		{
			name:  "remux movie",
			title: "Scenario.Movie.2026.2160p.Remux.TrueHD.Atmos.7.1.x265-DEEP",
			want: ParsedRelease{
				ReleaseGroup:  "DEEP",
				Year:          "2026",
				Source:        "Remux",
				Resolution:    "2160p",
				VideoCodec:    "x265",
				AudioCodec:    "TrueHD/Atmos",
				AudioChannels: "7.1",
				QualityID:     "remux-2160p",
			},
		},
		{
			name:  "dvd movie",
			title: "Scenario.Movie.2026.DVD-R.Xvid.AC3.5.1.Extended-GRP",
			want: ParsedRelease{
				ReleaseGroup:  "GRP",
				Year:          "2026",
				Source:        "DVD-R",
				VideoCodec:    "Xvid",
				AudioCodec:    "DD",
				AudioChannels: "5.1",
				Edition:       "Extended",
				QualityID:     "dvd-r",
			},
		},
		{
			name:       "episode",
			title:      "Scenario.Show.S02E03.720p.HDTV.AAC.2.0.H264.v2.PROPER.MULTi-GRP",
			wantSeason: int32Pointer(2),
			wantEp:     int32Pointer(3),
			wantLangs:  []string{"Multiple"},
			want: ParsedRelease{
				ReleaseGroup:  "GRP",
				Source:        "HDTV",
				Resolution:    "720p",
				VideoCodec:    "x264",
				AudioCodec:    "AAC",
				AudioChannels: "2.0",
				Version:       "v2",
				Proper:        true,
				QualityID:     "hdtv-720p",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseReleaseFileName(tc.title)
			assertString(t, "release group", got.ReleaseGroup, tc.want.ReleaseGroup)
			assertString(t, "year", got.Year, tc.want.Year)
			assertString(t, "source", got.Source, tc.want.Source)
			assertString(t, "resolution", got.Resolution, tc.want.Resolution)
			assertString(t, "video codec", got.VideoCodec, tc.want.VideoCodec)
			assertString(t, "audio codec", got.AudioCodec, tc.want.AudioCodec)
			assertString(t, "audio channels", got.AudioChannels, tc.want.AudioChannels)
			assertString(t, "edition", got.Edition, tc.want.Edition)
			assertString(t, "version", got.Version, tc.want.Version)
			assertString(t, "quality id", got.QualityID, tc.want.QualityID)
			if got.Proper != tc.want.Proper {
				t.Fatalf("proper = %v, want %v", got.Proper, tc.want.Proper)
			}
			assertInt32Pointer(t, "season", got.SeasonNumber, tc.wantSeason)
			assertInt32Pointer(t, "episode", got.EpisodeNumber, tc.wantEp)
			assertStrings(t, "languages", got.Languages, tc.wantLangs)
		})
	}
}

func int32Pointer(value int32) *int32 {
	return &value
}

func assertString(t *testing.T, label string, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %q, want %q", label, got, want)
	}
}

func assertInt32Pointer(t *testing.T, label string, got *int32, want *int32) {
	t.Helper()
	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil || *got != *want {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
}

func assertStrings(t *testing.T, label string, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %#v, want %#v", label, got, want)
	}
	for index, value := range want {
		if got[index] != value {
			t.Fatalf("%s[%d] = %q, want %q", label, index, got[index], value)
		}
	}
}
