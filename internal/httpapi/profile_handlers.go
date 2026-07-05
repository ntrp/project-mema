package httpapi

import (
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentSession, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	profile, err := s.settings.GetUserProfile(r.Context(), currentSession.userID())
	if err != nil {
		writeUserError(w, err, "Could not load profile")
		return
	}
	writeJSON(w, http.StatusOK, userProfileResponse(profile))
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	currentSession, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	var body UserProfileUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, valid := userProfileInput(w, body)
	if !valid {
		return
	}

	profile, err := s.settings.UpdateUserProfile(r.Context(), currentSession.userID(), input)
	if err != nil {
		writeUserError(w, err, "Could not update profile")
		return
	}
	writeJSON(w, http.StatusOK, userProfileResponse(profile))
}

func userProfileInput(
	w http.ResponseWriter,
	request UserProfileUpdateRequest,
) (storage.UserProfileInput, bool) {
	displayName := strings.TrimSpace(request.DisplayName)
	pictureURL := strings.TrimSpace(request.PictureUrl)
	if len(displayName) > 200 {
		writeError(w, http.StatusBadRequest, "invalid_display_name", "Name must be 200 characters or fewer")
		return storage.UserProfileInput{}, false
	}
	if len(pictureURL) > 2000 {
		writeError(w, http.StatusBadRequest, "invalid_picture_url", "Picture URL must be 2000 characters or fewer")
		return storage.UserProfileInput{}, false
	}
	return storage.UserProfileInput{DisplayName: displayName, PictureURL: pictureURL}, true
}

func userProfileResponse(profile storage.UserProfile) UserProfile {
	return UserProfile{
		Id:          openapi_types.UUID(profile.ID),
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		PictureUrl:  profile.PictureURL,
		Role:        UserRole(profile.Role),
		UpdatedAt:   profile.UpdatedAt,
	}
}

func userSummaryFromProfile(profile storage.UserProfile) UserSummary {
	return UserSummary{
		Id:          openapi_types.UUID(profile.ID),
		Username:    profile.Username,
		DisplayName: optionalString(profile.DisplayName),
		PictureUrl:  optionalString(profile.PictureURL),
		Role:        UserRole(profile.Role),
	}
}
