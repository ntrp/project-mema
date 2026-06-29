package storage

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MediaProfile struct {
	ID         string
	Name       string
	QualityIDs []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MediaProfileInput struct {
	Name       string
	QualityIDs []string
}

func (s *SettingsStore) EnsureDefaultMediaProfiles(ctx context.Context) error {
	for _, profile := range defaultMediaProfiles() {
		if err := s.ensureMediaProfile(ctx, profile); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingsStore) ListMediaProfiles(ctx context.Context) ([]MediaProfile, error) {
	if err := s.EnsureDefaultMediaProfiles(ctx); err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		select
			p.id,
			p.name,
			coalesce(array_agg(q.quality_id order by q.sort_order, q.quality_id) filter (where q.quality_id is not null), '{}') as quality_ids,
			p.created_at,
			p.updated_at
		from app.media_profiles p
		left join app.media_profile_qualities q on q.profile_id = p.id
		group by p.id, p.name, p.created_at, p.updated_at
		order by lower(p.name)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := []MediaProfile{}
	for rows.Next() {
		profile, err := scanMediaProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, rows.Err()
}

func (s *SettingsStore) MediaProfileExists(ctx context.Context, id string) (bool, error) {
	if err := s.EnsureDefaultMediaProfiles(ctx); err != nil {
		return false, err
	}
	var exists bool
	err := s.pool.QueryRow(ctx, `select exists(select 1 from app.media_profiles where id = $1)`, id).Scan(&exists)
	return exists, err
}

func (s *SettingsStore) CreateMediaProfile(ctx context.Context, input MediaProfileInput) (MediaProfile, error) {
	name := normalizeMediaProfileName(input.Name)
	if name == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	qualityIDs, err := normalizeProfileQualityIDs(input.QualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	id := mediaProfileSlug(name)
	if id == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	return s.saveMediaProfile(ctx, id, name, qualityIDs, true)
}

func (s *SettingsStore) UpdateMediaProfile(ctx context.Context, id string, input MediaProfileInput) (MediaProfile, error) {
	name := normalizeMediaProfileName(input.Name)
	if name == "" || strings.TrimSpace(id) == "" {
		return MediaProfile{}, ErrInvalidInput
	}
	qualityIDs, err := normalizeProfileQualityIDs(input.QualityIDs)
	if err != nil {
		return MediaProfile{}, err
	}
	return s.saveMediaProfile(ctx, strings.TrimSpace(id), name, qualityIDs, false)
}

func (s *SettingsStore) DeleteMediaProfile(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `delete from app.media_profiles where id = $1`, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) saveMediaProfile(
	ctx context.Context,
	id string,
	name string,
	qualityIDs []string,
	create bool,
) (MediaProfile, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaProfile{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if create {
		if _, err := tx.Exec(ctx, `
			insert into app.media_profiles (id, name)
			values ($1, $2)
		`, id, name); err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
	} else {
		tag, err := tx.Exec(ctx, `
			update app.media_profiles
			set name = $2, updated_at = now()
			where id = $1
		`, id, name)
		if err != nil {
			return MediaProfile{}, normalizeMediaProfileWriteError(err)
		}
		if tag.RowsAffected() == 0 {
			return MediaProfile{}, ErrNotFound
		}
	}

	if err := replaceMediaProfileQualities(ctx, tx, id, qualityIDs); err != nil {
		return MediaProfile{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaProfile{}, normalizeMediaProfileWriteError(err)
	}
	return s.getMediaProfile(ctx, id)
}

func (s *SettingsStore) getMediaProfile(ctx context.Context, id string) (MediaProfile, error) {
	profile, err := scanMediaProfile(s.pool.QueryRow(ctx, `
		select
			p.id,
			p.name,
			coalesce(array_agg(q.quality_id order by q.sort_order, q.quality_id) filter (where q.quality_id is not null), '{}') as quality_ids,
			p.created_at,
			p.updated_at
		from app.media_profiles p
		left join app.media_profile_qualities q on q.profile_id = p.id
		where p.id = $1
		group by p.id, p.name, p.created_at, p.updated_at
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaProfile{}, ErrNotFound
	}
	return profile, err
}

func (s *SettingsStore) ensureMediaProfile(ctx context.Context, profile MediaProfile) error {
	var insertedID string
	err := s.pool.QueryRow(ctx, `
		insert into app.media_profiles (id, name)
		values ($1, $2)
		on conflict do nothing
		returning id
	`, profile.ID, profile.Name).Scan(&insertedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return normalizeMediaProfileWriteError(err)
	}
	return replaceMediaProfileQualities(ctx, s.pool, insertedID, profile.QualityIDs)
}

type mediaProfileQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func replaceMediaProfileQualities(
	ctx context.Context,
	q mediaProfileQuerier,
	profileID string,
	qualityIDs []string,
) error {
	if _, err := q.Exec(ctx, `delete from app.media_profile_qualities where profile_id = $1`, profileID); err != nil {
		return err
	}
	definitions := QualitySizeDefinitionMap()
	for index, qualityID := range qualityIDs {
		definition := definitions[qualityID]
		sortOrder := int32(index)
		if definition.ID != "" {
			sortOrder = definition.SortOrder
		}
		if _, err := q.Exec(ctx, `
			insert into app.media_profile_qualities (profile_id, quality_id, sort_order)
			values ($1, $2, $3)
		`, profileID, qualityID, sortOrder); err != nil {
			return err
		}
	}
	return nil
}

func defaultMediaProfiles() []MediaProfile {
	return []MediaProfile{
		{ID: "any", Name: "Any acceptable release", QualityIDs: qualityDefinitionIDs()},
		{ID: "hd-1080p", Name: "HD 1080p", QualityIDs: []string{
			"hdtv-720p", "webdl-720p", "webrip-720p", "bluray-720p",
			"hdtv-1080p", "webdl-1080p", "webrip-1080p", "bluray-1080p",
		}},
		{ID: "uhd-4k", Name: "UHD 4K", QualityIDs: []string{
			"webdl-1080p", "webrip-1080p", "bluray-1080p", "remux-1080p",
			"hdtv-2160p", "webdl-2160p", "webrip-2160p", "bluray-2160p", "remux-2160p",
		}},
		{ID: "anime-1080p", Name: "Anime 1080p", QualityIDs: []string{
			"webdl-720p", "webrip-720p", "bluray-720p",
			"webdl-1080p", "webrip-1080p", "bluray-1080p",
		}},
	}
}

func qualityDefinitionIDs() []string {
	definitions := QualitySizeDefinitions()
	ids := make([]string, 0, len(definitions))
	for _, definition := range definitions {
		ids = append(ids, definition.ID)
	}
	return ids
}

func normalizeProfileQualityIDs(values []string) ([]string, error) {
	definitions := QualitySizeDefinitionMap()
	seen := map[string]struct{}{}
	qualityIDs := []string{}
	for _, value := range values {
		qualityID := strings.TrimSpace(value)
		if qualityID == "" {
			continue
		}
		if _, ok := definitions[qualityID]; !ok {
			return nil, ErrInvalidInput
		}
		if _, ok := seen[qualityID]; ok {
			continue
		}
		seen[qualityID] = struct{}{}
		qualityIDs = append(qualityIDs, qualityID)
	}
	return qualityIDs, nil
}

func normalizeMediaProfileName(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

var nonProfileSlugCharacter = regexp.MustCompile(`[^a-z0-9]+`)

func mediaProfileSlug(name string) string {
	slug := strings.Trim(nonProfileSlugCharacter.ReplaceAllString(strings.ToLower(name), "-"), "-")
	if len(slug) > 80 {
		slug = strings.Trim(slug[:80], "-")
	}
	return slug
}

func scanMediaProfile(row pgx.Row) (MediaProfile, error) {
	var profile MediaProfile
	err := row.Scan(&profile.ID, &profile.Name, &profile.QualityIDs, &profile.CreatedAt, &profile.UpdatedAt)
	return profile, err
}

func normalizeMediaProfileWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23503") {
		return ErrInvalidInput
	}
	return err
}
