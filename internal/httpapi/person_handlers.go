package httpapi

import (
	"net/http"
	"strings"

	"media-manager/internal/metadata"
)

func (s *Server) GetPersonDetails(w http.ResponseWriter, r *http.Request, providerType MetadataProviderType, personID string) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if !providerType.Valid() || strings.TrimSpace(personID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_person_request", "Metadata provider and person id are required")
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	provider, ok := metadataProviderByType(providers, string(providerType))
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Metadata provider is not configured")
		return
	}

	details, err := s.metadata.PersonDetails(r.Context(), metadataProviderConfig(provider), personID)
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, personDetailsResponse(providerType, details))
}

func personDetailsResponse(provider MetadataProviderType, details metadata.PersonDetails) PersonDetails {
	aliases := append([]string(nil), details.AlsoKnownAs...)
	appearances := make([]PersonAppearance, 0, len(details.Appearances))
	for _, appearance := range details.Appearances {
		appearances = append(appearances, PersonAppearance{
			Title:            appearance.Title,
			Type:             MediaType(appearance.Type),
			Year:             appearance.Year,
			ExternalProvider: MetadataProviderType(appearance.ExternalProvider),
			ExternalId:       appearance.ExternalID,
			Overview:         appearance.Overview,
			PosterPath:       appearance.PosterPath,
			BackdropPath:     appearance.BackdropPath,
			Role:             appearance.Role,
			ReleaseDate:      appearance.ReleaseDate,
		})
	}
	return PersonDetails{
		Id:           details.ID,
		Provider:     provider,
		Name:         details.Name,
		Biography:    details.Biography,
		Birthday:     details.Birthday,
		Deathday:     details.Deathday,
		PlaceOfBirth: details.PlaceOfBirth,
		ProfilePath:  details.ProfilePath,
		AlsoKnownAs:  &aliases,
		Appearances:  appearances,
	}
}
