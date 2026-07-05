package storage

import (
	"context"
	"errors"
	"strconv"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) FindMonitoredMediaMatch(ctx context.Context, title string, year string) (MediaItem, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return MediaItem{}, ErrNotFound
	}
	yearValue, err := parsedMediaYear(year)
	if err != nil {
		return MediaItem{}, ErrNotFound
	}
	row, err := storagegen.New(s.pool).FindMonitoredMediaMatch(ctx, storagegen.FindMonitoredMediaMatchParams{
		Title: title,
		Year:  int4Value(yearValue),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	return mediaItemFromMatchRow(row), err
}

func parsedMediaYear(year string) (*int32, error) {
	year = strings.TrimSpace(year)
	if year == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		return nil, err
	}
	parsed := int32(value)
	return &parsed, nil
}
