package httpapi

import (
	"encoding/json"

	"media-manager/internal/storage"
)

func dlnaRendererProfileListResponse(profiles []storage.DLNARendererProfile) DLNARendererProfileListResponse {
	response := DLNARendererProfileListResponse{Profiles: make([]DLNARendererProfile, 0, len(profiles))}
	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, dlnaRendererProfileResponse(profile))
	}
	return response
}

func dlnaRendererProfileResponse(profile storage.DLNARendererProfile) DLNARendererProfile {
	return DLNARendererProfile{
		Id:               profile.ID,
		Name:             profile.Name,
		Vendor:           profile.Vendor,
		DeviceClass:      profile.DeviceClass,
		Source:           DLNARendererProfileSource(profile.Source),
		SourceVersion:    profile.SourceVersion,
		Customized:       profile.Customized,
		Enabled:          profile.Enabled,
		Priority:         profile.Priority,
		IconKey:          profile.IconKey,
		Notes:            profile.Notes,
		MatchRules:       dlnaJSONResponse(profile.MatchRules),
		CapabilityRules:  dlnaJSONResponse(profile.CapabilityRules),
		DeliverySettings: dlnaJSONResponse(profile.DeliverySettings),
		DlnaFlags:        dlnaJSONResponse(profile.DLNAFlags),
		SubtitleRules:    dlnaJSONResponse(profile.SubtitleRules),
		ArtworkRules:     dlnaJSONResponse(profile.ArtworkRules),
		MetadataRules:    dlnaJSONResponse(profile.MetadataRules),
		Quirks:           dlnaJSONResponse(profile.Quirks),
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}
}

func dlnaRendererProfileInput(body DLNARendererProfileRequest) (storage.DLNARendererProfileInput, error) {
	matchRules, err := dlnaJSONRequest(body.MatchRules)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	capabilityRules, err := dlnaJSONRequest(body.CapabilityRules)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	deliverySettings, err := dlnaJSONRequest(body.DeliverySettings)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	dlnaFlags, err := dlnaJSONRequest(body.DlnaFlags)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	subtitleRules, err := dlnaJSONRequest(body.SubtitleRules)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	artworkRules, err := dlnaJSONRequest(body.ArtworkRules)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	metadataRules, err := dlnaJSONRequest(body.MetadataRules)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	quirks, err := dlnaJSONRequest(body.Quirks)
	if err != nil {
		return storage.DLNARendererProfileInput{}, err
	}
	return storage.DLNARendererProfileInput{
		Name:             body.Name,
		Vendor:           body.Vendor,
		DeviceClass:      body.DeviceClass,
		Enabled:          body.Enabled,
		Priority:         body.Priority,
		IconKey:          body.IconKey,
		Notes:            body.Notes,
		MatchRules:       matchRules,
		CapabilityRules:  capabilityRules,
		DeliverySettings: deliverySettings,
		DLNAFlags:        dlnaFlags,
		SubtitleRules:    subtitleRules,
		ArtworkRules:     artworkRules,
		MetadataRules:    metadataRules,
		Quirks:           quirks,
	}, nil
}

func dlnaRendererProfileCreateInput(
	body DLNARendererProfileCreateRequest,
) (storage.DLNARendererProfileInput, error) {
	return dlnaRendererProfileInput(DLNARendererProfileRequest{
		Name:             body.Name,
		Vendor:           body.Vendor,
		DeviceClass:      body.DeviceClass,
		Enabled:          body.Enabled,
		Priority:         body.Priority,
		IconKey:          body.IconKey,
		Notes:            body.Notes,
		MatchRules:       body.MatchRules,
		CapabilityRules:  body.CapabilityRules,
		DeliverySettings: body.DeliverySettings,
		DlnaFlags:        body.DlnaFlags,
		SubtitleRules:    body.SubtitleRules,
		ArtworkRules:     body.ArtworkRules,
		MetadataRules:    body.MetadataRules,
		Quirks:           body.Quirks,
	})
}

func dlnaJSONRequest(value DLNAJsonObject) ([]byte, error) {
	return json.Marshal(value)
}

func dlnaJSONResponse(value []byte) DLNAJsonObject {
	result := DLNAJsonObject{}
	_ = json.Unmarshal(value, &result)
	return result
}
