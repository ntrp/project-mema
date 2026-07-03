package storage

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings013StorageUsersLifecycle(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-013")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	adminPassword := "secret-" + suffix

	admin, err := store.CreateUser(ctx, UserInput{
		Username: "admin-" + suffix,
		Password: &adminPassword,
		Role:     "admin",
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if !VerifyPassword(adminPassword, admin.PasswordHash) {
		t.Fatal("created admin password hash should verify")
	}

	viewerPassword := "viewer-" + suffix
	viewer, err := store.CreateUser(ctx, UserInput{
		Username: "viewer-" + suffix,
		Password: &viewerPassword,
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("create viewer: %v", err)
	}

	renamedPassword := "new-" + suffix
	updated, err := store.UpdateUser(ctx, viewer.ID, UserInput{
		Username: "viewer-updated-" + suffix,
		Password: &renamedPassword,
		Role:     "admin",
	})
	if err != nil {
		t.Fatalf("update viewer: %v", err)
	}
	if updated.Username != "viewer-updated-"+suffix || updated.Role != "admin" {
		t.Fatalf("updated user = %#v", updated)
	}
	if !VerifyPassword(renamedPassword, updated.PasswordHash) {
		t.Fatal("updated password hash should verify")
	}

	found, err := store.GetUserByUsername(ctx, " VIEWER-updated-"+suffix+" ")
	if err != nil {
		t.Fatalf("get user by username: %v", err)
	}
	if found.ID != updated.ID {
		t.Fatalf("found user id = %s, want %s", found.ID, updated.ID)
	}

	if _, err := store.CreateUser(ctx, UserInput{Username: updated.Username, Password: &viewerPassword, Role: "user"}); !errors.Is(err, ErrDuplicateUser) {
		t.Fatalf("expected duplicate user error, got %v", err)
	}
	if err := store.DeleteUser(ctx, updated.ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if _, err := store.GetUser(ctx, updated.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted user to be missing, got %v", err)
	}
}
