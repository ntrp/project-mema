package storage

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

type UserProfile struct {
	ID          uuid.UUID
	Username    string
	DisplayName string
	PictureURL  string
	Role        string
	UpdatedAt   time.Time
}

type UserProfileInput struct {
	DisplayName string
	PictureURL  string
}

func (s *SettingsStore) GetUserProfile(ctx context.Context, id uuid.UUID) (UserProfile, error) {
	row, err := storagegen.New(s.pool).GetUserProfile(ctx, id)
	if err == pgx.ErrNoRows {
		return UserProfile{}, ErrNotFound
	}
	if err != nil {
		return UserProfile{}, err
	}
	return userProfileFromGetRow(row), nil
}

func (s *SettingsStore) UpdateUserProfile(
	ctx context.Context,
	id uuid.UUID,
	input UserProfileInput,
) (UserProfile, error) {
	row, err := storagegen.New(s.pool).UpdateUserProfile(ctx, storagegen.UpdateUserProfileParams{
		ID:          id,
		DisplayName: strings.TrimSpace(input.DisplayName),
		PictureUrl:  strings.TrimSpace(input.PictureURL),
	})
	if err == pgx.ErrNoRows {
		return UserProfile{}, ErrNotFound
	}
	if err != nil {
		return UserProfile{}, err
	}
	return userProfileFromUpdateRow(row), nil
}

func userProfileFromGetRow(row storagegen.GetUserProfileRow) UserProfile {
	return UserProfile{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		PictureURL:  row.PictureUrl,
		Role:        row.Role,
		UpdatedAt:   row.UpdatedAt,
	}
}

func userProfileFromUpdateRow(row storagegen.UpdateUserProfileRow) UserProfile {
	return UserProfile{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		PictureURL:  row.PictureUrl,
		Role:        row.Role,
		UpdatedAt:   row.UpdatedAt,
	}
}
