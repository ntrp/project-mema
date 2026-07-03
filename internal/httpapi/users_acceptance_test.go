package httpapi

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings024AdminManagesApplicationUsers(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-024")
	username := "scenario-user-" + uuid.NewString()[:8]
	renamed := username + "-renamed"
	updatedPassword := "updated-password-123"

	var created ManagedUser
	client.doJSON(t, http.MethodPost, "/settings/users", UserCreateRequest{
		Username: username,
		Password: "scenario-password-123",
		Role:     User,
	}, http.StatusCreated, &created)
	if created.Username != username || created.Role != User {
		t.Fatalf("created user = %#v", created)
	}

	var updated ManagedUser
	client.doJSON(t, http.MethodPut, "/settings/users/"+created.Id.String(), UserUpdateRequest{
		Username: renamed,
		Password: &updatedPassword,
		Role:     Admin,
	}, http.StatusOK, &updated)
	if updated.Username != renamed || updated.Role != Admin {
		t.Fatalf("updated user = %#v", updated)
	}

	var listed UserListResponse
	client.doJSON(t, http.MethodGet, "/settings/users", nil, http.StatusOK, &listed)
	if !managedUserListHas(listed.Users, updated.Id.String(), renamed) {
		t.Fatalf("updated user not listed: %#v", listed.Users)
	}

	client.doJSON(t, http.MethodDelete, "/settings/users/"+updated.Id.String(), nil, http.StatusNoContent, nil)
	var afterDelete UserListResponse
	client.doJSON(t, http.MethodGet, "/settings/users", nil, http.StatusOK, &afterDelete)
	if managedUserListHas(afterDelete.Users, updated.Id.String(), renamed) {
		t.Fatalf("deleted user still listed: %#v", afterDelete.Users)
	}
}

func managedUserListHas(users []ManagedUser, id string, username string) bool {
	for _, user := range users {
		if user.Id.String() == id && user.Username == username {
			return true
		}
	}
	return false
}
