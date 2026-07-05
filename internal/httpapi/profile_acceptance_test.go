package httpapi

import (
	"net/http"
	"testing"
)

func TestScenarioSCNAuth003UserUpdatesProfile(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-AUTH-003")

	var initial UserProfile
	client.doJSON(t, http.MethodGet, "/profile", nil, http.StatusOK, &initial)
	if initial.Username != "admin" || initial.Role != UserRoleAdmin {
		t.Fatalf("initial profile = %#v", initial)
	}

	request := UserProfileUpdateRequest{
		DisplayName: "Scenario UserRoleAdmin",
		PictureUrl:  "https://example.test/avatar.png",
	}
	var updated UserProfile
	client.doJSON(t, http.MethodPut, "/profile", request, http.StatusOK, &updated)
	if updated.DisplayName != request.DisplayName || updated.PictureUrl != request.PictureUrl {
		t.Fatalf("updated profile = %#v", updated)
	}

	var session SessionResponse
	client.doJSON(t, http.MethodGet, "/auth/session", nil, http.StatusOK, &session)
	if session.User == nil ||
		session.User.DisplayName == nil ||
		*session.User.DisplayName != request.DisplayName {
		t.Fatalf("session user was not refreshed: %#v", session.User)
	}
}
