package storage

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	storagegen "media-manager/internal/storage/generated"
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
	rows, err := storagegen.New(s.pool).ListLanguages(ctx)
	if err != nil {
		return nil, err
	}

	languages := make([]Language, 0, len(rows))
	for _, row := range rows {
		language, err := languageFromRow(row)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}
	return languages, nil
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
	row, err := storagegen.New(s.pool).UpsertLanguage(ctx, storagegen.UpsertLanguageParams{
		Code:        next.Code,
		DisplayName: next.DisplayName,
		Aliases:     aliases,
	})
	if err != nil {
		return Language{}, normalizeLanguageWriteError(err)
	}
	language, err := languageFromRow(row)
	return language, normalizeLanguageWriteError(err)
}

func (s *SettingsStore) DeleteLanguage(ctx context.Context, code string) error {
	rowsAffected, err := storagegen.New(s.pool).DeleteLanguage(ctx, normalizeLanguageCode(code))
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
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

func languageFromRow(row storagegen.AppLanguage) (Language, error) {
	language := Language{
		Code:        row.Code,
		DisplayName: row.DisplayName,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if len(row.Aliases) > 0 {
		if err := json.Unmarshal(row.Aliases, &language.Aliases); err != nil {
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
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23514") {
		return ErrInvalidInput
	}
	return err
}
