package storage

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	profile, err := scanUserProfile(s.pool.QueryRow(ctx, `
		select id, username, display_name, picture_url, role, updated_at
		from app.users
		where id = $1
	`, id))
	if err == pgx.ErrNoRows {
		return UserProfile{}, ErrNotFound
	}
	return profile, err
}

func (s *SettingsStore) UpdateUserProfile(
	ctx context.Context,
	id uuid.UUID,
	input UserProfileInput,
) (UserProfile, error) {
	profile, err := scanUserProfile(s.pool.QueryRow(ctx, `
		update app.users
		set display_name = $2,
			picture_url = $3,
			updated_at = now()
		where id = $1
		returning id, username, display_name, picture_url, role, updated_at
	`, id, strings.TrimSpace(input.DisplayName), strings.TrimSpace(input.PictureURL)))
	if err == pgx.ErrNoRows {
		return UserProfile{}, ErrNotFound
	}
	return profile, err
}

func scanUserProfile(row pgx.Row) (UserProfile, error) {
	var profile UserProfile
	err := row.Scan(
		&profile.ID,
		&profile.Username,
		&profile.DisplayName,
		&profile.PictureURL,
		&profile.Role,
		&profile.UpdatedAt,
	)
	return profile, err
}
