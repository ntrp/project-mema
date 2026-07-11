package httpapi

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings025AdminManagesDLNARendererProfiles(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-025")
	suffix := uuid.NewString()[:8]
	cloneID := "api-clone-" + suffix
	importID := "api-import-" + suffix

	var listed DLNARendererProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/dlna/profiles", nil, http.StatusOK, &listed)
	if len(listed.Profiles) == 0 {
		t.Fatal("expected seeded DLNA profiles")
	}

	var generic DLNARendererProfile
	client.doJSON(t, http.MethodGet, "/settings/dlna/profiles/generic", nil, http.StatusOK, &generic)
	seedGeneric := generic
	client.doJSON(t, http.MethodPost, "/settings/dlna/profiles/generic/clone", DLNARendererProfileCloneRequest{
		Id:   cloneID,
		Name: "API Clone",
	}, http.StatusCreated, &generic)
	if generic.Id != cloneID || !generic.Customized {
		t.Fatalf("cloned profile = %#v", generic)
	}

	update := dlnaProfileRequestFromResponse(generic)
	update.Notes = "updated through API"
	update.Priority = 222
	var updated DLNARendererProfile
	client.doJSON(t, http.MethodPut, "/settings/dlna/profiles/"+cloneID, update, http.StatusOK, &updated)
	if updated.Notes != "updated through API" || updated.Priority != 222 {
		t.Fatalf("updated profile = %#v", updated)
	}

	var exported DLNARendererProfile
	client.doJSON(t, http.MethodGet, "/settings/dlna/profiles/"+cloneID+"/export", nil, http.StatusOK, &exported)
	importRequest := dlnaProfileCreateFromResponse(importID, exported)
	importRequest.Name = "API Import"
	var imported DLNARendererProfile
	client.doJSON(t, http.MethodPost, "/settings/dlna/profiles/import", importRequest, http.StatusOK, &imported)
	if imported.Id != importID || imported.Name != "API Import" {
		t.Fatalf("imported profile = %#v", imported)
	}

	resetRequest := dlnaProfileRequestFromResponse(seedGeneric)
	resetRequest.Notes = "temporary generic edit"
	client.doJSON(t, http.MethodPut, "/settings/dlna/profiles/generic", resetRequest, http.StatusOK, &generic)
	client.doJSON(t, http.MethodPost, "/settings/dlna/profiles/generic/reset", nil, http.StatusOK, &generic)
	if generic.Customized || generic.Notes == "temporary generic edit" {
		t.Fatalf("reset profile = %#v", generic)
	}
	client.doJSON(t, http.MethodDelete, "/settings/dlna/profiles/generic", nil, http.StatusBadRequest, nil)

	update = dlnaProfileRequestFromResponse(generic)
	update.Notes = "temporary restore edit"
	client.doJSON(t, http.MethodPut, "/settings/dlna/profiles/generic", update, http.StatusOK, &generic)
	client.doJSON(t, http.MethodPost, "/settings/dlna/profiles/restore", nil, http.StatusNoContent, nil)
	client.doJSON(t, http.MethodGet, "/settings/dlna/profiles/generic", nil, http.StatusOK, &generic)
	if generic.Customized || generic.Notes == "temporary restore edit" {
		t.Fatalf("restored seeded profile = %#v", generic)
	}

	client.doJSON(t, http.MethodPost, "/settings/dlna/profiles", invalidDLNARegexProfile("bad-regex-"+suffix), http.StatusBadRequest, nil)

	ip := "192.0.2.88"
	var override DLNARendererDeviceOverride
	client.doJSON(t, http.MethodPost, "/settings/dlna/device-overrides", DLNARendererDeviceOverrideRequest{
		IpAddress:               &ip,
		ProfileId:               cloneID,
		DisplayName:             "API Device",
		Allowed:                 true,
		DeliveryPolicyOverrides: DLNAJsonObject{},
		Notes:                   "",
	}, http.StatusOK, &override)
	if override.ProfileId != cloneID || override.IpAddress == nil || *override.IpAddress != ip {
		t.Fatalf("override = %#v", override)
	}

	var overrides DLNARendererDeviceOverrideListResponse
	client.doJSON(t, http.MethodGet, "/settings/dlna/device-overrides", nil, http.StatusOK, &overrides)
	if len(overrides.Overrides) == 0 {
		t.Fatal("expected saved override")
	}
	var devices DLNARecentDeviceListResponse
	client.doJSON(t, http.MethodGet, "/settings/dlna/recent-devices", nil, http.StatusOK, &devices)
	if devices.Devices == nil {
		t.Fatalf("recent devices should be present: %#v", devices)
	}

	client.doJSON(t, http.MethodDelete, "/settings/dlna/device-overrides/"+override.Id.String(), nil, http.StatusNoContent, nil)
	client.doJSON(t, http.MethodDelete, "/settings/dlna/profiles/"+cloneID, nil, http.StatusNoContent, nil)
	client.doJSON(t, http.MethodDelete, "/settings/dlna/profiles/"+importID, nil, http.StatusNoContent, nil)
}

func dlnaProfileRequestFromResponse(profile DLNARendererProfile) DLNARendererProfileRequest {
	return DLNARendererProfileRequest{
		Name:             profile.Name,
		Vendor:           profile.Vendor,
		DeviceClass:      profile.DeviceClass,
		Enabled:          profile.Enabled,
		Priority:         profile.Priority,
		IconKey:          profile.IconKey,
		Notes:            profile.Notes,
		MatchRules:       profile.MatchRules,
		CapabilityRules:  profile.CapabilityRules,
		DeliverySettings: profile.DeliverySettings,
		DlnaFlags:        profile.DlnaFlags,
		SubtitleRules:    profile.SubtitleRules,
		ArtworkRules:     profile.ArtworkRules,
		MetadataRules:    profile.MetadataRules,
		Quirks:           profile.Quirks,
	}
}

func dlnaProfileCreateFromResponse(id string, profile DLNARendererProfile) DLNARendererProfileCreateRequest {
	request := DLNARendererProfileCreateRequest{Id: id}
	base := dlnaProfileRequestFromResponse(profile)
	request.Name = base.Name
	request.Vendor = base.Vendor
	request.DeviceClass = base.DeviceClass
	request.Enabled = base.Enabled
	request.Priority = base.Priority
	request.IconKey = base.IconKey
	request.Notes = base.Notes
	request.MatchRules = base.MatchRules
	request.CapabilityRules = base.CapabilityRules
	request.DeliverySettings = base.DeliverySettings
	request.DlnaFlags = base.DlnaFlags
	request.SubtitleRules = base.SubtitleRules
	request.ArtworkRules = base.ArtworkRules
	request.MetadataRules = base.MetadataRules
	request.Quirks = base.Quirks
	return request
}

func invalidDLNARegexProfile(id string) DLNARendererProfileCreateRequest {
	request := dlnaProfileCreateFromResponse(id, DLNARendererProfile{
		Name:             "Bad Regex",
		Vendor:           "Test",
		DeviceClass:      "software",
		Enabled:          true,
		Priority:         10,
		IconKey:          "",
		Notes:            "",
		MatchRules:       DLNAJsonObject{"tokens": []any{map[string]any{"kind": "regex", "value": "["}}},
		CapabilityRules:  DLNAJsonObject{},
		DeliverySettings: DLNAJsonObject{},
		DlnaFlags:        DLNAJsonObject{},
		SubtitleRules:    DLNAJsonObject{},
		ArtworkRules:     DLNAJsonObject{},
		MetadataRules:    DLNAJsonObject{},
		Quirks:           DLNAJsonObject{},
	})
	return request
}
