package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaProfile struct {
	ID                                string
	Name                              string
	IsDefault                         bool
	FinalContainer                    string
	UpgradesAllowed                   bool
	UpgradeUntilQualityID             *string
	MinimumCustomFormatScore          int32
	UpgradeUntilCustomFormatScore     int32
	MinimumCustomFormatScoreIncrement int32
	RemoveUnwantedAudio               bool
	RemoveUnwantedSubtitles           bool
	PreferredProtocol                 string
	SeriesPackPreference              string
	VideoTarget                       MediaProfileVideoTarget
	AudioTargets                      []MediaProfileAudioTarget
	SubtitleTargets                   []MediaProfileSubtitleTarget
	QualityIDs                        []string
	CustomFormatScores                []MediaProfileCustomFormatScore
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
}

type MediaProfileInput struct {
	Name                              string
	IsDefault                         bool
	FinalContainer                    string
	UpgradesAllowed                   bool
	UpgradeUntilQualityID             *string
	MinimumCustomFormatScore          int32
	UpgradeUntilCustomFormatScore     int32
	MinimumCustomFormatScoreIncrement int32
	RemoveUnwantedAudio               bool
	RemoveUnwantedSubtitles           bool
	PreferredProtocol                 string
	SeriesPackPreference              string
	VideoTarget                       MediaProfileVideoTarget
	AudioTargets                      []MediaProfileAudioTarget
	SubtitleTargets                   []MediaProfileSubtitleTarget
	QualityIDs                        []string
	CustomFormatScores                []MediaProfileCustomFormatScore
}

type MediaProfileVideoTarget struct {
	Codecs              []string
	CodecRequired       bool
	CodecScore          int32
	HDRFormats          []string
	HDRRequired         bool
	HDRScore            int32
	PixelFormats        []string
	PixelFormatRequired bool
	PixelFormatScore    int32
}

type MediaProfileAudioTarget struct {
	LanguageID           string
	Score                int32
	Required             bool
	Codecs               []string
	Channels             []string
	MinimumBitrateKbps   *int32
	PreferredBitrateKbps *int32
	LossyTranscodePolicy string
}

type MediaProfileCustomFormatScore struct {
	CustomFormatID uuid.UUID
	Score          int32
}

type MediaProfileSubtitleTarget struct {
	LanguageID string
	Score      int32
	Required   bool
	Source     string
	Formats    []string
}
