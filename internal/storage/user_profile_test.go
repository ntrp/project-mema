package storage

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestUserProfileUsesGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)
	password := "profile-secret"
	user, err := store.CreateUser(ctx, UserInput{
		Username: "profile-user-" + uuid.NewString(),
		Password: &password,
		Role:     "user",
	})
	if err != nil {
		t.Fatal(err)
	}

	profile, err := store.GetUserProfile(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if profile.ID != user.ID || profile.Username != user.Username || profile.Role != "user" {
		t.Fatalf("profile = %#v, user = %#v", profile, user)
	}

	updated, err := store.UpdateUserProfile(ctx, user.ID, UserProfileInput{
		DisplayName: " Scenario User ",
		PictureURL:  " https://example.test/avatar.png ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.DisplayName != "Scenario User" || updated.PictureURL != "https://example.test/avatar.png" {
		t.Fatalf("updated profile = %#v", updated)
	}

	if _, err := store.GetUserProfile(ctx, uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing profile error = %v, want %v", err, ErrNotFound)
	}
	if _, err := store.UpdateUserProfile(ctx, uuid.New(), UserProfileInput{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing profile update error = %v, want %v", err, ErrNotFound)
	}
}
