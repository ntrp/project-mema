package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"crypto/pbkdf2"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	passwordHashAlgorithm = "pbkdf2-sha256"
	passwordHashKeyLength = 32
	passwordHashSaltBytes = 16
	passwordHashIter      = 210000
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	DisplayName  string
	PictureURL   string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserInput struct {
	Username string
	Password *string
	Role     string
}

func (s *SettingsStore) EnsureDefaultAdminUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	id := uuid.New()
	return storagegen.New(s.pool).EnsureDefaultAdminUser(ctx, storagegen.EnsureDefaultAdminUserParams{
		ID:           id,
		Username:     username,
		PasswordHash: hash,
	})
}

func (s *SettingsStore) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := storagegen.New(s.pool).ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, userFromRow(row))
	}
	return users, nil
}

func (s *SettingsStore) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row, err := storagegen.New(s.pool).GetUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return userFromRow(row), err
}

func (s *SettingsStore) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row, err := storagegen.New(s.pool).GetUserByUsername(ctx, strings.TrimSpace(username))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return userFromRow(row), err
}

func (s *SettingsStore) CreateUser(ctx context.Context, input UserInput) (User, error) {
	if input.Password == nil {
		return User{}, errors.New("password is required")
	}
	hash, err := HashPassword(*input.Password)
	if err != nil {
		return User{}, err
	}
	id := uuid.New()
	row, err := storagegen.New(s.pool).CreateUser(ctx, storagegen.CreateUserParams{
		ID:           id,
		Username:     input.Username,
		PasswordHash: hash,
		Role:         input.Role,
	})
	if isUniqueViolation(err) {
		return User{}, ErrDuplicateUser
	}
	return userFromRow(row), err
}

func (s *SettingsStore) UpdateUser(ctx context.Context, id uuid.UUID, input UserInput) (User, error) {
	if input.Role != "admin" {
		if err := s.ensureAdminCanChange(ctx, id); err != nil {
			return User{}, err
		}
	}

	var hash *string
	if input.Password != nil && strings.TrimSpace(*input.Password) != "" {
		value, err := HashPassword(*input.Password)
		if err != nil {
			return User{}, err
		}
		hash = &value
	}

	queries := storagegen.New(s.pool)
	var row storagegen.AppUser
	var err error
	if hash == nil {
		row, err = queries.UpdateUser(ctx, storagegen.UpdateUserParams{
			ID:       id,
			Username: input.Username,
			Role:     input.Role,
		})
	} else {
		row, err = queries.UpdateUserWithPassword(ctx, storagegen.UpdateUserWithPasswordParams{
			ID:           id,
			Username:     input.Username,
			PasswordHash: *hash,
			Role:         input.Role,
		})
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if isUniqueViolation(err) {
		return User{}, ErrDuplicateUser
	}
	return userFromRow(row), err
}

func (s *SettingsStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.ensureAdminCanChange(ctx, id); err != nil {
		return err
	}
	rows, err := storagegen.New(s.pool).DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ensureAdminCanChange(ctx context.Context, id uuid.UUID) error {
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return err
	}
	if user.Role != "admin" {
		return nil
	}
	count, err := storagegen.New(s.pool).CountAdminUsers(ctx)
	if err != nil {
		return err
	}
	if count <= 1 {
		return ErrLastAdmin
	}
	return nil
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordHashSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, passwordHashIter, passwordHashKeyLength)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s$%d$%s$%s",
		passwordHashAlgorithm,
		passwordHashIter,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func VerifyPassword(password string, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != passwordHashAlgorithm {
		return false
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil || iter <= 0 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, iter, len(expected))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(key, expected) == 1
}

func userFromRow(row storagegen.AppUser) User {
	return User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		DisplayName:  row.DisplayName,
		PictureURL:   row.PictureUrl,
		Role:         row.Role,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
