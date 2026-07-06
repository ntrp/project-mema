package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaProfile struct {
	ID                                string
	Name                              string
	IsDefault                         bool
	QualityIDs                        []string
	UpgradesAllowed                   bool
	UpgradeUntilQualityID             *string
	MinimumCustomFormatScore          int32
	UpgradeUntilCustomFormatScore     int32
	MinimumCustomFormatScoreIncrement int32
	RemoveNonEnabledLanguages         bool
	RemoveNonEnabledSubtitleLanguages bool
	PreferredProtocol                 string
	SeriesPackPreference              string
	TargetLanguages                   []string
	TargetLanguageScores              []MediaProfileLanguageScore
	SubtitleLanguages                 []MediaProfileSubtitleLanguage
	ComponentTargets                  []MediaProfileComponentTarget
	CustomFormatScores                []MediaProfileCustomFormatScore
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
}

type MediaProfileInput struct {
	Name                              string
	IsDefault                         bool
	QualityIDs                        []string
	UpgradesAllowed                   bool
	UpgradeUntilQualityID             *string
	MinimumCustomFormatScore          int32
	UpgradeUntilCustomFormatScore     int32
	MinimumCustomFormatScoreIncrement int32
	RemoveNonEnabledLanguages         bool
	RemoveNonEnabledSubtitleLanguages bool
	PreferredProtocol                 string
	SeriesPackPreference              string
	TargetLanguages                   []string
	TargetLanguageScores              []MediaProfileLanguageScore
	SubtitleLanguages                 []MediaProfileSubtitleLanguage
	ComponentTargets                  []MediaProfileComponentTarget
	CustomFormatScores                []MediaProfileCustomFormatScore
}

type MediaProfileLanguageScore struct {
	LanguageID string
	Score      int32
	Required   bool
}

type MediaProfileCustomFormatScore struct {
	CustomFormatID uuid.UUID
	Score          int32
}

type MediaProfileSubtitleLanguage struct {
	LanguageID   string
	Score        int32
	Required     bool
	SubtitleType string
}

type MediaProfileComponentTarget struct {
	ID               uuid.UUID
	ComponentType    string
	Required         bool
	LanguageID       *string
	Codec            *string
	Channels         *string
	Source           string
	FallbackBehavior string
}
