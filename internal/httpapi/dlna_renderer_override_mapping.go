package httpapi

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func dlnaRendererOverrideListResponse(
	overrides []storage.DLNARendererDeviceOverride,
) DLNARendererDeviceOverrideListResponse {
	response := DLNARendererDeviceOverrideListResponse{
		Overrides: make([]DLNARendererDeviceOverride, 0, len(overrides)),
	}
	for _, override := range overrides {
		response.Overrides = append(response.Overrides, dlnaRendererOverrideResponse(override))
	}
	return response
}

func dlnaRendererOverrideResponse(override storage.DLNARendererDeviceOverride) DLNARendererDeviceOverride {
	return DLNARendererDeviceOverride{
		Id:                      openapi_types.UUID(override.ID),
		RendererUuid:            override.RendererUUID,
		IpAddress:               override.IPAddress,
		ProfileId:               override.ProfileID,
		DisplayName:             override.DisplayName,
		Allowed:                 override.Allowed,
		DeliveryPolicyOverrides: dlnaJSONResponse(override.DeliveryPolicyOverrides),
		Notes:                   override.Notes,
		CreatedAt:               override.CreatedAt,
		UpdatedAt:               override.UpdatedAt,
	}
}

func dlnaRendererOverrideInput(
	body DLNARendererDeviceOverrideRequest,
) (storage.DLNARendererDeviceOverrideInput, error) {
	policy, err := dlnaJSONRequest(body.DeliveryPolicyOverrides)
	if err != nil {
		return storage.DLNARendererDeviceOverrideInput{}, err
	}
	id := uuid.Nil
	if body.Id != nil {
		id = uuid.UUID(*body.Id)
	}
	return storage.DLNARendererDeviceOverrideInput{
		ID:                      id,
		RendererUUID:            body.RendererUuid,
		IPAddress:               body.IpAddress,
		ProfileID:               body.ProfileId,
		DisplayName:             body.DisplayName,
		Allowed:                 body.Allowed,
		DeliveryPolicyOverrides: policy,
		Notes:                   body.Notes,
	}, nil
}
