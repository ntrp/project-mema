package storage

import (
	"context"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserSession struct {
	ID          string
	UserID      uuid.UUID
	Username    string
	DisplayName string
	PictureURL  string
	Role        string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *SettingsStore) CreateSession(ctx context.Context, id string, userID uuid.UUID, expiresAt time.Time) error {
	return storagegen.New(s.pool).CreateSession(ctx, storagegen.CreateSessionParams{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
}

func (s *SettingsStore) GetSession(ctx context.Context, id string, now time.Time) (UserSession, error) {
	row, err := storagegen.New(s.pool).GetSession(ctx, storagegen.GetSessionParams{
		ID:        id,
		ExpiresAt: now,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return UserSession{}, ErrNotFound
	}
	return userSessionFromRow(row), err
}

func (s *SettingsStore) DeleteSession(ctx context.Context, id string) error {
	_, err := storagegen.New(s.pool).DeleteSession(ctx, id)
	return err
}

func (s *SettingsStore) DeleteExpiredSessions(ctx context.Context, now time.Time) error {
	_, err := storagegen.New(s.pool).DeleteExpiredSessions(ctx, now)
	return err
}

func userSessionFromRow(row storagegen.GetSessionRow) UserSession {
	return UserSession{
		ID:          row.ID,
		UserID:      row.UserID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		PictureURL:  row.PictureUrl,
		Role:        row.Role,
		ExpiresAt:   row.ExpiresAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
