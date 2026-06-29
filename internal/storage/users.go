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
	_, err = s.pool.Exec(ctx, `
		insert into app.users (id, username, password_hash, role)
		select $1, $2, $3, 'admin'
		where not exists (select 1 from app.users where username = $2)
	`, id, username, hash)
	return err
}

func (s *SettingsStore) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.pool.Query(ctx, `
		select id, username, password_hash, role, created_at, updated_at
		from app.users
		order by username asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (s *SettingsStore) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := scanUser(s.pool.QueryRow(ctx, `
		select id, username, password_hash, role, created_at, updated_at
		from app.users
		where id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (s *SettingsStore) GetUserByUsername(ctx context.Context, username string) (User, error) {
	user, err := scanUser(s.pool.QueryRow(ctx, `
		select id, username, password_hash, role, created_at, updated_at
		from app.users
		where lower(username) = lower($1)
	`, strings.TrimSpace(username)))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
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
	user, err := scanUser(s.pool.QueryRow(ctx, `
		insert into app.users (id, username, password_hash, role)
		values ($1, $2, $3, $4)
		returning id, username, password_hash, role, created_at, updated_at
	`, id, input.Username, hash, input.Role))
	if isUniqueViolation(err) {
		return User{}, ErrDuplicateUser
	}
	return user, err
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

	user, err := scanUser(s.pool.QueryRow(ctx, `
		update app.users
		set username = $2,
			password_hash = coalesce($3, password_hash),
			role = $4,
			updated_at = now()
		where id = $1
		returning id, username, password_hash, role, created_at, updated_at
	`, id, input.Username, hash, input.Role))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if isUniqueViolation(err) {
		return User{}, ErrDuplicateUser
	}
	return user, err
}

func (s *SettingsStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.ensureAdminCanChange(ctx, id); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `delete from app.users where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
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
	var count int
	if err := s.pool.QueryRow(ctx, `select count(*) from app.users where role = 'admin'`).Scan(&count); err != nil {
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

func scanUser(row pgx.Row) (User, error) {
	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
