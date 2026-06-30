package storage

import (
	"context"
	"strconv"
	"strings"
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
	return scanMediaItemRow(s.pool.QueryRow(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		where m.monitored = true
			and m.quality_profile_id is not null
			and lower(trim(regexp_replace(m.title, '[^[:alnum:]]+', ' ', 'g'))) =
				lower(trim(regexp_replace($1, '[^[:alnum:]]+', ' ', 'g')))
			and ($2::integer is null or m.year = $2)
		order by
			case when m.year = $2 then 0 when m.year is null then 1 else 2 end,
			m.updated_at desc
		limit 1
	`, title, yearValue))
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
