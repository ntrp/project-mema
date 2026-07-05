package httpapi

import (
	"net/http"
	"testing"
)

func TestScenarioSCNSettings023AdminManagesMediaProfilesAndQualitySizes(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-023")

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 || len(profiles.Profiles[0].QualityIds) == 0 {
		t.Fatalf("expected seeded media profile qualities: %#v", profiles)
	}
	qualityIDs := profiles.Profiles[0].QualityIds

	var created MediaProfile
	client.doJSON(t, http.MethodPost, "/settings/profiles", mediaProfileRequest("Scenario Profile", qualityIDs), http.StatusCreated, &created)
	if created.Name != "Scenario Profile" || len(created.QualityIds) != len(qualityIDs) {
		t.Fatalf("created media profile = %#v", created)
	}
	if len(created.SubtitleLanguages) != 1 || created.SubtitleLanguages[0].SubtitleType != MediaProfileSubtitleLanguageSubtitleTypeAny {
		t.Fatalf("created media profile subtitles = %#v", created.SubtitleLanguages)
	}

	var updated MediaProfile
	client.doJSON(t, http.MethodPut, "/settings/profiles/"+created.Id, mediaProfileRequest("Updated Scenario Profile", qualityIDs), http.StatusOK, &updated)
	if updated.Name != "Updated Scenario Profile" || updated.PreferredProtocol != MediaProfilePreferredProtocolUsenet {
		t.Fatalf("updated media profile = %#v", updated)
	}

	var listed MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &listed)
	if !mediaProfileListHas(listed.Profiles, updated.Id, "Updated Scenario Profile") {
		t.Fatalf("updated media profile not listed: %#v", listed.Profiles)
	}

	var qualitySizes QualitySizeSettingsResponse
	client.doJSON(t, http.MethodGet, "/settings/quality-sizes", nil, http.StatusOK, &qualitySizes)
	if len(qualitySizes.Qualities) == 0 {
		t.Fatal("expected seeded quality size settings")
	}
	qualityUpdate := QualitySizeSettingsUpdateRequest{Qualities: []QualitySizeSettingRequest{{
		QualityId:                qualitySizes.Qualities[0].QualityId,
		MinimumSizeMbPerMinute:   1.25,
		PreferredSizeMbPerMinute: float64Ptr(2.5),
		MaximumSizeMbPerMinute:   float64Ptr(4.75),
	}}}
	var updatedSizes QualitySizeSettingsResponse
	client.doJSON(t, http.MethodPut, "/settings/quality-sizes", qualityUpdate, http.StatusOK, &updatedSizes)
	if !qualitySizeListHas(updatedSizes.Qualities, qualityUpdate.Qualities[0].QualityId, 1.25) {
		t.Fatalf("updated quality sizes = %#v", updatedSizes.Qualities)
	}

	client.doJSON(t, http.MethodDelete, "/settings/profiles/"+updated.Id, nil, http.StatusNoContent, nil)
}

func mediaProfileRequest(name string, qualityIDs []string) MediaProfileRequest {
	return MediaProfileRequest{
		Name:                              name,
		QualityIds:                        append([]string(nil), qualityIDs...),
		UpgradesAllowed:                   true,
		UpgradeUntilQualityId:             stringPtr(qualityIDs[len(qualityIDs)-1]),
		MinimumCustomFormatScore:          0,
		UpgradeUntilCustomFormatScore:     50,
		MinimumCustomFormatScoreIncrement: 1,
		RemoveNonEnabledLanguages:         true,
		PreferredProtocol:                 MediaProfileRequestPreferredProtocolUsenet,
		SeriesPackPreference:              MediaProfileRequestSeriesPackPreferencePreferPacks,
		TargetLanguages:                   []string{"en"},
		TargetLanguageScores: []MediaProfileLanguageScore{{
			LanguageId: "en",
			Score:      100,
			Required:   true,
		}},
		SubtitleLanguages: []MediaProfileSubtitleLanguage{{
			LanguageId:   "en",
			Required:     true,
			SubtitleType: MediaProfileSubtitleLanguageSubtitleTypeAny,
		}},
		CustomFormatScores: []MediaProfileCustomFormatScore{},
	}
}

func mediaProfileListHas(profiles []MediaProfile, id string, name string) bool {
	for _, profile := range profiles {
		if profile.Id == id && profile.Name == name {
			return true
		}
	}
	return false
}

func qualitySizeListHas(settings []QualitySizeSetting, qualityID string, minimum float64) bool {
	for _, setting := range settings {
		if setting.QualityId == qualityID && setting.MinimumSizeMbPerMinute == minimum {
			return true
		}
	}
	return false
}

func float64Ptr(value float64) *float64 {
	return &value
}
