package storage

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Language struct {
	Code        string
	DisplayName string
	Aliases     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LanguageInput struct {
	Code        string
	DisplayName string
	Aliases     []string
}

var languageCodePattern = regexp.MustCompile(`^[A-Z0-9-]{2,8}$`)

func (s *SettingsStore) ListLanguages(ctx context.Context) ([]Language, error) {
	rows, err := s.pool.Query(ctx, `
		select code, display_name, aliases, created_at, updated_at
		from app.languages
		order by lower(display_name), code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	languages := []Language{}
	for rows.Next() {
		language, err := scanLanguage(rows)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}
	return languages, rows.Err()
}

func (s *SettingsStore) SaveLanguage(ctx context.Context, code string, input LanguageInput) (Language, error) {
	next, err := normalizeLanguageInput(input, code == "")
	if err != nil {
		return Language{}, err
	}
	if code != "" {
		next.Code = normalizeLanguageCode(code)
		if !languageCodePattern.MatchString(next.Code) {
			return Language{}, ErrInvalidInput
		}
	}
	aliases, err := json.Marshal(next.Aliases)
	if err != nil {
		return Language{}, err
	}
	language, err := scanLanguage(s.pool.QueryRow(ctx, `
		insert into app.languages (code, display_name, aliases)
		values ($1, $2, $3)
		on conflict (code) do update
		set display_name = excluded.display_name,
			aliases = excluded.aliases,
			updated_at = now()
		returning code, display_name, aliases, created_at, updated_at
	`, next.Code, next.DisplayName, aliases))
	return language, normalizeLanguageWriteError(err)
}

func (s *SettingsStore) DeleteLanguage(ctx context.Context, code string) error {
	tag, err := s.pool.Exec(ctx, `delete from app.languages where code = $1`, normalizeLanguageCode(code))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func normalizeLanguageInput(input LanguageInput, includeCode bool) (LanguageInput, error) {
	code := normalizeLanguageCode(input.Code)
	if includeCode && !languageCodePattern.MatchString(code) {
		return LanguageInput{}, ErrInvalidInput
	}
	displayName := strings.Join(strings.Fields(input.DisplayName), " ")
	if displayName == "" {
		return LanguageInput{}, ErrInvalidInput
	}
	return LanguageInput{
		Code:        code,
		DisplayName: displayName,
		Aliases:     normalizeLanguageAliases(input.Aliases, code, displayName),
	}, nil
}

func normalizeLanguageCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeLanguageAliases(values []string, code string, displayName string) []string {
	seen := map[string]struct{}{}
	aliases := []string{}
	for _, value := range append([]string{code, displayName}, values...) {
		alias := strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
		key := strings.ToLower(alias)
		if alias == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		aliases = append(aliases, alias)
	}
	return aliases
}

func scanLanguage(row pgx.Row) (Language, error) {
	var language Language
	var aliases []byte
	err := row.Scan(
		&language.Code,
		&language.DisplayName,
		&aliases,
		&language.CreatedAt,
		&language.UpdatedAt,
	)
	if err != nil {
		return Language{}, err
	}
	if len(aliases) > 0 {
		if err := json.Unmarshal(aliases, &language.Aliases); err != nil {
			return Language{}, err
		}
	}
	if language.Aliases == nil {
		language.Aliases = []string{}
	}
	return language, nil
}

func normalizeLanguageWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrInvalidInput
	}
	return err
}
